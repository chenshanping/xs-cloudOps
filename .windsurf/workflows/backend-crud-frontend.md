---
description: Standardize CRUD work for the go-base workspace. Use when creating/refactoring admin modules, DB tables, Go backend APIs, Vue frontend pages, or Excel import/export features.
---

# Backend CRUD Frontend

Discover the live module shape first, then extend it.

## Steps

Follow these steps in order. Do not skip any step.

1. **Discover** — read a neighboring module to learn the live shape. Inspect at least:
   - `server/router/modules/` (pick a similar module)
   - `server/api/v1/`, `server/service/`, `server/model/`
   - `web/src/views/admin/system/user/index.vue` (for v-permission pattern)
   - `server/initialize/db_tables.go` (for menu + button bootstrap pattern, search `ensureExamMenus`)

2. **Backend model + request** — create or update model in `server/model/<module>.go`, request types in `server/model/request/<module>_request.go`.
   - Business tables that reference uploads MUST store `xxx_file_id` / `xxx_file_ids`, not uploaded URL columns.

3. **Backend service** — create or update service in `server/service/<module>/`. Business logic in service, not handlers.
   - Service `Create` / `Update` paths that bind uploaded files MUST call `filesvc.Reference.ReplaceRefs(...)`.
   - Service `Delete` / `Clear` paths that remove business data MUST call `filesvc.Reference.ClearRefs(...)`.
   - Do NOT add module-specific scans/checks to `server/service/file/file.go`.

4. **Backend API handlers** — create or update in `server/api/v1/<module>.go`. Use `response.BadRequest`, `response.Fail`, `response.OkWithData`, `response.OkWithPage`.

5. **Backend routes** — register in `server/router/modules/<module>.go`. Register service singleton in `server/service/service.go` if new module.

6. **Backend verify** — run:
// turbo
   `go build ./...` (in `server/`)

7. **Menu & button permission bootstrap** — in `server/initialize/db_tables.go`:
   - Add `ensureXxxMenus()` to `ensureBuiltInData()` if not exists
   - Create root menu (Type=1) + page menus (Type=2) via `FirstOrCreate` by `permission`
   - **For EACH page menu, create Type=3 button permission menus** for every action: add, edit, delete, import, export, copy, batchDelete, batchEdit, etc.
   - Collect ALL menu IDs (root + pages + buttons) and call `grantMenusToRoleCodes(menuIDs, []string{"admin", "system_admin"})`
   - Register all API entries via `ensureApiAccessForRoleCodes`
   - Permission key format: `<module>:<resource>:<action>` (e.g. `exam:question:add`)

8. **Backend verify again** — run:
// turbo
   `go build ./...` (in `server/`)

9. **Frontend API** — create or update in `web/src/api/<module>.ts`.

10. **Frontend pages** — create or update in `web/src/views/admin/<module>/`. Follow the Frontend Design rules in the Reference section below. Use Drawer for create/edit, not Modal.
    - Admin pages are tool pages, not landing pages. Do not add marketing/hero sections, slogans, decorative chips, English promo copy, or large top banners unless explicitly requested.
    - Before introducing a page header, compare neighboring admin pages. If they use compact card/table/tabs layouts, keep the new page compact too.

11. **Add `v-permission` to EVERY action button** — this is mandatory, not optional:
    - Toolbar: add, import, export, batch edit, batch delete
    - Table action column: edit, copy, delete
    - Pattern: `v-permission="'module:resource:action'"`

12. **Frontend verify** — lightweight checks (do NOT run `npm run build` — see Token Budget Rule in `AGENTS.md`):
    - Read back edited files to confirm structure and imports.
    - `grep_search` for new API paths / component names to confirm no broken references.
    - `find_by_name` to confirm new files exist where expected.
    - Hand off to user with a click-through list for the running dev server.

13. **Self-check** — verify all items:
    - [ ] Menus registered in `ensureBuiltInData` (upgrade path)?
    - [ ] APIs registered via `ensureApiAccessForRoleCodes`?
    - [ ] **Type=3 button permissions created for EVERY action button?**
    - [ ] **`v-permission` added to EVERY action button in the frontend?**
    - [ ] **Permission keys match between backend Type=3 menus and frontend `v-permission`?**
    - [ ] No action button without `v-permission`?
    - [ ] `go build` passes?

## Rules

- Do NOT assume Spring Boot paths or XTMS conventions
- Do NOT assume each module has its own directory — many are flat files
- Prefer the current neighboring module over memory
- No database foreign keys
- Button permissions are NOT optional — every UI action must have a Type=3 menu entry + `v-permission`
- Stop and ask when blocked, don't guess

---

## Reference: Backend Conventions

- API handlers in `server/api/v1/<module>.go`
- Services in `server/service/<module>.go`
- Models in `server/model/<module>.go`
- Routes in `server/router/modules/<module>.go`
- Response helpers: `response.BadRequest`, `response.Fail`, `response.OkWithData`, `response.OkWithPage`

## Reference: Frontend Conventions

- Use `ProTable` for standard search + table pages
- Use `useTableColumns(...)` for permission-aware action columns
- Default create/edit flows to Drawer (not Modal)
- Reuse shared components: `AvatarUpload`, `ImageUpload`, `FileUpload`, `FilePreview`
- Support dark mode: use `useUiStore().isDark` and semantic CSS variables
- Upload components should keep preview URLs in local UI state; API payloads should submit only file IDs.

### v-permission Pattern

