## Context

当前聊天链路是 `web/src/components/AIChat.vue` -> `web/src/store/ai.ts` -> `/api/v1/ai/chat` 或 `/api/v1/ai/chat/stream` -> `server/service/ai_conversation.go` -> `server/service/ai_client.go`。前端“联网搜索”开关已经能传到后端，但后端只是把 `enable_search` 和提示词透传给上游 OpenAI 兼容接口，没有真实执行网页检索，所以模型仍会基于旧知识作答。

这次变更跨越前后端聊天链路并接入新的外部系统 Exa MCP，属于高风险集成改动。约束条件是：继续复用现有聊天 API 与页面入口，不新增数据库表，不破坏已存在的 AI 配置/模型选择逻辑，并且要保证流式与非流式路径行为一致。

## Goals / Non-Goals

**Goals:**
- 在开启“联网搜索”时，由后端先调用 Exa MCP 的 `web_search_exa`，拿到真实网页检索结果。
- 将检索结果压缩成受控的来源摘要，作为系统上下文注入当前模型请求，优先回答时效性问题。
- 当检索结果不足以支持结论时，模型回答必须明确说明“暂未确认”，而不是猜测补全。
- 把本轮使用的来源附加到最终回答中，让用户能直接看到和点开来源链接。
- 让普通聊天和流式聊天共用同一套 Exa grounding 逻辑，避免两条路径再次分叉。

**Non-Goals:**
- 不改 AI 配置页 UI，不为 Exa 单独新增管理页面。
- 不引入多搜索供应商抽象层，也不实现搜索供应商切换。
- 不新增数据库 schema，也不持久化单独的结构化 citation 表。
- 不实现后台定时同步、缓存索引、自动刷新搜索结果。

## Decisions

### 1. 使用标准库实现最小 MCP Streamable HTTP 客户端

选择在 `server/service/` 内新增一个小型 Exa MCP 客户端，直接用 Go 标准库发送 JSON-RPC over HTTP，而不是引入第三方 Go MCP SDK。

原因：
- Exa MCP 已验证支持标准 Streamable HTTP：`initialize`、`notifications/initialized`、`tools/list`、`tools/call`。
- 本次只需要同步调用两个工具，完整 SDK 的连接管理、双向请求、资源订阅都属于无关复杂度。
- 标准库实现更容易控制超时、错误脱敏、请求体裁剪和单元测试。

备选方案：
- 引入 Go MCP SDK：能力更全，但本次只用少量 JSON-RPC 请求，依赖成本和升级风险不成比例。
- 继续依赖上游模型自己的 search 选项：无法保证真实检索发生，也无法控制来源质量。

### 2. 保持无状态调用，兼容 MCP 初始化生命周期

实现时每次搜索都发送一次 `initialize`，随后发送 `notifications/initialized`，再发 `tools/call`。如果 Exa 返回 `Mcp-Session-Id`，客户端会在本次搜索内透传；如果没有返回，则按无状态流程继续。

原因：
- Exa MCP 当前对直接 POST 请求兼容良好，且服务端公开标记为 stateless。
- 遵循 MCP 生命周期可以降低未来服务端严格化后的兼容风险。
- 以“单次搜索一个短生命周期 client”处理，避免在当前仓库里引入长连接会话管理。

备选方案：
- 跳过初始化直接 `tools/call`：当前可能可用，但不符合 MCP 生命周期，兼容性差。
- 长驻 session 复用：状态更多、失败面更大，对当前同步聊天收益有限。

### 3. 搜索结果先做服务端压缩，再注入模型

后端不把 Exa 原始长文本直接喂给模型，而是先提取前几条结果的 `title`、`url`、`published`、高亮摘要，裁剪成受控的 grounding 文本块，并在系统消息中明确“只能依据这些来源回答时效问题”。

原因：
- 降低 token 膨胀和提示注入风险。
- 便于统一附加来源列表、控制最大结果数和最大字符数。
- 让模型侧输入稳定，可测试。

备选方案：
- 直接传 Exa 全量文本：实现最简单，但 token 不稳定，且更容易带入噪声。
- 搜索后再逐条 `web_fetch_exa` 深抓：准确性更高，但请求量和时延明显上升，本次先不做。

### 4. 来源以回答附录形式持久化进消息内容

不新增数据库字段保存结构化 citations，而是在最终助手回答后由服务端追加一段统一格式的“来源” Markdown 列表，例如 `来源：1. [标题](URL)`。这样流式与非流式返回、消息入库、历史会话重开都会保留来源。

原因：
- 满足“不新增 schema”的约束。
- 现有前端已支持 Markdown 链接展示，无需再改消息存储结构。
- 比只在当前 SSE 返回临时 sources 更稳，避免刷新会话后来源丢失。

备选方案：
- 新增 `sources` 响应字段但不落库：当前会话能看，刷新后丢失。
- 新增数据库字段：会超出本次 scope。

### 5. 聊天入口与权限不变

继续复用现有 `AI聊天` 页面和 `/api/v1/ai/chat`、`/api/v1/ai/chat/stream` 路由，不新增权限点。Exa search 只是原聊天功能在 `enable_search=true` 下的增强行为。

原因：
- 用户要求聚焦现有聊天行为，而不是额外做搜索平台页面。
- 现有权限模型已经覆盖 AI 聊天访问，无需引入新的 admin 配置权限。

## Risks / Trade-offs

- [Exa MCP 返回文本结构变化] → 解析逻辑基于明确的 `Title:` / `URL:` / `Highlights:` 段落，同时保留降级分支：解析失败时仍能把原始检索摘要作为 grounding 文本使用。
- [外部搜索增加延迟] → 限制结果数、限制字符数、设置请求超时；只在 `enable_search=true` 时触发。
- [搜索结果不权威或互相冲突] → 提示词要求优先采用官方/协会/赛事主办方来源；当结果冲突或没有权威来源时必须明确不确定。
- [无 Exa API Key 或上游限流] → 默认先支持匿名访问 Hosted MCP；若请求失败，向用户返回明确搜索失败信息，不静默回退到猜测回答。
- [来源附录影响回答排版] → 采用固定简短 Markdown 附录格式，只输出少量来源，避免淹没主体答案。

## Migration Plan

1. 在后端增加 Exa MCP 客户端与 grounded prompt 构造逻辑。
2. 将普通聊天与流式聊天统一切到“先 search，再 call model”的流程。
3. 让最终助手消息自动附加来源附录并沿用现有消息保存逻辑。
4. 运行 `cd server && go test ./...`，再运行前端构建验证来源展示没有破坏聊天页面。
5. 若上线后 Exa 搜索异常，可通过回滚代码版本恢复原聊天逻辑；由于未改 schema、未改路由，回滚成本低。

## Open Questions

- 当前版本先默认走 Hosted Exa MCP 的匿名能力；如果后续需要更高配额，再决定是否把 Exa API Key 纳入 AI 配置页。
- 当前版本使用 `web_search_exa` 即可完成闭环；是否补 `web_fetch_exa` 做二次抓取，留待后续根据真实问题质量再决定。
