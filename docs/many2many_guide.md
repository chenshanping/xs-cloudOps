# 多对多关系使用指南

## 概述
本文档说明代码生成器中多对多关系的配置和工作原理。

## 1. 配置示例

### 场景：试卷与题目的多对多关系
- **试卷表 (psy_paper)**：一份试卷包含多道题目
- **题目表 (psy_question)**：一道题目可以被多份试卷使用
- **中间表 (paper_question)**：存储试卷和题目的关联关系

### 代码生成器配置
在"关联关系"标签页添加：
- **关联类型**: 多对多 (many2many)
- **关联表**: psy_question
- **显示字段**: name
- **中间表**: paper_question
- **注释**: 试卷题目

## 2. 生成的代码结构

### 2.1 Model（模型）
```go
type PsyPaper struct {
    BaseModel
    Name         string        `json:"name"`
    // ... 其他字段
    PsyQuestions []PsyQuestion `json:"psy_questions" gorm:"many2many:paper_question;"`
}
```

### 2.2 中间表SQL
```sql
CREATE TABLE IF NOT EXISTS `paper_question` (
  `psy_paper_id` BIGINT UNSIGNED NOT NULL COMMENT '试卷表ID',
  `psy_question_id` BIGINT UNSIGNED NOT NULL COMMENT '试卷题目ID',
  PRIMARY KEY (`psy_paper_id`, `psy_question_id`),
  INDEX `idx_psy_paper_id` (`psy_paper_id`),
  INDEX `idx_psy_question_id` (`psy_question_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.3 Request（请求）
```go
type CreatePsyPaperRequest struct {
    Name            string `json:"name"`
    // ... 其他字段
    PsyQuestionsIds []uint `json:"psy_questions_ids"` // 题目ID列表
}
```

### 2.4 Service（服务层）

#### 创建操作
```go
func (s *PsyPaperService) CreatePsyPaper(req *request.CreatePsyPaperRequest, userID uint) error {
    data := model.PsyPaper{
        Name: req.Name,
        // ... 其他字段
    }
    
    return global.DB.Transaction(func(tx *gorm.DB) error {
        // 1. 创建主记录
        if err := tx.Create(&data).Error; err != nil {
            return err
        }
        
        // 2. 处理多对多关联
        if len(req.PsyQuestionsIds) > 0 {
            var psy_questions []model.PsyQuestion
            if err := tx.Where("id IN ?", req.PsyQuestionsIds).Find(&psy_questions).Error; err != nil {
                return err
            }
            // 3. 自动插入中间表数据
            if err := tx.Model(&data).Association("PsyQuestions").Replace(psy_questions); err != nil {
                return err
            }
        }
        return nil
    })
}
```

#### 更新操作
```go
func (s *PsyPaperService) UpdatePsyPaper(id uint, req *request.UpdatePsyPaperRequest) error {
    var data model.PsyPaper
    if err := global.DB.First(&data, id).Error; err != nil {
        return errors.New("数据不存在")
    }
    
    updates := map[string]interface{}{
        "name": req.Name,
        // ... 其他字段
    }
    
    return global.DB.Transaction(func(tx *gorm.DB) error {
        // 1. 更新主记录
        if err := tx.Model(&data).Updates(updates).Error; err != nil {
            return err
        }
        
        // 2. 更新关联（Replace会自动清除旧关联并建立新关联）
        var psy_questions []model.PsyQuestion
        if len(req.PsyQuestionsIds) > 0 {
            if err := tx.Where("id IN ?", req.PsyQuestionsIds).Find(&psy_questions).Error; err != nil {
                return err
            }
        }
        if err := tx.Model(&data).Association("PsyQuestions").Replace(psy_questions); err != nil {
            return err
        }
        return nil
    })
}
```

#### 删除操作
```go
func (s *PsyPaperService) DeletePsyPaper(id uint) error {
    var data model.PsyPaper
    if err := global.DB.First(&data, id).Error; err != nil {
        return errors.New("数据不存在")
    }
    
    return global.DB.Transaction(func(tx *gorm.DB) error {
        // 1. 清除中间表关联
        if err := tx.Model(&data).Association("PsyQuestions").Clear(); err != nil {
            return err
        }
        // 2. 删除主记录
        return tx.Delete(&data).Error
    })
}
```

## 3. 工作原理

### 3.1 GORM自动管理中间表

当使用 `Association().Replace()` 时，GORM会：

1. **创建时**：
   - 插入试卷记录到 `psy_paper` 表
   - 根据传入的题目ID列表，在 `paper_question` 中插入关联记录

2. **更新时**：
   - 删除 `paper_question` 中该试卷的所有旧关联
   - 插入新的关联记录

3. **删除时**：
   - 使用 `Association().Clear()` 删除中间表中的关联
   - 删除主记录

### 3.2 数据流示例

#### 创建试卷并关联3道题目
```json
POST /api/psy_paper
{
  "name": "心理测试问卷A",
  "psy_questions_ids": [1, 2, 3]
}
```

**数据库操作**：
1. `psy_paper` 表插入：`{id: 10, name: "心理测试问卷A"}`
2. `paper_question` 表插入：
   - `{psy_paper_id: 10, psy_question_id: 1}`
   - `{psy_paper_id: 10, psy_question_id: 2}`
   - `{psy_paper_id: 10, psy_question_id: 3}`

#### 查询试卷及关联题目
```go
var paper model.PsyPaper
global.DB.Preload("PsyQuestions").First(&paper, 10)
```

返回：
```json
{
  "id": 10,
  "name": "心理测试问卷A",
  "psy_questions": [
    {"id": 1, "name": "题目1"},
    {"id": 2, "name": "题目2"},
    {"id": 3, "name": "题目3"}
  ]
}
```

## 4. 前端交互

### 4.1 表单组件
```vue
<a-form-item label="选择题目" name="psy_questions_ids">
  <a-select 
    v-model:value="formState.psy_questions_ids" 
    mode="multiple"
    placeholder="请选择题目" 
    allow-clear
  >
    <a-select-option v-for="item in psyQuestionsOptions" :key="item.id" :value="item.id">
      {{ item.name }}
    </a-select-option>
  </a-select>
</a-form-item>
```

### 4.2 提交数据格式
```typescript
{
  name: "心理测试问卷A",
  psy_questions_ids: [1, 2, 3]  // 题目ID数组
}
```

## 5. 注意事项

1. **中间表由GORM自动管理**，无需手动操作
2. **中间表命名规范**：`表1_表2`（如 `paper_question`）
3. **使用事务**确保数据一致性
4. **删除主记录前需清除关联**，否则会留下孤立的中间表记录
5. **Preload加载**：查询时需要关联数据，记得添加 `Preload("PsyQuestions")`

## 6. 常见问题

### Q: 为什么中间表没有数据？
A: 检查以下几点：
- 中间表是否已创建（执行SQL文件）
- GORM标签是否正确：`gorm:"many2many:paper_question;"`
- 是否使用了事务并调用了 `Association().Replace()`

### Q: 如何只添加关联而不替换？
A: 使用 `Association().Append()` 代替 `Replace()`

### Q: 如何移除部分关联？
A: 使用 `Association().Delete()`

## 7. 总结

多对多关系的核心：
- ✅ **Model定义**：`gorm:"many2many:中间表名"`
- ✅ **中间表SQL**：自动生成
- ✅ **Service层**：使用事务 + `Association()` API
- ✅ **前端**：传递ID数组

GORM会自动处理中间表的所有操作，开发者只需要关注业务逻辑！
