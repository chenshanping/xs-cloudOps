# 导入导出功能 - 完整实现说明

## 功能概述

代码生成器现已完整支持Excel导入导出功能，包括后端API、前端UI、权限控制等全套功能。

## 已实现的功能

### 1. 后端实现 ✅

#### Excel工具类
- **文件位置**: `go-base-server/utils/excel.go`
- **功能**:
  - ExcelExporter - 导出器（设置表头、添加数据行、生成文件）
  - ExcelImporter - 导入器（读取Excel、解析数据）
  - ParseCellValue - 类型转换（字符串→各种数据类型）

#### 字典服务扩展
- **文件位置**: `go-base-server/service/dict.go`
- **新增方法**:
  - `GetDictLabel(dictType, value)` - 值转标签（用于导出）
  - `GetDictValue(dictType, label)` - 标签转值（用于导入）

#### API接口（模板自动生成）
- **文件**: `go-base-server/generator/template/api.tpl`
- **接口**:
  - `GET /{module}/export` - 导出数据
  - `POST /{module}/import` - 导入数据
  - `GET /{module}/template` - 下载模板

#### Service层（模板自动生成）
- **文件**: `go-base-server/generator/template/service.tpl`
- **新增方法**:
  - `GetAll{ModelName}(req)` - 获取所有数据（不分页，用于导出）

#### 路由注册（模板自动生成）
- **文件**: `go-base-server/generator/template/router.tpl`
- **路由**:
  - 导出路由
  - 导入路由
  - 模板下载路由

### 2. 前端实现 ✅

#### API接口（模板自动生成）
- **文件**: `go-base-server/generator/template/frontend_api.tpl`
- **接口**:
  ```typescript
  export{ModelName}(params)      // 导出
  import{ModelName}(file)        // 导入
  downloadTemplate{ModelName}()  // 下载模板
  ```

#### UI组件（模板自动生成）
- **文件**: `go-base-server/generator/template/frontend_view.tpl`
- **按钮**:
  - 导出按钮（带权限控制）
  - 导入按钮（Upload组件，带权限控制）
  - 下载模板按钮（带权限控制）

#### 处理方法（模板自动生成）
- `handleExport()` - 导出数据，自动下载Excel
- `handleImport(file)` - 导入数据，显示结果
- `handleDownloadTemplate()` - 下载模板

### 3. 权限控制 ✅

#### 菜单SQL（模板自动生成）
- **文件**: `go-base-server/generator/template/menu_sql.tpl`
- **权限**:
  - `{permission}:export` - 导出权限
  - `{permission}:import` - 导入权限

#### 前端权限指令
```vue
v-permission="'{permission}:export'"  // 导出按钮
v-permission="'{permission}:import'"  // 导入按钮
```

### 4. 智能数据转换 ✅

#### 数据字典转换
- **导出**: 值 → 标签文本
- **导入**: 标签文本 → 值
- **示例**: 状态字段 1→"启用", "启用"→1

#### 开关字段转换
- **导出**: 值 → 配置的文本
- **导入**: 文本 → 值
- **示例**: 1→"是", "是"→1

#### 关联关系处理
- **导出**: 显示关联对象的display_field
- **导入**: 根据display_field查询ID
- **示例**: category_id → "电子产品", "电子产品" → 查询ID

#### 时间字段格式化
- **导出**: 格式化为 "2006-01-02 15:04:05"
- **导入**: 支持多种格式自动解析

#### 数据隔离
- **导出**: 自动应用数据隔离规则
- **导入**: 自动设置created_by

## 使用流程

### 1. 生成代码
使用代码生成器生成模块时，自动包含导入导出功能。

### 2. 执行菜单SQL
```sql
-- 执行生成的菜单SQL文件
source go-base-server/sql/{module}_menu.sql
```

### 3. 配置角色权限
在角色管理中勾选：
- 导出权限
- 导入权限

### 4. 前端使用
用户登录后，在列表页面可以看到：
- 导出按钮
- 导入按钮
- 下载模板按钮

## 完整示例

### 后端API调用
```bash
# 导出数据
GET /api/v1/product/export?status=1&name=测试
Authorization: Bearer {token}

# 导入数据
POST /api/v1/product/import
Authorization: Bearer {token}
Content-Type: multipart/form-data
file: [Excel文件]

# 下载模板
GET /api/v1/product/template
Authorization: Bearer {token}
```

