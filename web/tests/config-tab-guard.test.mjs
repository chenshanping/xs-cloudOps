import test from 'node:test'
import assert from 'node:assert/strict'
import {
  cloneFromSnapshot,
  createSnapshot,
  hasDirtyConfigTabs,
  isSnapshotDirty,
} from '../dist-tests-config-tab-guard/src/views/admin/system/config/config-tab-guard.js'

test('isSnapshotDirty returns false for unchanged form state', () => {
  const snapshot = createSnapshot({
    sys_name: 'Go RBAC Admin',
    front_mode: 'full',
  })

  assert.equal(isSnapshotDirty(snapshot, {
    sys_name: 'Go RBAC Admin',
    front_mode: 'full',
  }), false)
})

test('isSnapshotDirty returns true after nested state changes', () => {
  const snapshot = createSnapshot({
    storage_type: 'minio',
    config: {
      endpoint: '127.0.0.1:9000',
      use_ssl: false,
    },
  })

  assert.equal(isSnapshotDirty(snapshot, {
    storage_type: 'minio',
    config: {
      endpoint: '127.0.0.1:9000',
      use_ssl: true,
    },
  }), true)
})

test('cloneFromSnapshot restores the original draft payload', () => {
  const original = {
    login_title: '欢迎回来',
    features: [{ icon: 'CheckCircleOutlined', title: '智能分析' }],
  }
  const snapshot = createSnapshot(original)
  const restored = cloneFromSnapshot(snapshot)

  assert.deepEqual(restored, original)
  assert.notEqual(restored, original)
})

test('hasDirtyConfigTabs detects whether any config tab is dirty', () => {
  assert.equal(hasDirtyConfigTabs({
    basic: false,
    file: false,
    login: false,
    email: false,
    ai: false,
  }), false)

  assert.equal(hasDirtyConfigTabs({
    basic: false,
    file: true,
    login: false,
    email: false,
    ai: false,
  }), true)
})
