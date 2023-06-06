package agent

import (
	"context"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/agent_model"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/common/utils"
	"crazy_server/pkg/proto/agent"
	utils2 "crazy_server/pkg/utils"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
)

// 推广员下属成员购买咖豆(推广员商城) - 互娱回调
func (rpc *AgentServer) ChessPurchaseBeanNotify(ctx context.Context, req *agent.ChessPurchaseBeanNotifyReq) (*agent.ChessPurchaseBeanNotifyResp, error) {
	resp := &agent.ChessPurchaseBeanNotifyResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}
	log.Info("", fmt.Sprintf("start 推广员下属成员购买咖豆(推广员商城) - 互娱回调, 订单号(%s),新生支付订单号(%s),", req.OrderNo, req.NcountOrderNo), utils2.JsonFormat(req))

	// 加锁
	lockKey := fmt.Sprintf("ChessPurchaseBeanNotify:%s", req.OrderNo)
	if err := utils.Lock(ctx, lockKey); err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "操作加锁失败"
		return resp, nil
	}
	defer utils.UnLock(ctx, lockKey)

	//校验订单号
	orderInfo, err := imdb.GetOrderByOrderNo(req.OrderNo)
	if err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "订单不存在"
		return resp, nil
	}

	//校验订单状态、已处理
	if orderInfo.PayStatus == 1 {
		return resp, nil
	}

	//处理下属成员购买咖豆后逻辑
	err = handelChessPurchaseBeanLogic(orderInfo, req.NcountOrderNo)
	if err != nil {
		log.Error("", fmt.Sprintf("处理下属成员购买咖豆后逻辑失败,订单号(%s),err:%s", req.OrderNo, err.Error()))
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = err.Error()
		return resp, nil
	}

	//解冻推广员咖豆
	_ = rocksCache.DelAgentBeanBalance(ctx, orderInfo.UserId, orderInfo.ChessUserId)
	return resp, nil
}

// 处理下属成员购买商城咖豆逻辑
func handelChessPurchaseBeanLogic(info *db.TAgentBeanRechargeOrder, ncountOrderNo string) error {
	//开启事务
	tx := db.DB.AgentMysqlDB.DefaultGormDB().Begin()

	//1、更新订单状态
	err := tx.Table("t_agent_bean_recharge_order").Where("id = ?", info.Id).Updates(map[string]interface{}{
		"ncount_order_no": ncountOrderNo,
		"pay_status":      1,
		"pay_time":        time.Now().Unix(),
	}).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新订单状态失败")
	}

	//2、扣除推广员咖豆数、增加账户余额
	err = tx.Table("t_agent_account").Where("user_id = ? and bean_balance >= ?", info.UserId, info.Number).UpdateColumns(map[string]interface{}{
		"bean_balance": gorm.Expr(" bean_balance - ? ", info.Number+int64(info.GiveNumber)), // 减去(购买数 + 赠送数)
		"balance":      gorm.Expr(" balance + ? ", info.Amount),
	}).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "扣除推广员咖豆数、增加账户余额失败")
	}

	//3、增加余额变更日志
	balanceRecord := &db.TAgentAccountRecord{
		OrderNo:           info.OrderNo,
		UserId:            info.UserId,
		Type:              1,
		BusinessType:      imdb.AccountBusinessTypeShop,
		ChessUserId:       info.ChessUserId,
		ChessUserNickname: info.ChessUserNickname,
		Describe:          fmt.Sprintf("%sID%d 购买%d咖豆", info.ChessUserNickname, info.ChessUserId, info.Number),
		Amount:            info.Amount,
		Status:            1,
		Day:               time.Now().Format("2006-01-02"),
		Month:             time.Now().Format("2006-01"),
		CreatedTime:       time.Now(),
		UpdatedTime:       time.Now(),
		DB:                tx,
	}
	err = tx.Table("t_agent_account_record").Create(&balanceRecord).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "增加余额变更日志失败")
	}

	//4、增加咖豆变更日志
	beanRecord := &db.TAgentBeanAccountRecord{
		OrderNo:           info.OrderNo,
		UserId:            info.UserId,
		Type:              1,
		BusinessType:      imdb.BeanAccountBusinessTypeSale,
		ChessUserId:       info.ChessUserId,
		ChessUserNickname: info.ChessUserNickname,
		Describe:          fmt.Sprintf("%sID%d 购买%d咖豆", info.ChessUserNickname, info.ChessUserId, info.Number),
		Amount:            info.Amount,
		Number:            info.Number,
		GiveNumber:        info.GiveNumber,
		Day:               time.Now().Format("2006-01-02"),
		CreatedTime:       time.Now(),
		UpdatedTime:       time.Now(),
		DB:                tx,
	}
	err = tx.Table("t_agent_bean_account_record").Create(&beanRecord).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "增加咖豆变更日志失败")
	}

	//5、更新成员贡献度
	err = tx.Table("t_agent_member").Where("user_id = ? and chess_user_id = ?", info.UserId, info.ChessUserId).UpdateColumn("contribution", gorm.Expr(" contribution + ? ", info.Amount)).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新成员贡献度失败")
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		log.NewError("", "handelMemberPurchaseBeanLogic commit error:", err, "tx.Rollback().Error :", tx.Rollback().Error)
		return errors.Wrap(err, "事务提交失败")
	}
	return nil
}

