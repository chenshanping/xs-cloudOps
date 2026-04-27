## Why

当前系统还没有部门、组织树、用户归属部门和角色数据范围模型，用户管理与后续业务列表也没有统一的数据权限过滤底座。一旦部门参与用户数据权限，它就会进入用户、角色、查询过滤和权限校验的核心链路，因此必须按主系统正式能力建设，而不能做成可随时拔掉的插件。

## What Changes

- 新增部门基础模型，建立 `sys_dept` 树结构，并为用户增加 `dept_id` 归属。
- 为角色增加 `data_scope`，并补充自定义部门范围关联，形成统一数据权限语义。
- 新增后端部门管理模块，提供部门树、详情、创建、编辑、删除接口，以及必要的校验和删除限制。
- 新增统一的数据权限解析与查询过滤入口，并先接入用户管理相关接口，防止直接调接口绕过范围控制。
- 新增前端部门管理页面，按当前后台动态菜单模式接入；创建/编辑交互使用 `Drawer`，表单内容拆到局部 `components/`。
- 新增部门管理菜单、按钮权限、API 元数据和可选的 `dept_module_enabled` 配置开关；该开关只控制菜单/页面暴露，不影响底层部门与数据权限模型存在。

### Out of Scope

- 不在本轮把所有业务表和所有列表页都改为部门数据权限过滤。
- 不引入真正可拔插的底层权限插件机制。
- 不重构现有 Casbin 菜单/API 功能权限模型，只在其之上补齐数据权限底座。
- 不把组织机构扩展成包含岗位、编制、汇报线等更复杂的人事系统。

## Capabilities

### New Capabilities

- `department-management`: 后台可以维护部门树、控制部门模块暴露，并对删除/层级关系做安全限制。
- `department-data-scope`: 系统根据角色数据范围和用户部门归属统一过滤用户管理数据与受保护操作。

### Modified Capabilities

- None. Repository has no existing OpenSpec capability for department or data-scope behavior yet.

## Impact

- Backend:
  - `server/model/`
  - `server/model/request/`
  - `server/service/`
  - `server/api/v1/`
  - `server/router/modules/`
  - `server/initialize/db_tables.go`
- Frontend:
  - `web/src/api/`
  - `web/src/types/`
  - `web/src/views/system/dept/`
  - `web/src/views/system/user/`
  - `web/src/views/system/role/`
  - `web/src/store/config.ts`
- Persistence:
  - `sys_dept`
  - `sys_role_dept`
  - `sys_user.dept_id`
  - `sys_role.data_scope`
- Upgrade / seed:
  - incremental SQL scripts and built-in menu/API/config bootstrap
- Rollback:
  - 回滚菜单、页面与接口暴露
  - 回滚用户管理对数据权限的接入
  - 保留或单独清理新增表字段需通过数据库升级脚本处理，不能直接删除线上数据结构
