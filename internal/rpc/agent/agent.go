package agent

import (
	"context"
	chessApi "crazy_server/pkg/agent"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/agent_model"
	rocksCache "crazy_server/pkg/common/db/rocks_cache"
	"crazy_server/pkg/common/log"
	"crazy_server/pkg/proto/agent"
	"github.com/spf13/cast"
	"math/rand"
	"time"
)

// 推广员申请提交
func (rpc *AgentServer) AgentApply(_ context.Context, req *agent.AgentApplyReq) (*agent.AgentApplyResp, error) {
	resp := &agent.AgentApplyResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}

	//查询是否已申请
	_, err := imdb.GetApplyByChessUserId(req.ChessUserId)
	if err == nil {
		resp.CommonResp.Msg = "已提交申请，请勿重复提交"
		resp.CommonResp.Code = 400
		return resp, nil
	}

	//申请数据入库
	err = imdb.AgentApply(&db.TAgentApplyRecord{
		UserId:      req.UserId,
		ChessUserId: req.ChessUserId,
		Name:        req.Name,
		Mobile:      req.Mobile,
	})

	if err != nil {
		log.Error(req.OperationId, "推广员申请提交数据入库失败:%s", err.Error())
		resp.CommonResp.Msg = "申请数据保存失败"
		resp.CommonResp.Code = 400
	}

	return resp, nil
}

// 绑定推广员
func (rpc *AgentServer) BindAgentNumber(_ context.Context, req *agent.BindAgentNumberReq) (*agent.BindAgentNumberResp, error) {
	resp := &agent.BindAgentNumberResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}

	//查询推广员是否存在
	agentInfo, err := imdb.GetAgentByAgentNumber(req.AgentNumber)
	if err != nil || agentInfo.OpenStatus == 0 {
		resp.CommonResp.Msg = "请输入正确的推广员ID"
		resp.CommonResp.Code = 400
		return resp, nil
	}

	//判断是否已绑定其他推广员
	_, err = imdb.AgentNumberByChessUserId(req.ChessUserId)
	if err == nil {
		resp.CommonResp.Msg = "该用户已绑定其他推广员"
		resp.CommonResp.Code = 400
		return resp, nil
	}

	//绑定推广员
	err = imdb.BindAgentNumber(&db.TAgentMember{
		UserId:        agentInfo.UserId,
		AgentNumber:   req.AgentNumber,
		ChessUserId:   req.ChessUserId,
		ChessNickname: req.ChessNickname,
		Day:           time.Now().Format("2006-01-02"),
	})

	if err != nil {
		log.Error(req.OperationId, "绑定推广员数据入库失败:%s", err.Error())
		resp.CommonResp.Msg = "绑定推广员失败"
		resp.CommonResp.Code = 400
	}

	return resp, nil
}

// 获取当前用户的推广员信息以及绑定关系
func (rpc *AgentServer) GetUserAgentInfo(_ context.Context, req *agent.GetUserAgentInfoReq) (*agent.GetUserAgentInfoResp, error) {
	resp := &agent.GetUserAgentInfoResp{
		IsAgent:         false,
		AgentNumber:     0,
		AgentName:       "",
		BindAgentNumber: 0,
	}

	//是否为推广员
	info, err := imdb.GetAgentByChessUserId(req.ChessUserId)
	if err == nil && info.OpenStatus == 1 {
		resp.IsAgent = true
		resp.AgentName = info.Name
		resp.AgentNumber = info.AgentNumber
	} else {
		//是否申请
		_, err = imdb.GetApplyByChessUserId(req.ChessUserId)
		if err == nil {
			resp.IsApply = true
		}
	}

	//是否绑定推广员
	agentMember, err := imdb.AgentNumberByChessUserId(req.ChessUserId)
	if err == nil {
		resp.BindAgentNumber = agentMember.AgentNumber
	}

	return resp, nil
}

