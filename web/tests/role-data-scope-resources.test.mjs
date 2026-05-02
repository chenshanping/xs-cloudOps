import test from 'node:test'
import assert from 'node:assert/strict'
import { mkdtempSync, readFileSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { resolve, join } from 'node:path'
import ts from 'typescript'

function readUtf8(relativePath) {
  return readFileSync(resolve(process.cwd(), relativePath), 'utf8')
}

async function loadTsModule(relativePath) {
  const sourcePath = resolve(process.cwd(), relativePath)
  const source = readFileSync(sourcePath, 'utf8')
  const output = ts.transpileModule(source, {
    compilerOptions: {
      module: ts.ModuleKind.ESNext,
      target: ts.ScriptTarget.ES2020
    },
    fileName: sourcePath
  })

  const tempDir = mkdtempSync(join(tmpdir(), 'role-data-scope-resources-'))
  const tempFile = join(tempDir, 'module.mjs')
  writeFileSync(tempFile, output.outputText, 'utf8')
  return import(`${new URL(`file://${tempFile}`).href}?t=${Date.now()}`)
}

test('buildRoleFeatureDataScopeForm builds editable items from known resources and role scopes', async () => {
  const helpers = await loadTsModule('src/views/admin/system/role/components/dataScopeResources.ts')

  const resources = [
    { code: 'system:user-management', label: '用户管理', description: 'desc', owner_fields: ['dept_id'] },
    { code: 'system:dept-management', label: '部门管理', description: 'desc', owner_fields: ['dept_id'] }
  ]
  const scopes = [
    {
      resource_code: 'system:user-management',
      data_scope: 2,
      depts: [{ id: 11 }, { id: 12 }]
    },
    {
      resource_code: 'legacy:archived-resource',
      data_scope: 5,
      depts: [{ id: 99 }]
    }
  ]

  assert.deepEqual(helpers.buildRoleFeatureDataScopeForm(resources, scopes), [
    {
      resource_code: 'system:user-management',
      data_scope: 2,
      dept_ids: [11, 12]
    },
    {
      resource_code: 'system:dept-management',
      data_scope: 0,
      dept_ids: []
    }
  ])
})

test('unknown resource scopes are split out and preserved in final payload', async () => {
  const helpers = await loadTsModule('src/views/admin/system/role/components/dataScopeResources.ts')

  const resources = [
    { code: 'system:user-management', label: '用户管理', description: 'desc', owner_fields: ['dept_id'] }
  ]
  const scopes = [
    {
      resource_code: 'system:user-management',
      data_scope: 2,
      depts: [{ id: 21 }]
    },
    {
      resource_code: 'legacy:archived-resource',
      data_scope: 5,
      depts: [{ id: 88 }]
    }
  ]

  const { knownScopes, unknownScopes } = helpers.splitKnownAndUnknownFeatureDataScopes(resources, scopes)
  assert.deepEqual(knownScopes, [scopes[0]])
  assert.deepEqual(unknownScopes, [
    {
      resource_code: 'legacy:archived-resource',
      data_scope: 5,
      dept_ids: [88]
    }
  ])

  const payload = helpers.buildRoleFeatureDataScopePayload(
    [
      {
        resource_code: 'system:user-management',
        data_scope: 2,
        dept_ids: [31, 32]
      }
    ],
    unknownScopes
  )

  assert.deepEqual(payload, [
    {
      resource_code: 'system:user-management',
      data_scope: 2,
      dept_ids: [31, 32]
    },
    {
      resource_code: 'legacy:archived-resource',
      data_scope: 5,
      dept_ids: [88]
    }
  ])
})

test('resource-specific scope options hide unsupported self scope', async () => {
  const helpers = await loadTsModule('src/views/admin/system/role/components/dataScopeResources.ts')

  const userResource = {
    code: 'system:user-management',
    label: '用户管理',
    description: 'desc',
    owner_fields: ['dept_id', 'created_by']
  }
  const deptResource = {
    code: 'system:dept-management',
    label: '部门管理',
    description: 'desc',
    owner_fields: ['dept_id']
  }

  assert.deepEqual(
    helpers.getSupportedFeatureScopeOptions(userResource).map(item => item.value),
    [0, 1, 2, 3, 4, 5]
  )
  assert.deepEqual(
    helpers.getSupportedFeatureScopeOptions(deptResource).map(item => item.value),
    [0, 1, 2, 3, 4]
  )
})

test('useRolePermissionDrawer loads data scope resources from backend API', () => {
  const source = readUtf8('src/views/admin/system/role/components/useRolePermissionDrawer.ts')

  assert.match(source, /getDataScopeResources/)
  assert.match(
    source,
    /Promise\.all\(\[\s*fetchMenuTree\(requestToken\)\s*,\s*fetchAllApis\(requestToken\)\s*,\s*fetchDataScopeResources\(requestToken\)\s*\]\)/,
    'expected drawer open flow to fetch data scope resources in parallel'
  )
})

test('useRolePermissionDrawer resets stale interactive state and guards async loading', () => {
  const source = readUtf8('src/views/admin/system/role/components/useRolePermissionDrawer.ts')

  assert.match(source, /const\s+permissionLoading\s*=\s*ref\(false\)/)
  assert.match(source, /const\s+permissionRequestToken\s*=\s*ref\(0\)/)
  assert.match(source, /const\s+resetPermissionState\s*=\s*\(\)\s*=>\s*\{/)
  assert.match(source, /selectedMenuKeys\.value\s*=\s*\[\]/)
  assert.match(source, /checkedApiIds\.value\s*=\s*\[\]/)
  assert.match(source, /defaultDataScope\.value\s*=\s*1/)
  assert.match(source, /resourceDefinitions\.value\s*=\s*\[\]/)
  assert.match(source, /featureDataScopes\.value\s*=\s*\[\]/)
  assert.match(
    source,
    /const\s+requestToken\s*=\s*\+\+permissionRequestToken\.value[\s\S]*permissionLoading\.value\s*=\s*true[\s\S]*resetPermissionState\(\)[\s\S]*await Promise\.all\(\[\s*fetchMenuTree\(requestToken\)\s*,\s*fetchAllApis\(requestToken\)\s*,\s*fetchDataScopeResources\(requestToken\)\s*\]\)/,
    'expected opening flow to synchronously enter loading mode, reset stale state, and pass request token through loaders'
  )
  assert.match(
    source,
    /if\s*\(requestToken\s*!==\s*permissionRequestToken\.value\)\s*\{\s*return\s*\}/,
    'expected stale async responses to be ignored by request token guard'
  )
})

test('RolePermissionDrawer disables save and suspends content while permissions load', () => {
  const source = readUtf8('src/views/admin/system/role/components/RolePermissionDrawer.vue')

  assert.match(source, /permissionLoading/)
  assert.match(source, /:spinning="permissionLoading"/)
  assert.match(source, /:disabled="permissionLoading"/)
})

test('unknown feature data scopes are preserved in the save path', () => {
  const source = readUtf8('src/views/admin/system/role/components/useRolePermissionDrawer.ts')

  assert.match(source, /const\s+unknownFeatureDataScopes\s*=\s*ref<[^>]+>\(\[\]\)/)
  assert.match(
    source,
    /splitKnownAndUnknownFeatureDataScopes\(\s*scopeResources\s*,\s*res\.data\.feature_data_scopes\s*\|\|\s*\[\]\s*\)/,
    'expected unknown scopes to be separated from rendered resources during load'
  )
  assert.match(
    source,
    /buildRoleFeatureDataScopePayload\(\s*featureDataScopes\.value\s*,\s*unknownFeatureDataScopes\.value\s*\)/,
    'expected save payload builder to merge rendered scopes with unknown preserved scopes'
  )
})

test('dataScopeResources no longer hardcodes role feature scope resource array', () => {
  const source = readUtf8('src/views/admin/system/role/components/dataScopeResources.ts')

  assert.doesNotMatch(source, /ROLE_FEATURE_SCOPE_RESOURCES[\s\S]*=\s*\[/)
})

test('RolePermissionDataScopePanel renders scope options from resource capabilities', () => {
  const source = readUtf8('src/views/admin/system/role/components/RolePermissionDataScopePanel.vue')

  assert.match(source, /getSupportedFeatureScopeOptions/)
  assert.doesNotMatch(source, /const\s+scopeOptions\s*=\s*FEATURE_SCOPE_OPTIONS/)
})
