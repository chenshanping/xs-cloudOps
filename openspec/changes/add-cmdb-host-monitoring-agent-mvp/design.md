## Context

`caelor` 的 CMDB 模块当前能管理主机、分组、SSH 凭据，能通过 WebSocket 打开 SSH 终端、测试 SSH 连通性，但缺乏对主机运行态的**持续观测**。用户希望借鉴 `E:\go_project\deviops` 的监控设计，但**不希望一次性把 Prometheus/Pushgateway 引入当前项目**。

`deviops` 的做法可以概括为：

- 后端**动态编译** Go agent（main.go 字符串模板 + go build）
- Agent 用 `gopsutil` 采集，15s push 到 Pushgateway
- Agent 额外每 60s POST 心跳到后端 `/api/v1/monitor/agent/heartbeat`
- 后端通过 PromQL 查询 Prometheus 获取指标
- 心跳 token 全局共享；Agent 按 hostname/IP 反查主机
- 有独立 `agent` 表记录 pid、last_heartbeat、安装进度

我们不直接搬这套架构，因为：

- **额外 Ops 成本**：Prometheus + Pushgateway 要你自己部署
- **Go 工具链依赖**：后端服务器要装 go 才能编译 agent
- **弱鉴权**：全局共享 token 无法吊销
- **一次性耦合过多**：采集 + 时序存储 + 查询 + 曲线 UI 一把梭，难复盘

## Goals / Non-Goals

**Goals:**

- 让已配置 SSH 凭据的主机能一键部署监控 agent。
- 后端持久化每台主机最新的一份 CPU/内存/磁盘/负载/运行时长快照。
- 心跳超时能自动把主机标为离线。
- 前端列表 + 详情能看到当前资源使用情况与最后心跳。
- Agent 与后端之间使用 **per-host token**，可随时吊销。
- Agent 交付物是**预编译二进制**，安装链路不依赖服务器上装 Go。
- 设计保留延伸到阶段 2（Prometheus + 曲线）的可能，但不在本轮落地。

**Non-Goals:**

- 不把指标写入 Prometheus。
- 不做历史曲线、不做图表 UI。
- 不做告警规则 / 通知通道（邮件/飞书/Webhook）。
- 不做 agent 自升级、灰度、版本管理。
- 不做多租户或跨项目共享的 agent 网关。
- 不做进程 TOP-N、网络流量、TCP 端口监控（留阶段 2 扩展）。
- 不做 agent 侧的 `/metrics` HTTP 端点（只主动上报）。

## Decisions

### 1. Agent 以独立项目交付，预编译二进制 + SSH 推送安装

**决策：**

- 新增独立 Go 模块 `caelor-agent`，依赖 `gopsutil/v3`。
- 后端发布流程同时 cross-build 出 `linux/amd64`、`linux/arm64` 两个二进制，放入 `server/assets/agent/`。
- 后端部署接口流程：
  1. 按主机 `OS/arch` 选对应二进制。
  2. 通过 CMDB `ssh_credential` + `ssh_ip` SCP 到 `/usr/local/bin/caelor-agent`。
  3. 写入 systemd unit `/etc/systemd/system/caelor-agent.service`，env 中带 `REPORT_URL`、`AGENT_TOKEN`。
  4. `systemctl enable --now caelor-agent`。
- 卸载接口执行 `systemctl disable --now` + 文件清理。

**原因：**

- 不需要服务器装 Go。
- 复用 CMDB 已有 SSH 凭据，没有额外的登录机制。
- systemd 是主流 Linux 发行版默认存在，不增加运维前置条件。

**备选方案：**

- **后端动态编译（deviops 方案）**：灵活，但强制在后端部署环境维护 Go 工具链。
- **让 agent 主动下载**：需要再开一条文件分发通道，实现成本更高，且要解决首个 token 如何下发的问题。

### 2. 指标协议：JSON POST，不引入 Prometheus

**决策：**

- Agent → 后端使用两个 REST 接口：
  - `POST /api/v1/cmdb/agents/report`：指标 + 心跳合并一次请求。
  - 请求体：
    ```json
    {
      "token": "...",
      "hostname": "...",
      "pid": 1234,
      "version": "0.1.0",
      "collected_at": "2026-05-03T03:15:00Z",
      "metric": {
        "cpu_percent": 12.3,
        "mem_total_bytes": 17179869184,
        "mem_used_bytes": 5368709120,
        "disk_usage": [{"mount": "/", "total": 0, "used": 0, "percent": 0}],
        "load_1": 0.42,
        "load_5": 0.33,
        "load_15": 0.20,
        "uptime_seconds": 12345
      }
    }
    ```
- 后端直接把快照 `upsert` 进 `cmdb_host_metric`；不保存历史记录。
- 心跳语义：成功收到上报即等价于一次心跳。

**原因：**

- 前期看快照已足够回答“这机器活不活、累不累”。
- 不引入 Prometheus 就不用管 `instance` 标签、scrape 配置、查询端点可用性。
- 一个 endpoint 同时承担心跳 + 指标，agent 逻辑极简。

**备选方案：**

- **直接上 Pushgateway**：曲线能力天然具备，但强依赖外部服务。
- **自研时序表**：自然扩展到曲线，但会立刻把本轮 scope 撑大，违反 AGENTS.md “scope discipline”。