```vue
<!-- Toolbar -->
<a-button v-permission="'module:resource:add'" @click="handleAdd">新增</a-button>
<a-button v-permission="'module:resource:import'" @click="handleImport">导入</a-button>
<a-dropdown v-permission="'module:resource:export'">...</a-dropdown>
<a-button v-if="selected.length > 0" v-permission="'module:resource:batchDelete'" danger>批量删除</a-button>

<!-- Table action column -->
<a-button v-permission="'module:resource:edit'" @click="handleEdit(record)">编辑</a-button>
<a-popconfirm @confirm="handleDelete(record.id)">
  <a-button v-permission="'module:resource:delete'" danger>删除</a-button>
</a-popconfirm>
```

### Button Permission Bootstrap Pattern

```go
if pageMenuID > 0 {
    buttons := []model.SysMenu{
        {ParentID: pageMenuID, Name: "新增", Sort: 1, Type: 3, Permission: "module:resource:add", Status: 1},
        {ParentID: pageMenuID, Name: "编辑", Sort: 2, Type: 3, Permission: "module:resource:edit", Status: 1},
        {ParentID: pageMenuID, Name: "删除", Sort: 3, Type: 3, Permission: "module:resource:delete", Status: 1},
    }
    for _, def := range buttons {
        menu := def
        global.DB.Where("permission = ?", menu.Permission).
            Attrs(model.SysMenu{ParentID: menu.ParentID, Name: menu.Name, Sort: menu.Sort, Type: menu.Type, Status: menu.Status, Hidden: 0}).
            FirstOrCreate(&menu)
        menuIDs = append(menuIDs, menu.ID)
    }
}

### Frontend Design Skill

When building or refactoring frontend pages, apply these design principles:

- **Admin constraint first**: `frontend-design` is used to improve hierarchy, spacing, loading/error states, density, and dark-mode consistency inside the existing Ant Design Vue admin style. It must not override project admin conventions.
- **Layout**: White card sections on `#f5f5f5` background, `border-radius: 8px`, light `box-shadow`. Group related form fields into visual sections.
- **Drawer/Form**: `layout="vertical"`, header on `#fafafa` with `2px solid #f0f0f0` border.
- **Detail views**: Card-based, not flat `a-descriptions`. Top info bar with tags.
- **Interactive elements**: Circle badges for option keys, large toggle buttons for binary choices, numbered indices for ordered items.
- **Typography & color**: `13px` body, `12px` labels `#8c8c8c`, `15px` titles `#262626`. Green `#52c41a` success, orange `#fa8c16` warnings, blue `#1677ff` primary.
- **No generic aesthetics**: Every page must have visual sections, clear hierarchy, polished micro-interactions, but admin pages must avoid landing-page/marketing-style hero blocks, slogans, decorative chips, and promotional top sections unless explicitly requested.

## Reference: Menu & API Bootstrap

### Fresh Install (`initDefaultData`)

```go
xxxMgmt := model.SysMenu{ParentID: 0, Name: "XXX管理", Path: "/xxx", Component: "Layout", Icon: "...", Sort: N, Type: 1, Permission: "xxx", Status: 1}
global.DB.Create(&xxxMgmt)
xxxMenus := []model.SysMenu{
    {ParentID: xxxMgmt.ID, Name: "子菜单", Path: "/xxx/sub", Component: "xxx/sub/index", Sort: 1, Type: 2, Permission: "xxx:sub:list", Status: 1},
}
global.DB.Create(&xxxMenus)
```

### Upgrade (`ensureBuiltInData`)

Follow `ensureExamMenus` pattern:
- Root menu → `FirstOrCreate` by `permission`
- Sub-menus → `FirstOrCreate` by `permission`, `Attrs` for create-only fields
- Button permissions (Type=3) → one per action, under each sub-menu
- `grantMenusToRoleCodes(menuIDs, []string{"admin", "system_admin"})` — ALL IDs including buttons
- API entries → `ensureApiAccessForRoleCodes(api, []string{"admin", "system_admin"})`

### Key Rules

- Never skip either path (fresh + upgrade)
- `FirstOrCreate + Attrs` — never overwrite existing customizations
- Menu `Permission` is unique key
- API `Path + Method` is unique key
- Always grant to `admin` + `system_admin`

## Reference: Excel Import / Export / Template

### Reference Files

- Framework: `server/utils/excel.go`
- Example: `server/service/user/user_excel.go`
- Handlers: `server/api/v1/user.go` (GetUserImportTemplate, ImportUsers, ExportUsers)
- Frontend: `web/src/api/user.ts`, `web/src/views/admin/system/user/index.vue`
- Result modal: `web/src/views/admin/system/user/components/ImportResultModal.vue`

### Backend: Three Functions Per Module

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
        {Header: "列A", Key: "field_a", Required: true, Type: "string", MaxLen: 50},
        {Header: "列B", Key: "field_b", Type: "string", Enum: []string{"选项1", "选项2"}},
    }
}
```

### Backend: Key Rules

- **Template**: Include context name in sheet/filename, one example row, `AddDataValidation` for enums
- **Import**: Context ID as param not column, `utils.ValidateImport`, batch dedup, file ≤10MB
- **Export**: Context ID + optional `ids`, filename with context name
- **Dict mapping**: Keep in `_excel.go`, align with dict values

### Frontend: API + Toolbar

```ts
export function downloadXxxImportTemplate(contextId?: number) {
  return request.get('/xxx/import-template', { params: { dept_id: contextId }, responseType: 'blob' })
}
export function importXxx(file: File, contextId: number) {
  const formData = new FormData(); formData.append('file', file); formData.append('dept_id', String(contextId))
  return request.post<any, ApiResponse<ImportResult>>('/xxx/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
}
export function exportXxx(contextId: number, ids?: number[]) {
  const params: Record<string, any> = { dept_id: contextId }
  if (ids?.length) params.ids = ids.join(',')
  return request.get('/xxx/export', { params, responseType: 'blob' })
}
```

- All import/export buttons use `v-permission`
- Reuse `ImportResultModal.vue` for import results
