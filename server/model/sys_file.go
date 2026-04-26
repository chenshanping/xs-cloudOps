package model

// SysFile 文件记录
type SysFile struct {
	BaseModel
	Name       string      `json:"name" gorm:"size:255;comment:文件名称"`
	Path       string      `json:"path" gorm:"size:500;comment:文件路径"`
	URL        string      `json:"url" gorm:"size:500;comment:访问URL"`
	Size       int64       `json:"size" gorm:"comment:文件大小(字节)"`
	Ext        string      `json:"ext" gorm:"size:20;comment:文件扩展名"`
	MimeType   string      `json:"mime_type" gorm:"size:100;comment:MIME类型"`
	MD5        string      `json:"md5" gorm:"size:32;index;comment:文件MD5"`
	StorageID  uint        `json:"storage_id" gorm:"comment:存储配置ID"`
	Storage    *SysStorage `json:"storage" gorm:"foreignKey:StorageID;references:ID"`
	UploaderID uint        `json:"uploader_id" gorm:"comment:上传者ID"`
	Status     int         `json:"status" gorm:"default:1;comment:状态 0删除 1正常"`
}

func (SysFile) TableName() string {
	return "sys_file"
}

// SysFileChunk 文件分片记录（用于断点续传）
type SysFileChunk struct {
	BaseModel
	UploadID    string `json:"upload_id" gorm:"size:100;index;comment:上传ID"`
	FileHash    string `json:"file_hash" gorm:"size:32;index;comment:文件MD5"`
	ChunkIndex  int    `json:"chunk_index" gorm:"comment:分片索引"`
	ChunkSize   int64  `json:"chunk_size" gorm:"comment:分片大小"`
	ChunkHash   string `json:"chunk_hash" gorm:"size:32;comment:分片MD5"`
	StorageID   uint   `json:"storage_id" gorm:"comment:存储配置ID"`
	StoragePath string `json:"storage_path" gorm:"size:500;comment:分片存储路径"`
	Status      int    `json:"status" gorm:"default:0;comment:状态 0上传中 1已完成"`
}

func (SysFileChunk) TableName() string {
	return "sys_file_chunk"
}