### 前端使用
```vue
<template>
  <!-- 工具栏按钮（自动生成） -->
  <a-button @click="handleExport" v-permission="'product:export'">
    <DownloadOutlined /> 导出
  </a-button>
  
  <a-upload
    :show-upload-list="false"
    :before-upload="handleImport"
    accept=".xlsx,.xls"
    v-permission="'product:import'"
  >
    <a-button><UploadOutlined /> 导入</a-button>
  </a-upload>
  
  <a-button @click="handleDownloadTemplate" v-permission="'product:import'">
    <FileExcelOutlined /> 下载模板
  </a-button>
</template>

<script setup lang="ts">
// 导出（自动生成）
const handleExport = async () => {
  const res = await exportProduct(searchForm)
  // 自动下载Excel文件
}

// 导入（自动生成）
const handleImport = async (file: File) => {
  const res = await importProduct(file)
  // 显示导入结果
  return false
}

// 下载模板（自动生成）
const handleDownloadTemplate = async () => {
  const res = await downloadTemplateProduct()
  // 自动下载模板文件
}
</script>
```

### Excel模板格式
| 产品名称(必填) | 数量(必填) | 价格(必填) | 状态(选填) | 产品类型(必填) |
|--------------|----------|----------|----------|-------------|
| 示例产品      | 100      | 99.9     | 启用      | 电子产品     |

## 文件清单

### 后端文件
- ✅ `utils/excel.go` - Excel工具类
- ✅ `service/dict.go` - 字典服务（新增方法）
- ✅ `generator/template/api.tpl` - API模板（含导入导出）
- ✅ `generator/template/service.tpl` - Service模板（含GetAll方法）
- ✅ `generator/template/request.tpl` - Request模板（含QueryRequest）
- ✅ `generator/template/router.tpl` - Router模板（含路由注册）
- ✅ `generator/template/menu_sql.tpl` - 菜单SQL模板（含权限）

### 前端文件
- ✅ `generator/template/frontend_api.tpl` - API模板（含导入导出接口）
- ✅ `generator/template/frontend_view.tpl` - 视图模板（含UI和方法）

### 文档文件
- ✅ `docs/import_export_guide.md` - 功能说明文档
- ✅ `docs/import_export_test.md` - 测试指南
- ✅ `docs/import_export_complete.md` - 完整实现说明（本文档）

## 依赖库

### 后端
```bash
go get github.com/xuri/excelize/v2
```

### 前端
无需额外依赖，使用原有的axios和ant-design-vue。

## 特性总结

✅ 自动生成完整的导入导出功能  
✅ 智能数据转换（字典、关联、开关、时间）  
✅ 权限控制（菜单权限+前端指令）  
✅ 错误处理（逐行验证，详细错误信息）  
✅ 数据隔离支持  
✅ 带示例的导入模板  
✅ 前后端完整集成  
✅ 开箱即用，无需额外配置  

## 注意事项

1. **文件字段不支持**: image、file、upload等文件类型字段会被跳过
2. **many2many不支持导入**: 多对多关联关系暂不支持通过Excel导入
3. **唯一性校验**: 导入时会执行唯一性校验
4. **数据量限制**: 建议单次导入不超过1000条数据
5. **权限配置**: 需要在角色管理中配置导入导出权限

## 后续优化建议

1. 支持文件字段的URL导入
2. 支持many2many关联的导入
3. 添加导入进度条
4. 支持大文件分批导入
5. 添加导入历史记录
6. 支持自定义导出字段选择

## 技术栈

- **后端**: Go + Gin + GORM + Excelize
- **前端**: Vue 3 + TypeScript + Ant Design Vue
- **数据库**: MySQL


## 操作日志记录

### 记录内容

导入导出操作会自动记录到`sys_operation_log`表中，包含以下信息：

#### 导出操作
- **路径**: `/api/v1/{module}/export`
- **方法**: GET
- **摘要**: 导出{模块名}
- **请求**: 查询参数（搜索条件）
- **响应**: `[文件下载]`
- **状态码**: 200
- **耗时**: 实际处理时间（毫秒）

