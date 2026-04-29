## ADDED Requirements

### Requirement: Admin role permission assignment must preserve separate menu and API authorization
The system MUST let an authorized administrator assign role menu permissions and role API permissions as separate grant sets inside the admin role permission workflow.

#### Scenario: Role permission UI shows separate authorization areas
- **GIVEN** an administrator opens role permission assignment for an existing role
- **WHEN** the permission drawer loads
- **THEN** the system MUST expose a menu authorization area and a separate API authorization area

#### Scenario: Menu authorization does not implicitly grant API access
- **GIVEN** a role is granted one or more menu or button permissions
- **WHEN** no corresponding API permissions are assigned
- **THEN** the system MUST NOT treat the menu grant as backend API access

#### Scenario: API authorization does not implicitly grant menu visibility
- **GIVEN** a role is granted one or more API permissions
- **WHEN** no corresponding menu or button permissions are assigned
- **THEN** the system MUST NOT treat the API grant as frontend menu or button visibility

### Requirement: Role permission UI must explain dual-track permission semantics
The system MUST clearly communicate that menu permissions control frontend visibility while API permissions control backend interface access.

#### Scenario: Menu authorization area explains frontend scope
- **GIVEN** an administrator views the menu authorization area
- **WHEN** the area is rendered
- **THEN** the system MUST indicate that the selected menu permissions control frontend menus or buttons

#### Scenario: API authorization area explains backend scope
- **GIVEN** an administrator views the API authorization area
- **WHEN** the area is rendered
- **THEN** the system MUST indicate that the selected API permissions control backend interface access

### Requirement: Button-only menu grants must still produce usable user menus
The system MUST preserve enough menu ancestry so that button-only or child-only menu grants still result in a usable visible menu path for authorized users.

#### Scenario: Button-only grant keeps parent menu chain
- **GIVEN** a role is granted a button permission under a page menu
- **WHEN** a user with that role requests current user info
- **THEN** the returned user menu tree MUST include the required parent menu chain for that granted button

#### Scenario: Button-only grant still exposes permission code
- **GIVEN** a role is granted a button permission with a non-empty permission code
- **WHEN** a user with that role requests current user info
- **THEN** the returned frontend permission list MUST include that button permission code

### Requirement: Role API permission changes must take effect immediately
The system MUST synchronize role API authorization changes to Casbin runtime enforcement without requiring a service restart.

#### Scenario: Added API permission works immediately
- **GIVEN** an administrator assigns one or more API permissions to a role
- **WHEN** a user with that role calls one of those interfaces after the save succeeds
- **THEN** the system MUST allow the request immediately according to the new role API binding

#### Scenario: Removed API permission is denied immediately
- **GIVEN** an administrator removes one or more API permissions from a role
- **WHEN** a user with that role calls one of those interfaces after the save succeeds
- **THEN** the system MUST deny the request immediately according to the updated role API binding

### Requirement: Reopening role permission assignment must reflect saved menu and API bindings
The system MUST reload and display the currently saved role menu and role API bindings when an administrator reopens role permission assignment.

#### Scenario: Reopen permission drawer shows saved menu selections
- **GIVEN** a role already has saved menu permissions
- **WHEN** an administrator reopens role permission assignment for that role
- **THEN** the system MUST display those saved menu selections

#### Scenario: Reopen permission drawer shows saved API selections
- **GIVEN** a role already has saved API permissions
- **WHEN** an administrator reopens role permission assignment for that role
- **THEN** the system MUST display those saved API selections
