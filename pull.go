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
	var buffers = make([]byte, 0)
	for _, name := range names {
		buffers = append(buffers, m.PullOne(name)...)
	}
	return buffers
}

func (m *Manger) PullAppConfigs() []byte {
	return m.pull(m.appKey)
}

func (m *Manger) pull(name string, opts ...clientv3.OpOption) []byte {
	var (
		getResp *clientv3.GetResponse
		err     error
	)
	fmt.Println("----", name)
	if getResp, err = m.kv.Get(context.TODO(), name, opts...); err != nil {
		fmt.Println(err)
	}

	if len(getResp.Kvs) == 0 {
		fmt.Println("not found the key: ", name)
		return []byte{}
	}
	return getResp.Kvs[0].Value
}
