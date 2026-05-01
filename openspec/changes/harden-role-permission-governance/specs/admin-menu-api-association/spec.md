## ADDED Requirements

### Requirement: Menu and button permissions may define associated backend APIs
The system MUST let an authorized administrator maintain explicit backend API associations for menu and button permissions.

#### Scenario: Menu management exposes associated API configuration
- **GIVEN** an administrator edits a menu or button permission item
- **WHEN** the item is eligible for backend interaction
- **THEN** the system MUST allow the administrator to view and update the associated API set for that menu item

#### Scenario: Directory menu does not require API associations
- **GIVEN** an administrator edits a directory-type menu
- **WHEN** the directory is displayed in menu management
- **THEN** the system MUST NOT require backend API associations for that directory item

### Requirement: Role menus inherit associated APIs without collapsing menu and API semantics
The system MUST treat menu-associated APIs as inherited backend access for roles that hold the menu while preserving direct API grants as a separate authorization source.

#### Scenario: Menu grant inherits associated APIs
- **GIVEN** a menu has one or more associated APIs
- **AND** a role is granted that menu
- **WHEN** the system computes the role's effective backend API access
- **THEN** the associated APIs MUST be included in the role's effective API authorization

#### Scenario: Direct API grant remains separately manageable
- **GIVEN** a role has a direct API grant that is not inherited from any selected menu
- **WHEN** the role permission assignment is reloaded
- **THEN** the system MUST preserve and display that direct API grant as a separately managed authorization

#### Scenario: Inherited API alone does not imply frontend visibility
- **GIVEN** a role gains backend API access only through a menu association model
- **WHEN** the role lacks the corresponding menu or button grant
- **THEN** the system MUST NOT treat inherited backend API access as frontend menu or button visibility
