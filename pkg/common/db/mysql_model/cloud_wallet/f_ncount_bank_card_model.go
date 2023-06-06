package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"time"
)

// 获取银行卡列表
func GetUserBankcardByUserId(userID string) (list []*db.FNcountBankCard, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_bank_card").Where("user_id = ? and is_bind = ?", userID, 1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return
}

// 绑定用户银行卡
func BindUserBankcard(info *db.FNcountBankCard) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_bank_card").Create(info).Error
	if err != nil {
		return err
	}
	return nil
}

// 绑定用户银行卡确认
func BindUserBankcardConfirm(bankcardId int32, userId, bindCardAgrNo, bankCode string) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_bank_card").Where("id = ? and user_id = ?", bankcardId, userId).Updates(map[string]interface{}{
		"bind_card_agr_no": bindCardAgrNo, "bank_code": bankCode, "is_bind": 1, "updated_time": time.Now(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// 解绑用户银行卡
func UnBindUserBankcard(bankcardId int32, userId string) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_bank_card").Where("id = ? and user_id = ?", bankcardId, userId).Updates(map[string]interface{}{
		"is_delete": 1, "is_bind": 0, "updated_time": time.Now(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// 获取绑定的银行卡信息ById
func GetNcountBankCardById(id int32, userId string) (info *db.FNcountBankCard, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_bank_card").Where("id = ? and user_id = ?", id, userId).First(&info).Error
	if err != nil {
		return nil, err
	}
	return
}

// 获取绑定的银行卡信息ByBindCardAgrNo
func GetNcountBankCardByBindCardAgrNo(bindCardAgrNo, userId string) (info *db.FNcountBankCard, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_bank_card").Where("bind_card_agr_no = ? and user_id = ?", bindCardAgrNo, userId).First(&info).Error
	if err != nil {
		return nil, err
	}
	return
}
