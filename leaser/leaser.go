package leaser

import (
	"github.com/coreos/etcd/clientv3"
)

const (
	CONST_HOST_LEN = 10
	CONST_TIMEOUT = 5

)

type LeaserConf struct {

	TimeOut  int64
	Stop     chan bool
	SvrHost  []string
	SvrName  string
	SvrInfo  string

	LeaseId  clientv3.LeaseID
	Client   *clientv3.Client
	LeaseCh  <-chan *clientv3.LeaseKeepAliveResponse
}
