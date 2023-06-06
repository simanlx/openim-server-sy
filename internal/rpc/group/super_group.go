package group

import (
	"context"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	cp "crazy_server/pkg/common/utils"
	pbGroup "crazy_server/pkg/proto/group"
	commonPb "crazy_server/pkg/proto/sdk_ws"
	"crazy_server/pkg/utils"

	"github.com/go-redis/redis/v8"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetJoinedSuperGroupListResp{CommonResp: &pbGroup.CommonResp{}}
	groupIDList, err := rocksCache.GetJoinedSuperGroupListFromCache(req.UserID)
	if err != nil {
		if err == redis.Nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroupByUserID nil ", err.Error(), req.UserID)
			return resp, nil
		}
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroupByUserID failed ", err.Error(), req.UserID)
		//resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		//resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	for _, groupID := range groupIDList {
		groupInfoFromCache, err := rocksCache.GetGroupInfoFromCache(groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupInfoByGroupID failed", groupID, err.Error())
			continue
		}
		groupInfo := &commonPb.GroupInfo{}
		if err := utils.CopyStructFields(groupInfo, groupInfoFromCache); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		}
		groupMemberIDList, err := rocksCache.GetGroupMemberIDListFromCache(groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroup failed", groupID, err.Error())
			continue
		}
		groupInfo.MemberCount = uint32(len(groupMemberIDList))
		resp.GroupList = append(resp.GroupList, groupInfo)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *groupServer) GetSuperGroupsInfo(_ context.Context, req *pbGroup.GetSuperGroupsInfoReq) (resp *pbGroup.GetSuperGroupsInfoResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbGroup.GetSuperGroupsInfoResp{CommonResp: &pbGroup.CommonResp{}}
	groupsInfoList := make([]*commonPb.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromRedis, err := rocksCache.GetGroupInfoFromCache(groupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), groupID)
			continue
		}
		var groupInfo commonPb.GroupInfo
		cp.GroupDBCopyOpenIM(&groupInfo, groupInfoFromRedis)
		groupsInfoList = append(groupsInfoList, &groupInfo)
	}
	resp.GroupInfoList = groupsInfoList
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}
