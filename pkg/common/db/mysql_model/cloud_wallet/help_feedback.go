package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"time"
)

// 反馈

func InsertHelpFeedback(feedback *db.HelpFeedback) (err error) {
	//err := db.DB.MysqlDB.DefaultGormDB().Table("f_error_log").Create(log).Error
	feedback.AddTime = time.Now()
	feedback.UpdateTime = time.Now()
	err = db.DB.MysqlDB.DefaultGormDB().Table("help_feedback").Create(feedback).Error
	return err
}

// 我的反馈
func SelectHelpFeedbackByUserID(userID string) (feedbacks []db.HelpFeedback, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("help_feedback").Where("user_id = ?", userID).Find(&feedbacks).Error
	return feedbacks, err
}
