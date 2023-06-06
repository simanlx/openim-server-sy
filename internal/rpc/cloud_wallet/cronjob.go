package cloud_wallet

import (
	"bytes"
	"crazy_server/pkg/cloud_wallet/ncount"
	imdb "crazy_server/pkg/common/db/mysql_model/cloud_wallet"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/contrive_msg"
	"crazy_server/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var NotifyChannel = make(chan string, 100)

func StarCorn() {
	// 启动消费者
	HandleThirdPayNotifyLogic()
	HandleRedPacketReturn() // 处理红包退回

	// 启动生产者
	StartThirdPayNotifyTicker() // 处理支付回调
}

// 启动定时任务 ： 每分钟查询过去2个小时所有订单
// 进行回调，主要是针对第三方支付- 新互娱支付
func StartThirdPayNotifyTicker() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic ,定时任务：查询历史第三方支付订单失败，err: ", err)
			}
		}()
		for {
			// 每分钟查询过去2个小时所有订单
			start_time := time.Now().Add(-time.Hour * 2).Format("2006-01-02 15:04:05")
			end_time := time.Now().Format("2006-01-02 15:04:05")
			result, err := imdb.GetThirdPayOrderListByTime(start_time, end_time)
			if err != nil {
				log.Error("查询订单失败，err: ", err)
				time.Sleep(time.Minute)
				continue
			}
			operation := utils.OperationIDGenerator()
			log.Info(operation, "第三方支付回调任务-master：", len(result), "个订单")
			// 将订单数量写入channel通道
			for _, v := range result {

				// 1.如果是notify_count  =1 ,间隔时间为30秒
				if v.LastNotifyTime.Add(time.Second*30).Unix() > time.Now().Unix() && v.NotifyCount == 1 {
					continue
				}
				// 2.如果是notify_count  =2 ,间隔时间为5分钟
				if v.LastNotifyTime.Add(time.Minute*5).Unix() > time.Now().Unix() && v.NotifyCount == 2 {
					continue
				}
				// 3.如果是notify_count  =3 ,间隔时间为30分钟
				if v.LastNotifyTime.Add(time.Minute*30).Unix() > time.Now().Unix() && v.NotifyCount == 3 {
					continue
				}
				// 4.如果是notify_count  =4 ,间隔时间为30分钟
				if v.LastNotifyTime.Add(time.Minute*30).Unix() > time.Now().Unix() && v.NotifyCount == 4 {
					continue
				}
				// 5.如果是notify_count  =5 ,间隔时间为30分钟
				if v.LastNotifyTime.Add(time.Minute*30).Unix() > time.Now().Unix() && v.NotifyCount == 5 {
					continue
				}
				notifyThirdPay(v.NcountOrderNo)
			}
			time.Sleep(time.Minute)
		}
	}()
}

// 就是针对回调实际处理的函数
func HandleThirdPayNotifyLogic() {
	for i := 0; i < 50; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Error("Panic 通知失败，err: ", err)
				}
			}()
			for {
				MerOrderID := <-NotifyChannel // 我们平台生成的订单号
				Operation := "竞技回调merOrderID ：" + MerOrderID
				// 查询订单信息
				err, payOrder := imdb.GetThirdPayJdnMerOrderID(MerOrderID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						log.Error(Operation, "没有订单ID ：", MerOrderID)
					} else {
						log.Error(Operation, "查询订单失败 ：", err)
					}
					continue
				}
				log.Info("第三方支付回调任务-son：", payOrder.MerOrderNo, "订单状态：", payOrder.Status, "回调次数：", payOrder.NotifyCount)
				// 如果请求超过5次就不再请求
				if payOrder.NotifyCount >= 5 {
					continue
				}

				// 如果超过两个小时后就不再请求
				if payOrder.AddTime.Add(time.Hour*2).Unix() < time.Now().Unix() {
					continue
				}

				// 检查上次通知时间到现在的时间间隔
				// 第一次是及时回调
				// 第二次是间隔5秒 ： 由线程休眠实现
				// 第三次是间隔30秒 ：数据库触发
				// 第四次后就是间隔5分钟
				// 第五次后就是间隔30分钟
				content := map[string]interface{}{
					"OrderID":    payOrder.OrderNo,
					"MerOrderID": payOrder.MerOrderNo,
					"Status":     payOrder.Status, // 100 是未支付，200是支付成功，300是支付失败
					"CreateTime": payOrder.AddTime.Format("2006-01-02 15:04:05"),
					"PayTime":    payOrder.PayTime.Format("2006-01-02 15:04:05"),
					"Amount":     payOrder.Amount, // 支付金额，以分为单位，整数
				}
				// 转换成为json
				body, err := json.Marshal(content)
				if err != nil {
					log.Error(Operation, "json转换失败，err: ", err)
					continue
				}
				// 发送请求
				err = HttpPost(payOrder.NotifyUrl, body)
				if err != nil {
					log.Error(Operation, "发送请求失败，err: ", err)
					// 修改订单状态
					err = imdb.UpdateThirdPayOrderCallback(0, int(payOrder.NotifyCount)+1, MerOrderID)
					if err != nil {
						log.Error(Operation, "修改订单状态失败，err: ", err)
					}
					continue
				}
				// 修改订单状态
				err = imdb.UpdateThirdPayOrderCallback(1, int(payOrder.NotifyCount)+1, MerOrderID)
				if err != nil {
					log.Error(Operation, "修改订单状态失败，err: ", err)
				}
			}

		}()
	}
}

