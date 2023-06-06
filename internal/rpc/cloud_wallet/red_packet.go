package cloud_wallet

import (
	"context"
	ncount "crazy_server/pkg/cloud_wallet/ncount"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/db"
	commonDB "crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	"crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	imdb2 "crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/contrive_msg"
	pb "crazy_server/pkg/proto/cloud_wallet"
	"crazy_server/pkg/tools/redpacket"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// 发送红包接口
func (rpc *CloudWalletServer) SendRedPacket(ctx context.Context, in *pb.SendRedPacketReq) (*pb.SendRedPacketResp, error) {
	handler := &handlerSendRedPacket{
		OperateID:  in.GetOperationID(),
		merOrderID: ncount.GetMerOrderID(),
		count:      rpc.count,
	}
	return handler.SendRedPacket(in)
}

type handlerSendRedPacket struct {
	OperateID  string
	merOrderID string
	count      ncount.NCounter
}

// 发送红包
func (h *handlerSendRedPacket) SendRedPacket(req *pb.SendRedPacketReq) (*pb.SendRedPacketResp, error) {
	var (
		result = &pb.SendRedPacketResp{
			CommonResp: &pb.CommonResp{
				ErrCode: 0,
				ErrMsg:  "发送成功",
			},
		}
	)
	// 1. 校验参数
	if err := h.validateParam(req); err != nil {
		return nil, err
	}

	// ========================================= 验证发送用户的信息=========================================

	//  缓存
	userAC, err := rocksCache.GetUserAccountInfoFromCache(req.UserId)
	if err != nil {
		log.Error(req.OperationID, "查询用户信息失败", zap.Error(err))
		result.CommonResp.ErrMsg = "用户未实名"
		result.CommonResp.ErrCode = 400
		return result, nil
	}

	if req.Password != userAC.PaymentPassword {
		result.CommonResp.ErrMsg = "支付密码错误"
		result.CommonResp.ErrCode = 400
		return result, nil
	}

	// ========================================= 查看发送用户ID是否存在 =========================================
	if req.PacketType == 1 {
		// cache
		user, err := rocksCache.GetUserInfoFromCache(req.RecvID)
		if err != nil {
			if err == sql.ErrNoRows {
				result.CommonResp.ErrMsg = "您发送红包的用户不存在"
				result.CommonResp.ErrCode = 400
				return result, nil
			}
			return nil, errors.New("查询用户信息失败")
		}
		if user.UserID == "" {
			result.CommonResp.ErrMsg = "您发送红包的用户不存在"
			result.CommonResp.ErrCode = 400
			return result, nil
		}
	} else {
		// cache
		group, err := rocksCache.GetGroupInfoFromCache(req.RecvID)
		if err != nil {
			if err == sql.ErrNoRows {
				result.CommonResp.ErrMsg = "群聊不存在"
				result.CommonResp.ErrCode = 400
				return result, nil
			}
			return nil, errors.New("查询群信息失败")
		}
		if group.GroupID == "" {
			result.CommonResp.ErrMsg = "群聊不存在"
			result.CommonResp.ErrCode = 400
			return result, nil
		}

		// mysql
		ok := imdb2.IsExistGroupMember(req.RecvID, req.UserId)
		if !ok {
			result.CommonResp.ErrMsg = "您不在当前群聊中"
			result.CommonResp.ErrCode = 400
			return result, nil
		}
	}

	// ========================================= 计算发送红包总金额 =========================================
	var amount int64 = 0

	if req.PacketType == 1 || req.IsExclusive == 1 || req.Number == 1 {
		// 这种就是一个红包
		amount = req.Amount
	} else {
		// 走的群聊红包
		if req.IsLucky != 1 {
			// 普通红包
			amount = req.Amount * int64(req.Number)
		} else {
			amount = req.Amount
		}
	}

	totalAmount := cast.ToString(cast.ToFloat64(amount) / 100)

	// ========================================= 这里是用户转账的金额 =========================================
	fmt.Println("\n发送总金额为： ", amount, totalAmount, "\n")

	res := &pb.SendRedPacketResp{
		RedPacketID: "",
	}

	// 3. 判断支付类型
	if req.SendType == 1 {
		transferMssgae, RedPacketID, err := h.walletTransfer(userAC, req, amount, totalAmount)
		if err != nil {
			log.Error(req.OperationID, "转账失败", err)
			return nil, err
		}

		if transferMssgae != "" {
			result.CommonResp.ErrMsg = transferMssgae
			result.CommonResp.ErrCode = 400
			return result, nil
		}

		res.RedPacketID = RedPacketID

		// 回调处理红包
		err = HandleSendPacketResult(RedPacketID, req.OperationID)
		if err != nil {
			log.Error(req.OperationID, "HandleSendPacketResult error", zap.Error(err))
			return nil, err
		}
	} else {
		merorderID := ncount.GetMerOrderID()
		redpacket, err := h.recordRedPacket(req, amount, userAC.PacketAccountId, merorderID)
		if err != nil {
			log.Error(req.OperationID, "record red packet error", zap.Error(err))
			return nil, err
		}
		res.RedPacketID = redpacket.PacketID
		// 走银行卡转账
		ncountID, err := h.bankTransfer(redpacket.PacketID, req, merorderID)
		// 这里是调用银行卡转账接口
		if err != nil {
			log.Error(req.OperationID, "bankTransfer error", zap.Error(err))
			return nil, err
		}

		redpacket.NcountOrderID = ncountID
		redpacket.Status = 1
		redpacket.Remark = "银行卡发起转账，待短信确认"

		// 更新红包状态
		err = imdb.UpdateRedPacketInfo(redpacket.PacketID, redpacket)
		if err != nil {
			log.Error(req.OperationID, "record red packet error", zap.Error(err))
			return nil, err
		}
	}
	return res, nil
}

func (h *handlerSendRedPacket) validateParam(req *pb.SendRedPacketReq) error {
	if len(req.UserId) <= 0 {
		return errors.New("user_id 不能为空")
	}

	// 检测红包类型
	if req.PacketType != 1 && req.PacketType != 2 {
		return errors.New("red_packet_type 错误输入 ")
	}

	// 检测是否为幸运红包
	if req.IsLucky != 1 && req.IsLucky != 0 {
		return errors.New("is_lucky 错误输入 ")
	}

	// 检测是否为专属红包
	if req.IsExclusive != 1 && req.IsExclusive != 0 {
		return errors.New("is_exclusive 错误输入 ")
	}

	// 专属红包必须要专属用户id
	if req.IsExclusive == 1 && req.ExclusiveUserID == "" {
		return errors.New("是专属红包就必须存在ExclusiveUserID")
	}
	// 红包必须要标题
	if req.PacketTitle == "" {
		return errors.New("red_packet_title 红包title不能为空")
	}

	// 红包金额必须大于0
	if req.Amount <= 0 {
		return errors.New(fmt.Sprintf("red_packet_amount 红包金额必须为大于0 , %v", req.Amount/100))
	}

	// 红包个数必须大于0
	if req.Number <= 0 {
		return errors.New("red_packet_number 个数必须大于0")
	}

	if req.IsExclusive == 1 && req.PacketType != 2 {
		return errors.New("IsExclusive 属性红包必须是PacketType = 2 ")
	}

	// 检测发送方式
	if req.SendType != 1 && req.SendType != 2 {
		return errors.New("send_type 发送方式输入错误 ")
	}

	if req.RecvID == "" {
		return errors.New("RecvID 不能为空")
	}

	// 如果是单个红包
	if req.PacketType == 1 || req.IsExclusive == 1 {
		if req.Number != 1 {
			return errors.New("number : 数量错误,当req.PacketType =1 或者 req.IsExclusive = 1 时，number 必须为1")
		}
	}

	// 检测金额 和 红包个数是否合理
	if req.Amount < int64(req.Number) {
		return errors.New("红包金额不能小于红包个数")
	}

	// 如果红包是lucky红包
	if req.IsLucky == 1 {
		if req.IsExclusive == 1 || req.ExclusiveUserID != "" {
			return errors.New("is_lucky 为1的时候 和 is_exclusive 和 exclusive_user_id 必须为0和空")
		}
	}

	return nil
}

// 验证业务上的逻辑错误
func (h *handlerSendRedPacket) checkGroupPacketState(req *pb.SendRedPacketReq) string {
	// 1.用户是否在群里
	if req.PacketType == 2 {
		ok := imdb2.IsExistGroupMember(req.RecvID, req.UserId)
		if !ok {
			return "用户不在群里"
		}
	}
	return ""
}

func (h handlerSendRedPacket) createRedpacketID(packetType int32, UserID string) string {
	rand.Seed(time.Now().UnixNano())
	redID := fmt.Sprintf("%v%v%v%v", packetType, UserID, time.Now().Unix(), rand.Intn(100000))
	return redID
}

// 创建红包信息
func (h *handlerSendRedPacket) recordRedPacket(in *pb.SendRedPacketReq, amount int64, packetID /*发红包的用户ID*/, merOrderID string) (*db.FPacket /* red packet ID */, error) {

	redPacket := &db.FPacket{
		PacketID:             h.createRedpacketID(in.PacketType, in.UserId),
		UserID:               in.UserId,
		UserRedpacketAccount: packetID,
		PacketType:           in.PacketType,
		IsLucky:              in.IsLucky,
		ExclusiveUserID:      in.ExclusiveUserID,
		PacketTitle:          in.PacketTitle,
		Amount:               in.Amount,
		Number:               in.Number,
		TotalAmount:          amount,
		MerOrderID:           merOrderID,
		OperateID:            h.OperateID,
		SendType:             in.SendType,
		BindCardAgrNo:        in.BindCardAgrNo,
		RecvID:               in.RecvID, // 接收ID
		Remain:               int64(in.Number),
		RemainAmout:          amount,
		ExpireTime:           time.Now().Unix() + 60*60*24,
		CreatedTime:          time.Now().Unix(),
		UpdatedTime:          time.Now().Unix(),
		Status:               0, // 红包被创建，但是还未掉第三方的内容
		IsExclusive:          in.IsExclusive,
	}

	// 这里创建用户的时候
	err := imdb.RedPacketCreateData(redPacket)
	if err != nil {
		return nil, errors.Wrap(err, "创建红包信息错误")
	}

	return redPacket, nil
}

// @return Param 提示错误消息
// @return Param 红包ID
// @return Param 错误
func (h *handlerSendRedPacket) walletTransfer(fncount *db.FNcountAccount, in *pb.SendRedPacketReq, tAmount int64, totalAmount string) (string, string, error) {
	var res string
	// 1. 获取用户的钱包账户
	merOrderID := h.merOrderID
	req := &ncount.TransferReq{
		MerOrderId: merOrderID,
		TransferMsgCipher: ncount.TransferMsgCipher{
			PayUserId:     fncount.MainAccountId,
			ReceiveUserId: fncount.PacketAccountId,
			TranAmount:    totalAmount, //分转元
		},
	}

	escap := time.Now()
	transferResult, err := h.count.Transfer(req)
	log.Info(in.OperationID, "transfer req", req, "耗费时间:", time.Since(escap))
	fmt.Printf("\n 第三方调用耗时： %v \n", time.Since(escap))
	if err != nil {
		log.Error(in.OperationID, "调用新生支付出现错误", transferResult)
		res = "第三方支付出现网络错误,请稍后重试，操作：" + h.OperateID
		return res, "", nil
	}

	// ========================================================下面是成功返回  =========================================
	if transferResult.ResultCode != ncount.ResultCodeSuccess {
		remark := "余额发送红包失败： "
		co, _ := json.Marshal(transferResult)
		err := imdb.CreateErrorLog(remark, in.OperationID, merOrderID, transferResult.ErrorMsg, transferResult.ErrorCode, string(co))
		if err != nil {
			log.Error(in.OperationID, "创建错误日志失败", err)
		}
		res = "第三方操作转账失败：" + transferResult.ErrorMsg + "，操作ID：" + h.OperateID
		return res, "", nil
	}
	// ====================================================== 创建红包信息 ==========================================

	redPacket := &db.FPacket{
		PacketID:             h.createRedpacketID(in.PacketType, in.UserId),
		UserID:               in.UserId,
		SubmitTime:           transferResult.OrderDate,
		UserRedpacketAccount: fncount.PacketAccountId, // 发红包的用户 红包账户
		PacketType:           in.PacketType,
		IsLucky:              in.IsLucky,
		ExclusiveUserID:      in.ExclusiveUserID,
		PacketTitle:          in.PacketTitle,
		Amount:               in.Amount,
		Number:               in.Number,
		TotalAmount:          tAmount, // 红包总金额 后端自己计算的
		MerOrderID:           h.merOrderID,
		OperateID:            h.OperateID,
		SendType:             in.SendType,
		BindCardAgrNo:        in.BindCardAgrNo,
		RecvID:               in.RecvID, // 接收ID
		Remain:               int64(in.Number),
		RemainAmout:          tAmount,
		ExpireTime:           time.Now().Unix() + 60*60*24,
		CreatedTime:          time.Now().Unix(),
		UpdatedTime:          time.Now().Unix(),
		Status:               1, // 此时红包已经被创建，所以直接赋值就好
		IsExclusive:          in.IsExclusive,
	}
	err = imdb.RedPacketCreateData(redPacket)
	if err != nil {
		//todo  这里会出现扣钱但是发不出红包的问题，需要人工介入
		return "", "", errors.Wrap(err, "创建红包信息错误")
	}

	//增加账户变更日志
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error(in.OperationID, "增加账户变更日志失败:Panic", zap.Any("err", err))
			}
		}()
		err = AddNcountTradeLog(BusinessTypeBalanceSendPacket, int32(in.Amount), in.UserId, fncount.MainAccountId, merOrderID, transferResult.NcountOrderId, redPacket.PacketID)
		if err != nil {
			log.Error(in.OperationID, "增加账户变更日志失败", zap.Error(err))
		}
	}()

	return res, redPacket.PacketID, nil
}

