## Why

当前用户管理只支持单个用户状态切换，批量处理场景需要重复逐条操作，效率低且容易遗漏。同时，“身份”按钮当前固定展示，但业务上暂时希望默认隐藏，并由系统参数控制是否显示。

现在补齐这两个能力，可以让后台用户管理在批量启用/禁用场景下具备完整闭环，并把“身份”按钮是否展示收敛到系统配置，避免后续反复改页面代码。

## What Changes

- 新增后台用户管理批量状态修改能力，支持批量启用和批量禁用选中用户。
- 新增后端批量状态接口，并在服务层增加安全校验：
  - 禁止操作空列表
  - 禁止禁用当前操作者自己
  - 禁止批量修改受保护管理员账号
  - 批量禁用后让目标用户 token 立即失效
- 在后台用户管理页面增加批量启用、批量禁用工具栏操作和二次确认反馈。
- 新增系统配置项，用于控制后台用户管理页“身份”按钮是否显示，默认隐藏。
- 在系统参数设置页面增加该开关，允许管理员按需开启。
- 复用现有操作日志中间件记录批量状态请求体，确保审计能看到操作者和目标用户 ID 列表。

### Out of Scope

- 不调整前台个人中心里的身份认证流程。
- 不改动身份注册表、身份档案数据结构或身份审核逻辑。
- 不引入新的代码生成器、权限系统重构或用户管理大改版。

## Capabilities

### New Capabilities

- `admin-user-controls`: 管理员可以在用户管理中批量启用/禁用用户，并通过系统参数控制“身份”按钮是否显示。

### Modified Capabilities

- None. Repository has no existing OpenSpec capability for user admin behavior yet.

## Impact

- Backend:
  - `server/api/v1/user.go`
  - `server/service/user.go`
  - `server/model/request/request.go`
  - `server/router/modules/user.go`
  - `server/initialize/db_tables.go` or equivalent config seed path for default config values
- Frontend:
  - `web/src/api/user.ts`
  - `web/src/views/system/user/index.vue`
  - `web/src/views/system/config/components/SystemConfig.vue`
  - `web/src/store/config.ts`
- API:
  - New batch user status endpoint
- Audit:
  - Existing operation log middleware will capture operator identity and request body
- Rollback:
  - Remove the batch status endpoint and toolbar actions
  - Remove the new config toggle from system config UI and defaults
  - Restore the identity button to fixed hidden behavior if the feature needs to be backed out
