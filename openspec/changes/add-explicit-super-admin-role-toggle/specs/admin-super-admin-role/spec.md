## ADDED Requirements

### Requirement: Roles can be explicitly marked as super admin
The system MUST let an authorized administrator explicitly mark a role as a super admin role in role management.

#### Scenario: Save role as super admin
- **GIVEN** an administrator edits or creates a role
- **WHEN** the administrator enables the super admin switch and saves
- **THEN** the system MUST persist that role as a super admin role

#### Scenario: Reopen role shows persisted super admin state
- **GIVEN** a role has a saved super admin state
- **WHEN** an administrator reopens that role in role management
- **THEN** the system MUST display the persisted super admin state

### Requirement: Super admin access must be driven by explicit role state
The system MUST determine all-access behavior from the explicit super admin role state instead of relying on role code, role ID, username, or user ID.

#### Scenario: Ordinary admin-coded role does not bypass API access
- **GIVEN** a role has code `admin` but is not marked as super admin
- **WHEN** a user with only that role calls an API that is not assigned to the role and is not otherwise whitelisted
- **THEN** the system MUST reject the request

#### Scenario: Ordinary admin-coded role does not bypass frontend button permissions
- **GIVEN** a role has code `admin` but is not marked as super admin
- **WHEN** a user with only that role requests current frontend permission codes
- **THEN** the system MUST return only the role's assigned permission codes instead of a global wildcard permission

### Requirement: Super admin roles grant full menu and API access
The system MUST treat any user who holds at least one super admin role as fully authorized for menu visibility, frontend permission codes, and backend API access.

#### Scenario: Super admin role bypasses unassigned API restriction
- **GIVEN** a user holds a role marked as super admin
- **WHEN** that role has not been explicitly assigned a target protected API
- **THEN** the system MUST still allow the request

#### Scenario: Super admin role exposes frontend wildcard permissions
- **GIVEN** a user holds a role marked as super admin
- **WHEN** the user requests current frontend permission codes
- **THEN** the system MUST return a wildcard-style full permission result

#### Scenario: Super admin role sees full available menu tree
- **GIVEN** a user holds a role marked as super admin
- **WHEN** the user requests current visible menus
- **THEN** the system MUST return the full available menu tree subject only to global feature toggles that intentionally hide modules

### Requirement: Disabling super admin must restore ordinary role behavior
The system MUST return a role to ordinary menu/API authorization behavior when its super admin flag is disabled.

#### Scenario: Disabled super admin no longer bypasses
- **GIVEN** a role was previously marked as super admin and later has the flag disabled
- **WHEN** a user with only that role accesses an unassigned protected API
- **THEN** the system MUST enforce ordinary role API permissions and reject the request if not assigned

#### Scenario: Disabled super admin preserves prior explicit assignments
- **GIVEN** a role has existing menu or API assignments and the administrator disables the super admin flag
- **WHEN** the role is used afterwards
- **THEN** the system MUST continue to honor the preserved explicit assignments without requiring the administrator to rebind them
