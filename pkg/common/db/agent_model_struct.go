package db

import (
	"gorm.io/gorm"
	"time"
)

// 推广员申请表
type TAgentApplyRecord struct {
	Id          int64     `gorm:"column:id" json:"id"`
	UserId      string    `gorm:"column:user_id" json:"user_id"`             //用户id
	ChessUserId int64     `gorm:"column:chess_user_id" json:"chess_user_id"` //互娱用户id
	Name        string    `gorm:"column:name" json:"name"`                   //推广员姓名
	Mobile      string    `gorm:"column:mobile" json:"mobile"`               //推广员电话
	AuditStatus int32     `gorm:"column:audit_status" json:"audit_status"`   //审核状态
	Remark      string    `gorm:"column:remark" json:"remark"`               //备注
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
}

func (TAgentApplyRecord) TableName() string {
	return "t_agent_apply_record"
}

// 推广员账户表
type TAgentAccount struct {
	Id          int64     `gorm:"column:id" json:"id"`
	UserId      string    `gorm:"column:user_id" json:"user_id"`             //用户id
	Name        string    `gorm:"column:name" json:"name"`                   //推广员姓名
	Mobile      string    `gorm:"column:mobile" json:"mobile"`               //推广员电话
	ChessUserId int64     `gorm:"column:chess_user_id" json:"chess_user_id"` //互娱用户id
	AgentNumber int32     `gorm:"column:agent_number" json:"agent_number"`   //推广员编号
	Balance     int64     `gorm:"column:balance" json:"balance"`             //余额(单位:分)
	BeanBalance int64     `gorm:"column:bean_balance" json:"bean_balance"`   //咖豆余额
	OpenStatus  int32     `gorm:"column:open_status" json:"open_status"`     //开通状态(1开通、0关闭)
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
	DB          *gorm.DB  `gorm:"-" json:"-"`
}

func (TAgentAccount) TableName() string {
	return "t_agent_account"
}

// 推广员余额账户变更记录表
type TAgentAccountRecord struct {
	Id                int64     `json:"id"`
	OrderNo           string    `json:"order_no"`            // 订单号
	UserId            string    `json:"user_id"`             // 用户id
	Type              int32     `json:"type"`                // 收支类型(1收入、2支出)
	BusinessType      int32     `json:"business_type"`       // 业务类型(见枚举)
	Describe          string    `json:"describe"`            // 描述
	Amount            int32     `json:"amount"`              // 金额(单位:分)
	ChessUserId       int64     `json:"chess_user_id"`       // 互娱用户id
	ChessUserNickname string    `json:"chess_user_nickname"` // 互娱用户昵称
	Day               string    `json:"day"`                 // 天
	Month             string    `json:"month"`               // 月
	Status            int32     `json:"status"`              // 状态
	CreatedTime       time.Time `json:"created_time"`
	UpdatedTime       time.Time `json:"updated_time"`
	DB                *gorm.DB  `gorm:"-" json:"-"`
}

func (TAgentAccountRecord) TableName() string {
	return "t_agent_account_record"
}

// 推广员咖豆账户变更记录表
type TAgentBeanAccountRecord struct {
	Id                int64     `json:"id"`
	OrderNo           string    `json:"order_no"`            // 订单号
	UserId            string    `json:"user_id"`             // 用户id
	Type              int32     `json:"type"`                // 收支类型(1收入、2支出)
	BusinessType      int32     `json:"business_type"`       // 业务类型(见枚举)
	ChessUserId       int64     `json:"chess_user_id"`       // 互娱用户id
	ChessUserNickname string    `json:"chess_user_nickname"` // 互娱用户昵称
	Describe          string    `json:"describe"`            // 描述
	Amount            int32     `json:"amount"`              // 金额(单位分)
	Number            int64     `json:"number"`              // 数量
	GiveNumber        int32     `json:"give_number"`         // 赠送数量
	Day               string    `json:"day"`                 // 天
	CreatedTime       time.Time `json:"created_time"`
	UpdatedTime       time.Time `json:"updated_time"`
	DB                *gorm.DB  `gorm:"-" json:"-"`
}

func (TAgentBeanAccountRecord) TableName() string {
	return "t_agent_bean_account_record"
}

