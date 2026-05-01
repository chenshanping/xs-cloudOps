## Context

当前仓库的角色权限模型已经稳定为双轨结构：

- `Roles.Menus` 决定前端菜单树和按钮权限码
- `Roles.Apis` 决定后端 API 访问能力
- `casbin_rule` 负责运行时接口放行
- 角色数据权限由 `SysRole.DataScope` 与 `SysRoleDataScope` 单独治理

这套结构对后台项目是可行的，但企业化问题集中在“治理”而不是“能不能跑”：

- 保存动作不是原子的
- 菜单和 API 双维护成本高
- 白名单和日志规则缺少长期治理约束
- 运行时授权来源不够直观，缺少统一回归基线

本设计目标是保留双轨模型的职责边界，同时把角色权限治理从“能用”收敛成“可长期维护、可审计、可回滚”。

约束：

- 继续使用 Gin + GORM + Casbin + JWT + Redis
- 不替换 `sys_role_api`
- 不改成统一权限码直校
- 前端继续沿用当前 Drawer 与 `components/` 拆分模式
- 优先复用现有角色权限页、菜单管理页和权限测试结构

## Goals / Non-Goals

**Goals:**

- 用一个保存入口原子提交角色菜单权限、角色 API 权限、角色数据权限
- 为菜单/按钮建立显式 API 关联关系，降低长期双维护
- 明确角色 API 的最终生效集合，并保证 Casbin 立即同步
- 收紧接口白名单，使绕过授权的规则最小化且可审计
- 给操作日志加统一脱敏和更稳健的异步写入策略
- 为权限治理行为补齐可重复验证的回归测试

**Non-Goals:**

- 不把前端权限和后端权限压缩成同一个直接校验源
- 不做多租户权限域模型
- 不做权限审批流程、审批工单或授权申请流
- 不在本次 change 中重构全部业务模块权限

## Decisions

### 1. 角色权限改为单接口原子保存

决策：

- 新增统一角色权限保存接口，单次提交中包含：
  - `menu_ids`
  - `direct_api_ids`
  - `feature_data_scopes`
- 服务层使用单事务完成：
  - 角色菜单更新
  - 角色直接 API 更新
  - 角色数据权限更新
  - 角色最终 Casbin 策略重建
  - 相关缓存失效

原因：

- 当前前端并发调用多个保存接口，容易出现部分成功
- 企业权限变更必须尽量避免脏状态

备选方案：

- 保持多接口保存，仅增强前端错误提示：不能解决数据库中间态问题

### 2. 新增菜单到 API 的显式关联层

决策：

- 为菜单和按钮维护显式 API 关联集合
- 这个关联表达“授权此菜单时，角色默认应拥有这些 API 访问能力”
- 关联关系由菜单管理能力维护，而不是在角色页临时推导

原因：

- 方案 B 的核心是降低双维护，而不是继续仅靠解释来约束管理员
- 角色页只负责选角色拥有的菜单和直接 API，不适合承担菜单/API 关系建模职责

备选方案：

- 纯前端按路径或关键字猜测菜单对应 API：不稳定，不可审计
- 直接取消 `sys_role_api`：会失去“无菜单入口但需要单独开放接口”的能力

### 3. 角色最终 API 权限由“直接授权 + 菜单继承”组成

决策：

- 角色最终接口权限集合 =
  - `sys_role_api` 中的直接 API
  - 角色已选菜单通过菜单 API 关联继承得到的 API
- Casbin 同步时使用最终并集，而不是只使用直接 API
- 角色详情回显时区分：
  - 直接 API
  - 由菜单继承而来的 API

原因：

- 既能降低双维护，又保留高级精细授权能力
- 不破坏现有双轨模型：“菜单负责前端显示，API 负责后端接口”仍成立
- 继承只是治理层补齐，不代表菜单权限直接等同于接口权限

备选方案：

- 菜单继承 API 后不保留直接 API：会丢掉特例授权能力
- 把继承 API 直接写回 `sys_role_api`：会混淆“直接授权”和“继承授权”的来源

### 4. 白名单只保留显式路由，不保留后缀规则

决策：

- 删除基于 `/my` 之类路径后缀的通配放行
- 仅保留显式白名单路由
- 如后续确需扩展，要求在白名单定义中显式声明具体路由模板

原因：

- 企业权限治理中，后缀型白名单难以审计且容易被误扩展
- 放行规则应可枚举、可评审、可测试

