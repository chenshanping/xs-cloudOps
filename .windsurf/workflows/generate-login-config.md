---
description: Generate login page configuration (system name, slogan, features, etc.) from a system functional description. Use when setting up a new project or rebranding.
---

# Generate Login Page Config

Given a system's functional description, generate complete login page configuration values.

## When To Use

- Setting up a new go-base project instance
- Rebranding or customizing the login page
- User says "帮我生成登录页配置" or provides system functionality and asks for login page content

## Input Required

Ask the user for:
1. **System functional description** — what does the system do? (e.g. "医院挂号管理系统", "企业OA办公平台")
2. **Target audience** — who uses it? (optional, infer from description if not given)
3. **Tone** — professional / friendly / tech / minimal (optional, default: professional)

## Output Format

Generate the following fields as a ready-to-paste configuration block:

```
系统名称 (sys_name):         <中文系统名称，简短有力，4-8字>
英文系统名称:                 <English System Name, Title Case, 2-4 words>
登录页标题 (login_title):     <欢迎语，2-6字，如"欢迎回来">
标语 (login_slogan):          <一句话定位，8-15字，突出系统核心价值>
副标题 (login_subtitle):      <补充说明，6-12字>
描述 (login_desc):            <1-2句话，30-50字，描述系统能力和价值>

特性标签 (login_features):
  1. icon: <icon_name> (<中文标签>)  title: <2-4字>  desc: <6-10字>
  2. icon: <icon_name> (<中文标签>)  title: <2-4字>  desc: <6-10字>
  3. icon: <icon_name> (<中文标签>)  title: <2-4字>  desc: <6-10字>
  4. icon: <icon_name> (<中文标签>)  title: <2-4字>  desc: <6-10字>
```

## Available Icons (24 icons)

Only use these icon names (registered in AuthLayout.vue + LoginRegisterConfig.vue):

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

## Rules

- **System name**: Concise, memorable, reflects core domain. Avoid generic names like "管理系统".
- **English name**: Professional, no abbreviations unless well-known (e.g. OA, CRM, ERP).
- **Slogan**: One phrase that captures the system's unique value. Avoid clichés.
- **Features**: Each feature should highlight a distinct capability. Use different icons for variety.
- **Tone consistency**: All text should feel like it belongs to the same brand.
- Do NOT just copy the defaults — tailor everything to the user's system description.

## Example

**Input**: "这是一个医院智能挂号和排班管理系统，支持患者在线预约、医生排班管理、科室管理和数据统计"

**Output**:

```
系统名称:     智慧医通
英文系统名称:  MediFlow
登录页标题:    欢迎使用
标语:          智慧医疗 高效就诊
副标题:        医院智能挂号与排班管理平台
描述:          整合预约挂号、医生排班、科室管理与数据分析，让医疗资源调度更智能，患者就诊更便捷。

特性标签:
  1. icon: RocketOutlined       title: 在线预约  desc: 患者自助挂号预约
  2. icon: TeamOutlined         title: 排班管理  desc: 智能医生排班调度
  3. icon: LineChartOutlined    title: 数据统计  desc: 就诊数据可视分析
  4. icon: SafetyOutlined       title: 信息安全  desc: 患者隐私合规保障
```

## Background Image Prompt

After generating the config, also output **3 background image prompts** the user can paste into AI image generators (豆包 / MidJourney / DALL-E).

### Prompt Rules

- Style: match the system tone — tech/business/medical/education etc.
- Color: derive from the system domain (e.g. blue-purple for tech, green for medical, warm for education).
- Composition: **16:9**, leave one side (left or right) with dark/low-detail area for login form overlay.
- Must include: `无文字`, `无水印`, `高清`.
- Provide 3 variants:
  1. **Abstract gradient** — safest, most versatile, geometric lines + particles.
  2. **Scene-based** — relates to the system domain (city for business, hospital for medical, campus for education).
  3. **Minimal texture** — frosted glass / mesh gradient / subtle pattern, ultra-clean.

