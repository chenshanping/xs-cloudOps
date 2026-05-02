# Data Scope Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Harden the existing data-scope capability into a reusable foundation so future business CRUD modules can plug into the same resource-code and ownership-field model without redesigning data permissions.

**Architecture:** Keep the existing `sys_role.data_scope + sys_role_data_scope` model, but remove the remaining hardcoded assumptions. The backend becomes the source of truth for supported data-scope resources and provides a generic record-scope helper for `dept_id` / `created_by` tables; the frontend role-permission drawer consumes resource metadata dynamically; onboarding guidance and regression tests lock the contract for future modules.

**Tech Stack:** Go, Gin, Gorm, SQLite-backed integration tests, Vue 3, TypeScript, Ant Design Vue, Node test scripts

---

### Task 1: Make Backend Resource Registry the Single Source of Truth

**Files:**
- Create: `server/service/core/data_scope_resource_test.go`
- Modify: `server/service/core/data_scope_resource.go`
- Modify: `server/api/v1/role.go`
- Modify: `server/router/modules/role.go`
- Modify: `server/model/response/common.go` or an adjacent response type file if a typed response struct is preferred
- Test: `server/service/core/data_scope_resource_test.go`

- [ ] **Step 1: Write the failing backend resource-registry tests**

```go
package core

import "testing"

func TestSupportedDataScopeResourcesExposeStableMetadata(t *testing.T) {
	resources := SupportedDataScopeResources()
	if len(resources) < 2 {
		t.Fatalf("expected built-in resources")
	}

	userResource := resources[0]
	if userResource.Code != DataScopeResourceUserManagement {
		t.Fatalf("unexpected first resource code: %s", userResource.Code)
	}
	if userResource.Label == "" || userResource.Description == "" {
		t.Fatalf("resource metadata should be complete: %+v", userResource)
	}
	if len(userResource.OwnerFields) == 0 {
		t.Fatalf("resource should declare owner fields")
	}
}

func TestSupportedDataScopeResourcesReturnCopy(t *testing.T) {
	resources := SupportedDataScopeResources()
	resources[0].Label = "mutated"

	fresh := SupportedDataScopeResources()
	if fresh[0].Label == "mutated" {
		t.Fatalf("supported resources should return a defensive copy")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./service/core -run "TestSupportedDataScopeResources"`

Expected: FAIL because the current `DataScopeResource` only exposes `Code` and `Label`, and there is no richer metadata contract yet.

- [ ] **Step 3: Implement richer resource metadata and resource-list endpoint**

```go
package core

type DataScopeOwnerField string

const (
	DataScopeOwnerFieldDeptID    DataScopeOwnerField = "dept_id"
	DataScopeOwnerFieldCreatedBy DataScopeOwnerField = "created_by"
)

type DataScopeResource struct {
	Code        string               `json:"code"`
	Label       string               `json:"label"`
	Description string               `json:"description"`
	OwnerFields []DataScopeOwnerField `json:"owner_fields"`
}

var supportedDataScopeResources = []DataScopeResource{
	{
		Code:        DataScopeResourceUserManagement,
		Label:       "用户管理",
		Description: "控制用户列表及关联用户操作可见的数据范围。",
		OwnerFields: []DataScopeOwnerField{DataScopeOwnerFieldDeptID, DataScopeOwnerFieldCreatedBy},
	},
	{
		Code:        DataScopeResourceDeptManagement,
		Label:       "部门管理",
		Description: "控制部门树、可管理部门及部门统计可见的数据范围。",
		OwnerFields: []DataScopeOwnerField{DataScopeOwnerFieldDeptID},
	},
}
```

```go
// server/api/v1/role.go
func (a *RoleApi) GetDataScopeResources(c *gin.Context) {
	response.OkWithData(c, core.SupportedDataScopeResources())
}
```

```go
// server/router/modules/role.go
R(rg, "GET", "/roles/data-scope-resources", m.Name(), "数据权限资源列表", v1.Role.GetDataScopeResources, registry.WithAuth())
```

- [ ] **Step 4: Run tests to verify it passes**

