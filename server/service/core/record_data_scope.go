package core

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// RecordDataScopeBinding defines how a generic business table maps ownership fields.
type RecordDataScopeBinding struct {
	TableAlias      string
	DeptColumn      string
	CreatedByColumn string
	SelfColumn      string
}

// ApplyRecordDataScope applies a generic record-level data scope for business tables.
func ApplyRecordDataScope(db *gorm.DB, scope *UserDataScope, binding RecordDataScopeBinding) *gorm.DB {
	if scope == nil {
		return db.Where("1 = 0")
	}
	if scope.All {
		return db
	}

	conditions := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)

	if len(scope.DeptIDs) > 0 && binding.DeptColumn != "" {
		conditions = append(conditions, fmt.Sprintf("%s IN ?", binding.qualify(binding.DeptColumn)))
		args = append(args, scope.DeptIDs)
	}
	if len(scope.CreatorIDs) > 0 && binding.CreatedByColumn != "" {
		conditions = append(conditions, fmt.Sprintf("%s IN ?", binding.qualify(binding.CreatedByColumn)))
		args = append(args, scope.CreatorIDs)
	}
	if scope.AllowSelf {
		if selfColumn := binding.resolveSelfColumn(); selfColumn != "" {
			conditions = append(conditions, fmt.Sprintf("%s = ?", binding.qualify(selfColumn)))
			args = append(args, scope.OperatorID)
		}
	}

	if len(conditions) == 0 {
		return db.Where("1 = 0")
	}

	return db.Where("("+strings.Join(conditions, " OR ")+")", args...)
}

func (binding RecordDataScopeBinding) resolveSelfColumn() string {
	if binding.SelfColumn != "" {
		return binding.SelfColumn
	}
	return binding.CreatedByColumn
}

func (binding RecordDataScopeBinding) qualify(column string) string {
	if binding.TableAlias == "" || strings.Contains(column, ".") {
		return column
	}
	return binding.TableAlias + "." + column
}
