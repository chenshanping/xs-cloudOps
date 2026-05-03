## ADDED Requirements

### Requirement: Administrators can deploy a monitoring agent onto a host that already has a valid SSH credential

The system MUST allow an administrator with `cmdb:host:monitor:deploy` permission to install the platform-provided monitoring agent onto a CMDB host by reusing that host's configured SSH credential, and MUST persist a per-host agent token, status, and installation metadata before any agent traffic is accepted.

#### Scenario: Successful agent deployment

- **GIVEN** a CMDB host has a non-zero `ssh_credential_id`, a reachable `ssh_ip`, and no existing running agent
- **WHEN** an administrator invokes the deploy-monitor action
- **THEN** the system MUST:
  - Generate a cryptographically-random token for that host
  - Push the platform-matched agent binary to the host over SSH
  - Install and start the agent with `REPORT_URL` and `AGENT_TOKEN` environment variables
  - Persist an agent record with `status = running` and `installed_at = now`
  - Never return the token to the frontend

#### Scenario: Deployment fails with partial remote state

- **GIVEN** an administrator invokes deploy-monitor
- **WHEN** any deployment step (copy, unit install, start) fails
- **THEN** the system MUST:
  - Attempt to clean up any remote artifacts it created in that attempt
  - Remove the agent row created for that attempt, if any
  - Return a failure response explaining the failing step

### Requirement: Administrators can uninstall the monitoring agent for a host

The system MUST allow an administrator with `cmdb:host:monitor:deploy` permission to uninstall a host's monitoring agent, and MUST revoke the previously-issued token so that any later report attempt using that token is rejected.

#### Scenario: Uninstall removes agent and revokes token

- **GIVEN** a host has an installed agent
- **WHEN** an administrator invokes the uninstall-monitor action
- **THEN** the system MUST stop and remove the remote agent, MUST delete or mark its agent record as `revoked`, and MUST reject any subsequent report request using the previous token

### Requirement: The public agent report endpoint must authenticate every request by per-host token

The system MUST expose a single public endpoint that accepts combined heartbeat + metric reports from the agent, and MUST reject any report whose token does not match an active `cmdb_host_agent` row before doing any business processing.

#### Scenario: Valid report is accepted and snapshot is upserted

- **GIVEN** an active agent sends a well-formed report with its assigned token
- **WHEN** the server receives the report
- **THEN** the server MUST upsert the latest metric snapshot for that host and update `last_heartbeat_at` to the server-side time of receipt

#### Scenario: Report with unknown or revoked token is rejected

- **GIVEN** a request uses a token that does not match any active agent record
- **WHEN** the server receives the report
- **THEN** the server MUST reject the request, MUST NOT persist any metric or heartbeat data, and MUST NOT leak whether the token ever existed

#### Scenario: Oversized or malformed body is rejected before persistence

- **GIVEN** a report request whose body exceeds the configured size limit or fails JSON validation
- **WHEN** the server processes the request
- **THEN** the server MUST reject the request without persisting any partial data

### Requirement: The system must maintain only the latest snapshot per host, not time-series history

The system MUST keep at most one metric snapshot row per host, and MUST overwrite it on each accepted report. The system MUST NOT accumulate historical metric rows in this capability.

#### Scenario: Repeated reports update the same row

- **GIVEN** a host is already represented in the metric snapshot table
- **WHEN** a new valid report for that host is accepted
- **THEN** the system MUST update the existing row in place rather than inserting a new row

### Requirement: Host online/offline status must track agent heartbeat with a bounded staleness window

The system MUST automatically mark a host as offline when its agent has not successfully reported for longer than the configured heartbeat timeout, and MUST return that host to online state when a valid report is accepted again. Hosts without any deployed agent MUST NOT be affected by this behavior.

#### Scenario: Heartbeat timeout flips host to offline

- **GIVEN** a host has an active agent record and its `last_heartbeat_at` is older than the configured timeout
- **WHEN** the offline checker runs
- **THEN** the system MUST set the agent record to `offline` and MUST set `cmdb_host.status` to offline for that host

#### Scenario: Successful report restores online state

- **GIVEN** a host whose agent record is `offline` or whose `cmdb_host.status` is 2 or 3
- **WHEN** a valid report from that host's agent is accepted
- **THEN** the system MUST set the agent record back to `running` and MUST set `cmdb_host.status` to online

#### Scenario: Hosts without an agent are not touched

- **GIVEN** a host that has never had an agent installed
- **WHEN** the offline checker runs
- **THEN** the system MUST NOT change that host's `status` based on agent heartbeat logic

### Requirement: Administrators can read the latest monitoring snapshot for a host

The system MUST expose an admin-authenticated endpoint that returns the current agent status plus the latest metric snapshot for a host, and MUST gracefully represent hosts that have never been deployed an agent.

#### Scenario: Read snapshot for a monitored host

- **GIVEN** a host has an active agent and a stored snapshot
- **WHEN** an administrator with `cmdb:host:monitor:view` permission reads the host monitor endpoint
- **THEN** the system MUST return the agent status, `last_heartbeat_at`, and the latest metric snapshot

#### Scenario: Read snapshot for a host without agent

- **GIVEN** a host has no agent record
- **WHEN** an administrator reads the host monitor endpoint
- **THEN** the system MUST return a response that clearly indicates "未部署 agent" without fabricating metric values and MUST NOT error
