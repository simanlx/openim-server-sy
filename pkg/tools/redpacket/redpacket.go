package redpacket

import (
	"math/rand"
	"time"
)

func GetRedPacket(amount, number int) []int {
	if amount < number {
		return nil
	}
	if number == amount {
		result := []int{}
		for i := 0; i < number; i++ {
			result = append(result, 1)
		}
		return result
	}

	remain := amount
	var result []int
	sum := 0
	for i := 0; i < number; i++ {
		x := DoubleAverage(number-i, remain)
		//金额减去
		remain -= x
		//发了多少钱
		sum += x

		result = append(result, x)
	}
	return result
}

//二倍均值算法,count剩余个数,amount剩余金额
func DoubleAverage(count, amount int) int {
	//最小钱
	min := 1
	if count == 1 {
		//返回剩余金额
		return amount
	}
	//计算最大可用金额,min最小是1分钱,减去的min,下面会加上,避免出现0分钱
	max := amount - min*count
	//计算最大可用平均值
	avg := max / count
	//二倍均值基础加上最小金额,防止0出现,作为上限
	avg2 := 2*avg + min
	//随机红包金额序列元素,把二倍均值作为随机的最大数
	rand.Seed(time.Now().UnixNano())
	//加min是为了避免出现0值,上面也减去了min
	x := rand.Intn(avg2) + min
	return x
}
