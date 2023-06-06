package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// 获取APP的版本
func GetFVersion(versionCode string) (db.FVersion, error) {
	var fversion db.FVersion
	// 获取最新的版本
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_version").Where("version_code = ?", versionCode).First(&fversion)
	return fversion, result.Error
}

// 获取最新红包信息
func GetLastedFVersion() (db.FVersion, error) {
	var fversion db.FVersion
	// 获取最新的版本
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_version").Order("id desc").First(&fversion)
	return fversion, result.Error
}

// 获取最新版本 by appType
func LatestVersionByAppType(appType int32) (info *db.FVersion, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("f_version").Where("app_type = ? and status = ?", appType, 1).Take(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}
