package ncount

import (
	"github.com/pkg/errors"
)

// 提现接口

/*
tranAmount 支付金额 1-10 格式：数字（以元为单位） 不可 例如：100.21
userId 提现用户 编号 12 商户用户在新账通平台 的唯一标识 不可 例 如 ： 1
bindCardAgrNo 绑卡协 议号 30 当业务类型为 08 时，必 填，业务类型为 09 时为 空 可空
notifyUrl 异步通知 地址 1-255 交易完成后，异步通知商 户地址
paymentTerminal Info 付款方终 端信息 0-100 付款方交易终端设备类 型、编码等 注：包含两个子域（终端 设备类型，终端设备编 码），子域值间用“|”隔 开，格式如下：子域值 1|子域值 2 详情参见 6.9.1 终端及设 备信息说明 不可
deviceInfo 设备信息 0-200 绑卡设备信息，如：IP、 MAC 、 IMEI 、 IMSI 、 ICCID、WifiMAC、GPS 注：包含七个子域（交易 设 备 IP 、 交 易 设 备 MAC、交易设备 IMEI、 交易设备 IMSI、交易设 备 ICCD 、 交 易 设 备 WIFI MAC、交易设备 GPS），子域值间用“|” 隔开，格式如下：子域 值 1|子域值 2|子域值 3| 子域值 4|子域值 5|子域 值 6|子域值 7 详情参见 6.9.3 设备信息 不可
*/
type WithdrawMsgCipher struct {
	BusinessType    string  `json:"businessType" binding:"required"`    // 业务类型
	TranAmount      float32 `json:"tranAmount" binding:"required"`      // 支付金额
	UserId          string  `json:"userId" binding:"required"`          // 提现用户编号
	BindCardAgrNo   string  `json:"bindCardAgrNo" binding:"required"`   // 绑卡协议号
	NotifyUrl       string  `json:"notifyUrl" binding:"required"`       // 异步通知地址
	PaymentTerminal string  `json:"paymentTerminal" binding:"required"` // 付款方终端信息
	DeviceInfo      string  `json:"deviceInfo" binding:"required"`      // 设备信息
}

func (w *WithdrawMsgCipher) Valid() error {
	if w.TranAmount <= 0 {
		return errors.New("支付金额不能为空")
	}
	if w.UserId == "" {
		return errors.New("提现用户编号不能为空")
	}
	if w.BindCardAgrNo == "" {
		return errors.New("绑卡协议号不能为空")
	}
	if w.NotifyUrl == "" {
		return errors.New("异步通知地址不能为空")
	}
	if w.PaymentTerminal == "" {
		return errors.New("付款方终端信息不能为空")
	}
	if w.DeviceInfo == "" {
		return errors.New("设备信息不能为空")
	}
	return nil
}

type WithdrawReq struct {
	MerOrderID string            `json:"merOrderId" binding:"required"` // 商户订单号
	MsgCipher  WithdrawMsgCipher `json:"msgCipher" binding:"required"`  // 消息密文
}

func (w *WithdrawReq) Valid() error {
	if w.MerOrderID == "" {
		return errors.New("商户订单号不能为空")
	}
	return w.MsgCipher.Valid()
}

/*
resultCode 处理结果码 4 附录二 resultCode 9999
errorCode 异常代码 1-10 附录一 errorCode
errorMsg 异常描述 1-200 中文、字母、数字
ncountOrderId 订单号 32 新账通平台订单号
orderDate 平台订单日期 8 YYYYMMDD signValue 签名字符串 将报文信息用
serviceAmount 服务费
*/
type WithdrawResp struct {
	*BaseReturnParam
	ResultCode    string `json:"resultCode"`    // 处理结果码
	ErrorCode     string `json:"errorCode"`     // 异常代码
	ErrorMsg      string `json:"errorMsg"`      // 异常描述
	NcountOrderId string `json:"ncountOrderId"` // 订单号
	OrderDate     string `json:"orderDate"`     // 平台订单日期
	SignValue     string `json:"signValue"`     // 签名字符串
	ServiceAmount string `json:"serviceAmount"` // 服务费
}
