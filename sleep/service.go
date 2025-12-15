package main

import (
	"fmt"

	"gorm.io/gorm"
)

func CreateSleepUser(db *gorm.DB, name string) (*SleepUser, error) {

	sleepuser := SleepUser{Name: name}
	if err := db.Create(&sleepuser).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	return &sleepuser, nil
}

func AddReason(db *gorm.DB, userID uint, reason string) error {
	reasonObj := &Reason{
		UserID: userID,
		Reason: reason,
	}
	result := db.Create(reasonObj)
	return result.Error
}
