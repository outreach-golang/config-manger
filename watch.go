package config_manger

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
)

// watch 一个配置文件下 单个key
func (m *Manger) WatchOne(name string, changes chan<- *clientv3.Event) {
	var key = fmt.Sprintf("%s/%s", m.appKey, name)
	m.watch(key, changes)
}

//watch 一个 配置文件下多个key
func (m *Manger) WatchKvs(names []string, changes chan<- *clientv3.Event) {
	for _, name := range names {
		m.WatchOne(name, changes)
	}
}

// watch 整个配置文件变化
func (m *Manger) WatchOverAll(changes chan<- *clientv3.Event) {
	m.watch(m.appKey, changes, clientv3.WithPrefix())
}

func (m *Manger) watch(key string, changes chan<- *clientv3.Event, opts ...clientv3.OpOption) {
	var (
		watchRespChan clientv3.WatchChan
		watchResp     clientv3.WatchResponse
		event         *clientv3.Event
		//pfxWch        chan *clientv3.Event
	)
	//开始监听
	watchRespChan = m.watcher.Watch(context.TODO(), key, opts...)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			//close(pfxWch)
			wg.Done()
		}()
		for watchResp = range watchRespChan {
			for _, event = range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					changes <- event
					fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
				case mvccpb.DELETE:
					fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
				}
			}
		}
	}()
}
