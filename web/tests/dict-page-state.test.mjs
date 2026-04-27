import test from 'node:test'
import assert from 'node:assert/strict'
import {
  canApplyDictDataResponse,
  filterDictTypes,
  reconcileSelectedType,
} from '../dist-tests-dict-page-state/src/views/admin/system/dict/dict-page-state.js'

const typeList = [
  { id: 1, name: '性别', type: 'sys_gender', status: 1, remark: '' },
  { id: 2, name: '用户状态', type: 'sys_user_status', status: 1, remark: '' },
  { id: 3, name: '审批状态', type: 'approval_status', status: 0, remark: '' },
]

test('filters dict types by name or type code', () => {
  assert.deepEqual(
    filterDictTypes(typeList, 'status').map(item => item.id),
    [2, 3],
  )
  assert.deepEqual(
    filterDictTypes(typeList, '性别').map(item => item.id),
    [1],
  )
})

test('keeps full list when search text is empty', () => {
  assert.deepEqual(filterDictTypes(typeList, '').map(item => item.id), [1, 2, 3])
})

test('reconciles selected type with refreshed list data when record still exists', () => {
  const selected = { id: 2, name: '旧名称', type: 'sys_user_status', status: 0, remark: 'old' }

  assert.deepEqual(reconcileSelectedType(typeList, selected), typeList[1])
})

test('keeps previous selection when current filtered list does not contain it', () => {
  const selected = { id: 2, name: '用户状态', type: 'sys_user_status', status: 1, remark: '' }

  assert.deepEqual(reconcileSelectedType([typeList[0]], selected), selected)
})

test('only applies dict data responses for the current selected type', () => {
  assert.equal(canApplyDictDataResponse('sys_gender', 'sys_gender'), true)
  assert.equal(canApplyDictDataResponse('approval_status', 'sys_gender'), false)
  assert.equal(canApplyDictDataResponse('sys_gender', ''), false)
})
