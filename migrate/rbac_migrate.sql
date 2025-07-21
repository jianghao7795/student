-- 创建角色表
CREATE TABLE `roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 NOT NULL COMMENT '角色名称',
  `description` varchar(500) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '角色描述',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态：1-启用，0-禁用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_name` (`name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 创建权限表
CREATE TABLE `permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 NOT NULL COMMENT '权限名称',
  `resource` varchar(200) CHARACTER SET utf8mb4 NOT NULL COMMENT '资源路径',
  `action` varchar(50) CHARACTER SET utf8mb4 NOT NULL COMMENT '操作类型：GET,POST,PUT,DELETE等',
  `description` varchar(500) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '权限描述',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态：1-启用，0-禁用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_permission_resource_action` (`resource`, `action`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 创建用户角色关联表
CREATE TABLE `user_roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `role_id` int(11) NOT NULL COMMENT '角色ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`),
  CONSTRAINT `fk_user_roles_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 创建角色权限关联表
CREATE TABLE `role_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL COMMENT '角色ID',
  `permission_id` int(11) NOT NULL COMMENT '权限ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`),
  CONSTRAINT `fk_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 插入基础角色数据
INSERT INTO `roles` (`name`, `description`, `status`) VALUES
('admin', '系统管理员，拥有所有权限', 1),
('manager', '部门经理，拥有部门管理权限', 1),
('user', '普通用户，拥有基本权限', 1),
('guest', '访客，只有查看权限', 1);

-- 插入基础权限数据
INSERT INTO `permissions` (`name`, `resource`, `action`, `description`, `status`) VALUES
-- 用户管理权限
('user:read', '/api/v1/users', 'GET', '查看用户列表', 1),
('user:create', '/api/v1/users', 'POST', '创建用户', 1),
('user:update', '/api/v1/users/*', 'PUT', '更新用户信息', 1),
('user:delete', '/api/v1/users/*', 'DELETE', '删除用户', 1),
('user:detail', '/api/v1/users/*', 'GET', '查看用户详情', 1),

-- 学生管理权限
('student:read', '/api/v1/students', 'GET', '查看学生列表', 1),
('student:create', '/api/v1/students', 'POST', '创建学生', 1),
('student:update', '/api/v1/students/*', 'PUT', '更新学生信息', 1),
('student:delete', '/api/v1/students/*', 'DELETE', '删除学生', 1),
('student:detail', '/api/v1/students/*', 'GET', '查看学生详情', 1),

-- 角色管理权限
('role:read', '/api/v1/roles', 'GET', '查看角色列表', 1),
('role:create', '/api/v1/roles', 'POST', '创建角色', 1),
('role:update', '/api/v1/roles/*', 'PUT', '更新角色信息', 1),
('role:delete', '/api/v1/roles/*', 'DELETE', '删除角色', 1),
('role:detail', '/api/v1/roles/*', 'GET', '查看角色详情', 1),

-- 权限管理权限
('permission:read', '/api/v1/permissions', 'GET', '查看权限列表', 1),
('permission:create', '/api/v1/permissions', 'POST', '创建权限', 1),
('permission:update', '/api/v1/permissions/*', 'PUT', '更新权限信息', 1),
('permission:delete', '/api/v1/permissions/*', 'DELETE', '删除权限', 1),
('permission:detail', '/api/v1/permissions/*', 'GET', '查看权限详情', 1);

-- 为admin角色分配所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) 
SELECT 1, id FROM `permissions` WHERE status = 1;

-- 为manager角色分配用户和学生管理权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) 
SELECT 2, id FROM `permissions` WHERE name IN ('user:read', 'user:create', 'user:update', 'user:detail', 'student:read', 'student:create', 'student:update', 'student:detail');

-- 为user角色分配基本查看权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) 
SELECT 3, id FROM `permissions` WHERE name IN ('user:read', 'user:detail', 'student:read', 'student:detail');

-- 为guest角色分配只读权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) 
SELECT 4, id FROM `permissions` WHERE name IN ('user:read', 'student:read');

-- 为现有用户分配角色
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES
(1, 1), -- admin用户分配admin角色
(2, 2), -- zhangsan用户分配manager角色
(3, 3), -- lisi用户分配user角色
(4, 4), -- wangwu用户分配guest角色
(5, 3); -- zhaoliu用户分配user角色 