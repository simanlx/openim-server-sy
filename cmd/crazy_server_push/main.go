package main

import (
	"crazy_server/internal/push/logic"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImPushPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "PushrometheusPort default listen port")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog(constant.LogFileName)
	fmt.Println("start push rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	logic.Init(*rpcPort)
	logic.Run(*prometheusPort)
	wg.Wait()
}
