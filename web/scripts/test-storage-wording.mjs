import { execFileSync } from 'node:child_process'

const nodeCommand = process.execPath

execFileSync(nodeCommand, ['--test', 'tests/storage-wording.test.mjs'], {
  stdio: 'inherit',
})
