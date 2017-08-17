package client

import (
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/clientselector"
	//"github.com/smallnest/rpcx/codec"
	"time"
	"context"
)

/**
 * 获取服务
 */
func Call(etcdURL, serviceName, methodName string, args, reply interface{}) error {
	// RandomSelect RoundRobin WeightedRoundRobin ConsistentHash
	//NewEtcdV3ClientSelector  NewEtcdClientSelector
	s := clientselector.NewEtcdV3ClientSelector([]string{etcdURL}, "/rpcx/"+serviceName, time.Minute, rpcx.RandomSelect, 10*time.Second)
	client := rpcx.NewClient(s)
	//client.ClientCodecFunc = codec.NewGobClientCodec

	// Failfast Failover Failtry Broadcast Forking
	client.FailMode = rpcx.Failover
	err := client.Call(context.Background(), serviceName+"."+methodName, args, &reply)
	client.Close()
	return err
}
