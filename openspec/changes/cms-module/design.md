## Context

go-base 当前是一个通用 RBAC 后台管理框架，具备完善的用户/角色/菜单/权限体系。后端采用模块自注册模式（`RouterModule` + `init()`），前端菜单由 DB 驱动、路由动态生成。

CMS 模块需要利用这套现有机制，作为一组新文件加入项目，不改动框架核心。模块关注**内容管理**（栏目、文章、单页、Banner、前台导航），并提供公开只读 API 供外部前台消费。

### 当前模块结构参考

```
model/sys_menu.go          ← 模型文件（平铺）
service/menu/              ← 服务层（子目录）
api/v1/menu.go             ← Handler（平铺）
router/modules/menu.go     ← 路由（init 自注册）
```

## Goals / Non-Goals

**Goals:**

- CMS 模块完全遵循现有代码组织模式，无学习成本
- 删除 CMS 相关文件即可完全移除模块，零残留
- 后台 CRUD 页面体验与现有系统管理页面一致（Ant Design Vue + Drawer）
- 公开 API 支持前台项目（Nuxt/Next/SPA）独立消费
- 文章支持富文本编辑、草稿/发布/下架状态流转
- 栏目支持树形无限层级
- 权限与现有 Casbin 体系集成，CMS 菜单/API 受角色控制

**Non-Goals:**

- 不做前台渲染层（SSR/SSG）
- 不做可视化页面编辑器/模板拖拽
- 不做评论系统
- 不做 SEO 管理（sitemap/meta）
- 不做多站点/多租户
- 不做文章版本管理

## Decisions

### D1: 数据模型设计

所有 CMS 表使用 `cms_` 前缀，与系统表 `sys_` 区分。

**cms_category（栏目）**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | 主键 |
| parent_id | uint | 父栏目 ID（0=顶级） |
| name | varchar(100) | 栏目名称 |
| slug | varchar(100) UNIQUE | URL 标识 |
| cover | varchar(500) | 封面图 |
| description | text | 描述 |
| sort | int | 排序（越小越前） |
| status | tinyint | 1=启用, 0=禁用 |
| created_at, updated_at | datetime | 时间戳 |

**cms_article（文章）**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | 主键 |
| category_id | uint | 所属栏目 |
| title | varchar(200) | 标题 |
| slug | varchar(200) UNIQUE | URL 标识 |
| summary | varchar(500) | 摘要 |
| content | longtext | 正文（HTML） |
| cover | varchar(500) | 封面图 |
| author | varchar(100) | 作者 |
| status | tinyint | 0=草稿, 1=已发布, 2=已下架 |
| is_top | tinyint | 是否置顶 |
| views | int | 浏览量 |
| published_at | datetime NULL | 发布时间 |
| created_at, updated_at | datetime | 时间戳 |

**cms_page（单页）**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | 主键 |
| title | varchar(200) | 页面标题 |
| slug | varchar(200) UNIQUE | URL 标识（如 about-us） |
| content | longtext | 页面内容（HTML） |
| status | tinyint | 1=启用, 0=禁用 |
| sort | int | 排序 |
| created_at, updated_at | datetime | 时间戳 |

**cms_banner（轮播图）**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | 主键 |
| title | varchar(200) | 标题 |
| image | varchar(500) | 图片 URL |
| link | varchar(500) | 跳转链接 |
| position | varchar(50) | 位置标识（home_top, sidebar 等） |
| sort | int | 排序 |
| status | tinyint | 1=启用, 0=禁用 |
| created_at, updated_at | datetime | 时间戳 |

**cms_navigation（前台导航）**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | 主键 |
| parent_id | uint | 父节点 |
| name | varchar(100) | 导航名称 |
| link | varchar(500) | 链接 |
| target | varchar(20) | 打开方式（_self, _blank） |
| sort | int | 排序 |
| status | tinyint | 1=启用, 0=禁用 |
| created_at, updated_at | datetime | 时间戳 |

**选择理由**：不使用外键（遵循 AGENTS.md 规则），栏目/导航的树形通过 `parent_id` + 应用层递归实现，与现有 `sys_menu` 模式一致。

**备选方案**：使用 closure table 或 materialized path 做树形。但现有 `sys_menu` 和 `sys_dept` 都用 `parent_id` 递归，保持一致更好。

### D2: 后端分层

```
server/
├── model/
│   ├── cms_category.go
│   ├── cms_article.go
│   ├── cms_page.go
│   ├── cms_banner.go
│   └── cms_navigation.go
├── model/request/
│   └── cms.go              # CMS 请求参数结构体
├── service/cms/
│   ├── category.go
│   ├── article.go
│   ├── page.go
│   ├── banner.go
│   └── navigation.go
├── api/v1/
│   └── cms.go              # 所有 CMS handler
└── router/modules/
    └── cms.go              # init() 自注册
```

**选择理由**：完全复用现有 `sys_*` 模块的分层模式，handler 放一个文件（参考 `dict.go` 含 dict + dict_data 两组 CRUD）。

**CMS 服务不注册到 `service/service.go` facade**——handler 直接导入 `service/cms` 包。这样删除 CMS 文件时不需要改 `service.go`，增强可插拔性。

### D3: API 路径设计

**管理端 API**（需认证 + Casbin）：

