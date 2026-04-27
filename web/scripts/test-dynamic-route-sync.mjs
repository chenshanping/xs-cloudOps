import assert from 'node:assert/strict'
import { computeDynamicRouteSyncPlan } from '../src/router/dynamic-route-sync.js'

assert.deepEqual(
  computeDynamicRouteSyncPlan(['Route_1'], ['Route_1', 'Route_2']),
  {
    removeNames: [],
    replaceNames: ['Route_1'],
    addNames: ['Route_2'],
  },
)

assert.deepEqual(
  computeDynamicRouteSyncPlan(['Route_1', 'Route_2'], ['Route_2']),
  {
    removeNames: ['Route_1'],
    replaceNames: ['Route_2'],
    addNames: [],
  },
)

assert.deepEqual(
  computeDynamicRouteSyncPlan(['Route_1'], ['Route_1']),
  {
    removeNames: [],
    replaceNames: ['Route_1'],
    addNames: [],
  },
)

console.log('dynamic route sync tests passed')
