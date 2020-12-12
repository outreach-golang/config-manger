package config_manger

import (
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

type Manger struct {
	config        clientv3.Config
	client        *clientv3.Client
	kv            clientv3.KV
	watcher       clientv3.Watcher
	getResp       *clientv3.GetResponse
	watchRespChan <-chan clientv3.WatchResponse
	watchResp     clientv3.WatchResponse
	app           string
	env           string
	appKey        string
}

func NewManger(endpoints []string, ops ...Option) (*Manger, error) {
	var manger *Manger
	if len(endpoints) == 0 {
		return nil, errors.New("必须传入etcd endpoints")
	}

	manger = DefaultManger(endpoints)

	for _, op := range ops {
		op(manger)
	}

	if manger.app == "" || manger.env == "" {
		return nil, errors.New("项目app, env 必填")
	}

	manger.appKey = fmt.Sprintf("/%s/%s", manger.app, manger.env)

	return manger, nil
}
