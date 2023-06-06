package agent_model

import (
	"crazy_server/pkg/common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// 获取申请记录
func GetApplyByChessUserId(chessUserId int64) (info *db.TAgentApplyRecord, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_apply_record").Where("chess_user_id = ?", chessUserId).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}

// 获取申请记录
func GetApplyById(id int32) (info *db.TAgentApplyRecord, err error) {
	err = db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_apply_record").Where("id = ?", id).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "")
	}
	return
}

// 更新审核状态
func UpApplyAuditStatus(id int32) error {
	return db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_apply_record").Where("id = ?", id).Updates(map[string]interface{}{
		"audit_status": 1,
		"updated_time": time.Now(),
	}).Error
}

// 申请
func AgentApply(info *db.TAgentApplyRecord) error {
	info.CreatedTime = time.Now()
	info.UpdatedTime = time.Now()
	err := db.DB.AgentMysqlDB.DefaultGormDB().Table("t_agent_apply_record").Create(info).Error
	return err
}
