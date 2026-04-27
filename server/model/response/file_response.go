package response

type FileMigrationItem struct {
	FileID            uint   `json:"file_id"`
	FileName          string `json:"file_name"`
	SourceStorageType string `json:"source_storage_type"`
	TargetStorageType string `json:"target_storage_type"`
	OldURL            string `json:"old_url"`
	NewURL            string `json:"new_url"`
	Action            string `json:"action"`
	Message           string `json:"message"`
}

type FileMigrationResult struct {
	TargetStorageType  string              `json:"target_storage_type"`
	TotalCount         int                 `json:"total_count"`
	TotalSize          int64               `json:"total_size"`
	PendingCount       int                 `json:"pending_count"`
	PendingSize        int64               `json:"pending_size"`
	SkippedCount       int                 `json:"skipped_count"`
	SkippedSize        int64               `json:"skipped_size"`
	ConflictCount      int                 `json:"conflict_count"`
	ConflictSize       int64               `json:"conflict_size"`
	MissingSourceCount int                 `json:"missing_source_count"`
	MissingSourceSize  int64               `json:"missing_source_size"`
	MigratedCount      int                 `json:"migrated_count"`
	FailedCount        int                 `json:"failed_count"`
	WarningCount       int                 `json:"warning_count"`
	Items              []FileMigrationItem `json:"items"`
}

type FileMigrationTaskStatus struct {
	TaskID            string              `json:"task_id"`
	Status            string              `json:"status"`
	Message           string              `json:"message"`
	Scope             string              `json:"scope"`
	SourceStorageType string              `json:"source_storage_type"`
	TargetStorageType string              `json:"target_storage_type"`
	TotalCount        int                 `json:"total_count"`
	TotalSize         int64               `json:"total_size"`
	PendingCount      int                 `json:"pending_count"`
	PendingSize       int64               `json:"pending_size"`
	SkippedCount      int                 `json:"skipped_count"`
	SkippedSize       int64               `json:"skipped_size"`
	ConflictCount     int                 `json:"conflict_count"`
	ConflictSize      int64               `json:"conflict_size"`
	MissingSourceCount int                `json:"missing_source_count"`
	MissingSourceSize int64               `json:"missing_source_size"`
	ProcessedCount    int                 `json:"processed_count"`
	ProcessedSize     int64               `json:"processed_size"`
	MigratedCount     int                 `json:"migrated_count"`
	FailedCount       int                 `json:"failed_count"`
	WarningCount      int                 `json:"warning_count"`
	CurrentFileID     uint                `json:"current_file_id"`
	CurrentFileName   string              `json:"current_file_name"`
	StartedAt         string              `json:"started_at"`
	FinishedAt        string              `json:"finished_at"`
	Items             []FileMigrationItem `json:"items"`
}
