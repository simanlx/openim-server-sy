package main

import (
	"crazy_server/internal/rpc/organization"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	promePkg "crazy_server/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImOrganizationPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcOrganizationPort from cmd,default 11200 as port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.OrganizationPrometheusPort[0], "organizationPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start organization rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := organization.NewServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
