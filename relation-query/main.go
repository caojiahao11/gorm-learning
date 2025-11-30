package main

import (
	"fmt"
	"gorm-learning/common"

	"gorm.io/gorm"
)

// -------------------------- 数据创建层（专门负责新增数据）--------------------------
// CreateUserWithIDCard 创建用户并关联身份证（一对一）
func CreateUserWithIDCard(db *gorm.DB, userName string, idCardNumber string) (*User, error) {
	user := User{Name: userName}
	// 先创建用户
	if err := db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	// 关联身份证
	idCard := IDCard{UserID: user.ID, Number: idCardNumber}
	if err := db.Create(&idCard).Error; err != nil {
		return nil, fmt.Errorf("创建身份证失败: %v", err)
	}

	// 把身份证关联回用户（可选，方便上层使用）
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

// -------------------------- 数据查询层（专门负责查询数据）--------------------------
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

// -------------------------- 测试数据初始化（方便快速生成测试数据）--------------------------
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

// -------------------------- 主函数（负责流程调度）--------------------------
func main() {
	// 1. 初始化数据库
	db := common.InitDB()
	db.AutoMigrate(&User{}, &IDCard{}, &Order{})

	// 2. 初始化测试数据（创建用户+身份证+订单）
	if err := InitTestData(db); err != nil {
		panic(fmt.Sprintf("初始化测试数据失败: %v", err))
	}

	// 3. 测试查询功能
	testQuery(db)
}

// testQuery 专门测试查询逻辑
func testQuery(db *gorm.DB) {
	// 查询用户1的信息（带身份证）
	user1, err := GetUserWithIDCard(db, 1)
	if err != nil {
		fmt.Println("查询用户1失败:", err)
		return
	}
	fmt.Printf("\n用户1信息：\n姓名：%s，身份证：%s\n", user1.Name, user1.IDCard.Number)

	// 查询用户2的最近2条订单
	user2, err := GetUserWithLatestOrders(db, 2, 2)
	if err != nil {
		fmt.Println("查询用户2订单失败:", err)
		return
	}
	fmt.Printf("\n用户2最近2条订单：\n")
	for _, order := range user2.Orders {
		fmt.Printf("- 商品：%s，价格：%.2f\n", order.Product, order.Price)
	}

	// 查询所有用户及关联数据
	allUsers, err := GetAllUsersWithRelations(db)
	if err != nil {
		fmt.Println("查询所有用户失败:", err)
		return
	}
	fmt.Printf("\n所有用户及订单数量：\n")
	for _, user := range allUsers {
		fmt.Printf("- 用户：%s，订单数：%d\n", user.Name, len(user.Orders))
	}
}
