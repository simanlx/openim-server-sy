package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"testing"
)

// 测试发送红包消息
func TestSendRedPacketMsg(t *testing.T) {
	type args struct {
		redpacketInfo *db.FPacket
		operationID   string
		clickUserID   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendRedPacketMsg(tt.args.redpacketInfo, tt.args.operationID, tt.args.clickUserID...); (err != nil) != tt.wantErr {
				t.Errorf("SendRedPacketMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
