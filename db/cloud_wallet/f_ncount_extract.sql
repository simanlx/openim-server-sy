CREATE TABLE `f_ncount_extract`
(
    `id`                int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id`           int(11) NOT NULL COMMENT '用户id',
    `bank_card_id`      int(11) NOT NULL COMMENT '银行卡id',
    `card_owner`        varchar(15)    NOT NULL COMMENT '持卡者',
    `bank_card_type`    tinyint(1) NOT NULL COMMENT '银行卡类型',
    `bank_card_number`  varchar(30)    NOT NULL COMMENT '银行卡号',
    `withdrawal_amount` decimal(10, 2) NOT NULL COMMENT '提现金额',
    `third_order_no`    varchar(100) DEFAULT NULL COMMENT '第三方订单号',
    `receipt_status`    tinyint(1) DEFAULT '0' COMMENT '到账状态(0未到账、1已到账)',
    `receipt_time`      datetime     DEFAULT NULL COMMENT '到账时间',
    `created_time`      datetime     DEFAULT NULL,
    `updated_time`      datetime     DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='提现记录';