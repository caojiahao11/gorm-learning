package main

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Age       int    `gorm:"default:18"`
	Email     string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted int `gorm:"softDelete:flag"` // 软删除字段
}