Run: `go test ./service/core -run "TestSupportedDataScopeResources"`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/service/core/data_scope_resource.go server/service/core/data_scope_resource_test.go server/api/v1/role.go server/router/modules/role.go
git commit -m "feat: expose data scope resource registry metadata"
```

### Task 2: Add a Generic Record-Level Data Scope Helper for Future Business Modules

**Files:**
- Create: `server/service/core/record_data_scope.go`
- Modify: `server/service/core/data_scope.go`
- Modify: `server/tests/data_scope_test.go`
- Test: `server/tests/data_scope_test.go`

- [ ] **Step 1: Write the failing generic record-scope integration test**

```go
func TestApplyRecordDataScopeUsesDeptAndCreatorOwnership(t *testing.T) {
	db := setupDataScopeTestDB(t)
	createRoleFeatureScopeTables(t, db)

	type BizRecord struct {
		ID        uint `gorm:"primarykey"`
		Name      string
		DeptID    uint
		CreatedBy uint
	}
	if err := db.AutoMigrate(&BizRecord{}); err != nil {
		t.Fatalf("auto migrate biz records: %v", err)
	}

	// seed dept tree, operator role, and three records:
	// one in-scope dept, one created-by self, one out of scope.
	// then assert ApplyRecordDataScope only returns the expected rows.
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./tests -run "TestApplyRecordDataScopeUsesDeptAndCreatorOwnership"`

Expected: FAIL because there is no generic helper yet for `dept_id` / `created_by` business tables.

- [ ] **Step 3: Implement the generic helper without changing the validated user/dept behavior**

```go
package core

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type RecordDataScopeBinding struct {
	TableAlias      string
	DeptColumn      string
	CreatedByColumn string
	SelfColumn      string
}

func ApplyRecordDataScope(db *gorm.DB, scope *UserDataScope, binding RecordDataScopeBinding) *gorm.DB {
	if scope == nil {
		return db.Where("1 = 0")
	}
	if scope.All {
		return db
	}

	conditions := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)

	if len(scope.DeptIDs) > 0 && binding.DeptColumn != "" {
		conditions = append(conditions, fmt.Sprintf("%s.%s IN ?", binding.TableAlias, binding.DeptColumn))
		args = append(args, scope.DeptIDs)
	}
	if len(scope.CreatorIDs) > 0 && binding.CreatedByColumn != "" {
		conditions = append(conditions, fmt.Sprintf("%s.%s IN ?", binding.TableAlias, binding.CreatedByColumn))
		args = append(args, scope.CreatorIDs)
	}
	if scope.AllowSelf && binding.SelfColumn != "" {
		conditions = append(conditions, fmt.Sprintf("%s.%s = ?", binding.TableAlias, binding.SelfColumn))
		args = append(args, scope.OperatorID)
	}

	if len(conditions) == 0 {
		return db.Where("1 = 0")
	}
	return db.Where(strings.Join(conditions, " OR "), args...)
}
```

Implementation notes:
- Keep `ApplyUserDataScope` and current user/dept management behavior intact.
- `ApplyRecordDataScope` is the new default helper for future business tables, not a rewrite of existing validated modules.
- For business modules following this foundation, pass `CreatedByColumn: "created_by"` and leave `SelfColumn` empty unless a module explicitly uses `owner_user_id`.

- [ ] **Step 4: Run tests to verify it passes**

Run: `go test ./tests -run "TestApplyRecordDataScopeUsesDeptAndCreatorOwnership|TestUserServiceGetUserListUsesCreatorForSelfScopeInUserManagement"`

Expected: PASS, and existing user-management self-scope tests stay green.

- [ ] **Step 5: Commit**

```bash
git add server/service/core/record_data_scope.go server/service/core/data_scope.go server/tests/data_scope_test.go
git commit -m "feat: add generic record data scope helper"
```

### Task 3: Make the Role Permission Drawer Load Resource Metadata from Backend

**Files:**
- Modify: `web/src/api/role.ts`
- Modify: `web/src/views/admin/system/role/components/dataScopeResources.ts`
- Modify: `web/src/views/admin/system/role/components/useRolePermissionDrawer.ts`
- Modify: `web/src/views/admin/system/role/components/RolePermissionDataScopePanel.vue`
- Create: `web/tests/role-data-scope-resources.test.mjs`
- Create: `web/scripts/test-role-data-scope-resources.mjs`
- Modify: `web/package.json`
- Test: `web/tests/role-data-scope-resources.test.mjs`

- [ ] **Step 1: Write the failing frontend resource-loading test**

```js
import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

