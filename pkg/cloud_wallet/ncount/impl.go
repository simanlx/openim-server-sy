package ncount

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type counter struct {
	notifyQuickPayConfirmURL string // 快捷支付回调地址
	notifyRefundURL          string // 退款异步通知接口
	notifyWithdrawURL        string // 提现异步通知接口
}

func NewCounter() NCounter {
	return &counter{}
}

// ======================================= 账户 =======================================

// 创建用户新生账户 ：单账户
func (c *counter) NewAccount(req *NewAccountReq) (*NewAccountResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Vaild(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.MsgCipherText)
	if err != nil {
		return nil, err
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := Encrpt(data, PUBLIC_KEY)
	if err != nil {
		return nil, err
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.OrderID, string(cipher), "R010")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(NewAccountURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &NewAccountResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 查询用户账户信息 ：单账户
func (c *counter) CheckUserAccountInfo(req *CheckUserAccountReq) (*CheckUserAccountResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := Encrpt(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.OrderID, string(cipher), "Q001")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(checkUserAccountURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &CheckUserAccountResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// ======================================= 银行卡管理 =======================================

// 绑定用户银行卡
func (c *counter) BindCard(req *BindCardReq) (*BindCardResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.BindCardMsgCipherText)
	fmt.Println(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := RsaEncryptBlock(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "R007")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}

	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}

	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(bindCardURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &BindCardResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 银行卡确认接口
func (c *counter) BindCardConfirm(req *BindCardConfirmReq) (*BindCardConfirmResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.BindCardConfirmMsgCipherText)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := RsaEncryptBlock(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "R008")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(bindCardConfirmURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &BindCardConfirmResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 解绑银行卡接口
func (c *counter) UnbindCard(req *UnBindCardReq) (*UnBindCardResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.UnBindCardMsgCipher)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := RsaEncryptBlock(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}

	fmt.Println("req.UnBindCardMsgCipher", string(data))
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "R009")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(unbindCardURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &UnBindCardResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 获取用户账户信息接口
func (c *counter) CheckUserAccountDetail(req *CheckUserAccountDetailReq) (*CheckUserAccountDetailResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.CheckUserAccountDetailMsgCipher)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := Encrpt(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "Q004")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(checkUserAccountDetailURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &CheckUserAccountDetailResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 查询用户交易详情
func (c *counter) CheckUserAccountTrans(req *CheckUserAccountTransReq) (*CheckUserAccountTransResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.CheckUserAccountTransMsgCipher)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := Encrpt(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "Q002")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(checkUserAccountTransURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &CheckUserAccountTransResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 快捷支付接口 ： 为了充值到后台
func (c *counter) QuickPayOrderOther(req *NAccountQuickPayOtherOther) (*QuickPayOrderResp, error) {
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := RsaEncryptBlock(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(GetMerOrderID(), string(cipher), "T007")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(quickPayOrderURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &QuickPayOrderResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 快捷支付接口
func (c *counter) QuickPayOrder(req *QuickPayOrderReq) (*QuickPayOrderResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.QuickPayMsgCipher)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := RsaEncryptBlock(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "T007")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(quickPayOrderURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &QuickPayOrderResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 快捷支付确认接口
func (c *counter) QuickPayConfirm(req *QuickPayConfirmReq) (*QuickPayConfirmResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.QuickPayConfirmMsgCipher)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := RsaEncryptBlock(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "T008")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(quickPayConfirmURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &QuickPayConfirmResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 转账接口
func (c *counter) Transfer(req *TransferReq) (*TransferResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.TransferMsgCipher)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := Encrpt(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "T003")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(transferURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &TransferResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 退款接口
func (c *counter) Refund(req *RefundReq) (*RefundResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.RefundMsgCipher)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := Encrpt(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderId, string(cipher), "T005")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}
	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(refundURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &RefundResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}

// 提现接口
func (c *counter) Withdraw(req *WithdrawReq) (*WithdrawResp, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}
	if err := req.Valid(); err != nil {
		return nil, errors.Wrap(err, "req.Vaild")
	}
	// 1. 将报文信息转换为 JSON 格式
	data, err := json.Marshal(req.MsgCipher)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	// 2. 将 JSON 格式的报文信息用平台公钥 RSA 加密后 base64 的编码值
	cipher, err := RsaEncryptBlock(data, PUBLIC_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Encrpt")
	}
	// 3.version=[]tranCode=[]merId=[]merOrderId=[]submitTime=[]msgCiphertext=[]signType=[]
	// signValue= version
	// 2. 使用RSA进行私钥签名
	// 3. 签名后的二进制转Base64编码
	body := NewNAccountBaseParam(req.MerOrderID, string(cipher), "T002")
	err, str := body.flushSignValue()
	if err != nil {
		return nil, errors.Wrap(err, "flushSignValue")
	}
	// 4. 使用RSA进行私钥签名
	sign, err := Sign([]byte(str), PRIVATE_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "Sign")
	}

	body.SignValue = sign
	content := body.Form()
	result, err := httpPost(withdrawURL, content)
	if err != nil {
		return nil, errors.Wrap(err, "httpPost")
	}
	reply := &WithdrawResp{}
	err = json.Unmarshal(result, reply)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return reply, nil
}
