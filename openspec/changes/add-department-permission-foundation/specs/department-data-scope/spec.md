## ADDED Requirements

### Requirement: Roles define user data scope
The system MUST allow each role to define a data scope that determines which users an operator can view and manage.

#### Scenario: Save role with department data scope
- **GIVEN** an administrator edits a role
- **WHEN** the administrator selects a supported data scope and saves the role
- **THEN** the system MUST persist the role data scope

#### Scenario: Save custom department scope
- **GIVEN** an administrator selects the custom department data scope for a role
- **WHEN** the administrator submits one or more departments for that role
- **THEN** the system MUST persist the selected department range for that role

### Requirement: User management list must enforce data scope
The system MUST filter the admin user list by the current operator's effective data scope.

#### Scenario: Filter by current department
- **GIVEN** the operator's effective data scope is current department only
- **WHEN** the operator requests the user list
- **THEN** the system MUST return only users belonging to the operator's own department

#### Scenario: Filter by current department and descendants
- **GIVEN** the operator's effective data scope is current department and descendants
- **WHEN** the operator requests the user list
- **THEN** the system MUST return only users belonging to the operator's department or any descendant department

#### Scenario: Filter by self only
- **GIVEN** the operator's effective data scope is self only
- **WHEN** the operator requests the user list
- **THEN** the system MUST return only the operator's own user record

### Requirement: User management mutations must enforce manageable scope
The system MUST reject user management operations that target users or departments outside the operator's manageable scope.

#### Scenario: Reject out-of-scope user detail
- **GIVEN** a target user is outside the operator's effective data scope
- **WHEN** the operator requests that user detail
- **THEN** the system MUST reject the request

#### Scenario: Reject assigning out-of-scope department
- **GIVEN** an operator cannot manage a target department under the effective data scope rules
- **WHEN** the operator creates or updates a user with that department assignment
- **THEN** the system MUST reject the request

#### Scenario: Reject out-of-scope mutation
- **GIVEN** one or more target users are outside the operator's effective data scope
- **WHEN** the operator submits edit, delete, status, reset password, force offline, or profile-view actions
- **THEN** the system MUST reject the request and MUST NOT partially mutate the targets

### Requirement: Existing data remains compatible after upgrade
The system MUST preserve current effective access for existing users and roles after the department foundation is introduced.

#### Scenario: Existing role keeps broad access
- **GIVEN** a pre-existing role from before the upgrade
- **WHEN** the upgrade initializes the new data scope field
- **THEN** the role MUST default to full data scope unless explicitly changed later

#### Scenario: Existing user gets valid department
- **GIVEN** a pre-existing user without a department assignment before the upgrade
- **WHEN** the upgrade initializes department data
- **THEN** the user MUST be assigned to the default root department so department-based filtering has a deterministic result
