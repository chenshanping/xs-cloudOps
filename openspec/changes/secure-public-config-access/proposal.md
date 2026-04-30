## Why

当前 `go-base` 的 `POST /api/v1/configs/keys` 是公开路由，但它会按客户端提交的任意 `keys` 直接查询并返回 `sys_config` 中的配置值。

这会带来两个直接问题：

- 匿名请求可以读取本应只在后台使用的敏感配置，例如 SMTP 凭据、对象存储配置、默认密码等。
- 前端启动阶段把登录页品牌配置和后台管理配置混在同一批 key 中加载，导致公开场景与后台场景的边界不清晰。

这次变更的目标不是重构整套配置中心，而是先把高风险入口收口到安全边界内，同时保持登录页、前台模式和后台配置页的现有可用性。

## What Changes

- 保留 `POST /api/v1/configs/keys` 作为现有配置批量读取入口，但改变它的访问语义：
  - 匿名请求只能读取服务端明确允许公开的配置键
  - 携带有效登录 Token 的请求可以继续读取后台所需的受保护配置键
- 将前端配置加载拆成两类：
  - 公开配置：应用启动、登录页、未登录前台使用
  - 后台配置：登录后后台页面再补充加载
- 为公开配置白名单补齐回归测试，防止未来把新的敏感键再次暴露到匿名接口

## Scope Boundaries

- 只覆盖 `sys_config` 的 HTTP 读取边界和前端配置加载方式
- 不修改后台配置编辑接口
- 不新增新的配置存储表或加密机制
- 不在本轮把 `/configs/keys` 拆成全新 public/private 两条独立接口
- 不扩展到限流、中间件级可选认证、配置项加密托管等系统性安全改造

## Affected Areas

- Backend:
  - `server/api/v1/config.go`
  - `server/service/configsvc/config.go`
  - `server/tests/`
- Frontend:
  - `web/src/store/config.ts`
  - 启动与登录后配置刷新调用链
- Specs:
  - 新增一个描述公开/受保护配置读取边界的 capability

## Rollback

- 后端移除公开白名单校验，恢复匿名请求按任意 key 读取配置
- 前端恢复单一配置 key 集合和统一加载逻辑

## Out of Scope

- 新增独立 `/public/configs` 接口
- 配置值加密存储或脱敏展示体系
- 配置访问审计日志和频率限制
- 后台配置模块 UI 改版
