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
5. Save and preview the login page
