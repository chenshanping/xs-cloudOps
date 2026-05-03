# Windsurf + OpenSpec + Superpowers 使用说明

本仓库已经接入项目级 AI 工作流，目标不是“让 AI 直接写代码”，而是让需求探索、规格留痕、实现执行、验证归档都有固定入口。

## 先看什么

开始任何非 trivial 任务前，先看：

1. `AGENTS.md`
2. `openspec/config.yaml`
3. `docs/superpowers/README.md`

## 这套组合怎么分工

- Superpowers：负责探索、计划、执行、验证的方法
- OpenSpec：负责 proposal / design / spec / tasks 的落盘与追踪
- Windsurf（Cascade）：负责读仓库、改代码、跑命令、验证结果

一句话：

先用 Superpowers 想清楚，再用 OpenSpec 锁住边界，最后回到 Superpowers + Windsurf 执行。

## 本仓库推荐流程

### 轻量任务

直接做：

- 文案调整
- 小型文档补充
- 不改变行为的微小配置修正

### 标准任务

按下面顺序执行（在 Windsurf Cascade 里输入斜杠命令即可）：

1. `/brainstorming`
2. `/openspec-propose` 或 `/openspec-explore`
3. `/writing-plans`
4. `/openspec-apply-change`
5. `/openspec-archive-change`

## 本仓库已放入的本地工作流

位置：`.windsurf/workflows/`

已接入（对应 Cascade 斜杠命令）：

- **OpenSpec**：`/openspec-explore`、`/openspec-propose`、`/openspec-apply-change`、`/openspec-archive-change`
- **Superpowers**：`/brainstorming`、`/writing-plans`、`/executing-plans`、`/systematic-debugging`、`/finishing-branch`
- **项目专用**：`/backend-crud-frontend`

## 建议提示词

### 1. 只做需求探索

```text
请按本仓库 AGENTS.md 执行。
先不要写代码，也不要直接开 OpenSpec。
/brainstorming
基于当前 caelor 项目做需求探索、方案比较和推荐设计。
```

### 2. 设计确认后生成规范

```text
设计已确认。
/openspec-propose
为这个需求创建 change，补齐 proposal、design、delta spec 和 tasks，完成后暂停，不要开始编码。
```

### 3. 开始执行

```text
OpenSpec 已确认。
/writing-plans 生成实现计划，
然后 /openspec-apply-change 执行并验证。
```

### 4. 归档

```text
请确认代码、测试、tasks、spec 一致。
验证通过后 /openspec-archive-change 归档这个 change。
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

```bash
cd web
npm run build
npm run typecheck
```

说明：

- 当前 `web` 的 `typecheck` 在某些环境下可能仍受已知 `vue-tsc` 工具链问题影响。
- 如果出现该已知环境错误，要明确报告，不允许假装“验证通过”。
