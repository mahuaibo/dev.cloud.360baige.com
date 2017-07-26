package client

import (
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/clientselector"
	"time"
	"context"
)

/**
 * 获取服务
 */
func Call(etcdURL, serviceName, methodName string, args, reply interface{}) error {
	// RandomSelect RoundRobin WeightedRoundRobin ConsistentHash
	//NewEtcdV3ClientSelector  NewEtcdClientSelector
	s := clientselector.NewEtcdClientSelector([]string{etcdURL}, "/rpcx/"+serviceName, time.Minute, rpcx.RandomSelect, 10*time.Second)
	client := rpcx.NewClient(s)
	// Failfast Failover Failtry Broadcast Forking
	client.FailMode = rpcx.Failover
	err := client.Call(context.Background(),serviceName+"."+methodName, args, &reply)
	client.Close()
	return err
}
