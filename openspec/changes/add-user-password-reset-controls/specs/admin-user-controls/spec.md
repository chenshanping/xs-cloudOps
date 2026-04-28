## MODIFIED Requirements

### Requirement: Admin can batch change user status
The system MUST clear stale list selections after admin batch actions or list-context changes that would otherwise make the previous selection ambiguous.

#### Scenario: Clear selected rows after search or tree filter change
- **GIVEN** an administrator has selected one or more users in the admin user list
- **WHEN** the administrator performs a search, resets filters, changes page, or switches the department tree filter
- **THEN** the system MUST clear the previous row selection before the next batch action can be submitted

#### Scenario: Clear selected rows after successful batch action
- **GIVEN** an administrator has selected one or more users in the admin user list
- **WHEN** the administrator successfully completes a batch enable, batch disable, batch delete, or batch reset password action
- **THEN** the system MUST clear the previous row selection

### Requirement: Admin can hide or show the profile button from system settings
The system MUST also provide a configurable default password for admin-triggered password resets.

#### Scenario: Default reset password config exists
- **GIVEN** the system is using default configuration values
- **WHEN** an administrator opens the system basic configuration page
- **THEN** the system MUST expose a `user_default_password` setting with a default fallback value

## ADDED Requirements

### Requirement: User management department tree defaults to expanded navigation
The system MUST show the user management department tree in a fully expanded state on initial page load and MUST let administrators expand or collapse the full tree explicitly.

#### Scenario: Initial page load expands all department nodes
- **GIVEN** the administrator opens the admin user management page
- **WHEN** the department tree is rendered for the first time
- **THEN** the system MUST expand all currently visible department nodes by default

#### Scenario: Collapse all department nodes
- **GIVEN** the department tree is visible on the admin user management page
- **WHEN** the administrator triggers the collapse-all control
- **THEN** the system MUST collapse the department tree to the root-level view

#### Scenario: Expand all department nodes
- **GIVEN** the department tree is visible on the admin user management page
- **WHEN** the administrator triggers the expand-all control
- **THEN** the system MUST expand all currently visible department nodes again

### Requirement: Unassigned department users are visually emphasized
The system MUST render unassigned department state in a visually emphasized warning style inside the admin user management page.

#### Scenario: Unassigned user row shows warning tag
- **GIVEN** a user record has no bound department
- **WHEN** the administrator views the user list
- **THEN** the department column MUST render a red warning tag indicating the user is unassigned

#### Scenario: Unassigned tree node shows warning style
- **GIVEN** the user management department tree contains the unassigned users node
- **WHEN** the administrator views the tree
- **THEN** the unassigned node MUST be rendered in a red warning style distinct from normal department nodes

### Requirement: Admin can reset passwords using configured default password
The system MUST reset admin-managed user passwords using the configured default password instead of a frontend hard-coded value.

#### Scenario: Single reset uses configured default password
- **GIVEN** the system configuration contains a `user_default_password` value
- **WHEN** an administrator confirms a single-user reset password action
- **THEN** the system MUST reset the target user password to the configured default password

#### Scenario: Single reset falls back to built-in default
- **GIVEN** the system configuration does not contain a `user_default_password` value
- **WHEN** an administrator confirms a single-user reset password action
- **THEN** the system MUST reset the target user password to the built-in fallback default password

### Requirement: Admin can batch reset passwords
The system MUST allow an authorized administrator to batch reset passwords for selected manageable users from the admin user management page.

#### Scenario: Batch reset selected users
- **GIVEN** an administrator has selected one or more users in the admin user list
- **WHEN** the administrator confirms a batch reset password action
- **THEN** the system MUST reset all selected manageable users to the configured default password in one request

#### Scenario: Reject batch reset outside manageable scope
- **GIVEN** one or more selected user IDs are outside the current operator's manageable scope
- **WHEN** the operator submits a batch reset password request
- **THEN** the system MUST reject the request and MUST NOT reset any selected user password

### Requirement: Reset password actions require explicit confirmation
The system MUST require explicit confirmation before executing single or batch reset password actions from the admin user management page.

#### Scenario: Confirm single reset password
- **GIVEN** an administrator chooses the single-user reset password action
- **WHEN** the confirmation dialog is shown
- **THEN** the dialog MUST clearly state that the user password will be reset to the configured default password

#### Scenario: Confirm batch reset password
- **GIVEN** an administrator chooses the batch reset password action
- **WHEN** the confirmation dialog is shown
- **THEN** the dialog MUST clearly state that all selected user passwords will be reset to the configured default password
