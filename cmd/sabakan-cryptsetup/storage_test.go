package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	backingFileSize = 3 * 1024 * 1024 // >2MiB (metadata size)
)

func createPseudoDevice(t *testing.T, dir, backing, device string, deviceMap map[string]*storageDevice) {
	f, err := os.Create(filepath.Join(dir, backing))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = f.Truncate(backingFileSize)
	if err != nil {
		t.Fatal(err)
	}

	devDir := filepath.Join(dir, "devices")
	err = os.MkdirAll(devDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	devPath := filepath.Join(devDir, device)
	err = os.Symlink(f.Name(), devPath)
	if err != nil {
		t.Fatal(err)
	}

	deviceMap[device] = &storageDevice{
		byPath:   devPath,
		realPath: f.Name(),
	}
}

func setupTestStorage(t *testing.T, dir string) (map[string]*storageDevice, *devfsType) {
	deviceMap := make(map[string]*storageDevice)

	createPseudoDevice(t, dir, "f1", "nvme-1", deviceMap)
	createPseudoDevice(t, dir, "f1", "nvme-1-dup", deviceMap)
	createPseudoDevice(t, dir, "f2", "nvme-1-part1", deviceMap)
	createPseudoDevice(t, dir, "f3", "nvme-2", deviceMap)
	createPseudoDevice(t, dir, "f4", "sata-1", deviceMap)
	createPseudoDevice(t, dir, "f5", "sata-2", deviceMap)

	testDevfs := &devfsType{path: filepath.Join(dir, "devices")}

	return deviceMap, testDevfs
}

func sameDevices(x, y []*storageDevice) bool {
	xMap := make(map[string]*storageDevice)
	for _, d := range x {
		xMap[d.realPath] = d
	}

	yMap := make(map[string]*storageDevice)
	for _, d := range y {
		yMap[d.realPath] = d
	}

	return reflect.DeepEqual(xMap, yMap)
}

func testDetectStorageDevices(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "sabakan-cryptsetup-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	deviceMap, testDevfs := setupTestStorage(t, d)

	devices, err := testDevfs.detectStorageDevices(context.Background(), []string{"nvme-*", "*-1"})
	if err != nil {
		t.Fatal(err)
	}
	expected := []*storageDevice{deviceMap["nvme-1"], deviceMap["nvme-2"], deviceMap["sata-1"]}
	if !sameDevices(devices, expected) {
		t.Error("detected wrong storage devices")
		for _, device := range devices {
			t.Log(device)
		}
	}
}

func allZero(b []byte) bool {
	for _, c := range b {
		if c != 0 {
			return false
		}
	}
	return true
}

func xorBytes(a, b []byte) []byte {
	l := len(a)
	if l > len(b) {
		l = len(b)
	}

	ret := make([]byte, l)
	for i := 0; i < l; i++ {
		ret[i] = a[i] ^ b[i]
	}

	return ret
}

func testEncrypt(t *testing.T) {
	t.Parallel()

	d, err := ioutil.TempDir("", "sabakan-cryptsetup-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	deviceMap, testDevfs := setupTestStorage(t, d)

	device := deviceMap["nvme-1"]
	rand.Seed(1)
	err = device.encrypt(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if device.id == nil {
		t.Error("device ID not set")
	}
	if allZero(device.id) {
		t.Error("device ID all zero")
	}
	if device.key == nil {
		t.Error("device key not set")
	}
	if allZero(device.key) {
		t.Error("device key all zero")
	}

	f, err := os.Open(device.byPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	readMagic := make([]byte, len(magic))
	_, err = io.ReadFull(f, readMagic)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(readMagic, magic) {
		t.Error("invalid magic:", readMagic)
	}

	readOpad := make([]byte, keyBytes)
	_, err = io.ReadFull(f, readOpad)
	if err != nil {
		t.Fatal(err)
	}
	if allZero(readOpad) {
		t.Error("opad all zero")
	}
	if allZero(xorBytes(readOpad, device.key)) {
		t.Error("encryption key all zero")
	}

	readID := make([]byte, idBytes)
	_, err = io.ReadFull(f, readID)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(readID, device.id) {
		t.Error("ID mismatch", readID, device.id)
	}

	device2 := deviceMap["nvme-2"]
	rand.Seed(1)
	err = device2.encrypt(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(device.id, device2.id) {
		t.Error("device ID not random")
	}

	f2, err := os.Open(device2.byPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	readOpad2 := make([]byte, keyBytes)
	_, err = f2.Seek(int64(len(magic)), 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.ReadFull(f2, readOpad2)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(readOpad, readOpad2) {
		t.Error("opad not random")
	}
	if bytes.Equal(xorBytes(readOpad, device.key), xorBytes(readOpad2, device2.key)) {
		t.Error("encryption key not random")
	}

	detected, err := testDevfs.detectStorageDevices(context.Background(), []string{"nvme-1"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(detected[0].id, readID) {
		t.Error("failed to detect ID")
	}
}

func TestStorage(t *testing.T) {
	t.Run("Detect", testDetectStorageDevices)
	t.Run("Encrypt", testEncrypt)
}
