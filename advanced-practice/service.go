package main

import (
	"fmt"

	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, name, email string, initialBalance float64) error {
	return db.Transaction(func(tx *gorm.DB) error {

		user := &User{Name: name, Email: email}
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		wallet := &Wallet{UserID: user.ID, Balance: initialBalance}
		if err := tx.Create(wallet).Error; err != nil {
			return fmt.Errorf("failed to create wallet: %w", err)
		}

		role := &Role{UserID: user.ID, Name: "member"}
		if err := tx.Create(role).Error; err != nil {
			return fmt.Errorf("failed to assign role: %w", err)
		}

		fmt.Println("User registered successfully with ID:", user.ID)
		return nil
	})
}
