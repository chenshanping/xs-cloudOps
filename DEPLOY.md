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
│   ├── config-test.yaml
│   ├── rbac_model.conf
│   └── sql/         # 增量升级脚本（仅升级已有环境时按需执行）
└── web/             # 前端 Vite 项目
    ├── src/
    ├── index.html
    ├── package.json
    └── vite.config.ts
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
- `config-test.yaml` （测试环境配置文件，可按需改名）
- `rbac_model.conf` （Casbin 权限模型配置）
- `sql/` （可选，已有环境升级时按需执行的增量脚本）

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
├── config-test.yaml   # 测试环境配置文件
├── rbac_model.conf    # Casbin 权限模型
└── sql/               # 可选：已有环境升级时按需执行的增量脚本
```

#### 4.2 准备测试环境配置文件

可以直接复制仓库中的 `server/config-test.yaml` 到服务器，再按测试环境修改 MySQL 和 Redis 连接信息：

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

先在 MySQL 中创建业务库，例如 `go_rbac_admin`。

首次启动后端时，程序会自动完成以下动作：

- 自动建表和补齐缺失字段
- 自动写入内置菜单、API、配置等基础数据
- 当库里还没有用户时，自动创建默认管理员账号

因此，新测试环境不需要再手工导入旧的初始化 SQL。

如果你是在已有环境上升级，才需要按版本变更执行 `server/sql/` 目录中的对应增量脚本。

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

    # API 反向代理到后端，注意 proxy_pass 末尾不要加 /
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 6. 1Panel 部署后基础连通性检查

- 后端接口：`http://服务器IP:8080`
- 前端页面：`http://服务器IP` 或绑定的域名

---

## 二、Docker Compose 部署运行手册

当前仓库的 Docker 栈固定为：`mysql + redis + server + web`。

- `mysql`：MySQL 8.4，数据保存在命名卷 `mysql_data`
- `redis`：Redis 7.4，启用密码与 AOF，数据保存在命名卷 `redis_data`
- `server`：Go 后端，对宿主机暴露 `9000`
- `web`：Nginx 托管前端，对宿主机暴露 `8081`

### 1. 启动前确认

在仓库根目录执行以下命令：

```bash
docker compose up -d --build
```

首次启动说明：

- 首次启动会拉取镜像、构建前后端镜像，并初始化 MySQL/Redis，耗时会明显更长。
- `server` 依赖 `mysql` 和 `redis` 健康检查通过后才会启动，`web` 会等待 `server` 健康检查通过。
- 如果本机已经存在旧的 `mysql_data` 卷，新的 `MYSQL_USER`、`MYSQL_PASSWORD`、`MYSQL_DATABASE` 环境变量不会回填到旧数据目录；此时应按下文“重置并清空数据”流程清理旧卷后再重启。

### 2. 状态与日志

查看服务状态：

```bash
docker compose ps
```

查看全部日志：

```bash
docker compose logs --tail=100
```

查看单个服务日志：

```bash
docker compose logs --tail=100 mysql
docker compose logs --tail=100 redis
docker compose logs --tail=100 server
docker compose logs --tail=100 web
```

持续跟随日志：

```bash
docker compose logs -f server
docker compose logs -f web
```

### 3. 访问地址

- 宿主机本地检查：
  - 前端首页：`http://127.0.0.1:8081`
  - 后端 API 基础地址：`http://127.0.0.1:9000/api/v1`
  - 健康检查：`http://127.0.0.1:9000/api/v1/health`
- 外部测试机访问形式：
  - 前端首页：`http://<server-ip>:8081`
  - 健康检查：`http://<server-ip>:9000/api/v1/health`

说明：

- 前端通过同源 `/api/*` 反向代理到 `server:9000`，浏览器访问前端时无需额外改 API 地址。
- `mysql` 和 `redis` 未对宿主机开放端口，默认仅供 Compose 内部服务使用。

### 4. 默认管理员账号

首次初始化完成后，系统会自动创建默认管理员：

