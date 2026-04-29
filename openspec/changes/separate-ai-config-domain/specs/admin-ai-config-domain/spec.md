## ADDED Requirements

### Requirement: Admin AI config MUST be exposed through AI module APIs
The system MUST expose backend AI configuration management through AI module APIs rather than through generic system-config APIs.

#### Scenario: Load AI config from AI module endpoint
- **WHEN** an authorized administrator opens the AI config page
- **THEN** the frontend MUST load AI configuration from an AI module endpoint
- **AND** the page MUST NOT require a generic system-config read API as its primary load dependency

#### Scenario: Save AI config through AI module endpoint
- **WHEN** an authorized administrator saves AI configuration changes
- **THEN** the frontend MUST persist the configuration through an AI module endpoint
- **AND** the page MUST NOT require a generic system-config write API as its primary save dependency

### Requirement: AI config page permissions MUST align with AI module ownership
The system MUST allow the AI config page and its related management actions to be authorized within AI module permissions rather than depending on generic config-management API permissions.

#### Scenario: Grant AI config without config-management API dependency
- **WHEN** a role is granted the AI config page and its required AI management APIs
- **THEN** that role MUST be able to load and save the AI config page without separately granting generic config-management APIs for the same page workflow

#### Scenario: Keep AI model discovery and testing within AI module
- **WHEN** a role is authorized to manage AI configuration
- **THEN** AI model discovery and AI model testing APIs MUST remain grouped under AI module ownership
- **AND** those actions MUST NOT be presented as config-management APIs for the AI config workflow

### Requirement: Phase 1 MUST preserve existing ai_config persistence compatibility
The system MUST preserve compatibility with already-saved `ai_config` data during the first-phase domain realignment.

#### Scenario: Existing ai_config remains readable
- **GIVEN** AI configuration has already been saved in the existing `ai_config` persistence shape
- **WHEN** the first-phase AI module endpoint loads configuration
- **THEN** the saved provider and model data MUST remain readable without migration

#### Scenario: First-phase save remains backward-compatible
- **WHEN** an administrator saves AI configuration through the new AI module endpoint
- **THEN** the resulting persisted configuration MUST remain compatible with the existing `ai_config` structure
- **AND** no provider/model table migration MUST be required in this phase
