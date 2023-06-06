package testAPI

import (
	"crazy_server/pkg/base_info/notify"
	"crazy_server/pkg/base_info/redpacket_struct"
	"encoding/json"
	"fmt"
	"testing"
)

const BaseURL = "http://127.0.0.1:10002"
const BaseURL2 = "http://127.0.0.1:10002"

func TestPostAPI(t *testing.T) {
	type args struct {
		url       string
		construct func() []byte
	}
	tests := []struct {
		name       string
		args       args
		httpCode   int
		resultCode int
	}{
		{
			name: "提现回调接口:参数是不存在的:走成功逻辑",
			args: args{
				url: BaseURL + "/cloudWallet/charge_account_callback",
				construct: func() []byte {
					req := notify.ChargeNotifyReq{
						Version:         "1.0.0",
						TranCode:        "1001",
						MerOrderId:      "1234567890",
						MerId:           "1234567890",
						MerAttach:       "",
						Charset:         "UTF-8",
						SignType:        "MD5",
						ResultCode:      "",
						ErrorCode:       "",
						ErrorMsg:        "",
						OrderId:         "10086",
						TranAmount:      "",
						SubmitTime:      "",
						TranFinishTime:  "",
						BusinessType:    "",
						FeeAmount:       "",
						BankOrderId:     "",
						RealBankOrderId: "",
						DivideAcctDtl:   "",
						SignValue:       "",
					}
					content, err := json.Marshal(req)
					if err != nil {
						panic(err)
					}
					return content
				},
			},
			httpCode:   200,
			resultCode: 200,
		},
		{
			name: "转账失败回调：走失败回调逻辑",
			args: args{
				url: BaseURL + "/cloudWallet/charge_account_callback",
				construct: func() []byte {
					req := notify.ChargeNotifyReq{
						Version:    "1.0.0",
						TranCode:   "1001",
						MerOrderId: "1234567890",
						MerId:      "1234567890",
						MerAttach:  "",
						Charset:    "UTF-8",
						SignType:   "MD5",
						ResultCode: "",
						ErrorCode:  "4444",
						ErrorMsg:   "",
					}
					content, err := json.Marshal(req)
					if err != nil {
						panic(err)
					}
					return content
				},
			},
			httpCode:   200,
			resultCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpcode, errcode, err := PostAPI(tt.args.url, tt.args.construct)
			if err != nil {
				t.Errorf("PostAPI() error = %v", err)
				return
			}
			if httpcode != tt.httpCode {
				t.Errorf("PostAPI() httpcode = %v, want %v", httpcode, tt.httpCode)
			}
			if errcode != tt.resultCode {
				t.Errorf("PostAPI() errcode = %v, want %v", errcode, tt.resultCode)
			}
		})
	}
}

// 测试发送红包接口
func TestPostAPI_CloudWallet_SendRedPacket(t *testing.T) {
	type args struct {
		url       string
		construct func() []byte
	}
	tests := []struct {
		name       string
		args       args
		httpCode   int
		resultCode int
	}{
		{
			name: "测试发送红包接口",
			args: args{
				url: BaseURL + "/cloudWallet/send_red_packet",
				construct: func() []byte {
					req := &redpacket_struct.SendRedPacket{
						PacketType:      1,             //红包类型(1个人红包、2群红包)
						IsLucky:         0,             //是否为拼手气红包
						IsExclusive:     0,             //是否为专属红包(0不是、1是)
						ExclusiveUserID: "0",           //专属红包接收者 和isExclusive
						PacketTitle:     "新年快乐",        //红包标题
						Amount:          1,             //红包金额 单位：元
						Number:          10,            //红包个数
						SendType:        1,             //发送方式(1钱包余额、2银行卡)
						BankCardID:      1,             //银行卡id
						OperationID:     "11111111111", //链路跟踪id
						RecvID:          "10081",
					}
					content, err := json.Marshal(req)
					if err != nil {
						panic(err)
					}
					fmt.Println(string(content))
					return content
				},
			},
			httpCode:   200,
			resultCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpcode, errcode, err := PostAPI(tt.args.url, tt.args.construct)
			if err != nil {
				t.Errorf("PostAPI() error = %v", err)
				return
			}
			if httpcode != tt.httpCode {
				t.Errorf("PostAPI() httpcode = %v, want %v", httpcode, tt.httpCode)
			}
			if errcode != tt.resultCode {
				t.Errorf("PostAPI() errcode = %v, want %v", errcode, tt.resultCode)
			}
		})
	}
}

// 测试发送自定义消息接口
func TestPostAPI_SendMessage(t *testing.T) {
	type args struct {
		url       string
		construct func() []byte
	}
	tests := []struct {
		name       string
		args       args
		httpCode   int
		resultCode int
	}{
		{
			name: "发送自定义消息结构",
			args: args{
				url: BaseURL + "/cloudWallet/send_red_packet",
				construct: func() []byte {
					req := &redpacket_struct.SendRedPacket{
						PacketType:      1,             //红包类型(1个人红包、2群红包)
						IsLucky:         0,             //是否为拼手气红包
						IsExclusive:     0,             //是否为专属红包(0不是、1是)
						ExclusiveUserID: "0",           //专属红包接收者 和isExclusive
						PacketTitle:     "新年快乐",        //红包标题
						Amount:          1,             //红包金额 单位：元
						Number:          10,            //红包个数
						SendType:        1,             //发送方式(1钱包余额、2银行卡)
						BankCardID:      1,             //银行卡id
						OperationID:     "11111111111", //链路跟踪id
						RecvID:          "10081",
					}
					content, err := json.Marshal(req)
					if err != nil {
						panic(err)
					}
					fmt.Println(string(content))
					return content
				},
			},
			httpCode:   200,
			resultCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpcode, errcode, err := PostAPI(tt.args.url, tt.args.construct)
			if err != nil {
				t.Errorf("PostAPI() error = %v", err)
				return
			}
			if httpcode != tt.httpCode {
				t.Errorf("PostAPI() httpcode = %v, want %v", httpcode, tt.httpCode)
			}
			if errcode != tt.resultCode {
				t.Errorf("PostAPI() errcode = %v, want %v", errcode, tt.resultCode)
			}
		})
	}
}
