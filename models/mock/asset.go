package mock

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/cybozu-go/sabakan"
)

type assetDriver struct {
	mu     sync.Mutex
	assets map[string]*sabakan.Asset
	data   map[string][]byte
	lastID int
}

func newAssetDriver() *assetDriver {
	return &assetDriver{
		assets: make(map[string]*sabakan.Asset),
		data:   make(map[string][]byte),
	}
}

func (d *assetDriver) GetIndex(ctx context.Context) ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	ret := make([]string, 0, len(d.assets))
	for k := range d.assets {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret, nil
}

func (d *assetDriver) GetInfo(ctx context.Context, name string) (*sabakan.Asset, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	asset, ok := d.assets[name]
	if !ok {
		return nil, sabakan.ErrNotFound
	}

	return asset, nil
}

func (d *assetDriver) Put(ctx context.Context, name, contentType string,
	r io.Reader) (*sabakan.AssetStatus, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	asset, ok := d.assets[name]
	if ok {
		return d.updateAsset(ctx, asset, contentType, r)
	}

	return d.newAsset(ctx, name, contentType, r)
}

func (d *assetDriver) Get(ctx context.Context, name string,
	f func(modtime time.Time, contentType string, content io.ReadSeeker)) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	asset, ok := d.assets[name]
	if !ok {
		return sabakan.ErrNotFound
	}

	f(asset.Date, asset.ContentType, bytes.NewReader(d.data[name]))

	return nil
}

func (d *assetDriver) Delete(ctx context.Context, name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, ok := d.assets[name]
	if !ok {
		return sabakan.ErrNotFound
	}

	delete(d.assets, name)
	delete(d.data, name)

	return nil
}

func (d *assetDriver) newAsset(ctx context.Context, name, contentType string,
	r io.Reader) (*sabakan.AssetStatus, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	d.lastID++
	id := d.lastID

	h := sha256.New()
	h.Write(data)
	sum := hex.EncodeToString(h.Sum(nil))

	asset := &sabakan.Asset{
		Name:        name,
		ID:          id,
		ContentType: contentType,
		Date:        time.Now().UTC(),
		Sha256:      sum,
		URLs:        nil,
		Version:     1,
		Exists:      true,
	}

	d.assets[name] = asset
	d.data[name] = data

	status := &sabakan.AssetStatus{
		Status:  http.StatusCreated,
		Version: 1,
		ID:      id,
	}

	return status, nil
}

func (d *assetDriver) updateAsset(ctx context.Context, asset *sabakan.Asset,
	contentType string, r io.Reader) (*sabakan.AssetStatus, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	d.lastID++
	id := d.lastID

	h := sha256.New()
	h.Write(data)
	sum := hex.EncodeToString(h.Sum(nil))

	asset.ID = id
	asset.ContentType = contentType
	asset.Date = time.Now().UTC()
	asset.Sha256 = sum
	asset.Version++

	d.data[asset.Name] = data

	status := &sabakan.AssetStatus{
		Status:  http.StatusOK,
		Version: asset.Version,
		ID:      id,
	}

	return status, nil
}
