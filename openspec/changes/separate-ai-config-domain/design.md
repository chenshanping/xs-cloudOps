## Context

当前仓库已经完成了 AI 配置页的 provider/model 交互优化，但页面主保存链路仍落在 `sys_config` 通用配置接口上。实际效果是：

- 菜单归属在 AI 模块
- 页面保存权限归属在 `配置管理`
- 模型测试与模型发现接口又归属在 AI 模块

这会让角色授权必须同时勾 AI 菜单权限、配置管理 API 权限和 AI API 权限，页面语义与权限语义长期不一致。与此同时，`ai_config` 仍以单条 JSON 存在 `sys_config.value(TEXT)` 中，已经暴露出容量、并发和局部更新边界，但这些问题不适合和接口归属修正绑在同一轮处理。

## Goals / Non-Goals

**Goals:**

- 在第一阶段把 AI 配置的对外接口归属统一到 AI 模块。
- 让后台 AI 配置页面只依赖 AI 模块接口完成加载、保存、测试和模型发现。
- 保留现有 `ai_config` JSON 存储，避免本轮触发 SQL 迁移和数据搬迁。
- 让 AI 配置相关 API 元数据、Casbin 授权和角色权限展示更符合 AI 模块直觉。

**Non-Goals:**

- 不拆 `ai_provider` / `ai_model` 表。
- 不在本轮解决 `TEXT` 容量、局部更新、并发覆盖和结构化查询问题。
- 不把 AI 配置迁移到统一配置中心或新的 JSON 列类型。
- 不重做现有 AI 配置页的主交互布局。

## Decisions

### 1. 第一阶段只改接口归属，不改底层存储

决策：

- 新增 AI 模块专属配置接口，例如 `GET /api/v1/ai/config` 与 `PUT /api/v1/ai/config`。
- 由 AI 模块服务层负责读写 `sys_config.key = ai_config`，前端不再直接使用通用 `配置管理` 接口读写 AI 配置。

原因：

- 先解决“菜单/API 归属错位”的主要矛盾。
- 不引入 SQL 迁移，可以降低当前分支的风险和回滚成本。

备选方案：

- 直接拆 provider/model 表：能一次性解决更多问题，但会同时引入接口、权限、页面、迁移四类变化，当前轮次风险过高。

### 2. AI 配置读写接口沿用 AI 模块权限，不再借道配置管理

决策：

- AI 配置页的读取、保存、测试模型、拉取模型都统一挂在 AI 模块权限名义下。
- 初始化补齐逻辑把这些接口元数据和 admin/system_admin 角色授权收口到 AI 模块。

原因：

- 用户在角色授权时能直接按 AI 模块理解权限，不再需要跨到“配置管理”里找 AI 页面所需的主保存接口。
- 角色权限抽屉的页面聚合逻辑也会更自然。

备选方案：

- 继续让 AI 页面保存走 `配置管理`，只在前端做分组映射补丁：能暂时缓解展示问题，但不能修复真实权限归属错位。

### 3. 第一阶段保留现有 AI 配置 JSON 结构

决策：

- `ai_config` 继续包含 `default_provider`、`providers[]`、`models[]`。
- AI 专属读写接口只作为兼容桥接层，不改变 JSON 结构和现有页面状态结构。

原因：

- 可以复用现有前端编辑逻辑和后端 `GetAIConfig` 读取逻辑。
- 让本轮聚焦在“域归属修正”，不扩展到“数据模型重构”。

备选方案：

- 在第一阶段先把 JSON 改成多条 config：仍然不能从根本上解决关系化能力，且会增加迁移复杂度。

### 4. 第二阶段单独处理结构化存储

决策：

- 在设计文档里明确记录第二阶段应单开 change。
- 第二阶段才评估 `ai_provider` / `ai_model` 以及关联字段、启停状态、排序和扩展字段。

原因：

- `sys_config.value(TEXT)` 的容量和并发问题是真问题，但与当前“权限/API 归属修正”不是同一个最小闭环。
- 分阶段能让第一阶段先交付可用、清晰的一致权限语义。

## Risks / Trade-offs

- [仍保留 `TEXT + JSON` 存储] → 第一阶段接受该限制，并在 change 中显式记录为后续拆表触发条件。
- [前后端会同时存在旧 config API 和新 AI API] → 只让 AI 配置页切到新接口，避免误用旧链路；后续可单独评估是否下线旧访问方式。
- [角色权限历史配置可能仍残留配置管理 API] → 新页面逻辑和推荐授权路径转向 AI 模块，但不强制清理历史角色。
- [现有测试覆盖主要围绕页面交互，不一定覆盖域归属] → 本轮必须新增 AI 配置专属接口和权限链路的后端/前端验证。

## Migration Plan

1. 新增 AI 模块专属配置读取和保存接口。
2. 前端 AI 配置页切换到 AI 模块接口。
3. 补齐 AI 配置相关 API 元数据和默认角色授权。
4. 验证旧的 `ai_config` 数据无需迁移即可继续使用。
5. 如需回滚：
   - 前端改回调用通用 `配置管理` 接口
   - 删除 AI 模块新增配置接口
   - 保留 `sys_config.ai_config` 数据原样不动

## Open Questions

- None for phase 1. The user has already confirmed:
  - phase 1 only realigns AI config endpoints and permissions
  - provider/model table split belongs to a separate future phase
