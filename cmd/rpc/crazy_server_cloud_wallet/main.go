package main

import (
	"crazy_server/internal/rpc/cloud_wallet"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/log"
	promePkg "crazy_server/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImCloudWalletPort
	log.NewPrivateLog("crazy_server_cloud_wallet")
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.CloudWalletPrometheusPort[0], "CloudWalletPrometheusPort default listen port")
	flag.Parse()
	rpcServer := cloud_wallet.NewRpcCloudWalletServer(*rpcPort)
	fmt.Println("start crazy_server_cloud_wallet server, address: ", *rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
