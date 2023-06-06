/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package logic

import (
	pusher "crazy_server/internal/push"
	fcm "crazy_server/internal/push/fcm"
	"crazy_server/internal/push/getui"
	jpush "crazy_server/internal/push/jpush"
	"crazy_server/internal/push/mobpush"
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/kafka"
	promePkg "crazy_server/pkg/common/prometheus"
	"crazy_server/pkg/statistics"
	"fmt"
)

var (
	rpcServer     RPCServer
	pushCh        PushConsumerHandler
	producer      *kafka.Producer
	offlinePusher pusher.OfflinePusher
	successCount  uint64
)

func Init(rpcPort int) {
	rpcServer.Init(rpcPort)
	pushCh.Init()

}
func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	statistics.NewStatistics(&successCount, config.Config.ModuleName.PushName, fmt.Sprintf("%d second push to msg_gateway count", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	if *config.Config.Push.Getui.Enable {
		offlinePusher = getui.GetuiClient
	}
	if config.Config.Push.Jpns.Enable {
		offlinePusher = jpush.JPushClient
	}

	if config.Config.Push.Fcm.Enable {
		offlinePusher = fcm.NewFcm()
	}

	if config.Config.Push.Mob.Enable {
		offlinePusher = mobpush.MobPushClient
	}
}

func initPrometheus() {
	promePkg.NewMsgOfflinePushSuccessCounter()
	promePkg.NewMsgOfflinePushFailedCounter()
}

func Run(promethuesPort int) {
	go rpcServer.run()
	go pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&pushCh)
	go func() {
		err := promePkg.StartPromeSrv(promethuesPort)
		if err != nil {
			panic(err)
		}
	}()
}
