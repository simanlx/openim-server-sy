package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"time"
)

// Path: pkg/common/db/mysql_model/cloud_wallet/third_withdraw.go

// 插入方法 ThirdWithdraw
func InsertThirdWithdraw(data *db.ThirdWithdraw) (err error) {
	data.AddTime = time.Now()
	data.UpdateTime = time.Now()
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_withdraw").Create(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 修改方法 ThirdWithdraw
func UpdateThirdWithdraw(data *db.ThirdWithdraw, id int64) (err error) {
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_withdraw").Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 通过第三方的订单ID查询订单记录
func GetThirdWithdrawByThirdOrderNo(orderNo string) (*db.ThirdWithdraw, error) {
	resp := &db.ThirdWithdraw{}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_withdraw").Where("third_order_id = ?", orderNo).Find(&resp)
	if result.Error != nil {
		return nil, result.Error
	}
	return resp, nil
}

// 查询方法 ThirdWithdraw
// 1. 查询一段时间内的提现记录
func GetThirdWithdrawByTime(startTime, endTime time.Time) (error, []*db.ThirdWithdraw) {
	resp := []*db.ThirdWithdraw{}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_withdraw").Where("status =? and add_time >= ? and add_time <= ?", 200, startTime, endTime).Find(&resp)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, resp
}
