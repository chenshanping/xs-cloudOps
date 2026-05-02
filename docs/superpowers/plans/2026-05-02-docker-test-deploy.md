# Docker Test Deploy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Land a reproducible single-machine Docker test deployment for `mysql + redis + server + web`, with first-start auto initialization, `/api` reverse proxying, persistent uploads/logs, and clear runbook docs.

**Architecture:** Keep application initialization inside the existing Go service, but make deployment-specific concerns explicit. The backend gains a selectable config file and a minimal health endpoint; the frontend test build switches to same-origin `/api`; `docker compose` becomes the source of truth for service topology, health checks, and mounted persistence.

**Tech Stack:** Go, Gin, Gorm, MySQL 8, Redis 7, Vue 3, Vite, Nginx, Docker, Docker Compose

---

### Task 1: Add Backend Config Selection and Health Endpoint

**Files:**
- Create: `server/api/v1/health.go`
- Create: `server/router/modules/health.go`
- Create: `server/initialize/config_test.go`
- Modify: `server/main.go`
- Modify: `server/initialize/config.go`
- Test: `server/initialize/config_test.go`

- [ ] **Step 1: Write the failing backend config-loading test**

```go
package initialize

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromPath(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.test.yaml")

	content := []byte(`
server:
  port: 9011
  mode: release
  host: 0.0.0.0:9011
mysql:
  host: mysql
  port: 3306
  username: tester
  password: secret
  dbname: go-base
  charset: utf8mb4
  max_idle_conns: 5
  max_open_conns: 10
redis:
  host: redis
  port: 6379
  password: ""
  db: 1
jwt:
  secret: demo
  expires: 7200
  refresh_window: 604800
  issuer: web
casbin:
  model_path: ./rbac_model.conf
log:
  level: info
  format: console
  directory: ./logs
`)

	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, _, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Server.Port != 9011 {
		t.Fatalf("unexpected port: %d", cfg.Server.Port)
	}
	if cfg.MySQL.Host != "mysql" {
		t.Fatalf("unexpected mysql host: %s", cfg.MySQL.Host)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./initialize -run TestLoadConfigFromPath`

Expected: FAIL because `initialize.InitConfig()` currently hardcodes `config.yaml` and does not expose a reusable loader.

- [ ] **Step 3: Implement config-path loading, CLI selection, and public health handler**

```go
// server/initialize/config.go
func LoadConfig(configPath string) (*config.Config, *viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, v, nil
}

func InitConfig(configPath string) {
	cfg, v, err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	global.Config = cfg
	global.Viper = v
}
```

```go
// server/main.go
func main() {
	configPath := flag.String("c", "config.yaml", "config file path")
	flag.Parse()

	initialize.InitConfig(*configPath)
	// keep the rest unchanged
}
```

```go
// server/api/v1/health.go
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthApi struct{}

func (a *HealthApi) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

var Health = new(HealthApi)
```

```go
// server/router/modules/health.go
package modules

import (
	v1 "server/api/v1"
)

func init() {
	RegisterModule(&HealthModule{})
}

type HealthModule struct{}

func (m *HealthModule) Name() string {
	return "基础服务"
}

func (m *HealthModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", v1.Health.GetHealth)
}

func (m *HealthModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {}
```

- [ ] **Step 4: Run tests to verify it passes**

Run:

```bash
go test ./initialize -run TestLoadConfigFromPath
go test ./... 
```

Expected:
- config loader test passes
- full backend test suite remains green

- [ ] **Step 5: Commit**

```bash
git add server/main.go server/initialize/config.go server/initialize/config_test.go server/api/v1/health.go server/router/modules/health.go
git commit -m "feat: support docker config selection and health check"
```

### Task 2: Make Frontend Test Build Use Same-Origin `/api` and Ship Nginx Runtime

