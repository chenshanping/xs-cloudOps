import test from 'node:test'
import assert from 'node:assert/strict'
import { resolveViewModulePath } from '../dist-tests-view-resolver/src/router/view-resolver.js'

test('maps auth views to the auth directory', () => {
  assert.equal(
    resolveViewModulePath('login'),
    '../views/auth/login/index.vue',
  )
  assert.equal(
    resolveViewModulePath('/forgot-password/'),
    '../views/auth/forgot-password/index.vue',
  )
})

test('keeps front views under the front directory', () => {
  assert.equal(
    resolveViewModulePath('front/profile'),
    '../views/front/profile/index.vue',
  )
})

test('maps fixed admin root views into admin directory', () => {
  assert.equal(
    resolveViewModulePath('dashboard'),
    '../views/admin/dashboard/index.vue',
  )
  assert.equal(
    resolveViewModulePath('ai/index'),
    '../views/admin/ai/index.vue',
  )
})

test('maps dynamic admin module views into admin directory', () => {
  assert.equal(
    resolveViewModulePath('system/user'),
    '../views/admin/system/user/index.vue',
  )
  assert.equal(
    resolveViewModulePath('system/dept/index'),
    '../views/admin/system/dept/index.vue',
  )
  assert.equal(
    resolveViewModulePath('monitor/login-log'),
    '../views/admin/monitor/login-log/index.vue',
  )
})

test('returns null for empty component paths', () => {
  assert.equal(resolveViewModulePath(''), null)
  assert.equal(resolveViewModulePath('/'), null)
})