// 银行卡转账
func (h *handlerSendRedPacket) bankTransfer(redPacketID string, in *pb.SendRedPacketReq, merOrderID string) (string, error) {
	//银行卡充值到红包账户
	ncountID, err := BankCardRechargePacketAccount(in.UserId, in.BindCardAgrNo, int32(in.Amount), redPacketID, merOrderID)
	if err != nil {
		return "", err
	}
	return ncountID, nil
	//// 如果转账成功，需要将红包状态修改为发送成功
	//err = imdb.UpdateRedPacketStatus(redPacketID, 1 /* 发送成功 */)
	//if err != nil {
	//	// todo 记录到死信队列中，后续监控处理， 如果转账成功，但是修改红包状态失败，需要人工介入
	//	log.Error(in.OperationID, zap.Error(err))
	//	return "", errors.Wrap(err, "修改红包状态失败 1")
	//}
	//return "", nil
}

// ===================================================================== 红包回调，红包发送成功后调用这个方法=================================

// 当用户发布红包发送成功的时候，调用这个回调函数进行发布红包的后续处理
func HandleSendPacketResult(redPacketID, OperateID string) error {
	//1. 查询红包信息
	redpacketInfo, err := imdb.GetRedPacketInfo(redPacketID)
	if err != nil {
		log.Error(OperateID, "get red packet info error", zap.Error(err))
		return err
	}
	if redpacketInfo == nil {
		log.Error(OperateID, "red packet info is nil")
		return errors.New("red packet info is nil")
	}

	// 2. 生成红包
	if redpacketInfo.PacketType == 2 || redpacketInfo.IsExclusive != 1 {
		// 群红包
		err = GroupPacket(redpacketInfo, redPacketID)
		if err != nil {
			return err
		}
	}

	if redpacketInfo.Status != 1 {
		// 3. 修改红包状态
		err = imdb.UpdateRedPacketStatus(redPacketID, imdb.RedPacketStatusNormal)
		if err != nil {
			log.Error(OperateID, "update red packet status error", zap.Error(err))
			return err
		}
	}

	// todo 发送红包消息
	freq := &contrive_msg.FPacket{
		PacketID:        redPacketID,
		UserID:          redpacketInfo.UserID,
		PacketType:      redpacketInfo.PacketType,
		IsLucky:         redpacketInfo.IsLucky,
		ExclusiveUserID: redpacketInfo.ExclusiveUserID,
		PacketTitle:     redpacketInfo.PacketTitle,
		Amount:          redpacketInfo.Amount,
		Number:          redpacketInfo.Number,
		ExpireTime:      redpacketInfo.ExpireTime,
		MerOrderID:      redpacketInfo.MerOrderID,
		OperateID:       redpacketInfo.OperateID,
		RecvID:          redpacketInfo.RecvID,
		CreatedTime:     redpacketInfo.CreatedTime,
		UpdatedTime:     redpacketInfo.UpdatedTime,
		IsExclusive:     redpacketInfo.IsExclusive,
	}

	// 发送红包消息
	return contrive_msg.SendSendRedPacket(freq, int(redpacketInfo.PacketType))
}

