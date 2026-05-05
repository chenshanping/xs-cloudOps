---
description: Standardize CRUD work for the go-base workspace. Use when creating/refactoring admin modules, DB tables, Go backend APIs, Vue frontend pages, or Excel import/export features.
---

# Backend CRUD Frontend

Discover the live module shape first, then extend it.

## Core Rule

- Do NOT assume Spring Boot paths or XTMS conventions
- Do NOT assume each module has its own directory — many are flat files
- Prefer the current neighboring module over memory

## Read First

Inspect these anchors before editing:
- `server/api/v1/user.go`, `server/service/user.go`, `server/model/sys_user.go`
- `server/router/modules/user.go`, `server/model/request/request.go`
- `web/src/api/user.ts`, `web/src/views/admin/system/user/index.vue`
- `web/src/components/ProTable.vue`, `web/src/utils/permission.ts`
- `web/src/types/index.ts`

## Backend Conventions

- API handlers in `server/api/v1/<module>.go`
- Services in `server/service/<module>.go`
- Models in `server/model/<module>.go`
- Routes in `server/router/modules/<module>.go`
- Response helpers: `response.BadRequest`, `response.Fail`, `response.OkWithData`, `response.OkWithPage`
- No database foreign keys

## Frontend Conventions

- Use `ProTable` for standard search + table pages
- Use `useTableColumns(...)` for permission-aware action columns
- Use `v-permission` on buttons, filter dropdown items in JS
- Default create/edit flows to Drawer (not Modal)
- Reuse shared components: `AvatarUpload`, `ImageUpload`, `FileUpload`, `FilePreview`
- Support dark mode: use `useUiStore().isDark` and semantic CSS variables

## Verification

// turbo
- Backend: `go build ./...` (in `server/`)
// turbo
- Frontend: `npm run build` (in `web/`)

## Excel Import / Export / Template

When a module needs Excel import, export, or template download, follow these rules.

### Reference Files

- Generic framework: `server/utils/excel.go` (ExcelExporter, ExcelImporter, ValidateImport, ImportField)
- Module example: `server/service/user/user_excel.go`
- API handlers: `server/api/v1/user.go` (GetUserImportTemplate, ImportUsers, ExportUsers)
- Frontend: `web/src/api/user.ts`, `web/src/views/admin/system/user/index.vue`
- Result modal: `web/src/views/admin/system/user/components/ImportResultModal.vue`

### Backend: Three Functions Per Module

Place in `server/service/<module>/<module>_excel.go`:

```go
func (s *XxxService) GetImportTemplate(contextID uint) ([]byte, string, error)
func (s *XxxService) ImportXxx(operatorID uint, contextID uint, fileData []byte) (*ImportXxxResult, error)
func (s *XxxService) ExportXxx(operatorID uint, contextID uint, ids []uint) ([]byte, string, error)
```

### Backend: Field Definition

```go
var xxxImportHeaders = []string{"列A", "列B"}

func xxxImportFields() []utils.ImportField {
    return []utils.ImportField{
        {Header: "列A", Key: "field_a", Required: true, Type: "string", MaxLen: 50, Validate: ...},
        {Header: "列B", Key: "field_b", Type: "string", Enum: []string{"选项1", "选项2"}},
    }
}
```

ImportField capabilities: Header, Key, Required, Type (string/int/uint/float64/time.Time), MaxLen, Enum, Validate func.
Validation order per cell: Required → MaxLen → Enum → Type parse → Custom Validate.

### Backend: Key Rules

- **Template**: Sheet name and filename include context name (e.g. `{部门名}_用户导入模板`). Include one example row. Add `AddDataValidation` dropdowns for enum columns.
- **Import**: Require context ID (e.g. dept_id) as param, NOT as an Excel column. Use `utils.ValidateImport` for generic validation, then batch uniqueness check + intra-file dedup. Set defaults (status=enabled, role, avatar from system config). Create records one by one, collect errors, continue on failure. File size limit 10MB in handler.
- **Export**: Require context ID + optional `ids []uint` for selective export. Query `WHERE context_id = ? [AND id IN ?]`. Filename: `{部门名}_用户导出.xlsx`. Same headers as import template.
- **Dict mapping**: Keep label↔value functions in the `_excel.go` file. Align with data dict values (e.g. 男=0, 女=1).

### Backend: API Handlers

```
GET  /<resources>/import-template?dept_id=xxx       → returns Excel blob
POST /<resources>/import         (form: file + dept_id) → returns ImportResult JSON
GET  /<resources>/export?dept_id=xxx&ids=1,2,3      → returns Excel blob
```

- Template handler: parse `dept_id` from query, call service, set `Content-Disposition: attachment; filename=<filename>`
- Import handler: parse `dept_id` from form value, validate non-zero, read file ≤10MB, call service
- Export handler: parse `dept_id` from query (required), parse `ids` from query (optional comma-separated), call service

### Backend: Permissions

Register in `server/initialize/db_tables.go`:
- Button permissions: `system:<module>:import`, `system:<module>:export`
- Use `grantMenuToRolesWithPermission` to inherit from parent list permission
- Use `ensureApiAccessInheritedFrom` for API Casbin policies

