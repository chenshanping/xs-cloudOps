## Why

当前后台 AI 配置页把平台字段和模型列表混在同一个折叠表单里，平台切换、模型维护、模型导入都不够直观。随着需要从第三方 OpenAI 兼容平台拉取官方模型列表并按需导入，本页已经从“简单表单编辑”变成一个涉及前后端交互、权限补齐和管理员工作流闭环的配置能力。

现在补齐这个能力，可以让管理员用更稳定的“平台列表 + 模型面板 + 管理 Drawer”方式维护 AI 配置，并通过后端代理安全地发现平台模型，而不改变现有 `ai_config` 的落库方式和统一保存心智。

## What Changes

- 重构后台 AI 配置页为“左侧平台列表 + 右侧已导入模型主面板”的布局。
- 新增平台信息编辑 Drawer，用于新增、编辑、删除平台以及设置默认平台。
- 新增平台模型管理 Drawer，支持使用当前编辑态的 `api_key` 与 `base_url` 拉取平台模型列表、搜索、勾选并导入。
- 新增后端代理接口 `POST /api/v1/ai/providers/models/fetch`，由后端调用第三方 OpenAI 兼容平台的模型列表接口并返回清洗后的模型数据。
- 导入行为采用追加去重；导入仅更新当前编辑态，仍由现有“保存配置”流程统一落库。
- 补齐新接口的 AI 模块权限和启动补齐逻辑，并保证启动补齐不覆盖管理员已自定义的菜单或元数据。

### Out of Scope

- 不新增数据库表，不拆分 `ai_config` JSON 结构。
- 不实现自动同步、定时刷新或平台模型缓存。
- 不改动 AI 对话主流程、聊天页面功能或模型调用协议。
- 不引入新的代码生成器、外部管理依赖或全新的系统配置框架。

## Capabilities

### New Capabilities
- `admin-ai-config-management`: 管理员可以通过平台列表、编辑 Drawer 和模型管理 Drawer 安全维护 AI 平台配置，并从 OpenAI 兼容平台拉取模型列表后按需导入到本地配置。

### Modified Capabilities
- None. Repository has no existing OpenSpec capability for admin AI config behavior yet.

## Impact

- Backend:
  - `server/api/v1/ai.go`
  - `server/service/ai_client.go`
  - `server/router/modules/ai.go`
  - `server/model/request/ai_request.go`
  - `server/initialize/db_tables.go`
- Frontend:
  - `web/src/views/admin/system/config/components/AIConfig.vue`
  - `web/src/views/admin/ai/config/index.vue`
  - local `components/` under the AI config page for the two Drawers
  - `web/src/api/config.ts`
- API:
  - New authenticated admin-side AI provider model fetch endpoint
- Security:
  - Third-party model discovery must be proxied through backend and must not log raw `api_key`
- Rollback:
  - Remove the provider model fetch endpoint and Drawer-based model import flow
  - Restore the previous inline AI config editing layout while keeping existing `ai_config` data intact
