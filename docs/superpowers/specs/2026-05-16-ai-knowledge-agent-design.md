# AI Knowledge Agent Design

## Context

This design defines an engineering-oriented knowledge agent architecture for `go-base`.

The goal is not to add another generic chat box. The goal is to build a maintainable private-knowledge agent system with:

- `go-base` as the primary business system
- `Eino` as the agent/workflow runtime
- a standard `MCP Server` for Claude Desktop / Cursor connectivity
- a document-processing pipeline for `PDF` / `Markdown` / `Excel`
- hybrid retrieval with `Qdrant + Elasticsearch + handwritten RRF`
- explicit task orchestration, auditability, and permission boundaries

This design reuses:

- `go-base` for user, permission, menu, config, and audit capabilities
- the existing AI provider/config foundation already present under `server/service/ai`

This design does **not**:

- migrate the project into a standalone AI platform
- expose direct database / Qdrant / ES access to external MCP clients
- copy the Python/LangGraph/MCP implementation style from `AI-CloudOps-aiops`

## Confirmed Decisions

- The main application remains `go-base`.
- The selected architecture is `Go orchestrator + Python document worker + standalone standard MCP server`.
- The `MCP Server` runs as an independent process, not embedded into the main `go-base` process.
- The main target capabilities are `Eino + standard MCP + document pipeline + hybrid retrieval + DAG`.
- `AI-CloudOps-aiops` is used only as an engineering-structure reference.
- `basjoo` is used only as a provider/task decomposition reference.

## Recommended Architecture

### 1. go-base Main Application

Responsibilities:

- knowledge-base business management
- task lifecycle management
- authorization and RBAC
- menu / API exposure for admin UI
- retrieval orchestration
- Eino graph execution
- audit logging

The main app remains the system of record.

### 2. Python Document Worker

Responsibilities:

- PDF parsing
- Markdown parsing
- Excel parsing
- structure normalization
- fixed-size chunking
- semantic chunking

The worker does not own business state. It receives parsing jobs and returns normalized output.

### 3. Retrieval Infrastructure

Responsibilities:

- `Qdrant` for dense vector retrieval
- `Elasticsearch` for sparse BM25 retrieval
- RRF merge in Go business layer
- optional rerank in Go business layer

Retrieval infrastructure is storage-level capability, not business orchestration.

### 4. Standard MCP Server

Responsibilities:

- expose standard MCP tools to Claude Desktop / Cursor
- call controlled knowledge services from `go-base`
- enforce access control and audit boundaries

The MCP layer is protocol adaptation, not business ownership.

## Directory Layout

### go-base server

- `server/api/v1/knowledge.go`
- `server/router/modules/knowledge.go`
- `server/model/knowledge/`
- `server/model/request/knowledge_request.go`
- `server/model/response/knowledge_response.go`
- `server/service/knowledge/`
- `server/service/retrieval/`
- `server/service/agent/`
- `server/service/docworker/`
- `server/service/searchindex/`
- `server/service/mcpaudit/`

### go-base web

- `web/src/views/admin/ai-knowledge/`
- `web/src/views/admin/ai-knowledge/components/`
- `web/src/api/knowledge.ts`
- `web/src/types/knowledge.ts`

### standalone document worker

- `workers/doc-parser/app/parsers/`
- `workers/doc-parser/app/chunkers/`
- `workers/doc-parser/app/pipeline/`
- `workers/doc-parser/app/main.py`

### standalone MCP server

- `mcp-server/cmd/server/main.go`
- `mcp-server/internal/server/`
- `mcp-server/internal/tools/`
- `mcp-server/internal/auth/`
- `mcp-server/internal/audit/`

## Data Model

### ai_knowledge_base

Fields:

- `id`
- `name`
- `code`
- `description`
- `status`
- `scope_type`
- `scope_id`
- `embedding_provider`
- `embedding_model`
- `rerank_provider`
- `rerank_model`
- `created_by`
- `updated_by`

Purpose:

- logical boundary for a knowledge corpus
- future separation by tenant, team, or module

Soft-delete rule:

- if `code` is unique, use physical delete strategy to avoid the repo's soft-delete + unique-index trap

### ai_knowledge_document

Fields:

- `id`
- `knowledge_base_id`
- `title`
- `source_type`
- `file_id`
- `mime_type`
- `content_hash`
- `doc_version`
- `parse_strategy_version`
- `chunk_strategy_version`
- `status`
- `error_message`
- `page_count`
- `token_count`
- `created_by`
- `updated_by`

Purpose:

- document source of record
- explicit versioning for re-import and rebuild control

### ai_knowledge_chunk

Fields:

- `id`
- `document_id`
- `knowledge_base_id`
- `chunk_no`
- `chunk_type`
- `heading_path`
- `page_no`
- `sheet_name`
- `text_content`
- `token_count`
- `char_count`
- `chunk_hash`
- `parent_chunk_id`
- `index_status`

Purpose:

- support both fixed and semantic chunks
- preserve PDF page information and Excel sheet information

### ai_knowledge_index_ref

Fields:

- `id`
- `chunk_id`
- `document_id`
- `knowledge_base_id`
- `vector_store`
- `vector_id`
- `search_store`
- `search_doc_id`
- `embedding_model`
- `index_version`
- `status`

Purpose:

- persistent mapping for Qdrant / Elasticsearch
- safe deletion, rebuild, and replay

### ai_knowledge_task

Fields:

- `id`
- `knowledge_base_id`
- `document_id`
- `task_type`
- `status`
- `payload_json`
- `result_json`
- `retry_count`
- `started_at`
- `finished_at`
- `error_message`
- `triggered_by`

Purpose:

- ingestion / rebuild / delete task tracking
- graph node execution visibility

