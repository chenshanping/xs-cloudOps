# go-base 部署说明

本文档包含两种部署方式：

1. 测试环境：通过 1Panel 直接部署前后端（非 Docker）
2. Docker + docker-compose 一键部署方案

---

## 一、测试环境通过 1Panel 部署（非 Docker）

### 1. 项目结构

```
go-base/
├── server/          # 后端 Go 服务
│   ├── main.go
│   ├── go.mod
│   ├── config.yaml
│   ├── rbac_model.conf
└── web/             # 前端 Vite 项目
│   ├── src/
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
└── go_rbac_admin.sql   # 数据库初始化脚本
```

### 2. 后端打包步骤（本地）

```bash
# 进入后端目录
cd server

# 交叉编译为 Linux 可执行文件（如果服务器是 Linux）
# Windows PowerShell:
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o server main.go

# Linux / macOS:
GOOS=linux GOARCH=amd64 go build -o server main.go
```

需要上传到服务器的文件：

- `server` （编译后的可执行文件）
- `config.yaml` （配置文件）
- `rbac_model.conf` （Casbin 权限模型配置）
- `sql/` （数据库初始化脚本目录）

### 3. 前端打包步骤（本地）

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 生产环境打包
npm run build:test
```

打包完成后，`dist/` 目录下的所有文件即为需要上传的静态资源。

### 4. 在 1Panel 上部署后端

#### 4.1 上传文件

将后端文件上传到服务器，例如：`/opt/go-base/server/`

```
/opt/go-base/server/
├── server          # 编译好的二进制文件
├── config-test.yaml        # 测试环境配置文件
├── rbac_model.conf         # Casbin 权限模型
└── sql/                    # 数据库初始化脚本（可选）
    ├── product.sql
    └── product_type.sql
```

#### 4.2 准备测试环境配置文件

创建 `config-test.yaml`，配置测试环境的 MySQL 和 Redis 连接信息：

```yaml
server:
  port: 8080
  mode: release
  host: 0.0.0.0:8080

mysql:
  host: 你的MySQL地址
  port: 3306
  username: root
  password: 你的密码
  dbname: go_rbac_admin
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: 你的Redis地址
  port: 6379
  password: "你的密码"
  db: 1

jwt:
  secret: your-jwt-secret
  expires: 7200
  issuer: server

casbin:
  model_path: ./rbac_model.conf

log:
  level: info
  format: console
  directory: ./logs
