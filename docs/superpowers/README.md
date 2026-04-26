# go-base Superpowers 工作流

这个目录是本仓库的 Superpowers 落地点。

## 目录

- `docs/superpowers/specs/`
  - 存放通过 `brainstorming` 产出的设计文档
- `docs/superpowers/plans/`
  - 存放通过 `writing-plans` 产出的实现计划

## 推荐执行顺序

### 阶段一：探索

目标：先把需求、边界、风险和方案讲清楚。

使用：

- `using-superpowers`
- `brainstorming`

产物：

- `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md`

### 阶段二：锁定规格

目标：把确认后的内容落到 OpenSpec 中，形成可追溯变更。

使用：

- `openspec-explore`
- `openspec-propose`

产物：

- `openspec/changes/<change-name>/proposal.md`
- `openspec/changes/<change-name>/design.md`
- `openspec/changes/<change-name>/specs/...`
- `openspec/changes/<change-name>/tasks.md`

### 阶段三：执行

目标：按规格和计划实施，而不是边写边改需求。

使用：

- `writing-plans`
- `using-git-worktrees`
- `subagent-driven-development` 或 `executing-plans`
- `test-driven-development`
- `verification-before-completion`
- `openspec-apply-change`

产物：

- `docs/superpowers/plans/YYYY-MM-DD-<feature-name>.md`
- 代码、测试、已打勾的 `tasks.md`

### 阶段四：收尾

目标：验证实现和规格一致，并归档变更。

使用：

- `verification-before-completion`
- `requesting-code-review`
- `finishing-a-development-branch`
- `openspec-archive-change`

## 什么时候必须走这套流程

- 新业务功能
- 接口行为变化
- 权限、审计、认证、数据一致性相关修改
- 跨后端/前端/数据库的联动改造
- 需要留痕、评审、可恢复上下文的工作

## 什么时候可以轻量化

- 文档润色
- 拼写修复
- 不改业务行为的小修正

## 本仓库专项约束

- 后端固定目录：`server/`
- 前端固定目录：`web/`
- 后台页面默认 Drawer 作为 create/edit 交互
- 不要重新引入代码生成器工作流
- 规格、计划、代码、验证结果不一致时，不允许归档
