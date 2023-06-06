package agent

import (
	"crazy_server/pkg/common/config"
	"crazy_server/pkg/common/http"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/common/utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"sort"
)

type RechargeUserGoldReq struct {
	OrderId     string `json:"order_id"`
	Uid         int64  `json:"uid"`
	AgentNumber int32  `json:"agent_number"`
	Num         int64  `json:"num"`
}

type ChessApiResp struct {
	Msg  string `json:"msg"`
	Code int64  `json:"code"`
}

type AgentChessUserListResp struct {
	Code int32            `json:"code"`
	Msg  string           `json:"msg"`
	Data []*ChessUserInfo `json:"data"`
}

type ChessUserInfo struct {
	Uid      int64  `json:"uid"`
	Nickname string `json:"nickname"`
	Gold     int64  `json:"gold"`
}

// 调用chess api 获取推广员下属成员列表
func ChessApiAgentChessMemberList(agentNumber int32) ([]*ChessUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v1/getUserListInfo?agent_number=%d", config.Config.Agent.ChessApiDomain, agentNumber))
	if err != nil {
		log.Error("", "请求chess api getUserListInfo", err.Error())
		return nil, errors.Wrap(err, "请求chess api getUserListInfo失败")
	}

	chessApiResp := &AgentChessUserListResp{Data: []*ChessUserInfo{}}
	_ = json.Unmarshal(resp, &chessApiResp)
	if chessApiResp.Code != 200 {
		errMsg := fmt.Sprintf("调用chess api 接口失败, err:%s", chessApiResp.Msg)
		log.Error("", "请求chess api getUserListInfo", errMsg)
		return nil, errors.New(errMsg)
	}
	return chessApiResp.Data, nil
}

// 获取推广员下属成员列表
func GetAgentChessMemberList(agentNumber, orderBy, page, size int32) (int64, []int64, map[int64]*ChessUserInfo) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("", "获取推广员下属成员列表GetAgentChessMemberList panic:参数->", agentNumber, orderBy, page, size)
			return
		}
	}()

	chessUserIds := make([]int64, 0)
	chessUserMap := map[int64]*ChessUserInfo{}

	agentChessMemberList, err := ChessApiAgentChessMemberList(agentNumber)
	memberTotal := int64(len(agentChessMemberList))
	if err != nil || memberTotal == 0 {
		return memberTotal, chessUserIds, chessUserMap
	}

	if orderBy == 1 {
		//咖豆倒序
		sort.Slice(agentChessMemberList, func(p, q int) bool {
			return agentChessMemberList[p].Gold > agentChessMemberList[q].Gold
		})

	} else if orderBy == 2 {
		//咖豆正序
		sort.Slice(agentChessMemberList, func(p, q int) bool {
			return agentChessMemberList[p].Gold < agentChessMemberList[q].Gold
		})
	} else {
		// 无排序、返回uid map
		for _, v := range agentChessMemberList {
			chessUserMap[v.Uid] = v
		}
		return memberTotal, chessUserIds, chessUserMap
	}

	//分页获取
	chessMemberList := SliceListPage(agentChessMemberList, int(page), int(size))
	for _, v := range chessMemberList {
		chessUserIds = append(chessUserIds, v.Uid)
		chessUserMap[v.Uid] = v
	}

	return memberTotal, chessUserIds, chessUserMap
}

// 分页
func SliceListPage(list []*ChessUserInfo, page, size int) []*ChessUserInfo {
	sliceLen := len(list)

	if size > sliceLen {
		return list
	}

	chessUserList := make([]*ChessUserInfo, 0)

	// 总页数计算
	pageCount := int(math.Ceil(float64(sliceLen) / float64(size)))
	if page > pageCount {
		return chessUserList
	}

	start := (page - 1) * size
	end := start + size
	if end > sliceLen {
		end = sliceLen
	}
	return list[start:end]
}

// 调用chess api 给用户加咖豆
func ChessApiGiveUserBean(orderNo string, agentNumber int32, chessUserId, beanNumber int64) error {
	data := RechargeUserGoldReq{
		OrderId:     orderNo,
		Uid:         chessUserId,
		AgentNumber: agentNumber,
		Num:         beanNumber,
	}
	resp, err := http.Post(config.Config.Agent.ChessApiDomain+"/v1/rechargeUserGold", data, 3)
	if err != nil {
		//企业微信告警
		go utils.AgentGiveMemberBeanWarn(agentNumber, chessUserId, beanNumber, err.Error())

		log.Error("", "请求chess api rechargeUserGold失败", err.Error())
		return errors.Wrap(err, "请求chess api rechargeUserGold失败")
	}

	chessApiResp := &ChessApiResp{}
	_ = json.Unmarshal(resp, &chessApiResp)
	if chessApiResp.Code != 200 {
		//企业微信告警
		go utils.AgentGiveMemberBeanWarn(agentNumber, chessUserId, beanNumber, err.Error())

		errMsg := fmt.Sprintf("调用chess api 给用户加咖豆失败, 订单号(%s),推广员编号(%d),互娱用户id(%d),咖豆数(%d),err:%s", orderNo, agentNumber, chessUserId, beanNumber, chessApiResp.Msg)
		log.Error("", "请求chess api rechargeUserGold失败", errMsg)
		return errors.New(errMsg)
	}
	return nil
}
