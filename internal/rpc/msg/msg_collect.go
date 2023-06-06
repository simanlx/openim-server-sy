package msg

import (
	"context"
	"crazy_server/pkg/common/db"
	imdb "crazy_server/pkg/common/db/mysql_model/im_mysql_model"
	pbMsg "crazy_server/pkg/proto/msg"
	"errors"
	"fmt"
)

// 消息收藏
func (rpc *rpcChat) MsgCollect(_ context.Context, req *pbMsg.MsgCollectReq) (resp *pbMsg.MsgCollectResp, err error) {
	resp = &pbMsg.MsgCollectResp{}

	//保存收藏数据
	err = imdb.InsertUserCollect(&db.UserCollect{
		UserId:         req.UserId,
		CollectType:    req.MsgType,
		CollectContent: req.Content,
	})

	if err != nil {
		return nil, errors.New("收藏失败")
	}

	return resp, nil
}

// 删除收藏消息
func (rpc *rpcChat) MsgCollectDel(_ context.Context, req *pbMsg.MsgCollectDelReq) (resp *pbMsg.MsgCollectDelResp, err error) {
	resp = &pbMsg.MsgCollectDelResp{}

	err = imdb.DelUserCollect(req.CollectId, req.UserId)
	fmt.Println("--DelUserCollect--", err)
	if err != nil {
		return nil, errors.New("删除收藏失败")
	}

	return resp, nil
}

// 收藏消息列表
func (rpc *rpcChat) MsgCollectList(_ context.Context, req *pbMsg.MsgCollectListReq) (resp *pbMsg.MsgCollectListResp, err error) {
	resp = &pbMsg.MsgCollectListResp{}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Size <= 0 {
		req.Size = 20
	}

	//获取分页列表数据
	list, total, err := imdb.FindUserCollectList(req.UserId, req.Keyword, req.MsgType, req.Page, req.Size)
	if err != nil {
		return resp, nil
	}

	resp.Total = total
	for _, v := range list {
		resp.MsgCollectList = append(resp.MsgCollectList, &pbMsg.MsgCollectList{
			MsgType:     v.CollectType,
			Content:     v.CollectContent,
			CollectTime: v.CreatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}