// 推广员成员购买咖豆回调(平台商城) - 互娱回调
func (rpc *AgentServer) PlatformPurchaseBeanNotify(ctx context.Context, req *agent.PlatformPurchaseBeanNotifyReq) (*agent.PlatformPurchaseBeanNotifyResp, error) {
	resp := &agent.PlatformPurchaseBeanNotifyResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}
	log.Info("", fmt.Sprintf("start 推广员成员购买咖豆回调(平台商城) - 互娱回调, 互娱订单号(%s),新生支付订单号(%s),", req.ChessOrderNo, req.NcountOrderNo), utils2.JsonFormat(req))

	// 加锁
	lockKey := fmt.Sprintf("PlatformPurchaseBeanNotify:%s", req.ChessOrderNo)
	if err := utils.Lock(ctx, lockKey); err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "操作加锁失败"
		return resp, nil
	}
	defer utils.UnLock(ctx, lockKey)

	//是否为下属成员
	agentMember, err := imdb.AgentNumberByChessUserId(req.ChessUserId)
	if err != nil || req.AgentNumber != agentMember.AgentNumber {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "该用户不是推广员编号下成员"
		return resp, nil
	}

	//校验订单号、已处理
	_, err = imdb.GetOrderByChessOrderNo(req.ChessOrderNo)
	if err == nil {
		return resp, nil
	}

	//计算平台充值返利
	contribution := computeRechargeRebate(req.Amount)
	if contribution == 0 {
		return resp, nil
	}

	orderNo := utils.GetOrderNo() //平台订单号
	orderInfo := &db.TAgentBeanRechargeOrder{
		BusinessType:      imdb.RechargeOrderBusinessTypeChess,
		UserId:            agentMember.UserId,
		ChessUserId:       req.ChessUserId,
		ChessUserNickname: agentMember.ChessNickname,
		OrderNo:           orderNo,
		ChessOrderNo:      req.ChessOrderNo,
		NcountOrderNo:     req.NcountOrderNo,
		Number:            req.BeanNumber,
		GiveNumber:        req.GiveBeanNumber,
		Amount:            req.Amount,
		PayTime:           time.Now().Unix(),
		PayStatus:         1,
		CreatedTime:       time.Now(),
		UpdatedTime:       time.Now(),
	}

	//处理下属成员购买平台咖豆后逻辑
	err = handelPlatformPurchaseBeanLogic(orderInfo, agentMember.Id, contribution)
	if err != nil {
		log.Error("", fmt.Sprintf("处理下属成员购买咖豆后逻辑失败,互娱订单号(%s),err:%s", req.ChessOrderNo, err.Error()))
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = err.Error()
		return resp, nil
	}

	//解冻推广员咖豆
	_ = rocksCache.DelAgentBeanBalance(ctx, orderInfo.UserId, orderInfo.ChessUserId)
	return resp, nil
}

// 处理下属成员购买平台咖豆逻辑
func handelPlatformPurchaseBeanLogic(info *db.TAgentBeanRechargeOrder, agentMemberId, contribution int32) error {
	//开启事务
	tx := db.DB.AgentMysqlDB.DefaultGormDB().Begin()

	//1、生成订单
	err := tx.Table("t_agent_bean_recharge_order").Create(info).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "生成订单失败")
	}

	//2、充值返利、增加账户余额
	err = tx.Table("t_agent_account").Where("user_id = ? ", info.UserId).UpdateColumn("balance", gorm.Expr(" balance + ? ", contribution)).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("充值返利(%d)、增加账户余额失败", contribution))
	}

	//3、增加余额变更日志
	balanceRecord := &db.TAgentAccountRecord{
		OrderNo:           info.OrderNo,
		UserId:            info.UserId,
		Type:              1,
		BusinessType:      imdb.AccountBusinessTypePay,
		ChessUserId:       info.ChessUserId,
		ChessUserNickname: info.ChessUserNickname,
		Describe:          fmt.Sprintf("%sID%d 购买%d咖豆", info.ChessUserNickname, info.ChessUserId, info.Number),
		Amount:            contribution,
		Status:            1,
		Day:               time.Now().Format("2006-01-02"),
		Month:             time.Now().Format("2006-01"),
		CreatedTime:       time.Now(),
		UpdatedTime:       time.Now(),
		DB:                tx,
	}
	err = tx.Table("t_agent_account_record").Create(&balanceRecord).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "增加余额变更日志失败")
	}

	//4、更新成员贡献度
	err = tx.Table("t_agent_member").Where("id = ?", agentMemberId).UpdateColumn("contribution", gorm.Expr(" contribution + ? ", contribution)).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新成员贡献度失败")
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		log.NewError("", "handelPlatformPurchaseBeanLogic commit error:", err, "tx.Rollback().Error :", tx.Rollback().Error)
		return errors.Wrap(err, "事务提交失败")
	}
	return nil
}