## Menu & API Bootstrap (MANDATORY for new modules)

Every new module MUST register its menus and APIs in `server/initialize/db_tables.go`. Two separate paths:

### 1. Fresh Install (`initDefaultData`)

Add root menu + sub-menus + all API entries to the existing block:

```go
// 根菜单
xxxMgmt := model.SysMenu{ParentID: 0, Name: "XXX管理", Path: "/xxx", Component: "Layout", Icon: "...", Sort: N, Type: 1, Permission: "xxx", Status: 1}
global.DB.Create(&xxxMgmt)

xxxMenus := []model.SysMenu{
    {ParentID: xxxMgmt.ID, Name: "子菜单", Path: "/xxx/sub", Component: "xxx/sub/index", Icon: "...", Sort: 1, Type: 2, Permission: "xxx:sub:list", Status: 1},
}
global.DB.Create(&xxxMenus)
// Include in allMenus and adminRole association
```

Add all module APIs to the `apis` slice so they get created and assigned to the admin role + Casbin policies.

### 2. Existing Install Upgrade (`ensureBuiltInData`)

Create an `ensureXxxMenus()` function following the `ensureAIToolMenus` pattern:

- Root menu: `FirstOrCreate` by `permission`
- Sub-menus: `FirstOrCreate` by `permission`, using `Attrs` for create-only fields
- Call `grantMenusToRoleCodes(menuIDs, []string{"admin", "system_admin"})`
- API entries: loop + `ensureApiAccessForRoleCodes(api, []string{"admin", "system_admin"})`
- Call the function from `ensureBuiltInData()`

### Key Rules

- **Never skip either path** — fresh installs use `initDefaultData`, upgrades use `ensureBuiltInData`
- Use `FirstOrCreate + Attrs` pattern — only create if missing, never overwrite existing customizations
- Menu `Permission` field is the unique key for dedup
- API `Path + Method` pair is the unique key for dedup
- Always grant to `admin` + `system_admin` role codes

### Frontend: API Functions

```ts
export function downloadXxxImportTemplate(contextId?: number) {
  return request.get('/xxx/import-template', { params: { dept_id: contextId }, responseType: 'blob' })
}
export function importXxx(file: File, contextId: number) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('dept_id', String(contextId))
  return request.post<any, ApiResponse<ImportResult>>('/xxx/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
export function exportXxx(contextId: number, ids?: number[]) {
  const params: Record<string, any> = { dept_id: contextId }
  if (ids && ids.length > 0) params.ids = ids.join(',')
  return request.get('/xxx/export', { params, responseType: 'blob' })
}
```

### Frontend: Toolbar Buttons

- **导入按钮**: `a-upload` with `:disabled="!hasContextSelected"`, wrapped in `a-tooltip` showing "请先选择部门"
- **下载模板**: `a-dropdown` with menu item, pass contextId to template API, filename = `{部门名}_导入模板.xlsx`
- **导出按钮**: `a-dropdown` with two items: 导出全部 / 导出选中(N), disabled when no context selected. "导出选中" disabled when `selectedRowKeys.length === 0`
- All buttons use `v-permission`

### Frontend: Import Result

Reuse `ImportResultModal.vue` — props: `open: boolean`, `result: ImportResult | null`. Shows success/warning/error status + error detail table (row, column, value, message).

### New Module Checklist

1. [ ] Models in `server/model/<module>.go`
2. [ ] Request types in `server/model/request/<module>_request.go`
3. [ ] Services in `server/service/<module>/`
4. [ ] API handlers in `server/api/v1/<module>.go`
5. [ ] Routes in `server/router/modules/<module>.go`
6. [ ] Register service singleton in `server/service/service.go`
7. [ ] AutoMigrate models in `server/initialize/db_tables.go` → `InitDBTables()`
8. [ ] **Menu bootstrap (fresh)**: add menus to `initDefaultData()`, include in admin role association
9. [ ] **Menu bootstrap (upgrade)**: add `ensureXxxMenus()` to `ensureBuiltInData()`
10. [ ] **API bootstrap (fresh)**: add API entries to `initDefaultData()` apis slice
11. [ ] **API bootstrap (upgrade)**: add `ensureApiAccessForRoleCodes` calls in `ensureXxxMenus()`
12. [ ] Frontend API in `web/src/api/<module>.ts`
13. [ ] Frontend view resolver: add prefix to `ADMIN_VIEW_PREFIXES` in `router/view-resolver.ts` if needed
14. [ ] Frontend pages in `web/src/views/admin/<module>/`
15. [ ] (Optional) Excel import/export: `server/service/<module>/<module>_excel.go`
16. [ ] Verify `go build ./...` + `npm run build`

## Self-Check

- Did I inspect a real neighboring module first?
- Did I follow current flat-file backend structure?
- Did I register menus in BOTH `initDefaultData` (fresh) AND `ensureBuiltInData` (upgrade)?
- Did I register APIs in BOTH `initDefaultData` (fresh) AND `ensureXxxMenus` (upgrade)?
- Did I reuse ProTable, permission helpers, shared components?
- Did I keep interactions usable (no placeholder-only buttons)?
- Did I verify dark mode for touched surfaces?
