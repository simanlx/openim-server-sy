package ncount

import (
	"errors"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

// ============================================================请求头部参数===========================
// 商户请求参数列表（POST）
type NAccountBaseParam struct {
	Version    string `json:"version" binding:"required"`
	TranCode   string `json:"tranCode" binding:"required"`
	MerId      string `json:"merId" binding:"required"`
	MerOrderId string `json:"merOrderId" binding:"required"`
	SubmitTime string `json:"submitTime" binding:"required"`
	MsgCipher  string `json:"msgCipher" binding:"required"`
	SignType   string `json:"signType" binding:"required"`
	SignValue  string `json:"signValue" binding:"required"`
	Charset    string `json:"charset" binding:"required"`
}

func NewNAccountBaseParam(merOrderID, msgCipher, tranCode string) *NAccountBaseParam {
	tim := time.Now()
	times := tim.Format("20060102150405")
	return &NAccountBaseParam{
		Version:    "1.0",
		TranCode:   tranCode,
		SignType:   "1",
		Charset:    "1",
		SubmitTime: times,
		MerId:      MER_USER_ID,
		MsgCipher:  msgCipher,
		MerOrderId: merOrderID,
	}
}

// str : version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
// signValue = version
func (n *NAccountBaseParam) flushSignValue() (error, string) {
	if n.Version == "" {
		return errors.New("version is empty"), ""
	}
	if n.TranCode == "" {
		return errors.New("tranCode is empty"), ""
	}
	if n.MerId == "" {
		return errors.New("merId is empty"), ""
	}
	if n.MerOrderId == "" {
		return errors.New("merOrderId is empty"), ""
	}
	if n.SubmitTime == "" {
		return errors.New("submitTime is empty"), ""
	}
	if n.MsgCipher == "" {
		return errors.New("msgCiphert is empty"), ""
	}
	if n.SignType == "" {
		return errors.New("signType is empty"), ""
	}
	var str = ""
	str += "version=[" + n.Version + "]"
	str += "tranCode=[" + n.TranCode + "]"
	str += "merId=[" + n.MerId + "]"
	str += "merOrderId=[" + n.MerOrderId + "]"
	str += "submitTime=[" + n.SubmitTime + "]"
	str += "msgCiphertext=[" + n.MsgCipher + "]"
	str += "signType=[" + n.SignType + "]"
	return nil, str
}

func (n *NAccountBaseParam) Form() url.Values {
	form := url.Values{}
	form.Add("version", n.Version)
	form.Add("tranCode", n.TranCode)
	form.Add("merId", n.MerId)
	form.Add("merOrderId", n.MerOrderId)
	form.Add("submitTime", n.SubmitTime)
	form.Add("msgCiphertext", n.MsgCipher)
	form.Add("signType", n.SignType)
	form.Add("signValue", n.SignValue)
	form.Add("charset", "1")
	return form
}

// ============================================================银行卡================================

/*
	indCardAgrNo 绑 卡 协 议 号 30
	bankCode 银行简码 详情参见附录三 银行简码 例如：ICBC
	cardNo 卡号掩码 1-30 格式：数字
*/

type NAccountBankCard struct {
	BindCardAgrNo string `json:"bindCardAgrNo" binding:"required"`
	BankCode      string `json:"bankCode" binding:"required"`
	CardNo        string `json:"cardNo" binding:"required"`
}

// ============================================================返回头部参数===========================
type BaseReturnParam struct {
	Version    string `json:"version"`
	TranCode   string `json:"tranCode"`
	MerOrderId string `json:"merOrderId"`
	MerId      string `json:"merId"`
	MerAttach  string `json:"merAttach"`
	Charset    string `json:"charset"`
	SignType   string `json:"signType"`
}

// ============================================================退款回调接口的内容===========================

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
refundInstalme ntAmt 当前退款订单应退 回的商户分期贴息 费用 0-12 注：仅在原交易支付请求下单时 分期期数不为空并且退款成功 后的异步通知报文中才返回 格式：数字（以元为单位）
*/

type NAccountRefundCallBack struct {
	BaseReturnParam
	ResultCode       string `json:"resultCode"`
	ErrorCode        string `json:"errorCode"`
	ErrorMsg         string `json:"errorMsg"`
	OrgMerOrderId    string `json:"orgMerOrderId"`
	RefundAmount     string `json:"refundAmount"`
	OrderAmount      string `json:"orderAmount"`
	TranFinishTime   string `json:"tranFinishTime"`
	NcountOrderId    string `json:"ncountOrderId"`
	Remark           string `json:"remark"`
	SingValue        string `json:"signValue"`
	RefundInstalment string `json:"refundInstalmentAmt"`
}

/*
resultCode 交易结果 4 详情参见 6.2 附录二 resultCode 0000
errorCode 异常代码 10 详情参见 6.1 附录一 errorCode
errorMsg 异常描述 200 中文、字母、数字
ncountOrderId 新账通订单 32 商户支付流程订单凭证
*/

type NQuickPayCallBack struct {
	BaseReturnParam
	ResultCode    string `json:"resultCode"`
	ErrorCode     string `json:"errorCode"`
	ErrorMsg      string `json:"errorMsg"`
	NcountOrderId string `json:"ncountOrderId"`
}

func GetMerOrderID() string {
	// 生成一串随机数
	// 时间戳 + 6位随机数
	tim := time.Now()
	times := tim.Format("20060102150405")
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(999999)
	merOrderID := times + strconv.Itoa(randNum)
	return merOrderID
}
