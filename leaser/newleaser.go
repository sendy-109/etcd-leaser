package leaser

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func (svrconf *LeaserConf)checkParam() error {

	if 	len(svrconf.SvrName) == 0 ||
		len(svrconf.SvrInfo) == 0 ||
		len(svrconf.SvrHost) == 0 {
		return errors.New("inval paramer")
	}

	for _,v := range svrconf.SvrHost {
		if len(v) < CONST_HOST_LEN {
			return errors.New("inval svrver host")
		}
	}

	if svrconf.TimeOut <= 0 {
		svrconf.TimeOut = CONST_TIMEOUT
	}

	svrconf.Stop = make(chan bool, 1)

	return nil
}

func (svrconf *LeaserConf)NewLeaser() error {

	//参数检查

	err := svrconf.checkParam()
	if err != nil {
		return err
	}

	err = svrconf.addLeaser(svrconf.SvrName, svrconf.SvrInfo)
	if err != nil {
		return err
	}

	return nil
}

func (svrconf *LeaserConf)addLeaser(key,value string )error{

	client, err := clientv3.New(clientv3.Config{
		Endpoints:	svrconf.SvrHost,
		DialTimeout: time.Duration(svrconf.TimeOut)*time.Second,
	})
	if err != nil {
		return err
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(svrconf.TimeOut)*time.Second)
	//defer func() {
	//	cancel()
	//}()
	lease, err := client.Grant(context.TODO(), int64(svrconf.TimeOut))
	if err != nil {
		return err
	}

	_, err = client.Put(context.TODO(), key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	svrconf.LeaseCh, err = client.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		return err
	}

	svrconf.Client = client
	svrconf.LeaseId = lease.ID

	err = svrconf.LeaserLoop()
	client.Close()
	return err

}

func (svrconf *LeaserConf)LeaserLoop() error{

	for {
		select {
		case <- svrconf.Stop:
			fmt.Println("stop lease")
			svrconf.Client.Revoke(context.TODO(), svrconf.LeaseId)
			return nil
		case <- svrconf.Client.Ctx().Done():
			return errors.New("server closed")
		case msg,ok := <- svrconf.LeaseCh:
			if !ok {
				if msg != nil {
					fmt.Println("msg ", msg.String())
				}
				fmt.Println("租约终止")
				svrconf.Client.Revoke(context.TODO(), svrconf.LeaseId)
				return errors.New("lease close")
			}
			if msg != nil {
				fmt.Printf("续约:%d, leadid:%d\n",msg.TTL, svrconf.LeaseId)
			}
		}
	}
	return nil
}