```

#### 4.3 初始化数据库

通过 1Panel 数据库管理或命令行执行 SQL 脚本：

```bash
mysql -h MySQL地址 -P 3306 -u root -p数据库密码 go_rbac_admin < /opt/go-base/server/sql/product.sql
mysql -h MySQL地址 -P 3306 -u root -p数据库密码 go_rbac_admin < /opt/go-base/server/sql/product_type.sql
```

#### 4.4 通过 1Panel 运行环境部署

1. 登录 1Panel 后台
2. 进入 **网站** → **运行环境** → **Go**
3. 创建运行环境，配置如下：

| 配置项 | 值 |
|--------|----|
| 名称 | server |
| 运行目录 | /opt/go-base/server |
| 启动命令 | `./server -c config-test.yaml` |
| 端口 | 8080 |

4. 确保二进制文件有执行权限：

```bash
chmod +x /opt/go-base/server/server
```

5. 启动服务

#### 4.5 开放端口

确保防火墙或安全组已开放 `8080` 端口。

### 5. 在 1Panel 上部署前端

#### 5.1 上传文件

将前端 `dist/` 目录内容上传到服务器，例如：`/opt/go-base/web/`

#### 5.2 创建 Nginx 站点

在 1Panel「网站管理」中新建静态站点：

- **站点目录**：`/opt/go-base/web`
- **绑定域名**：你的测试域名（或使用 IP + 端口）

#### 5.3 Nginx 配置示例（可选）

如需配置 SPA 路由和 API 反向代理：

```nginx
server {
    listen 80;
    server_name test.example.com;

    root /opt/go-base/web;
    index index.html;

    # SPA 路由支持
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 反向代理到后端
    location /api/ {
        proxy_pass http://127.0.0.1:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### 6. 验证部署

- 后端接口：`http://服务器IP:8080`
- 前端页面：`http://服务器IP` 或绑定的域名

---

## 二、Docker 部署方案

### 1. 后端 Dockerfile

在 `server/` 目录下创建 `Dockerfile`：

```dockerfile
# 构建阶段
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

# 运行阶段
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /app/server /app/server
COPY config.docker.yaml /app/config.docker.yaml
COPY rbac_model.conf /app/rbac_model.conf

EXPOSE 8080

CMD ["./server", "-c", "config.docker.yaml"]
```

### 2. 后端 Docker 配置文件

在 `server/` 目录下创建 `config.docker.yaml`：

```yaml
server:
  port: 8080
  mode: release
  host: 0.0.0.0:8080

mysql:
  host: mysql
  port: 3306
  username: go_base_user
  password: go_base_pass
  dbname: go_base
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: redis
  port: 6379
  password: ""
  db: 1

jwt:
  secret: your-jwt-secret
  expires: 7200
  issuer: server

casbin:
  model_path: ./rbac_model.conf

log:
  level: info
  format: console
  directory: ./logs
```

### 3. 前端 Dockerfile

在 `web/` 目录下创建 `Dockerfile`：

```dockerfile
# 构建阶段
FROM node:20-alpine AS build

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# 运行阶段
FROM nginx:1.27-alpine

COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### 4. 前端 Nginx 配置

在 `web/` 目录下创建 `nginx.conf`：

```nginx
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    # SPA 路由支持
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 反向代理
    location /api/ {
        proxy_pass http://server:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### 5. docker-compose.yml

在项目根目录 `go-base/` 下创建 `docker-compose.yml`：

```yaml
version: "3.9"

services:
  mysql:
    image: mysql:8.0
    container_name: go-base-mysql
    environment:
      MYSQL_ROOT_PASSWORD: root123456
      MYSQL_DATABASE: go_base
      MYSQL_USER: go_base_user
      MYSQL_PASSWORD: go_base_pass
    command: ["--default-authentication-plugin=mysql_native_password"]
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./server/sql:/docker-entrypoint-initdb.d
    networks:
      - go-base-net

  redis:
    image: redis:7-alpine
    container_name: go-base-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - go-base-net

  server:
    build:
      context: ./server
      dockerfile: Dockerfile
    container_name: server
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis
    networks:
      - go-base-net

  web:
    build:
      context: ./web
      dockerfile: Dockerfile
    container_name: web
    ports:
      - "8081:80"
    depends_on:
      - server
    networks:
      - go-base-net

networks:
  go-base-net:
    driver: bridge

volumes:
  mysql_data:
  redis_data:
```

### 6. Docker 部署命令

```bash
# 构建并启动所有服务
docker-compose build
docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看单个服务日志
docker-compose logs -f server
docker-compose logs -f web

# 停止并移除容器
docker-compose down

# 停止并移除容器及数据卷（会清空数据库）
docker-compose down -v
```

### 7. Docker 部署访问地址

- 后端接口：`http://localhost:8080`
- 前端页面：`http://localhost:8081`

---

## 三、常见问题

### 1. 端口冲突

如果端口被占用，修改 `docker-compose.yml` 或 `config.yaml` 中的端口映射。

### 2. 数据库连接失败

- 检查 MySQL 服务是否正常启动
- 检查配置文件中的数据库连接信息是否正确
- Docker 环境下，确保使用服务名（`mysql`）而非 `127.0.0.1`

### 3. 前端接口请求失败

- 检查前端打包时 API 地址配置是否正确
- 检查 Nginx 反向代理配置是否生效
- 检查后端服务是否正常运行
