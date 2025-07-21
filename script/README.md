# Scripts 目录

本目录包含项目相关的各种脚本文件。

## 脚本列表

### test_getme.sh

- **用途**: 测试 GetMe API 接口功能
- **功能**:
  - 自动登录获取 JWT token
  - 使用 token 调用 `/v1/user/me` 接口
  - 验证接口响应
- **使用方法**: `./test_getme.sh`

### run_tests.sh

- **用途**: 运行项目测试套件
- **功能**:
  - 运行单元测试
  - 运行集成测试
  - 运行端到端测试
- **使用方法**: `./run_tests.sh`

## 使用说明

1. 确保脚本具有执行权限：

   ```bash
   chmod +x script/*.sh
   ```

2. 在项目根目录下运行脚本：

   ```bash
   # 测试 GetMe 接口
   ./script/test_getme.sh

   # 运行测试套件
   ./script/run_tests.sh
   ```

## 注意事项

- 运行脚本前请确保项目已正确配置
- 某些脚本可能需要数据库连接或其他依赖服务
- 建议在测试环境中运行脚本
