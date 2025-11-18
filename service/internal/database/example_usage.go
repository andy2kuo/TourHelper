package database

// 本檔案提供 DBManager 的使用範例（使用 GORM DBResolver 插件）
// 這只是範例，不會在實際程式中執行

/*
使用範例：

package main

import (
	"log"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/database"
	"github.com/andy2kuo/TourHelper/internal/models"
	"gorm.io/gorm"
)

func main() {
	// 1. 載入設定
	cfg, err := config.Load("TourHelper", "dev", "1.0.0")
	if err != nil {
		log.Fatalf("載入設定失敗: %v", err)
	}

	// 2. 初始化資料庫連線
	if err := database.Init(cfg.Database); err != nil {
		log.Fatalf("初始化資料庫失敗: %v", err)
	}
	defer database.GetInstance().Close()

	// 3. 取得 DBManager 實例
	dbManager := database.GetInstance()

	// ============================================================
	// 基本使用：自動讀寫分離
	// ============================================================

	// 3.1 取得資料庫連線（GORM DBResolver 會自動處理讀寫分離）
	db := dbManager.GetDB() // 預設資料庫（通常是 "main"）

	// 寫入操作會自動使用 Master
	user := &models.User{
		Name:  "張三",
		Email: "zhang@example.com",
	}
	if err := db.Create(user).Error; err != nil {
		log.Printf("建立使用者失敗: %v", err)
	}

	// 讀取操作會自動使用 Slave（如無 Slave 則降級為 Master）
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("查詢使用者失敗: %v", err)
	}

	// ============================================================
	// 明確指定使用 Master 或 Slave
	// ============================================================

	// 4.1 明確使用 Master（強制寫入操作）
	masterDB := dbManager.GetMaster()
	if err := masterDB.Create(&models.Destination{
		Name:     "台北101",
		Location: "台北市信義區",
	}).Error; err != nil {
		log.Printf("建立景點失敗: %v", err)
	}

	// 4.2 明確使用 Slave（強制讀取操作）
	slaveDB := dbManager.GetSlave()
	var destinations []models.Destination
	if err := slaveDB.Where("city = ?", "台北").Find(&destinations).Error; err != nil {
		log.Printf("查詢景點失敗: %v", err)
	}

	// 4.3 指定特定的資料庫（根據名稱）
	mainDB := dbManager.GetDB("main")
	analyticsDB := dbManager.GetDB("analytics")

	// ============================================================
	// 根據 Schema 選擇資料庫
	// ============================================================

	// 5.1 使用 Schema 取得對應的資料庫
	analyticsDB := dbManager.GetBySchema("analytics")
	if err := analyticsDB.Create(&models.Analytics{
		Event: "user_visit",
		Data:  "{}",
	}).Error; err != nil {
		log.Printf("建立分析資料失敗: %v", err)
	}

	// ============================================================
	// 進階使用範例
	// ============================================================

	// 6.1 事務操作（自動使用 Master）
	err = db.Transaction(func(tx *gorm.DB) error {
		// 在事務中的所有操作都會使用 Master
		if err := tx.Create(&models.User{Name: "李四"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Destination{Name: "日月潭"}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("事務執行失敗: %v", err)
	}

	// 6.2 複雜查詢（自動使用 Slave）
	var count int64
	db.Model(&models.User{}).
		Where("created_at > ?", "2024-01-01").
		Count(&count)
	log.Printf("新使用者數量: %d", count)

	// 6.3 讀寫分離的最佳實踐
	// 寫入：會自動使用 Master
	newUser := &models.User{Name: "王五"}
	db.Create(newUser)

	// 讀取：會自動使用 Slave
	var foundUser models.User
	db.First(&foundUser, newUser.ID)

	// 如果需要立即讀取剛寫入的資料（避免主從延遲），使用 Master
	dbManager.GetMaster().First(&foundUser, newUser.ID)

	// 6.4 強制使用特定連線
	// 使用 dbresolver.Write 強制使用 Master
	import "gorm.io/plugin/dbresolver"
	db.Clauses(dbresolver.Write).First(&foundUser, newUser.ID)

	// 使用 dbresolver.Read 強制使用 Slave
	db.Clauses(dbresolver.Read).Find(&users)

// ============================================================
// 在 Service 層中的使用範例
// ============================================================

type UserService struct {
	dbManager *database.DBManager
}

func NewUserService() *UserService {
	return &UserService{
		dbManager: database.GetInstance(),
	}
}

// CreateUser 建立使用者（寫入操作）
func (s *UserService) CreateUser(user *models.User) error {
	// 方法 1：使用 GetDB() 讓 GORM 自動處理
	return s.dbManager.GetDB().Create(user).Error

	// 方法 2：明確使用 GetMaster()
	// return s.dbManager.GetMaster().Create(user).Error
}

// GetUserByID 根據 ID 取得使用者（讀取操作）
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	// 方法 1：使用 GetDB() 讓 GORM 自動使用 Slave
	err := s.dbManager.GetDB().First(&user, id).Error

	// 方法 2：明確使用 GetSlave()
	// err := s.dbManager.GetSlave().First(&user, id).Error

	return &user, err
}

// ListUsers 列出所有使用者（讀取操作）
func (s *UserService) ListUsers() ([]models.User, error) {
	var users []models.User
	err := s.dbManager.GetDB().Find(&users).Error
	return users, err
}

// UpdateUser 更新使用者（寫入操作）
func (s *UserService) UpdateUser(user *models.User) error {
	// 更新操作會自動使用 Master
	return s.dbManager.GetDB().Save(user).Error
}

// DeleteUser 刪除使用者（寫入操作）
func (s *UserService) DeleteUser(id uint) error {
	// 刪除操作會自動使用 Master
	return s.dbManager.GetDB().Delete(&models.User{}, id).Error
}

// ============================================================
// 設定檔範例
// ============================================================

/*
config.yaml 範例：

database:
  maxidleconns: 10
  maxopenconns: 100
  connmaxlifetime: 3600

  masters:
    # 主要資料庫
    - name: main
      host: localhost
      port: "3306"
      user: root
      password: "your_password"
      dbname: tourhelper
      charset: utf8mb4
      parsetime: true
      loc: Local

    # 分析資料庫（可選）
    - name: analytics
      host: analytics-db.example.com
      port: "3306"
      user: analytics_user
      password: "analytics_password"
      dbname: analytics
      charset: utf8mb4
      parsetime: true
      loc: Local
      schema: analytics  # 指定此 Master 處理的 Schema

  slaves:
    # Slave 1（權重較高）
    - name: slave1
      host: slave1.example.com
      port: "3306"
      user: readonly_user
      password: "readonly_password"
      dbname: tourhelper
      charset: utf8mb4
      parsetime: true
      loc: Local
      weight: 100
      mastername: main

    # Slave 2（權重較低）
    - name: slave2
      host: slave2.example.com
      port: "3306"
      user: readonly_user
      password: "readonly_password"
      dbname: tourhelper
      charset: utf8mb4
      parsetime: true
      loc: Local
      weight: 50
      mastername: main
*/