#### 导入操作
- **路径**: `/api/v1/{module}/import`
- **方法**: POST
- **摘要**: 导入{模块名}
- **请求**: `[文件上传] file: filename.xlsx (123.45KB)`
- **响应**: 导入结果JSON（成功数、失败数、错误详情）
- **状态码**: 200
- **耗时**: 实际处理时间（毫秒）

#### 下载模板操作
- **路径**: `/api/v1/{module}/template`
- **方法**: GET
- **摘要**: 下载导入模板
- **请求**: 空
- **响应**: `[文件下载]`
- **状态码**: 200
- **耗时**: 实际处理时间（毫秒）

### 查询操作日志

使用提供的SQL脚本查询操作日志：

```bash
# 执行查询脚本
mysql -u root -p your_database < docs/check_import_export_logs.sql
```

或在数据库管理工具中执行：

```sql
-- 查看最近的导入导出操作
SELECT 
    username,
    method,
    path,
    summary,
    status,
    latency,
    created_at
FROM sys_operation_log
WHERE path LIKE '%/export' 
   OR path LIKE '%/import' 
   OR path LIKE '%/template'
ORDER BY created_at DESC
LIMIT 20;
```

### 日志示例

#### 导出操作日志
```json
{
  "id": 1001,
  "user_id": 1,
  "username": "admin",
  "ip": "192.168.1.100",
  "method": "GET",
  "path": "/api/v1/product/export",
  "group": "产品管理",
  "summary": "导出产品",
  "request": "status=1&name=测试",
  "response": "[文件下载]",
  "status": 200,
  "business_code": 200,
  "latency": 1250,
  "created_at": "2024-02-14 10:30:00"
}
```

#### 导入操作日志
```json
{
  "id": 1002,
  "user_id": 1,
  "username": "admin",
  "ip": "192.168.1.100",
  "method": "POST",
  "path": "/api/v1/product/import",
  "group": "产品管理",
  "summary": "导入产品",
  "request": "[文件上传] file: products.xlsx (45.67KB)",
  "response": "{\"success_count\":10,\"fail_count\":2,\"total\":12,\"errors\":[...]}",
  "status": 200,
  "business_code": 200,
  "latency": 3500,
  "created_at": "2024-02-14 10:35:00"
}
```

### 日志分析

#### 统计导入导出次数
```sql
SELECT 
    username,
    COUNT(CASE WHEN summary LIKE '%导出%' THEN 1 END) as export_count,
    COUNT(CASE WHEN summary LIKE '%导入%' THEN 1 END) as import_count
FROM sys_operation_log
WHERE summary LIKE '%导出%' OR summary LIKE '%导入%'
GROUP BY username;
```

#### 分析操作耗时
```sql
SELECT 
    CASE 
        WHEN path LIKE '%/export' THEN '导出'
        WHEN path LIKE '%/import' THEN '导入'
    END as operation_type,
    AVG(latency) as avg_latency_ms,
    MAX(latency) as max_latency_ms
FROM sys_operation_log
WHERE path LIKE '%/export' OR path LIKE '%/import'
GROUP BY operation_type;
```

#### 查看失败的操作
```sql
SELECT 
    username,
    path,
    summary,
    response,
    created_at
FROM sys_operation_log
WHERE (path LIKE '%/export' OR path LIKE '%/import')
  AND (status != 200 OR business_code != 200)
ORDER BY created_at DESC;
```

### 日志保留策略

建议定期清理旧的操作日志：

```sql
-- 删除30天前的操作日志
DELETE FROM sys_operation_log 
WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

-- 或者只保留最近10000条记录
DELETE FROM sys_operation_log 
WHERE id NOT IN (
    SELECT id FROM (
        SELECT id FROM sys_operation_log 
        ORDER BY created_at DESC 
        LIMIT 10000
    ) tmp
);
```

### 日志监控

可以基于操作日志实现监控告警：

1. **导入失败率监控**：如果导入失败率超过阈值，发送告警
2. **操作耗时监控**：如果导出/导入耗时过长，发送告警
3. **异常操作监控**：如果某用户频繁导出大量数据，发送告警

### 审计追踪

操作日志可用于审计追踪：

- 谁在什么时间导出了哪些数据
- 谁导入了什么文件
- 导入了多少条数据
- 是否有失败的操作
- 操作的IP地址和User-Agent

这些信息对于安全审计和问题排查非常重要。
