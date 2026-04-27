import { execFileSync } from 'node:child_process'

const nodeCommand = process.execPath

execFileSync(nodeCommand, ['--test', 'tests/config-component-imports.test.mjs'], {
  stdio: 'inherit',
})
