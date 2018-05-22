package etcd

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/cybozu-go/sabakan"
)

var testIPAMConfig = sabakan.IPAMConfig{
	MaxNodesInRack:  28,
	NodeIPv4Pool:    "10.69.0.0/20",
	NodeRangeSize:   6,
	NodeRangeMask:   26,
	NodeIndexOffset: 3,
	NodeIPPerNode:   3,
	BMCIPv4Pool:     "10.72.16.0/20",
	BMCRangeSize:    5,
	BMCRangeMask:    20,
}

func testIPAMPutConfig(t *testing.T) {
	d, ch := testNewDriver(t)
	config := &testIPAMConfig
	err := d.putIPAMConfig(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	<-ch

	resp, err := d.client.Get(context.Background(), t.Name()+KeyIPAM)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Kvs) != 1 {
		t.Error("config was not saved")
	}
	var actual sabakan.IPAMConfig
	err = json.Unmarshal(resp.Kvs[0].Value, &actual)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, &actual) {
		t.Errorf("unexpected saved config %#v", actual)
	}

	err = d.Register(context.Background(), []*sabakan.Machine{{Serial: "1234abcd", Role: "worker"}})
	if err != nil {
		t.Fatal(err)
	}
	err = d.putIPAMConfig(context.Background(), config)
	if err == nil {
		t.Error("should be failed, because some machine is already registered")
	}
}

func testIPAMGetConfig(t *testing.T) {
	d, ch := testNewDriver(t)

	config := &testIPAMConfig

	bytes, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.client.Put(context.Background(), t.Name()+KeyIPAM, string(bytes))
	if err != nil {
		t.Fatal(err)
	}
	<-ch

	actual, err := d.getIPAMConfig()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, actual) {
		t.Errorf("unexpected loaded config %#v", actual)
	}
}

func TestIPAM(t *testing.T) {
	t.Run("Put", testIPAMPutConfig)
	t.Run("Get", testIPAMGetConfig)
}
