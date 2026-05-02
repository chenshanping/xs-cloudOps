import assert from 'node:assert/strict'
import { existsSync, readFileSync } from 'node:fs'

function readRelativeFile(relativePath) {
  const fileUrl = new URL(relativePath, import.meta.url)

  assert.ok(existsSync(fileUrl), `${relativePath} must exist`)

  return readFileSync(fileUrl, 'utf8')
}

function extractBlock(source, blockHeader) {
  const headerIndex = source.indexOf(blockHeader)
  assert.notEqual(headerIndex, -1, `Missing block: ${blockHeader}`)

  const openBraceIndex = source.indexOf('{', headerIndex)
  assert.notEqual(openBraceIndex, -1, `Missing opening brace for block: ${blockHeader}`)

  let depth = 0
  for (let index = openBraceIndex; index < source.length; index += 1) {
    const character = source[index]
    if (character === '{') {
      depth += 1
    } else if (character === '}') {
      depth -= 1
      if (depth === 0) {
        return source.slice(openBraceIndex + 1, index)
      }
    }
  }

  throw new Error(`Unclosed block: ${blockHeader}`)
}

const envTest = readRelativeFile('../.env.test')
const nginxConf = readRelativeFile('../nginx.conf')
const dockerfile = readRelativeFile('../Dockerfile')
const requestSource = readRelativeFile('../src/utils/request.ts')
const apiLocationBlock = extractBlock(nginxConf, 'location /api/')

assert.match(envTest, /^VITE_API_BASE_URL=\/api$/m, '.env.test must set VITE_API_BASE_URL=/api')
assert.match(
  requestSource,
  /baseURL:\s*'\/api\/v1'/,
  'request.ts must keep the deployed client on a same-origin /api/v1 baseURL',
)
assert.doesNotMatch(
  requestSource,
  /baseURL:\s*['"]https?:\/\//,
  'request.ts must not hardcode an absolute http(s) API host',
)
assert.match(
  apiLocationBlock,
  /proxy_pass\s+http:\/\/server:9000\s*;/,
  'the /api/ location block must preserve the /api prefix by proxying to http://server:9000 without a URI suffix',
)
assert.doesNotMatch(
  apiLocationBlock,
  /proxy_pass\s+http:\/\/server:9000\/\s*;/,
  'the /api/ location block must not use a trailing slash in proxy_pass',
)
assert.match(apiLocationBlock, /proxy_set_header\s+Host\s+\$host\s*;/, 'the /api/ location block must forward Host')
assert.match(apiLocationBlock, /proxy_set_header\s+X-Real-IP\s+\$remote_addr\s*;/, 'the /api/ location block must forward X-Real-IP')
assert.match(
  apiLocationBlock,
  /proxy_set_header\s+X-Forwarded-For\s+\$proxy_add_x_forwarded_for\s*;/,
  'the /api/ location block must forward X-Forwarded-For',
)
assert.match(
  apiLocationBlock,
  /proxy_set_header\s+X-Forwarded-Proto\s+\$scheme\s*;/,
  'the /api/ location block must forward X-Forwarded-Proto',
)
assert.match(dockerfile, /^FROM\s+node:.*\s+AS\s+builder$/m, 'Dockerfile must define a node builder stage')
assert.match(dockerfile, /^FROM\s+nginx:.*$/m, 'Dockerfile must define an nginx runtime stage')
assert.match(dockerfile, /^RUN\s+npm run build:test$/m, 'Dockerfile must run npm run build:test')
assert.match(
  dockerfile,
  /^COPY\s+nginx\.conf\s+\/etc\/nginx\/nginx\.conf$/m,
  'Dockerfile must copy the custom nginx.conf into the runtime image',
)
assert.match(
  dockerfile,
  /^COPY\s+--from=builder\s+\/app\/dist\s+\/usr\/share\/nginx\/html$/m,
  'Dockerfile must copy built assets from the builder stage into nginx html root',
)

console.log('docker test env smoke checks passed')