// 咖豆充值订单表
type TAgentBeanRechargeOrder struct {
	Id                int64     `json:"id"`
	BusinessType      int32     `json:"business_type"`       // 业务类型(见枚举)
	UserId            string    `json:"user_id"`             // 平台用户id
	ChessUserId       int64     `json:"chess_user_id"`       // 互娱用户id
	ChessUserNickname string    `json:"chess_user_nickname"` // 互娱用户昵称
	OrderNo           string    `json:"order_no"`            // 平台订单号
	ChessOrderNo      string    `json:"chess_order_no"`      // 互娱订单号
	NcountOrderNo     string    `json:"ncount_order_no"`     // 新生支付订单号
	Number            int64     `json:"number"`              // 充值数量
	GiveNumber        int32     `json:"give_number"`         // 赠送金额
	Amount            int32     `json:"amount"`              // 金额(单位:元)
	PayTime           int64     `json:"pay_time"`            // 支付时间
	PayStatus         int32     `json:"pay_status"`          // 支付状态
	CreatedTime       time.Time `json:"created_time"`
	UpdatedTime       time.Time `json:"updated_time"`
	DB                *gorm.DB  `gorm:"-" json:"-"`
}

func (TAgentBeanRechargeOrder) TableName() string {
	return "t_agent_bean_recharge_order"
}

// 推广咖豆商店配置表
type TAgentBeanShopConfig struct {
	Id             int32     `json:"id"`
	UserId         string    `json:"user_id"`          // 用户id
	BeanNumber     int64     `json:"bean_number"`      // 咖豆
	GiveBeanNumber int32     `json:"give_bean_number"` // 赠送咖豆(0不赠送)
	Amount         int32     `json:"amount"`           // 金额(单位:元)
	Status         int32     `json:"status"`           // 状态(0下架、1上架)
	CreatedTime    time.Time `json:"created_time"`
	UpdatedTime    time.Time `json:"updated_time"`
}

func (TAgentBeanShopConfig) TableName() string {
	return "t_agent_bean_shop_config"
}

// 推广员下属成员表
type TAgentMember struct {
	Id            int32     `json:"id"`
	UserId        string    `json:"user_id"`        // 用户id
	AgentNumber   int32     `json:"agent_number"`   // 推广员编号
	ChessUserId   int64     `json:"chess_user_id"`  // 互娱用户id
	ChessNickname string    `json:"chess_nickname"` // 互娱用户昵称
	Contribution  int64     `json:"contribution"`   // 成员贡献值
	Day           string    `json:"day"`            // 天
	CreatedTime   time.Time `json:"created_time"`
	UpdatedTime   time.Time `json:"updated_time"`
	DB            *gorm.DB  `gorm:"-" json:"-"`
}

func (TAgentMember) TableName() string {
	return "t_agent_member"
}

type TAgentPlatformConfig struct {
	Id          int32     `json:"id"`
	ConfigKey   string    `json:"config_key"`  // 配置类型
	ConfigData  string    `json:"config_data"` // 配置值
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

func (TAgentPlatformConfig) TableName() string {
	return "t_agent_platform_config"
}

type TAgentWithdraw struct {
	Id              int32     `json:"id"`
	OrderNo         string    `json:"order_no"`         // 平台订单号
	NcountOrderNo   string    `json:"ncount_order_no"`  // 新生支付订单号
	UserId          string    `json:"user_id"`          // 用户id
	AgentNumber     int32     `json:"agent_number"`     // 推广员编号
	BeforeBalance   int64     `json:"before_balance"`   // 提现前余额(单位:分)
	Balance         int32     `json:"balance"`          // 提现余额(单位:分)
	NcountBalance   int64     `json:"ncount_balance"`   // 新生支付到账金额(单位:分)
	TransferredTime int32     `json:"transferred_time"` // 到账时间
	Commission      int32     `json:"commission"`       // 手续费率千分几
	CommissionFee   int32     `json:"commission_fee"`   // 手续费(单位分)
	Status          int8      `json:"status"`           // 状态
	CreatedTime     time.Time `json:"created_time"`
	UpdatedTime     time.Time `json:"updated_time"`
}

func (TAgentWithdraw) TableName() string {
	return "t_agent_withdraw"
}
