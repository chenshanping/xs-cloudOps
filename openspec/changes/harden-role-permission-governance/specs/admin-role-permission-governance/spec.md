## ADDED Requirements

### Requirement: Admin role permissions must be saved atomically across menus, direct APIs, and data scopes
The system MUST let an authorized administrator save role menu permissions, direct API permissions, and data-scope permissions as a single atomic permission update.

#### Scenario: Unified save persists all role permission inputs together
- **GIVEN** an administrator edits role permissions
- **WHEN** the administrator saves menu selections, direct API selections, and data-scope selections together
- **THEN** the system MUST persist menu permissions, direct API permissions, and data-scope permissions as one successful role permission update

#### Scenario: Failed unified save does not partially apply role permissions
- **GIVEN** an administrator submits a role permission update
- **WHEN** any part of the role permission persistence or runtime policy synchronization fails
- **THEN** the system MUST leave the role's stored permissions and effective API authorization unchanged

### Requirement: Final role API access must combine direct and inherited grants immediately
The system MUST synchronize role API authorization changes to Casbin runtime enforcement without requiring a service restart.

#### Scenario: Final role API access includes direct API grants
- **GIVEN** a role is granted one or more direct API permissions
- **WHEN** a user with that role calls one of those interfaces after the save succeeds
- **THEN** the system MUST allow the request immediately according to the saved direct API binding

#### Scenario: Final role API access includes menu-inherited API grants
- **GIVEN** a role is granted one or more menus that have associated APIs
- **WHEN** a user with that role calls one of those associated interfaces after the save succeeds
- **THEN** the system MUST allow the request immediately according to the inherited API binding

#### Scenario: Removed inherited API is denied immediately
- **GIVEN** a role no longer holds a menu that provided an inherited API permission and the same API was not directly granted
- **WHEN** a user with that role calls that interface after the save succeeds
- **THEN** the system MUST deny the request immediately according to the updated effective API binding
