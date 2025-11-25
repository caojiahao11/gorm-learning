markdown
# GORM 学习作业仓库

这是一个用于学习 GORM（Go 语言 ORM 框架）的作业代码仓库，包含基础 CRUD、关联查询等核心功能的实践代码。


## 项目结构
gorm-learning /
├── basic-crud/ # 基础 CRUD 操作作业（创建 / 查询 / 更新 / 删除）
├── relation-query/ # 关联查询作业（待补充）
├── advanced-practice/ # 进阶实战作业（待补充）
├── performance-debug/ # 性能与调试作业（待补充）
├── common/ # 公共工具（数据库连接等）│
    └── db.go
├── go.mod # Go 模块配置
└── README.md # 项目说明


## 环境依赖
- Go 1.18+
- MySQL 5.7+/8.0+
- GORM v2：`gorm.io/gorm`
- MySQL 驱动：`gorm.io/driver/mysql`


## 运行方法
1. 克隆仓库：
   ```bash
   git clone https://github.com/caojiahao11/gorm-learning.git
   cd gorm-learning
配置数据库：
修改 common/db.go 中的 dsn，替换为你的 MySQL 账号、密码和数据库名：
go
运行
dsn := "root:你的密码@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
确保 MySQL 中已创建 test_db 数据库。
运行指定作业（以基础 CRUD 为例）：
bash
运行
cd basic-crud
go run main.go
作业内容说明
basic-crud：
创建操作：字段映射、默认值、唯一索引测试
查询操作：主键查、条件查、指定字段查对比
更新操作：零值更新、Update vs UpdateColumn、非零值更新
删除操作：软删除与物理删除对比
relation-query（待补充）：
一对一、一对多、多对多关联查询实践
advanced-practice（待补充）：
事务、钩子函数、批量操作等进阶功能
performance-debug（待补充）：
SQL 日志、性能优化、错误调试技巧
plaintext
