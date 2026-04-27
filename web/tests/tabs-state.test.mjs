import test from 'node:test'
import assert from 'node:assert/strict'
import {
  HOME_TAB,
  closeLeftTabs,
  closeOtherTabs,
  closeRightTabs,
  closeTab,
  createInitialTabsState,
  upsertTab,
} from '../dist-tests-tabs-state/src/store/tabs-state.js'

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

test('createInitialTabsState always preserves the home tab', () => {
  const state = createInitialTabsState([])
  assert.deepEqual(state.tabs, [HOME_TAB])
  assert.equal(state.activeKey, HOME_TAB.key)
})

test('upsertTab appends a new tab and updates active key', () => {
  const state = upsertTab(createInitialTabsState([]), reportTab)
  assert.deepEqual(state.tabs.map((tab) => tab.key), [HOME_TAB.key, reportTab.key])
  assert.equal(state.activeKey, reportTab.key)
})

test('upsertTab de-duplicates by fullPath key', () => {
  const first = upsertTab(createInitialTabsState([]), reportTab)
  const second = upsertTab(first, { ...reportTab, title: '报表更新' })
  assert.equal(second.tabs.length, 2)
  assert.equal(second.tabs[1].title, '报表更新')
})

test('closeTab chooses the right neighbor, then left neighbor, then home', () => {
  const state = {
    tabs: [HOME_TAB, reportTab, profileTab],
    activeKey: reportTab.key,
  }
  const next = closeTab(state, reportTab.key)
  assert.deepEqual(next.tabs.map((tab) => tab.key), [HOME_TAB.key, profileTab.key])
  assert.equal(next.activeKey, profileTab.key)
})

test('closeLeftTabs keeps the current tab and affix tabs', () => {
  const state = {
    tabs: [HOME_TAB, reportTab, profileTab],
    activeKey: profileTab.key,
  }
  const next = closeLeftTabs(state, profileTab.key)
  assert.deepEqual(next.tabs.map((tab) => tab.key), [HOME_TAB.key, profileTab.key])
})

test('closeRightTabs keeps the current tab and affix tabs', () => {
  const state = {
    tabs: [HOME_TAB, reportTab, profileTab],
    activeKey: reportTab.key,
  }
  const next = closeRightTabs(state, reportTab.key)
  assert.deepEqual(next.tabs.map((tab) => tab.key), [HOME_TAB.key, reportTab.key])
})

test('closeOtherTabs keeps the current tab and affix tabs only', () => {
  const state = {
    tabs: [HOME_TAB, reportTab, profileTab],
    activeKey: profileTab.key,
  }
  const next = closeOtherTabs(state, profileTab.key)
  assert.deepEqual(next.tabs.map((tab) => tab.key), [HOME_TAB.key, profileTab.key])
  assert.equal(next.activeKey, profileTab.key)
})
