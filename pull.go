package config_manger

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	mvccpb "github.com/coreos/etcd/mvcc/mvccpb"
)

func (m *Manger) PullOne(name string, configs map[string]string) {
	var key = fmt.Sprintf("%s/%s", m.appKey, name)
	m.pull(key, configs)
}

func (m *Manger) PullKvs(names []string, configs map[string]string) {
	for _, name := range names {
		m.PullOne(name, configs)
	}
}

func (m *Manger) PullOverAll(configs map[string]string) {
	m.pull(m.appKey, configs, clientv3.WithPrefix())
}

func (m *Manger) pull(name string, configs map[string]string, opts ...clientv3.OpOption) {
	var (
		getResp *clientv3.GetResponse
		kvs     *mvccpb.KeyValue
		err     error
	)
	if getResp, err = m.kv.Get(context.TODO(), name, opts...); err != nil {
		fmt.Println(err)
	}

	for _, kvs = range getResp.Kvs {
		key := string(kvs.Key)[len(m.appKey)+1:]
		configs[key] = string(kvs.Value)
	}
}
