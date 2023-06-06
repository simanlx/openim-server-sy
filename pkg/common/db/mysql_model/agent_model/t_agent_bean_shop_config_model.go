package agent_model

import (
	"crazy_server/pkg/common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// 获取推广员自定义商城咖豆配置
func GetAgentDiyShopBeanConfig(userId string) (data []*db.TAgentBeanShopConfig, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_shop_config").Where("user_id = ?", userId).Order("id asc").Find(&data).Error
	return
}

// 获取推广员自定义商城上架咖豆配置
func GetAgentDiyShopBeanOnlineConfig(userId string) (data []*db.TAgentBeanShopConfig, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_shop_config").Where("user_id = ? and status = ?", userId, 1).Order("bean_number asc").Find(&data).Error
	return
}

// 删除咖豆配置
func DelAgentDiyShopBeanConfig(userId string) error {
	return db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_shop_config").Where("user_id = ?", userId).Delete(&db.TAgentBeanShopConfig{}).Error
}

// 修改配置上下架状态
func UpAgentDiyShopBeanConfigStatus(userId string, configIds []int32, status int32) error {
	model := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_shop_config").Where("user_id = ?", userId)

	if len(configIds) > 0 {
		model = model.Where("id in (?)", configIds)
	}

	return model.Update("status", status).Error
}

// 批量插入推广员自定义商城咖豆配置
func InsertAgentDiyShopBeanConfigs(data []*db.TAgentBeanShopConfig) (err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_shop_config").Create(&data).Error
	return
}

// 获取推广员咖豆配置
func GetAgentBeanConfigById(userId string, configId int32) (info *db.TAgentBeanShopConfig, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_shop_config").Where("id = ? and user_id = ?", configId, userId).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}
