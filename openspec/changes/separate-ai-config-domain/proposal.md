## Why

当前后台 `AI配置` 菜单属于 AI 模块，但页面读写仍走通用 `配置管理` 接口，导致菜单权限、API 权限和页面归属长期割裂。继续沿用这条链路，会让角色授权、问题排查和后续数据结构演进都越来越难收口。

现在需要先把 AI 配置的接口归属收回 AI 模块，稳定菜单/API/权限语义，再把“是否拆 provider/model 表”留到单独第二阶段处理，避免把一次域归属修正和一次数据模型重构硬绑在同一轮。

## What Changes

- 新增 AI 模块专属的后台配置读写能力，由 AI 路由而不是通用 `配置管理` 路由对外提供 AI 配置读取、保存和相关管理接口。
- 将后台 AI 配置页面改为只依赖 AI 模块接口，不再依赖 `配置管理` API 作为主读写链路。
- 保留当前 `sys_config.key = ai_config` 的兼容存储方式，作为第一阶段的持久化桥接实现。
- 将 AI 配置菜单、AI 配置页面、AI 测试接口、平台模型发现接口统一收口为 AI 模块权限语义。
- 明确记录第二阶段边界：若需要解决模型规模、局部更新、并发覆盖和结构化查询问题，后续再单开 change 拆 `provider/model` 表。

### Out of Scope

- 本轮不新增数据库表，不拆分 `ai_config` JSON 结构。
- 本轮不迁移历史 AI 配置数据。
- 本轮不重构 AI 聊天主流程、会话存储或模型调用协议。
- 本轮不引入统一权限码或菜单/API 自动联动。

## Capabilities

### New Capabilities
- `admin-ai-config-domain`: 后台 AI 配置必须通过 AI 模块接口和 AI 模块权限对外暴露，并在第一阶段继续兼容现有 `ai_config` 存储。

### Modified Capabilities
- None. Repository source-of-truth specs do not yet include archived AI config behavior.

## Impact

- Backend:
  - `server/api/v1/ai.go`
  - `server/service/ai`
  - `server/service/configsvc`
  - `server/router/modules/ai.go`
  - `server/initialize/db_tables.go`
  - related request/response DTOs and tests
- Frontend:
  - AI 配置页及其本地 `components/`
  - `web/src/api/config.ts`
  - new or updated AI-specific frontend API module
- Permissions:
  - AI 配置相关角色授权将不再依赖 `配置管理` API
  - 角色权限聚合展示将以 AI 组接口为主，不再需要把 AI 页面挂到配置管理 API 下
- Rollback:
  - 恢复 AI 配置页对通用 `配置管理` 接口的依赖
  - 删除 AI 模块新增的专属配置读写接口
  - 保留现有 `ai_config` 数据，不做数据回滚
