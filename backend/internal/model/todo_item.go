package model

import (
	"time"
)

// TodoItem 待办事项模型
type TodoItem struct {
	Id          uint64    `xorm:"'id' pk autoincr comment('自增主键') BIGINT(20) UNSIGNED"`
	UserId      uint64    `xorm:"'user_id' comment('外键，引用users表的id') index BIGINT(20) UNSIGNED"`
	Title       string    `xorm:"'title' notnull comment('任务标题') VARCHAR(255)"`
	Description string    `xorm:"'description' null comment('任务描述') TEXT"`
	Priority    string    `xorm:"'priority' notnull default('medium') comment('优先级') ENUM('low','medium','high')"`
	DueDate     time.Time `xorm:"'due_date' null comment('截止日期') DATETIME"`
	Completed   bool      `xorm:"'completed' notnull default(false) comment('是否完成') BOOL"`
	CreatedAt   time.Time `xorm:"'created_at' notnull created comment('创建时间')"`
	UpdatedAt   time.Time `xorm:"'updated_at' notnull updated comment('更新时间')"`
	EnableFlag  bool      `xorm:"'enable_flag' notnull default(true) comment('启用标志，true为启用，false为禁用') BOOL"`
	// 关联字段（不映射到数据库）
	User *User `xorm:"-"`
}

// TableName 返回表名
func (t *TodoItem) TableName() string {
	return "todo_items_tab"
}

// BeforeInsert 插入前的钩子
func (t *TodoItem) BeforeInsert() {
	if t.Priority == "" {
		t.Priority = "medium"
	}
	t.EnableFlag = true
}
