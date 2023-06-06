package redpacket

import (
	"fmt"
	"testing"
)

func TestDoubleAverage(t *testing.T) {
	fmt.Println(GetRedPacket(2, 2))
	fmt.Println(GetRedPacket(40, 2))
	fmt.Println(GetRedPacket(50, 1))
	fmt.Println(GetRedPacket(50, 12))
	fmt.Println(GetRedPacket(60, 12))
	fmt.Println(GetRedPacket(2, 2))
}
