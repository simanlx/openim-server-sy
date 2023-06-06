package ncount

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNewCounter(t *testing.T) {
	tests := []struct {
		name string
		want NCounter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCounter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 绑定银行卡接口 ： 单元测试通过
func Test_counter_BindCard(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *BindCardReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				notifyQuickPayConfirmURL: "http://www.baidu.com",
				notifyRefundURL:          "http://www.baidu.com",
				notifyWithdrawURL:        "http://www.baidu.com",
			},
			args: args{
				req: &BindCardReq{
					MerOrderId: "afdafa",
					BindCardMsgCipherText: BindCardMsgCipherText{
						CardNo:            "",
						HolderName:        "",
						CardAvailableDate: "",
						Cvv2:              "",
						MobileNo:          "",
						IdentityType:      "1",
						IdentityCode:      "",
						UserId:            "",
					},
				},
			},
			want: "4444",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.BindCard(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BindCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BindCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 绑定确认接口 ：  单元测试通过
func Test_counter_BindCardConfirm(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *BindCardConfirmReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				notifyQuickPayConfirmURL: "http://www.baidu.com",
				notifyRefundURL:          "http://www.baidu.com",
				notifyWithdrawURL:        "http://www.baidu.com",
			},
			args: args{
				req: &BindCardConfirmReq{
					MerOrderId: "afdafa",
					BindCardConfirmMsgCipherText: BindCardConfirmMsgCipherText{
						NcountOrderId: "1234",   // 通联订单号 这个是之前的记录好
						SmsCode:       "123456", // 短信验证码
						MerUserIp:     "1234",   // 用户IP
					},
				},
			},
			want: "4444",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.BindCardConfirm(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BindCardConfirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ResultCode != tt.want {
				t.Errorf("BindCardConfirm() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_CheckUserAccountDetail(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *CheckUserAccountDetailReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CheckUserAccountDetailResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.CheckUserAccountDetail(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUserAccountDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckUserAccountDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_CheckUserAccountInfo(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *CheckUserAccountReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CheckUserAccountResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.CheckUserAccountInfo(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUserAccountInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckUserAccountInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_CheckUserAccountTrans(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *CheckUserAccountTransReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CheckUserAccountTransResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.CheckUserAccountTrans(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUserAccountTrans() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckUserAccountTrans() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 创建用户账户 单元测试通过
func Test_counter_NewAccount(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *NewAccountReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// 由于信息安全原因，身份证不方便留下，所以这里只能用测试的身份证号码，但是肯定会报错，所以这里只是测试一下请求是否成功
		{
			name: "测试非真实用户手机 ： 这个手机号和用户的身份证信息不对应",
			fields: fields{
				notifyQuickPayConfirmURL: "http://www.baidu.com",
				notifyRefundURL:          "http://www.baidu.com",
				notifyWithdrawURL:        "http://www.baidu.com",
			},
			args: args{
				req: &NewAccountReq{
					OrderID: "dsfsafdsa",
					MsgCipherText: &NewAccountMsgCipherText{
						MerUserId: "main_10086",
						Mobile:    "15282603386",
						UserName:  "沈晨曦",
						CertNo:    "5116231185554",
					},
				},
			},
			want: "4444", // 表示请求结果失败
		},
		{
			name: "真实手机用户 : 这个用户已经存在账号，所以会返回失败",
			fields: fields{
				notifyQuickPayConfirmURL: "http://www.baidu.com",
				notifyRefundURL:          "http://www.baidu.com",
				notifyWithdrawURL:        "http://www.baidu.com",
			},
			args: args{
				req: &NewAccountReq{
					OrderID: "ds_" + time.Now().Format("20060102150405"),
					MsgCipherText: &NewAccountMsgCipherText{
						MerUserId: "main_100861ss",
						Mobile:    "18566634004",
						UserName:  "沈晨曦",
						CertNo:    "511185554",
					},
				},
			},
			want: "4444", // 表示请求结果失败
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.NewAccount(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccount() error = %+v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ResultCode != tt.want {
				t.Errorf("NewAccount() got = %+v, want %v", got, tt.want)
			}
		})
	}
}

// 快捷支付接口测试 单元测试通过
func Test_counter_QuickPayOrder(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *QuickPayOrderReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *QuickPayOrderResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.QuickPayOrder(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuickPayOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuickPayOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_Refund(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *RefundReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RefundResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.Refund(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Refund() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Refund() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 转账接口
func Test_counter_Transfer(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *TransferReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "转账",
			fields: fields{
				notifyQuickPayConfirmURL: "http://www.baidu.com",
				notifyRefundURL:          "http://www.baidu.com",
				notifyWithdrawURL:        "http://www.baidu.com",
			},
			args: args{
				req: &TransferReq{
					MerOrderId: "",
					TransferMsgCipher: TransferMsgCipher{
						PayUserId:     "",
						ReceiveUserId: "",
						TranAmount:    "",
						BusinessType:  "", // 业务类型
					},
				},
			},
			want:    "4444",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.Transfer(tt.args.req)
			if err != nil && err != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transfer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_UnbindCard(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *UnBindCardReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UnBindCardResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.UnbindCard(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnbindCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnbindCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_Withdraw(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *WithdrawReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *WithdrawResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.Withdraw(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Withdraw() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 发送自定义红包消息
func Test_counter_SendRedPack(t *testing.T) {

}

// 在这里进行转账
func Test_counter_QuickPayOrderOther(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *NAccountQuickPayOtherOther
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *QuickPayOrderResp
		wantErr bool
	}{
		{
			name: "周荣亮-转账",
			fields: fields{
				notifyQuickPayConfirmURL: "",
				notifyRefundURL:          "",
				notifyWithdrawURL:        "",
			},
			args: args{
				req: &NAccountQuickPayOtherOther{
					TranAmount:    "0.5",
					PayType:       "2",
					CardNo:        "6217710312086812",
					HolderName:    "邹荣亮",
					MobileNo:      "13530210115",
					IdentityType:  "1",
					IdentityCode:  "430426199102086615",
					NotifyUrl:     "www.baidu.com",
					ReceiveUserId: MER_USER_ID,
					SubMerchantId: SUB_MERCHANT_ID,
				},
			},
		},
		{
			name: "声焕-转账",
			fields: fields{
				notifyQuickPayConfirmURL: "",
				notifyRefundURL:          "",
				notifyWithdrawURL:        "",
			},
			args: args{
				req: &NAccountQuickPayOtherOther{
					TranAmount:    "5",
					PayType:       "2",
					CardNo:        "6217680913281569", // 工商
					HolderName:    "黄声焕",
					MobileNo:      "13580779180",
					IdentityType:  "1",
					IdentityCode:  "445121199204134815",
					NotifyUrl:     "https://www.baidu.com",
					ReceiveUserId: MER_USER_ID,
					SubMerchantId: SUB_MERCHANT_ID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, _ := c.QuickPayOrderOther(tt.args.req)
			fmt.Printf("got:%+v\n", got)
		})
	}
}

func Test_counter_QuickPayConfirm(t *testing.T) {
	type fields struct {
		notifyQuickPayConfirmURL string
		notifyRefundURL          string
		notifyWithdrawURL        string
	}
	type args struct {
		req *QuickPayConfirmReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *QuickPayConfirmResp
		wantErr bool
	}{
		{
			name: "周荣亮-转账确认",
			fields: fields{
				notifyQuickPayConfirmURL: "",
				notifyRefundURL:          "",
				notifyWithdrawURL:        "",
			},
			args: args{
				req: &QuickPayConfirmReq{
					MerOrderId: GetMerOrderID(),
					QuickPayConfirmMsgCipher: QuickPayConfirmMsgCipher{
						NcountOrderId:        "2023052610846012", //订单ID
						SmsCode:              "794268",           //code
						PaymentTerminalInfo:  "02|AA01BB",
						ReceiverTerminalInfo: "01|00001|CN|469023",
						DeviceInfo:           "192.168.0.1|E1E2E3E4E5E6|123456789012345|20000|898600MFSSYYGXXXXXXP|H1H2H3H4H5H6|AABBCC",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				notifyQuickPayConfirmURL: tt.fields.notifyQuickPayConfirmURL,
				notifyRefundURL:          tt.fields.notifyRefundURL,
				notifyWithdrawURL:        tt.fields.notifyWithdrawURL,
			}
			got, err := c.QuickPayConfirm(tt.args.req)
			if err != nil {
				panic(err)
			}
			fmt.Printf("got:%+v\n", got)
		})
	}
}
