package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"github.com/pkg/errors"
)

/*
CREATE TABLE `f_ncount_account` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(64) NOT NULL COMMENT '用户id',
  `main_account_id` varchar(32) DEFAULT NULL COMMENT '主账号id',
  `packet_account_id` varchar(32) DEFAULT NULL COMMENT '红包账户id',
  `mobile` varchar(15) DEFAULT NULL COMMENT '手机号码',
  `real_auth` tinyint(1) DEFAULT '0' COMMENT '是否已实名认证',
  `realname` varchar(20) DEFAULT NULL COMMENT '真实姓名',
  `id_card` varchar(30) DEFAULT NULL COMMENT '身份证',
  `pay_switch` tinyint(4) DEFAULT '1' COMMENT '支付开关(0关闭、1默认开启)',
  `bod_pay_switch` tinyint(4) DEFAULT '0' COMMENT '指纹支付/人脸支付开关(0默认关闭、1开启)',
  `payment_password` varchar(32) DEFAULT NULL COMMENT '支付密码(md5加密)',
  `open_status` tinyint(4) DEFAULT '0' COMMENT '开通状态',
  `open_step` tinyint(4) DEFAULT '1' COMMENT '开通认证步骤(1身份证认证、2支付密码、3绑定银行卡或者刷脸)',
  `created_time` datetime DEFAULT NULL,
  `updated_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COMMENT='云钱包账户表';
*/

func FNcountAccountGetUserAccountID(userId string) (*db.FNcountAccount, error) {
	var account *db.FNcountAccount
	err := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_account").
		Select("main_account_id", "packet_account_id").Where("user_id = ?", userId).First(&account).Error
	if err != nil {
		return nil, errors.Wrap(err, "FNcountAccountGetUserAccountID error")
	}
	return account, nil
}

// 获取用户账户信息
func GetNcountAccountByUserId(userID string) (info *db.FNcountAccount, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_account").Where("user_id = ?", userID).First(&info).Error
	if err != nil {
		return nil, err
	}
	return
}

// 创建云钱包账户
func CreateNcountAccount(info *db.FNcountAccount) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_account").Create(&info).Error
	if err != nil {
		return err
	}
	return nil
}

// 更新账户信息
func UpdateNcountAccountField(userId string, m map[string]interface{}) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_account").Where("user_id = ?", userId).Updates(m).Error
	return err
}

// 单个身份证实名总数
func IdCardRealNameAuthNumber(idCard string) (count int64) {
	_ = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_account").Where("id_card = ?", idCard).Count(&count).Error
	return
}
