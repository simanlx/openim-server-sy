package msg

import (
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/db"
	"crazy_server/pkg/common/kafka"
	"crazy_server/pkg/common/log"
	promePkg "crazy_server/pkg/common/prometheus"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	"crazy_server/pkg/proto/msg"
	"crazy_server/pkg/utils"
	"net"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
)

type MessageWriter interface {
	SendMessage(m proto.Message, key string, operationID string) (int32, int64, error)
}
type rpcChat struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	messageWriter   MessageWriter
	//offlineProducer *kafka.Producer
	delMsgCh       chan deleteMsg
	dMessageLocker MessageLocker

	msg.UnimplementedMsgServer
}

type deleteMsg struct {
	UserID      string
	OpUserID    string
	SeqList     []uint32
	OperationID string
}

func NewRpcChatServer(port int) *rpcChat {
	log.NewPrivateLog(constant.LogFileName)
	rc := rpcChat{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImMsgName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
		dMessageLocker:  NewLockerMessage(),
	}
	rc.messageWriter = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	//rc.offlineProducer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschatOffline.Addr, config.Config.Kafka.Ws2mschatOffline.Topic)
	rc.delMsgCh = make(chan deleteMsg, 1000)
	return &rc
}

func (rpc *rpcChat) initPrometheus() {
	//sendMsgSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
	//	Name: "send_msg_success",
	//	Help: "The number of send msg success",
	//})
	//sendMsgFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
	//	Name: "send_msg_failed",
	//	Help: "The number of send msg failed",
	//})
	promePkg.NewMsgPullFromRedisSuccessCounter()
	promePkg.NewMsgPullFromRedisFailedCounter()
	promePkg.NewMsgPullFromMongoSuccessCounter()
	promePkg.NewMsgPullFromMongoFailedCounter()

	promePkg.NewSingleChatMsgRecvSuccessCounter()
	promePkg.NewGroupChatMsgRecvSuccessCounter()
	promePkg.NewWorkSuperGroupChatMsgRecvSuccessCounter()

	promePkg.NewSingleChatMsgProcessSuccessCounter()
	promePkg.NewSingleChatMsgProcessFailedCounter()
	promePkg.NewGroupChatMsgProcessSuccessCounter()
	promePkg.NewGroupChatMsgProcessFailedCounter()
	promePkg.NewWorkSuperGroupChatMsgProcessSuccessCounter()
	promePkg.NewWorkSuperGroupChatMsgProcessFailedCounter()
}

func (rpc *rpcChat) Run() {
	log.Info("", "rpcChat init...")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(rpc.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + rpc.rpcRegisterName)
	}
	log.Info("", "listen network success, address ", address)
	recvSize := 1024 * 1024 * 30
	sendSize := 1024 * 1024 * 30
	var grpcOpts = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(recvSize),
		grpc.MaxSendMsgSize(sendSize),
	}
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()

	rpcRegisterIP := config.Config.RpcRegisterIP
	msg.RegisterMsgServer(srv, rpc)
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log.Error("", "register rpcChat to etcd failed ", err.Error())
		panic(utils.Wrap(err, "register chat module  rpc to etcd err"))
	}
	go rpc.runCh()
	rpc.initPrometheus()
	err = srv.Serve(listener)
	if err != nil {
		log.Error("", "rpc rpcChat failed ", err.Error())
		return
	}
	log.Info("", "rpc rpcChat init success")
}

func (rpc *rpcChat) runCh() {
	log.NewInfo("", "start del msg chan ")
	for {
		select {
		case msg := <-rpc.delMsgCh:
			log.NewInfo(msg.OperationID, utils.GetSelfFuncName(), "delmsgch recv new: ", msg)
			if len(msg.SeqList) > 0 {
				db.DB.DelMsgFromCache(msg.UserID, msg.SeqList, msg.OperationID)
				DeleteMessageNotification(msg.OpUserID, msg.UserID, msg.SeqList, msg.OperationID)
			}
		}
	}
}
