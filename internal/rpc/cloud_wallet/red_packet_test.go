package cloud_wallet

import (
	"testing"
)

// 测试红包回调接口
func TestHandleSendPacketResult(t *testing.T) {
	type args struct {
		redPacketID string
		OperateID   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "一个找不到红包ID的测试",
			args: args{
				redPacketID: "12234532432",
				OperateID:   "1111111111",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HandleSendPacketResult(tt.args.redPacketID, tt.args.OperateID); (err != nil) != tt.wantErr {
				t.Errorf("HandleSendPacketResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
