CREATE TABLE `f_packet`
(
    `id`           int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id`      int(11) NOT NULL COMMENT '红包发起者',
    `packet_type`  tinyint(1) NOT NULL COMMENT '红包类型(1个人红包、2群红包)',
    `is_lucky`     tinyint(1) DEFAULT '0' COMMENT '是否为拼手气红包',
    `is_exclusive` tinyint(1) DEFAULT '0' COMMENT '是否为专属红包(0不是、1是)',
    `packet_title` varchar(100)   NOT NULL COMMENT '红包标题',
    `amount`       decimal(10, 2) NOT NULL COMMENT '红包金额',
    `number`       tinyint(3) NOT NULL COMMENT '红包个数',
    `expire_time`  datetime DEFAULT NULL COMMENT '红包过期时间',
    `created_time` datetime DEFAULT NULL,
    `updated_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY            `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户红包表';