// 推广员主页信息
func (rpc *AgentServer) AgentMainInfo(ctx context.Context, req *agent.AgentMainInfoReq) (*agent.AgentMainInfoResp, error) {
	resp := &agent.AgentMainInfoResp{}

	//获取推广员信息
	info, err := imdb.GetAgentByUserId(req.UserId)
	if err != nil || info.OpenStatus == 0 {
		return resp, nil
	}

	resp.AgentNumber = info.AgentNumber
	resp.AgentName = info.Name
	resp.Balance = info.Balance
	resp.BeanBalance = info.BeanBalance

	//累计收益、今日收益
	statIncome, _ := imdb.StatAgentIncomeData(req.UserId)
	resp.TodayIncome = statIncome.TodayIncome
	resp.AccumulatedIncome = statIncome.AccumulatedIncome

	//绑定的下属成员
	statMember, _ := imdb.StatAgentMemberData(req.UserId)
	resp.TodayBindUser = statMember.TodayBindUser
	resp.AccumulatedBindUser = statMember.AccumulatedBindUser

	//获取提现手续费比例(‰)千分之几
	commission := rocksCache.GetPlatformValueConfigCache("withdrawal_commission")
	resp.Commission = cast.ToInt32(commission)

	//提现次数
	todayWithdrawalNumber := rocksCache.GetWithdrawalNumber(ctx, req.UserId)
	resp.WithdrawalNumber = 3 - todayWithdrawalNumber //默认三次
	return resp, nil
}

// 账户明细收益趋势图
func (rpc *AgentServer) AgentAccountIncomeChart(_ context.Context, req *agent.AgentAccountIncomeChartReq) (*agent.AgentAccountIncomeChartResp, error) {
	resp := &agent.AgentAccountIncomeChartResp{}

	//获取收益统计数据
	chartData, _ := imdb.AccountIncomeChart(req.UserId, req.DateType)
	if len(chartData) > 0 {
		resp.IncomeChartData = make([]*agent.IncomeChartData, 0)
		for _, v := range chartData {
			resp.IncomeChartData = append(resp.IncomeChartData, &agent.IncomeChartData{
				Date:   v.Date,
				Income: v.Income,
			})
		}
	}

	return resp, nil
}

// 账户明细详情列表
func (rpc *AgentServer) AgentAccountRecordList(_ context.Context, req *agent.AgentAccountRecordListReq) (*agent.AgentAccountRecordListResp, error) {
	resp := &agent.AgentAccountRecordListResp{Total: 0, AccountRecordList: []*agent.AccountRecordList{}}

	//搜索用户
	chessUserIds := make([]int64, 0)
	if len(req.Keyword) > 0 {
		chessUserIds, _ = imdb.FindAgentMemberIds(req.UserId, req.Keyword)
		if len(chessUserIds) == 0 {
			return resp, nil
		}
	}

	//获取收益统计数据
	list, count, _ := imdb.AccountIncomeList(req.UserId, req.Date, req.BusinessType, req.Page, req.Size, chessUserIds)
	resp.Total = count
	if len(list) > 0 {
		for _, v := range list {
			resp.AccountRecordList = append(resp.AccountRecordList, &agent.AccountRecordList{
				BusinessType:      v.BusinessType,
				Amount:            v.Amount,
				CreatedTime:       v.CreatedTime.Format("2006-01-02 15:04:05"),
				Type:              v.Type,
				ChessUserId:       v.ChessUserId,
				ChessUserNickname: v.ChessUserNickname,
			})
		}
	}

	return resp, nil
}

