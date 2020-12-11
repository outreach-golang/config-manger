package config_manger

import (
	"errors"
	"github.com/coreos/etcd/clientv3"
	"go.uber.org/zap"
)

type Manger struct {
	config             clientv3.Config
	client             *clientv3.Client
	kv                 clientv3.KV
	watcher            clientv3.Watcher
	getResp            *clientv3.GetResponse
	watchStartRevision int64
	watchRespChan      <-chan clientv3.WatchResponse
	watchResp          clientv3.WatchResponse
	appKey             string
	Log                *zap.Logger
	err                error
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

	if manger.appKey == "" {
		return nil, errors.New("项目apps name必填")
	}

	return manger, nil
}
