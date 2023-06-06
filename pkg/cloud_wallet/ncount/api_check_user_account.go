package ncount

import (
	"github.com/pkg/errors"
)

type CheckUserAccountReq struct {
	OrderID string `json:"orderId" binding:"required"` // 商户订单号
	UserID  string `json:"userId" binding:"required"`  // 商户的用户ID
}

func (c *CheckUserAccountReq) Valid() error {
	if c.OrderID == "" {
		return errors.New("orderId is empty")
	}
	if c.UserID == "" {
		return errors.New("userId is empty")
	}
	return nil
}

/*
version 版本号 同上送
tranCode 交易代码 同上送
merOrderId 商 户 订 单 号 同上送
merId 商户 ID 同上送
merAttach 附加数据 同上送
charset 编码方式 同上送
signType 签名类型 同上送
resultCode 处 理 结 果 码 4 格式：数字 详 情 参 见 附 录 二
resultCode 9999
errorCode 异常代码 1-10 格式：数字，字母 详 情 参 见 附 录 一 errorCode
errorMsg 异常描述 1-200 中文、字母、数字
failReason 失败原因 中文、字母、数字
userId 用户编号 12 同上送
outUserId 用户标识 开户时商户
*/
type CheckUserAccountResp struct {
	Version    string `json:"version" binding:"required"`
	TranCode   string `json:"tranCode" binding:"required"`
	MerOrderId string `json:"merOrderId" binding:"required"`
	MerId      string `json:"merId" binding:"required"`
	MerAttach  string `json:"merAttach" binding:"required"`
	Charset    string `json:"charset" binding:"required"`
	SignType   string `json:"signtype" binding:"required"`

	ResultCode string `json:"resultCode" binding:"required"`
	ErrorCode  string `json:"errorCode" binding:"required"`
	ErrorMsg   string `json:"errorMsg" binding:"required"`
	UserId     string `json:"userId" binding:"required"`
	OutUserId  string `json:"outUserId" binding:"required"`

	/*
		userStat 用户状态 2 00：正常 01：待激活 02：锁定 03：开户失败 99：注销
		auditStat 审核状态 2 00：审核通过 01：待审核 02：审核不通过 03：无需审核
		authStat 实名状态 2 00：实名认证成功 01：待认证 02：实名认证失败 03：无需实名认证 04：认证超时
		balAmount 帐户余额 格式：数字，单位：元
		unclearAmount 待 清 算 余 额 格式：数字，单位：元
		unclearSumAmount 待 清 算 余 额汇总 格式：数字，单位：元；平 台用户时出现
		availableBalance 可用余额 格式：数字，单位：元
		unsettleBalance 待 结 转 余 额 格式：数字，单位：元
		bindCardAgrNoLi st 绑 卡 协 议 号列表 例 如 ： ["20180822000000 0123","2018082200 00000118"]
	*/
	UserStat          string `json:"userStat" binding:"required"`
	AuditStat         string `json:"auditStat" binding:"required"`
	AuthStat          string `json:"authStat" binding:"required"`
	BalAmount         string `json:"balAmount" binding:"required"`
	UnclearAmount     string `json:"unclearAmount" binding:"required"`
	UnclearSumAmount  string `json:"unclearSumAmount" binding:"required"`
	AvailableBalance  string `json:"availableBalance" binding:"required"`
	UnsettleBalance   string `json:"unsettleBalance" binding:"required"`
	BindCardAgrNoList string `json:"bindCardAgrNoList" binding:"required"`

	/*
		indCardAgrNo 绑 卡 协 议 号 30
		bankCode 银行简码 详情参见附录三 银行简码 例如：ICBC
		cardNo 卡号掩码 1-30 格式：数字
	*/
	//BindCards interface{} `json:"bindCards" binding:"required"`

	SignValue string `json:"signValue" binding:"required"`
}
