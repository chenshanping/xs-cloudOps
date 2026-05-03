## Why

go-base 是一个通用后台管理框架，当前缺少内容管理能力。用户希望将 go-base 用于企业官网场景时，需要在后台管理栏目、文章、单页、Banner 等内容，并通过 API 供前台（独立项目或内置前台）消费。

CMS 作为一个**可插拔模块**加入框架，利用现有的 `RouterModule` 自注册机制、DB 驱动菜单、Casbin 权限体系，不修改框架核心代码。模块可以通过添加/删除文件的方式启用/禁用。

## What Changes

- **新增 CMS 数据模型**：栏目（树形）、文章、单页、Banner、前台导航，均以 `cms_` 前缀命名
- **新增 CMS 后端服务层**：`service/cms/` 目录，提供 CRUD + 树形操作 + 发布管理
- **新增 CMS API Handler**：`api/v1/cms.go`，提供管理端 API 和公开读取 API
- **新增 CMS 路由模块**：`router/modules/cms.go`，通过 `init()` 自注册
- **新增 CMS 表迁移**：在 `initialize/` 中注册 CMS 表的 AutoMigrate
- **新增 CMS 后台管理页面**：`web/src/views/admin/cms/` 下的栏目管理、文章管理、单页管理、Banner 管理
- **新增 CMS 内置菜单**：初始化时插入 CMS 菜单和 API 权限记录（create-only 模式）
- **新增 CMS 公开 API**：供前台项目（Nuxt/Next 或内置前台）消费的只读接口

## Capabilities

### New Capabilities

- `cms-category`: 栏目管理——树形栏目的 CRUD、排序、启用/禁用，支持多级嵌套
- `cms-article`: 文章管理——文章的 CRUD、富文本编辑、发布/草稿/下架状态流转、置顶、按栏目筛选
- `cms-page`: 单页管理——固定页面（关于我们、联系方式等）的 CRUD，slug 唯一标识
- `cms-banner`: Banner 管理——轮播图的 CRUD、排序、位置管理、启用/禁用
- `cms-navigation`: 前台导航管理——前台网站导航菜单的树形 CRUD，独立于后台菜单体系
- `cms-public-api`: CMS 公开读取 API——无需认证的只读接口，供前台项目消费

### Modified Capabilities

（无。CMS 模块是纯新增，不修改现有功能的行为规范。）

## Impact

### 后端

- **新增文件**：`model/cms_*.go`(5)、`service/cms/*.go`(5+)、`api/v1/cms.go`(1)、`router/modules/cms.go`(1)
- **微调文件**：`initialize/gorm.go`（增加 CMS 表迁移行）、`initialize/menu.go`/`initialize/api.go`（增加 CMS 内置菜单和 API 记录）
- **不修改**：`router/router.go`、`main.go`、`global/global.go`、`service/service.go`（CMS 服务不注册到 facade，由 handler 直接实例化）

### 前端

- **新增文件**：`views/admin/cms/` 下 4 个管理模块（category、article、page、banner），每个含 index.vue + components/
- **新增 API**：`api/cms.ts`
- **不修改**：路由系统（菜单 DB 驱动，自动生效）、布局、权限体系

### 数据库

- 新增 5 张表（`cms_category`、`cms_article`、`cms_page`、`cms_banner`、`cms_navigation`）
- 需要提供 `server/sql/` 下的建表升级脚本

### 依赖

- 无新依赖。富文本编辑器使用前端现有依赖或按需引入（如 wangeditor / tiptap）

### 回滚方案

- 删除所有 `cms_` 前缀的文件和 DB 菜单/API 记录即可完全移除模块
- 数据库表可保留或手动 DROP

## 不在范围内

- 前台官网的渲染层（SSR/SSG）——由独立前端项目或后续模块处理
- 模板引擎/页面编辑器——首期不做可视化拖拽编排
- SEO 优化（sitemap、meta 管理）——后续迭代
- 评论系统——后续迭代
- 多站点/多租户——后续迭代