// 给群发的红包
func GroupPacket(req *db.FPacket, redpacketID string) error {

	// 在群红包这里，
	var err error
	if req.IsLucky == 1 {
		// 如果说是手气红包，分散放入红包池
		err = spareRedPacket(redpacketID, int(req.TotalAmount), int(req.Number))
	} else {
		// 这里是平均红包
		err = spareEqualRedPacket(redpacketID, int(req.Amount), int(req.Number))
	}
	if err != nil {
		log.Error(req.OperateID, zap.Error(err))
		return err
	}
	return err
}

// 将红包放入红包池
func spareRedPacket(packetID string, amount, number int) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("spareRedPacket panic", zap.Any("err", err))
		}
	}()
	// 将发送的红包进行计算 : 注意  amount > number
	if amount < number {
		return errors.New(" amount 必须大于 number")
	}
	result := redpacket.GetRedPacket(amount, number)
	err := commonDB.DB.SetRedPacket(packetID, result)
	if err != nil {
		return err
	}
	return nil
}

// amount = 3 ,number =3
func spareEqualRedPacket(packetID string, amount, number int) error {
	result := []int{}
	for i := 0; i < number; i++ {
		result = append(result, amount)
	}
	fmt.Println("\n", "这里是收到的result", result, "\n")
	// 将发送的红包进行计算
	err := commonDB.DB.SetRedPacket(packetID, result)
	if err != nil {
		return err
	}
	return nil
}

