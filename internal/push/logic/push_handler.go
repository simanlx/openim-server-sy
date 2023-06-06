/*
** description("").
** copyright('crazy_server,www.crazy_server.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/13 10:33).
 */
package logic

import (
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	kfk "crazy_server/pkg/common/kafka"
	"crazy_server/pkg/common/log"
	pbChat "crazy_server/pkg/proto/msg"
	pbPush "crazy_server/pkg/proto/push"
	"crazy_server/pkg/utils"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type fcb func(msg []byte)

type PushConsumerHandler struct {
	msgHandle         map[string]fcb
	pushConsumerGroup *kfk.MConsumerGroup
}

func (ms *PushConsumerHandler) Init() {
	ms.msgHandle = make(map[string]fcb)
	ms.msgHandle[config.Config.Kafka.Ms2pschat.Topic] = ms.handleMs2PsChat
	ms.pushConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ms2pschat.Topic}, config.Config.Kafka.Ms2pschat.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
}
func (ms *PushConsumerHandler) handleMs2PsChat(msg []byte) {
	log.NewDebug("", "msg come from kafka  And push!!!", "msg", string(msg))
	msgFromMQ := pbChat.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg, &msgFromMQ); err != nil {
		log.Error("", "push Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	pbData := &pbPush.PushMsgReq{
		OperationID:  msgFromMQ.OperationID,
		MsgData:      msgFromMQ.MsgData,
		PushToUserID: msgFromMQ.PushToUserID,
	}
	sec := msgFromMQ.MsgData.SendTime / 1000
	nowSec := utils.GetCurrentTimestampBySecond()
	if nowSec-sec > 10 {
		return
	}
	switch msgFromMQ.MsgData.SessionType {
	case constant.SuperGroupChatType:
		// 这里是发送给超级群
		MsgToSuperGroupUser(pbData)
	default:
		// 这里是发送给用户
		MsgToUser(pbData)
	}
	//Call push module to send message to the user
	//MsgToUser((*pbPush.PushMsgReq)(&msgFromMQ))
}
func (PushConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (PushConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (ms *PushConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
		ms.msgHandle[msg.Topic](msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}
