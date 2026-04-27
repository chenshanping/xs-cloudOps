import test from 'node:test'
import assert from 'node:assert/strict'
import {
  createStorageDraftsState,
  getActiveStorageConfig,
  buildStorageSavePayload,
  changeActiveStorageType,
} from '../dist-tests-storage-config-state/src/views/admin/system/config/storage-config-state.js'

test('createStorageDraftsState loads saved config for each storage type independently', () => {
  const state = createStorageDraftsState({
    storage_type: 'minio',
    storage_local_config: '{"base_path":"local-a","base_url":"/api/v1/upload"}',
    storage_minio_config: '{"endpoint":"127.0.0.1:9000","access_key_id":"ak","secret_access_key":"sk","bucket_name":"demo","use_ssl":false}',
  })

  assert.equal(state.activeType, 'minio')
  assert.equal(state.drafts.local.base_path, 'local-a')
  assert.equal(state.drafts.minio.endpoint, '127.0.0.1:9000')
})

test('changeActiveStorageType keeps previous draft untouched', () => {
  let state = createStorageDraftsState({
    storage_type: 'local',
    storage_local_config: '{"base_path":"local-a","base_url":"/api/v1/upload"}',
  })

  state.drafts.local.base_path = 'local-b'
  state = changeActiveStorageType(state, 'minio')
  state.drafts.minio.endpoint = '127.0.0.1:9000'
  state = changeActiveStorageType(state, 'local')

  assert.equal(getActiveStorageConfig(state).base_path, 'local-b')
  assert.equal(state.drafts.minio.endpoint, '127.0.0.1:9000')
})

test('buildStorageSavePayload updates only active type config key', () => {
  const state = createStorageDraftsState({
    storage_type: 'local',
    storage_local_config: '{"base_path":"uploads","base_url":"/api/v1/upload"}',
  })
  state.drafts.minio.endpoint = '127.0.0.1:9000'
  state.drafts.minio.access_key_id = 'ak'
  state.drafts.minio.secret_access_key = 'sk'
  state.drafts.minio.bucket_name = 'demo'

  const payload = buildStorageSavePayload(changeActiveStorageType(state, 'minio'))

  assert.equal(payload.storage_type, 'minio')
  assert.ok(payload.storage_minio_config)
  assert.equal(payload.storage_local_config, undefined)
})