**Files:**
- Create: `web/Dockerfile`
- Create: `web/nginx.conf`
- Create: `web/.dockerignore`
- Create: `web/scripts/test-docker-test-env.mjs`
- Modify: `web/.env.test`
- Modify: `web/package.json`
- Test: `web/scripts/test-docker-test-env.mjs`

- [ ] **Step 1: Write the failing frontend deployment-smoke test**

```js
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import test from 'node:test'

test('test env uses same-origin api base', () => {
  const envFile = readFileSync(resolve(process.cwd(), '.env.test'), 'utf8')
  assert.match(envFile, /VITE_API_BASE_URL=\/api/)
})

test('nginx proxies api traffic to backend container', () => {
  const nginxFile = readFileSync(resolve(process.cwd(), 'nginx.conf'), 'utf8')
  assert.match(nginxFile, /location \/api\//)
  assert.match(nginxFile, /proxy_pass http:\/\/server:9000\//)
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `node scripts/test-docker-test-env.mjs`

Expected: FAIL because `.env.test` still points to a hardcoded external IP and `web/nginx.conf` does not exist yet.

- [ ] **Step 3: Implement same-origin test env, Nginx proxy, and frontend Docker image**

```env
# web/.env.test
VITE_APP_ENV=test
VITE_API_BASE_URL=/api
```

```nginx
# web/nginx.conf
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://server:9000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```dockerfile
# web/Dockerfile
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build:test

FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

```dockerignore
# web/.dockerignore
node_modules
dist
.git
```

```json
// web/package.json
{
  "scripts": {
    "test:docker-test-env": "node scripts/test-docker-test-env.mjs"
  }
}
```

- [ ] **Step 4: Run tests to verify it passes**

Run:

```bash
npm run test:docker-test-env
npm run build:test
```

Expected:
- deployment-smoke test passes
- test-mode frontend build succeeds

- [ ] **Step 5: Commit**

```bash
git add web/.env.test web/Dockerfile web/nginx.conf web/.dockerignore web/scripts/test-docker-test-env.mjs web/package.json
git commit -m "feat: add docker-ready frontend test build"
```

### Task 3: Add Backend Docker Assets and Root Compose Orchestration

**Files:**
- Create: `server/Dockerfile`
- Create: `server/.dockerignore`
- Create: `server/config.docker.yaml`
- Create: `docker-compose.yml`
- Test: `docker-compose.yml`

- [ ] **Step 1: Write the failing compose-structure smoke test**

```yaml
# target expectations for docker-compose.yml
services:
  mysql:
  redis:
  server:
  web:
```

Create a simple validation script:

```powershell
$content = Get-Content -Raw docker-compose.yml
if ($content -notmatch 'services:' -or $content -notmatch 'mysql:' -or $content -notmatch 'server:' -or $content -notmatch 'web:') {
  throw 'docker-compose.yml missing required services'
}
if ($content -notmatch 'condition:\s*service_healthy') {
  throw 'docker-compose.yml missing health-based startup ordering'
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```powershell
powershell -Command "$content = Get-Content -Raw docker-compose.yml; if ($content -notmatch 'services:' -or $content -notmatch 'mysql:' -or $content -notmatch 'server:' -or $content -notmatch 'web:') { throw 'docker-compose.yml missing required services' }; if ($content -notmatch 'condition:\s*service_healthy') { throw 'docker-compose.yml missing health-based startup ordering' }"
```

Expected: FAIL because the root `docker-compose.yml` does not exist yet.

- [ ] **Step 3: Implement backend Docker build, docker config, and service orchestration**

```dockerfile
# server/Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl && \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /app/server /app/server
COPY config.docker.yaml /app/config.docker.yaml
COPY rbac_model.conf /app/rbac_model.conf

EXPOSE 9000
CMD ["./server", "-c", "config.docker.yaml"]
```

```yaml
# server/config.docker.yaml
server:
  port: 9000
  mode: release
  host: 0.0.0.0:9000
mysql:
  host: mysql
  port: 3306
  username: go_base_user
  password: go_base_pass
  dbname: go-base
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100
redis:
  host: redis
  port: 6379
  password: ""
  db: 1
