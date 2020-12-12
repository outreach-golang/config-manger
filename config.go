package config_manger

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type Option func(manger *Manger)

func DefaultManger(endpoints []string) *Manger {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		watcher clientv3.Watcher
		kv      clientv3.KV
		err     error
	)

	config = clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return nil
	}

	//watcher
	watcher = clientv3.NewWatcher(client)

	// KV
	kv = clientv3.NewKV(client)

	return &Manger{
		client:  client,
		kv:      kv,
		watcher: watcher,
	}
}

func SetApp(app string) Option {
	return func(m *Manger) {
		m.app = app
	}
}

func SetEnv(env string) Option {
	return func(m *Manger) {
		m.env = env
	}
}
