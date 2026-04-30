## Context

当前项目在应用启动、登录页、路由守卫和后台页面中都复用了 `configStore.loadConfigs()`。这条链路会通过公开的 `POST /api/v1/configs/keys` 加载一整组配置键。

现状问题不在“前端显不显示”，而在“后端匿名请求可以读什么”：

- `server/router/modules/config.go` 将 `/configs/keys` 注册为公开接口
- `server/api/v1/config.go` 直接将请求中的 `keys` 透传给 `service.Config.GetConfigsByKeys`
- `web/src/store/config.ts` 当前把登录页品牌配置和后台敏感配置写在同一个 `CONFIG_KEYS` 常量中

本次设计必须保证三件事同时成立：

- 登录前依然能拿到品牌、登录页和前台模式相关配置
- 普通登录用户访问公开批量配置接口时，也不能读取 SMTP、对象存储、默认密码等后台配置
- 运营侧后续调整公开配置范围时，不需要再改后端代码

## Decisions

### 1. 保留现有公开接口路径，但统一收敛到服务端白名单

决策：

- 不直接把 `/configs/keys` 改成私有接口
- `POST /configs/keys` 对所有调用方都只返回服务端判定为公开的配置键
- 公开范围不再硬编码在接口逻辑里，而是优先读取数据库配置 `public_config_keys`

原因：

- 这是当前风险最低的修复方式，既能立即收口漏洞，又不会打断登录页和应用启动
- 安全边界必须放在后端，不能依赖前端筛选
- 普通用户带有效 Token 也不能因为“已登录”自动升级为可读敏感配置

备选方案：

- 直接改私有接口：会破坏未登录阶段的配置加载，排除
- 立即拆成两条新接口：方向正确，但本轮改动面更大，排除

### 2. 公开白名单改为后台可配置，但保留后端硬性保护

决策：

- 新增系统配置项 `public_config_keys`，值为 JSON 数组，由后台配置页维护
- 后端读取该配置作为公开 allowlist
- 后端再叠加一层 `never public` 强制拒绝集合，敏感配置即使被误填进白名单也不会公开

原因：

- 满足“以后改公开配置范围不必改代码”的诉求
- 防止后台误操作把邮箱密码、对象存储密钥、AI Key 等真正敏感配置暴露出去

### 3. 后台完整配置继续走现有私有配置接口

决策：

- 不新增新的私有批量配置接口
- 后台完整配置继续通过现有 `GET /configs` 私有接口获取
- 配置页在需要完整配置时显式调用私有接口；其他页面不再默认全局拉取后台配置

原因：

- 现有私有配置接口已经挂在 `JWTAuth + CasbinAuth` 链路下
- 复用现有权限模型，避免扩大 API 面
- 避免普通已登录用户因为全局初始化动作而隐式请求敏感配置

### 4. 前端拆分公开配置集和后台配置集

决策：

- 将 `web/src/store/config.ts` 的单一 `CONFIG_KEYS` 拆成 `PUBLIC_CONFIG_KEYS` 与 `ADMIN_CONFIG_KEYS`
- `loadConfigs()` 默认只加载公开配置
- 只有系统配置页等确实需要完整配置的页面，才通过私有 `GET /configs` 补充加载后台配置
- 路由守卫和普通后台页面不再默认全局拉取后台敏感配置

原因：

- 只有后端白名单还不够；如果前端仍把敏感配置作为全局启动依赖，职责边界依然混乱
- 拆分 key 集可以把“哪些配置用于公开启动”固定在前端实现里，但真正的安全判断仍以后端为准

### 5. 公开配置白名单使用 allowlist，而不是黑名单

决策：

- 后端只维护允许匿名读取的固定 key 集
- 不做 `password` / `secret` / `token` 之类的关键字黑名单过滤

原因：

- 黑名单容易漏网且不可维护
- allowlist 更符合当前登录页和品牌配置的稳定场景

## Backend Design

- 在 `server/service/configsvc/config.go`：
  - 增加默认公开配置键集合
  - 增加 `public_config_keys` 配置读取与 JSON 解析
  - 增加 `never public` 强制拒绝集合
  - 统一批量过滤并构造成 `map[string]model.SysConfig`
- 在 `server/api/v1/config.go`：
  - `GetConfigsByKeys` 只返回公开白名单命中的配置
  - 不再根据 Token 有效性提升读取范围
- 保持 `GET /configs`、`GET /configs/key/:key`、`PUT /configs/batch` 等后台私有接口原有鉴权链路不变
- 在启动补齐逻辑中补充 `public_config_keys` 默认配置，并保证已有自定义值不被覆盖

## Frontend Design

- `web/src/store/config.ts`
  - 拆分公开和后台配置 key 常量
  - 默认加载公开配置
  - 当页面需要完整后台配置时，通过私有 `GET /configs` 补充加载，并保证多次调用不会丢失已加载的公开配置
- `web/src/views/admin/system/config/components/SystemConfig.vue`
  - 增加“公开配置键”多行输入框
  - 以一行一个键的方式维护 `public_config_keys`
- 现有启动链路继续能读取：
  - `sys_name`
  - `sys_logo`
  - `register_logo`
  - `login_*`
  - `enable_register`
  - `front_mode`
- `user_profile_button_visible` 保留在公开配置集，避免普通用户管理页对后台敏感配置接口形成依赖
- 邮箱、文件存储、默认密码等后台专用配置只在系统配置页等私有场景补充读取

## Persistence

- 不新增表结构
- 不修改 `sys_config` 表字段
- 不新增 SQL 脚本

## Validation

- 请求体仍要求 JSON 格式的 `keys []string`
- 公开接口不报“缺少权限”来暴露内部规则，只返回允许公开的配置结果
- 有效 Token 不会提升公开接口的读取权限
- 完整后台配置仍依赖现有私有接口和 Casbin 权限链路

## Risks

- 如果遗漏某个登录前必需的公开 key，登录页或前台模式可能退回默认值
- 如果后台把真正需要公开的键从 `public_config_keys` 中删掉，对应公开页面会回退到默认值
- 如果普通后台页面仍误用完整配置加载路径，会遇到私有配置权限限制

## Verification

- 后端测试覆盖：
  - 公开接口按 `public_config_keys` 返回公开配置
  - 白名单中误填敏感键时，敏感值仍不会公开
  - 有效登录 Token 调用公开接口也不会读取到敏感配置
  - 私有 `GET /configs` 仍受权限链路保护
- 前端验证覆盖：
  - 应用启动和登录页仍可正常展示品牌与登录配置
  - 系统配置页可维护 `public_config_keys`
  - 非配置页面不再因为全局初始化动作请求后台敏感配置
