# Go Base

基于 Go + Gin + JWT + GORM + Casbin + Vue 3 + TypeScript + Ant Design Vue 的企业级后台管理系统

> 本仓库已接入项目级 Codex + OpenSpec + Superpowers 工作流。开始非 trivial 任务前，请先阅读 [AGENTS.md](AGENTS.md) 和 [README_AGENT.md](README_AGENT.md)。

## 技术栈

### 后端
- Go 1.24+
- Gin (Web 框架)
- GORM (ORM 框架)
- JWT (身份认证 + 自动续期)
- Casbin (RBAC 权限控制)
- MySQL 8.0+
- Redis
- Zap + Lumberjack (结构化日志 + 文件轮转，ELK 就绪)
- Viper (配置管理)

### 前端
- Vue 3
- TypeScript
- Ant Design Vue
- Vite
- Pinia (状态管理)
- Axios
- Vue Router

## 功能特性

### 基础权限
- ✅ 用户管理：增删改查、状态控制、角色分配、批量操作
- ✅ 角色管理：增删改查、菜单 / API 权限分配、超级管理员标识
- ✅ 部门管理：树形部门结构、数据权限范围
- ✅ 菜单管理：树形菜单、按钮级权限
- ✅ API 管理：接口注册与权限控制
- ✅ JWT 认证：Token 自动续期、多终端管理
- ✅ Casbin 权限：基于 RBAC 的细粒度权限控制

### 系统功能
- ✅ 字典管理：数据字典的类型与条目管理
- ✅ 系统配置：Logo、登录页、邮箱、安全策略等在线配置
- ✅ 文件管理：上传 / 预览 / 删除，支持多存储后端
- ✅ 文件存储：本地存储、阿里云 OSS、腾讯云 COS、MinIO
- ✅ 文件迁移：跨存储类型迁移、同类型跨桶迁移、迁移历史记录
- ✅ 操作日志：记录用户操作日志
- ✅ 登录日志：记录用户登录日志
- ✅ 结构化日志：Zap JSON 日志 + Lumberjack 轮转，支持 Filebeat/ELK 采集
- ✅ 验证码：数字验证码、滑动验证码
- ✅ 邮件服务：SMTP 邮件发送

### AI 功能
- ✅ AI 对话：多模型对话、流式输出
- ✅ AI 配置：多 Provider / Model 管理、远程模型导入
- ✅ 可视化图表：ECharts 数据可视化

### 登录 / 注册
- ✅ 登录页定制：背景图、标语、特性列表等可视化配置
- ✅ 用户注册：可开关的注册功能、邮箱验证

## 快速开始

### 前置要求

- Go 1.24+
- Node.js 18+
- MySQL 8.0+
- Redis

### 后端启动

```bash
cd server

# 安装依赖
go mod download

# 修改配置
# 编辑 config.yaml 配置MySQL、Redis连接和日志参数

# 启动服务
go run main.go
```

后端默认运行在 `http://localhost:8080`

### 前端启动

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端默认运行在 `http://localhost:3000`

## 默认账号

- 用户名：`admin`
- 密码：`123456`

## 项目结构

```
go-base/
├── server/                 # 后端代码
│   ├── api/v1/            # API 处理层
│   ├── config/            # 配置结构
│   ├── global/            # 全局变量
│   ├── initialize/        # 初始化 (DB、路由、日志等)
│   ├── middleware/         # 中间件 (JWT、Casbin、日志等)
│   ├── model/             # 数据模型 + 请求/响应结构
│   ├── router/modules/    # 路由模块
│   ├── service/           # 业务逻辑层
│   │   ├── ai/           # AI 对话服务
│   │   ├── auth/         # 认证服务
│   │   ├── file/         # 文件管理服务
│   │   ├── oss/          # 对象存储客户端
│   │   ├── storagesvc/   # 存储配置服务
│   │   ├── configsvc/    # 系统配置服务
│   │   └── ...           # 其他业务模块
│   ├── sql/               # 增量 SQL 升级脚本
│   ├── utils/             # 工具函数
│   ├── config.yaml        # 配置文件
│   ├── rbac_model.conf    # Casbin 模型
│   └── main.go            # 入口文件
├── web/                   # 前端代码
│   ├── src/
│   │   ├── api/          # API 请求
│   │   ├── components/   # 公共组件
│   │   ├── layouts/      # 布局组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # 状态管理 (Pinia)
│   │   ├── types/        # TypeScript 类型定义
│   │   ├── utils/        # 工具函数
│   │   ├── views/        # 页面组件
│   │   └── main.ts       # 入口文件
│   └── package.json
├── openspec/              # OpenSpec 变更管理
├── docs/                  # 文档与设计
├── go-base.sql            # 数据库基线快照
└── README.md
```

