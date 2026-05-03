## 1. OpenSpec alignment

- [ ] 1.1 Add proposal, design, tasks, and capability spec for `cmdb-host-monitoring-mvp`.
- [ ] 1.2 Keep phase boundary explicit: phase 1 ships agent + JSON report + snapshot only; Prometheus / time-series / alerts / charts remain out of scope.

## 2. Database & models

- [ ] 2.1 Add `server/sql/add_cmdb_host_monitoring.sql`: idempotent creation of `cmdb_host_agent` and `cmdb_host_metric`, no foreign keys, indexes on `host_id` and `last_heartbeat_at`; verified against MySQL 8 rerun.
- [ ] 2.2 Add `server/model/cmdb_host_agent.go` with fields `HostID`, `Token`, `Version`, `Status`, `Pid`, `LastHeartbeatAt`, `InstalledAt`, `Platform`.
- [ ] 2.3 Add `server/model/cmdb_host_metric.go` with snapshot fields for CPU %, memory, disk JSON, load, uptime, `CollectedAt` and unique key on `HostID`.
- [ ] 2.4 Extend `server/initialize/db_tables.go` AutoMigrate list to include the new models, without overwriting customized metadata on restart.

## 3. Agent codebase

- [ ] 3.1 Introduce `apps/caelor-agent/` as an independent Go module using `github.com/shirou/gopsutil/v3`.
- [ ] 3.2 Implement a single `report` loop: every 15s collect CPU %, memory total/used, disk usage per mount, load averages, uptime, and POST JSON to configured `REPORT_URL` with bearer-like token.
- [ ] 3.3 Read `REPORT_URL` and `AGENT_TOKEN` from environment variables only; fail fast if missing. Do not persist token to disk beyond systemd unit file (owned by root, mode 0600).
- [ ] 3.4 Handle transient HTTP errors with exponential backoff capped at 60s; never crash the process on upload failure.
- [ ] 3.5 Provide a Makefile / build script producing `linux/amd64` and `linux/arm64` binaries into `server/assets/agent/` for deployment use.

## 4. Backend service & API

- [ ] 4.1 Add `server/service/cmdb/agent.go` with install / uninstall / heartbeat-miss workflows; install must generate a cryptographically-random per-host token and upsert `cmdb_host_agent`.
- [ ] 4.2 Add `server/service/cmdb/metric.go` responsible for snapshot upsert and host status transition (in-memory; no history retention).
- [ ] 4.3 Add `server/api/v1/cmdb_monitor.go`:
  - `POST /api/v1/cmdb/agents/report` as public route, validating token, size-limiting body, and rejecting unknown tokens.
  - `POST /api/v1/cmdb/hosts/:id/monitor/deploy` as admin-authenticated install endpoint.
  - `DELETE /api/v1/cmdb/hosts/:id/monitor` for uninstall.
  - `GET /api/v1/cmdb/hosts/:id/monitor` returning snapshot + agent status.
- [ ] 4.4 Register the new routes in `server/router/modules/cmdb.go`; public route placement MUST NOT attach JWT middleware.
- [ ] 4.5 Add startup metadata and default grants for `cmdb:host:monitor:deploy` and `cmdb:host:monitor:view` in `server/initialize/db_tables.go` without overwriting customized built-in metadata.

## 5. Deployment plumbing

- [ ] 5.1 Implement SSH-based install procedure: SCP the platform-matched binary to `/usr/local/bin/caelor-agent`, render systemd unit with env vars, `systemctl enable --now`.
- [ ] 5.2 Implement uninstall procedure: `systemctl disable --now caelor-agent`, remove unit + binary, mark `cmdb_host_agent.status = revoked`.
- [ ] 5.3 Guard deploy/uninstall paths to reuse existing `ssh_credential` dial logic without duplicating SSH client setup.
- [ ] 5.4 Error paths must clean up partial state: if SCP or systemctl fails, remove any remote artifacts and delete the `cmdb_host_agent` row created in that attempt.

## 6. Offline checker

- [ ] 6.1 Implement a 1-minute cron task that flips `cmdb_host_agent.status = offline` and sets `cmdb_host.status = 3` for any agent whose `last_heartbeat_at < now - 2min`.
- [ ] 6.2 Ensure that a successful report restores `cmdb_host_agent.status = running` and bumps `cmdb_host.status` from 2 or 3 back to 1.
- [ ] 6.3 Do not touch `cmdb_host.status` for hosts that have never been deployed an agent.

## 7. Frontend integration

- [ ] 7.1 Extend `web/src/types/cmdb.ts` with `CmdbHostAgent`, `CmdbHostMetric`, and an aggregate type returned by the monitor endpoint.
- [ ] 7.2 Add API wrappers in `web/src/api/cmdb.ts`: `deployMonitor`, `uninstallMonitor`, `getHostMonitor`.
- [ ] 7.3 Add a monitor cell (CPU %, memory %, last heartbeat) to `web/src/views/admin/cmdb/host/index.vue` without breaking existing columns.
- [ ] 7.4 Extend `HostDetailDrawer.vue` with a "系统监控" panel showing current snapshot + agent status; gracefully handle "未部署 agent" state.
- [ ] 7.5 Add `部署监控 / 卸载监控` buttons gated by `cmdb:host:monitor:deploy`; successful action refreshes host list.
- [ ] 7.6 Respect existing dark mode variables and Ant Design Vue patterns already used by host pages.

## 8. Security & safety

- [ ] 8.1 `POST /api/v1/cmdb/agents/report` MUST reject missing/invalid token before reading the full body.
- [ ] 8.2 Body size MUST be capped (e.g. 64KB) to avoid malicious clients blowing up the server.
- [ ] 8.3 Agent tokens MUST NEVER be returned by any admin API; only their existence/status is surfaced.
- [ ] 8.4 Logs for the public report route MUST redact tokens.

## 9. Verification

- [ ] 9.1 Run `cd server && go test ./...`.
- [ ] 9.2 Run `cd web && npm run build`.
- [ ] 9.3 Run `cd web && npm run typecheck`, or explicitly report the known existing vue-tsc environment issue if it still blocks verification.
- [ ] 9.4 Manual verification on one Linux host: deploy agent → observe snapshot refresh within 30s → kill agent → observe offline within 2 minutes → uninstall leaves no leftover files, no agent row, no systemd unit.
- [ ] 9.5 Confirm SQL upgrade script runs twice in a row on an existing DB without error.
