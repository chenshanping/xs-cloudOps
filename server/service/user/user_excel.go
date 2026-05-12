package user

import (
	"fmt"

	"gorm.io/gorm"

	"server/global"
	"server/model"
	"server/service/filesvc"
	"server/utils"
)

// ==================== 用户导入字段定义 ====================

var userImportHeaders = []string{"账号", "用户名称", "性别", "邮箱", "手机号"}

func userImportFields() []utils.ImportField {
	return []utils.ImportField{
		{
			Header:   "账号",
			Key:      "username",
			Required: true,
			Type:     "string",
			MaxLen:   50,
			Validate: func(value string, row int) error {
				if len(value) < 2 {
					return fmt.Errorf("第%d行【账号】长度不能少于2个字符", row)
				}
				return nil
			},
		},
		{
			Header: "用户名称",
			Key:    "nickname",
			Type:   "string",
			MaxLen: 50,
		},
		{
			Header: "性别",
			Key:    "gender",
			Type:   "string",
			Enum:   []string{"男", "女"},
		},
		{
			Header: "邮箱",
			Key:    "email",
			Type:   "string",
			MaxLen: 100,
			Validate: func(value string, row int) error {
				return validateOptionalEmail(value)
			},
		},
		{
			Header: "手机号",
			Key:    "phone",
			Type:   "string",
			MaxLen: 20,
			Validate: func(value string, row int) error {
				return validateOptionalMainlandPhone(value)
			},
		},
	}
}

// ==================== 导入模板 ====================

// GetImportTemplate 生成用户导入模板
func (s *UserService) GetImportTemplate(deptID uint) ([]byte, string, error) {
	deptName := lookupDeptName(deptID)
	sheetName := "用户导入模板"
	filename := "用户导入模板.xlsx"
	if deptName != "" {
		sheetName = deptName + "_用户导入模板"
		filename = deptName + "_用户导入模板.xlsx"
	}

	exporter := utils.NewExcelExporter(sheetName)
	if err := exporter.SetHeaders(userImportHeaders); err != nil {
		return nil, "", err
	}

	// 添加示例数据
	if err := exporter.AddRow([]interface{}{"zhangsan", "张三", "男", "zhangsan@example.com", "13800138000"}); err != nil {
		return nil, "", err
	}

	// 性别下拉（第3列，索引2）
	if err := exporter.AddDataValidation(2, []string{"男", "女"}, 2, 1000); err != nil {
		return nil, "", err
	}

	buf, err := exporter.SaveToBuffer()
	if err != nil {
		return nil, "", err
	}
	return buf, filename, nil
}

// ==================== 导入用户 ====================

// ImportUsersResult 导入用户结果
type ImportUsersResult struct {
	TotalCount   int                 `json:"total_count"`
	SuccessCount int                 `json:"success_count"`
	FailedCount  int                 `json:"failed_count"`
	Errors       []utils.ImportError `json:"errors"`
}

