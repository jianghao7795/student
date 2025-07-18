-- 创建错误信息表
CREATE TABLE `errors` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `error_code` int(11) NOT NULL COMMENT '错误码',
  `error_type` varchar(50) CHARACTER SET utf8mb4 NOT NULL COMMENT '错误类型',
  `error_message` varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '错误消息',
  `error_description` text CHARACTER SET utf8mb4 COMMENT '错误描述',
  `solution` text CHARACTER SET utf8mb4 COMMENT '解决方案',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_error_code_type` (`error_code`, `error_type`),
  KEY `idx_error_type` (`error_type`),
  KEY `idx_error_code` (`error_code`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='错误信息表';

-- 插入一些自定义错误信息
INSERT INTO `errors` (`error_code`, `error_type`, `error_message`, `error_description`, `solution`) VALUES
(2001, 'STUDENT', 'Student name already exists', '学生姓名已存在', '请使用不同的学生姓名'),
(2002, 'STUDENT', 'Student age invalid', '学生年龄无效', '请检查学生年龄是否在有效范围内'),
(2003, 'USER', 'Username already exists', '用户名已存在', '请使用不同的用户名'),
(2004, 'USER', 'Email already exists', '邮箱已存在', '请使用不同的邮箱地址'),
(2005, 'AUTH', 'Password too weak', '密码强度不够', '请使用包含字母、数字和特殊字符的强密码'),
(2006, 'VALIDATION', 'Invalid phone number', '手机号格式无效', '请使用正确的手机号格式'),
(2007, 'DATABASE', 'Database connection failed', '数据库连接失败', '请检查数据库配置或联系管理员'),
(2008, 'EXTERNAL', 'External service unavailable', '外部服务不可用', '请稍后重试或联系技术支持');