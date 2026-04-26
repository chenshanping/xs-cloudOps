# Go RBAC Admin

基于 Go + Gin + JWT + Gorm + Casbin + Vue3 + TypeScript + Ant Design Vue 的后台权限管理系统

> 本仓库已接入项目级 Codex + OpenSpec + Superpowers 工作流。开始非 trivial 任务前，请先阅读 [AGENTS.md](/E:/go_project/go-base/AGENTS.md) 和 [README_AGENT.md](/E:/go_project/go-base/README_AGENT.md)。

## 技术栈

### 后端
- Go 1.21+
- Gin (Web框架)
- GORM (ORM框架)
- JWT (身份认证)
- Casbin (权限控制)
- MySQL 8.0+
- Redis
- Zap (日志)
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

- ✅ 用户管理：用户的增删改查、状态控制、角色分配
- ✅ 角色管理：角色的增删改查、权限分配
- ✅ 菜单管理：菜单的树形管理
- ✅ API管理：API接口的管理
- ✅ JWT认证：基于JWT的身份认证
- ✅ Casbin权限：基于Casbin的RBAC权限控制
- ✅ 操作日志：记录用户操作日志
- ✅ 登录日志：记录用户登录日志

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis

### 后端启动

```bash
cd server

# 安装依赖
go mod download

# 修改配置
# 编辑 config.yaml 配置MySQL和Redis连接信息

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
go-rbac-admin/
├── server/                 # 后端代码
│   ├── api/               # API处理层
│   ├── config/            # 配置
│   ├── global/            # 全局变量
│   ├── initialize/        # 初始化
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── router/            # 路由
│   ├── service/           # 业务逻辑层
│   ├── utils/             # 工具函数
│   ├── config.yaml        # 配置文件
│   ├── rbac_model.conf    # Casbin模型
│   └── main.go            # 入口文件
├── web/                   # 前端代码
│   ├── src/
│   │   ├── api/          # API请求
│   │   ├── layouts/      # 布局组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # 状态管理
│   │   ├── types/        # 类型定义
│   │   ├── utils/        # 工具函数
│   │   ├── views/        # 页面组件
│   │   └── main.ts       # 入口文件
│   └── package.json
└── README.md
```

## API文档

API基础路径：`/api/v1`

### 认证接口

- `POST /auth/login` - 登录
- `POST /auth/logout` - 登出
- `POST /auth/refresh` - 刷新Token
- `GET /auth/userinfo` - 获取当前用户信息

### 用户管理

- `GET /users` - 用户列表
- `GET /users/:id` - 用户详情
- `POST /users` - 创建用户
- `PUT /users/:id` - 更新用户
- `DELETE /users/:id` - 删除用户
- `PUT /users/:id/status` - 修改用户状态
- `PUT /users/:id/password` - 重置密码

### 角色管理

- `GET /roles` - 角色列表
- `GET /roles/:id` - 角色详情
- `POST /roles` - 创建角色
- `PUT /roles/:id` - 更新角色
- `DELETE /roles/:id` - 删除角色
- `PUT /roles/:id/menus` - 分配菜单权限
- `PUT /roles/:id/apis` - 分配API权限

### 菜单管理

- `GET /menus` - 菜单列表(树形)
- `GET /menus/:id` - 菜单详情
- `POST /menus` - 创建菜单
- `PUT /menus/:id` - 更新菜单
- `DELETE /menus/:id` - 删除菜单

### API管理

- `GET /apis` - API列表
- `GET /apis/:id` - API详情
- `POST /apis` - 创建API
- `PUT /apis/:id` - 更新API
- `DELETE /apis/:id` - 删除API

### 日志管理

- `GET /logs/operation` - 操作日志列表
- `GET /logs/login` - 登录日志列表

## License

MIT