// ImportUsers 导入用户（deptID 由前端选中的部门传入）
func (s *UserService) ImportUsers(operatorID uint, deptID uint, fileData []byte) (*ImportUsersResult, error) {
	if deptID == 0 {
		return nil, fmt.Errorf("请先选择部门再导入")
	}

	importer, err := utils.NewExcelImporter(fileData)
	if err != nil {
		return nil, fmt.Errorf("文件格式错误，请上传xlsx格式的Excel文件")
	}

	// 通用校验
	fields := userImportFields()
	result, err := utils.ValidateImport(importer, fields)
	if err != nil {
		return nil, err
	}

	// 全部失败则直接返回校验结果
	if result.SuccessCount == 0 {
		return &ImportUsersResult{
			TotalCount:   result.TotalCount,
			SuccessCount: 0,
			FailedCount:  result.FailedCount,
			Errors:       result.Errors,
		}, nil
	}

	// 批量唯一性检查：账号
	usernames := make([]string, 0, len(result.Data))
	for _, row := range result.Data {
		if u, ok := row["username"].(string); ok {
			usernames = append(usernames, u)
		}
	}
	existingUsernames, err := findExistingUsernames(usernames)
	if err != nil {
		return nil, fmt.Errorf("检查账号唯一性失败: %w", err)
	}

	// Excel内账号去重
	seenUsernames := make(map[string]bool)

	// 获取默认密码
	defaultPassword := s.managedUserDefaultPassword()
	hashedPassword, err := utils.HashPassword(defaultPassword)
	if err != nil {
		return nil, fmt.Errorf("生成默认密码失败: %w", err)
	}

	// 获取默认头像
	defaultAvatarFileID, _ := resolveRegisterLogoAvatarFileID()

	// 查找默认角色（code=user）
	var defaultRole model.SysRole
	hasDefaultRole := false
	if err := global.DB.Where("code = ? AND status = 1", "user").First(&defaultRole).Error; err == nil {
		hasDefaultRole = true
	}

	// 逐条入库
	importResult := &ImportUsersResult{
		TotalCount: result.TotalCount,
		Errors:     result.Errors,
	}

	for _, row := range result.Data {
		username := row["username"].(string)

		// 检查Excel内重复
		if seenUsernames[username] {
			importResult.FailedCount++
			importResult.Errors = append(importResult.Errors, utils.ImportError{
				Row:     0,
				Column:  "账号",
				Value:   username,
				Message: fmt.Sprintf("账号【%s】在导入文件中重复", username),
			})
			continue
		}
		seenUsernames[username] = true

		// 检查数据库已存在
		if existingUsernames[username] {
			importResult.FailedCount++
			importResult.Errors = append(importResult.Errors, utils.ImportError{
				Row:     0,
				Column:  "账号",
				Value:   username,
				Message: fmt.Sprintf("账号【%s】已存在", username),
			})
			continue
		}

		user := model.SysUser{
			Username:  username,
			Password:  hashedPassword,
			Nickname:  safeString(row["nickname"]),
			Gender:    genderLabelToValue(safeString(row["gender"])),
			Email:     safeString(row["email"]),
			Phone:     safeString(row["phone"]),
			Status:    1,
			DeptID:    deptID,
			CreatedBy: operatorID,
		}

		if defaultAvatarFileID > 0 {
			user.AvatarFileID = defaultAvatarFileID
		}

		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			createQuery := tx
			if user.AvatarFileID == 0 {
				createQuery = createQuery.Omit("AvatarFileID")
			}
			if err := createQuery.Create(&user).Error; err != nil {
				return err
			}

			if hasDefaultRole {
				if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
					return err
				}
			}

			if user.AvatarFileID > 0 {
				if err := filesvc.Reference.ReplaceRefs(tx, "sys_user", user.ID, []filesvc.FileRef{{
					FileID: user.AvatarFileID,
					Field:  "avatar",
				}}); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			importResult.FailedCount++
			importResult.Errors = append(importResult.Errors, utils.ImportError{
				Row:     0,
				Column:  "账号",
				Value:   username,
				Message: fmt.Sprintf("创建用户【%s】失败: %s", username, err.Error()),
			})
			continue
		}

		importResult.SuccessCount++
	}

	importResult.FailedCount += result.FailedCount
	return importResult, nil
}

// ==================== 导出用户 ====================

// ExportUsers 导出指定部门的用户列表，支持可选 userIDs 选择性导出
func (s *UserService) ExportUsers(operatorID uint, deptID uint, userIDs []uint) ([]byte, string, error) {
	if deptID == 0 {
		return nil, "", fmt.Errorf("请先选择部门再导出")
	}

	var users []model.SysUser
	db := global.DB.Where("dept_id = ?", deptID)
	if len(userIDs) > 0 {
		db = db.Where("id IN ?", userIDs)
	}
	if err := db.Order("id ASC").Find(&users).Error; err != nil {
		return nil, "", err
	}

	deptName := lookupDeptName(deptID)
	sheetName := "用户数据"
	filename := "用户导出.xlsx"
	if deptName != "" {
		sheetName = deptName + "_用户数据"
		filename = deptName + "_用户导出.xlsx"
	}

	exporter := utils.NewExcelExporter(sheetName)
	headers := []string{"账号", "用户名称", "性别", "邮箱", "手机号"}
	if err := exporter.SetHeaders(headers); err != nil {
		return nil, "", err
	}

	for _, u := range users {
		row := []interface{}{
			u.Username,
			u.Nickname,
			genderValueToLabel(u.Gender),
			u.Email,
			u.Phone,
		}
		if err := exporter.AddRow(row); err != nil {
			return nil, "", err
		}
	}

	buf, err := exporter.SaveToBuffer()
	if err != nil {
		return nil, "", err
	}

	return buf, filename, nil
}

// ==================== 辅助函数 ====================

func lookupDeptName(deptID uint) string {
	if deptID == 0 {
		return ""
	}
	var dept model.SysDept
	if err := global.DB.Select("name").Where("id = ?", deptID).First(&dept).Error; err != nil {
		return ""
	}
	return dept.Name
}

func findExistingUsernames(usernames []string) (map[string]bool, error) {
	if len(usernames) == 0 {
		return map[string]bool{}, nil
	}
	var existing []string
	if err := global.DB.Model(&model.SysUser{}).Where("username IN ?", usernames).Pluck("username", &existing).Error; err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(existing))
	for _, u := range existing {
		result[u] = true
	}
	return result, nil
}

func safeString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func genderLabelToValue(label string) int {
	switch label {
	case "男":
		return 0
	case "女":
		return 1
	default:
		return 0
	}
}

func genderValueToLabel(value int) string {
	switch value {
	case 0:
		return "男"
	case 1:
		return "女"
	default:
		return "男"
	}
}