// 推广下属用户列表
func (rpc *AgentServer) AgentMemberList(_ context.Context, req *agent.AgentMemberListReq) (*agent.AgentMemberListResp, error) {
	resp := &agent.AgentMemberListResp{Total: 0, AgentMemberList: []*agent.AgentMemberList{}}

	//获取推广员信息
	info, err := imdb.GetAgentByUserId(req.UserId)
	if err != nil || info.OpenStatus == 0 {
		return resp, nil
	}

	//req.OrderBy 排序(0默认-绑定时间倒序,1咖豆倒序,2咖豆正序,3贡献值倒序,4贡献值正序)
	//根据推广员编号获取所有下属成员
	total, chessUserIds, agentMemberList := chessApi.GetAgentChessMemberList(info.AgentNumber, req.OrderBy, req.Page, req.Size)
	if len(agentMemberList) == 0 {
		return resp, nil
	}

	//获取条件列表数据
	list, err := imdb.FindAgentMemberList(req.UserId, req.Keyword, chessUserIds, req.OrderBy, req.Page, req.Size)
	if err != nil {
		return resp, nil
	}

	resp.Total = total
	if len(list) > 0 {
		//按咖豆排序
		if len(chessUserIds) > 0 {
			// 下属成员贡献度
			userContribution := make(map[int64]int64, 0)
			for _, v := range list {
				userContribution[v.ChessUserId] = v.Contribution
			}

			for _, chessUserId := range chessUserIds {
				var contribution int64 = 0
				if c, ok := userContribution[chessUserId]; ok {
					contribution = c
				}

				member, _ := agentMemberList[chessUserId]
				resp.AgentMemberList = append(resp.AgentMemberList, &agent.AgentMemberList{
					ChessUserId:     member.Uid,
					ChessNickname:   member.Nickname,
					ChessBeanNumber: member.Gold,
					Contribution:    contribution,
				})
			}
		} else {
			for _, v := range list {
				chessNickname := v.ChessNickname
				var chessBeanNumber int64 = 0
				if member, ok := agentMemberList[v.ChessUserId]; ok {
					chessNickname = member.Nickname
					chessBeanNumber = member.Gold

					// 更新冗余昵称数据
					//if v.ChessNickname != member.Nickname {
					//
					//}
				}
				resp.AgentMemberList = append(resp.AgentMemberList, &agent.AgentMemberList{
					ChessUserId:     v.ChessUserId,
					ChessNickname:   chessNickname,
					ChessBeanNumber: chessBeanNumber,
					Contribution:    v.Contribution,
				})
			}
		}
	}
	return resp, nil
}

// 开通推广员
func (rpc *AgentServer) OpenAgent(_ context.Context, req *agent.OpenAgentReq) (*agent.OpenAgentResp, error) {
	resp := &agent.OpenAgentResp{CommonResp: &agent.CommonResp{Code: 0, Msg: ""}}

	//获取推广员申请记录
	info, err := imdb.GetApplyById(req.ApplyId)
	if err != nil {
		resp.CommonResp.Msg = "推广员申请记录不存在"
		resp.CommonResp.Code = 400
		return resp, nil
	}

	//已审批
	if info.AuditStatus == 1 {
		return resp, nil
	}

	//创建推广员账户
	err = imdb.CreateAgentAccount(&db.TAgentAccount{
		UserId:      info.UserId,
		Name:        info.Name,
		Mobile:      info.Mobile,
		ChessUserId: info.ChessUserId,
		AgentNumber: CreateAgentNumber(),
		OpenStatus:  1,
	})
	if err != nil {
		log.Error("", "开通推广员失败:%s", err.Error())
		resp.CommonResp.Msg = "开通推广员失败"
		resp.CommonResp.Code = 400
		return resp, nil
	}

	//更新审核状态
	err = imdb.UpApplyAuditStatus(req.ApplyId)
	if err != nil {
		log.Error("", "更新推广员申请记录失败:%s", err.Error())
		resp.CommonResp.Msg = "开通推广员失败"
		resp.CommonResp.Code = 400
		return resp, nil
	}

	return resp, nil
}

// 生成推广员编号
func CreateAgentNumber() int32 {
	rand.Seed(time.Now().UnixNano())
	agentNumber := int32(100000 + rand.Intn(900000))

	_, err := imdb.GetAgentByAgentNumber(agentNumber)
	if err != nil {
		return agentNumber
	}

	return int32(100000 + rand.Intn(900000))
}

// 获取推广员开通状态
func (rpc *AgentServer) GetAgentOpenStatus(ctx context.Context, req *agent.GetAgentOpenStatusReq) (*agent.GetAgentOpenStatusResp, error) {
	resp := &agent.GetAgentOpenStatusResp{}

	//获取推广员信息
	agentInfo, err := imdb.GetAgentByUserId(req.UserId)
	if err != nil || agentInfo.OpenStatus == 0 {
		return resp, nil
	}

	resp.AgentOpenStatus = true

	return resp, nil
}
