CREATE TABLE `f_ncount_bank_card` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT '主键',
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `mobile` varchar(15) DEFAULT NULL COMMENT '手机号码',
  `card_owner` varchar(10) DEFAULT NULL COMMENT '持卡者名字',
  `bank_card_type` tinyint(1) DEFAULT '0' COMMENT '银行卡类型(1-...枚举)',
  `bank_card_number` varchar(30) DEFAULT NULL COMMENT '银行卡号',
  `is_delete` tinyint(1) DEFAULT '0' COMMENT '是否删除(0未删除，1已删除)',
  `created_time` datetime DEFAULT NULL,
  `updated_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户银行卡绑定表';