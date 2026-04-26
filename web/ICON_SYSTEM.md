# 图标系统说明

## 概述

该项目使用图标系统包括两部分：
1. **Ant Design Vue 官方图标** - 自动生成，包含 Outlined、Filled、TwoTone 三种样式
2. **自定义 SVG 图标** - 从 `src/assets/icons/` 目录自动扫描

## 官方图标

### 自动生成

官方图标列表由脚本自动生成：

```bash
npm run icons
```

这会扫描 `@ant-design/icons-vue` 包并生成 `src/utils/ant-icons.ts` 文件。

总共包含：
- **Outlined**: 421 个
- **Filled**: 218 个  
- **TwoTone**: 150 个
- **总计**: 789 个图标

### 生成的文件

生成的 `src/utils/ant-icons.ts` 包含：

```typescript
export const antIconsByType = {
  Outlined: [...],
  Filled: [...],
  TwoTone: [...]
}

export const allAntIcons = [
  { name: 'dashboard', type: 'Outlined', full: 'dashboardOutlined' },
  // ...
]
```

### 使用方式

在代码中使用：

```typescript
import { allAntIcons } from '@/utils/ant-icons'

// 获取所有图标
const icons = allAntIcons

// 按类型过滤
const outlinedOnly = allAntIcons.filter(icon => icon.type === 'Outlined')
```

## 自定义 SVG 图标

### 命名规则

自定义 SVG 图标应放在 `src/assets/icons/` 目录中，命名规则：
- 使用小写英文和下划线或中划线
- 示例：`phone.svg`, `user-add.svg`, `menu_icon.svg`

### 添加新图标

1. 将 SVG 文件放入 `src/assets/icons/` 目录
2. 重启开发服务器，图标会自动出现在 IconSelect 组件的"自定义SVG"选项卡中

### 存储格式

在数据库中的存储格式为：
- 官方图标：`official-dashboardOutlined`
- 自定义图标：`custom-phone`

## IconSelect 组件

### 功能

- **两个选项卡**：
  - Ant Design Icons：支持搜索和按类型过滤（Outlined/Filled/TwoTone）
  - 自定义SVG：支持搜索
- **图标预览**：显示已选择的图标

### 使用示例

```vue
<template>
  <IconSelect v-model="selectedIcon" />
</template>

<script setup>
import { ref } from 'vue'
import IconSelect from '@/components/IconSelect.vue'

const selectedIcon = ref('official-dashboardOutlined')
</script>
```

## 集成流程

### 开发流程

```bash
# 安装依赖
npm install

# 生成图标列表（npm run dev 会自动执行）
npm run icons

# 启动开发服务器
npm run dev
```

### 构建流程

```bash
# 构建时会自动生成最新的图标列表
npm run build
```

## 故障排除

### 图标不显示

1. 检查图标文件是否在 `src/assets/icons/` 目录中
2. 确保文件名没有特殊字符
3. 重启开发服务器

### 官方图标列表不更新

运行 `npm run icons` 手动生成

### 自定义 SVG 不加载

1. 确保使用正确的 SVG 格式
2. 检查浏览器控制台是否有错误
3. 尝试清理缓存并重新加载

## 技术细节

### 官方图标生成脚本

- **位置**：`scripts/generate-icons.mjs`
- **输入**：`node_modules/@ant-design/icons-vue`
- **输出**：`src/utils/ant-icons.ts`
- **方式**：扫描 npm 包文件系统，提取图标名称

### 自定义图标加载

- **方式**：使用 `import.meta.glob()` 动态扫描
- **路径**：`/src/assets/icons/*.svg`
- **去重**：使用 `Set` 防止重复
  o{
  {ID: "deepseek-v3.2", Name: "DeepSeek-V3.2", Description: "DeepSeek最新模型,支持联网和思考"},
  {ID: "kimi-k2-thinking", Name: "kimi-k2-thinking", Description: "kimi-k2-thinking模型是月之暗面提供的具有通用 Agentic能力和推理能力的思考模型，它擅长深度推理，并可通过多步工具调用，帮助解决各类难题"},
  {ID: "glm-4.7", Name: "GLM 4.7", Description: "智谱最新旗舰，具备更强的编程能力与更稳定的多步骤推理/执行能力。总参数355B，支持长程任务规划、编码、工具协同，问答自然、写作沉浸、创意角色扮演能力强"},
  {ID: "qwen3-vl-plus-2025-12-19", Name: "通义千问3-VL-Plus", Description: "Qwen3系列视觉理解模型，实现思考模式和非思考模式的有效融合。相较于9月23日快照，在推理及分析任务、风格控制上表现更优；同时拥有更低的延时和更快的响应速度。此版本为2025年12月19日快照版本"},
  {ID: "qwen3-max-2026-01-23", Name: "通义千问3-Max", Description: "通义千问3系列Max模型，相较2025年9月23日快照，此版本实现思考模式和非思考模式的有效融合，模型整体效果得到全方位的大幅度提升。在思考模式下，同时发布Web搜索、Web信息提取和代码解释器工具能力，使得模型在慢思考的同时，能够通过引入外部工具，以更高的准确性解决更有难度的问题。此版本为2026年1月23日快照"},
  }