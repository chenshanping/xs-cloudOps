# 改进：统一使用options接口

## 改进内容

在导入导出功能中，关联字段（belongsTo）现在统一使用options接口，而不是直接查询数据库表。

## 改进前后对比

### 改进前

**下载模板：** 直接查询表
```go
var relatedList []model.ProductType
global.DB.Find(&relatedList)
// 问题：包含所有数据，包括禁用的
```

**导入验证：** 直接查询表
```go
var related model.ProductType
global.DB.Where("name = ?", cellValue).First(&related)
// 问题：可以导入禁用的产品类型
```

### 改进后

**下载模板：** 使用options接口
```go
optionsList, _ := service.ProductType.GetProductTypeOptions("name", "", "", false, 0)
// ✅ 自动过滤status=1
```

**导入验证：** 使用options接口
```go
optionsList, _ := service.ProductType.GetProductTypeOptions("name", "", "", false, 0)
// ✅ 自动过滤status=1
// ✅ 验证失败时显示所有有效值
```

## 优点

### 1. 自动过滤条件

options接口已经包含了业务逻辑：
- ✅ 过滤 `status=1`（只显示启用的）
- ✅ 可以扩展其他过滤条件
- ✅ 不需要在每个地方重复判断

### 2. 统一逻辑

三个地方使用同一个接口：
- ✅ 前端下拉框：调用options接口
- ✅ 下载模板：调用options接口
- ✅ 导入验证：调用options接口

### 3. 更好的错误提示

导入失败时显示所有有效值：

**改进前：**
```
第2行: 产品类型"测试"不存在，请先创建该产品类型
```

**改进后：**
```
第2行: 产品类型"测试"不存在，有效值: 类型A, 类型B, 类型C
```

用户可以直接看到应该填什么值。

### 4. 性能优化

options接口可以优化查询：
- 只查询需要的字段（id, name）
- 可以添加缓存
- 可以添加分页（如果选项太多）

## 使用场景

### 场景1：产品类型被禁用

1. 管理员禁用了某个产品类型（status=0）
2. 用户下载产品导入模板
3. 下拉列表中不显示该类型 ✅
4. 用户无法选择禁用的类型
5. 即使手动输入，导入时也会验证失败 ✅

### 场景2：导入旧数据

1. 用户有一个旧的Excel文件
2. 文件中包含已禁用的产品类型
3. 导入时验证失败
4. 错误提示显示当前有效的类型列表
5. 用户可以快速修正数据

### 场景3：数据一致性

1. 前端下拉框显示的选项
2. 下载模板的下拉验证
3. 导入时的验证规则
4. 三者完全一致 ✅

## 配置说明

### 启用options接口

在生成器配置中，为关联字段设置 `use_options_api: true`：

```json
{
  "columns": [
    {
      "column_name": "product_type_id",
      "related_table": "product_type",
      "display_field": "name",
      "use_options_api": true  // ✅ 启用options接口
    }
  ]
}
```

### 不启用options接口

如果不启用（`use_options_api: false` 或未设置），会直接查询表：

```go
// 下载模板
var relatedList []model.ProductType
global.DB.Find(&relatedList)

// 导入验证
var related model.ProductType
global.DB.Where("name = ?", cellValue).First(&related)
```

**适用场景：**
- 关联表没有status字段
- 不需要过滤条件
- 关联表数据量很小

## 实现细节

### options接口返回格式

```go
[]map[string]interface{}{
    {
        "id": 1,
        "name": "类型A",
        "count": 10,  // 可选，统计关联数量
    },
    {
        "id": 2,
        "name": "类型B",
        "count": 5,
    },
}
```

### 下载模板代码

```go
{{- if .UseOptionsApi}}
// 使用options接口获取选项（已过滤status等条件）
optionsList, _ := service.{{.RelatedModel}}.Get{{.RelatedModel}}Options("{{.DisplayField}}", "", "", false, 0)
options := make([]string, 0, len(optionsList))
for _, item := range optionsList {
    if name, ok := item["name"].(string); ok {
        options = append(options, name)
    }
}
if len(options) > 0 && len(options) <= 100 {
    exporter.AddDataValidation(colIndex, options, 2, 1000)
}
{{- else}}
// 直接查询关联表
var relatedList []model.{{.RelatedModel}}
if err := global.DB.Find(&relatedList).Error; err == nil {
    options := make([]string, 0, len(relatedList))
    for _, item := range relatedList {
        options = append(options, item.{{ToPascalCase .DisplayField}})
    }
    if len(options) > 0 && len(options) <= 100 {
        exporter.AddDataValidation(colIndex, options, 2, 1000)
    }
}
{{- end}}
```

