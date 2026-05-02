import assert from 'node:assert/strict'
import { existsSync, readFileSync } from 'node:fs'

function readRelativeFile(relativePath) {
  const fileUrl = new URL(relativePath, import.meta.url)

  assert.ok(existsSync(fileUrl), `${relativePath} must exist`)

  return readFileSync(fileUrl, 'utf8')
}

const envTest = readRelativeFile('../.env.test')
const nginxConf = readRelativeFile('../nginx.conf')

assert.match(envTest, /^VITE_API_BASE_URL=\/api$/m, '.env.test must set VITE_API_BASE_URL=/api')
assert.match(
  nginxConf,
  /location\s+\/api\/\s*\{[\s\S]*proxy_pass\s+http:\/\/server:9000\/;/,
  'nginx.conf must proxy /api/ to http://server:9000/',
)

console.log('docker test env smoke checks passed')
