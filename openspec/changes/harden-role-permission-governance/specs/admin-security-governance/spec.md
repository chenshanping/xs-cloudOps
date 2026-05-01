## ADDED Requirements

### Requirement: Casbin bypass rules must use explicit allowlisted routes
The system MUST restrict authenticated authorization bypasses to explicit allowlisted routes.

#### Scenario: Explicit allowlisted route still bypasses role API authorization
- **GIVEN** a route is configured as an explicit authenticated whitelist route
- **WHEN** a logged-in user calls that route
- **THEN** the system MUST allow the request without requiring role API authorization

#### Scenario: Suffix-only route match does not bypass authorization
- **GIVEN** a route is not explicitly allowlisted
- **WHEN** its path happens to share a suffix with another route pattern
- **THEN** the system MUST still enforce normal role API authorization

### Requirement: Operation logs must sanitize sensitive request and response data
The system MUST prevent common sensitive values from being stored in operation log payload bodies.

#### Scenario: Sensitive request fields are masked before persistence
- **GIVEN** a logged request contains sensitive fields such as password, token, authorization, secret, or key material
- **WHEN** the operation log entry is persisted
- **THEN** the system MUST store masked values instead of the original sensitive content

#### Scenario: Sanitized logging still records useful audit metadata
- **GIVEN** the system writes an operation log for a protected request
- **WHEN** the log is later viewed for audit or troubleshooting
- **THEN** the log entry MUST still preserve route, actor, status, latency, and non-sensitive request context needed for diagnosis
