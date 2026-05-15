# Codex + OpenSpec + 项目 Skills 使用说明

本仓库已经切到项目级 Codex 工作方式。目标不是“让 AI 直接写代码”，而是让需求探索、规格留痕、实现执行、验证归档都有固定入口。

## 先看什么

开始任何非 trivial 任务前，先看：

1. `AGENTS.md`
2. `openspec/config.yaml`
3. `docs/superpowers/README.md`

## 这套组合怎么分工

- 项目 Skills：负责探索、计划、执行、验证的方法入口
- OpenSpec：负责 proposal / design / spec / tasks 的落盘与追踪
- Codex：负责读仓库、改代码、跑命令、验证结果

一句话：

先用项目 Skills 想清楚，再用 OpenSpec 锁住边界，最后回到 Codex 执行。

## 本仓库推荐流程

### 轻量任务

直接做：

- 文案调整
- 小型文档补充
- 不改变行为的微小配置修正

### 标准任务

按下面顺序执行（可以直接在对话里显式写 `$skill-name`，也可以依赖 Codex 自动触发）：

1. `$go-base-brainstorming`
2. `$go-base-openspec-propose` 或 `$go-base-openspec-explore`
3. `$go-base-writing-plans`
4. `$go-base-openspec-apply-change`
5. `$go-base-openspec-archive-change`

## 本仓库已放入的项目内 Skills

位置：`.codex/skills/`

已接入（显式调用时使用 `$skill-name`）：

- **OpenSpec**：`$go-base-openspec-explore`、`$go-base-openspec-propose`、`$go-base-openspec-apply-change`、`$go-base-openspec-archive-change`
- **计划与执行**：`$go-base-brainstorming`、`$go-base-writing-plans`、`$go-base-executing-plans`、`$go-base-systematic-debugging`、`$go-base-finishing-branch`
- **项目专用**：`$go-base-backend-crud-frontend`、`$go-base-sql-upgrade-guardrails`、`$go-base-file-reference-guardrails`、`$go-base-generate-login-config`

## 建议提示词

### 1. 只做需求探索

```text
请按本仓库 AGENTS.md 执行。
先不要写代码，也不要直接开 OpenSpec。
$go-base-brainstorming
基于当前 caelor 项目做需求探索、方案比较和推荐设计。
```

### 2. 设计确认后生成规范

```text
设计已确认。
$go-base-openspec-propose
为这个需求创建 change，补齐 proposal、design、delta spec 和 tasks，完成后暂停，不要开始编码。
```

### 3. 开始执行

```text
OpenSpec 已确认。
$go-base-writing-plans 生成实现计划，
然后 $go-base-openspec-apply-change 执行并验证。
```

### 4. 归档

```text
请确认代码、测试、tasks、spec 一致。
验证通过后 $go-base-openspec-archive-change 归档这个 change。
```

## 本仓库约定

- 后端目录固定为 `server/`
- 前端目录固定为 `web/`
- 已移除代码生成器思路，不要重新走生成器式开发
- 后台 create/edit 默认使用 `Drawer`
- 非 trivial 的 Drawer/弹窗内容拆到页面本地 `components/`

## 常用验证命令

### Backend

```bash
cd server
go test ./...
```

### Frontend

前端默认不跑 `npm run build` / `npm run typecheck` 作为代理验证。

- 先读回改动文件确认结构
- 再做定点引用搜索
- 最后交给用户在 dev server 做点击验证
- 只有用户明确要求时，才跑 `build` 或 `typecheck`
