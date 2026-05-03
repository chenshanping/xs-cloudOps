## ADDED Requirements

### Requirement: Article CRUD

The system SHALL provide article management with rich text content, category association, and cover image support.

#### Scenario: Create article as draft
- **WHEN** admin sends POST `/api/v1/cms/articles` with `title`, `content`, `category_id`
- **THEN** system creates the article with `status=0` (draft) and returns the created record

#### Scenario: Update article
- **WHEN** admin sends PUT `/api/v1/cms/articles/:id` with updated fields
- **THEN** system updates the article and returns the updated record

#### Scenario: Get article list with pagination
- **WHEN** admin sends GET `/api/v1/cms/articles` with `page`, `page_size`, optional `category_id`, `status`, `keyword`
- **THEN** system returns paginated article list with total count, sorted by `is_top` DESC then `created_at` DESC

#### Scenario: Get article detail
- **WHEN** admin sends GET `/api/v1/cms/articles/:id`
- **THEN** system returns full article detail including content

#### Scenario: Delete article
- **WHEN** admin sends DELETE `/api/v1/cms/articles/:id`
- **THEN** system deletes the article permanently

### Requirement: Article status workflow

Articles SHALL have three statuses: draft (0), published (1), and unpublished/archived (2). Status transitions are explicit.

#### Scenario: Publish a draft article
- **WHEN** admin sends PUT `/api/v1/cms/articles/:id/status` with `status=1`
- **THEN** system sets status to published and sets `published_at` to current time (if not already set)

#### Scenario: Unpublish a published article
- **WHEN** admin sends PUT `/api/v1/cms/articles/:id/status` with `status=2`
- **THEN** system sets status to unpublished, article no longer appears in public API

#### Scenario: Re-publish an unpublished article
- **WHEN** admin sends PUT `/api/v1/cms/articles/:id/status` with `status=1`
- **THEN** system sets status back to published, `published_at` retains its original value

### Requirement: Article top/pin feature

Articles SHALL support a top/pin flag that makes them appear first in list queries.

#### Scenario: Pin an article
- **WHEN** admin updates article with `is_top=1`
- **THEN** the article appears before non-pinned articles in list queries regardless of creation time

#### Scenario: Unpin an article
- **WHEN** admin updates article with `is_top=0`
- **THEN** the article follows normal sorting order

### Requirement: Article slug auto-generation

Articles SHALL have a unique `slug` for URL-friendly identification.

#### Scenario: Create article without slug
- **WHEN** admin creates an article without providing `slug`
- **THEN** system auto-generates slug from the title (pinyin or transliteration)

#### Scenario: Create article with custom slug
- **WHEN** admin creates an article with a specific `slug`
- **THEN** system uses the provided slug if it is unique

#### Scenario: Duplicate slug
- **WHEN** admin creates or updates an article with a slug that already exists
- **THEN** system returns a 400 error indicating slug duplication
