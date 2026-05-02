# server/tests

这里集中放置后端集成测试，优先包含以下类型：

- 依赖内存数据库的服务级测试
- 依赖本地文件系统或对象存储假实现的测试
- 跨多个 service 子包协作的测试

不放在这里的测试通常是白盒测试，原因是它们需要访问对应子包的未导出实现。
这类测试继续与源码同目录放置，例如：

- `server/service/ai/`
- `server/service/auth/`
- `server/service/user/`

常用单项运行命令：

- 角色权限 smoke 测试：`server\scripts\test_role_permission_smoke.bat`
- 跨平台等价命令：`go test -run TestRolePermissionSmoke -count=1 ./tests`
