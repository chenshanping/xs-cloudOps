# List Page Pattern

## Standard Structure

For a standard admin CRUD list page:

1. define base columns first
2. build the action column separately
3. call `useTableColumns(baseColumns, actionColumn, rowActionPermissions)`
4. render through `ProTable`
5. keep search, toolbar, and table in one compact page

## Header Rule

For a standard CRUD list page:

- default to a compact layout that starts directly with `ProTable`
- keep the page focused on search, toolbar, and table actions
- do not add decorative hero headers, marketing copy, or summary cards unless the user explicitly asks for them

## Operation Column Rule

Always separate:

- toolbar permissions
- row button permissions
- operation column visibility

Correct pattern:

```ts
const columns = useTableColumns(
  [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '用户名', dataIndex: 'username', key: 'username' },
  ],
  { title: '操作', key: 'action', width: 200 },
  ['system:user:edit', 'system:user:delete']
)
```

Incorrect pattern:

- using add-only permissions to decide whether the action column exists
- rendering an operation column when the user has no row-level action permission

That creates an empty `操作` column.

## Toolbar Rule

Toolbar buttons stay on their own checks:

- `v-permission="'system:user:add'"`
- `v-permission="'system:user:delete'"`
- `v-permission="'system:user:export'"`

Do not couple toolbar-only permissions to action-column rendering.

## Reuse Targets

Prefer these shared files before writing page-local replacements:

- `web/src/components/ProTable.vue`
- `web/src/components/AvatarUpload.vue`
- `web/src/components/ImageUpload.vue`
- `web/src/components/FileUpload.vue`
- `web/src/components/FilePreview.vue`
- `web/src/utils/permission.ts`
- `web/src/directives/permission.ts`
- `web/src/views/system/user/index.vue`
- `web/src/views/system/role/index.vue`

## Closure Rule

Before finishing a CRUD page:

- every newly added button should have a real handler
- every newly added image or file display should have preview or open behavior
- every exposed CRUD action should complete a usable local loop instead of stopping at a placeholder
