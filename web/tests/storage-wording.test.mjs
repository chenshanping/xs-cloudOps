import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

function readUtf8(relativePath) {
  return readFileSync(resolve(process.cwd(), relativePath), 'utf8')
}

test('FileSettings clarifies default upload storage only affects new uploads', () => {
  const source = readUtf8('src/views/admin/system/config/components/FileSettings.vue')

  assert.match(source, /默认上传存储/)
  assert.match(source, /仅影响后续新上传文件的默认存储位置/)
  assert.match(source, /不会自动迁移已有历史文件/)
  assert.match(source, /文件管理\s*>\s*文件迁移/)
})

test('FileMigrationDrawer clarifies migration direction and relation to default storage', () => {
  const source = readUtf8('src/views/admin/system/file/components/FileMigrationDrawer.vue')

  assert.match(source, /文件迁移只处理已上传的历史文件/)
  assert.match(source, /不会修改系统当前的默认上传存储配置/)
  assert.match(source, /源存储与目标存储相同，无需迁移/)
  assert.match(source, /目标存储已是当前默认上传位置/)
})