// 银行卡充值到红包账户
func BankCardRechargePacketAccount(userId, bindCardAgrNo string, amount int32, packetID, merOrderId string) (string, error) {
	//获取用户账户信息
	accountInfo, err := imdb.GetNcountAccountByUserId(userId)
	if err != nil || accountInfo.Id <= 0 {
		return "", errors.New("账户信息不存在")
	}

	accountResp, err := ncount.NewCounter().QuickPayOrder(&ncount.QuickPayOrderReq{
		MerOrderId: merOrderId,
		QuickPayMsgCipher: ncount.QuickPayMsgCipher{
			PayType:       "3", //绑卡协议号充值
			TranAmount:    cast.ToString(cast.ToFloat64(amount) / 100),
			NotifyUrl:     config.Config.Ncount.Notify.RechargeNotifyUrl,
			BindCardAgrNo: bindCardAgrNo,
			ReceiveUserId: accountInfo.PacketAccountId, //收款账户
			UserId:        accountInfo.MainAccountId,
			SubMerchantId: ncount.SUB_MERCHANT_ID, // 子商户编号
		}})
	if err != nil {
		return "", errors.New(fmt.Sprintf("充值失败(%s)", err.Error()))
	} else {
		if accountResp.ResultCode != "0000" {
			return "", errors.New(fmt.Sprintf("充值失败 (%s,%s)", accountResp.ErrorCode, accountResp.ErrorMsg))
		}
	}

	//增加账户变更日志
	err = AddNcountTradeLog(BusinessTypeBankcardSendPacket, amount, userId, accountInfo.MainAccountId, merOrderId, accountResp.NcountOrderId, packetID)
	if err != nil {
		return "", errors.New(fmt.Sprintf("增加账户变更日志失败(%s)", err.Error()))
	}

	return accountResp.NcountOrderId, nil
}

