// relation-query/main.go
package main

import (
	"fmt"
	"gorm-learning/common"

	"gorm.io/gorm"
)

func main() {
	// 1. 初始化数据库（调用公共模块）
	db := common.InitDB()
	// 自动迁移模型
	db.AutoMigrate(&User{}, &IDCard{}, &Order{}, &Role{})

	// 2. 初始化测试数据
	if err := InitTestData(db); err != nil {
		panic(fmt.Sprintf("初始化测试数据失败: %v", err))
	}

	// 3. 测试查询功能
	testQuery(db)
}

// testQuery 测试查询逻辑
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
		// 查询所有用的所有角色
		roles, _ := GetUserRoles(db, user.ID)
		fmt.Println("用户角色:", roles) // 输出：[{1 管理员} {2 编辑}]
	}

	queryUsersByRoleName(db, "管理员")
}
