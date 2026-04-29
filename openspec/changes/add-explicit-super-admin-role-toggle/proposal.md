## Why

当前项目里“超管”语义来自隐式特判：历史上既出现过 `admin` 角色名硬编码放行，也出现过按角色 ID 直接返回全部前端权限。这种方式有两个问题：

- 行为不可见，管理员无法在角色管理里明确知道哪个角色是超管
- 角色名、角色 ID 与“是否超管”耦合，修改角色绑定后容易出现越权或误判

用户已经明确新的目标：在角色管理里提供一个显式开关，由管理员决定某个角色是否为超管。只有显式开启的角色，才拥有全部菜单与全部 API 的直接访问能力；普通角色即使编码是 `admin`，也必须继续按绑定的菜单/API 权限生效。

## What Changes

- 为角色新增显式超管标记字段，例如 `is_super_admin`
- 角色新增/编辑 Drawer 增加“是否超管”开关
- 角色详情、角色列表、角色请求/响应结构补充超管标记
- 后端菜单权限、前端按钮权限、Casbin API 鉴权统一改为读取显式超管标记，而不是依赖 `admin` 名称或角色 ID
- 角色权限抽屉继续保留菜单权限与 API 权限分离；但当角色被标记为超管时，系统必须把它视为拥有全部菜单和全部 API
- 启动补齐、升级脚本、默认内置角色逻辑需要明确处理历史 `admin` / `system_admin` 的迁移口径

### Out of Scope

- 不重构现有 Casbin 主链路
- 不改“双轨权限”模型，菜单权限和 API 权限仍然分离存储
- 不新增统一权限码或自动菜单/API 映射
- 不把数据权限范围 `data_scope` 改造成超管语义来源
- 不处理业务模块级别的更多特殊角色体系

## Capabilities

### New Capabilities

- `admin-super-admin-role`: 管理员可以显式标记某个角色为超管角色，超管角色直接拥有全部菜单和全部 API 访问能力

### Modified Capabilities

- `admin-role-permission-management`: 角色权限分配逻辑增加显式超管角色语义，但仍保持菜单/API 分离模型

## Impact

- Backend:
  - `server/model/sys_role.go`
  - `server/model/request/role_request.go`
  - `server/api/v1/role.go`
  - `server/service/role/role.go`
  - `server/service/menu/menu.go`
  - `server/middleware/casbin.go`
  - role/bootstrap/Casbin related tests and startup repair paths
- Frontend:
  - `web/src/types/index.ts`
  - `web/src/api/role.ts`
  - `web/src/views/admin/system/role/components/RoleFormDrawer.vue`
  - `web/src/views/admin/system/role/hooks/useRolePage.ts`
  - role list rendering where super admin status should be visible if needed
- Persistence:
  - `sys_role` requires a persistent super admin flag
  - upgrade path must be duplicate-safe and compatible with existing environments
- Rollback:
  - revert to ordinary role-only behavior by disabling or removing explicit super admin handling
  - rollback must not depend on restoring hidden `admin`-name hardcoding
- Specs:
  - add a new capability spec for explicit super admin roles
  - update role permission management behavior where it references effective all-access handling
