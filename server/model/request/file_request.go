package request

type FileMigrationRequest struct {
	Scope             string               `json:"scope"`
	IDs               []uint               `json:"ids"`
	SourceStorageType string               `json:"source_storage_type" binding:"required"`
	TargetStorageType string               `json:"target_storage_type" binding:"required"`
	Filters           FileMigrationFilters `json:"filters"`
}

type FileMigrationFilters struct {
	Name string `json:"name"`
	Ext  string `json:"ext"`
}