## 日志配置

日志同时输出到**文件（JSON 格式）**和**控制台**，文件日志可直接被 Filebeat 采集送入 ELK。

```yaml
log:
  level: info           # 日志级别: debug, info, warn, error
  format: console       # 控制台格式: console(人类可读), json
  directory: ./logs     # 日志文件目录
  filename: app.log     # 日志文件名
  max_size: 100         # 单文件最大 MB
  max_backups: 5        # 保留旧文件数量
  max_age: 30           # 保留天数
  compress: true        # 压缩旧日志
  stdout: true          # 同时输出到控制台（容器部署建议开启）
```

- **传统部署**：Filebeat 指向 `./logs/app.log`，文件始终为 JSON 格式
- **容器部署**：设 `format: json`，通过 Docker/K8s 日志驱动采集 stdout

## API 文档

API 基础路径：`/api/v1`

### 认证接口
- `POST /auth/login` - 登录
- `POST /auth/logout` - 登出
- `POST /auth/refresh` - 刷新 Token
- `GET /auth/userinfo` - 获取当前用户信息

### 用户管理
- `GET /users` - 用户列表
- `POST /users` - 创建用户
- `PUT /users/:id` - 更新用户
- `DELETE /users/:id` - 删除用户
- `PUT /users/:id/status` - 修改用户状态
- `PUT /users/:id/password` - 重置密码

### 角色管理
- `GET /roles` - 角色列表
- `POST /roles` - 创建角色
- `PUT /roles/:id` - 更新角色
- `DELETE /roles/:id` - 删除角色
- `PUT /roles/:id/menus` - 分配菜单权限
- `PUT /roles/:id/apis` - 分配 API 权限

### 部门管理
- `GET /depts` - 部门列表 (树形)
- `POST /depts` - 创建部门
- `PUT /depts/:id` - 更新部门
- `DELETE /depts/:id` - 删除部门

### 菜单管理
- `GET /menus` - 菜单列表 (树形)
- `POST /menus` - 创建菜单
- `PUT /menus/:id` - 更新菜单
- `DELETE /menus/:id` - 删除菜单

### API 管理
- `GET /apis` - API 列表
- `POST /apis` - 创建 API
- `PUT /apis/:id` - 更新 API
- `DELETE /apis/:id` - 删除 API

### 字典管理
- `GET /dict/types` - 字典类型列表
- `POST /dict/types` - 创建字典类型
- `GET /dict/data` - 字典条目列表
- `POST /dict/data` - 创建字典条目

### 文件管理
- `POST /file/upload` - 上传文件
- `GET /files` - 文件列表
- `DELETE /files/:id` - 删除文件
- `POST /file/migration/preview` - 迁移预览
- `POST /file/migration/execute` - 执行迁移
- `GET /file/migration/status` - 迁移状态

### 系统配置
- `GET /configs` - 配置列表
- `PUT /configs` - 批量更新配置
- `GET /configs/keys` - 按键获取配置 (公开)

### 日志管理
- `GET /logs/operation` - 操作日志列表
- `GET /logs/login` - 登录日志列表

### AI
- `POST /ai/chat` - AI 对话 (流式)
- `GET /ai/providers` - Provider 列表
- `GET /ai/models` - 模型列表

## Docker 部署

参考 [DEPLOY.md](DEPLOY.md) 和 `docker-compose.yml`。

## License

MIT
