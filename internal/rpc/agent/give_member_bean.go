package agent

import (
	"context"
	chessApi "crazy_server/pkg/agent"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/agent_model"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/common/utils"
	"crazy_server/pkg/proto/agent"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// 赠送下属成员咖豆
func (rpc *AgentServer) AgentGiveMemberBean(ctx context.Context, req *agent.AgentGiveMemberBeanReq) (*agent.AgentGiveMemberBeanResp, error) {
	resp := &agent.AgentGiveMemberBeanResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}
	log.Info(req.OperationId, fmt.Sprintf("start 推广员(%s),赠送下属成员(%d)咖豆,赠送咖豆数(%d)", req.UserId, req.ChessUserId, req.BeanNumber))

	// 加锁
	lockKey := fmt.Sprintf("AgentGiveMemberBean:%d", req.ChessUserId)
	if err := utils.Lock(ctx, lockKey); err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "操作加锁失败"
		return resp, nil
	}
	defer utils.UnLock(ctx, lockKey)

	//获取推广员信息
	info, err := imdb.GetAgentByUserId(req.UserId)
	if err != nil || info.OpenStatus == 0 {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "推广员信息有误"
		return resp, nil
	}

	//是否为推广员下属成员
	agentMember, err := imdb.AgentNumberByChessUserId(req.ChessUserId)
	if err != nil || agentMember.UserId != req.UserId {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "不是推广员下属成员"
		return resp, nil
	}

	//冻结的咖豆额度
	freezeBeanBalance := rocksCache.GetAgentFreezeBeanBalance(ctx, info.UserId)

	//校验推广员咖豆余额
	if info.BeanBalance < (req.BeanNumber + freezeBeanBalance) {
		log.Error(req.OperationId, fmt.Sprintf("推广员(%d),赠送下属成员(%d)咖豆,咖豆余额不足,咖豆余额(%d),赠送咖豆(%d),冻结咖豆(%d)", info.AgentNumber, req.ChessUserId, info.BeanBalance, req.BeanNumber, freezeBeanBalance))
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "咖豆余额不足"
		return resp, nil
	}

	// 赠送新互娱用户咖豆
	err = GiveChessUserBean(req.UserId, info.AgentNumber, info.Id, req.BeanNumber, req.ChessUserId, agentMember.ChessNickname)
	if err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = err.Error()
		log.Error(req.OperationId, fmt.Sprintf("推广员(%d),赠送下属成员(%d)咖豆,操作失败:%s", info.AgentNumber, req.ChessUserId, err.Error()))
		return resp, nil
	}

	return resp, nil
}

// 赠送互娱用户咖豆
func GiveChessUserBean(userId string, agentNumber int32, agentId, beanNumber, chessUserId int64, chessNickname string) error {
	//开启事务
	tx := db.DB.AgentMysqlDB.DefaultGormDB().Begin()

	//1、扣除推广员咖豆
	err := tx.Table("t_agent_account").Where("id = ? and bean_balance >= ?", agentId, beanNumber).UpdateColumn("bean_balance", gorm.Expr(" bean_balance - ? ", beanNumber)).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "扣除推广员咖豆失败")
	}

	//2、增加咖豆变更日志
	orderNo := utils.GetOrderNo()
	record := &db.TAgentBeanAccountRecord{
		OrderNo:           orderNo,
		UserId:            userId,
		Type:              2,
		BusinessType:      imdb.BeanAccountBusinessTypeGive,
		ChessUserId:       chessUserId,
		ChessUserNickname: chessNickname,
		Describe:          fmt.Sprintf("赠送给%sID%d %d咖豆", chessNickname, chessUserId, beanNumber),
		Amount:            0,
		Number:            beanNumber,
		GiveNumber:        0,
		Day:               time.Now().Format("2006-01-02"),
		CreatedTime:       time.Now(),
		UpdatedTime:       time.Now(),
		DB:                tx,
	}
	err = tx.Table("t_agent_bean_account_record").Create(&record).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "增加咖豆变更日志失败")
	}

	//3、调用chess api 给用户加咖豆、企业微信告警
	if !GiveUserBeanRetry(orderNo, agentNumber, chessUserId, beanNumber, []int{0, 3}) {
		tx.Rollback()
		return errors.New("调用chess api 给用户加咖豆失败")
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		log.NewError("", "GiveChessUserBean commit error:", err, "tx.Rollback().Error :", tx.Rollback().Error)
		return errors.Wrap(err, "事务提交失败")
	}
	return nil
}

// 赠送用户咖豆(重试)
func GiveUserBeanRetry(orderNo string, agentNumber int32, chessUserId, beanNumber int64, intervals []int) bool {
	var retryCh = make(chan bool)
	index := 0

	for {
		go time.AfterFunc(time.Duration(intervals[index])*time.Second, func() {
			err := chessApi.ChessApiGiveUserBean(orderNo, agentNumber, chessUserId, beanNumber)
			log.Info("", "ChessApiGiveUserBean err:", err, orderNo, agentNumber, chessUserId, beanNumber)
			if err == nil {
				retryCh <- true
			} else {
				retryCh <- false
			}
		})

		if <-retryCh {
			return true
		}

		if len(intervals)-1 == index {
			log.Info("", "GiveUserBeanRetry ---次数索引-index:", index)
			return false
		}

		index++
	}
}
