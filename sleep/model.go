package main

import "gorm.io/gorm"

type SleepUser struct {
	gorm.Model
	Name    string
	Reasons []Reason `gorm:"foreignKey:UserID"` // 添加外键标签
}

type Reason struct {
	gorm.Model
	UserID uint   // 外键，指向用户
	Reason string // 不学习的理由
}
