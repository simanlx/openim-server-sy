CREATE TABLE `f_ncount_account` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int(11) NOT NULL COMMENT '用户id',
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
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT='云钱包账户表';