## ADDED Requirements

### Requirement: 公开认证接口必须具备频率限制

The system SHALL enforce Redis-backed rate limits on public authentication and recovery endpoints, with code-defined thresholds and dimensions appropriate to each endpoint.

#### Scenario: Register endpoint is rate limited

- **WHEN** a client repeatedly calls `POST /api/v1/auth/register` beyond the allowed threshold for the same IP / username / email within the window
- **THEN** the system returns a non-success response indicating the operation is too frequent

#### Scenario: Email code endpoint is rate limited

- **WHEN** a client repeatedly calls `POST /api/v1/auth/send-email-code` for the same email or from the same IP beyond the allowed threshold
- **THEN** the system rejects further sends within the window

### Requirement: 公开认证接口不得泄露滑动验证码答案

The system SHALL NOT expose any field on a public captcha challenge endpoint that directly reveals the server-side success position or equivalent verification answer.

#### Scenario: Slider captcha configuration is present

- **WHEN** the login captcha type is configured as `slider`
- **THEN** the public captcha config response falls back to a safe supported type instead of exposing an insecure slider challenge

### Requirement: 密码恢复接口不得通过响应暴露邮箱是否注册

The system SHALL avoid revealing whether an email address corresponds to an existing account during the email-based password reset flow.

#### Scenario: Email does not exist

- **WHEN** a client submits a valid email verification code for an email address that does not belong to any user and calls `POST /api/v1/auth/reset-password-by-email`
- **THEN** the system returns the same outward-facing success message used for an existing account
- **AND** no user record is modified

### Requirement: Token 刷新必须验证白名单一致性

The system SHALL require a refresh token request to match the current Redis whitelist token for that user, even when the JWT itself is within the configured refresh window.

#### Scenario: Current token is refreshed

- **WHEN** a client calls `POST /api/v1/auth/refresh` with the current whitelisted token within the refresh window
- **THEN** the system returns a new token

#### Scenario: Replaced old token cannot refresh

- **WHEN** a client calls `POST /api/v1/auth/refresh` with an older token that is no longer the user's current whitelisted token
- **THEN** the system rejects the refresh request
