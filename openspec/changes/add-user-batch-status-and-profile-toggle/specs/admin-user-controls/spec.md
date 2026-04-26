## ADDED Requirements

### Requirement: Admin can batch change user status
The system MUST allow an authorized administrator to batch enable or batch disable users from the admin user management page.

#### Scenario: Batch disable selected users
- **GIVEN** an administrator has selected one or more users in the admin user list
- **WHEN** the administrator confirms a batch disable action
- **THEN** the system MUST disable all selected users in one request

#### Scenario: Batch enable selected users
- **GIVEN** an administrator has selected one or more users in the admin user list
- **WHEN** the administrator confirms a batch enable action
- **THEN** the system MUST enable all selected users in one request

### Requirement: Batch disable must protect operator and protected admin accounts
The system MUST reject a batch disable request that attempts to disable the current operator or a protected administrator account.

#### Scenario: Reject disabling current operator
- **GIVEN** the current operator is included in the selected user IDs
- **WHEN** the operator submits a batch disable request
- **THEN** the system MUST reject the request and MUST NOT disable any selected user

#### Scenario: Reject disabling protected admin account
- **GIVEN** the selected user IDs include a protected administrator account
- **WHEN** the operator submits a batch status request
- **THEN** the system MUST reject the request and MUST NOT change any selected user status

### Requirement: Batch disable must invalidate active login tokens
The system MUST invalidate active login tokens for users that are successfully batch disabled.

#### Scenario: Disabled user is forced offline
- **GIVEN** a selected user is currently logged in
- **WHEN** the administrator successfully batch disables that user
- **THEN** the user's active token MUST be invalidated so subsequent authenticated requests are rejected

### Requirement: Batch status changes must be auditable
The system MUST record the operator identity and requested target user IDs for batch status changes in the existing operation log flow.

#### Scenario: Batch status request is logged
- **GIVEN** an administrator sends a batch user status request
- **WHEN** the request completes
- **THEN** the operation log MUST contain the operator identity, request path, and request body including target user IDs and status

### Requirement: Admin can hide or show the profile button from system settings
The system MUST provide a system configuration switch that controls whether the “身份” button is shown on the admin user management page, and the default state MUST be hidden.

#### Scenario: Default hidden profile button
- **GIVEN** the system is using default configuration values
- **WHEN** an administrator opens the admin user management page
- **THEN** the “身份” button MUST be hidden

#### Scenario: Show profile button after enabling config
- **GIVEN** the profile button visibility setting is enabled in system settings
- **WHEN** an administrator opens the admin user management page
- **THEN** the “身份” button MUST be shown for each user row
