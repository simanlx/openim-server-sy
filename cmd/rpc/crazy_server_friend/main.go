package main

import (
	"crazy_server/internal/rpc/friend"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	promePkg "crazy_server/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImFriendPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcFriendPort from cmd,default 12000 as port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.FriendPrometheusPort[0], "friendPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start friend rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := friend.NewFriendServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
