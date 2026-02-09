package model

import "time"

type SysLoginLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"comment:用户ID"`
	Username  string    `json:"username" gorm:"size:50;comment:用户名"`
	IP        string    `json:"ip" gorm:"size:50;comment:IP地址"`
	Location  string    `json:"location" gorm:"size:100;comment:登录地点"`
	Browser   string    `json:"browser" gorm:"size:100;comment:浏览器"`
	OS        string    `json:"os" gorm:"size:100;comment:操作系统"`
	Status    int       `json:"status" gorm:"comment:状态 1成功 0失败"`
	Msg       string    `json:"msg" gorm:"size:255;comment:消息"`
	CreatedAt time.Time `json:"created_at"`
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}
