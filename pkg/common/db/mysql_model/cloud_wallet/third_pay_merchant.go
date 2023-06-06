package cloud_wallet

import "crazy_server/pkg/common/db"

// 查询merchat

func GetMerchant(MerchantId string) (*db.ThirdPayMerchant, error) {
	resp := &db.ThirdPayMerchant{}
	result := db.DB.MysqlDB.DefaultGormDB().Table("third_pay_merchant").Where("merchant_id = ?", MerchantId).Find(resp)
	if result.Error != nil {
		return nil, result.Error
	}
	return resp, nil
}
