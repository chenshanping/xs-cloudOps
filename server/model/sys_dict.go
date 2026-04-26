package model

// SysDictType 字典类型
type SysDictType struct {
	BaseModel
	Name   string `json:"name" gorm:"size:100;not null;comment:字典名称"`
	Type   string `json:"type" gorm:"size:100;not null;uniqueIndex;comment:字典类型"`
	Status int    `json:"status" gorm:"default:1;comment:状态(1:正常,0:停用)"`
	Remark string `json:"remark" gorm:"size:500;comment:备注"`
}

func (SysDictType) TableName() string {
	return "sys_dict_type"
}

// SysDictData 字典数据
type SysDictData struct {
	BaseModel
	DictType  string `json:"dict_type" gorm:"size:100;not null;index;comment:字典类型"`
	Label     string `json:"label" gorm:"size:100;not null;comment:字典标签"`
	Value     string `json:"value" gorm:"size:100;not null;comment:字典键值"`
	Sort      int    `json:"sort" gorm:"default:0;comment:排序"`
	Status    int    `json:"status" gorm:"default:1;comment:状态(1:正常,0:停用)"`
	TagType   string `json:"tag_type" gorm:"size:50;comment:标签类型(success/info/warning/error)"`
	IsDefault int    `json:"is_default" gorm:"default:0;comment:是否默认(1:是,0:否)"`
	Remark    string `json:"remark" gorm:"size:500;comment:备注"`
}

func (SysDictData) TableName() string {
	return "sys_dict_data"
}
