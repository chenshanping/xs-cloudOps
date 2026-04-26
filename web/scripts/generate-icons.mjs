#!/usr/bin/env node

/**
 * 生成 Ant Design Vue 图标列表
 * 通过扫描 @ant-design/icons-vue 包中的导出
 */

import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const projectRoot = path.resolve(__dirname, '..')
const nodeModulesPath = path.join(projectRoot, 'node_modules', '@ant-design', 'icons-vue')

// 读取 icons-vue 的 lib 目录
const libPath = path.join(nodeModulesPath, 'lib', 'icons')

if (!fs.existsSync(libPath)) {
  console.error(`Error: Cannot find icons directory at ${libPath}`)
  process.exit(1)
}

// 收集所有图标信息
const allIcons = []

// 扫描所有文件
const files = fs.readdirSync(libPath)
files.forEach(file => {
  if (file.endsWith('.js') || file.endsWith('.ts')) {
    const iconName = path.basename(file, path.extname(file))
    
    // 提取完整名称（保留大小写）
    if (iconName.endsWith('Outlined')) {
      const baseName = iconName.replace(/Outlined$/, '')
      allIcons.push({
        name: baseName,
        type: 'Outlined',
        full: iconName
      })
    } else if (iconName.endsWith('Filled')) {
      const baseName = iconName.replace(/Filled$/, '')
      allIcons.push({
        name: baseName,
        type: 'Filled',
        full: iconName
      })
    } else if (iconName.endsWith('TwoTone')) {
      const baseName = iconName.replace(/TwoTone$/, '')
      allIcons.push({
        name: baseName,
        type: 'TwoTone',
        full: iconName
      })
    }
  }
})

// 统计
const outlined = allIcons.filter(icon => icon.type === 'Outlined').length
const filled = allIcons.filter(icon => icon.type === 'Filled').length
const twoTone = allIcons.filter(icon => icon.type === 'TwoTone').length

// 生成 TypeScript 文件
const outputPath = path.join(projectRoot, 'src', 'utils', 'ant-icons.ts')
const content = `/**
 * Auto-generated Ant Design Vue icons list
 * Generated from @ant-design/icons-vue
 */

export interface AntIcon {
  name: string
  type: 'Outlined' | 'Filled' | 'TwoTone'
  full: string
}

export const allAntIcons: AntIcon[] = ${JSON.stringify(allIcons, null, 2)}
`

fs.writeFileSync(outputPath, content, 'utf-8')
console.log(`✓ Generated icon list at ${outputPath}`)
console.log(`  - Outlined: ${outlined}`)
console.log(`  - Filled: ${filled}`)
console.log(`  - TwoTone: ${twoTone}`)
console.log(`  - Total: ${allIcons.length}`)
