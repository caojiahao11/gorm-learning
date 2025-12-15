// relation-query/service.go
package main

import (
	"fmt"
	"log"

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

func CreateUserWithIDCard1(db *gorm.DB, userName string, idCardNumber string) (*User, error) {
	user := User{Name: userName}
	if err := db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败哦: %v", err)
	}
	idCard := IDCard{UserID: user.ID, Number: idCardNumber}
	if err := db.Create(&idCard).Error; err != nil {
		return nil, fmt.Errorf("创建身份证失败: %v", err)
	}
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

// GetUserRoles 查询用户所有角色
func GetUserRoles(db *gorm.DB, userID uint) ([]Role, error) {
	var user User
	if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return user.Roles, nil
}

func queryUsersByRoleName(db *gorm.DB, roleName string) {
	// 关键1：用切片 []Role 接收所有匹配的同名角色
	var roles []Role

	// 关键2：Find(&roles) 查所有同名角色，并预加载每个角色的用户
	if err := db.Preload("Users").Where("name = ?", roleName).Find(&roles).Error; err != nil {
		log.Fatalf("查询角色【%s】失败：%v", roleName, err)
	}

	// 无匹配角色的情况
	if len(roles) == 0 {
		log.Printf("未找到名称为【%s】的角色", roleName)
		return
	}

	// 关键3：汇总所有角色下的用户（用map去重，避免同一用户关联多个同名角色时重复）
	userMap := make(map[uint]User)
	for _, role := range roles {
		for _, user := range role.Users {
			userMap[user.ID] = user // 以用户ID为key，自动去重
		}
	}

	// 输出最终结果
	log.Printf("\n拥有角色【%s】的所有用户（共%d个）：", roleName, len(userMap))
	if len(userMap) == 0 {
		log.Printf("该角色下暂无关联用户")
		return
	}
	for _, user := range userMap {
		log.Printf("用户ID：%d，用户名：%s", user.ID, user.Name)
	}
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
	user4, err := CreateUserWithIDCard(db, "姬如雪", "1234567890")
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
	if err := BatchCreateOrders(db, user4.ID, orderProducts); err != nil {
		return err
	}

	// ========== 1. 先创建测试角色（多对多关联的角色数据） ==========
	roles := []Role{
		{Name: "管理员"},
		{Name: "普通用户"},
		{Name: "VIP用户"},
	}
	if err := db.Create(&roles).Error; err != nil {
		return fmt.Errorf("创建测试角色失败: %v", err)
	}

	// ========== 3. 给用户分配角色（多对多关联） ==========
	// 给user1分配「管理员」和「VIP用户」角色
	if err := db.Model(&user1).Association("Roles").Append(roles[0], roles[2]); err != nil {
		return fmt.Errorf("给李四分配角色失败: %v", err)
	}
	// 给user2分配「普通用户」角色
	if err := db.Model(&user2).Association("Roles").Append(roles[0], roles[1]); err != nil {
		return fmt.Errorf("给张三分配角色失败: %v", err)
	}
	// 给user3分配「VIP用户」角色
	if err := db.Model(&user3).Association("Roles").Append(roles[2]); err != nil {
		return fmt.Errorf("给王五分配角色失败: %v", err)
	}

	fmt.Println("测试数据初始化完成！")
	return nil
}
