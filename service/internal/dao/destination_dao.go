package dao

import (
	"gorm.io/gorm"
)

// DestinationDAO 景點資料庫操作介面
type DestinationDAO interface {
	// TODO: 實作 CRUD 方法
}

// destinationDAO 景點資料庫操作實作
type destinationDAO struct {
	db *gorm.DB
}

// NewDestinationDAO 建立景點 DAO
func NewDestinationDAO(db *gorm.DB) DestinationDAO {
	return &destinationDAO{db: db}
}
