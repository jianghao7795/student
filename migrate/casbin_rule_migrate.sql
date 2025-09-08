-- 创建casbin_rule表
-- 这个表用于存储Casbin的权限策略和角色分配规则
CREATE TABLE `casbin_rule` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL COMMENT 'policy type: p or g',
  `v0` varchar(100) DEFAULT NULL COMMENT 'first parameter',
  `v1` varchar(100) DEFAULT NULL COMMENT 'second parameter',
  `v2` varchar(100) DEFAULT NULL COMMENT 'third parameter',
  `v3` varchar(100) DEFAULT NULL COMMENT 'fourth parameter',
  `v4` varchar(100) DEFAULT NULL COMMENT 'fifth parameter',
  `v5` varchar(100) DEFAULT NULL COMMENT 'sixth parameter',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`),
  KEY `idx_ptype` (`ptype`),
  KEY `idx_v0` (`v0`),
  KEY `idx_v1` (`v1`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Casbin policy rules table';

-- 插入策略规则 (p = 策略)
-- 从rbac_policy.csv转换而来
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
-- 管理员权限
('p', 'admin', '/v1/users', '*'),
('p', 'admin', '/v1/students', '*'),

-- 用户权限
('p', 'user', '/v1/users', 'GET'),
('p', 'user', '/v1/students', 'GET'),

-- 访客权限
('p', 'guest', '/v1/users', 'GET');

-- 插入角色分配规则 (g = 角色分组)
-- 从rbac_policy.csv转换而来
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`) VALUES
-- 用户角色分配
('g', '1', 'admin'),
('g', '2', 'user'),
('g', '3', 'guest');

