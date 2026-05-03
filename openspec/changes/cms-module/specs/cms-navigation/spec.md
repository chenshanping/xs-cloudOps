## ADDED Requirements

### Requirement: Navigation tree CRUD

The system SHALL provide tree-structured navigation management for the public-facing website, independent of the admin menu system (`sys_menu`).

#### Scenario: Create a navigation item
- **WHEN** admin sends POST `/api/v1/cms/navigations` with `name`, `link`, `parent_id`, `target`
- **THEN** system creates the navigation item and returns the created record

#### Scenario: List navigations as tree
- **WHEN** admin sends GET `/api/v1/cms/navigations`
- **THEN** system returns all navigation items organized as a tree structure

#### Scenario: Update navigation
- **WHEN** admin sends PUT `/api/v1/cms/navigations/:id` with updated fields
- **THEN** system updates the item and returns the updated record

#### Scenario: Delete navigation with children
- **WHEN** admin sends DELETE `/api/v1/cms/navigations/:id` and the item has children
- **THEN** system returns a 400 error indicating the item has children

#### Scenario: Delete leaf navigation
- **WHEN** admin sends DELETE `/api/v1/cms/navigations/:id` and the item has no children
- **THEN** system deletes the item successfully

### Requirement: Navigation status control

Navigation items SHALL have a status field. Only enabled items appear in public API.

#### Scenario: Disable a navigation item
- **WHEN** admin disables a navigation item
- **THEN** the item and its children no longer appear in public navigation tree