### 3. 鉴权：per-host token，而不是全局共享 token

**决策：**

- 安装 agent 时由后端为该主机随机生成 token（`crypto/rand` 32 字节 hex），写入 `cmdb_host_agent.token`，同时通过环境变量下发到 agent systemd unit。
- `/api/v1/cmdb/agents/report` 为**公开路由**，不走 JWT，但必须做 token 校验，并在进入业务前按 token 查到 `host_id`。
- 吊销 = 在 DB 里把 `cmdb_host_agent.status` 设为 `revoked` 或直接删除该行，agent 下次上报会被拒绝。
- 卸载 agent 同时删除/置吊销该行。

**原因：**

- 对齐 AGENTS.md 后端安全原则：高风险入口采用后端策略而不是共享 secret。
- 便于吊销、审计、未来做 IP 白名单/CIDR 限制。

**备选方案：**

- **全局共享 token（deviops 方案）**：实现最简，但任何一个 agent 泄露即全军覆没。
- **mTLS**：更安全，但让 agent 部署和 CA 管理变复杂，本轮不值得。

### 4. 离线判定：定时 cron，不修改 status 现有业务语义

**决策：**

- 新增后台定时任务，频率 1 分钟，判定：`cmdb_host_agent.last_heartbeat_at < now - 2min` 的主机。
- 命中后：
  - `cmdb_host_agent.status` = `offline`
  - 如果 `cmdb_host.status` 仍为 1（在线），改为 3（离线）。
- 成功心跳时：
  - `cmdb_host_agent.status` = `running`
  - `cmdb_host.status` 如果是 2 或 3，则改为 1。

**原因：**

- 和现有“SSH 测试通过 = status=1”不冲突，只是多了一条更快的自动更新路径。
- 心跳频率 15s + 超时 120s 是业界常见的“容忍一次网络抖动”的平衡点。

**备选方案：**

- **让 cmdb_host 的 status 完全由 agent 接管**：会把无 agent 主机打成离线，与现有语义冲突。
- **不做离线判定，只展示 last_heartbeat**：体验差，用户需要自己算。

### 5. 存储只保留最新快照，不做时序

**决策：**

- `cmdb_host_metric` 以 `host_id` 为主键唯一约束，每次上报 `ON DUPLICATE KEY UPDATE`。
- 磁盘多挂载点以 JSON 存在一列里，避免每挂载点一行的复杂度。

**原因：**

- 快照够用，表体积可控。
- 阶段 2 再接入 Prometheus 或另起历史表，不冲突。

**备选方案：**

- **每次上报插一行**：表会无限膨胀，需要额外 TTL/归档策略，超出 MVP 必要性。

## Risks / Trade-offs

- **仅有快照，失去历史数据观察能力**：阶段 1 接受，前端明确文案“当前快照”，阶段 2 再补。
- **Agent 部署失败可能留脏数据**：必须在部署接口的异常路径中清理 `cmdb_host_agent` 行和远端文件，避免“看起来装过但其实没跑”。
- **per-host token 下发链路**：token 一次性通过 SSH env 下发；agent 将其持久化到 systemd unit 文件（权限 0600）。后端**永远不展示 token 明文**。
- **定时 cron 并发**：1 分钟一次、一次扫全量，数据量小的阶段足够；后续若主机数过千再考虑分片。
- **systemd 依赖**：不支持 systemd 的极小/老旧发行版暂不覆盖，提示用户手工部署或阶段 2 提供 init.d 兼容。
- **指标上报失败无重试队列**：阶段 1 直接丢弃本次上报，15s 后下一次；不值得为 MVP 做本地 WAL。
- **公开上报路由滥用风险**：必须强制 token 校验 + 请求大小限制 + 速率限制（复用现有中间件，若无则最简 IP-level limiter 兜底）。

## Migration Plan

1. 编写 SQL 升级脚本新增两张表，与当前 `add_cmdb_host_management.sql` 并列，不改原脚本。
2. 后端 AutoMigrate 补齐 `CmdbHostAgent`, `CmdbHostMetric`。
3. 新增 `caelor-agent` 子项目，能在本地跑通 `go run ./cmd/caelor-agent`。
4. 后端实现 `report` 接口 + token 校验 + 快照 upsert。
5. 后端实现部署/卸载接口，走 `ssh_credential`。
6. 后端实现离线 cron。
7. 后端注册菜单/API 元数据，初始化默认授权。
8. 前端接入列表列、详情面板、操作按钮。
9. 手动集成验证：部署到 1 台测试机 → 看快照刷新 → kill agent → 2 分钟内变离线 → 卸载干净。
10. 如需回滚：
    - 停 offline cron
    - 下线 `/api/v1/cmdb/agents/report` 路由
    - 前端关闭监控相关列与按钮
    - 两张新表数据**不回滚**，保留供后续阶段使用

## Open Questions

- 本轮 agent 只覆盖 Linux 还是同时发 `darwin/amd64` 给本地开发机测试？当前倾向只发 `linux/{amd64,arm64}`，后续按需补。
- `caelor-agent` 源码放本仓库 `apps/caelor-agent/` 还是新仓库？当前倾向本仓库子目录，便于一体化发布。
