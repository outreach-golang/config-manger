package config_manger

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

func (m *Manger) PullConfig(name string) []byte {
	var key = fmt.Sprintf("%s/%s", m.appKey, name)
	return m.pull(key)
}

func (m *Manger) PullAppConfigs() []byte {
	return m.pull(m.appKey)
}

func (m *Manger) pull(name string, opts ...clientv3.OpOption) []byte {
	var (
		getResp *clientv3.GetResponse
		err     error
	)
	if getResp, err = m.kv.Get(context.TODO(), name, opts...); err != nil {
		fmt.Println(err)
	}

	if len(getResp.Kvs) == 0 {
		fmt.Println("key: %s not found in etcd", name)
	}
	return getResp.Kvs[0].Value
}
