package main

import (
	"gorm-learning/common"
	"log"
)

func main() {
	db := common.InitDB()
	db.AutoMigrate(&User{}, &Wallet{}, &Role{})

	err := RegisterUser(db, "Alice", "alice@example.com", -50)
	if err != nil {
		log.Println("Registration failed:", err)
	} else {
		log.Println("Registration succeeded")
	}
}
