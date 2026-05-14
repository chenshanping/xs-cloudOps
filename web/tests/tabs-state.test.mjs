import test from 'node:test'
import assert from 'node:assert/strict'
import {
  closeAllTabs,
  closeLeftTabs,
  closeOtherTabs,
  closeRightTabs,
  closeTab,
  createInitialTabsState,
  upsertTab,
} from '../dist-tests-tabs-state/src/store/tabs-state.js'

const homeTab = {
  key: '/dashboard',
  path: '/dashboard',
  fullPath: '/dashboard',
  title: '首页',
  name: 'Route_1',
  affix: true,
}

const reportTab = {
  key: '/system/report?page=1',
  path: '/system/report',
  fullPath: '/system/report?page=1',
  title: '报表',
  name: 'ReportPage',
  affix: false,
}

const profileTab = {
  key: '/profile',
  path: '/profile',
  fullPath: '/profile',
  title: '个人中心',
  name: 'Profile',
  affix: false,
}

test('createInitialTabsState starts empty until a route opens a tab', () => {
  const state = createInitialTabsState([])
  assert.deepEqual(state.tabs, [])
  assert.equal(state.activeKey, '')
})

test('upsertTab appends a new tab and updates active key', () => {
  const state = upsertTab(createInitialTabsState([]), reportTab)
  assert.deepEqual(state.tabs.map((tab) => tab.key), [reportTab.key])
  assert.equal(state.activeKey, reportTab.key)
})

test('upsertTab de-duplicates by fullPath key', () => {
  const first = upsertTab(createInitialTabsState([]), reportTab)
  const second = upsertTab(first, { ...reportTab, title: '报表更新' })
  assert.equal(second.tabs.length, 1)
  assert.equal(second.tabs[0].title, '报表更新')
})

test('closeTab chooses the right neighbor, then left neighbor, then home', () => {
  const state = {
    tabs: [homeTab, reportTab, profileTab],
    activeKey: reportTab.key,
  }
  const next = closeTab(state, reportTab.key)
  assert.deepEqual(next.tabs.map((tab) => tab.key), [homeTab.key, profileTab.key])
  assert.equal(next.activeKey, profileTab.key)
})

test('closeLeftTabs keeps the current tab and affix tabs', () => {
  const state = {
    tabs: [homeTab, reportTab, profileTab],
    activeKey: profileTab.key,
  }
  const next = closeLeftTabs(state, profileTab.key)
  assert.deepEqual(next.tabs.map((tab) => tab.key), [homeTab.key, profileTab.key])
})

test('closeRightTabs keeps the current tab and affix tabs', () => {
  const state = {
    tabs: [homeTab, reportTab, profileTab],
    activeKey: reportTab.key,
  }
  const next = closeRightTabs(state, reportTab.key)
  assert.deepEqual(next.tabs.map((tab) => tab.key), [homeTab.key, reportTab.key])
})

test('closeOtherTabs keeps the current tab and affix tabs only', () => {
  const state = {
    tabs: [homeTab, reportTab, profileTab],
    activeKey: profileTab.key,
  }
  const next = closeOtherTabs(state, profileTab.key)
  assert.deepEqual(next.tabs.map((tab) => tab.key), [homeTab.key, profileTab.key])
  assert.equal(next.activeKey, profileTab.key)
})

test('closeAllTabs clears static tabs and lets the current route reopen the real home tab', () => {
  const next = closeAllTabs()
  assert.deepEqual(next.tabs, [])
  assert.equal(next.activeKey, '')
})
