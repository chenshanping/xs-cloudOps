## ADDED Requirements

### Requirement: Page CRUD

The system SHALL provide single-page management for static content pages (e.g., About Us, Contact).

#### Scenario: Create a page
- **WHEN** admin sends POST `/api/v1/cms/pages` with `title`, `slug`, `content`
- **THEN** system creates the page and returns the created record

#### Scenario: Slug uniqueness
- **WHEN** admin creates a page with a `slug` that already exists
- **THEN** system returns a 400 error indicating slug is duplicated

#### Scenario: List pages
- **WHEN** admin sends GET `/api/v1/cms/pages`
- **THEN** system returns all pages sorted by `sort` ASC

#### Scenario: Get page detail
- **WHEN** admin sends GET `/api/v1/cms/pages/:id`
- **THEN** system returns full page detail including content

#### Scenario: Update page
- **WHEN** admin sends PUT `/api/v1/cms/pages/:id` with updated fields
- **THEN** system updates the page and returns the updated record

#### Scenario: Delete page
- **WHEN** admin sends DELETE `/api/v1/cms/pages/:id`
- **THEN** system deletes the page permanently

### Requirement: Page status control

Pages SHALL have a status field (enabled/disabled). Only enabled pages appear in public API responses.

#### Scenario: Disable a page
- **WHEN** admin updates a page's status to disabled
- **THEN** the page no longer appears in public API responses

#### Scenario: Public access to disabled page
- **WHEN** public API requests a page by slug that is disabled
- **THEN** system returns a 404 error
