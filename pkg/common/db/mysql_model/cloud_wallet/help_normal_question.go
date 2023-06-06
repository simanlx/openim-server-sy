package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"gorm.io/gorm"
)

// 获取常见问题
func GetHelpNormalQuestion() (questions []db.HelpNormalQuestion, err error) {
	// 限制5条 利用ord排序
	err = db.DB.MysqlDB.DefaultGormDB().Table("help_normal_question").Order("ord desc").Limit(5).Find(&questions).Error
	return questions, err
}

// 获取所有常见问题
func GetAllHelpNormalQuestion() (questions []db.HelpNormalQuestion, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("help_normal_question").Find(&questions).Error
	return questions, err
}

// 反馈常见问题
func UpdateHelpNormalQuestionFeedback(QuestionID, solved int64) (err error) {
	if solved == 0 {
		// 将soveld += 1
		err = db.DB.MysqlDB.DefaultGormDB().Table("help_normal_question").Where("id = ?", QuestionID).Update("solved", gorm.Expr("solved + ?", 1)).Error
	} else {
		// 将soveld += 1
		err = db.DB.MysqlDB.DefaultGormDB().Table("help_normal_question").Where("id = ?", QuestionID).Update("unsolved", gorm.Expr("unsolved + ?", 1)).Error
	}
	return err
}
