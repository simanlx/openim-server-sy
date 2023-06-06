package main

import (
	rpc "crazy_server/internal/rpc/office"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	promePkg "crazy_server/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImOfficePort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.OfficePrometheusPort[0], "officePrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start office rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpc.NewOfficeServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
