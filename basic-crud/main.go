package main

import (
	"fmt"
	"gorm-learning/common" // 导入公共包W
	"gorm.io/gorm"
)

func main() {
	db := common.InitDB()
	// 自动迁移表结构
	db.AutoMigrate(&User{})

	// 1. 创建操作测试
	//testCreate(db)
	// 2. 查询操作测试
	//testQuery(db)
	// 3. 更新操作测试
	//testUpdate(db)
	// 4. 删除操作测试
	testDelete(db)
}

// 每个测试逻辑拆成单独函数（清晰）
func testCreate(db *gorm.DB) { /* 创建测试代码 */
	// 3. 测试1：未传Age，是否自动填充默认值18
	fmt.Println("=== 测试1：未传Age，验证默认值 ===")
	user1 := User{
		Name:  "张三",
		Email: "zhangsan@test.com",
	}
	result := db.Create(&user1) // 创建用户
	if result.Error != nil {
		fmt.Println("创建失败：", result.Error)
	} else {
		fmt.Printf("创建成功，用户ID：%d，Age：%d\n", user1.ID, user1.Age) // 查看Age是否为18
	}

	// 4. 测试2：重复插入相同Email，是否触发唯一索引报错
	fmt.Println("\n=== 测试2：重复Email，验证唯一索引 ===")
	user2 := User{
		Name:  "李四",
		Email: "zhangsan@test.com", // 和user1的Email重复
	}
	result = db.Create(&user2)
	if result.Error != nil {
		fmt.Println("预期报错：", result.Error) // 应该触发唯一索引冲突错误
	} else {
		fmt.Println("创建成功（不符合预期）")
	}
}
func testQuery(db *gorm.DB) {
	// 1. 方式1：主键查询（db.First(&user, 1)）
	fmt.Println("=== 方式1：主键查询（db.First(&user, 1)）===")
	var user1 User
	db.First(&user1, 1)                 // 根据主键ID=1查询
	fmt.Printf("user1: %+v\n\n", user1) // 打印所有字段

	// 2. 方式2：条件查询（db.Where("name = ?", "张三").Take(&user)）
	fmt.Println("=== 方式2：条件查询（Where+Take）===")
	var user2 User
	db.Where("name = ?", "张三").Take(&user2) // 根据姓名条件查询
	fmt.Printf("user2: %+v\n\n", user2)

	// 3. 方式3：指定字段查询（Select+Find）
	fmt.Println("=== 方式3：指定字段查询（Select+Find）===")
	var user3 User
	db.Model(&User{}).Select("name", "age").Find(&user3, "id = ?", 1) // 只查name和age字段
	fmt.Printf("user3: %+v\n", user3)
	fmt.Println("未指定的字段（如ID/Email/CreatedAt）为类型零值：", user3.ID, user3.Email, user3.CreatedAt)

}
func testUpdate(db *gorm.DB) { /* 更新测试代码 */

	// 先创建测试用户（ID=1）
	var user User
	if err := db.First(&user, 1).Error; err != nil {
		user = User{Name: "张三", Age: 20, Email: "zhangsan@test.com"}
		db.Create(&user)
		fmt.Println("已创建测试用户：", user)
	}

	fmt.Println("\n=== 测试1：db.Save(&user) 更新Age为0 ===")
	user.Age = 0
	db.Save(&user)
	db.First(&user, 1)
	fmt.Printf("更新后Age：%d（预期0），UpdatedAt：%s\n", user.Age, user.UpdatedAt)

	fmt.Println("\n=== 测试2：Update vs UpdateColumn 对比 ===")
	oldUpdatedAt := user.UpdatedAt
	// Update更新
	db.Model(&user).Update("age", 5)
	db.First(&user, 1)
	fmt.Printf("Update后Age：%d，UpdatedAt是否变化：%v\n", user.Age, user.UpdatedAt != oldUpdatedAt)

	// UpdateColumn更新
	oldUpdatedAt = user.UpdatedAt
	db.Model(&user).UpdateColumn("age", 8)
	db.First(&user, 1)
	fmt.Printf("UpdateColumn后Age：%d，UpdatedAt是否变化：%v\n", user.Age, user.UpdatedAt != oldUpdatedAt)

	fmt.Println("\n=== 测试3：Updates非零值更新 ===")
	db.Model(&user).Updates(User{Name: "李四", Age: 0})
	db.First(&user, 1)
	fmt.Printf("更新后Name：%s（预期李四），Age：%d（预期8，0被忽略）\n", user.Name, user.Age)
}
func testDelete(db *gorm.DB) { /* 删除测试代码 */
	// 创建测试用户（确保有数据可删）
	var user User
	if err := db.First(&user, 1).Error; err != nil {
		user = User{Name: "张三", Age: 20, Email: "zhangsan@test.com"}
		db.Create(&user)
		fmt.Println("已创建测试用户：", user)
	}

	fmt.Println("\n=== 测试1：软删除 db.Delete(&user) ===")
	// 执行软删除
	db.Delete(&user)
	// 普通查询（自动过滤软删除数据）
	var softDelUser User
	normalQuery := db.First(&softDelUser, 1)
	fmt.Printf("普通查询是否找到：%v（预期未找到）\n", normalQuery.Error == nil)
	// 带Unscoped查询（显示软删除数据）
	var unscopedUser User
	db.Unscoped().First(&unscopedUser, 1)
	fmt.Printf("软删除后数据是否存在：%v，IsDeleted标记：%d\n", unscopedUser.ID != 0, unscopedUser.IsDeleted)

	fmt.Println("\n=== 测试2：物理删除 db.Unscoped().Delete(&user) ===")
	// 先恢复软删除数据（用于测试物理删除）
	db.Unscoped().Model(&user).Update("is_deleted", 0)
	// 执行物理删除
	db.Unscoped().Delete(&user)
	// 再次查询（包括软删除数据）
	var physicalDelUser User
	physicalQuery := db.Unscoped().First(&physicalDelUser, 1)
	fmt.Printf("物理删除后是否找到数据：%v（预期未找到）\n", physicalQuery.Error == nil)

}
