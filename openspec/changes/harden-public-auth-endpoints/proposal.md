## Why

当前仓库的公开认证接口已存在一处已修复问题（关闭注册后仍可直接注册）和一处已修复问题（用户名重置密码验证码绕过），但公开接口面上仍有 4 类未处理风险：

- 滑动验证码接口直接暴露 `target_x`，导致挑战答案可被前端或脚本直接读取
- 注册、发验证码、密码重置、刷新 Token 等公开接口缺少频率限制
- 邮箱找回密码接口会暴露“该邮箱未注册”，存在账号枚举风险
- `auth/refresh` 刷新逻辑对白名单一致性校验不完整，弱化当前单用户单有效 Token 模型

这些问题都落在认证、验证码、会话边界上，属于高风险行为修复，必须以 code-defined policy 方式统一处理，而不是停留在文档提醒。

## What Changes

- **新增** 公开认证接口频率限制能力，按接口场景对 IP / 用户名 / 邮箱 / Token 维度做 Redis 固定窗口限流
- **修改** 滑动验证码策略：当前版本不再对外提供可直接恢复答案的滑动验证码挑战；登录验证码类型若配置为 `slider`，后端安全降级到 `digit`
- **修改** `POST /api/v1/auth/reset-password-by-email`：对不存在邮箱返回统一成功文案，避免账号枚举
- **修改** JWT 白名单 TTL 与刷新校验逻辑：刷新前无条件校验白名单一致性，并让白名单存活时间覆盖刷新窗口
- **新增** 后端回归测试，覆盖限流、邮箱枚举保护、刷新白名单保护、滑动验证码降级行为
- **修改** 前端验证码消费：移除对 `target_x` 的依赖类型，并适配后端对 slider 的安全降级

## Capabilities

### Modified Capabilities

- `public-auth-security`: 公开认证接口在验证码、限流、账号枚举控制、Token 刷新一致性上遵循统一安全边界。

## Impact

**后端：**
- `server/service/auth/`：新增公开认证限流与会话保护逻辑
- `server/service/captcha/`：调整滑动验证码类型返回策略
- `server/api/v1/auth.go`、`server/api/v1/captcha.go`：接入限流、枚举保护和 slider 降级
- `server/utils/jwt.go`：修复刷新白名单一致性与 TTL 策略

**前端：**
- `web/src/api/captcha.ts`
- `web/src/components/SliderCaptcha.vue`
- `web/src/views/auth/login/index.vue`

**安全策略：**
- 频率限制阈值写在代码中，不走系统配置页面
- 不新增数据库配置开关，不允许管理员动态放宽认证边界

## Out of Scope

- 不引入新的验证码供应商
- 不新增短信验证码、设备指纹、风控平台
- 不做全项目所有接口的统一限流中间件，只覆盖公开认证接口
- 不在本次变更中重写前端验证码 UI 体系，只做安全降级与兼容
