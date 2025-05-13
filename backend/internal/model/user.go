package model

import (
	"time"
)

// User 用户模型
type User struct {
	Id           uint64    `xorm:"'id' pk autoincr comment('用户唯一标识符，自增') BIGINT(20) UNSIGNED"`
	Username     string    `xorm:"'username' notnull unique comment('用户名，必须唯一') VARCHAR(255)"`
	PasswordHash string    `xorm:"'password_hash' notnull comment('密码哈希') VARCHAR(255)"`
	Role         string    `xorm:"'role' notnull comment('用户角色') VARCHAR(20)"`
	Status       string    `xorm:"'status' notnull default('enabled') comment('用户状态') ENUM('enabled','disabled')"`
	Nickname     string    `xorm:"'nickname' null comment('用户昵称') VARCHAR(255)"`
	Email        string    `xorm:"'email' null comment('用户电子邮件') VARCHAR(255)"`
	AvatarUrl    string    `xorm:"'avatar_url' null comment('用户头像URL') VARCHAR(1024)"`
	CreatedAt    time.Time `xorm:"'created_at' notnull created comment('记录创建时间')"`
	UpdatedAt    time.Time `xorm:"'updated_at' notnull updated comment('记录更新时间')"`

	// 关联字段（不映射到数据库）
	TodoItems []*TodoItem `xorm:"-"`
}

// TableName 返回表名
func (u *User) TableName() string {
	return "users_tab"
}

// BeforeInsert 插入前的钩子
func (u *User) BeforeInsert() {
	if u.Status == "" {
		u.Status = "enabled"
	}
}
