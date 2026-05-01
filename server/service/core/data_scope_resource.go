package core

type DataScopeResource struct {
	Code        string   `json:"code"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	OwnerFields []string `json:"owner_fields"`
}

const (
	DataScopeResourceUserManagement = "system:user-management"
	DataScopeResourceDeptManagement = "system:dept-management"

	DataScopeOwnerFieldDeptID    = "dept_id"
	DataScopeOwnerFieldCreatedBy = "created_by"
)

var supportedDataScopeResources = []DataScopeResource{
	{
		Code:        DataScopeResourceUserManagement,
		Label:       "用户管理",
		Description: "用户管理数据权限资源，支持按部门和创建人限定访问范围。",
		OwnerFields: []string{DataScopeOwnerFieldDeptID, DataScopeOwnerFieldCreatedBy},
	},
	{
		Code:        DataScopeResourceDeptManagement,
		Label:       "部门管理",
		Description: "部门管理数据权限资源，支持按部门范围限定可管理部门。",
		OwnerFields: []string{DataScopeOwnerFieldDeptID},
	},
}

func SupportedDataScopeResources() []DataScopeResource {
	resources := make([]DataScopeResource, len(supportedDataScopeResources))
	for i, resource := range supportedDataScopeResources {
		resources[i] = cloneDataScopeResource(resource)
	}
	return resources
}

func IsSupportedDataScopeResource(resourceCode string) bool {
	for _, resource := range supportedDataScopeResources {
		if resource.Code == resourceCode {
			return true
		}
	}
	return false
}

func cloneDataScopeResource(resource DataScopeResource) DataScopeResource {
	cloned := resource
	if resource.OwnerFields != nil {
		cloned.OwnerFields = make([]string, len(resource.OwnerFields))
		copy(cloned.OwnerFields, resource.OwnerFields)
	}
	return cloned
}