// 红包确认接口
func (rpc *CloudWalletServer) SendRedPacketConfirm(ctx context.Context, req *pb.SendRedPacketConfirmReq) (*pb.SendRedPacketConfirmResp, error) {
	var (
		resp = &pb.SendRedPacketConfirmResp{
			CommonResp: &pb.CommonResp{
				ErrCode: 0,
				ErrMsg:  "确认成功",
			},
		}
	)
	// 获取红包信息
	redpacketInfo, err := imdb.GetRedPacketInfo(req.RedPacketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.CommonResp.ErrCode = 400
			resp.CommonResp.ErrMsg = "红包ID不存在"
			return resp, nil
		}
		return nil, err
	}

	// 发送红包确认
	nc := NewNcountPay()
	payresult := nc.payComfirm(redpacketInfo.NcountOrderID, req.Code)
	if payresult.ErrCode != 0 {
		resp.CommonResp.ErrCode = 400
		resp.CommonResp.ErrMsg = "Notice:" + payresult.ErrMsg
		return resp, nil
	}

	redpacketInfo.Remark = "银行卡支付已经确认，等待第三方回调"

	// 更新红包信息
	err = imdb.UpdateRedPacketInfo(redpacketInfo.PacketID, redpacketInfo)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 红包领取明细
func (rpc *CloudWalletServer) RedPacketReceiveDetail(_ context.Context, req *pb.RedPacketReceiveDetailReq) (*pb.RedPacketReceiveDetailResp, error) {
	//查询时间转换
	sTime, _ := time.ParseInLocation("2006-01-02", req.StartTime, time.Local)
	eTime, _ := time.ParseInLocation("2006-01-02", req.EndTime, time.Local)

	//获取列表数据
	list, _ := imdb.FindReceiveRedPacketList(req.UserId, sTime.Unix(), eTime.Unix()+86399)

	receiveList := make([]*pb.RedPacketReceiveDetail, 0)
	for _, v := range list {
		receiveList = append(receiveList, &pb.RedPacketReceiveDetail{
			PacketId:    v.PacketId,
			Amount:      v.Amount,
			PacketTitle: v.PacketTitle,
			ReceiveTime: time.Unix(v.ReceiveTime, 0).Format("2006-01-02 15:04:05"),
			PacketType:  v.PacketType,
			IsLucky:     v.IsLucky,
		})
	}

	return &pb.RedPacketReceiveDetailResp{
		RedPacketReceiveDetail: receiveList,
	}, nil
}

