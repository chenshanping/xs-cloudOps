package core

type DataScopeResource struct {
	Code  string
	Label string
}

const (
	DataScopeResourceUserManagement = "system:user-management"
	DataScopeResourceDeptManagement = "system:dept-management"
)

var supportedDataScopeResources = []DataScopeResource{
	{Code: DataScopeResourceUserManagement, Label: "用户管理"},
	{Code: DataScopeResourceDeptManagement, Label: "部门管理"},
}

func SupportedDataScopeResources() []DataScopeResource {
	resources := make([]DataScopeResource, len(supportedDataScopeResources))
	copy(resources, supportedDataScopeResources)
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
