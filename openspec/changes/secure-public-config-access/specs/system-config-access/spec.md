## ADDED Requirements

### Requirement: Anonymous clients can load only public configuration keys
The system MUST allow anonymous clients to read only a fixed allowlist of public configuration keys through the batch config read endpoint.

#### Scenario: Anonymous request reads public branding keys
- **GIVEN** a client is not authenticated
- **WHEN** the client submits a batch config read request containing only public branding or login-page keys
- **THEN** the system MUST return those public configuration values

#### Scenario: Anonymous request includes protected keys
- **GIVEN** a client is not authenticated
- **WHEN** the client submits a batch config read request containing one or more protected configuration keys
- **THEN** the system MUST NOT return any protected configuration values from that request

### Requirement: Public batch config access never escalates on login state
The system MUST NOT allow a valid login token to elevate the public batch config endpoint into a sensitive config read path.

#### Scenario: Valid token still cannot read protected keys from public endpoint
- **GIVEN** a client sends a valid authenticated request
- **WHEN** the client submits a batch config read request containing protected configuration keys to the public batch config endpoint
- **THEN** the system MUST return only keys that are publicly allowed

#### Scenario: Invalid token does not change public batch config behavior
- **GIVEN** a client sends an invalid or expired token
- **WHEN** the client submits a batch config read request containing protected configuration keys
- **THEN** the system MUST behave the same as an anonymous request for batch config access

### Requirement: Public config allowlist is backend-configurable and default-safe
The system MUST allow operators to configure which config keys are publicly readable without requiring code changes, while still preventing hard-sensitive keys from being exposed.

#### Scenario: Public allowlist is read from system config
- **GIVEN** the system stores a `public_config_keys` config value
- **WHEN** the public batch config endpoint filters requested keys
- **THEN** the system MUST use that stored allowlist as the primary source of truth

#### Scenario: Sensitive key remains private even if whitelisted by mistake
- **GIVEN** `public_config_keys` includes a sensitive key such as a password or storage secret
- **WHEN** a client requests that key through the public batch config endpoint
- **THEN** the system MUST NOT return that sensitive configuration value

### Requirement: Protected config remains available only through existing private config APIs
The system MUST continue to expose protected configuration values only through authenticated, authorized private config APIs.

#### Scenario: Config admin reads protected values from private config API
- **GIVEN** a client has valid access to the private config API
- **WHEN** the client requests the private config list
- **THEN** the system MUST return protected configuration values needed by the admin console

#### Scenario: Ordinary authenticated user cannot use private config API without permission
- **GIVEN** a client is authenticated but lacks config API permission
- **WHEN** the client calls the private config list API
- **THEN** the system MUST reject the request

### Requirement: Public application bootstrap uses only public configuration keys
The system MUST separate public bootstrap configuration loading from authenticated admin configuration loading so that unauthenticated startup paths no longer request protected config keys.

#### Scenario: Login page bootstrap requests only public config
- **GIVEN** the application is loading before the user has authenticated
- **WHEN** the frontend requests bootstrap configuration for the login page or public shell
- **THEN** the request MUST contain only public configuration keys

#### Scenario: Admin console loads additional protected config after authentication
- **GIVEN** the user has authenticated successfully and entered a protected settings screen
- **WHEN** the admin console loads configuration needed for protected settings screens
- **THEN** the frontend MAY request additional protected configuration values through the existing private config API
