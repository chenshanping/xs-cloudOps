import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

function getVueImports(source) {
  const match = source.match(/import\s*\{([^}]*)\}\s*from\s*'vue'/)
  if (!match) {
    return []
  }
  return match[1]
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean)
}

test('LoginRegisterConfig imports computed when using computed()', () => {
  const filePath = resolve(process.cwd(), 'src/views/admin/system/config/components/LoginRegisterConfig.vue')
  const source = readFileSync(filePath, 'utf8')
  const vueImports = getVueImports(source)

  assert.match(source, /computed\(/)
  assert.ok(vueImports.includes('computed'), 'expected computed to be imported from vue')
})