备选方案：

- 保留后缀白名单并补文档：风险仍然存在

### 5. 操作日志引入脱敏和受控异步写入

决策：

- 统一脱敏请求体和响应体中的敏感字段，如：
  - password
  - token
  - authorization
  - secret
  - access_key / secret_key
  - email code / captcha code
- 将当前“每请求启动 goroutine 直接写库”收敛为受控缓冲策略
- 对超长响应、二进制内容和敏感配置写入继续保留截断与替换策略

原因：

- 企业项目里日志默认可被运维、审计、排障工具读取，敏感信息不能裸写
- 无控制 goroutine 直写数据库不利于高并发稳定性

备选方案：

- 仅增加字段长度截断：不能解决泄露问题

## Persistence Design

新增一层菜单到 API 的持久化关系。

建议：

- 新增 `sys_menu_api` 关系表
- 使用多对多结构表达一个菜单或按钮关联多个 API

保留现有表：

- `sys_role_menu`
- `sys_role_api`
- `casbin_rule`

语义划分：

- `sys_role_menu`: 角色显式拥有的菜单/按钮
- `sys_role_api`: 角色显式拥有的直接 API
- `sys_menu_api`: 菜单默认继承的 API
- `casbin_rule`: 运行时实际放行策略，由角色最终 API 集合同步得出

## Backend Flow

### Role permission save flow

```text
Admin saves role permissions
        |
        v
PUT /roles/:id/permissions
        |
        v
Transaction:
  1. Replace role menus
  2. Replace direct role APIs
  3. Replace feature data scopes
  4. Load menu-inherited APIs
  5. Build final API union
  6. Rebuild Casbin rules for role
  7. Clear user cache by role
        |
        v
Committed or fully rolled back
```

### Menu API association flow

```text
Admin edits menu
        |
        v
Assign related APIs to menu/button
        |
        v
Save sys_menu_api mappings
        |
        v
Future role menu selections inherit these APIs
```

## Frontend Design

### Role permission drawer

- 保留当前 Drawer
- 保存入口改为单次提交
- 角色页继续区分菜单权限区域和 API 权限区域
- API 区域需要能区分：
  - 直接授予
  - 菜单继承
- 对继承 API 显示来源菜单，避免管理员误删或误判

### Menu management

- 在现有菜单管理抽屉中补一个“关联 API”区域
- 只对类型为菜单/按钮的项开放
- 目录类型默认不管理关联 API
- 交互仍保持紧凑，不增加独立新页面

## Security and Audit

- 白名单规则必须最小化
- 操作日志必须对敏感字段脱敏
- 菜单 API 关联变更、角色权限总保存变更都应进入操作日志
- 角色权限保存必须避免部分成功带来的授权异常窗口

## Validation

后端验证重点：

- 单接口保存成功时，菜单、直接 API、数据权限、Casbin 规则全部同步
- 单接口保存任一子步骤失败时，数据库与 Casbin 规则不应出现部分提交
- 菜单继承 API 与直接 API 合并后，接口权限立即生效
- 删除菜单继承来源后，未被直接授予的 API 权限立即失效

前端验证重点：

- 角色权限抽屉重新打开时，菜单、直接 API、继承 API 状态回显正确
- 菜单管理抽屉可维护关联 API
- 保存失败时，错误反馈明确且不会误导为全部成功

## Risks / Trade-offs

- [新增 `sys_menu_api` 关系后，菜单管理复杂度上升] → 接受，用显式关联换长期治理稳定性
- [角色 API 回显需要区分直接/继承来源] → 前端复杂度上升，但这是管理员可诊断性的必要代价
- [统一保存接口会触及现有角色权限抽屉交互] → 范围可控，因为仍复用现有页面
- [日志脱敏可能影响个别问题排查信息量] → 优先保护敏感数据，再通过摘要和元信息补足排障能力

## Migration Plan

1. 增加菜单 API 关联持久化结构和读取能力
2. 增加菜单 API 管理接口和菜单管理页入口
3. 增加角色权限统一保存接口
4. 调整角色权限前端保存流为单请求提交
5. 调整 Casbin 策略同步来源为直接 API 与继承 API 并集
6. 收紧白名单
7. 调整日志脱敏与异步写入策略
8. 补齐回归测试并完成验证

## Open Questions

- None for this change scope. The direction is fixed: keep Casbin and the dual-track model, but add menu API association and governance hardening on top.
