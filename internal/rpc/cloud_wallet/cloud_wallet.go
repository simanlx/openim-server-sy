package cloud_wallet

import (
	"context"
	"crazy_server/pkg/cloud_wallet/ncount"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	promePkg "crazy_server/pkg/common/prometheus"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	"crazy_server/pkg/proto/chat"
	"crazy_server/pkg/proto/cloud_wallet"
	"crazy_server/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	UserMainAccountPrefix   = "main_"   //主账户前缀 完整账户id + 用户id
	UserPacketAccountPrefix = "packet_" //红包账户前缀 完整账户id + 用户id
)

type CloudWalletServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string

	// 依赖钱包服务
	count ncount.NCounter
	cloud_wallet.UnsafeCloudWalletServiceServer
}

func NewRpcCloudWalletServer(port int) *CloudWalletServer {
	log.NewPrivateLog("crazy_server_cloud_wallet")
	StarCorn()
	return &CloudWalletServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImCloudWalletName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
		count:           ncount.NewCounter(),
	}
}

func (rpc *CloudWalletServer) Run() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, "rpc auth start...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(rpc.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + rpc.rpcRegisterName)
	}
	log.NewInfo(operationID, "listen network success, ", address, listener)
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()

	//service registers with etcd
	cloud_wallet.RegisterCloudWalletServiceServer(srv, rpc)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)

	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log.NewError(operationID, "RegisterEtcd failed ", err.Error(),
			rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))

	}
	log.NewInfo(operationID, "RegisterAuthServer ok ", rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError(operationID, "Serve failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "rpc auth ok")
}

// 获取云账户信息
func (rpc *CloudWalletServer) UserNcountAccount(_ context.Context, req *cloud_wallet.UserNcountAccountReq) (*cloud_wallet.UserNcountAccountResp, error) {
	resp := &cloud_wallet.UserNcountAccountResp{Step: 0, BalAmount: "0", CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	//获取用户账户信息
	accountInfo, err := imdb.GetNcountAccountByUserId(req.UserId)
	if err != nil || accountInfo.Id <= 0 {
		return resp, nil
	}

	//调新生支付接口，获取用户信息
	accountResp, err := rpc.count.CheckUserAccountInfo(&ncount.CheckUserAccountReq{
		OrderID: ncount.GetMerOrderID(),
		UserID:  accountInfo.MainAccountId,
	})

	log.Info(req.OperationID, "获取云账户信息-CheckUserAccountInfo->", utils.JsonFormat(accountResp), err)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("查询账户信息失败(%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	} else {
		if accountResp.ResultCode != ncount.ResultCodeSuccess {
			resp.CommonResp.ErrMsg = fmt.Sprintf("查询账户信息失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg)
			resp.CommonResp.ErrCode = 400
			return resp, nil
		}
	}

	//绑定的银行卡列表
	bindCardsList := make([]*cloud_wallet.BindCardsList, 0)
	if len(accountResp.BindCardAgrNoList) > 0 {
		//获取用户银行卡信息列表
		bindCardAgrNoBank := map[string]*db.FNcountBankCard{}
		bankcardList, _ := imdb.GetUserBankcardByUserId(req.UserId)
		for _, v := range bankcardList {
			bindCardAgrNoBank[v.BindCardAgrNo] = v
		}

		bindCards := make([]ncount.NAccountBankCard, 0)
		err = json.Unmarshal([]byte(accountResp.BindCardAgrNoList), &bindCards)
		if err == nil {
			for _, v := range bindCards {
				mobile := ""
				if bc, ok := bindCardAgrNoBank[v.BindCardAgrNo]; ok {
					mobile = bc.Mobile
				}

				bindCardsList = append(bindCardsList, &cloud_wallet.BindCardsList{
					BankCode:      v.BankCode,
					CardNo:        v.CardNo,
					BindCardAgrNo: v.BindCardAgrNo,
					Mobile:        mobile,
				})
			}
		}

		accountInfo.OpenStep = 3 //状态-绑定了银行卡
	}

	//删除缓存
	_ = rocksCache.DeleteAccountInfoFromCache(req.UserId)

	//精度
	balAmount, _ := decimal.NewFromString(accountResp.BalAmount)
	return &cloud_wallet.UserNcountAccountResp{
		Step:             accountInfo.OpenStep,
		IdCard:           accountInfo.IdCard,
		RealName:         accountInfo.RealName,
		AccountStatus:    accountInfo.OpenStatus,
		BalAmount:        balAmount.Mul(decimal.NewFromInt(100)).String(), //转换为分
		AvailableBalance: accountResp.AvailableBalance,
		BindCardsList:    bindCardsList,
	}, nil
}

// 身份证实名认证
func (rpc *CloudWalletServer) IdCardRealNameAuth(_ context.Context, req *cloud_wallet.IdCardRealNameAuthReq) (*cloud_wallet.IdCardRealNameAuthResp, error) {
	resp := &cloud_wallet.IdCardRealNameAuthResp{Step: 0, CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	//获取用户账户信息
	accountInfo, err := imdb.GetNcountAccountByUserId(req.UserId)
	if accountInfo != nil && accountInfo.Id > 0 {
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = "已实名认证,请勿重复操作"
		return resp, nil
	}

	//一个身份证最多只能实名2个用户
	authNumber := imdb.IdCardRealNameAuthNumber(req.IdCard)
	if authNumber >= 2 {
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = "身份证已被其他用户实名"
		return resp, nil
	}

	//组装数据
	info := &db.FNcountAccount{
		UserID:      req.UserId,
		Mobile:      req.Mobile,
		RealName:    req.RealName,
		IdCard:      req.IdCard,
		OpenStep:    1,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	//调新生支付接口-开户
	errGroup := new(errgroup.Group)
	accountIds := []string{
		fmt.Sprintf("%s%s", UserMainAccountPrefix, info.UserID),
		fmt.Sprintf("%s%s", UserPacketAccountPrefix, info.UserID),
	}
	for _, account := range accountIds {
		id := account
		errGroup.Go(func() error {
			accountResp, err := rpc.count.NewAccount(&ncount.NewAccountReq{
				OrderID: ncount.GetMerOrderID(),
				MsgCipherText: &ncount.NewAccountMsgCipherText{
					MerUserId: id,
					Mobile:    info.Mobile,
					UserName:  info.RealName,
					CertNo:    info.IdCard,
				},
			})

			log.Info(req.OperationID, "实名认证-NewAccount->", utils.JsonFormat(accountResp), err)
			if err != nil {
				return errors.New(fmt.Sprintf("实名认证失败(%s)", err.Error()))
			} else {
				if accountResp.ResultCode != ncount.ResultCodeSuccess {
					return errors.New(fmt.Sprintf("实名认证失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg))
				}
			}

			//主账户
			if id == fmt.Sprintf("%s%s", UserMainAccountPrefix, info.UserID) {
				info.MainAccountId = accountResp.UserId
			} else {
				info.PacketAccountId = accountResp.UserId
			}

			return nil
		})
	}

	if err = errGroup.Wait(); err != nil {
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	//实名数据入库
	err = imdb.CreateNcountAccount(info)
	if err != nil {
		log.Error(req.OperationID, "实名认证数据入库失败:%s", err.Error())
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = fmt.Sprintf("实名认证数据入库失败:%s", err.Error())
		return resp, nil
	}

	resp.Step = 1
	return resp, nil
}

// 校验用户支付密码
func (rpc *CloudWalletServer) CheckPaymentSecret(_ context.Context, req *cloud_wallet.CheckPaymentSecretReq) (*cloud_wallet.CheckPaymentSecretResp, error) {
	resp := &cloud_wallet.CheckPaymentSecretResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	//获取用户账户信息
	accountInfo, err := imdb.GetNcountAccountByUserId(req.UserId)
	if err != nil || accountInfo.Id <= 0 {
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = "账户信息不存在"
		return resp, nil
	}

	//验证支付密码
	if len(accountInfo.PaymentPassword) == 0 || req.PaymentSecret != accountInfo.PaymentPassword {
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = "支付密码错误"
		return resp, nil
	}
	return resp, nil
}

// 设置用户支付密码
func (rpc *CloudWalletServer) SetPaymentSecret(_ context.Context, req *cloud_wallet.SetPaymentSecretReq) (*cloud_wallet.SetPaymentSecretResp, error) {
	resp := &cloud_wallet.SetPaymentSecretResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	if req.Type == 2 {
		//忘记支付密码、校验验证码
		if !RpcForgetPayPasswordVerifyCode(req.UserId, req.Code, req.OperationID) {
			resp.CommonResp.ErrCode = 400
			resp.CommonResp.ErrMsg = "验证码错误"
			return resp, nil
		}
	} else {
		//获取用户账户信息
		accountInfo, err := imdb.GetNcountAccountByUserId(req.UserId)
		if err != nil || accountInfo.Id <= 0 {
			resp.CommonResp.ErrCode = 400
			resp.CommonResp.ErrMsg = "账户信息不存在"
			return resp, nil
		}

		if req.Type == 1 {
			if len(accountInfo.PaymentPassword) > 1 {
				resp.CommonResp.ErrCode = 400
				resp.CommonResp.ErrMsg = "调用接口错误"
				return resp, nil
			}
		} else {
			if req.OriginalPaymentSecret != accountInfo.PaymentPassword {
				resp.CommonResp.ErrCode = 400
				resp.CommonResp.ErrMsg = "原支付密码错误"
				return resp, nil
			}
		}
	}

	//修改支付密码
	err := imdb.UpdateNcountAccountField(req.UserId, map[string]interface{}{"payment_password": req.PaymentSecret, "open_step": 2})
	if err != nil {
		log.Error(req.OperationID, "保存支付密码失败:%s", err.Error())
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = fmt.Sprintf("保存支付密码失败,err:%s", err.Error())
		return resp, nil
	}

	//删除缓存
	_ = rocksCache.DeleteAccountInfoFromCache(req.UserId)

	resp.Step = 2
	return resp, nil
}

// rpc 调用Open-IM-Enterprise --> ForgetPayPasswordVerifyCode 接口
func RpcForgetPayPasswordVerifyCode(userId, code, operationID string) bool {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImChatName, operationID)
	if etcdConn == nil {
		errMsg := operationID + "getcdv3.GetDefaultConn RpcForgetPayPasswordVerifyCode == nil"
		log.NewError(operationID, errMsg)
		return false
	}

	chatReq := &chat.ForgetPayPasswordVerifyCodeReq{
		UserId:      userId,
		Code:        code,
		OperationID: operationID,
	}

	client := chat.NewChatClient(etcdConn)
	rpcResp, _ := client.ForgetPayPasswordVerifyCode(context.Background(), chatReq)
	if rpcResp.CommonResp == nil || rpcResp.CommonResp.ErrCode == 0 {
		return true
	}
	log.NewError(operationID, "client.ForgetPayPasswordVerifyCode 验证失败:", rpcResp.CommonResp.ErrMsg)

	return false
}

// 云钱包收支明细
func (rpc *CloudWalletServer) CloudWalletRecordList(_ context.Context, req *cloud_wallet.CloudWalletRecordListReq) (*cloud_wallet.CloudWalletRecordListResp, error) {
	resp := &cloud_wallet.CloudWalletRecordListResp{}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Size <= 0 {
		req.Size = 20
	}

	// 获取总量
	tradein, tradOut, err := imdb.GetNcountTradeTotal(req.UserId, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	resp.Totalincome = tradein
	resp.Totalpay = tradOut

	//条件获取列表数据
	list, count, err := imdb.FindNcountTradeList(req.UserId, req.StartTime, req.EndTime, req.Page, req.Size)
	if err != nil {
		return resp, nil
	}

	recordList := make([]*cloud_wallet.RecordList, 0)
	for _, v := range list {
		recordList = append(recordList, &cloud_wallet.RecordList{
			Id:                v.ID,
			Describe:          v.Describe,
			Amount:            v.Amount,
			CreatedTime:       v.CreatedTime.Format("2006-01-02 15:04:05"),
			RelevancePacketId: v.PacketID,
			AfterAmount:       v.AfterAmount,
			Type:              v.Type,
		})
	}

	resp.Total = int32(count)
	resp.RecordList = recordList
	return resp, nil
}

// 绑定用户银行卡
func (rpc *CloudWalletServer) BindUserBankcard(_ context.Context, req *cloud_wallet.BindUserBankcardReq) (*cloud_wallet.BindUserBankcardResp, error) {
	resp := &cloud_wallet.BindUserBankcardResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	//获取用户账户信息
	accountInfo, err := imdb.GetNcountAccountByUserId(req.UserId)
	if err != nil || accountInfo.Id <= 0 {
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = fmt.Sprintf("查询账户数据失败 %s,error:%s", req.UserId, err.Error())
		return resp, nil
	}

	merOrderId := ncount.GetMerOrderID()
	accountResp, err := rpc.count.BindCard(&ncount.BindCardReq{
		MerOrderId: merOrderId,
		BindCardMsgCipherText: ncount.BindCardMsgCipherText{
			CardNo:       req.BankCardNumber,
			HolderName:   req.CardOwner,
			MobileNo:     req.Mobile,
			IdentityType: "1",
			IdentityCode: accountInfo.IdCard,
			UserId:       accountInfo.MainAccountId,
		},
	})

	log.Info(req.OperationID, "绑定银行卡-BindUserBankcard->", utils.JsonFormat(accountResp), err)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("绑定银行卡失败(%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	} else {
		if accountResp.ResultCode != ncount.ResultCodeSuccess {
			resp.CommonResp.ErrMsg = fmt.Sprintf("绑定银行卡失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg)
			resp.CommonResp.ErrCode = 400
			return resp, nil
		}
	}

	info := &db.FNcountBankCard{
		UserId:            req.UserId,
		MerOrderId:        merOrderId,
		NcountOrderId:     accountResp.NcountOrderId,
		NcountUserId:      accountInfo.MainAccountId,
		Mobile:            req.Mobile,
		CardOwner:         req.CardOwner,
		BankCardNumber:    req.BankCardNumber,
		Cvv2:              req.Cvv2,
		CardAvailableDate: req.CardAvailableDate,
		CreatedTime:       time.Now(),
		UpdatedTime:       time.Now(),
	}

	//数据入库
	err = imdb.BindUserBankcard(info)
	if err != nil {
		log.Error(req.OperationID, "银行卡数据入库失败:%s", err.Error())
		resp.CommonResp.ErrMsg = "银行卡数据入库失败"
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	resp.BankCardId = info.Id
	return resp, nil
}

// 绑定用户银行卡确认code
func (rpc *CloudWalletServer) BindUserBankcardConfirm(_ context.Context, req *cloud_wallet.BindUserBankcardConfirmReq) (*cloud_wallet.BindUserBankcardConfirmResp, error) {
	resp := &cloud_wallet.BindUserBankcardConfirmResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	//获取绑定的银行卡信息
	bankCardInfo, err := imdb.GetNcountBankCardById(req.BankCardId, req.UserId)
	if err != nil || bankCardInfo.Id <= 0 {
		resp.CommonResp.ErrMsg = fmt.Sprintf("查询银行卡数据失败,error:%s", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	//已绑定
	if bankCardInfo.IsBind == 1 {
		return resp, err
	}

	//新生支付确定接口
	accountResp, err := rpc.count.BindCardConfirm(&ncount.BindCardConfirmReq{
		MerOrderId: ncount.GetMerOrderID(),
		BindCardConfirmMsgCipherText: ncount.BindCardConfirmMsgCipherText{
			NcountOrderId: bankCardInfo.NcountOrderId,
			SmsCode:       req.SmsCode,
			MerUserIp:     req.MerUserIp,
		},
	})

	log.Info(req.OperationID, "绑定用户银行卡确认-BindUserBankcardConfirm->", utils.JsonFormat(accountResp), err)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("绑定用户银行卡确认失败(%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	} else {
		if accountResp.ResultCode != ncount.ResultCodeSuccess {
			resp.CommonResp.ErrMsg = fmt.Sprintf("绑定用户银行卡确认失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg)
			resp.CommonResp.ErrCode = 400
			return resp, nil
		}
	}

	//更新数据
	err = imdb.BindUserBankcardConfirm(bankCardInfo.Id, req.UserId, accountResp.BindCardAgrNo, accountResp.BankCode)
	if err != nil {
		log.Error(req.OperationID, "更新银行卡数据失败:%s", err.Error())
		resp.CommonResp.ErrMsg = fmt.Sprintf("更新银行卡数据失败 (%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	resp.BankCardId = bankCardInfo.Id
	return resp, err
}

// 解绑用户银行卡
func (rpc *CloudWalletServer) UnBindingUserBankcard(_ context.Context, req *cloud_wallet.UnBindingUserBankcardReq) (*cloud_wallet.UnBindingUserBankcardResp, error) {
	resp := &cloud_wallet.UnBindingUserBankcardResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	//获取绑定的银行卡信息
	bankCardInfo, err := imdb.GetNcountBankCardByBindCardAgrNo(req.BindCardAgrNo, req.UserId)
	if err != nil || bankCardInfo.Id <= 0 {
		resp.CommonResp.ErrMsg = fmt.Sprintf("查询银行卡数据失败,error:%s", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	//新生支付确定接口
	accountResp, err := rpc.count.UnbindCard(&ncount.UnBindCardReq{
		MerOrderId: ncount.GetMerOrderID(),
		UnBindCardMsgCipher: ncount.UnBindCardMsgCipher{
			OriBindCardAgrNo: bankCardInfo.BindCardAgrNo,
			UserId:           bankCardInfo.NcountUserId,
		},
	})

	log.Info(req.OperationID, "解绑银行卡-UnBindingUserBankcard->", utils.JsonFormat(accountResp), err)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("解绑银行卡失败(%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	} else {
		if accountResp.ResultCode != ncount.ResultCodeSuccess {
			resp.CommonResp.ErrMsg = fmt.Sprintf("解绑银行卡失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg)
			resp.CommonResp.ErrCode = 400
			return resp, nil
		}
	}

	//更新数据库
	err = imdb.UnBindUserBankcard(bankCardInfo.Id, req.UserId)
	if err != nil {
		log.Error(req.OperationID, "更新银行卡数据失败:%s", err.Error())
		resp.CommonResp.ErrMsg = fmt.Sprintf("更新银行卡数据失败 (%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	return &cloud_wallet.UnBindingUserBankcardResp{}, err
}

// 银行卡充值
func (rpc *CloudWalletServer) UserRecharge(_ context.Context, req *cloud_wallet.UserRechargeReq) (*cloud_wallet.UserRechargeResp, error) {
	resp := &cloud_wallet.UserRechargeResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	// 获取银行卡信息
	bankCardInfo, err := imdb.GetNcountBankCardByBindCardAgrNo(req.BindCardAgrNo, req.UserId)
	if err != nil {
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = fmt.Sprintf("获取银行卡信息失败%s", err.Error())
		return resp, nil
	}

	//充值支付
	merOrderID := ncount.GetMerOrderID()
	accountResp, err := rpc.count.QuickPayOrder(&ncount.QuickPayOrderReq{
		MerOrderId: merOrderID,
		QuickPayMsgCipher: ncount.QuickPayMsgCipher{
			PayType:       "3", //绑卡协议号充值
			TranAmount:    cast.ToString(cast.ToFloat64(req.Amount) / 100),
			NotifyUrl:     config.Config.Ncount.Notify.RechargeNotifyUrl,
			BindCardAgrNo: bankCardInfo.BindCardAgrNo,
			ReceiveUserId: bankCardInfo.NcountUserId, //收款账户
			UserId:        bankCardInfo.NcountUserId,
			SubMerchantId: ncount.SUB_MERCHANT_ID, // 子商户编号
		}})

	log.Info(req.OperationID, "银行卡充值-UserRecharge->", utils.JsonFormat(accountResp), err)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("银行卡充值失败(%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	} else {
		if accountResp.ResultCode != ncount.ResultCodeSuccess {
			resp.CommonResp.ErrMsg = fmt.Sprintf("银行卡充值失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg)
			resp.CommonResp.ErrCode = 400
			return resp, nil
		}
	}

	//增加账户变更日志
	err = AddNcountTradeLog(BusinessTypeBankcardRecharge, req.Amount, req.UserId, bankCardInfo.NcountUserId, merOrderID, accountResp.NcountOrderId, "")
	if err != nil {
		log.Error(req.OperationID, "增加账户变更日志失败[%s]", err.Error(), "参数：", BusinessTypeBankcardRecharge, req.Amount, req.UserId, bankCardInfo.NcountUserId, accountResp.NcountOrderId)
	}

	return &cloud_wallet.UserRechargeResp{
		OrderNo: accountResp.NcountOrderId,
	}, nil
}

// 账户充值code 确认
func (rpc *CloudWalletServer) UserRechargeConfirm(_ context.Context, req *cloud_wallet.UserRechargeConfirmReq) (*cloud_wallet.UserRechargeConfirmResp, error) {
	resp := &cloud_wallet.UserRechargeConfirmResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	// 获取记录信息
	tradeInfo, err := imdb.GetFNcountTradeByOrderNo(req.MerOrderId, req.UserId)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("获取充值记录信息失败%s", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	//新生支付确认接口
	accountResp, err := rpc.count.QuickPayConfirm(&ncount.QuickPayConfirmReq{
		MerOrderId: ncount.GetMerOrderID(),
		QuickPayConfirmMsgCipher: ncount.QuickPayConfirmMsgCipher{
			NcountOrderId:        tradeInfo.ThirdOrderNo,
			SmsCode:              req.SmsCode,
			PaymentTerminalInfo:  "02|AA01BB",
			ReceiverTerminalInfo: "01|00001|CN|469023",
			DeviceInfo:           "192.168.0.1|E1E2E3E4E5E6|123456789012345|20000|898600MFSSYYGXXXXXXP|H1H2H3H4H5H6|AABBCC",
		},
	})

	log.Info(req.OperationID, "充值确认-UserRechargeConfirm->", utils.JsonFormat(accountResp), err)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("充值确认失败(%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	} else {
		if accountResp.ResultCode == ncount.ResultCodeFail {
			resp.CommonResp.ErrMsg = fmt.Sprintf("充值确认失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg)
			resp.CommonResp.ErrCode = 400
			return resp, nil
		}
	}

	return resp, nil
}

// 提现
func (rpc *CloudWalletServer) UserWithdrawal(_ context.Context, req *cloud_wallet.DrawAccountReq) (*cloud_wallet.DrawAccountResp, error) {
	resp := &cloud_wallet.DrawAccountResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	//获取用户账户信息
	accountInfo, err := imdb.GetNcountAccountByUserId(req.UserId)
	if err != nil || accountInfo.Id <= 0 {
		resp.CommonResp.ErrMsg = fmt.Sprintf("查询账户数据失败 %s,error:%s", req.UserId, err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	//验证支付密码
	if len(accountInfo.PaymentPassword) < 10 || req.PaymentPassword != accountInfo.PaymentPassword {
		resp.CommonResp.ErrMsg = "支付密码错误"
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	// 获取银行卡信息
	bankCardInfo, err := imdb.GetNcountBankCardByBindCardAgrNo(req.BindCardAgrNo, req.UserId)
	if err != nil {
		resp.CommonResp.ErrMsg = "获取银行卡信息错误"
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	//调用新生支付提现接口
	merOrderID := ncount.GetMerOrderID()
	accountResp, err := rpc.count.Withdraw(&ncount.WithdrawReq{
		MerOrderID: merOrderID,
		MsgCipher: ncount.WithdrawMsgCipher{
			BusinessType:    "08",
			TranAmount:      cast.ToFloat32(cast.ToFloat64(req.Amount) / 100),
			UserId:          bankCardInfo.NcountUserId,
			BindCardAgrNo:   req.BindCardAgrNo,
			NotifyUrl:       config.Config.Ncount.Notify.WithdrawNotifyUrl,
			PaymentTerminal: "02|AA01BB",
			DeviceInfo:      "192.168.0.1|E1E2E3E4E5E6|123456789012345|20000|898600MFSSYYGXXXXXXP|H1H2H3H4H5H6|AABBCC",
		},
	})

	log.Info(req.OperationID, "银行卡提现-UserWithdrawal->", utils.JsonFormat(accountResp), err)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("银行卡提现失败(%s)", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	} else {
		if accountResp.ResultCode == ncount.ResultCodeFail {
			resp.CommonResp.ErrMsg = fmt.Sprintf("银行卡提现失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg)
			resp.CommonResp.ErrCode = 400
			return resp, nil
		}
	}

	//增加账户变更日志
	err = AddNcountTradeLog(BusinessTypeBankcardWithdrawal, req.Amount, req.UserId, bankCardInfo.NcountUserId, merOrderID, accountResp.NcountOrderId, "")
	if err != nil {
		log.Error(req.OperationID, "增加账户变更日志失败[%s]", err.Error(), "参数：", BusinessTypeBankcardWithdrawal, req.Amount, req.UserId, bankCardInfo.NcountUserId, accountResp.NcountOrderId)
	}

	resp.OrderNo = accountResp.NcountOrderId
	return resp, nil
}

// 删除云钱包明细
func (rpc *CloudWalletServer) CloudWalletRecordDel(_ context.Context, req *cloud_wallet.CloudWalletRecordDelReq) (*cloud_wallet.CloudWalletRecordDelResp, error) {
	resp := &cloud_wallet.CloudWalletRecordDelResp{CommonResp: &cloud_wallet.CommonResp{ErrCode: 0, ErrMsg: ""}}

	//软删除记录
	err := imdb.DelNcountTradeRecord(req.DelType, req.RecordId, req.UserId)
	if err != nil {
		resp.CommonResp.ErrMsg = fmt.Sprintf("删除记录失败,%s", err.Error())
		resp.CommonResp.ErrCode = 400
		return resp, nil
	}

	return resp, nil
}
