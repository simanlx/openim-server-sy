package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"time"
)

// 第三方支付(创建一个第三方登陆)
func InsertThirdPayOrder(tp *db.ThirdPayOrder) error {
	tp.AddTime = time.Now()
	tp.EditTime = time.Now()
	tp.LastNotifyTime = time.Now()
	tp.PayTime = time.Now()
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_order").Create(tp)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 修改第三方支付
func UpdateThirdPayOrder(tp *db.ThirdPayOrder, id int64) error {
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_order").Where("id = ?", id).Updates(tp)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 记录本次回调次数和回调时间 以及结果
func UpdateThirdPayOrderCallback(isNotify, notifyCount int, OrderID string) error {
	tp := &db.ThirdPayOrder{
		IsNotify:       int32(isNotify),
		NotifyCount:    int32(notifyCount),
		LastNotifyTime: time.Now(),
	}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_order").Where("ncount_order_no = ?", OrderID).Updates(tp)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 查询订单号： 通过商户订单号查询
func GetThirdPayMerOrderNO(MerOrderNo string) (error, *db.ThirdPayOrder) {
	resp := &db.ThirdPayOrder{}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_order").Where("mer_order_no = ?", MerOrderNo).Find(resp)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, resp
}

// 查询订单号： 通过家等你订单号查询
func GetThirdPayOrderNo(OrderNo string) (error, *db.ThirdPayOrder) {
	resp := &db.ThirdPayOrder{}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_order").Where("order_no = ?", OrderNo).Find(resp)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, resp
}

// 查询订单号： 通过新生支付的商户订单号查询
func GetThirdPayJdnMerOrderID(MerOrderID string) (error, *db.ThirdPayOrder) {
	resp := &db.ThirdPayOrder{}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_order").Where("ncount_order_no = ?", MerOrderID).Find(resp)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, resp
}

// 查询订单号： 通过新生支付的商户订单号查询
func GetThirdPayNcountMerOrderID(MerOrderID string) (error, *db.ThirdPayOrder) {
	resp := &db.ThirdPayOrder{}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_order").Where("ncount_ture_order_no = ?", MerOrderID).Find(resp)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, resp
}

// 查询一定时间内的所有订单 ，2 回调次数小于5次的,找一百条
// 20个goroutine ，平均每个请求处理时间0.5 秒， 一分钟内可以处理 20*60*2 = 2400 个请求
// 保证了第一次回调的时候，不会等待太久
// 状态为200 的订单
func GetThirdPayOrderListByTime(startTime, endTime string) ([]*db.ThirdPayOrder, error) {
	resp := []*db.ThirdPayOrder{}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_order").Where("add_time >= ? and add_time <= ? and is_notify = 0 and notify_count < 5 and status =200 ", startTime, endTime).Limit(100).Find(&resp)
	if result.Error != nil {
		return nil, result.Error
	}
	return resp, nil
}
