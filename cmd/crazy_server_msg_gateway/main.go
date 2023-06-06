package main

import (
	"crazy_server/internal/msg_gateway/gate"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	defaultRpcPorts := config.Config.RpcPort.OpenImMessageGatewayPort
	defaultWsPorts := config.Config.LongConnSvr.WebsocketPort
	defaultPromePorts := config.Config.Prometheus.MessageGatewayPrometheusPort
	rpcPort := flag.Int("rpc_port", defaultRpcPorts[0], "rpc listening port")
	wsPort := flag.Int("ws_port", defaultWsPorts[0], "ws listening port")
	prometheusPort := flag.Int("prometheus_port", defaultPromePorts[0], "PushrometheusPort default listen port")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("start rpc/msg_gateway server, port: ", *rpcPort, *wsPort, *prometheusPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	gate.Init(*rpcPort, *wsPort)
	gate.Run(*prometheusPort)
	wg.Wait()
}
