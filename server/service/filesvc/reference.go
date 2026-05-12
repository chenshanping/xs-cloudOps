package filesvc

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"

	"server/model"
)

type FileRef struct {
	FileID uint
	Field  string
}

type ReferenceService struct{}

var Reference = &ReferenceService{}

func (s *ReferenceService) ReplaceRefs(tx *gorm.DB, refTable string, refID uint, refs []FileRef) error {
	if err := s.validateTarget(refTable, refID); err != nil {
		return err
	}

	if err := s.ClearRefs(tx, refTable, refID); err != nil {
		return err
	}

	normalized := normalizeFileRefs(refs)
	if len(normalized) == 0 {
		return nil
	}

	rows := make([]model.SysFileReference, 0, len(normalized))
	for _, ref := range normalized {
		rows = append(rows, model.SysFileReference{
			FileID:   ref.FileID,
			RefTable: refTable,
			RefID:    refID,
			RefField: ref.Field,
		})
	}

	return tx.Create(&rows).Error
}

func (s *ReferenceService) ClearRefs(tx *gorm.DB, refTable string, refID uint) error {
	if err := s.validateTarget(refTable, refID); err != nil {
		return err
	}
	return tx.Where("ref_table = ? AND ref_id = ?", refTable, refID).Delete(&model.SysFileReference{}).Error
}

func (s *ReferenceService) GetReferenceCounts(tx *gorm.DB, fileIDs []uint) (map[uint]int64, error) {
	result := make(map[uint]int64, len(fileIDs))
	if len(fileIDs) == 0 {
		return result, nil
	}

	type countRow struct {
		FileID uint  `gorm:"column:file_id"`
		Count  int64 `gorm:"column:count"`
	}

	var rows []countRow
	if err := tx.Model(&model.SysFileReference{}).
		Select("file_id, COUNT(*) AS count").
		Where("file_id IN ?", normalizeIDs(fileIDs)).
		Group("file_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.FileID] = row.Count
	}
	return result, nil
}

func (s *ReferenceService) EnsureFileNotReferenced(tx *gorm.DB, fileID uint) error {
	var ref model.SysFileReference
	if err := tx.Where("file_id = ?", fileID).Order("id ASC").First(&ref).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return fmt.Errorf("文件正在被引用：%s，无法删除", describeReference(ref))
}

func (s *ReferenceService) validateTarget(refTable string, refID uint) error {
	if strings.TrimSpace(refTable) == "" {
		return errors.New("文件引用表名不能为空")
	}
	if refID == 0 {
		return errors.New("文件引用记录ID不能为空")
	}
	return nil
}

func normalizeFileRefs(refs []FileRef) []FileRef {
	seen := make(map[string]struct{}, len(refs))
	result := make([]FileRef, 0, len(refs))
	for _, ref := range refs {
		ref.Field = strings.TrimSpace(ref.Field)
		if ref.FileID == 0 || ref.Field == "" {
			continue
		}

		key := fmt.Sprintf("%s#%d", ref.Field, ref.FileID)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, ref)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Field == result[j].Field {
			return result[i].FileID < result[j].FileID
		}
		return result[i].Field < result[j].Field
	})
	return result
}

func normalizeIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func describeReference(ref model.SysFileReference) string {
	switch ref.RefTable {
	case "sys_user":
		if ref.RefField == "avatar" {
			return "用户头像正在使用"
		}
	case "sys_config":
		if label, ok := configImageRefLabels[ref.RefField]; ok {
			return fmt.Sprintf("系统配置[%s]正在使用", label)
		}
	case "ai_message":
		if ref.RefField == "attachment" {
			return "AI对话附件正在使用"
		}
	}

	return fmt.Sprintf("%s[%d].%s 正在使用", ref.RefTable, ref.RefID, ref.RefField)
}

var configImageRefLabels = map[string]string{
	"sys_logo_file_id":          "系统 Logo",
	"register_logo_file_id":     "注册默认头像",
	"login_bg_image_file_id":    "登录页背景图",
	"slider_captcha_bg_file_id": "滑动验证码背景",
}