```
POST   /api/v1/cms/categories          创建栏目
PUT    /api/v1/cms/categories/:id       更新栏目
DELETE /api/v1/cms/categories/:id       删除栏目
GET    /api/v1/cms/categories           栏目列表（树形）

POST   /api/v1/cms/articles             创建文章
PUT    /api/v1/cms/articles/:id         更新文章
DELETE /api/v1/cms/articles/:id         删除文章
GET    /api/v1/cms/articles             文章列表（分页+筛选）
GET    /api/v1/cms/articles/:id         文章详情
PUT    /api/v1/cms/articles/:id/status  更新文章状态（发布/下架）

POST   /api/v1/cms/pages                创建单页
PUT    /api/v1/cms/pages/:id            更新单页
DELETE /api/v1/cms/pages/:id            删除单页
GET    /api/v1/cms/pages                单页列表
GET    /api/v1/cms/pages/:id            单页详情

POST   /api/v1/cms/banners              创建 Banner
PUT    /api/v1/cms/banners/:id          更新 Banner
DELETE /api/v1/cms/banners/:id          删除 Banner
GET    /api/v1/cms/banners              Banner 列表

POST   /api/v1/cms/navigations          创建导航
PUT    /api/v1/cms/navigations/:id      更新导航
DELETE /api/v1/cms/navigations/:id      删除导航
GET    /api/v1/cms/navigations          导航列表（树形）
```

**公开 API**（无需认证）：

```
GET    /api/v1/public/cms/categories           栏目树（仅启用）
GET    /api/v1/public/cms/categories/:slug     栏目下文章列表
GET    /api/v1/public/cms/articles             文章列表（分页，仅已发布）
GET    /api/v1/public/cms/articles/:slug       文章详情（同时 +1 浏览量）
GET    /api/v1/public/cms/pages/:slug          单页详情
GET    /api/v1/public/cms/banners              Banner 列表（按 position 筛选，仅启用）
GET    /api/v1/public/cms/navigations          导航树（仅启用）
```

**选择理由**：管理端和公开端分别走 `private` 和 `public` 路由组，复用现有的认证/Casbin 中间件链。公开 API 使用 `/public/cms/` 前缀，与管理端 `/cms/` 隔离。

### D4: 前端页面设计

```
web/src/
├── api/cms.ts                         # CMS API 请求
├── views/admin/cms/
│   ├── category/
│   │   ├── index.vue                  # 栏目树形管理页
│   │   └── components/
│   │       └── CategoryFormDrawer.vue # 创建/编辑栏目抽屉
│   ├── article/
│   │   ├── index.vue                  # 文章列表页
│   │   └── components/
│   │       └── ArticleFormDrawer.vue  # 创建/编辑文章抽屉
│   ├── page/
│   │   ├── index.vue                  # 单页列表
│   │   └── components/
│   │       └── PageFormDrawer.vue     # 创建/编辑单页抽屉
│   └── banner/
│       ├── index.vue                  # Banner 列表
│       └── components/
│           └── BannerFormDrawer.vue   # 创建/编辑 Banner 抽屉
```

- 所有编辑交互使用 **Drawer**（遵循 AGENTS.md）
- 复用现有 `FileUpload` 组件处理封面图上传
- 文章编辑器：引入 wangEditor 5（Vue3 版本），轻量且支持富文本，不依赖外部服务
- 栏目/导航使用 `a-tree` 组件展示树形

### D5: 菜单与权限初始化

在 `initialize/` 中使用 `FirstOrCreate + Attrs` 模式插入 CMS 内置菜单：

```
内容管理 (目录, sort: 30)
├── 栏目管理 (页面, component: views/admin/cms/category/index.vue)
├── 文章管理 (页面, component: views/admin/cms/article/index.vue)
├── 单页管理 (页面, component: views/admin/cms/page/index.vue)
├── Banner管理 (页面, component: views/admin/cms/banner/index.vue)
└── 导航管理 (页面, component: views/admin/cms/navigation/index.vue - 可选后续)
```

同时注册对应的 API 权限记录到 `sys_api` 表。

**不强制分配给任何角色**——管理员通过角色管理自行分配。

### D6: 富文本编辑器选择

| 选项 | 优点 | 缺点 |
|------|------|------|
| **wangEditor 5** ✅ | 中文社区活跃、Vue3 原生支持、轻量 | 插件生态不如 tiptap |
| tiptap | 高度可扩展、schema 强 | 配置复杂、体积较大 |
| TinyMCE | 功能最全 | 需要 API Key（云版本）或体积巨大 |

选择 wangEditor 5：安装 `@wangeditor/editor` + `@wangeditor/editor-for-vue`，约 200KB gzip，满足文章/单页编辑需求。

## Risks / Trade-offs

- **文章内容存储为 HTML** → 后续如果需要多端适配（小程序等），可能需要额外转换。首期用 HTML 满足 Web 场景。
- **slug 唯一约束** → 用户需要理解 slug 概念。解决：前端自动从标题生成 slug，允许手动修改。
- **富文本图片** → 文章中的图片上传复用现有 `file` 模块上传接口，wangEditor 配置自定义上传函数。
- **大量文章时树形栏目查询性能** → 首期 parent_id 递归足够，文章量到万级以上再考虑优化。

## Migration Plan

1. 后端：添加 model → service → api → router 文件，`go build` 验证
2. 数据库：AutoMigrate 自动建表 + `server/sql/create_cms_tables.sql` 升级脚本
3. 前端：添加 api → views → 菜单初始化后自动生效
4. 回滚：删除所有 `cms_` 前缀文件 + 删除 DB 菜单/API 记录 + DROP 表

## Open Questions

1. 导航管理（cms_navigation）是否在首期实现？还是等有实际前台项目后再加？
2. 文章是否需要标签（tag）功能？如果需要，要增加 `cms_tag` + `cms_article_tag` 关联表。
3. 是否需要文章的定时发布功能？首期建议不做，后续通过定时任务模块追加。
