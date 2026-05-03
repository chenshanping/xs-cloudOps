## 1. 数据模型（Backend Models）

- [ ] 1.1 创建 `server/model/cms_category.go` — CmsCategory 结构体（id, parent_id, name, slug, cover, description, sort, status）
- [ ] 1.2 创建 `server/model/cms_article.go` — CmsArticle 结构体（id, category_id, title, slug, summary, content, cover, author, status, is_top, views, published_at）
- [ ] 1.3 创建 `server/model/cms_page.go` — CmsPage 结构体（id, title, slug, content, status, sort）
- [ ] 1.4 创建 `server/model/cms_banner.go` — CmsBanner 结构体（id, title, image, link, position, sort, status）
- [ ] 1.5 创建 `server/model/cms_navigation.go` — CmsNavigation 结构体（id, parent_id, name, link, target, sort, status）
- [ ] 1.6 创建 `server/model/request/cms.go` — CMS 请求参数结构体（列表筛选、创建/更新入参）
- [ ] 1.7 在 `server/initialize/gorm.go` 中添加 CMS 表 AutoMigrate
- [ ] 1.8 验证：`cd server && go build ./...`

## 2. 数据库升级脚本

- [ ] 2.1 创建 `server/sql/create_cms_tables.sql` — 5 张 CMS 表的建表语句（幂等）
- [ ] 2.2 检查 `go-base.sql` 基线是否需要注意，但不修改基线

## 3. 服务层（Backend Services）

- [ ] 3.1 创建 `server/service/cms/category.go` — 栏目 CRUD + 树形构建 + slug 唯一校验 + 删除前检查（子节点/文章关联）
- [ ] 3.2 创建 `server/service/cms/article.go` — 文章 CRUD + 分页 + 状态流转 + slug 自动生成 + 浏览量递增
- [ ] 3.3 创建 `server/service/cms/page.go` — 单页 CRUD + slug 唯一校验
- [ ] 3.4 创建 `server/service/cms/banner.go` — Banner CRUD + 按 position 筛选
- [ ] 3.5 创建 `server/service/cms/navigation.go` — 导航 CRUD + 树形构建 + 删除前检查
- [ ] 3.6 验证：`cd server && go build ./...`

## 4. API Handler（Backend API）

- [ ] 4.1 创建 `server/api/v1/cms.go` — 管理端 handler（Category/Article/Page/Banner/Navigation CRUD）
- [ ] 4.2 在 cms.go 中添加公开 API handler（Public 只读接口：列表/详情/树形）
- [ ] 4.3 验证：`cd server && go build ./...`

## 5. 路由注册（Backend Router）

- [ ] 5.1 创建 `server/router/modules/cms.go` — init() 自注册，RegisterPublicRoutes 注册公开 API，RegisterPrivateRoutes 注册管理 API
- [ ] 5.2 验证：`cd server && go build ./...`

## 6. 菜单与权限初始化（Backend Initialize）

- [ ] 6.1 在 `server/initialize/` 中添加 CMS 内置菜单数据（"内容管理"目录 + 5 个子菜单页面），使用 FirstOrCreate + Attrs 模式
- [ ] 6.2 在 `server/initialize/` 中添加 CMS API 权限记录到 sys_api 表
- [ ] 6.3 创建 `server/sql/insert_cms_menu_api.sql` — CMS 菜单和 API 权限记录的升级脚本（幂等）
- [ ] 6.4 验证：`cd server && go test ./...`

## 7. 前端 API 层

- [ ] 7.1 创建 `web/src/api/cms.ts` — CMS 管理端 API 请求函数（category/article/page/banner/navigation CRUD）
- [ ] 7.2 验证：`cd web && npm run build`

## 8. 前端富文本编辑器

- [ ] 8.1 安装 `@wangeditor/editor` 和 `@wangeditor/editor-for-vue` 依赖
- [ ] 8.2 创建通用富文本编辑器组件 `web/src/components/RichEditor/index.vue`（封装 wangEditor，支持 v-model、图片上传对接现有 file 接口）
- [ ] 8.3 验证：`cd web && npm run build`

## 9. 前端栏目管理页面

- [ ] 9.1 创建 `web/src/views/admin/cms/category/index.vue` — 树形栏目列表页（a-tree 组件 + 操作按钮）
- [ ] 9.2 创建 `web/src/views/admin/cms/category/components/CategoryFormDrawer.vue` — 创建/编辑栏目抽屉
- [ ] 9.3 验证：`cd web && npm run build`

## 10. 前端文章管理页面

- [ ] 10.1 创建 `web/src/views/admin/cms/article/index.vue` — 文章列表页（分页表格 + 栏目筛选 + 状态筛选 + 关键词搜索）
- [ ] 10.2 创建 `web/src/views/admin/cms/article/components/ArticleFormDrawer.vue` — 创建/编辑文章抽屉（含富文本编辑器、封面上传、栏目选择）
- [ ] 10.3 验证：`cd web && npm run build`

## 11. 前端单页管理页面

- [ ] 11.1 创建 `web/src/views/admin/cms/page/index.vue` — 单页列表
- [ ] 11.2 创建 `web/src/views/admin/cms/page/components/PageFormDrawer.vue` — 创建/编辑单页抽屉（含富文本编辑器）
- [ ] 11.3 验证：`cd web && npm run build`

## 12. 前端 Banner 管理页面

- [ ] 12.1 创建 `web/src/views/admin/cms/banner/index.vue` — Banner 列表（图片预览 + position 筛选）
- [ ] 12.2 创建 `web/src/views/admin/cms/banner/components/BannerFormDrawer.vue` — 创建/编辑 Banner 抽屉
- [ ] 12.3 验证：`cd web && npm run build`

## 13. 前端导航管理页面（可选首期）

- [ ] 13.1 创建 `web/src/views/admin/cms/navigation/index.vue` — 导航树形管理
- [ ] 13.2 创建 `web/src/views/admin/cms/navigation/components/NavigationFormDrawer.vue` — 创建/编辑导航抽屉
- [ ] 13.3 验证：`cd web && npm run build`

## 14. 集成验证

- [ ] 14.1 后端全量测试：`cd server && go test ./...`
- [ ] 14.2 前端构建验证：`cd web && npm run build`
- [ ] 14.3 手动验证：启动服务，登录后台，分配 CMS 菜单给角色，验证栏目/文章 CRUD 流程
- [ ] 14.4 手动验证：公开 API 无需认证可读取已发布文章和启用栏目
- [ ] 14.5 提交代码并推送
