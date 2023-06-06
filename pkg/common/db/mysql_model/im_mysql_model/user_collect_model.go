package im_mysql_model

import (
	"crazy_server/pkg/common/db"
	"fmt"
	"time"
)

// 收藏数据入库
func InsertUserCollect(info *db.UserCollect) (err error) {
	info.CreatedTime = time.Now()
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_collect").Create(info).Error
	return
}

// 删除收藏数据
func DelUserCollect(collectId int32, userId string) (err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_collect").Where("id = ? and user_id = ?", collectId, userId).Delete(&db.UserCollect{}).Error
	return
}

// 获取收藏数据列表
func FindUserCollectList(userId, keyword string, msgType, page, size int32) (list []*db.UserCollect, count int64, err error) {
	model := db.DB.MysqlDB.DefaultGormDB().Table("user_collect").Where(" user_id = ? ", userId)

	if msgType > 0 {
		model = model.Where("collect_type = ?", msgType)
	}

	if len(keyword) > 0 {
		model = model.Where("collect_content like ?", fmt.Sprintf("%%%s%%", keyword))
	}

	err = model.Count(&count).Limit(int(size)).Offset(int(size * (page - 1))).Order("id desc").Find(&list).Error
	return
}
