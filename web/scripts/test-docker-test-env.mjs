import assert from 'node:assert/strict'
import { existsSync, readFileSync } from 'node:fs'

function readRelativeFile(relativePath) {
  const fileUrl = new URL(relativePath, import.meta.url)

  assert.ok(existsSync(fileUrl), `${relativePath} must exist`)

  return readFileSync(fileUrl, 'utf8')
}

const envTest = readRelativeFile('../.env.test')
const nginxConf = readRelativeFile('../nginx.conf')
const dockerfile = readRelativeFile('../Dockerfile')

assert.match(envTest, /^VITE_API_BASE_URL=\/api$/m, '.env.test must set VITE_API_BASE_URL=/api')
assert.match(
  nginxConf,
  /location\s+\/api\/\s*\{[\s\S]*proxy_pass\s+http:\/\/server:9000\s*;/,
  'nginx.conf must preserve the /api prefix by proxying to http://server:9000 without a URI suffix',
)
assert.doesNotMatch(
  nginxConf,
  /location\s+\/api\/\s*\{[\s\S]*proxy_pass\s+http:\/\/server:9000\/\s*;/,
  'nginx.conf must not use a trailing slash in proxy_pass for /api/',
)
assert.match(dockerfile, /^FROM\s+node:.*\s+AS\s+builder$/m, 'Dockerfile must define a node builder stage')
assert.match(dockerfile, /^FROM\s+nginx:.*$/m, 'Dockerfile must define an nginx runtime stage')
assert.match(dockerfile, /^RUN\s+npm run build:test$/m, 'Dockerfile must run npm run build:test')
assert.match(
  dockerfile,
  /^COPY\s+--from=builder\s+\/app\/dist\s+\/usr\/share\/nginx\/html$/m,
  'Dockerfile must copy built assets from the builder stage into nginx html root',
)

console.log('docker test env smoke checks passed')
