package leaser

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func (svrconf *LeaserConf)NewWatcher() error {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:	svrconf.SvrHost,
		DialTimeout: time.Duration(svrconf.TimeOut)*time.Second,
	})
	if err != nil {
		return err
	}

	wch := client.Watch(context.Background(), svrconf.SvrName, clientv3.WithPrefix())
	fmt.Printf("new watcher:%v\n", svrconf.SvrName)
	for leasers := range wch {
		for _, ev := range leasers.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("new leaser:%v svrname:%s info:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			case clientv3.EventTypeDelete:
				fmt.Printf("del leaser:%v svrname:%s info:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
	return nil
}

