package role

import (
	"strings"
	"testing"

	"server/model"
	"server/model/request"
)

func TestNormalizeRoleFeatureDataScopeAssignmentsRejectsUnsupportedResource(t *testing.T) {
	assignments := []request.RoleFeatureDataScopeAssignment{
		{
			ResourceCode: "biz:unsupported-resource",
			DataScope:    model.DataScopeSelf,
		},
	}

	_, err := normalizeRoleFeatureDataScopeAssignments(assignments)
	if err == nil {
		t.Fatal("normalizeRoleFeatureDataScopeAssignments should reject unsupported resource_code")
	}
	if !strings.Contains(err.Error(), "不支持的数据权限资源") {
		t.Fatalf("unexpected error: %v", err)
	}
}
