CREATE TABLE `f_trade`
(
    `id`               int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id`          int(11) NOT NULL COMMENT '用户id',
    `payment_platform` tinyint(1) NOT NULL COMMENT '支付平台(1云钱包、2支付宝、3微信、4银行卡)',
    `type`             tinyint(1) NOT NULL COMMENT '类型(1充值、2提现、3红包支出)',
    `amount`           decimal(10, 2) NOT NULL COMMENT '变更金额',
    `befer_amount`     decimal(10, 2) DEFAULT NULL COMMENT '变更前金额',
    `after_amount`     decimal(10, 2) DEFAULT NULL COMMENT '变更后金额',
    `third_order_no`   varchar(100)   DEFAULT NULL COMMENT '第三方订单号',
    `created_time`     datetime       DEFAULT NULL,
    `updated_time`     datetime       DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY                `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户账户变更表';