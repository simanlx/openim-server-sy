package agent_model

import (
	"crazy_server/pkg/common/db"
	"time"
)

const (
	BeanAccountBusinessTypeSale = 1 //出售
	BeanAccountBusinessTypePay  = 2 //购买
	BeanAccountBusinessTypeGive = 3 //赠送
)

// 获取推广员今日出售咖豆数
func GetAgentTodaySalesNumber(userId string) int64 {
	var info *db.TAgentBeanAccountRecord
	today := time.Now().Format("2006-01-02")
	_ = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_account_record").Select("sum(number) number").
		Where("user_id = ? and business_type = ? and `day` = ?", userId, BeanAccountBusinessTypeSale, today).Scan(&info).Error

	if info != nil {
		return info.Number
	}

	return 0
}

// 账户明细变更列表
func BeanAccountRecordList(userId, date string, businessType, page, size int32, chessUserIds []int64) (list []*db.TAgentBeanAccountRecord, count int64, err error) {
	model := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_account_record").Where("user_id = ? and day = ?", userId, date)

	if businessType > 0 {
		model = model.Where("business_type = ?", businessType)
	}

	if len(chessUserIds) > 0 {
		model = model.Where("chess_user_id in (?)", chessUserIds)
	}

	err = model.Count(&count).Limit(int(size)).Offset(int(size * (page - 1))).Order("id desc").Find(&list).Error
	return
}