jwt:
  secret: go-base-jwt-secret-key
  expires: 7200
  refresh_window: 604800
  issuer: web
casbin:
  model_path: ./rbac_model.conf
log:
  level: info
  format: console
  directory: ./logs
```

```yaml
# docker-compose.yml
services:
  mysql:
    image: mysql:8.0
    container_name: go-base-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: root123456
      MYSQL_DATABASE: go-base
      MYSQL_USER: go_base_user
      MYSQL_PASSWORD: go_base_pass
    command: ["--default-authentication-plugin=mysql_native_password"]
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "127.0.0.1", "-proot123456"]
      interval: 10s
      timeout: 5s
      retries: 12
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    container_name: go-base-redis
    restart: unless-stopped
    command: ["redis-server", "--appendonly", "yes"]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 12
    volumes:
      - redis_data:/data

  server:
    build:
      context: ./server
      dockerfile: Dockerfile
    container_name: go-base-server
    restart: unless-stopped
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "9000:9000"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://127.0.0.1:9000/api/v1/health"]
      interval: 15s
      timeout: 5s
      retries: 10
    volumes:
      - ./server/uploads:/app/uploads
      - ./server/logs:/app/logs

  web:
    build:
      context: ./web
      dockerfile: Dockerfile
    container_name: go-base-web
    restart: unless-stopped
    depends_on:
      server:
        condition: service_healthy
    ports:
      - "8081:80"

volumes:
  mysql_data:
  redis_data:
```

```dockerignore
# server/.dockerignore
logs
uploads
.git
```

- [ ] **Step 4: Run validations to verify it passes**

Run:

```bash
docker compose config
```

Expected:
- compose file parses successfully
- service dependencies and health checks are visible in normalized output

If Docker is unavailable in the local session, record the blocker explicitly and still verify YAML structure with the scripted smoke check.

- [ ] **Step 5: Commit**

```bash
git add server/Dockerfile server/.dockerignore server/config.docker.yaml docker-compose.yml
git commit -m "feat: add docker test stack orchestration"
```

### Task 4: Update Deployment Runbook and Run End-to-End Verification

**Files:**
- Modify: `DEPLOY.md`
- Test: backend, frontend, and docker verification commands

- [ ] **Step 1: Update the deployment document with the real test-environment workflow**

Replace the old example-heavy Docker section with the actual project runbook:

```md
## Docker 测试环境部署

### 启动
docker compose up -d --build

### 查看状态
docker compose ps

### 查看日志
docker compose logs -f server
docker compose logs -f web

### 访问地址
- 前端：http://服务器IP:8081
- 后端健康检查：http://服务器IP:9000/api/v1/health

### 重建但保留数据
docker compose down
docker compose up -d --build

### 清空数据重置
docker compose down -v
docker compose up -d --build
```

- [ ] **Step 2: Run code-level verification before docker smoke**

Run:

```bash
cd server && go test ./...
cd ../web && npm run test:docker-test-env
cd ../web && npm run build:test
```

Expected:
- backend tests pass
- frontend docker env smoke test passes
- frontend test build passes

- [ ] **Step 3: Run deployment verification**

Run:

```bash
docker compose up -d --build
docker compose ps
docker compose logs --tail=100 server
```

Expected:
- `mysql`, `redis`, `server`, `web` all show running
- `server` no longer reports config-path or DB connection errors
- first start auto-initializes schema and default data

- [ ] **Step 4: Run the user-facing smoke test**

Run:

```bash
curl http://127.0.0.1:9000/api/v1/health
```

Expected:
- returns HTTP 200 with `{"status":"ok"}`

Then open `http://127.0.0.1:8081` in a browser and verify:
- login page loads
- default `admin / 123456` can sign in
- frontend API traffic goes through `/api`

- [ ] **Step 5: Commit**

```bash
git add DEPLOY.md
git commit -m "docs: finalize docker test deployment runbook"
```
