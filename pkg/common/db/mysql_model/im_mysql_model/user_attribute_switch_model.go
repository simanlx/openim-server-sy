package im_mysql_model

import (
	"crazy_server/pkg/common/db"
	"time"
)

// 设置|获取用户属性开关配置
func GetUserAttributeSwitch(userId string) (*db.UserAttributeSwitch, error) {
	info := db.UserAttributeSwitch{
		UserId:                userId,
		AddFriendVerifySwitch: 1,
		AddFriendGroupSwitch:  1,
		AddFriendQrcodeSwitch: 1,
		AddFriendCardSwitch:   1,
		UpdatedTime:           time.Now(),
	}
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_attribute_switch").Where("user_id = ?", userId).FirstOrCreate(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// 设置用户属性开关配置
func SetUserAttributeSwitch(id int32, data map[string]interface{}) (err error) {
	data["updated_time"] = time.Now().Format("2006-01-02 15:04:05")
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_attribute_switch").Where("id = ?", id).Updates(data).Error
	return
}