### Prompt Template

Each prompt should follow this structure:

> [场景/主体描述]，[色调描述]，[风格关键词]，适合作为后台管理系统登录页背景，16:9宽幅构图，[留白方向]侧留出空间放置登录表单，无文字，无水印，高清

### Example (for 智慧医通)

```
方案A — 抽象科技流线:
深蓝渐变背景，半透明几何网格线条从左下角扩散，带有微光粒子效果，右侧融入医疗元素的极简线条图标（心电图、听诊器），科技感、商务、简洁大气，适合作为后台管理系统登录页背景，16:9宽幅构图，左侧留白区域用于放置登录表单，无文字，无水印，高清

方案B — 医疗场景:
现代化医院大厅，柔和自然光，蓝绿色调，前景虚化，远处走廊延伸感，叠加半透明数据图表全息效果，适合后台系统登录页背景，16:9，左侧偏暗留出空间，无文字，无水印，高清

方案C — 纯净渐变:
深色渐变背景，从左侧深靛蓝过渡到右侧薄荷绿，表面布满细腻磨砂质感和稀疏微光粒子，左下角有淡色几何圆环装饰，极简高级，适合SaaS后台登录页背景，16:9，左侧大面积留白，无文字，无水印，高清
```

## Logo Prompt

Also output **3 logo prompts** for AI image generation.

### Prompt Rules

- Must include: `无文字`, `透明背景`, `1:1`, `高清`, `矢量风格` or `扁平设计`.
- Derive the icon concept from the system name and domain.
- Provide 3 variants:
  1. **Letter-based** — initials of the English name combined with domain symbol.
  2. **Abstract symbol** — domain-related shapes merged into a geometric mark.
  3. **Name-concept** — visual metaphor of the Chinese system name (e.g. "云商管家" → cloud + shop).

### Prompt Template

> 设计一个现代[风格]Logo，[主体描述]，[配色描述]，扁平化设计，透明背景，适合作为后台管理系统图标，1:1正方形，无文字，高清

### Example (for 智慧医通 / MediFlow)

```
方案A — 字母组合:
设计一个现代科技风格Logo，字母M和F组合（MediFlow缩写），融入十字医疗符号，渐变色从深蓝到薄荷绿，扁平化设计，透明背景，适合作为后台管理系统图标，1:1正方形，无文字，高清

方案B — 抽象符号:
极简Logo设计，心电图波形与数据节点融合成环形标识，线条流畅，蓝绿渐变配色，扁平矢量风格，透明背景，适合SaaS医疗管理平台，1:1正方形，无文字，高清

方案C — 名称概念:
极简Logo，一个圆润的灯泡造型内部融合十字医疗符号和连接线路，象征"智慧医通"，配色为靛蓝渐变到青色，扁平矢量，透明背景，1:1，无文字，高清
```

### Logo Usage Tips

After generating, tell the user:
- Export PNG with transparent background
- Prepare two sizes: **128×128** (sidebar logo) and **32×32** (favicon)
- If the generator adds unwanted text, retry with `纯图形，不要任何文字和字母` appended

## After Generation

If the user confirms, help them update the system config:

1. Go to **系统配置 → 登录与注册** page
2. Fill in the generated values
3. For `login_features`, the JSON format is:
```json
[
  {"icon": "RocketOutlined", "title": "在线预约", "desc": "患者自助挂号预约"},
  {"icon": "TeamOutlined", "title": "排班管理", "desc": "智能医生排班调度"},
  {"icon": "LineChartOutlined", "title": "数据统计", "desc": "就诊数据可视分析"},
  {"icon": "SafetyOutlined", "title": "信息安全", "desc": "患者隐私合规保障"}
]
```
4. Optionally update `sys_name` in **系统配置 → 基础设置**
5. Upload the generated logo to **系统配置 → 基础设置 → 系统Logo**
6. Upload the background image to **系统配置 → 登录与注册 → 登录背景图**
7. Save and preview the login page