### 导入验证代码

```go
{{- if .UseOptionsApi}}
// 使用options接口验证（已过滤status等条件）
optionsList, _ := service.{{.RelatedModel}}.Get{{.RelatedModel}}Options("{{.DisplayField}}", "", "", false, 0)
found := false
var relatedID uint
for _, item := range optionsList {
    if name, ok := item["name"].(string); ok && name == cellValue {
        if id, ok := item["id"].(uint); ok {
            relatedID = id
            found = true
            break
        }
        // 处理id可能是float64的情况（JSON解析）
        if id, ok := item["id"].(float64); ok {
            relatedID = uint(id)
            found = true
            break
        }
    }
}
if found {
    createReq.{{.ForeignKey | ToPascalCase}} = relatedID
} else {
    // 获取所有有效值
    validValues := make([]string, 0, len(optionsList))
    for _, item := range optionsList {
        if name, ok := item["name"].(string); ok {
            validValues = append(validValues, name)
        }
    }
    if len(validValues) > 0 {
        rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}\"%s\"不存在，有效值: %s", cellValue, strings.Join(validValues, ", ")))
    } else {
        rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}\"%s\"不存在，请先创建该{{.Comment}}", cellValue))
    }
    hasError = true
}
{{- else}}
// 直接查询关联表
var related model.{{.RelatedModel}}
if err := global.DB.Where("{{.DisplayField}} = ?", cellValue).First(&related).Error; err == nil {
    createReq.{{.ForeignKey | ToPascalCase}} = related.ID
} else {
    rowErrors = append(rowErrors, fmt.Sprintf("{{.Comment}}\"%s\"不存在，请先创建该{{.Comment}}", cellValue))
    hasError = true
}
{{- end}}
```

## 注意事项

### 1. ID类型转换

options接口返回的map中，id可能是 `uint` 或 `float64`（JSON解析）：

```go
if id, ok := item["id"].(uint); ok {
    relatedID = id
} else if id, ok := item["id"].(float64); ok {
    relatedID = uint(id)
}
```

### 2. 选项数量限制

Excel下拉验证有限制，只在选项数量 ≤ 100 时添加：

```go
if len(options) > 0 && len(options) <= 100 {
    exporter.AddDataValidation(colIndex, options, 2, 1000)
}
```

如果选项太多，只显示批注说明，不添加下拉验证。

### 3. 空值列表

如果options接口返回空列表（没有启用的数据）：

**下载模板：** 不添加下拉验证，只显示批注

**导入验证：** 提示"请先创建该XXX"

### 4. 性能考虑

每次下载模板或导入时都会调用options接口：

**优化建议：**
- options接口添加缓存（Redis）
- 缓存时间：5-10分钟
- 数据变更时清除缓存

## 测试建议

### 测试用例1：禁用产品类型

1. 创建产品类型A（status=1）
2. 创建产品类型B（status=0，禁用）
3. 下载产品导入模板
4. 预期：下拉列表只有类型A
5. 手动输入类型B并导入
6. 预期：验证失败，提示"有效值: 类型A"

### 测试用例2：有效值提示

1. 创建3个启用的产品类型
2. 下载模板
3. 填写一个不存在的类型名称
4. 导入
5. 预期：错误提示显示所有3个有效值

### 测试用例3：直接查询模式

1. 配置关联字段 `use_options_api: false`
2. 生成代码
3. 下载模板
4. 预期：下拉列表包含所有产品类型（包括禁用的）

## 总结

统一使用options接口的好处：

1. ✅ **自动过滤**：status=1等条件自动应用
2. ✅ **统一逻辑**：前端、下载模板、导入验证使用同一接口
3. ✅ **更好提示**：显示所有有效值列表
4. ✅ **易于维护**：业务逻辑集中在options接口
5. ✅ **可扩展**：可以添加更多过滤条件

现在导入导出功能更加完善和用户友好了！🎉
