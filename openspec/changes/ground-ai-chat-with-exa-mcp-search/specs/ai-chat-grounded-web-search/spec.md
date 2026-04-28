## ADDED Requirements

### Requirement: Search-enabled chat SHALL ground answers with Exa web search
When a user enables联网搜索 in AI chat, the system MUST perform Exa-backed web search before generating the assistant answer and MUST use the returned search evidence as the primary basis for answering time-sensitive questions.

#### Scenario: Search is executed before answering
- **GIVEN** a user sends an AI chat message with `enable_search=true`
- **WHEN** the backend prepares the model request
- **THEN** it MUST call Exa MCP search before calling the configured LLM
- **AND** it MUST inject the search evidence into the model context for that turn

#### Scenario: Search-disabled chat keeps existing behavior
- **GIVEN** a user sends an AI chat message with `enable_search=false`
- **WHEN** the backend prepares the model request
- **THEN** it MUST NOT call Exa MCP
- **AND** it MUST preserve the existing non-search chat path

### Requirement: Search-grounded chat SHALL not guess when evidence is insufficient
When search is enabled and the system cannot confirm an answer from the returned evidence, the assistant MUST explicitly state that the answer is not yet confirmed instead of fabricating or completing the answer from stale knowledge.

#### Scenario: No authoritative evidence is available
- **GIVEN** a user asks about a time-sensitive fact with `enable_search=true`
- **AND** the search results do not provide a trustworthy or sufficient answer
- **WHEN** the assistant generates the reply
- **THEN** the reply MUST clearly state that the answer is not confirmed
- **AND** the reply MUST NOT invent names, schedules, scores, or other facts

#### Scenario: Search evidence conflicts with older model knowledge
- **GIVEN** a user asks about a time-sensitive fact with `enable_search=true`
- **AND** the search results conflict with what the underlying model may have learned previously
- **WHEN** the assistant generates the reply
- **THEN** the reply MUST prioritize the search evidence over stale model knowledge

### Requirement: Search-grounded chat SHALL expose visible sources in the assistant reply
When search is enabled and Exa returns usable results, the assistant reply MUST include visible source links so the user can inspect the evidence used for the answer.

#### Scenario: Successful grounded answer shows sources
- **GIVEN** a user sends an AI chat message with `enable_search=true`
- **AND** Exa returns usable search results
- **WHEN** the assistant reply is returned to the frontend
- **THEN** the reply MUST include a visible source section with one or more clickable links

#### Scenario: Saved conversation retains source visibility
- **GIVEN** a search-grounded assistant reply has been saved into the conversation history
- **WHEN** the user reopens that conversation later
- **THEN** the source section MUST still be visible as part of the saved reply content
