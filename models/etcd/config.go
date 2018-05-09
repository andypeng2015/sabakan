package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"path"

	"github.com/coreos/etcd/clientv3"
	"github.com/cybozu-go/sabakan"
)

// PutConfig implements sabakan.ConfigModel
func (d *Driver) PutConfig(ctx context.Context, config *sabakan.IPAMConfig) error {
	key := path.Join(d.prefix, KeyMachines)
	resp, err := d.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if resp.Count != 0 {
		return errors.New("machine already exists")
	}

	j, err := json.Marshal(config)
	if err != nil {
		return err
	}

	key = path.Join(d.prefix, KeyConfig)
	_, err = d.client.Put(ctx, key, string(j))
	return err
}

// GetConfig implements sabakan.ConfigModel
func (d *Driver) GetConfig(ctx context.Context) (*sabakan.IPAMConfig, error) {
	key := path.Join(d.prefix, KeyConfig)
	resp, err := d.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	var config sabakan.IPAMConfig
	err = json.Unmarshal(resp.Kvs[0].Value, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}