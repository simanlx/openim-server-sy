package agent_model

import (
	"crazy_server/pkg/common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

const (
	RechargeOrderBusinessTypeWeb   = 1 //h5
	RechargeOrderBusinessTypeChess = 2 //互娱app
)

// 创建购买咖豆订单
func CreatePurchaseBeanOrder(info *db.TAgentBeanRechargeOrder) error {
	info.CreatedTime = time.Now()
	info.UpdatedTime = time.Now()
	err := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_recharge_order").Create(info).Error
	return err
}

// 获取咖豆购买订单by order_no
func GetOrderByOrderNo(orderNo string) (info *db.TAgentBeanRechargeOrder, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_recharge_order").Where("order_no = ?", orderNo).Take(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}

// 获取咖豆购买订单by chess_order_no
func GetOrderByChessOrderNo(chessOrderNo string) (info *db.TAgentBeanRechargeOrder, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_recharge_order").Where("chess_order_no = ?", chessOrderNo).Take(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}

// 支付成功更新订单状态
func PaySuccessPurchaseBeanOrderStatus(id int64, ncountOrderNo string) error {
	return db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_bean_recharge_order").Where("id = ?", id).Updates(map[string]interface{}{
		"ncount_order_no": ncountOrderNo,
		"pay_status":      1,
		"pay_time":        time.Now().Unix(),
	}).Error
}
