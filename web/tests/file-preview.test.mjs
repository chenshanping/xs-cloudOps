import test from 'node:test'
import assert from 'node:assert/strict'
import {
  FILE_PREVIEW_SIZE_LIMIT,
  getFilePreviewDescriptor,
} from '../dist-tests-file-preview/src/components/file-preview-utils.js'

test('supports online preview for pdf, docx, xlsx and pptx', () => {
  assert.equal(getFilePreviewDescriptor({ ext: 'pdf' }).kind, 'pdf')
  assert.equal(getFilePreviewDescriptor({ ext: 'docx' }).kind, 'docx')
  assert.equal(getFilePreviewDescriptor({ ext: 'xlsx' }).kind, 'excel')
  assert.equal(getFilePreviewDescriptor({ ext: 'pptx' }).kind, 'pptx')
})

test('marks legacy office formats as download only', () => {
  const doc = getFilePreviewDescriptor({ ext: 'doc' })
  assert.equal(doc.kind, 'unsupported')
  assert.equal(doc.reason, 'legacy-office')
})

test('supports xls and ppt as previewable formats', () => {
  assert.equal(getFilePreviewDescriptor({ ext: 'xls' }).kind, 'excel')
  assert.equal(getFilePreviewDescriptor({ ext: 'ppt' }).kind, 'pptx')
  assert.equal(getFilePreviewDescriptor({ ext: 'csv' }).kind, 'excel')
})

test('supports code and markdown file kinds', () => {
  assert.equal(getFilePreviewDescriptor({ ext: 'md' }).kind, 'markdown')
  assert.equal(getFilePreviewDescriptor({ ext: 'js' }).kind, 'code')
  assert.equal(getFilePreviewDescriptor({ ext: 'go' }).kind, 'code')
  assert.equal(getFilePreviewDescriptor({ ext: 'json' }).kind, 'code')
})

test('supports epub file kind', () => {
  assert.equal(getFilePreviewDescriptor({ ext: 'epub' }).kind, 'epub')
})

test('blocks oversized office preview and keeps download fallback', () => {
  const result = getFilePreviewDescriptor({
    ext: 'xlsx',
    size: FILE_PREVIEW_SIZE_LIMIT + 1,
  })

  assert.equal(result.kind, 'unsupported')
  assert.equal(result.reason, 'too-large')
})

test('recognizes txt as text kind', () => {
  assert.equal(getFilePreviewDescriptor({ ext: 'txt' }).kind, 'text')
})
