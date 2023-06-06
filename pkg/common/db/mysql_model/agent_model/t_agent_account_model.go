package agent_model

import (
	"crazy_server/pkg/common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// 获取推广员信息AgentNumber
func GetAgentByAgentNumber(agentNumber int32) (info *db.TAgentAccount, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_account").Where("agent_number = ?", agentNumber).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}

// 获取推广员信息ByChessUserId
func GetAgentByChessUserId(chessUserId int64) (info *db.TAgentAccount, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_account").Where("chess_user_id = ?", chessUserId).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}

// 获取推广员信息ByUserId
func GetAgentByUserId(userId string) (info *db.TAgentAccount, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_account").Where("user_id = ?", userId).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}

// 创建推广员账户
func CreateAgentAccount(info *db.TAgentAccount) error {
	info.CreatedTime = time.Now()
	info.UpdatedTime = time.Now()
	err := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_account").Create(info).Error
	return err
}