### ai_mcp_tool_audit

Fields:

- `id`
- `tool_name`
- `client_type`
- `client_id`
- `user_id`
- `knowledge_base_id`
- `request_summary`
- `response_status`
- `latency_ms`
- `error_message`
- `created_at`

Purpose:

- audit MCP invocations from Claude Desktop / Cursor / internal clients
- record summaries and status, not raw sensitive payloads by default

## State Transitions

### Document status

- `uploaded`
- `queued`
- `parsing`
- `parsed`
- `chunking`
- `chunked`
- `indexing`
- `ready`
- `failed`
- `archived`

### Task status

- `pending`
- `running`
- `success`
- `failed`
- `cancelled`

### Chunk index status

- `pending`
- `vector_ready`
- `search_ready`
- `ready`
- `failed`

## Processing Pipeline

### Ingestion pipeline

Logical sequence:

- upload
- normalize
- parse
- clean
- fixed chunking
- semantic chunking
- embed
- index to Qdrant
- index to Elasticsearch
- finalize

Rules:

- PDF keeps page numbers
- Markdown keeps heading/code/table structure
- Excel keeps sheet/table/range context
- chunking output is versioned
- content hash and strategy versions are part of idempotency

## Eino Graph Design

### IngestionGraph

Nodes:

1. `validate_document`
2. `deduplicate_document`
3. `dispatch_parse_task`
4. `collect_parse_result`
5. `chunk_document`
6. `persist_chunks`
7. `build_embeddings`
8. `index_vector_store`
9. `index_search_store`
10. `finalize_document`

Execution boundary:

- the Python worker handles parsing and normalization-heavy work
- Go remains the owner of chunk persistence, embedding orchestration, and index state

### QueryGraph

Nodes:

1. `normalize_query`
2. `classify_query`
3. `rewrite_query`
4. `retrieve_dense`
5. `retrieve_sparse`
6. `merge_rrf`
7. `rerank_candidates`
8. `build_answer_context`
9. `generate_answer`
10. `post_process_answer`

Execution boundary:

- this graph answers knowledge questions only
- it does not perform direct write actions

### ToolGraph

Nodes:

1. `classify_intent`
2. `decide_tool_or_rag`
3. `invoke_tool`
4. `merge_tool_result`
5. `generate_tool_answer`

Execution boundary:

- ToolGraph never directly accesses DB/Qdrant/ES
- tool invocation must go through service and registry boundaries

## Retrieval Design

### Dense retrieval

- storage: `Qdrant`
- output: topK dense candidates

### Sparse retrieval

- storage: `Elasticsearch`
- output: topK BM25 candidates

### Merge

- handwritten `RRF`
- de-duplication
- same-document grouping
- baseline business filtering

### Rerank

- optional in phase 1
- keep interface from the start

### Answer generation

- answer must include citations
- answer must explicitly state uncertainty when evidence is insufficient

## MCP Tool Design

### Phase 1 tools

1. `search_knowledge`
- input: `knowledge_base`, `query`, `top_k`
- output: chunk summaries, document title, page/sheet metadata, score

2. `get_document_chunk`
- input: `chunk_id`
- output: full chunk content and metadata

3. `get_document_outline`
- input: `document_id`
- output: heading tree, page numbers, sheet information

4. `list_knowledge_bases`
- input: optional scope filters
- output: visible knowledge bases for the caller

5. `get_ingestion_task_status`
- input: `task_id`
- output: ingestion / indexing task state

### Phase 2 tools

- `search_runbook`
- `search_config_policy`
- `query_cmdb_resource`

### MCP constraints

- tool results must be structured JSON
- all tools require permission checks
- all tools require audit logging
- no direct SQL / direct ES / direct Qdrant access from MCP clients

## Phase 1 Scope

### Included

1. knowledge-base management models and APIs
2. document upload and task creation
3. Python worker for PDF / Markdown / Excel parsing
4. fixed + semantic chunking
5. Qdrant + Elasticsearch dual indexing
6. QueryGraph for knowledge answering
7. standard MCP server with first-batch knowledge tools
8. Claude Desktop / Cursor connectivity verification

### Explicitly excluded

1. multi-agent collaboration
2. generalized ops auto-execution platform
3. unrestricted tool execution
4. direct MCP-driven write operations against business systems
5. system-wide refactor of the existing AI module

## Acceptance Criteria

### Functional

1. an admin can create a knowledge base
2. an admin can upload PDF / Markdown / Excel documents
3. the system can parse and chunk these documents successfully
4. chunks are written to MySQL and indexed into both Qdrant and Elasticsearch
5. QueryGraph returns answers with citations
6. Claude Desktop / Cursor can call MCP tools against the private knowledge base

### Engineering

1. ingestion tasks are observable and retryable
2. index rebuild is idempotent
3. permission checks are enforced for MCP tools
4. MCP tool invocations are audited
5. failures are diagnosable by task/document/index status

### Non-goals for phase 1

1. no promise of full autonomous workflow execution
2. no business-write tools exposed through MCP
3. no complex version-history UI beyond explicit document version fields

## Risks

1. pure-Go document parsing would slow delivery and reduce format quality
2. protocol concerns and business concerns can become entangled if MCP is embedded into the main app
3. hybrid retrieval quality will drift if chunk/version/index versioning is not enforced
4. delete/rebuild consistency across MySQL/Qdrant/ES requires explicit idempotent design

## Recommendation

Proceed with phase 1 using:

- `go-base` as the orchestrator and business owner
- Python worker for document parsing and chunk preparation
- `Eino` for graph execution
- `Qdrant + Elasticsearch + handwritten RRF`
- standalone standard MCP server for external AI clients

This is the smallest architecture that is still engineering-grade, extensible, and aligned with the current repository constraints.
