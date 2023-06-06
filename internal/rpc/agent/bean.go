package agent

import (
	"context"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/agent_model"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/proto/agent"
	"time"
)

// 获取平台咖豆商城配置
func (rpc *AgentServer) PlatformBeanShopConfig(_ context.Context, req *agent.PlatformBeanShopConfigReq) (*agent.PlatformBeanShopConfigResp, error) {
	resp := &agent.PlatformBeanShopConfigResp{}

	//获取平台咖豆redis缓存配置
	beanConfig, err := rocksCache.GetAgentPlatformBeanConfigCache()
	if err != nil {
		log.Error(req.OperationId, "获取平台咖豆商城配置缓存-GetAgentPlatformBeanConfigCache err :", err.Error())
		return resp, nil
	}

	resp.BeanShopConfig = make([]*agent.BeanShopConfig, 0)
	for _, v := range beanConfig {
		resp.BeanShopConfig = append(resp.BeanShopConfig, &agent.BeanShopConfig{
			ConfigId:       v.ConfigId,
			BeanNumber:     v.BeanNumber,
			GiveBeanNumber: v.GiveBeanNumber,
			Amount:         v.Amount,
			Status:         1,
		})
	}

	return resp, nil
}

// 推广员游戏咖豆商城配置
func (rpc *AgentServer) AgentGameBeanShopConfig(_ context.Context, req *agent.AgentGameBeanShopConfigReq) (*agent.AgentGameBeanShopConfigResp, error) {
	resp := &agent.AgentGameBeanShopConfigResp{}

	//获取推广员信息
	agentInfo, err := imdb.GetAgentByAgentNumber(req.AgentNumber)
	if err != nil || agentInfo.OpenStatus == 0 {
		return resp, nil
	}

	//获取推广员自定义咖豆配置
	configList, _ := imdb.GetAgentDiyShopBeanOnlineConfig(agentInfo.UserId)
	if len(configList) > 0 {
		resp.BeanShopConfig = make([]*agent.BeanShopConfig, 0)
		for _, v := range configList {
			//推广员咖豆余额 < 配置咖豆(购买值+赠送值)、不返回
			if agentInfo.BeanBalance < (v.BeanNumber + int64(v.GiveBeanNumber)) {
				continue
			}

			resp.BeanShopConfig = append(resp.BeanShopConfig, &agent.BeanShopConfig{
				ConfigId:       v.Id,
				BeanNumber:     v.BeanNumber,
				GiveBeanNumber: v.GiveBeanNumber,
				Amount:         v.Amount,
			})
		}
	}

	return resp, nil
}

// 推广员自定义咖豆商城配置
func (rpc *AgentServer) AgentDiyBeanShopConfig(_ context.Context, req *agent.AgentDiyBeanShopConfigReq) (*agent.AgentDiyBeanShopConfigResp, error) {
	resp := &agent.AgentDiyBeanShopConfigResp{BeanShopConfig: []*agent.BeanShopConfig{}, TodaySales: 0}

	//查询推广员是否存在
	agentInfo, err := imdb.GetAgentByUserId(req.UserId)
	if err != nil || agentInfo.OpenStatus == 0 {
		return resp, nil
	}

	//获取今日出售咖豆数
	resp.TodaySales = imdb.GetAgentTodaySalesNumber(req.UserId)

	//获取推广员自定义咖豆配置
	configList, _ := imdb.GetAgentDiyShopBeanConfig(req.UserId)
	if len(configList) > 0 {
		for _, v := range configList {
			status := v.Status
			// 上架状态、判断咖豆
			if v.Status == 1 && ((v.BeanNumber + int64(v.GiveBeanNumber)) > agentInfo.BeanBalance) {
				status = 0
			}

			resp.BeanShopConfig = append(resp.BeanShopConfig, &agent.BeanShopConfig{
				ConfigId:       v.Id,
				BeanNumber:     v.BeanNumber,
				GiveBeanNumber: v.GiveBeanNumber,
				Amount:         v.Amount,
				Status:         status,
			})
		}
	}

	return resp, nil
}

