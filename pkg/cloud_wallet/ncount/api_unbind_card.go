package ncount

import (
	"github.com/pkg/errors"
)

/*
oriBindCardAgrN o 原 绑 卡 协议号 30 商户绑卡确认后获得的 绑卡协议号 不可 例 如 ： 11000000111
userId 用户 ID 1-32 格式：数字，字母，下 划线，竖划线，中划线 不可 例如：102121
*/
type UnBindCardMsgCipher struct {
	OriBindCardAgrNo string `json:"oriBindCardAgrNo" binding:"required"`
	UserId           string `json:"userId" binding:"required"`
}

func (u *UnBindCardMsgCipher) Valid() error {
	if u.OriBindCardAgrNo == "" {
		return errors.New("oriBindCardAgrN is empty")
	}
	if u.UserId == "" {
		return errors.New("userId is empty")
	}
	return nil
}

type UnBindCardReq struct {
	MerOrderId          string `json:"merOrderId" binding:"required"`
	UnBindCardMsgCipher UnBindCardMsgCipher
}

func (u *UnBindCardReq) Valid() error {
	if u.MerOrderId == "" {
		return errors.New("merOrderId is empty")
	}
	if err := u.UnBindCardMsgCipher.Valid(); err != nil {
		return err
	}
	return nil
}

/*
version 版本号 同上送
tranCode 交易代码 同上送
merOrderId 商户订单号 同上送
merId 商户 ID 同上送
merAttach 附加数据 同上送
charset 编码方式 同上送
signType 签名类型 同上送
resultCode 处理结果码 4 详情参见附录二resultCode 9999
errorCode 异常代码 1-10 详情参见附录一 errorCode
errorMsg 异常描述 1-200 中文、字母、数字
signValue 签名字符串 将报文信息用
*/

type UnBindCardResp struct {
	Version    string `json:"version" binding:"required"`
	TranCode   string `json:"tranCode" binding:"required"`
	MerOrderId string `json:"merOrderId" binding:"required"`
	MerId      string `json:"merId" binding:"required"`
	MerAttach  string `json:"merAttach" binding:"required"`
	Charset    string `json:"charset" binding:"required"`
	SignType   string `json:"signType" binding:"required"`
	ResultCode string `json:"resultCode" binding:"required"`
	ErrorCode  string `json:"errorCode" binding:"required"`
	ErrorMsg   string `json:"errorMsg" binding:"required"`
	SignValue  string `json:"signValue" binding:"required"`
}
