# 文档目录

本目录包含项目的所有文档文件。

## 文档结构

### API 相关文档

#### API_VERSIONING_CHANGES.md

- **描述**: API 版本化更改总结
- **内容**:
  - 接口路径从 `/api/v1` 改为 `/v1` 的详细说明
  - 修改的文件列表
  - 影响范围和兼容性说明

#### GET_ME_API_USAGE.md

- **描述**: GetMe API 接口使用指南
- **内容**:
  - 接口详细说明
  - 使用示例（curl、JavaScript、Python）
  - 错误处理和安全注意事项

### RBAC 权限系统文档

#### RBAC_README.md

- **描述**: RBAC 权限系统主要文档
- **内容**:
  - 系统概述和架构
  - API 接口说明
  - 使用示例和配置说明

#### RBAC_IMPLEMENTATION_SUMMARY.md

- **描述**: RBAC 实现总结
- **内容**:
  - 实现概述
  - 核心组件说明
  - 技术特点

#### RBAC_POLICY_README.md

- **描述**: RBAC 策略配置文档
- **内容**:
  - 策略文件格式
  - 配置示例
  - 管理方法

#### rbac_model_documentation.md

- **描述**: RBAC 模型详细文档
- **内容**:
  - Casbin 模型配置详解
  - 权限检查逻辑
  - 示例和最佳实践

#### rbac_model_quick_reference.md

- **描述**: RBAC 模型快速参考
- **内容**:
  - 模型语法速查
  - 常用配置示例
  - 故障排除

#### rbac_policy_documentation.md

- **描述**: RBAC 策略详细文档
- **内容**:
  - 策略文件格式详解
  - 权限定义方法
  - 高级配置选项

#### rbac_policy_quick_reference.md

- **描述**: RBAC 策略快速参考
- **内容**:
  - 策略语法速查
  - 常用权限配置
  - 快速配置指南

#### rbac_policy_template.csv

- **描述**: RBAC 策略模板文件
- **用途**: 作为创建新策略文件的模板

### 技术迁移文档

#### INTERFACE_TO_ANY_MIGRATION.md

- **描述**: interface{} 到 any 的迁移总结
- **内容**:
  - 迁移概述
  - 修改的文件列表
  - 技术说明和验证结果

## 文档使用指南

### 新用户入门

1. 首先阅读 `RBAC_README.md` 了解系统架构
2. 查看 `API_VERSIONING_CHANGES.md` 了解 API 规范
3. 参考 `GET_ME_API_USAGE.md` 学习 API 使用

### 开发者参考

1. `RBAC_IMPLEMENTATION_SUMMARY.md` - 了解实现细节
2. `rbac_model_documentation.md` - 深入理解权限模型
3. `rbac_policy_documentation.md` - 掌握策略配置

### 快速查找

1. `rbac_model_quick_reference.md` - 模型语法速查
2. `rbac_policy_quick_reference.md` - 策略配置速查
3. `INTERFACE_TO_ANY_MIGRATION.md` - 代码迁移参考

## 文档维护

- 所有文档使用 Markdown 格式
- 保持文档与代码同步更新
- 新增功能时及时更新相关文档
- 定期检查和更新文档链接