// 咖豆账户明细详情列表
func (rpc *AgentServer) AgentBeanAccountRecordList(_ context.Context, req *agent.AgentBeanAccountRecordListReq) (*agent.AgentBeanAccountRecordListResp, error) {
	resp := &agent.AgentBeanAccountRecordListResp{BeanRecordList: []*agent.BeanRecordList{}, Total: 0}

	//搜索用户
	chessUserIds := make([]int64, 0)
	if len(req.Keyword) > 0 {
		chessUserIds, _ = imdb.FindAgentMemberIds(req.UserId, req.Keyword)
		if len(chessUserIds) == 0 {
			return resp, nil
		}
	}

	list, count, err := imdb.BeanAccountRecordList(req.UserId, req.Date, req.BusinessType, req.Page, req.Size, chessUserIds)
	if err != nil {
		return resp, nil
	}

	resp.Total = count
	for _, v := range list {
		resp.BeanRecordList = append(resp.BeanRecordList, &agent.BeanRecordList{
			Type:         v.Type,
			BusinessType: v.BusinessType,
			Amount:       v.Amount,
			Number:       v.Number + int64(v.GiveNumber),
			//GiveNumber:        v.GiveNumber,
			ChessUserId:       v.ChessUserId,
			ChessUserNickname: v.ChessUserNickname,
			CreatedTime:       v.CreatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

// 咖豆管理上下架
func (rpc *AgentServer) AgentBeanShopUpStatus(_ context.Context, req *agent.AgentBeanShopUpStatusReq) (*agent.AgentBeanShopUpStatusResp, error) {
	resp := &agent.AgentBeanShopUpStatusResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}

	//查询推广员是否存在
	agentInfo, err := imdb.GetAgentByUserId(req.UserId)
	if err != nil || agentInfo.OpenStatus == 0 {
		resp.CommonResp.Msg = "推广员信息不存在"
		resp.CommonResp.Code = 400
		return resp, nil
	}

	if req.IsAll == 1 {
		if req.Status == 1 {
			configList, _ := imdb.GetAgentDiyShopBeanConfig(req.UserId)

			openConfigIds := make([]int32, 0)
			closeConfigIds := make([]int32, 0)

			//校验配置
			for _, v := range configList {
				// 判断咖豆
				if (v.BeanNumber + int64(v.GiveBeanNumber)) > agentInfo.BeanBalance {
					closeConfigIds = append(closeConfigIds, v.Id)
				} else {
					openConfigIds = append(openConfigIds, v.Id)
				}
			}

			if len(openConfigIds) > 0 {
				err = imdb.UpAgentDiyShopBeanConfigStatus(req.UserId, openConfigIds, 1)
			}

			if len(closeConfigIds) > 0 {
				err = imdb.UpAgentDiyShopBeanConfigStatus(req.UserId, closeConfigIds, 0)
			}

		} else {
			err = imdb.UpAgentDiyShopBeanConfigStatus(req.UserId, []int32{}, req.Status)
		}
	} else {
		if req.Status == 1 {
			configInfo, _ := imdb.GetAgentBeanConfigById(agentInfo.UserId, req.ConfigId)
			// 判断咖豆
			if (configInfo.BeanNumber + int64(configInfo.GiveBeanNumber)) > agentInfo.BeanBalance {
				resp.CommonResp.Msg = "咖豆余额不足，不能上架"
				resp.CommonResp.Code = 400
				return resp, nil
			}
		}

		err = imdb.UpAgentDiyShopBeanConfigStatus(req.UserId, []int32{req.ConfigId}, req.Status)
	}

	if err != nil {
		resp.CommonResp.Code = 400
		resp.CommonResp.Msg = "更新上下架失败"
		log.Error(req.OperationId, "更新上下架失败:%s", err.Error())
	}

	return resp, nil
}

// 咖豆管理(新增、编辑)
func (rpc *AgentServer) AgentBeanShopUpdate(_ context.Context, req *agent.AgentBeanShopUpdateReq) (*agent.AgentBeanShopUpdateResp, error) {
	resp := &agent.AgentBeanShopUpdateResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}

	//查询推广员是否存在
	agentInfo, err := imdb.GetAgentByUserId(req.UserId)
	if err != nil || agentInfo.OpenStatus == 0 {
		resp.CommonResp.Msg = "推广员信息不存在"
		resp.CommonResp.Code = 400
		return resp, nil
	}

	//删除历史咖豆配置
	_ = imdb.DelAgentDiyShopBeanConfig(req.UserId)

	//批量插入新配置
	if len(req.BeanShopConfig) > 0 {
		data := make([]*db.TAgentBeanShopConfig, 0)
		for _, v := range req.BeanShopConfig {
			// 上架状态、判断咖豆
			if v.Status == 1 && ((v.BeanNumber + int64(v.GiveBeanNumber)) > agentInfo.BeanBalance) {
				v.Status = 0
			}

			data = append(data, &db.TAgentBeanShopConfig{
				UserId:         req.UserId,
				BeanNumber:     v.BeanNumber,
				GiveBeanNumber: v.GiveBeanNumber,
				Amount:         v.Amount,
				Status:         v.Status,
				CreatedTime:    time.Now(),
				UpdatedTime:    time.Now(),
			})
		}

		err = imdb.InsertAgentDiyShopBeanConfigs(data)
		if err != nil {
			resp.CommonResp.Code = 400
			resp.CommonResp.Msg = "更新失败"
			log.Error(req.OperationId, "咖豆管理(新增、编辑)更新失败:%s", err.Error())
		}
	}

	return resp, nil
}