- 用户名：`admin`
- 密码：`123456`

首次登录后请立即修改默认密码。

### 5. 保留数据的重建流程

以下命令会重建镜像并重启服务，但保留 `mysql_data`、`redis_data`、`server/uploads`、`server/logs` 中的数据：

```bash
docker compose up -d --build
```

如果希望先停再起，也可以执行：

```bash
docker compose down
docker compose up -d --build
```

`docker compose down` 不会删除命名卷，因此数据库与 Redis 数据会保留。

### 6. 重置并清空数据

如果需要回到全新初始化状态，执行：

```bash
docker compose down -v
```

如需同时清空后端本地挂载目录中的上传文件与日志，再删除：

```bash
rm -rf server/uploads server/logs
```

```powershell
Remove-Item -Recurse -Force server\uploads, server\logs
```

说明：

- `docker compose down -v` 会删除 `mysql_data` 和 `redis_data`，下次启动会重新初始化数据库、Redis 与默认管理员数据。
- 旧卷导致账号密码不匹配、初始化数据不符合当前 compose 配置时，应优先使用此流程。

### 7. 常用运维命令

停止服务：

```bash
docker compose down
```

仅重启后端：

```bash
docker compose restart server
```

重新拉起后端并重建：

```bash
docker compose up -d --build server
docker compose restart web
```

说明：

- `web/nginx.conf` 当前通过 `proxy_pass http://server:9000;` 转发到后端
- 如果后端容器被重建，建议同时重启一次 `web`，避免 Nginx 仍持有旧的上游解析结果
- 如果前端资源或 Nginx 配置有变更，再单独执行 `docker compose up -d --build web`

### 8. 本工作树实际验证结果

以下结果来自本工作树 `codex/docker-test-deploy` 的实际执行，不代表所有环境都能直接复现。

实际尝试的命令：

```bash
cd server && go test ./...
cd ../web && npm run test:docker-test-env
cd ../web && npm run build:test
docker compose up -d --build
docker compose ps
docker compose logs --tail=100 server
curl http://127.0.0.1:9000/api/v1/health
```

已通过：

- `cd server && go test ./...`
- `cd ../web && npm run test:docker-test-env`
- `cd ../web && npm run build:test`

未通过：

- `docker compose up -d --build`
  - 失败原因是 Docker Hub 拉取 `mysql:8.4` 时超时，实际报错为：
  - `failed to resolve reference "docker.io/library/mysql:8.4": failed to do request: Head "https://registry-1.docker.io/v2/library/mysql/manifests/8.4"`
  - `dial tcp 159.138.20.20:443: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond`
- `docker compose ps`
  - 命令本身执行成功，但因为镜像未拉取完成，Compose 栈没有服务启动。
- `docker compose logs --tail=100 server`
  - 命令本身执行成功，但没有 `server` 容器日志可读。
- `curl http://127.0.0.1:9000/api/v1/health`
  - 本次没有得到 Compose 栈健康结果；在该次验证中 `docker compose up -d --build` 已失败，且宿主机 `9000` 端口当时被其他本地进程占用，返回的不是这套 Compose 服务的健康响应。

结论：

- 本工作树已经完成后端测试、前端 Docker 环境脚本检查、前端 test 构建验证。
- 端到端 Compose 启动、容器健康检查、以及基于该栈的 UI 登录验证未完成，原因是镜像未能从 Docker Hub 拉取成功。

---

## 三、常见问题

### 1. 端口冲突

如果端口被占用，修改 `docker-compose.yml` 的宿主机端口映射，或修改 Docker 镜像内使用的 `server/config.docker.yaml`。

### 2. 数据库连接失败

- 检查 MySQL 服务是否正常启动
- 检查配置文件中的数据库连接信息是否正确
- Docker 环境下，确保使用服务名（`mysql`）而非 `127.0.0.1`

### 3. 前端接口请求失败

- 检查前端打包时 API 地址配置是否正确
- 检查 Nginx 反向代理配置是否生效
- 检查后端服务是否正常运行