// 计算平台充值返利
func computeRechargeRebate(amount int32) int32 {
	if amount == 0 {
		return 0
	}

	//充值返现比例、千分比
	rechargeRebate := rocksCache.GetPlatformValueConfigCache("recharge_rebate")
	rechargeRebateDecimal, _ := decimal.NewFromString(rechargeRebate)
	if rechargeRebateDecimal.IsZero() {
		return 0
	}

	// (支付金额 * 返现比例) / 千分比
	amountDecimal := decimal.NewFromInt32(amount)
	rebateAmount, _ := amountDecimal.Mul(rechargeRebateDecimal).Div(decimal.NewFromInt(1000)).RoundFloor(0).Float64() //(a * b) / 1000 取整
	log.Info("", "计算平台充值返利 computeRechargeRebate,充值金额,返利比例,返利金额 ", amount, rechargeRebate, rebateAmount)

	return cast.ToInt32(rebateAmount)
}

// 推广员充值咖豆 - 新生支付回调
func (rpc *AgentServer) RechargeNotify(ctx context.Context, req *agent.RechargeNotifyReq) (*agent.RechargeNotifyResp, error) {
	resp := &agent.RechargeNotifyResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}
	log.Info("", fmt.Sprintf("推广员充值咖豆 - 新生支付回调, 订单号(%s),新生支付订单号(%s),", req.OrderNo, req.NcountOrderNo), utils2.JsonFormat(req))

	// 加锁
	lockKey := fmt.Sprintf("RechargeNotify:%s", req.OrderNo)
	if err := utils.Lock(ctx, lockKey); err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "操作加锁失败"
		return resp, nil
	}
	defer utils.UnLock(ctx, lockKey)

	//校验订单号
	orderInfo, err := imdb.GetOrderByOrderNo(req.OrderNo)
	if err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "订单不存在"
		return resp, nil
	}

	//校验订单状态、已处理
	if orderInfo.PayStatus == 1 {
		return resp, nil
	}

	//校验金额
	if orderInfo.Amount > req.Amount {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = fmt.Sprintf("订单支付金额错误,订单金额(%d),支付金额(%d)", orderInfo.Amount, req.Amount)
		log.Error("", fmt.Sprintf("推广员充值咖豆-新生支付回调订单号:(%s),err:%s", req.OrderNo, resp.CommonResp.Msg))
		return resp, nil
	}

	//处理充值咖豆逻辑：1、改订单状态 、2给推广员增加咖豆、3记录咖豆账户变更日志
	err = handelRechargeNotifyLogic(orderInfo, req.NcountOrderNo, req.PayTime)
	if err != nil {
		log.Error("", fmt.Sprintf("推广员充值咖豆-新生支付回调订单号:(%s),err:%s", req.OrderNo, err.Error()))
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}

// 处理充值咖豆逻辑：1、改订单状态 、2给推广员增加咖豆、3记录咖豆账户变更日志
func handelRechargeNotifyLogic(info *db.TAgentBeanRechargeOrder, ncountOrderNo, payTime string) error {
	//开启事务
	tx := db.DB.AgentMysqlDB.DefaultGormDB().Begin()

	//1、更新订单状态
	var payT int64
	payTimeInt, err := time.Parse("2006-01-02 15:04:05", payTime)
	if err == nil {
		payT = payTimeInt.Unix()
	}

	err = tx.Table("t_agent_bean_recharge_order").Where("id = ?", info.Id).Updates(map[string]interface{}{
		"ncount_order_no": ncountOrderNo,
		"pay_status":      1,
		"pay_time":        payT,
	}).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新订单状态失败")
	}

	//2、推广员增加咖豆
	err = tx.Table("t_agent_account").Where("user_id = ? ", info.UserId).UpdateColumn("bean_balance", gorm.Expr(" bean_balance + ? ", info.Number+int64(info.GiveNumber))).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("给推广员(%s)增加咖豆(%d)失败：%s", info.UserId, info.Number, err.Error()))
	}

	//3、记录咖豆账户变更日志
	beanRecord := &db.TAgentBeanAccountRecord{
		OrderNo:           info.OrderNo,
		UserId:            info.UserId,
		Type:              2,
		BusinessType:      imdb.BeanAccountBusinessTypePay,
		ChessUserId:       info.ChessUserId,
		ChessUserNickname: info.ChessUserNickname,
		Describe:          fmt.Sprintf("购买%d咖豆", info.Number),
		Amount:            info.Amount,
		Number:            info.Number,
		GiveNumber:        info.GiveNumber,
		Day:               time.Now().Format("2006-01-02"),
		CreatedTime:       time.Now(),
		UpdatedTime:       time.Now(),
		DB:                tx,
	}
	err = tx.Table("t_agent_bean_account_record").Create(&beanRecord).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "增加咖豆变更日志失败")
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		log.NewError("", "handelRechargeNotifyLogic commit error:", err, "tx.Rollback().Error :", tx.Rollback().Error)
		return errors.Wrap(err, "事务提交失败")
	}
	return nil
}
