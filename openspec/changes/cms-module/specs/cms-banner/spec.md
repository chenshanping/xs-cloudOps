## ADDED Requirements

### Requirement: Banner CRUD

The system SHALL provide banner/carousel management with image, link, and position support.

#### Scenario: Create a banner
- **WHEN** admin sends POST `/api/v1/cms/banners` with `title`, `image`, `link`, `position`
- **THEN** system creates the banner and returns the created record

#### Scenario: List banners
- **WHEN** admin sends GET `/api/v1/cms/banners` with optional `position` filter
- **THEN** system returns banners sorted by `sort` ASC, optionally filtered by position

#### Scenario: Update banner
- **WHEN** admin sends PUT `/api/v1/cms/banners/:id` with updated fields
- **THEN** system updates the banner and returns the updated record

#### Scenario: Delete banner
- **WHEN** admin sends DELETE `/api/v1/cms/banners/:id`
- **THEN** system deletes the banner permanently

### Requirement: Banner status and position

Banners SHALL have a status field and a position identifier for placement on different page areas.

#### Scenario: Filter by position
- **WHEN** public API requests banners with `position=home_top`
- **THEN** system returns only enabled banners with that position, sorted by `sort` ASC

#### Scenario: Disable a banner
- **WHEN** admin updates a banner's status to disabled
- **THEN** the banner no longer appears in public API responses
