package cloud_wallet

import (
	pb "crazy_server/pkg/proto/cloud_wallet"
)

var CodeMessage = map[pb.CloudWalletErrCode]string{

	/*	ServerError = 9999 ; // 9999 服务器错误
		// 1000- 1999 云钱包相关错误
		PacketStatusIsCreate  = 1000; // 红包状态是创建
		PacketStatusIsFinish  = 1001; // 红包状态是完成
		PacketStatusIsExpire  = 1002; // 红包状态是过期
		PacketStatusIsExclusive  = 1003; // 红包状态是独占
		PacketStatusIsReceived  = 1004; // 红包状态是取消*/

	pb.CloudWalletErrCode_ServerError:             "服务器错误",
	pb.CloudWalletErrCode_PacketStatusIsCreate:    "红包状态是创建",
	pb.CloudWalletErrCode_PacketStatusIsFinish:    "红包状态是完成",
	pb.CloudWalletErrCode_PacketStatusIsExpire:    "红包状态是过期",
	pb.CloudWalletErrCode_PacketStatusIsExclusive: "红包是专属红包",
	pb.CloudWalletErrCode_PacketStatusIsReceived:  "红包状态是取消",
}

func CompleteRespMsg(resp *pb.CommonResp) {
	resp.ErrMsg = CodeMessage[resp.ErrCode]
}

func GetErrMsg(code pb.CloudWalletErrCode) string {
	return CodeMessage[code]
}
