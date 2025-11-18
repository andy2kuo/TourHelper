package dao

import (
	"sync"

	"gorm.io/gorm"
)

// DAO 集中管理所有 DAO 實例
type DAO struct {
	User        UserDAO
	Destination DestinationDAO
	// 未來可以新增其他 DAO，例如：
	// Tag         TagDAO
	// Preference  PreferenceDAO
}

var (
	instance *DAO
	once     sync.Once
	db       *gorm.DB
)

// SetDB 設定資料庫連線
func SetDB(database *gorm.DB) {
	db = database
}

// Get 取得 DAO 實例（單例模式）
func Get() *DAO {
	once.Do(func() {
		instance = &DAO{
			User:        NewUserDAO(db),
			Destination: NewDestinationDAO(db),
			// 初始化其他 DAO
		}
	})
	return instance
}