// ============================================================
// 注意事項
// ============================================================

/*
1. GORM DBResolver 自動讀寫分離：
   - 使用 GetDB() 取得連線，GORM 會自動判斷操作類型
   - 寫入操作（INSERT, UPDATE, DELETE）自動使用 Master
   - 讀取操作（SELECT）自動使用 Slave（如無 Slave 則降級為 Master）
   - 事務操作自動使用 Master

2. 明確指定連線：
   - 使用 GetMaster() 強制使用 Master（適合需要立即讀取剛寫入資料的場景）
   - 使用 GetSlave() 強制使用 Slave（適合明確的唯讀操作）
   - 使用 Clauses(dbresolver.Write) 或 Clauses(dbresolver.Read) 精確控制

3. 主從同步延遲：
   - Slave 資料可能比 Master 延遲幾毫秒到幾秒
   - 對於即時性要求高的讀取（例如剛寫入的資料），使用 GetMaster()
   - 一般查詢使用 GetDB() 即可，讓 GORM 自動選擇 Slave

4. 負載平衡：
   - DBResolver 使用 RandomPolicy 隨機選擇 Slave
   - 未來可擴展為自訂策略（如加權隨機、輪詢等）

5. 錯誤處理：
   - DBResolver 會自動處理 Slave 不可用的情況，降級為 Master
   - 建議在生產環境中加入健康檢查和監控

6. 連線池設定：
   - 可為每個 Master/Slave 單獨設定連線池參數
   - 未設定時使用全域設定

7. Schema 分離：
   - 使用 GetBySchema() 根據業務邏輯分離資料庫
   - 適合多租戶或微服務架構

8. DBResolver 進階功能：
   - 支援多個 Source（多個 Master）
   - 支援多個 Replica（多個 Slave）
   - 支援自訂負載平衡策略
   - 支援為特定 Model 或 Table 指定不同的 Resolver
*/
