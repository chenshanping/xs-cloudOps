## ADDED Requirements

### Requirement: Category tree CRUD

The system SHALL provide tree-structured category management for CMS content. Categories support unlimited nesting via `parent_id`. Each category has a unique `slug` for URL identification.

#### Scenario: Create a top-level category
- **WHEN** admin sends POST `/api/v1/cms/categories` with `name`, `slug`, `parent_id=0`
- **THEN** system creates the category and returns the created record with assigned ID

#### Scenario: Create a child category
- **WHEN** admin sends POST `/api/v1/cms/categories` with `parent_id` pointing to an existing category
- **THEN** system creates the category as a child of the specified parent

#### Scenario: Slug uniqueness
- **WHEN** admin creates a category with a `slug` that already exists
- **THEN** system returns a 400 error indicating slug is duplicated

#### Scenario: List categories as tree
- **WHEN** admin sends GET `/api/v1/cms/categories`
- **THEN** system returns all categories organized as a tree structure with `children` arrays

#### Scenario: Update category
- **WHEN** admin sends PUT `/api/v1/cms/categories/:id` with updated fields
- **THEN** system updates the category and returns the updated record

#### Scenario: Delete category with children
- **WHEN** admin sends DELETE `/api/v1/cms/categories/:id` and the category has children
- **THEN** system returns a 400 error indicating the category has children and cannot be deleted

#### Scenario: Delete category with articles
- **WHEN** admin sends DELETE `/api/v1/cms/categories/:id` and articles exist under this category
- **THEN** system returns a 400 error indicating the category has associated articles

#### Scenario: Delete empty category
- **WHEN** admin sends DELETE `/api/v1/cms/categories/:id` and no children or articles exist
- **THEN** system deletes the category successfully

### Requirement: Category status control

Categories SHALL have a status field (enabled/disabled). Disabled categories and their articles are hidden from public API responses.

#### Scenario: Disable a category
- **WHEN** admin updates a category's status to disabled
- **THEN** the category and its articles no longer appear in public API responses
- **THEN** the category remains visible in admin management pages

#### Scenario: Enable a category
- **WHEN** admin updates a category's status to enabled
- **THEN** the category and its published articles become visible in public API responses
