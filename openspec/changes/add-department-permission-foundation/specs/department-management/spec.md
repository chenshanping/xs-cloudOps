## ADDED Requirements

### Requirement: Admin can manage a department tree
The system MUST allow an authorized administrator to view, create, edit, and delete departments in a tree structure from the admin console.

#### Scenario: View department tree
- **GIVEN** an administrator has permission to access department management
- **WHEN** the administrator opens the department management page
- **THEN** the system MUST return and display the full department tree in sort order

#### Scenario: Create child department
- **GIVEN** an administrator is editing a valid parent department
- **WHEN** the administrator submits a new child department with required fields
- **THEN** the system MUST create the department under that parent and show it in the tree

#### Scenario: Edit department
- **GIVEN** an administrator selects an existing department
- **WHEN** the administrator updates the department name, sort, status, remark, or parent
- **THEN** the system MUST persist the change and return the updated tree data

### Requirement: Department hierarchy must stay valid
The system MUST prevent invalid department tree mutations that would create cycles or orphaned hierarchy state.

#### Scenario: Reject cycle
- **GIVEN** a department already has one or more descendants
- **WHEN** an administrator attempts to move that department under one of its descendants
- **THEN** the system MUST reject the request and MUST NOT change the tree

#### Scenario: Reject missing parent
- **GIVEN** an administrator submits a department update with a non-existent parent department
- **WHEN** the request is validated
- **THEN** the system MUST reject the request and MUST NOT persist the change

### Requirement: Department deletion must enforce usage constraints
The system MUST reject department deletion when the department still has child departments or assigned users.

#### Scenario: Reject deleting department with children
- **GIVEN** a department has one or more child departments
- **WHEN** an administrator requests deletion
- **THEN** the system MUST reject the delete request

#### Scenario: Reject deleting department with assigned users
- **GIVEN** one or more users are assigned to the target department
- **WHEN** an administrator requests deletion
- **THEN** the system MUST reject the delete request

### Requirement: Department module exposure can be hidden without removing the foundation
The system MUST support hiding department management menus and page exposure by configuration without removing the underlying department and data-scope model.

#### Scenario: Hide department menu by config
- **GIVEN** the department module display config is disabled
- **WHEN** a user refreshes dynamic menus
- **THEN** the department management menu MUST NOT be exposed in the visible navigation

#### Scenario: Keep core model active while menu hidden
- **GIVEN** the department module display config is disabled
- **WHEN** user management or role data-scope logic executes
- **THEN** the system MUST continue using department and data-scope model data internally
