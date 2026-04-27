import { execFileSync } from 'node:child_process'
import { existsSync, rmSync } from 'node:fs'
import { fileURLToPath } from 'node:url'

const distDir = new URL('../dist-tests', import.meta.url)
const nodeCommand = process.execPath
const tscScript = fileURLToPath(
  new URL('../node_modules/typescript/bin/tsc', import.meta.url),
)

function cleanup() {
  if (existsSync(distDir)) {
    rmSync(distDir, { recursive: true, force: true })
  }
}

try {
  cleanup()

  execFileSync(
    nodeCommand,
    [
      tscScript,
      'src/router/view-resolver.ts',
      '--module',
      'ESNext',
      '--moduleResolution',
      'bundler',
      '--target',
      'ES2020',
      '--outDir',
      'dist-tests',
      '--rootDir',
      '.',
    ],
    { stdio: 'inherit' },
  )

  execFileSync(nodeCommand, ['--test', 'tests/view-resolver.test.mjs'], {
    stdio: 'inherit',
  })
} finally {
  cleanup()
}
