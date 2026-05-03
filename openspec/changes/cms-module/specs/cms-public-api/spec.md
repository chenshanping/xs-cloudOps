## ADDED Requirements

### Requirement: Public category API

The system SHALL provide unauthenticated read-only access to enabled CMS categories.

#### Scenario: Get public category tree
- **WHEN** any client sends GET `/api/v1/public/cms/categories`
- **THEN** system returns only enabled categories as a tree, without requiring authentication

#### Scenario: Get articles by category slug
- **WHEN** any client sends GET `/api/v1/public/cms/categories/:slug` with `page`, `page_size`
- **THEN** system returns paginated published articles under that category (including child categories)
- **THEN** only articles with `status=1` (published) are returned

### Requirement: Public article API

The system SHALL provide unauthenticated read-only access to published CMS articles.

#### Scenario: List published articles
- **WHEN** any client sends GET `/api/v1/public/cms/articles` with `page`, `page_size`, optional `category_id`
- **THEN** system returns paginated published articles sorted by `is_top` DESC then `published_at` DESC
- **THEN** only articles with `status=1` are returned

#### Scenario: Get article detail by slug
- **WHEN** any client sends GET `/api/v1/public/cms/articles/:slug`
- **THEN** system returns the full article detail if published
- **THEN** system increments the article's `views` count by 1

#### Scenario: Access unpublished article
- **WHEN** any client sends GET `/api/v1/public/cms/articles/:slug` for a draft or unpublished article
- **THEN** system returns a 404 error

### Requirement: Public page API

The system SHALL provide unauthenticated read-only access to enabled CMS pages.

#### Scenario: Get page by slug
- **WHEN** any client sends GET `/api/v1/public/cms/pages/:slug`
- **THEN** system returns the page content if enabled

#### Scenario: Access disabled page
- **WHEN** any client sends GET `/api/v1/public/cms/pages/:slug` for a disabled page
- **THEN** system returns a 404 error

### Requirement: Public banner API

The system SHALL provide unauthenticated read-only access to enabled banners.

#### Scenario: Get banners by position
- **WHEN** any client sends GET `/api/v1/public/cms/banners` with optional `position` parameter
- **THEN** system returns enabled banners filtered by position, sorted by `sort` ASC

### Requirement: Public navigation API

The system SHALL provide unauthenticated read-only access to the public navigation tree.

#### Scenario: Get navigation tree
- **WHEN** any client sends GET `/api/v1/public/cms/navigations`
- **THEN** system returns only enabled navigation items as a tree structure
