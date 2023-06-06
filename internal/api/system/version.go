package system

import (
	api "crazy_server/pkg/base_info"
	"crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	"crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// wgt版本
func WgtVersion(c *gin.Context) {
	var params api.WgtVersionReq
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	info, err := im_mysql_model.GetNewWgtVersion(params.AppId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": map[string]string{
			"app_id":  params.AppId,
			"version": "",
			"url":     "",
			"remarks": "",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": map[string]string{
		"app_id":  params.AppId,
		"version": info.Version,
		"url":     info.Url,
		"remarks": info.Remarks,
	}})
	return
}

// 家等你app最新版本
func LatestVersion(c *gin.Context) {
	params := &api.LatestVersionReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	info, err := cloud_wallet.LatestVersionByAppType(params.AppType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "app版本不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": map[string]string{
		"version": info.VersionCode,
		"url":     info.DownloadUrl,
		"content": info.UpdateContent,
	}})
	return

}
