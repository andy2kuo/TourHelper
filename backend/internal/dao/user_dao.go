package dao

import (
	"gorm.io/gorm"
)

// UserDAO 使用者資料庫操作介面
type UserDAO interface {
	// TODO: 實作 CRUD 方法
}

// userDAO 使用者資料庫操作實作
type userDAO struct {
	db *gorm.DB
}

// NewUserDAO 建立使用者 DAO
func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAO{db: db}
}
