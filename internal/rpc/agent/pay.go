package agent

import (
	"context"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/agent_model"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/common/utils"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	"crazy_server/pkg/proto/agent"
	rpc "crazy_server/pkg/proto/cloud_wallet"
	utils2 "crazy_server/pkg/utils"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// 互娱商城购买咖豆下单(预提交)
func (rpc *AgentServer) ChessShopPurchaseBean(ctx context.Context, req *agent.ChessShopPurchaseBeanReq) (*agent.ChessShopPurchaseBeanResp, error) {
	resp := &agent.ChessShopPurchaseBeanResp{OrderNo: "", CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}
	log.Info(req.OperationId, "start 互娱商城购买咖豆下单(预提交), 参数:", utils2.JsonFormat(req))

	// 加锁
	lockKey := fmt.Sprintf("ChessShopPurchaseBean:%d", req.ChessUserId)
	if err := utils.Lock(ctx, lockKey); err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "操作加锁失败"
		return resp, nil
	}
	defer utils.UnLock(ctx, lockKey)

	//获取推广员信息
	agentInfo, err := imdb.GetAgentByAgentNumber(req.AgentNumber)
	if err != nil || agentInfo.OpenStatus == 0 {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "推广员信息有误"
		return resp, nil
	}

	//校验咖豆配置
	configInfo, err := imdb.GetAgentBeanConfigById(agentInfo.UserId, req.ConfigId)
	if err != nil || configInfo.Status == 0 {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "咖豆配置有误"
		return resp, nil
	}

	//是否为下属成员
	agentMember, err := imdb.AgentNumberByChessUserId(req.ChessUserId)
	if err != nil || agentInfo.UserId != agentMember.UserId {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "该用户不是推广员下成员"
		return resp, nil
	}

	//冻结的咖豆额度
	freezeBeanBalance := rocksCache.GetAgentFreezeBeanBalance(ctx, agentInfo.UserId)
	//校验推广员咖豆余额 + 冻结部分
	if agentInfo.BeanBalance < (configInfo.BeanNumber + int64(configInfo.GiveBeanNumber) + freezeBeanBalance) {
		log.Error(req.OperationId, fmt.Sprintf("推广员(%d),下属成员(%d)购买咖豆,推广员咖豆余额不足,咖豆余额(%d),冻结咖豆(%d)", agentInfo.AgentNumber, req.ChessUserId, agentInfo.BeanBalance, freezeBeanBalance))
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "推广员咖豆余额不足"
		return resp, nil
	}

	orderNo := utils.GetOrderNo() //平台订单号
	//生成订单
	err = imdb.CreatePurchaseBeanOrder(&db.TAgentBeanRechargeOrder{
		BusinessType:      imdb.RechargeOrderBusinessTypeChess,
		UserId:            agentInfo.UserId,
		ChessUserId:       req.ChessUserId,
		ChessUserNickname: agentMember.ChessNickname,
		OrderNo:           orderNo,
		ChessOrderNo:      req.ChessOrderNo,
		Number:            configInfo.BeanNumber,
		GiveNumber:        configInfo.GiveBeanNumber,
		Amount:            configInfo.Amount,
	})

	if err != nil {
		log.Error(req.OperationId, "互娱商城购买咖豆下单(预提交) 生成订单失败。互娱订单号：", req.ChessOrderNo, ",err:", err.Error())
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "生成订单失败"
		return resp, nil
	}

	//冻结推广员咖豆
	_ = rocksCache.FreezeAgentBeanBalance(ctx, agentInfo.UserId, req.ChessUserId, configInfo.BeanNumber+int64(configInfo.GiveBeanNumber))

	resp.OrderNo = orderNo
	resp.GiveBeanNumber = configInfo.GiveBeanNumber
	resp.BeanNumber = configInfo.BeanNumber
	resp.ConfigId = configInfo.Id
	resp.Amount = configInfo.Amount
	return resp, nil
}

