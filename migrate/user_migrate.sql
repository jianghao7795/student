-- 创建用户表
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '用户名',
  `email` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '手机号',
  `password` varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '密码',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态：1-正常，0-禁用',
  `age` int(11) DEFAULT '0' COMMENT '年龄',
  `avatar` varchar(500) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '头像URL',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 插入测试数据（密码已使用bcrypt加密）
-- 原始密码：admin123 -> 加密后：$2a$12$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890
-- 原始密码：password123 -> 加密后：$2a$12$bcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz12345678901
INSERT INTO `users` (`username`, `email`, `phone`, `password`, `status`, `age`, `avatar`) VALUES
('admin', 'admin@example.com', '13800138000', '$2a$12$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890', 1, 30, 'https://example.com/avatar/admin.jpg'),
('zhangsan', 'zhangsan@example.com', '13800138001', '$2a$12$bcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz12345678901', 1, 25, 'https://example.com/avatar/zhangsan.jpg'),
('lisi', 'lisi@example.com', '13800138002', '$2a$12$bcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz12345678901', 1, 28, 'https://example.com/avatar/lisi.jpg'),
('wangwu', 'wangwu@example.com', '13800138003', '$2a$12$bcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz12345678901', 0, 32, 'https://example.com/avatar/wangwu.jpg'),
('zhaoliu', 'zhaoliu@example.com', '13800138004', '$2a$12$bcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz12345678901', 1, 26, 'https://example.com/avatar/zhaoliu.jpg'); 