import test from 'node:test'
import assert from 'node:assert/strict'
import {
  normalizeGenderDictOptions,
  resolveGenderLabel,
} from '../dist-tests-user-gender/src/views/admin/system/user/user-gender.js'

test('normalizes dict values to numbers for form binding and search', () => {
  const options = normalizeGenderDictOptions([
    { label: '未知', value: '0', tag_type: 'default' },
    { label: '男', value: '1', tag_type: 'processing' },
    { label: '女', value: '2', tag_type: 'pink' },
  ])

  assert.deepEqual(
    options.map(item => item.value),
    [0, 1, 2],
  )
})

test('resolves gender labels from normalized options', () => {
  const options = normalizeGenderDictOptions([
    { label: '未知', value: '0', tag_type: 'default' },
    { label: '男', value: '1', tag_type: 'processing' },
    { label: '女', value: '2', tag_type: 'pink' },
  ])

  assert.equal(resolveGenderLabel(options, 1), '男')
  assert.equal(resolveGenderLabel(options, 0), '未知')
  assert.equal(resolveGenderLabel(options, 9), '-')
})
