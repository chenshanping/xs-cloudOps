package model

type SysUser struct {
	BaseModel
	Username      string    `json:"username" gorm:"size:50;uniqueIndex;comment:用户名"`
	Password      string    `json:"-" gorm:"size:128;comment:密码"`
	Nickname      string    `json:"nickname" gorm:"size:50;comment:昵称"`
	Email         string    `json:"email" gorm:"size:100;comment:邮箱"`
	Phone         string    `json:"phone" gorm:"size:20;comment:手机号"`
	AvatarFileID  uint      `json:"avatar_file_id" gorm:"comment:头像文件ID"`
	AvatarFile    *SysFile  `json:"-" gorm:"foreignKey:AvatarFileID;references:ID"`
	AvatarFileURL string    `json:"avatar_file_url" gorm:"-"`
	Status        int       `json:"status" gorm:"default:1;comment:状态 1启用 0禁用"`
	DeptID        uint      `json:"dept_id" gorm:"default:0;index;comment:部门ID"`
	Dept          *SysDept  `json:"dept,omitempty" gorm:"foreignKey:DeptID;references:ID"`
	Roles         []SysRole `json:"roles" gorm:"many2many:sys_user_role;"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

// FillAvatarURL 填充头像URL
func (u *SysUser) FillAvatarURL() {
	if u.AvatarFile != nil {
		u.AvatarFileURL = u.AvatarFile.URL
	}
}
