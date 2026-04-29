## Why

当前 `go-base` 已经形成了稳定的双轨权限模型：

- 菜单权限负责前端菜单、页面和按钮可见性
- API 权限负责后端接口访问控制，且仍以 Casbin 为主执行层

但最近的真实使用暴露出几个问题：

- 管理员容易误以为“勾了按钮权限就应该同时拥有接口权限”
- 也容易误以为“勾了 API 权限就应该自动看到前端按钮”
- 角色重新打开权限抽屉时，如果菜单或 API 回显不完整，会被误判成“没有写入数据库”
- 按钮型菜单单独授权时，如果父级菜单链没有被保留，会导致 `userinfo.menus` 为空或页面入口缺失
- 修改 API 权限后如果 Casbin 运行时策略没有及时同步，会出现“数据库已写入但接口仍 403”

这次变更不改变双轨模型，而是把它正式收敛成可理解、可验证、可排错的仓库行为，降低后续继续误用和误判的成本。

## What Changes

- 保持当前双轨模型不变：
  - 菜单权限继续控制前端菜单与按钮显示
  - API 权限继续控制后端接口访问
  - Casbin 继续作为后端接口权限主链路
- 为角色权限分配补齐并固定以下行为：
  - 角色菜单与角色 API 仍在独立授权区域中管理
  - 权限界面明确提示两者职责边界，不做“自动互相授权”
  - 角色详情重新加载时必须准确回显已保存的菜单与 API 绑定
  - 按钮型菜单单独授权时，系统必须保留所需父级菜单链，保证用户菜单树可用
  - 角色 API 权限保存后必须立即同步 Casbin 策略，使接口权限立刻生效
- 将这套双轨规则写入 OpenSpec，作为后续权限相关改动的行为基线

### Out of Scope

- 不改成统一权限码模型
- 不移除 Casbin，也不降低 Casbin 在后端接口鉴权中的主地位
- 不做“菜单自动推导 API 权限”或“API 自动生成菜单权限”的隐式联动
- 不扩展到业务模块级权限体系重构，只覆盖当前后台系统角色权限分配链路
- 不改变已归档的部门管理、部门数据范围权限规则

## Capabilities

### New Capabilities

- `admin-role-permission-management`: 管理员可以在保持菜单权限与 API 权限分离的前提下，稳定分配、回显和验证角色权限。

### Modified Capabilities

- None. Repository currently has no main spec for this admin role permission behavior.

## Impact

- Backend:
  - `server/service/role`
  - `server/service/menu`
  - `server/middleware/casbin.go`
  - `server/service/auth/auth_flow.go`
- Frontend:
  - `web/src/views/admin/system/role/components/RolePermissionDrawer.vue`
  - related role page hooks / API wrappers as needed
- Persistence:
  - existing `sys_role_menu`, `sys_role_api`, and `casbin_rule` semantics remain in place
- Rollback:
  - remove the dual-track explanatory UI changes
  - revert any new role permission read/write hardening
  - keep the original separate menu/API model and existing data tables unchanged
- Specs:
  - no current archived spec needs to be replaced or archived for this change
  - this change adds a new capability spec describing the dual-track permission baseline
