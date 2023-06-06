package ncount

import (
	"github.com/pkg/errors"
)

/*
serialId 请求序列号 32 商户系统本次提交请求的 序列号， 每次提交必须唯 一。 建议使用本公司代号 加当前时间流水的方式。 例如:爱泉网络科技公司于 2011 年 1 月 1 日 12 时 30 分 50 秒，提交请求 值:aqpay201101011230 50 不可 例 如 ： aqpay2011 010112305 0
orgMerOrderId 商户原始 订单号 32 商户系统原始支付请求的 商户订单号 不可
orgSubmitTime 原订单支 付下单请 求时间 14 原订单支付下单时间，格 式 :YYYYMMDDHHMMS S 不可
orderAmount 原订单金 1-10 商户原订单交易金额，以元 不可
*/
type RefundMsgCipher struct {
	SerialId      string `json:"serialId" binding:"required"`
	OrgMerOrderId string `json:"orgMerOrderId" binding:"required"`
	OrgSubmitTime string `json:"orgSubmitTime" binding:"required"`
	OrderAmount   string `json:"orderAmount" binding:"required"`
	/*
		refundSource 退款资金 来源 1 退款资金来源 1. 新账通原路退款 2. 商户现金户 不可 目前只支持 1
		destType 退款目的 地类型 1 退款目的地的类型 1:原路退回 不可 目前只支持 1
		refundType 退款类型 1 退款类型 1:全额退款 2：部分退款 不可 03 消费类型 支持部分退， 04 担保下单 内扣不支持部 分退，担保下 单外扣暂时只 支持微信公众 号及 H5 交易 部分退款
		refundAmount 商户退款 金额 1-10 格式：数字（以元为单位） 不可
		notifyUrl 异步通知 地址 1-255 交易完成后，异步通知商户 地址 不可 例 如 ： https://www .xxx.com/res ponse.do
		remark 备注 50 备注信息 不可
		divideRefundDtl 担保下单 退款明细 担保下单外扣部分退款时 必填 可空 格 式 ： [{\"ledgerU serId\":\"1 23\",\"amo unt\":\"30\ "}]
	*/
	RefundSource    string `json:"refundSource" binding:"required"`
	DestType        string `json:"destType" binding:"required"`
	RefundType      string `json:"refundType" binding:"required"`
	RefundAmount    string `json:"refundAmount" binding:"required"`
	NotifyUrl       string `json:"notifyUrl" binding:"required"`
	Remark          string `json:"remark" binding:"required"`
	DivideRefundDtl string `json:"divideRefundDtl" binding:"required"`
}

func (r *RefundMsgCipher) Valid() error {
	if r.SerialId == "" {
		return errors.New("serialId is empty")
	}
	if r.OrgMerOrderId == "" {
		return errors.New("orgMerOrderId is empty")
	}
	if r.OrgSubmitTime == "" {
		return errors.New("orgSubmitTime is empty")
	}
	if r.OrderAmount == "" {
		return errors.New("orderAmount is empty")
	}

	if r.RefundSource == "" {
		return errors.New("refundSource is empty")
	}
	if r.DestType == "" {
		return errors.New("destType is empty")
	}
	if r.RefundType == "" {
		return errors.New("refundType is empty")
	}
	if r.RefundAmount == "" {
		return errors.New("refundAmount is empty")
	}
	if r.NotifyUrl == "" {
		return errors.New("notifyUrl is empty")
	}
	if r.Remark == "" {
		return errors.New("remark is empty")
	}
	if r.DivideRefundDtl == "" {
		return errors.New("divideRefundDtl is empty")
	}
	return nil
}

type RefundReq struct {
	MerOrderId      string `json:"merOrderId" binding:"required"`
	RefundMsgCipher RefundMsgCipher
}

func (t *RefundReq) Valid() error {
	if t.MerOrderId == "" {
		return errors.New("merOrderId is empty")
	}
	return t.RefundMsgCipher.Valid()
}

/*
resultCode 处理结果码 4 附录二 resultCode 0000
errorCode 异常代码 1-10 附录一 errorCode
errorMsg 异常描述 1-200 中文、字母、数字
orgMerOrderId 商户系统原始支付 请求的商户订单号 同上送
refundAmount 商户退款金额 同上送
orderAmount 原订单金额 同上送
tranFinishTime 交易完成时间 格式：YYYYMMDDHHMMSS
ncountOrderId 新账通退款流水号 新账通平台生成的退款订单号， 新账通平台受理商户退款订单 请求失败时，无退款订单号
remark 扩展字段 同上送
signValue 签名字符串 将报文信息用signType 域设置 的方式签名后生成的字符串
*/
type RefundResp struct {
	BaseReturnParam
	ResultCode     string `json:"resultCode" binding:"required"`
	ErrorCode      string `json:"errorCode" binding:"required"`
	ErrorMsg       string `json:"errorMsg" binding:"required"`
	OrgMerOrderId  string `json:"orgMerOrderId" binding:"required"`
	RefundAmount   string `json:"refundAmount" binding:"required"`
	OrderAmount    string `json:"orderAmount" binding:"required"`
	TranFinishTime string `json:"tranFinishTime" binding:"required"`
	NcountOrderId  string `json:"ncountOrderId" binding:"required"`
	Remark         string `json:"remark" binding:"required"`
	SignValue      string `json:"signValue" binding:"required"`
}
