package auth

import (
	"context"
	"crazy_server/pkg/common/constant"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	promePkg "crazy_server/pkg/common/prometheus"
	"crazy_server/pkg/common/token_verify"
	"crazy_server/pkg/grpc-etcdv3/getcdv3"
	pbAuth "crazy_server/pkg/proto/auth"
	pbRelay "crazy_server/pkg/proto/relay"
	"crazy_server/pkg/utils"
	"errors"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"crazy_server/pkg/common/config"

	"google.golang.org/grpc"
)

// 单点登录
func (rpc *rpcAuth) SingleLogin(ctx context.Context, req *pbAuth.SingleLoginReq) (*pbAuth.SingleLoginResp, error) {
	resp := &pbAuth.SingleLoginResp{CommonResp: &pbAuth.CommonResp{ErrCode: 0, ErrMsg: ""}}

	// 删除jwt token
	rocksCache.DelUserJwtToken(ctx, req.UserId)

	// 安卓
	if err := rpc.forceKickOff(req.UserId, 2, req.OperationID); err != nil {
		errMsg := req.OperationID + " SingleLogin android forceKickOff failed " + err.Error() + req.UserId + utils.Int32ToString(req.Platform)
		log.NewError(req.OperationID, errMsg)
	}

	// ios
	if err := rpc.forceKickOff(req.UserId, 1, req.OperationID); err != nil {
		errMsg := req.OperationID + " SingleLogin ios forceKickOff failed " + err.Error() + req.UserId + utils.Int32ToString(req.Platform)
		log.NewError(req.OperationID, errMsg)
	}

	return resp, nil
}

func (rpc *rpcAuth) UserRegister(ctx context.Context, req *pbAuth.UserRegisterReq) (*pbAuth.UserRegisterResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "UserRegister注册>", utils.JsonFormat(req))
	var user db.User
	utils.CopyStructFields(&user, req.UserInfo)
	// 如果用户逇信息中有生日，需要转换一下
	if req.UserInfo.BirthStr != "" {
		time, err := utils.TimeStringToTime(req.UserInfo.BirthStr)
		if err != nil {
			log.NewError(req.OperationID, "TimeStringToTime failed ", err.Error(), req.UserInfo.BirthStr)
			return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: "TimeStringToTime failed:" + err.Error()}}, nil
		}
		user.Birth = time
	}
	log.Debug(req.OperationID, "copy ", user, req.UserInfo)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "UserRegister注册-2>", utils.JsonFormat(user))

	// 创建用户账号
	err := imdb.UserRegister(user)
	if err != nil {
		errMsg := req.OperationID + " imdb.UserRegister failed " + err.Error() + user.UserID
		log.NewError(req.OperationID, errMsg, user)
		return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}

	//写入redis
	_ = rocksCache.SetCrazyUserToken(ctx, user.UserID, user.Unionid)

	promePkg.PromeInc(promePkg.UserRegisterCounter)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}})
	return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}}, nil
}

func (rpc *rpcAuth) UserToken(ctx context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	user, err := imdb.GetUserByUserID(req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, "not this user:", req.FromUserID, req.String())
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	//单点登录
	_, _ = rpc.SingleLogin(ctx, &pbAuth.SingleLoginReq{
		Platform:    req.Platform,
		UserId:      req.FromUserID,
		OperationID: req.OperationID,
	})

	//写入redis
	_ = rocksCache.SetCrazyUserToken(ctx, user.UserID, user.Unionid)

	tokens, expTime, err := token_verify.CreateToken(req.FromUserID, int(req.Platform))
	if err != nil {
		errMsg := req.OperationID + " token_verify.CreateToken failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	promePkg.PromeInc(promePkg.UserLoginCounter)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}, Token: tokens, ExpiredTime: expTime})
	return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}, Token: tokens, ExpiredTime: expTime}, nil
}

func (rpc *rpcAuth) ParseToken(_ context.Context, req *pbAuth.ParseTokenReq) (*pbAuth.ParseTokenResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	claims, err := token_verify.ParseToken(req.Token, req.OperationID)
	if err != nil {
		errMsg := "ParseToken failed " + err.Error() + req.OperationID + " token " + req.Token
		log.Error(req.OperationID, errMsg, "token:", req.Token)
		return &pbAuth.ParseTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: 4001, ErrMsg: errMsg}}, nil
	}
	resp := pbAuth.ParseTokenResp{CommonResp: &pbAuth.CommonResp{}, UserID: claims.UID, Platform: claims.Platform, ExpireTimeSeconds: uint32(claims.ExpiresAt.Unix())}
	log.Info(req.OperationID, utils.GetSelfFuncName(), " rpc return ", resp.String())
	return &resp, nil
}

func (rpc *rpcAuth) ForceLogout(_ context.Context, req *pbAuth.ForceLogoutReq) (*pbAuth.ForceLogoutResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " IsManagerUserID false " + req.OpUserID
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	//if err := token_verify.DeleteToken(req.FromUserID, int(req.Platform)); err != nil {
	//	errMsg := req.OperationID + " DeleteToken failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
	//	log.NewError(req.OperationID, errMsg)
	//	return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	//}
	if err := rpc.forceKickOff(req.FromUserID, req.Platform, req.OperationID); err != nil {
		errMsg := req.OperationID + " forceKickOff failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}})
	return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{}}, nil
}

func (rpc *rpcAuth) forceKickOff(userID string, platformID int32, operationID string) error {
	log.NewInfo(operationID, utils.GetSelfFuncName(), " args ", userID, platformID)
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), operationID)
	for _, v := range grpcCons {
		client := pbRelay.NewRelayClient(v)
		kickReq := &pbRelay.KickUserOfflineReq{OperationID: operationID, KickUserIDList: []string{userID}, PlatformID: platformID}
		log.NewInfo(operationID, "KickUserOffline ", client, kickReq.String())
		_, err := client.KickUserOffline(context.Background(), kickReq)
		return utils.Wrap(err, "")
	}
	return errors.New("no rpc node ")
}

type rpcAuth struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string

	pbAuth.UnsafeAuthServer
}

func NewRpcAuthServer(port int) *rpcAuth {
	log.NewPrivateLog(constant.LogFileName)
	return &rpcAuth{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAuthName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (rpc *rpcAuth) Run() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, "rpc auth start...")

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
	log.NewInfo(operationID, "listen network success, ", address, listener)
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		promePkg.NewUserRegisterCounter()
		promePkg.NewUserLoginCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()

	//service registers with etcd
	pbAuth.RegisterAuthServer(srv, rpc)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)

	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log.NewError(operationID, "RegisterEtcd failed ", err.Error(),
			rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))

	}
	log.NewInfo(operationID, "RegisterAuthServer ok ", rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError(operationID, "Serve failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "rpc auth ok")
}
