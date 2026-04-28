## Why

当前 AI 聊天里的“联网搜索”只是把 `enable_search` 和提示词透传给上游模型，没有真正执行外部检索、来源筛选和结果注入，所以面对 `2026 伦敦世乒赛男团名单` 这类强时效问题时，模型仍然会回退到旧知识并给出错误名单。现在需要把联网搜索改成后端先走 Exa MCP 检索，再把受控来源喂给模型，确保回答以检索结果为准。

这次改动只聚焦 AI 聊天搜索链路本身，不改 AI 配置页、数据库结构，也不引入多搜索供应商平台化抽象。目标是尽快补齐一个可验证、可回滚的最小闭环。

## What Changes

- 在 `server/` 新增 Exa MCP HTTP 客户端，用 MCP `initialize` / `tools/call` 访问 `https://mcp.exa.ai/mcp`。
- 当 AI 聊天开启“联网搜索”时，后端先调用 Exa `web_search_exa`，筛出少量搜索结果摘要与来源，再将其拼入对模型的系统上下文。
- 为聊天回答增加“来源”回传能力，支持普通请求与流式请求把来源数据返回给前端展示。
- 前端 AI 聊天气泡下展示本轮回答使用的来源链接，不改变原有聊天入口和配置心智。
- 保留当前 `ai_config` 与模型调用配置；Exa API Key 通过现有 AI 配置扩展字段或配置读取接入，不新增数据库表。
- 明确无权威来源时的降级行为：回答必须说明“暂未确认”，不能猜测补全。
- 回滚方式保持简单：保留原有非 Exa 搜索关闭路径，必要时可通过配置禁用 Exa 联网搜索逻辑并恢复为无搜索回答。

## Capabilities

### New Capabilities
- `ai-chat-grounded-web-search`: AI 聊天在启用联网搜索时，通过 Exa MCP 执行真实网页检索、注入受控来源，并向前端返回可展示的来源信息。

### Modified Capabilities
- None.

## Impact

- Backend: `server/service/ai_client.go`, `server/service/ai_conversation.go`, 新增 Exa MCP 搜索客户端与测试。
- Frontend: `web/src/api/ai.ts`, `web/src/store/ai.ts`, `web/src/components/AIChat.vue` 的聊天响应与来源展示。
- API: AI 聊天普通响应和流式响应会增加来源字段，但不更改现有请求入参。
- Dependencies: 继续使用标准库 HTTP 访问 MCP，不强制引入新的 Go MCP SDK。
- Out of scope: AI 配置页重构、数据库 schema 调整、多搜索引擎抽象、后台定时同步搜索索引。