// 红包详情
func (rpc *CloudWalletServer) RedPacketInfo(_ context.Context, req *pb.RedPacketInfoReq) (*pb.RedPacketInfoResp, error) {
	//获取红包记录
	redPacketInfo, err := imdb.GetRedPacketInfo(req.PacketId)
	if err != nil {
		return nil, errors.New("红包信息不存在")
	}

	//补充发红包人的用户信息
	nickname, faceUrl := "", ""
	userInfo, err := im_mysql_model.GetUserByUserID(redPacketInfo.UserID)
	if err == nil {
		nickname = userInfo.Nickname
		faceUrl = userInfo.FaceURL
	}

	info := &pb.RedPacketInfoResp{
		UserId:          redPacketInfo.UserID,
		PacketType:      redPacketInfo.PacketType,
		IsLucky:         redPacketInfo.IsLucky,
		IsExclusive:     redPacketInfo.IsExclusive,
		ExclusiveUserID: redPacketInfo.ExclusiveUserID,
		PacketTitle:     redPacketInfo.PacketTitle,
		Amount:          redPacketInfo.Amount,
		Number:          redPacketInfo.Number,
		ExpireTime:      redPacketInfo.ExpireTime,
		Remain:          redPacketInfo.Remain,
		Nickname:        nickname,
		FaceUrl:         faceUrl,
		ReceiveDetail:   make([]*pb.ReceiveDetail, 0),
	}

	//获取当前红包领取记录
	receiveList, _ := imdb.ReceiveListByPacketId(req.PacketId)
	for _, v := range receiveList {
		info.ReceiveDetail = append(info.ReceiveDetail, &pb.ReceiveDetail{
			UserId:      v.UserId,
			Amount:      v.Amount,
			Nickname:    v.Nickname,
			FaceUrl:     v.FaceUrl,
			ReceiveTime: time.Unix(v.ReceiveTime, 0).Format("01月02日 15:04"),
		})
	}

	return info, nil
}

