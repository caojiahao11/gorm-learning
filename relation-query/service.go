// relation-query/service.go
package main

import (
	"fmt"
	"gorm.io/gorm"
)

// -------------------------- 数据创建层 --------------------------
// CreateUserWithIDCard 创建用户并关联身份证（一对一）
func CreateUserWithIDCard(db *gorm.DB, userName string, idCardNumber string) (*User, error) {
	user := User{Name: userName}
	if err := db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	idCard := IDCard{UserID: user.ID, Number: idCardNumber}
	if err := db.Create(&idCard).Error; err != nil {
		return nil, fmt.Errorf("创建身份证失败: %v", err)
	}

	user.IDCard = idCard
	return &user, nil
}

// BatchCreateOrders 批量给用户创建订单（一对多）
func BatchCreateOrders(db *gorm.DB, userID uint, products []map[string]float64) error {
	var orders []Order
	for _, p := range products {
		for productName, price := range p {
			orders = append(orders, Order{
				UserID:  userID,
				Product: productName,
				Price:   price,
			})
		}
	}
	return AddOrderToUser(db, userID, orders)
}

// AddOrderToUser 给用户添加订单（基础方法）
func AddOrderToUser(db *gorm.DB, userID uint, orders []Order) error {
	for i := range orders {
		orders[i].UserID = userID
	}
	return db.Create(&orders).Error
}

// -------------------------- 数据查询层 --------------------------
// GetUserWithIDCard 查询用户及身份证信息
func GetUserWithIDCard(db *gorm.DB, userID uint) (*User, error) {
	var user User
	if err := db.Preload("IDCard").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return &user, nil
}

// GetUserWithLatestOrders 查询用户及最近N条订单
func GetUserWithLatestOrders(db *gorm.DB, userID uint, limit int) (*User, error) {
	var user User
	if err := db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(limit)
	}).First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("查询用户订单失败: %v", err)
	}
	return &user, nil
}

// GetAllUsersWithRelations 查询所有用户及关联数据（身份证+订单）
func GetAllUsersWithRelations(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Preload("IDCard").Preload("Orders").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询所有用户失败: %v", err)
	}
	return users, nil
}

// InitTestData 初始化测试数据
func InitTestData(db *gorm.DB) error {
	// 创建3个测试用户（带身份证）
	user1, err := CreateUserWithIDCard(db, "李四", "3206999999999")
	if err != nil {
		return err
	}
	user2, err := CreateUserWithIDCard(db, "张三", "110101199001011234")
	if err != nil {
		return err
	}
	user3, err := CreateUserWithIDCard(db, "王五", "440101199505056789")
	if err != nil {
		return err
	}

	// 定义测试订单数据
	orderProducts := []map[string]float64{
		{"macbook": 12000},
		{"book": 12},
		{"apple": 120},
	}

	// 给3个用户批量创建订单
	if err := BatchCreateOrders(db, user1.ID, orderProducts); err != nil {
		return err
	}
	if err := BatchCreateOrders(db, user2.ID, orderProducts); err != nil {
		return err
	}
	if err := BatchCreateOrders(db, user3.ID, orderProducts); err != nil {
		return err
	}

	fmt.Println("测试数据初始化完成！")
	return nil
}
