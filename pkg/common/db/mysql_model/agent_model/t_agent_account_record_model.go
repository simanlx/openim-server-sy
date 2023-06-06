package agent_model

import (
	"crazy_server/pkg/common/db"
	"time"
)

const (
	AccountBusinessTypePay      = 1 //推广-下属用户充值返利
	AccountBusinessTypeShop     = 2 //销售-商城出售咖豆收入
	AccountBusinessTypeWithdraw = 3 //提现
)

type SAgentIncomeData struct {
	TodayIncome       int64 `json:"today_income"`
	AccumulatedIncome int64 `json:"accumulated_income"`
}

type AccountIncomeChartData struct {
	Date   string `json:"date"`   //日期
	Income int64  `json:"income"` //收益值
}

// 创建账户余额变更日志
func CreateAccountRecord(info *db.TAgentAccountRecord) error {
	info.CreatedTime = time.Now()
	info.UpdatedTime = time.Now()
	err := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_account_record").Create(info).Error
	return err
}

// 统计推广员收益数据
func StatAgentIncomeData(userId string) (data *SAgentIncomeData, err error) {
	today := time.Now().Format("2006-01-02")
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_account_record").
		Select("sum(amount) accumulated_income,sum(if(`day`=?,amount,0)) today_income", today).
		Where("user_id = ? and business_type in (?,?)", userId, AccountBusinessTypePay, AccountBusinessTypeShop).Scan(&data).Error
	return
}

// 账户收益图表
func AccountIncomeChart(userId string, dateType int32) (data []*AccountIncomeChartData, err error) {
	model := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_account_record").
		Where("user_id = ? and business_type in (?,?)", userId, AccountBusinessTypePay, AccountBusinessTypeShop)

	if dateType == 2 {
		month := time.Now().AddDate(0, -5, 0).Format("2006-01")
		model = model.Select("sum(amount) income,month date").Where("month >= ?", month).Group("month")
	} else {
		day := time.Now().AddDate(0, 0, -6).Format("2006-01-02")
		model = model.Select("sum(amount) income,`day` date").Where("day >= ?", day).Group("day")
	}
	err = model.Scan(&data).Error
	return
}

// 账户明细变更列表
func AccountIncomeList(userId, date string, businessType, page, size int32, chessUserIds []int64) (list []*db.TAgentAccountRecord, count int64, err error) {
	model := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_account_record").Where("user_id = ? and day = ? and status = ?", userId, date, 1)

	if businessType > 0 {
		model = model.Where("business_type = ?", businessType)
	}

	if len(chessUserIds) > 0 {
		model = model.Where("chess_user_id in (?)", chessUserIds)
	}

	err = model.Count(&count).Limit(int(size)).Offset(int(size * (page - 1))).Order("id desc").Find(&list).Error
	return
}
