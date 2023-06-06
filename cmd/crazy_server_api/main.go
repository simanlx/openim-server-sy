package main

import (
	"crazy_server/internal/api"
	"crazy_server/internal/api/filter"
	apiThird "crazy_server/internal/api/third"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	"flag"
	"fmt"

	"strconv"

	//"syscall"
	"crazy_server/pkg/common/constant"
)

// @title open-IM-Server API
// @version 1.0
// @description  open-IM-Server 的API服务器文档, 文档中所有请求都有一个operationID字段用于链路追踪

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	router := api.NewGinRouter()
	log.NewPrivateLog("im_api")
	err := filter.StartFilter()
	if err != nil {
		panic(err)
	}
	// 注册配置中心
	go getcdv3.RegisterConf()
	// 第三方
	go apiThird.MinioInit()
	defaultPorts := config.Config.Api.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10002 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", constant.CurrentVersion)
	err = router.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		panic("api start failed " + err.Error())
	}
}
