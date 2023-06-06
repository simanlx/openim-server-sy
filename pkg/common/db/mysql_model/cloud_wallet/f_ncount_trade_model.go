package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

const (
	TradeTypeCharge       = iota + 1 // 充值
	TradeTypeWithdraw                // 提现
	TradeTypeRedPacketOut            // 红包支出
	TradeTypeRedPacketIn             // 红包收入
	TradeTypeTransferOut             // 转账支出
	TradeTypeTransferIn              // 转账收入
	TradeTypeRefund                  // 退款
)

func FNcountTradeCreateData(req *db.FNcountTrade) error {
	req.CreatedTime = time.Now()
	req.UpdatedTime = time.Now()
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_trade").Create(req)
	if result.Error != nil {
		return errors.Wrap(result.Error, "创建交易记录失败")
	}
	return nil
}

// 修改交易的状态
func FNcountTradeUpdateStatusbyThirdOrderNo(thirdOrderNo string) error {
	// 修改红包状态
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_trade").Where("third_order_no = ?", thirdOrderNo).Updates(map[string]interface{}{"ncount_status": 1, "updated_time": time.Now()})
	if result.Error != nil {
		return errors.Wrap(result.Error, "修改交易状态失败")
	}
	return nil
}

// 根据订单号查询记录
func GetThirdOrderNoRecord(thirdOrderNo string) (info *db.FNcountTrade, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_trade").Where("third_order_no = ?", thirdOrderNo).First(&info).Error
	return
}

// 获取充值记录信息
func GetFNcountTradeByOrderNo(orderNo, userId string) (info *db.FNcountTrade, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_trade").Where("third_order_no = ? and user_id = ?", orderNo, userId).First(&info).Error
	if err != nil {
		return nil, err
	}
	return
}

// 获取账户变更列表
func FindNcountTradeList(userId, startTime, endTime string, page, size int32) (list []*db.FNcountTrade, count int64, err error) {
	model := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_trade").
		Where(" user_id = ? and ncount_status = ? and is_delete = ?", userId, 1, 0)

	if len(startTime) > 0 {
		model = model.Where("created_time >= ?", fmt.Sprintf("%s 00:00:00", startTime))
	}

	if len(endTime) > 0 {
		model = model.Where("created_time <= ?", fmt.Sprintf("%s 23:59:59", endTime))
	}

	err = model.Count(&count).Limit(int(size)).Offset(int(size * (page - 1))).Order("id desc").Find(&list).Error
	return
}

// 获取账户一段时间内的：总支出和总收入
// 条件是：ncount_status 状态为1
// is_delete 为0
// type 为 1是支出，2是收入，主要查询这两个字段
func GetNcountTradeTotal(userId, startTime, endTime string) (int64, int64, error) {
	var result struct {
		TotalIn  int64 `json:"total_in"`
		TotalOut int64 `json:"total_out"`
	}

	temres := &result
	model := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_trade").
		Where(" user_id = ? and ncount_status = ? and is_delete = ?", userId, 1, 0)

	if len(startTime) > 0 {
		model = model.Where("created_time >= ?", fmt.Sprintf("%s 00:00:00", startTime))
	}

	if len(endTime) > 0 {
		model = model.Where("created_time <= ?", fmt.Sprintf("%s 23:59:59", endTime))
	}

	err := model.Select("sum(case when type = 1 then amount else 0 end) as total_in, sum(case when type = 2 then amount else 0 end) as total_out").Scan(&temres).Error

	return temres.TotalIn, temres.TotalOut, err
}

// 软删除账户记录
func DelNcountTradeRecord(delType, recordId int32, userId string) error {
	var err error
	if delType == 0 {
		err = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_trade").Where(" id = ? and user_id = ? ", recordId, userId).Update("is_delete", 1).Error
	} else {
		//删除全部记录
		err = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_trade").Where("user_id = ? ", userId).Update("is_delete", 1).Error
	}
	return err
}