// 推广员购买咖豆
func (rpc *AgentServer) AgentPurchaseBean(ctx context.Context, req *agent.AgentPurchaseBeanReq) (*agent.AgentPurchaseBeanResp, error) {
	resp := &agent.AgentPurchaseBeanResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}
	log.Info(req.OperationId, "start 推广员购买咖豆, 参数:", utils2.JsonFormat(req))

	// 加锁
	lockKey := fmt.Sprintf("AgentPurchaseBean:%s", req.UserId)
	if err := utils.Lock(ctx, lockKey); err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "操作加锁失败"
		return resp, nil
	}
	defer utils.UnLock(ctx, lockKey)

	configInfo, err := GetPlatformBeanConfigInfo(req.ConfigId)
	if err != nil || configInfo == nil {
		log.Error(req.OperationId, "获取平台咖豆商城配置缓存-GetAgentPlatformBeanConfigCache err :", err)
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "获取咖豆配置失败"
		return resp, nil
	}

	orderNo := utils.GetOrderNo() //平台订单号

	// rpc 调用新生支付下单接口
	ncountOrderNo, err := RpcCreateThirdPayOrder(ctx, orderNo, configInfo.Amount, req.OperationId)
	if err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = err.Error()
		return resp, nil
	}

	//创建充值咖豆订单
	err = imdb.CreatePurchaseBeanOrder(&db.TAgentBeanRechargeOrder{
		BusinessType:      imdb.RechargeOrderBusinessTypeWeb,
		UserId:            req.UserId,
		ChessUserId:       0,
		ChessUserNickname: "",
		OrderNo:           orderNo,
		ChessOrderNo:      "",
		NcountOrderNo:     ncountOrderNo,
		Number:            configInfo.BeanNumber,
		GiveNumber:        configInfo.GiveBeanNumber,
		Amount:            configInfo.Amount,
	})
	if err != nil {
		errMsg := fmt.Sprintf("处理推广员购买咖豆逻辑逻辑下单失败,推广员id(%s),err:%s", req.UserId, err.Error())
		log.Error("", errMsg)

		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = errMsg
		return resp, nil
	}

	//返回新生支付订单号给 app sdk
	resp.NcountOrderNo = ncountOrderNo

	return resp, nil
}

// rpc 调用新生支付 CreateThirdPayOrder 下单接口
func RpcCreateThirdPayOrder(ctx context.Context, orderNo string, amount int32, operationID string) (string, error) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCloudWalletName, operationID)
	if etcdConn == nil {
		errMsg := operationID + "getcdv3.GetDefaultConn CreateThirdPayOrder == nil"
		log.NewError(operationID, errMsg)
		return "", errors.New(errMsg)
	}

	//组装数据
	rpcReq := rpc.CreateThirdPayOrderReq{
		MerchantId:  config.Config.Agent.MerchantId, //商户号
		MerOrderId:  orderNo,
		NotifyUrl:   config.Config.Agent.AgentRechargeNotifyUrl,
		Amount:      amount,
		Remark:      "推广员充值咖豆",
		OperationID: operationID,
	}

	client := rpc.NewCloudWalletServiceClient(etcdConn)
	RpcResp, _ := client.CreateThirdPayOrder(ctx, &rpcReq)
	if RpcResp.CommonResp != nil && RpcResp.CommonResp.ErrCode != 0 {
		log.NewError(operationID, "client.CreateThirdPayOrder 调用失败:", RpcResp.CommonResp.ErrMsg)
		return "", errors.New(RpcResp.CommonResp.ErrMsg)
	}

	return RpcResp.OrderNo, nil
}

// 获取平台咖豆配置项
func GetPlatformBeanConfigInfo(configId int32) (*imdb.BeanShopConfig, error) {
	//获取平台咖豆redis缓存配置
	beanConfig, err := rocksCache.GetAgentPlatformBeanConfigCache()
	if err != nil || len(beanConfig) == 0 {
		return nil, errors.Wrap(err, "获取平台咖豆redis缓存配置失败")
	}

	for _, v := range beanConfig {
		if v.ConfigId == configId {
			return v, nil
		}
	}

	return nil, errors.Wrap(err, "获取平台咖豆redis缓存配置失败.")
}
