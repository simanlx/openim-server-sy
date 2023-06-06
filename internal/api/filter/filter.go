package filter

import (
	api "crazy_server/pkg/base_info"
	imdb "crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	"crazy_server/pkg/common/log"
	filter "github.com/antlinker/go-dirtyfilter"
	"github.com/antlinker/go-dirtyfilter/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

//package main
//
//import (
//  "fmt"
//
//  "github.com/antlinker/go-dirtyfilter"
//  "github.com/antlinker/go-dirtyfilter/store"
//)
//
//var (
//  filterText = `我是需要过滤的内容，内容为：**文@@件，需要过滤。。。`
//)
//
//func main() {
//  memStore, err := store.NewMemoryStore(store.MemoryConfig{
//    DataSource: []string{"文件"},
//  })
//  if err != nil {
//    panic(err)
//  }
//  filterManage := filter.NewDirtyManager(memStore)
//  result, err := filterManage.Filter().Filter(filterText, '*', '@')
//  if err != nil {
//    panic(err)
//  }
//  fmt.Println(result)
//}

var memoryStore *store.MemoryStore

func StartFilter() error {
	// 需要倒入过滤器文件
	list, err := imdb.GetAllBlockword()
	if err != nil {
		return err
	}
	var filterList []string
	for _, v := range list {
		filterList = append(filterList, v.Word)
	}
	memStore, err := store.NewMemoryStore(store.MemoryConfig{
		DataSource: filterList})
	if err != nil {
		return err
	}
	memoryStore = memStore
	log.Info("filter start success")
	return nil
}

func Filter(c *gin.Context) {
	// 这里是过滤器的逻辑
	var params api.FilterContentReq
	if err := c.BindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	filterManage := filter.NewDirtyManager(memoryStore)

	result, err := filterManage.Filter().Filter(params.Content, '*', '@')
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 200, "errMsg": "验证成功", "data": result})
}
