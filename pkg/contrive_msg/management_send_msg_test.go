package contrive_msg

import (
	"testing"
)

func TestNewManagementSendMsg_RedMsg(t *testing.T) {
	req := &FPacket{
		PacketID:    "110018168067881370276",
		UserID:      "1914080869",
		PacketType:  1,
		IsLucky:     1,
		PacketTitle: "新年快乐",
		OperateID:   "10000085",
		// steven 1914080869
		RecvID:      "1914080869",
		IsExclusive: 0,
	}
	SendSendRedPacket(req, 1)
}

// 模拟发送一个红包
func TestSendSendRedPacket(t *testing.T) {
	SendGrabPacket("1914080869", "1461819639", 1, "123", "100", "", "1111")
}
