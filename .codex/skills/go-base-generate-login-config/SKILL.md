---
name: go-base-generate-login-config
description: Generate project-accurate login and register config content for this go-base repository. Use when configuring 系统配置基础配置 and 登录与注册 for a new project, rebrand, or login-page copy refresh.
---

# Generate Login Config

Generate a **project-accurate** login page config package for this `go-base` repository.

This workflow is not a generic branding prompt. It must follow the actual config keys, UI tabs, file fields, and rendering behavior already present in this project.

## Project Reality Check

Before generating content, align with the real implementation in this repo:

- Basic config tab: `web/src/views/admin/system/config/components/SystemConfig.vue`
- Login/Register tab: `web/src/views/admin/system/config/components/LoginRegisterConfig.vue`
- Auth layout: `web/src/layouts/AuthLayout.vue`
- Login page form: `web/src/views/auth/login/index.vue`
- Public config allowlist: `server/service/configsvc/config.go`

## When To Use

- Setting up a new `go-base` project instance
- Rebranding the login page for a specific business domain
- User says `帮我生成登录页配置`
- User gives a business/system description and wants ready-to-fill config content

## What This Project Actually Renders

Current live login page behavior:

- `sys_name`: displayed on the left showcase area
- `sys_logo`: displayed on the left showcase area
- `login_title`: displayed above the login form
- `login_slogan`: displayed on the left showcase area
- `login_desc`: displayed on the left showcase area
- `login_features`: displayed on the left showcase area, **最多实际显示 4 条**
- `login_bg_image` or `login_bg_color`: used as left showcase background
- `enable_register`: controls whether register / forgot-password entry appears

Current implementation caveats:

- `login_subtitle` is configurable, but the current login page form text is still hardcoded as `登录您的账户`
- `login_images` and `login_images_max` are editable in admin config, but the current `AuthLayout.vue` does not render them
- There is **no actual config key** for “英文系统名称”; if generated, treat it as branding reference only, not a persisted config item

## Required Inputs

Ask only for what is missing.

Minimum required:

1. **System/business description**
   Example: `AI基金和股票分析后台`

Helpful optional inputs:

2. **Target audience**
   Example: `投研人员、运营人员、老板自己看盘`

3. **Tone**
   Example: `专业`, `科技`, `克制`, `金融`, `极简`

4. **Register entry**
   Choose one:
   - `开启注册`
   - `关闭注册`

5. **Visual direction**
   Choose one:
   - `偏深色科技`
   - `偏浅色商务`
   - `偏金融图表风`
   - `让我推荐`

## Output Rules

- Prefer output that maps directly to existing config keys
- Do **not** invent new config keys
- Do **not** tell the user to add backend fields unless they explicitly ask for code changes
- Keep Chinese copy concise and usable in admin UI
- `login_features` must match the actual JSON structure:

```json
[
  { "icon": "CheckCircleOutlined", "title": "标题", "desc": "描述" }
]
```

- Recommend `login_features_max = 4` unless the user explicitly wants more
- Because the live login page only shows 4 features, do not generate more than 4 feature items by default
- If `enable_register = true`, remind the user to upload a register default avatar in the login/register tab
- File/image settings in this project should be applied through admin upload components, which persist `*_file_id` config keys behind the scenes

## Available Icons

Only use icon names that are already registered in this project:

| Icon Name              | Meaning  | Best For                        |
|------------------------|----------|---------------------------------|
| CheckCircleOutlined    | 成功/可靠 | 质量保障、合规、审核             |
| SafetyOutlined         | 安全     | 安全、权限、加密、合规           |
| LineChartOutlined      | 图表     | 数据分析、统计、报表             |
| ThunderboltOutlined    | 闪电     | 高效、快速、实时                 |
| RocketOutlined         | 火箭     | 创新、提升、加速                 |
| SettingOutlined        | 设置     | 配置灵活、自定义、系统管理       |
| CloudOutlined          | 云       | 云服务、在线协作、SaaS           |
| TeamOutlined           | 团队     | 协作、多人、组织管理             |
| GlobalOutlined         | 全球     | 国际化、多语言、全球业务         |
| DashboardOutlined      | 仪表盘   | 监控面板、总览、运维             |
| DatabaseOutlined       | 数据库   | 数据存储、题库、知识库           |
| ApiOutlined            | 接口     | 开放平台、API、集成              |
| BulbOutlined           | 灯泡     | 智能、创意、AI、灵感             |
| BookOutlined           | 书本     | 教育、学习、文档、知识           |
| FileProtectOutlined    | 文件保护 | 文档管理、合同、审批             |
| SolutionOutlined       | 方案     | 解决方案、咨询、专业服务         |
| ExperimentOutlined     | 实验     | 实验、测试、实操、科研           |
| FundOutlined           | 趋势     | 金融、投资、增长、走势           |
| ApartmentOutlined      | 组织架构 | 组织管理、层级、部门             |
| ScheduleOutlined       | 日程     | 排班、考试安排、任务计划         |
| MobileOutlined         | 手机     | 移动端、APP、响应式              |
| LikeOutlined           | 点赞     | 评价、满意度、社区               |
| StarOutlined           | 星标     | 收藏、评分、精选                 |
| CrownOutlined          | 皇冠     | VIP、高级、专业版                |

