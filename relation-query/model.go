package main

import "gorm.io/gorm"

//1.一对一：用户与身份证
//- 定义 IDCard 结构体（ID/UserID/Number），与 User 关联（User has one IDCard）。
//- 实现：
//
//- 创建用户时同时创建身份证（db.Create(&user).Association("IDCard").Append(&idCard)）。
//- 查询用户并预加载身份证（Preload("IDCard")）。

type User struct {
	gorm.Model
	Name   string
	IDCard IDCard
	Orders []Order
	Roles  []Role `gorm:"many2many:user_role"`
}

// 一对一
type IDCard struct {
	gorm.Model
	Number string
	UserID uint
}

// 一对多
type Order struct {
	gorm.Model
	UserID  uint
	Product string
	Price   float64
}

type Role struct {
	gorm.Model
	Name  string
	Users []User `gorm:"many2many:user_role;"`
}
