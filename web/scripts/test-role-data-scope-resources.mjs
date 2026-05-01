import { execFileSync } from 'node:child_process'

const nodeCommand = process.execPath

execFileSync(nodeCommand, ['--test', 'tests/role-data-scope-resources.test.mjs'], {
  stdio: 'inherit',
})
