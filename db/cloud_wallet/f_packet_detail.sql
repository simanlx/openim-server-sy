CREATE TABLE `f_packet_detail`
(
    `id`           int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `packet_id`    int(11) NOT NULL COMMENT '红包id',
    `user_id`      int(11) DEFAULT NULL COMMENT '用户id',
    `amount`       decimal(10, 2) NOT NULL COMMENT '领取金额',
    `receive_time` datetime DEFAULT NULL COMMENT '领取时间',
    `is_receive`   tinyint(1) DEFAULT '0' COMMENT '是否已领取(0未领取、1已领取)',
    `is_exclusive` tinyint(1) DEFAULT '0' COMMENT '是否为专属红包(0不是、1是)',
    `created_time` datetime DEFAULT NULL,
    `updated_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY            `idx_packet_id` (`packet_id`) USING BTREE,
    KEY            `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='红包领取记录';