// 禁止群抢红包
func (rpc *CloudWalletServer) ForbidGroupRedPacket(ctx context.Context, req *pb.ForbidGroupRedPacketReq) (*pb.ForbidGroupRedPacketResp, error) {
	var (
		result = &pb.ForbidGroupRedPacketResp{
			CommonResp: &pb.CommonResp{
				ErrMsg:  "禁止群抢红包成功",
				ErrCode: 0,
			},
		}
	)
	// 查看用户是否为群主
	group, err := imdb2.GetGroupInfoByGroupID(req.GroupId)
	if (err != nil && errors.Is(err, sql.ErrNoRows)) || group.GroupID == "" {
		result.CommonResp.ErrCode = 400
		result.CommonResp.ErrMsg = "群信息不存在"
		return result, nil
	}

	// 如果存在群，且用户不是群主
	if group.CreatorUserID != req.UserId {
		result.CommonResp.ErrCode = 400
		result.CommonResp.ErrMsg = "您不是群主"
		return result, nil
	}

	// 禁止抢红包
	err = imdb2.UpdateGroupIsAllowRedPacket(req.GroupId, req.Forbid)
	if err != nil {
		log.Error(req.OperationID, "禁止群抢红包失败", err)
		return nil, err
	}

	// 如果ok 删除
	err = rocksCache.DelGroupInfoFromCache(req.GroupId)
	if err != nil {
		log.Error(req.OperationID, "删除群缓存", err)
		return nil, err
	}

	return result, nil
}

// 获取版本
func (rpc *CloudWalletServer) GetVersion(in context.Context, req *pb.GetVersionReq) (*pb.GetVersionResp, error) {
	var (
		resp = &pb.GetVersionResp{
			Version:       "",
			DownloadUrl:   "",
			UpdateContent: "",
			IsForceUpdate: 0,
			CommonResp: &pb.CommonResp{
				ErrCode: 0,
				ErrMsg:  "获取版本成功",
			},
		}
	)
	// 获取版本信息
	UserVersion, err := imdb.GetFVersion(req.Version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.CommonResp.ErrCode = 400
			resp.CommonResp.ErrMsg = "上报的版本号不存在"
		} else {
			return nil, err
		}
	}

	// 获取最新版本信息
	NewVersion, err := imdb.GetLastedFVersion()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.CommonResp.ErrCode = 400
			resp.CommonResp.ErrMsg = "版本库信息错误"
			return resp, nil
		}
		return nil, err
	}

	// 这里逻辑是如果当前版本远小于最新版本，就强制更新
	if NewVersion.IsForce == 1 {
		resp.Version = NewVersion.VersionCode
		resp.DownloadUrl = NewVersion.DownloadUrl
		resp.UpdateContent = NewVersion.UpdateContent
		resp.IsForceUpdate = NewVersion.IsForce
	}

	// 第二种情况： 如果当前版本落后时间太多，也强制更新 (3个月)
	if NewVersion.CreateTime-UserVersion.CreateTime > 60*60*24*30*3 {
		resp.Version = NewVersion.VersionCode
		resp.DownloadUrl = NewVersion.DownloadUrl
		resp.UpdateContent = NewVersion.UpdateContent
		resp.IsForceUpdate = 1 // 强制更新
	}

	// 第三种情况，如果中间存在多个版本，那么就强制更新 （3个以上的版本）
	return resp, nil
}

