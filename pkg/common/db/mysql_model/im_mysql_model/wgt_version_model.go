package im_mysql_model

import (
	"crazy_server/pkg/common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// 获取最新wgt版本
func GetNewWgtVersion(appId string) (info *db.AppWgtVersion, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("app_wgt_version").Where("app_id = ? and status = 1", appId).Order("version desc").Take(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}
