import test from 'node:test'
import assert from 'node:assert/strict'
import {
  FILE_PREVIEW_SIZE_LIMIT,
  getFilePreviewDescriptor,
  getPreferredPreviewMimeType,
} from '../dist-tests-file-preview/src/components/file-preview-utils.js'

test('supports online preview for pdf, docx, xlsx and pptx', () => {
  assert.equal(getFilePreviewDescriptor({ ext: 'pdf' }).kind, 'pdf')
  assert.equal(getFilePreviewDescriptor({ ext: 'docx' }).kind, 'docx')
  assert.equal(getFilePreviewDescriptor({ ext: 'xlsx' }).kind, 'excel')
  assert.equal(getFilePreviewDescriptor({ ext: 'pptx' }).kind, 'pptx')
})

test('marks legacy office formats as download only', () => {
  const doc = getFilePreviewDescriptor({ ext: 'doc' })
  const xls = getFilePreviewDescriptor({ ext: 'xls' })

  assert.equal(doc.kind, 'unsupported')
  assert.equal(doc.reason, 'legacy-office')
  assert.equal(xls.kind, 'unsupported')
  assert.equal(xls.reason, 'legacy-office')
})

test('blocks oversized office preview and keeps download fallback', () => {
  const result = getFilePreviewDescriptor({
    ext: 'xlsx',
    size: FILE_PREVIEW_SIZE_LIMIT + 1,
  })

  assert.equal(result.kind, 'unsupported')
  assert.equal(result.reason, 'too-large')
})

test('normalizes pdf blob mime type for preview', () => {
  const descriptor = getFilePreviewDescriptor({ ext: 'pdf', mimeType: 'application/octet-stream' })
  assert.equal(getPreferredPreviewMimeType(descriptor, 'application/octet-stream'), 'application/pdf')
})
