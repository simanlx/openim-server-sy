package main

import (
	"crazy_server/internal/rpc/user"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	promePkg "crazy_server/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImUserPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.UserPrometheusPort[0], "userPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start user rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion)
	rpcServer := user.NewUserServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