## Recommended Generation Process

### Step 1: Confirm the scope

Identify whether the user wants:

- just copy/content generation
- copy + background prompt
- copy + background + logo prompt
- actual code/config modification

If they only want content, do **not** drift into implementation.

### Step 2: Generate project-accurate config values

Output in this structure:

```text
一、基础配置（系统配置 → 基础配置）

sys_name:
<系统名称>

品牌参考（不入库）:
<可选英文名，仅作品牌灵感>

二、登录与注册（系统配置 → 登录与注册）

enable_register:
true | false

login_title:
<2-6字，例如：欢迎回来>

login_subtitle:
<可选；当前项目已可配置，但当前登录页主表单未实际展示>

login_slogan:
<8-16字，一句话定位>

login_desc:
<1-2句话，20-50字，说明系统价值>

login_bg_color:
<CSS 渐变字符串；如果用户后续上传背景图，可忽略此项>

login_features_max:
4

login_features:
[
  {"icon":"...","title":"...","desc":"..."},
  {"icon":"...","title":"...","desc":"..."},
  {"icon":"...","title":"...","desc":"..."},
  {"icon":"...","title":"...","desc":"..."}
]

三、当前项目渲染提醒

- 当前真实显示字段：...
- 当前仅可配置但未明显生效字段：...
```

### Step 3: Keep the copy grounded

Rules:

- `sys_name` should usually be 4-10 Chinese characters
- `login_title` should be short and direct
- `login_slogan` should express positioning, not marketing fluff
- `login_desc` should read like a serious admin system, not landing-page ad copy
- Avoid clichés such as:
  - `赋能增长`
  - `一站式闭环`
  - `数字化转型新引擎`
  unless the user explicitly wants that style

### Step 4: Generate optional image prompts

If the user wants image generation help, output:

- 3 background prompts
- 3 logo prompts

#### Background Prompt Rules

- Must include: `无文字`, `无水印`, `高清`
- Must be suitable for a **后台登录页左侧背景**
- Must mention **16:9**
- Must reserve one side for the form overlay
- Match the user’s domain and tone

Prompt template:

> [场景/主体描述]，[色调描述]，[风格关键词]，适合作为后台管理系统登录页背景，16:9宽幅构图，[左侧或右侧]保留低干扰区域用于放置登录表单，无文字，无水印，高清

#### Logo Prompt Rules

- Must include: `无文字`, `透明背景`, `1:1`, `高清`
- Prefer `扁平设计` or `矢量风格`
- Must fit sidebar/admin branding instead of poster-style illustration

Prompt template:

> 设计一个现代[风格]Logo，[主体描述]，[配色描述]，扁平设计，透明背景，适合作为后台管理系统图标，1:1正方形，无文字，高清

## Example Output Shape

```text
一、基础配置（系统配置 → 基础配置）

sys_name:
智投研判

品牌参考（不入库）:
AlphaScope

二、登录与注册（系统配置 → 登录与注册）

enable_register:
false

login_title:
欢迎登录

login_subtitle:
基金与股票分析后台

login_slogan:
聚焦数据 让判断更稳

login_desc:
整合A股与公募基金核心数据，支持行情浏览、指标查看与后台同步管理，帮助投研和运营快速获取关键市场信息。

login_bg_color:
linear-gradient(135deg, #0f1c3f 0%, #123f67 55%, #1d7a85 100%)

login_features_max:
4

login_features:
[
  {"icon":"FundOutlined","title":"市场走势","desc":"聚焦基金股票趋势"},
  {"icon":"LineChartOutlined","title":"数据分析","desc":"核心指标集中查看"},
  {"icon":"DashboardOutlined","title":"统一看板","desc":"后台信息快速总览"},
  {"icon":"SafetyOutlined","title":"权限清晰","desc":"管理操作边界明确"}
]

三、当前项目渲染提醒

- 当前真实显示字段：sys_name、sys_logo、login_title、login_slogan、login_desc、login_features、login_bg_image/login_bg_color、enable_register
- login_subtitle 当前可配置，但当前登录页主表单未直接显示
- login_images 当前后台可编辑，但 AuthLayout 暂未实际渲染
```

## After Generation

If the user confirms and only wants to apply values manually, guide them like this:

1. Open `系统配置 → 基础配置`
2. Fill `sys_name`
3. Upload the system logo through `系统Logo`
4. Open `系统配置 → 登录与注册`
5. Fill `login_title` / `login_slogan` / `login_desc`
6. Paste `login_features` content through the feature editor
7. Upload the login background image
8. If register is enabled, upload the register default avatar
9. Save and preview the login page

## If User Wants Actual Implementation

If the user asks to actually modify code or config defaults:

- first inspect the current neighboring implementation
- keep changes limited to existing config keys whenever possible
- do not invent new DB-configurable public keys unless explicitly requested
- if persisted defaults or bootstrap behavior are changed for existing installations, follow this repository’s SQL upgrade and initialize rules

## Hard Boundaries

- Do not generate secret/sensitive config values
- Do not advise putting API keys or auth bypass config into public login-page keys
- Do not claim a field is visible if the current project code does not actually render it
- Do not output generic “English system name” as if it were a stored config key
