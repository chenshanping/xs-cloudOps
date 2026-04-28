## Why

当前后台用户管理还存在三类直接影响使用和运维效率的缺口：

- 左侧部门树首次进入没有默认展开，也缺少显式的展开/收缩控制。
- “未绑定部门”在树与表格中没有清晰风险提示，只是普通文本。
- 重置密码仍依赖前端写死默认值，且只有单条操作，没有批量重置闭环。

同时，当前页面还存在两个容易误操作的局部问题：重置密码没有二次确认，以及搜索/翻页/切树后保留旧勾选，批量操作容易误选。

## What Changes

- 为用户管理页补齐左侧部门树默认全展开与“全部展开 / 全部收缩”控制。
- 将“未绑定部门”在用户管理页的树节点和表格列中统一改为红色标签强调。
- 新增系统配置键 `user_default_password`，作为用户管理单条/批量重置密码的默认来源。
- 新增后台用户批量重置密码能力，并将单条重置密码也改为服务端按系统配置执行。
- 为单条/批量重置密码增加二次确认，并在用户管理页相关交互后清空旧勾选。

## Scope Boundaries

- 只覆盖后台用户管理页、系统基础配置、相关用户接口与权限元数据。
- 不引入首次登录强制改密。
- 不改新增用户表单默认密码来源。
- 不扩展成完整密码策略或密码历史管理能力。
- 不做用户管理页的大范围结构重构。

## Affected Areas

- `web/src/views/admin/system/user/`
- `web/src/views/admin/system/config/components/SystemConfig.vue`
- `web/src/store/config.ts`
- `server/api/v1/user.go`
- `server/model/request/user_request.go`
- `server/router/modules/user.go`
- `server/service/user/user.go`
- `server/initialize/db_tables.go`
- `server/sql/`
- `openspec/specs/admin-user-controls`

## Rollback

- 前端移除用户管理页的树展开控制、批量重置按钮、红色未绑定标签强化与配置项输入。
- 后端移除批量重置密码接口，并把单条重置密码恢复为原先前端传入密码模式。
- 保留 `user_default_password` 配置键不消费，或后续单独清理。

## Out of Scope

- 首次登录强制修改密码
- 密码复杂度策略和历史密码校验
- 其他模块中的“未绑定部门”展示统一
- 新增用户默认密码改为系统配置
