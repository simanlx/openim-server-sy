package cloud_wallet

import (
	"crazy_server/pkg/common/db"
	"crazy_server/pkg/common/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

//CREATE TABLE `f_packet` (
//`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`packet_id` varchar(255) DEFAULT NULL COMMENT '红包ID',
//`submit_time` varchar(255) DEFAULT NULL COMMENT '下单时间，用于退款',
//`user_id` varchar(255) NOT NULL COMMENT '红包发起者',
//`user_redpacket_account` varchar(255) DEFAULT NULL COMMENT '发送红包的用户的账户',
//`packet_type` tinyint(1) NOT NULL COMMENT '红包类型(1个人红包、2群红包)',
//`is_lucky` tinyint(1) DEFAULT '0' COMMENT '是否为拼手气红包',
//`exclusive_user_id` varchar(255) DEFAULT '0' COMMENT '专属用户id',
//`packet_title` varchar(100) NOT NULL COMMENT '红包标题',
//`amount` int(11) NOT NULL COMMENT '红包金额',
//`number` tinyint(3) NOT NULL COMMENT '红包个数',
//`expire_time` int(11) DEFAULT NULL COMMENT '红包过期时间',
//`mer_order_id` varchar(255) DEFAULT NULL COMMENT '红包第三方的请求ID',
//`operate_id` varchar(255) DEFAULT NULL COMMENT '链路追踪ID',
//`recv_id` varchar(255) DEFAULT NULL COMMENT '被发送用户的ID',
//`send_type` tinyint(11) DEFAULT NULL COMMENT '红包发送方式： 1：钱包余额，2是银行卡',
//`bind_card_agr_no` varchar(255) DEFAULT NULL COMMENT '银行卡绑定协议号',
//`remain` int(11) DEFAULT NULL COMMENT '剩余红包数量',
//`remain_amout` int(11) NOT NULL DEFAULT '0' COMMENT '剩余红包金额',
//`lucky_user_id` varchar(255) NOT NULL DEFAULT '' COMMENT '最佳手气红包用户ID',
//`luck_user_amount` int(11) NOT NULL DEFAULT '0' COMMENT '最大红包的值： account amount  分为单位',
//`created_time` int(11) DEFAULT NULL,
//`updated_time` int(11) DEFAULT NULL,
//`status` tinyint(1) NOT NULL COMMENT '红包状态： 1 为创建 、2 为正常、3为异常',
//`is_exclusive` tinyint(1) NOT NULL COMMENT '是否为专属红包： 0为否，1为是',
//PRIMARY KEY (`id`),
//KEY `idx_user_id` (`user_id`) USING BTREE,
//KEY `idx_packet_id` (`packet_id`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=294 DEFAULT CHARSET=utf8mb4 COMMENT='用户红包表';

const (
	RedPacketStatusCreate   = iota // 0 创建未生效
	RedPacketStatusNormal          // 1 为红包正在领取中
	RedPacketStatusFinished        // 2为红包领取完毕
	RedPacketStatusExpired         // 3为红包过期
)

type RedPacketDetail struct {
	UserId          string `json:"user_id"`
	PacketType      int32  `json:"packet_type"`
	IsLucky         int32  `json:"is_lucky"`
	IsExclusive     int32  `json:"is_exclusive"`
	ExclusiveUserID string `json:"exclusiveUser_id"`
	PacketTitle     string `json:"packet_title"`
	Amount          int64  `json:"amount"`
	Number          int32  `json:"number"`
	ExpireTime      int64  `json:"expire_time"`
	Remain          int64  `json:"remain"`
	Nickname        string `json:"nickname"`
	FaceUrl         string `json:"face_url"`
}

// 保存到红包到数据库
func RedPacketCreateData(req *db.FPacket) error {
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_packet").Create(req)
	if result.Error != nil {
		return errors.Wrap(result.Error, "创建红包失败")
	}
	return nil
}

// 修改红包状态
func UpdateRedPacketStatus(packetID string, status int64) error {
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_packet").Where("packet_id = ?", packetID).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(result.Error, "修改红包状态失败")
	}
	return nil
}

// 根据红包ID获取红包信息
func GetRedPacketInfo(packetID string) (*db.FPacket, error) {
	var fPacket db.FPacket
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_packet").Where("packet_id = ?", packetID).First(&fPacket)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "获取红包信息失败")
	}
	return &fPacket, nil
}

// 通过红包ID获取红包发送者的红包账户
func SelectRedPacketSenderRedPacketAccountByPacketID(packetID string) (string, error) {
	// 查询到发送红包的信息
	var fPacket db.FPacket
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_packet").Where("packet_id = ?", packetID).First(&fPacket)
	if result.Error != nil {
		log.Debug("查询红包信息: ", packetID)
		return "", errors.Wrap(result.Error, "获取红包信息失败: "+packetID)
	}
	sendUserID := fPacket.UserID

	// 查询用户的红包账户 : 查找用户的ID
	var fAccount db.FNcountAccount
	result = db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_account").Where("user_id = ?", sendUserID).First(&fAccount)
	if result.Error != nil {
		return "", errors.Wrap(result.Error, "获取红包信息失败")
	}
	return fAccount.PacketAccountId, nil
}

// 通过红包ID查询到 发送者的用户ID
func SelectUserMainAccountByUserID(userID string) (string, error) {
	var fAccount db.FNcountAccount
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_ncount_account").Where("user_id = ?", userID).First(&fAccount)
	if result.Error != nil {
		return "", errors.Wrap(result.Error, "获取红包信息失败")
	}
	return fAccount.MainAccountId, nil
}

// 更新红包的剩余数量 （如果红包数量没有了）
func UpdateRedPacketRemain(packetID string) error {
	// 将红包数量减一
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_packet").Where("packet_id = ?", packetID).Update("remain", gorm.Expr("remain - ?", 1))
	if result.Error != nil {
		return errors.Wrap(result.Error, "修改红包状态失败")
	}
	return nil
}

// 修改红包信息
func UpdateRedPacketInfo(packetID string, req *db.FPacket) error {
	req.UpdatedTime = time.Now().Unix()
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_packet").Where("packet_id = ?", packetID).Updates(req)
	if result.Error != nil {
		return errors.Wrap(result.Error, "修改红包状态失败")
	}
	return nil
}

// 查询过期红包
func GetExpiredRedPacketList() ([]*db.FPacket, error) {
	var fPacketList []*db.FPacket
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_packet").Where("expire_time < ? and status = ?", time.Now().Unix(), RedPacketStatusNormal).Find(&fPacketList)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "查询过期红包失败")
	}
	return fPacketList, nil
}

// 查询过期的红包： 过期时间小于当前时间，状态为正常，红包剩余余额大于0， 每次查询100条
func GetExpiredRedPacketListByPage() ([]*db.FPacket, error) {
	var fPacketList []*db.FPacket
	result := db.DB.MysqlDB.DefaultGormDB().Table("f_packet").Where("expire_time < ? and status = ? and remain > ? and remain_amout > 0", time.Now().Unix(), RedPacketStatusNormal, 0).Limit(1).Find(&fPacketList)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "查询过期红包失败")
	}
	return fPacketList, nil
}
