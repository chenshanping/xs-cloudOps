## ADDED Requirements

### Requirement: Role permission assignment MUST support page-grouped same-screen selection
The system MUST let administrators configure menu permissions and API permissions inside one same-screen role permission drawer without requiring primary tab switching between menu and API areas.

#### Scenario: Drawer shows page-grouped permission blocks
- **WHEN** an administrator opens role permission assignment for a role
- **THEN** the system MUST show a same-screen permission editor with left-side menu navigation and right-side page-grouped permission blocks

#### Scenario: One page block contains both menu and API areas
- **WHEN** the administrator views a page permission block
- **THEN** the system MUST expose a menu or button permission area and a separate related API permission area within that same block

### Requirement: Same-screen grouping MUST NOT change dual-track authorization semantics
The system MUST preserve separate menu and API selection, saving, and runtime effect even when both are displayed together.

#### Scenario: Selecting a menu does not automatically select APIs
- **WHEN** an administrator selects one or more menus or buttons in the grouped drawer
- **THEN** the system MUST NOT automatically select related APIs

#### Scenario: Selecting an API does not automatically select menus
- **WHEN** an administrator selects one or more APIs in the grouped drawer
- **THEN** the system MUST NOT automatically select related menus or buttons

### Requirement: Drawer MUST support a single save action with split result handling
The system MUST provide one save action for the grouped drawer while continuing to persist menu permissions and API permissions through their existing separate write paths.

#### Scenario: Save succeeds on both tracks
- **WHEN** menu and API permission updates both succeed
- **THEN** the system MUST report overall success and close the drawer

#### Scenario: Save partially fails
- **WHEN** exactly one of the menu or API permission updates fails
- **THEN** the system MUST keep the drawer open and indicate which side failed

### Requirement: Grouped drawer MUST preserve role truth on reopen
The system MUST rebuild the grouped view from the current saved role menus and role APIs each time the permission drawer is reopened.

#### Scenario: Reopen shows saved menu selections
- **WHEN** an administrator reopens the role permission drawer after menu permissions were previously saved
- **THEN** the grouped view MUST show the saved menu or button selections

#### Scenario: Reopen shows saved API selections
- **WHEN** an administrator reopens the role permission drawer after API permissions were previously saved
- **THEN** the grouped view MUST show the saved API selections

### Requirement: Unmatched APIs MUST remain assignable
The system MUST keep APIs assignable even when the frontend cannot confidently group them under a specific page block.

#### Scenario: API cannot be matched to a page
- **WHEN** an API cannot be stably grouped to a page menu using the frontend grouping rules
- **THEN** the system MUST place that API in an uncategorized area that still allows direct selection
