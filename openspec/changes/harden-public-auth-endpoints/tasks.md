## 1. Spec & skill

- [ ] 1.1 新增 OpenSpec change `harden-public-auth-endpoints`
- [ ] 1.2 更新 `security-and-hardening` skill，补充：
  - 绑定校验错误需统一翻译为用户可读提示
  - 公开认证接口默认检查：服务端兜底、限流、账号枚举、会话白名单一致性

## 2. Backend hardening

- [ ] 2.1 增加公开认证接口 Redis 限流能力与测试
- [ ] 2.2 调整 `auth.go`：注册、发验证码、邮箱找回、用户名找回、refresh 接入限流
- [ ] 2.3 调整 `captcha.go` / `captcha service`：slider 安全降级，不再暴露 `target_x`
- [ ] 2.4 调整 `jwt.go`：白名单 TTL 覆盖 refresh window，刷新前无条件校验白名单一致性
- [ ] 2.5 调整邮箱找回密码接口统一成功响应，补测试

## 3. Frontend compatibility

- [ ] 3.1 更新 `web/src/api/captcha.ts` 的 slider 类型
- [ ] 3.2 确认登录页在后端返回 `digit` 时自动走普通图形验证码流程
- [ ] 3.3 清理前端对 `target_x` 的类型依赖

## 4. Verification

- [ ] 4.1 运行后端相关单测
- [ ] 4.2 运行 `go build ./...`
- [ ] 4.3 给用户明确列出登录 / 找回密码 / 注册 / refresh 的手工验收清单