test('role permission drawer no longer hardcodes built-in data-scope resources', () => {
  const source = readFileSync(
    resolve(process.cwd(), 'src/views/admin/system/role/components/useRolePermissionDrawer.ts'),
    'utf8'
  )

  assert.match(source, /getDataScopeResources\(/)
  assert.doesNotMatch(source, /ROLE_FEATURE_SCOPE_RESOURCES = \[/)
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm run test:role-data-scope-resources`

Expected: FAIL because the current drawer still depends on the hardcoded `ROLE_FEATURE_SCOPE_RESOURCES` constant.

- [ ] **Step 3: Implement dynamic loading with frontend fallback-safe shaping**

```ts
// web/src/api/role.ts
export interface DataScopeResource {
  code: string
  label: string
  description: string
  owner_fields: Array<'dept_id' | 'created_by'>
}

export function getDataScopeResources() {
  return request.get<any, ApiResponse<DataScopeResource[]>>('/roles/data-scope-resources')
}
```

```ts
// useRolePermissionDrawer.ts
const resourceDefinitions = ref<DataScopeResource[]>([])

const fetchDataScopeResources = async () => {
  const res = await getDataScopeResources()
  resourceDefinitions.value = res.data
}

watch(visible, async val => {
  if (!val) return
  await Promise.all([fetchMenuTree(), fetchAllApis(), fetchDataScopeResources()])
  featureDataScopes.value = buildRoleFeatureDataScopeForm(resourceDefinitions.value, role.feature_data_scopes)
})
```

```ts
// dataScopeResources.ts
export function buildRoleFeatureDataScopeForm(
  resources: DataScopeResourceDefinition[],
  scopes?: RoleFeatureDataScope[]
): RoleFeatureDataScopeFormItem[] {
  const scopeMap = new Map((scopes || []).map(scope => [scope.resource_code, scope]))
  return resources.map(resource => ({
    resource_code: resource.code,
    data_scope: scopeMap.get(resource.code)?.data_scope ?? 0,
    dept_ids: scopeMap.get(resource.code)?.depts?.map(item => item.id) || []
  }))
}
```

- [ ] **Step 4: Run tests to verify it passes**

Run:

```bash
npm run test:role-data-scope-resources
npm run test:config-component-imports
```

Expected:
- the new role data-scope resource test passes
- the existing component import smoke test still passes

- [ ] **Step 5: Commit**

```bash
git add web/src/api/role.ts web/src/views/admin/system/role/components/dataScopeResources.ts web/src/views/admin/system/role/components/useRolePermissionDrawer.ts web/src/views/admin/system/role/components/RolePermissionDataScopePanel.vue web/tests/role-data-scope-resources.test.mjs web/scripts/test-role-data-scope-resources.mjs web/package.json
git commit -m "feat: load role data scope resources dynamically"
```

### Task 4: Document the Business-Module Onboarding Contract and Lock It with Regression Coverage

**Files:**
- Create: `docs/superpowers/specs/2026-05-02-data-scope-module-onboarding.md`
- Modify: `server/tests/data_scope_test.go`
- Modify: `web/src/views/admin/system/role/components/RolePermissionDataScopePanel.vue` (only if field-help copy needs the onboarding contract wording)
- Test: `server/tests/data_scope_test.go`

- [ ] **Step 1: Write the failing regression test for unsupported resource registration**

```go
func TestNormalizeRoleFeatureDataScopeAssignmentsRejectsUnsupportedResource(t *testing.T) {
	_, err := normalizeRoleFeatureDataScopeAssignments([]request.RoleFeatureDataScopeAssignment{
		{
			ResourceCode: "biz:unknown-management",
			DataScope:    model.DataScopeDept,
		},
	})
	if err == nil {
		t.Fatalf("expected unsupported resource to be rejected")
	}
}
```

- [ ] **Step 2: Run test to verify it fails or is missing**

Run: `go test ./tests ./service/role -run "TestNormalizeRoleFeatureDataScopeAssignmentsRejectsUnsupportedResource"`

Expected: FAIL or missing-coverage signal, depending on where the helper is currently tested.

- [ ] **Step 3: Write the onboarding guide and finish the missing regression assertions**

```md
# Business Module Data Scope Onboarding

## Required inputs
- A stable `resource_code`
- At least one ownership field: `dept_id` or `created_by`
- One canonical query entry point that applies `ResolveUserDataScopeForResource`

## Required checks
- List uses the module resource code
- Detail uses the same resource code
- Edit/Delete/Batch operations use the same resource code
- Role management exposes the resource in the data-scope drawer
```

Testing additions in `server/tests/data_scope_test.go` should cover:
- unsupported resource code is rejected
- a resource with `created_by` only can still use the generic helper
- a resource with `dept_id` only can still use the generic helper

- [ ] **Step 4: Run tests to verify it passes**

Run:

```bash
go test ./tests -run "TestApplyRecordDataScope|TestNormalizeRoleFeatureDataScopeAssignmentsRejectsUnsupportedResource"
go test ./...
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add docs/superpowers/specs/2026-05-02-data-scope-module-onboarding.md server/tests/data_scope_test.go
git commit -m "docs: add business module data scope onboarding contract"
```

### Task 5: Final Verification

**Files:**
- Modify: none
- Test: existing backend and frontend verification commands only

- [ ] **Step 1: Run backend verification**

Run: `go test ./...`

Expected: PASS

- [ ] **Step 2: Run targeted frontend verification**

Run:

```bash
npm run test:role-data-scope-resources
npm run test:config-component-imports
```

Expected: PASS

- [ ] **Step 3: Run strongest available frontend structural check**

Run: `npm run typecheck`

Expected:
- If the known `vue-tsc` toolchain replacement error still occurs, record it explicitly as an existing environment issue.
- If it passes in this session, keep the output in the task notes.

- [ ] **Step 4: Summarize the finished contract**

Record in the PR or handoff summary:
- backend is the source of truth for `resource_code` metadata
- future modules use `ApplyRecordDataScope`
- onboarding requires `dept_id` and/or `created_by`
- role-permission drawer automatically reflects newly registered backend resources

- [ ] **Step 5: Commit or prepare merge-ready diff**

```bash
git status
git log --oneline -5
```

Expected: working tree only contains intentional plan-scope changes, or is fully committed if completing in one branch.
