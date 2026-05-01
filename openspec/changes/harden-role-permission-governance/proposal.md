## Why

当前仓库的角色权限体系已经具备后台企业项目的基础能力：

- 菜单权限控制前端菜单、页面和按钮可见性
- API 权限控制后端接口访问
- Casbin 继续作为后端接口执行层
- 数据权限继续控制资源可见范围

但从治理角度看，仍存在几类明显问题：

- 角色权限保存分成菜单、API、数据权限三个独立写入动作，存在部分成功、部分失败的中间状态
- 菜单权限与 API 权限长期双维护，管理员容易漏配或误判，造成前端显示和后端接口访问漂移
- Casbin 白名单包含基于路径后缀的宽松放行规则，不利于长期权限收敛和审计
- 操作日志默认记录请求体和响应体，缺少统一脱敏与更稳健的异步写入治理
- 当前测试已经覆盖菜单链保留、Casbin 立即生效等行为，但缺少“原子保存”和“菜单关联 API”一致性回归基线

本次变更不是改造为统一权限码架构，也不是替换 Casbin，而是在保留当前双轨模型的前提下，把角色权限治理补到更适合企业项目长期维护的状态。

## What Changes

- 新增一个角色权限治理加固 change，用于补齐当前双轨权限模型的企业化治理能力
- 将角色菜单权限、角色 API 权限、角色数据权限的保存收敛为单个后端原子提交接口
- 为菜单建立显式 API 关联关系，用于表达“这个菜单/按钮依赖哪些后端接口”
- 保留 `sys_role_api` 作为角色直接 API 授权来源，同时允许角色通过已选菜单继承菜单关联 API
- 保留 Casbin 作为后端执行层，但把角色 API 生效来源固定为：
  - 角色直接授予的 API
  - 角色菜单继承得到的关联 API
- 收紧 Casbin 白名单规则，移除仅靠路径后缀的宽泛豁免，改为可审计的显式放行清单
- 为操作日志增加统一脱敏规则，并将数据库写入改为更稳健的异步缓冲/批量治理模式
- 为角色权限治理补齐回归测试，覆盖原子保存、菜单继承 API、生效同步、失败回滚等行为

### Out of Scope

- 不替换 Casbin，不引入 Cerbos、Permify、Oso 等新权限引擎
- 不改造成“前后端完全同一套权限码直接校验”的统一权限模型
- 不取消 `sys_role_api`，也不把 API 授权完全折叠进菜单树
- 不扩展到业务模块级全量权限重构，只覆盖后台角色权限治理链路
- 不改变既有部门数据权限模型的语义，只调整其保存入口的原子性
- 不新增独立权限平台、权限审批流或多租户域权限模型

## Capabilities

### New Capabilities

- `admin-role-permission-governance`: 管理员可以以可回滚、可审计、可继承的方式治理角色菜单权限、角色 API 权限和角色数据权限。
- `admin-menu-api-association`: 管理员可以为菜单或按钮维护关联 API，并让角色在选中菜单时继承这些接口权限。
- `admin-security-governance`: 管理员可以在不放宽授权边界的前提下使用显式白名单和脱敏审计日志治理后台安全例外。

### Modified Capabilities

- None. This change is additive and builds on the current dual-track role permission baseline instead of editing an already-merged capability.

## Impact

- Backend:
  - `server/api/v1/role.go`
  - `server/service/role`
  - `server/service/menu`
  - `server/middleware/casbin.go`
  - `server/middleware/operation_log.go`
  - `server/model`
  - `server/router/modules/menu.go`
  - `server/router/modules/role.go`
- Frontend:
  - `web/src/views/admin/system/role/components/RolePermissionDrawer.vue`
  - `web/src/views/admin/system/role/components/useRolePermissionDrawer.ts`
  - `web/src/views/admin/system/menu/index.vue`
  - related local API wrappers and types as needed
- Persistence:
  - existing `sys_role_menu`, `sys_role_api`, `casbin_rule` remain
  - add a new menu-to-API association persistence structure for explicit menu API mappings
- Rollback:
  - revert single-submit permission save flow
  - remove menu API association management and inheritance logic
  - restore previous explicit-only role API authorization behavior
  - restore previous log handling and exact whitelist behavior if needed
- Specs:
  - current `refine-dual-track-role-permissions` remains valid as the dual-track baseline
  - this change extends that baseline with governance and inheritance behavior instead of replacing it
