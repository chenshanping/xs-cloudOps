## ADDED Requirements

### Requirement: Admin can manage AI providers from a master-detail admin page
The system MUST allow an authorized administrator to manage AI providers from an admin page that shows a provider list and a model-focused detail panel.

#### Scenario: View provider list and current provider models
- **GIVEN** the administrator opens the AI config page
- **WHEN** AI config data is loaded from the saved `ai_config`
- **THEN** the page MUST show a left-side provider list
- **AND** the page MUST show the selected provider's imported model list in the main detail panel

#### Scenario: Switch provider
- **GIVEN** the AI config contains multiple providers
- **WHEN** the administrator selects a different provider from the left-side list
- **THEN** the page MUST switch the detail panel to that provider without leaving the page

### Requirement: Admin can edit provider fields in a dedicated drawer
The system MUST allow an authorized administrator to create or edit provider fields through a dedicated provider edit drawer rather than inline editing in the main panel.

#### Scenario: Create provider from drawer
- **GIVEN** the administrator is on the AI config page
- **WHEN** the administrator opens the provider edit drawer and submits a new provider with required fields
- **THEN** the page MUST add that provider to the current editing state
- **AND** the new provider MUST NOT be persisted until the administrator saves the overall AI config

#### Scenario: Edit provider from drawer
- **GIVEN** the administrator selects an existing provider
- **WHEN** the administrator opens the provider edit drawer and updates provider fields
- **THEN** the page MUST update the provider in the current editing state
- **AND** the update MUST NOT be persisted until the administrator saves the overall AI config

### Requirement: Admin can fetch provider models through a backend proxy
The system MUST allow an authorized administrator to fetch a remote provider model list through a backend proxy endpoint using the current provider edit values.

#### Scenario: Fetch models with current provider values
- **GIVEN** the administrator has entered a provider `base_url` and `api_key` in the current editing state
- **WHEN** the administrator opens provider model management and requests the remote model list
- **THEN** the frontend MUST send those current values to a backend endpoint
- **AND** the backend MUST proxy the request to the remote OpenAI-compatible model-list endpoint

#### Scenario: Reject missing provider credentials
- **GIVEN** the current provider editing state does not contain the required `base_url` or `api_key`
- **WHEN** the administrator requests the remote model list
- **THEN** the system MUST reject the request
- **AND** the administrator MUST receive a readable validation error

### Requirement: Provider model fetch must support OpenAI-compatible base URL variants
The backend MUST support OpenAI-compatible base URL variants when discovering remote models.

#### Scenario: Base URL already ends with v1
- **WHEN** the fetch request contains a `base_url` ending in `/v1`
- **THEN** the backend MUST request `<base_url>/models`

#### Scenario: Base URL does not end with v1
- **WHEN** the fetch request contains a `base_url` that does not end in `/v1`
- **THEN** the backend MUST first try `<base_url>/v1/models`
- **AND** if that attempt fails, the backend MUST retry with `<base_url>/models`

### Requirement: Admin can import selected provider models without immediate persistence
The system MUST allow an authorized administrator to select remote provider models and import them into the current provider editing state without immediately writing to the database.

#### Scenario: Import selected models
- **GIVEN** the provider model management drawer shows a fetched remote model list
- **WHEN** the administrator selects one or more remote models and confirms import
- **THEN** the page MUST merge the selected models into the current provider's imported model list
- **AND** the merged result MUST remain in the current editing state until the administrator saves the overall AI config

#### Scenario: Imported models are not persisted before save
- **GIVEN** the administrator imported remote models but did not save the AI config
- **WHEN** the administrator refreshes the page or abandons the edit state
- **THEN** the unsaved imported models MUST NOT appear in persisted config

### Requirement: Model import must append by unique model id and preserve local edits
The system MUST treat model ID as the unique import key and MUST NOT overwrite locally edited display fields for models that already exist.

#### Scenario: Skip duplicate model ids on import
- **GIVEN** the current provider already contains a model with a given `id`
- **WHEN** the administrator imports a remote model with the same `id`
- **THEN** the system MUST NOT append a duplicate entry

#### Scenario: Preserve local display fields for existing model
- **GIVEN** the current provider already contains a model with a customized local `name` or `description`
- **WHEN** the administrator imports a remote model with the same `id`
- **THEN** the existing local `name` and `description` MUST remain unchanged

### Requirement: Provider model discovery must protect provider secrets
The system MUST protect provider secrets during remote model discovery.

#### Scenario: Backend proxy hides provider key from browser targets
- **WHEN** the administrator requests the remote model list
- **THEN** the browser MUST call only the local backend endpoint
- **AND** the browser MUST NOT directly call the third-party provider endpoint

#### Scenario: Error handling does not expose raw api key
- **GIVEN** the remote model discovery fails
- **WHEN** the backend records the failure or returns an error to the client
- **THEN** the system MUST NOT expose the raw `api_key` in logs or response messages
