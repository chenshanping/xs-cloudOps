## Why

当前角色权限抽屉已经稳定支持菜单权限与 API 权限分离保存，但管理员在实际分配权限时仍需要在两个主 Tab 间来回切换。对于“用户管理”“角色管理”这类页面，菜单、按钮和相关 API 本来属于同一业务上下文，拆成两个主视图会明显增加选择成本和漏配概率。

这次变更要解决的是权限选择效率，而不是权限语义本身。双轨模型、Casbin 主链路、现有保存接口和库表语义都保持不变，只把前端角色权限抽屉重组为更适合批量配置的同屏聚合视图。

## What Changes

- 新增一个独立的角色权限选择体验优化能力，不并入当前双轨权限稳定化 change
- 将角色权限抽屉从“菜单 Tab / API Tab”切换模式改为“按页面聚合的同屏配置”模式
- 保留左侧一级菜单导航，右侧按页面展示菜单/按钮权限区与相关 API 权限区
- 使用前端聚合规则把现有 API 归到页面上下文里，不新增后端映射字段或关系表
- 保持单一保存按钮，但继续通过现有 `assignMenus` 与 `assignApis` 分开提交
- 保存失败时明确区分菜单失败、API 失败或部分失败，不做自动回滚或自动联动

### Out of Scope

- 不改变菜单权限与 API 权限分离的技术决策
- 不新增 Casbin 之外的后端权限引擎
- 不做菜单/API 自动互相授权
- 不做统一权限码、统一权限树或菜单-API 关系表
- 不新增后端接口、数据库表或持久化映射结构
- 不归档或替换当前 `refine-dual-track-role-permissions` change

## Capabilities

### New Capabilities

- `admin-role-permission-selection-ux`: 管理员可以在保持菜单权限和 API 权限分离保存的前提下，以按页面聚合的同屏方式高效完成角色权限选择

### Modified Capabilities

- None

## Impact

- Frontend:
  - `web/src/views/admin/system/role/components/RolePermissionDrawer.vue`
  - related local view-model types or helpers inside the role permission drawer
- Backend:
  - no new endpoints
  - existing `GET /roles/:id`, `PUT /roles/:id/menus`, `PUT /roles/:id/apis`, `GET /menus`, `GET /apis/all` remain unchanged
- Persistence:
  - no schema change to `sys_role_menu`, `sys_role_api`, `casbin_rule`, `sys_menu`, or `sys_api`
- Rollback:
  - revert the drawer back to the current dual-Tab interaction without changing backend data
- Specs:
  - this change adds a new capability spec for role permission selection UX
