package cloud_wallet

import (
	"context"
	imdb "crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	"crazy_server/pkg/common/log"
	pb "crazy_server/pkg/proto/cloud_wallet"
	"fmt"
	"github.com/pkg/errors"
)

// 银行卡充值回调
func (s *CloudWalletServer) ChargeNotify(_ context.Context, req *pb.ChargeNotifyReq) (*pb.ChargeNotifyResp, error) {
	// 1.检查订单是否存在
	if req.NcountOrderId == "" {
		return nil, errors.New("订单号不能为空")
	}

	if req.ResultCode != "0000" {
		// 这里需要通知用户发送红包失败
	}

	if req.ResultCode == "0000" {
		// 需要发送code到所有群用户
	}

	// 查询记录
	tradeInfo, err := imdb.GetThirdOrderNoRecord(req.NcountOrderId)
	if err != nil {
		return nil, errors.New("交易记录不存在")
	}

	//已处理
	if tradeInfo.NcountStatus == 1 {
		return &pb.ChargeNotifyResp{}, nil
	}

	//校验订单金额

	// 修改订单状态
	err = imdb.FNcountTradeUpdateStatusbyThirdOrderNo(req.NcountOrderId)
	if err != nil {
		log.Error("修改订单状态失败", err, req)
		return nil, err
	}

	//处理红包逻辑
	if tradeInfo.PacketID != "" {
		if err := HandleSendPacketResult(tradeInfo.PacketID, ""); err != nil {
			fmt.Println("HandleSendPacketResult err", err)
			return nil, err
		}
	}

	return &pb.ChargeNotifyResp{}, nil
}

// 提现回调接口
func (s *CloudWalletServer) WithDrawNotify(_ context.Context, req *pb.DrawNotifyReq) (*pb.DrawNotifyResp, error) {
	// 1.检查订单是否存在
	if req.NcountOrderId == "" {
		return nil, errors.New("订单号不能为空")
	}

	// 查询记录
	tradeInfo, err := imdb.GetThirdOrderNoRecord(req.NcountOrderId)
	if err != nil {
		return nil, errors.New("交易记录不存在")
	}

	//已处理
	if tradeInfo.NcountStatus == 1 {
		return &pb.DrawNotifyResp{}, nil
	}

	// 修改订单状态
	err = imdb.FNcountTradeUpdateStatusbyThirdOrderNo(req.NcountOrderId)
	if err != nil {
		log.Error("修改订单状态失败", err, req)
		return nil, err
	}

	return &pb.DrawNotifyResp{}, nil
}
