package main

import (
	"fmt"
	"gorm-learning/common" // 导入公共包W
)

func main() {
	db := common.InitDB()
	// 应该同时迁移两个模型
	db.AutoMigrate(&SleepUser{}, &Reason{})

	user, err := CreateSleepUser(db, "曹佳豪")
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		return
	}
	// 为用户添加理由
	err = AddReason(db, user.ID, "太困了不想学习")
	if err != nil {
		fmt.Printf("添加理由失败: %v\n", err)
		return
	}
	// 为用户添加理由
	err = AddReason(db, user.ID, "太困了不想学习")
	if err != nil {
		fmt.Printf("添加理由失败: %v\n", err)
		return
	}

}
