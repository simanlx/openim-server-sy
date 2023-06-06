package agent_model

import (
	"crazy_server/pkg/common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type BeanShopConfig struct {
	ConfigId       int32 `json:"config_id"`        //配置id
	BeanNumber     int64 `json:"bean_number"`      //购买数量
	GiveBeanNumber int32 `json:"give_bean_number"` //赠送数量
	Amount         int32 `json:"amount"`           //金额(单位元)
	Status         int32 `json:"status"`           //状态
}

// 获取配置项
func GetPlatformConfigValue(configKey string) string {
	var info *db.TAgentPlatformConfig
	err := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_platform_config").Where("config_key = ?", configKey).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ""
	}

	return info.ConfigData
}
