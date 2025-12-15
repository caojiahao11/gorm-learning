package main

type User struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Email  string
	Wallet Wallet `gorm:"foreignKey:UserID"`
	Role   Role   `gorm:"foreignKey:UserID"`
}

type Wallet struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint
	Balance float64
}

type Role struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	Name   string
}
