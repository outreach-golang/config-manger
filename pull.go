package config_manger

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"strings"
)

func (m *Manger) PullOne(name string) []byte {
	var key string
	if strings.Contains(name, "/") {
		key = fmt.Sprintf("%s/%s", name, m.env)
	} else {
		key = fmt.Sprintf("/%s/%s/%s", m.app, name, m.env)
	}
	return m.pull(key)
}

func (m *Manger) PullKvs(names []string) []byte {
	var buffers []byte
	for _, name := range names {
		buffers = append(buffers, m.PullOne(name)...)
	}
	return buffers
}

func (m *Manger) PullAppConfigs() []byte {
	return m.pull(m.appKey, clientv3.WithPrefix())
}

func (m *Manger) pull(name string, opts ...clientv3.OpOption) []byte {
	var (
		value   []byte
		getResp *clientv3.GetResponse
		err     error
	)

	if getResp, err = m.kv.Get(context.TODO(), name, opts...); err != nil {
		fmt.Println(err)
	}

	if len(getResp.Kvs) == 0 {
		fmt.Println("not found the key: ", name)
		return []byte{}
	}

	for _, kv := range getResp.Kvs {
		value = append(value, kv.Value...)
	}
	return value
}