// 红包退还情况
func (rpc *CloudWalletServer) RefoundPacket(ctx context.Context, in *pb.RefoundPacketReq) (*pb.RefoundPacketResp, error) {
	var (
		resp = &pb.RefoundPacketResp{
			CommonResp: &pb.CommonResp{
				ErrCode: 0,
				ErrMsg:  "红包退还成功",
			},
		}
	)

	// todo 这里需要注意的是：红包并发处理问题，所以需要引入互斥锁
	// 查询红包状态为1(正常) ，但是超时时间小于当前的红包
	collection, err := imdb.GetExpiredRedPacketList()
	if err != nil {
		return nil, err
	}
	if len(collection) == 0 {
		return resp, nil
	}
	resp.ExpireList = int32(len(collection)) // 查询到的过期红包总数量

	// 红包实际逻辑
	reback := func(packet *db.FPacket, userAccount *db.FNcountAccount, OperationID string) error {

		// 获取实际的金额
		totalAmount := cast.ToString(cast.ToFloat64(packet.RemainAmout) / 100)
		merID := ncount.GetMerOrderID()
		packet.Status = 2 // 红包退回状态

		// 开始转账
		transerMsg := ncount.TransferReq{
			MerOrderId: merID,
			TransferMsgCipher: ncount.TransferMsgCipher{
				PayUserId:     userAccount.PacketAccountId,
				ReceiveUserId: userAccount.MainAccountId,
				TranAmount:    totalAmount,
			},
		}
		transferResult, err := rpc.count.Transfer(&transerMsg)
		if err != nil {
			log.Error(OperationID, "转账失败", err)
			return err
		}

		// ======================= 新生返回参数=======================
		if transferResult.ResultCode != ncount.ResultCodeSuccess {
			remark := "余额发送红包失败： "
			co, _ := json.Marshal(transferResult)
			err := imdb.CreateErrorLog(remark, in.OperationID, merID, transferResult.ErrorMsg, transferResult.ErrorCode, string(co))
			if err != nil {
				log.Error(in.OperationID, "创建错误日志失败", err)
				return err
			}
			packet.Status = 3 // 红包退回失败
		}
		// 修改红包状态
		err = imdb.UpdateRedPacketInfo(packet.PacketID, packet)
		if err != nil {
			log.Error(OperationID, "修改红包状态失败", err)
			return err
		}
		return nil // 红包退还成功
	}

	// 红包退回消息
	rebackMessage := func(packet *db.FPacket) error {
		return nil
	}

	for _, redPacket := range collection {

		if redPacket.Remain == 0 { // 红包剩余金额为0，且红包过期
			err = imdb.UpdateRedPacketStatus(redPacket.PacketID, 2)
			if err != nil {
				log.Error(in.OperationID, "修改红包状态失败", err)
				resp.RefundFailed++
				continue
			}
		}

		if redPacket.Remain > 0 && redPacket.RemainAmout > 0 { // 红包存在剩余金额，且红包过期
			// 1.调用第三方的转账接口：参照小q
			// 2.修改红包状态
			// 3.发送红包退还消息
			// 获取发送红包的信息
			userAcount, err := rocksCache.GetUserAccountInfoFromCache(redPacket.UserID)
			if err != nil {
				log.Error(in.OperationID, "获取用户信息失败", err)
				continue
			}

			err = reback(redPacket, userAcount, in.OperationID)
			if err != nil {
				log.Error(in.OperationID, "红包退还失败", err)
				resp.RefundFailed++
				continue
			}

			err = rebackMessage(redPacket)
			if err != nil {
				log.Error(in.OperationID, "红包退还失败", err)
				resp.RefundFailed++
				continue
			}

			// 判断此时红包状态
			if redPacket.Status == 2 {
				resp.RefundSuccess++ // 红包退还成功

				// 发送红包退还消息
			} else {
				resp.RefundFailed++
			}
		}

	}
	return nil, nil

}
