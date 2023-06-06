package main

import (
	"crazy_server/internal/rpc/agent"
	"crazy_server/pkg/common/config"
	promePkg "crazy_server/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImAgentPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.AgentPrometheusPort[0], "AgentPrometheusPort default listen port")
	flag.Parse()
	rpcServer := agent.NewRpcAgentServer(*rpcPort)
	fmt.Println("start crazy_server_agent server, address: ", *rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
