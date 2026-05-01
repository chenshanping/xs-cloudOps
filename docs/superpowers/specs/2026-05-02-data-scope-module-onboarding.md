# Business-Module Data Scope Onboarding Contract

## Status

Accepted

## Date

2026-05-02

## Purpose

This document locks the existing business-module data scope onboarding contract in `go-base`.
It does not introduce a new permission model.
It explains the minimum conditions a module must satisfy to reuse the current resource-level data scope framework without drift.

## Non-Goals

- Redesigning role data scope semantics
- Splitting one business module into multiple resource codes by default
- Adding frontend-only toggles or alternative permission paths

## Existing Framework Anchors

- Supported resource registry: `server/service/core/data_scope_resource.go`
- Resource scope resolution: `server/service/core/data_scope.go`
- Generic record filter helper: `server/service/core/record_data_scope.go`
- Role assignment normalization and guardrail: `server/service/role/role.go`

## Onboarding Contract

### 1. `resource_code` must be stable and backend-registered

Every business module that wants resource-level data scope must define one stable `resource_code` and keep using that exact code across backend and admin permission assignment.

Required rules:

- The code must be added to the backend-supported registry before role permission saving can use it.
- The code is a module contract, not a page-local string. Do not generate it dynamically.
- Unsupported codes are rejected during role permission normalization instead of being silently accepted.

Practical consequence:

- If a module uses `biz:order` for list access, then detail, edit, delete, and batch actions for that same module must also resolve scope with `biz:order`.
- If a future requirement truly needs different scope semantics per action, that is a permission-model change and must be reviewed separately. It is not part of normal onboarding.

### 2. The business table must expose at least one ownership column

The generic helper assumes business ownership can be expressed with one or both of these columns:

- `dept_id`
- `created_by`

Required rules:

- A module must have at least one of these columns before using the generic record scope helper.
- If both columns exist, bind both so department scope and creator/self scope can be combined.
- If only `created_by` exists, bind only `created_by`.
- If only `dept_id` exists, bind only `dept_id`.

This is the current contract already supported by `core.ApplyRecordDataScope`.
The regression suite now locks both single-column cases.

### 2.1 `dept_id`-only modules do not get `仅本人` for free

A `dept_id`-only module can safely onboard department-oriented data scopes such as:

- `全部`
- `本部门`
- `本部门及以下`
- `自定义部门`

But `DataScopeSelf` has stricter requirements.

Required rules:

- If the module needs stable `仅本人` semantics, it must provide `created_by`, or explicitly bind `SelfColumn`.
- If the module only provides `dept_id` and leaves both `CreatedByColumn` and `SelfColumn` empty, `AllowSelf` does not imply an automatic fallback.
- In that case, self scope is considered unsupported for that binding and the helper must fail closed.

Practical consequence:

- `dept_id`-only means "department scope is available".
- It does not mean "self scope is automatically supported".
- If onboarding needs `仅本人`, the schema or binding contract must be completed first instead of relying on implicit behavior.

### 3. List, detail, edit, delete, and batch operations must share one resource contract

A module is only considered onboarded when all direct record-access paths use the same `resource_code` and the same ownership interpretation.

Required rules:

- List queries must resolve scope with the module `resource_code` and apply the helper filter.
- Detail queries must use the same `resource_code`; do not bypass scope just because the request already has an ID.
- Edit and delete must use the same `resource_code`; do not rely on a prior list page having filtered the record.
- Batch operations must also reuse the same `resource_code` and filter the target set in the same scope model.

Forbidden pattern:

- List is scope-filtered, but detail or mutation paths query by ID only.

### 4. Onboarding must call the unified helper path

Do not re-implement ad hoc data scope SQL in each module.
The module must plug into the existing resolver and filtering helpers.

Required path:

1. Resolve the operator scope with `core.ResolveUserDataScopeForResource(operatorID, resourceCode)`.
2. Bind the business table ownership fields with `core.RecordDataScopeBinding`.
3. Apply the filter with `core.ApplyRecordDataScope(...)` for generic business records, or use the existing resource-aware core helpers for user and department management where those already exist.

Reference pattern:

```go
scope, err := core.ResolveUserDataScopeForResource(operatorID, resourceCode)
if err != nil {
	return err
}

query := core.ApplyRecordDataScope(db.Model(&BizRecord{}), scope, core.RecordDataScopeBinding{
	TableAlias:      "biz_record",
	DeptColumn:      "dept_id",
	CreatedByColumn: "created_by",
})
```

If a module only has one ownership column, leave the other binding field empty instead of inventing fallback SQL.
For `dept_id`-only bindings, this also means self scope will fail closed unless `SelfColumn` is explicitly provided.

## Onboarding Checklist

Before calling a module "data-scope onboarded", verify all of the following:

1. The module has one stable backend-registered `resource_code`.
2. The module table has `dept_id` and/or `created_by`.
3. List, detail, edit, delete, and batch paths all use the same `resource_code`.
4. The module resolves scope through `ResolveUserDataScopeForResource`.
5. If the module needs `仅本人`, it provides `created_by` or an explicit `SelfColumn`; otherwise self scope is treated as unsupported and fail-closed.
6. The module applies the unified helper or an existing core resource-aware helper instead of handwritten divergent filters.
7. Regression coverage exists for the module contract and any module-specific edge case.

## Regression Coverage Locked By Task 4

This task adds regression coverage for the shared onboarding contract itself:

- Unsupported `resource_code` is rejected by `normalizeRoleFeatureDataScopeAssignments`.
- `ApplyRecordDataScope` works when a table only exposes `created_by`.
- `ApplyRecordDataScope` works when a table only exposes `dept_id`.
- `ApplyRecordDataScope` fails closed for `dept_id`-only bindings when `AllowSelf` is requested without `created_by` or `SelfColumn`.

Those checks intentionally lock the current framework behavior without expanding the permission model.