// 处理红包退款
func HandleRedPacketReturn() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic ,定时任务：进行红包退回，err: ", err)
			}
		}()
		for {
			// 查询一段时间内的红包
			PacektSet, err := imdb.GetExpiredRedPacketListByPage()
			if err != nil {
				time.Sleep(time.Minute)
				continue
			}
			log.Info("红包退回脚本-master：定时退回", PacektSet)
			for _, v := range PacektSet {
				// 这里进行退款操作
				err = returnRedPacket(v.PacketID)
				if err != nil {
					log.Error("红包退回脚本-son：定时退回失败", err)
					continue
				}
			}
			time.Sleep(time.Minute)
		}
	}()
}

// 查询所有的红包进行退款
// 红包退款
func returnRedPacket(RedPacketID string) error {
	// 查询红包信息
	PacketInfo, err := imdb.GetRedPacketInfo(RedPacketID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	// 判断红包过期时间是否小于当前时间
	if PacketInfo.ExpireTime > time.Now().Unix() {
		// 红包未过期
		return nil
	}

	// 红包状态为正常的红包才进行退款
	if PacketInfo.Status != 1 {
		return nil
	}

	// 红包剩余金额必须大于0
	if PacketInfo.RemainAmout <= 0 {
		return nil
	}

	// 查询用户信息
	User, err := rocksCache.GetUserAccountInfoFromCache(PacketInfo.UserID)
	if err != nil {
		// 查询用户信息失败
		return err
	}

	// 将红包的剩余的金额退回到用户的账户中
	nc := NewNcountPay()
	merOrderID := ncount.GetMerOrderID()
	remainAmout := PacketInfo.RemainAmout
	totalAmount := cast.ToString(cast.ToFloat64(remainAmout) / 100)
	payResult := nc.payByBalance("红包退款", User.PacketAccountId, User.MainAccountId, merOrderID, totalAmount)
	if payResult.ErrCode == 0 {
		// 红包退款成功
		PacketInfo.Status = 100
		PacketInfo.Remark = "红包已退款"
		PacketInfo.RefoundAmout = PacketInfo.RemainAmout
	} else {
		// 红包退款失败 红包退款异常，需要进行数据库核查
		PacketInfo.Status = 200
		PacketInfo.Remark = "红包退款失败，需要进行数据核查:" + payResult.ErrMsg
		// 这里最好进行报警处理 todo 通知管理员
	}

	PacketInfo.UpdatedTime = time.Now().Unix()
	err = imdb.UpdateRedPacketInfo(PacketInfo.PacketID, PacketInfo)
	if err != nil {
		return err
	}
	if PacketInfo.Status == 200 {
		return nil
	}

	// 修改用户交易记录
	err = AddNcountTradeLog(BusinessTypePacketExpire, int32(remainAmout), PacketInfo.UserID, User.MainAccountId, merOrderID, payResult.NcountOrderID, "")
	if err != nil {
		return err
	}

	// 通知用户红包退款
	// OperationID, redPacketID, content string, sessionID int, SenderID, ReciveID
	err = contrive_msg.SendRebackMessage("红包退款消息", PacketInfo.PacketID, "红包过期，剩余金额退回到钱包", int(PacketInfo.PacketType), PacketInfo.UserID, PacketInfo.RecvID)
	if err != nil {
		return err
	}
	return err
}

func notifyThirdPay(MerOrderID string) {
	NotifyChannel <- MerOrderID
}

func HttpPost(Url string, content []byte) error {
	body := bytes.NewBuffer(content)
	resp, err := http.Post(Url, "application/json;charset=utf-8", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 判断状态码
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Info("第三方回调失败", fmt.Sprintf("请求参数：%s,HTTPCode:%v,响应参数：%s", content, resp.StatusCode, string(b)))
		return errors.New("返回响应HttpCode不为200," + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
