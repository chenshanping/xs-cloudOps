import assert from 'node:assert/strict'
import { shouldRestoreMigrationTaskOnOpen } from '../src/views/admin/system/file/components/file-migration-task-state.js'

assert.equal(shouldRestoreMigrationTaskOnOpen(null), false)
assert.equal(shouldRestoreMigrationTaskOnOpen({ status: 'SCANNING' }), true)
assert.equal(shouldRestoreMigrationTaskOnOpen({ status: 'RUNNING' }), true)
assert.equal(shouldRestoreMigrationTaskOnOpen({ status: 'SUCCESS' }), false)
assert.equal(shouldRestoreMigrationTaskOnOpen({ status: 'FAILED' }), false)

console.log('file migration task restore tests passed')
