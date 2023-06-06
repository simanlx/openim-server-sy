package main

import (
	"crazy_server/internal/cms_api"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/utils"
	"flag"
	"fmt"
	"strconv"

	"crazy_server/pkg/common/config"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	log.NewPrivateLog("im_cms_api")
	log.Info("0", "启动im_cms_api 服务 ", constant.CurrentVersion)
	router := cms_api.NewGinRouter()
	router.Use(utils.CorsHandler())
	defaultPorts := config.Config.CmsApi.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10006 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.Config.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start cms api server, address: ", address, ", OpenIM version: ", constant.CurrentVersion, "\n")
	router.Run(address)
}
