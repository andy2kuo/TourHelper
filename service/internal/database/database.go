package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/andy2kuo/TourHelper/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// DBManager 資料庫管理器（Master-Slave 架構）
// 使用 GORM DBResolver 插件實現讀寫分離
type DBManager struct {
	databases map[string]*gorm.DB // 多個資料庫實例，key 為 Master 名稱或 Schema
	mu        sync.RWMutex        // 讀寫鎖

	config config.DatabaseConfig
}

var (
	instance *DBManager
	once     sync.Once
)

// Init 初始化資料庫連線（單例模式）
func Init(cfg config.DatabaseConfig) error {
	var err error
	once.Do(func() {
		instance = &DBManager{
			databases: make(map[string]*gorm.DB),
			config:    cfg,
		}
		err = instance.initialize()
	})
	return err
}

// GetInstance 取得 DBManager 實例
func GetInstance() *DBManager {
	if instance == nil {
		panic("database not initialized, call Init() first")
	}
	return instance
}

// initialize 初始化所有資料庫連線
func (m *DBManager) initialize() error {
	if len(m.config.Masters) == 0 {
		return fmt.Errorf("至少需要設定一個 Master 資料庫")
	}

	// 將 Master 按照 MasterName 分組，並找出對應的 Slaves
	masterSlaveMap := m.groupSlavesByMaster()

	// 為每個 Master 建立資料庫實例（含 DBResolver）
	for _, masterCfg := range m.config.Masters {
		db, err := m.createDatabaseWithResolver(masterCfg, masterSlaveMap[masterCfg.Name])
		if err != nil {
			return fmt.Errorf("建立資料庫實例 [%s] 失敗: %w", masterCfg.Name, err)
		}

		// 使用 Master 名稱作為 key
		m.databases[masterCfg.Name] = db

		// 如果有指定 Schema，也用 Schema 作為 key
		if masterCfg.Schema != "" {
			m.databases[masterCfg.Schema] = db
		}
	}

	return nil
}

// groupSlavesByMaster 將 Slave 按照 MasterName 分組
func (m *DBManager) groupSlavesByMaster() map[string][]config.SlaveDBConfig {
	masterSlaveMap := make(map[string][]config.SlaveDBConfig)

	for _, slaveCfg := range m.config.Slaves {
		masterName := slaveCfg.MasterName
		if masterName == "" {
			// 如果未指定 MasterName，預設使用第一個 Master
			if len(m.config.Masters) > 0 {
				masterName = m.config.Masters[0].Name
			}
		}
		masterSlaveMap[masterName] = append(masterSlaveMap[masterName], slaveCfg)
	}

	return masterSlaveMap
}

// createDatabaseWithResolver 建立包含 DBResolver 的資料庫實例
func (m *DBManager) createDatabaseWithResolver(masterCfg config.MasterDBConfig, slaveCfgs []config.SlaveDBConfig) (*gorm.DB, error) {
	// 1. 建立 Master 連線
	masterDSN := m.buildDSN(masterCfg.Host, masterCfg.Port, masterCfg.User, masterCfg.Password,
		masterCfg.DBName, masterCfg.Charset, masterCfg.ParseTime, masterCfg.Loc)

	db, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("開啟 Master 連線失敗: %w", err)
	}

	// 2. 設定 Master 連線池
	if err := m.configureConnectionPool(db, masterCfg.MaxIdleConns, masterCfg.MaxOpenConns, masterCfg.ConnMaxLifetime); err != nil {
		return nil, fmt.Errorf("設定 Master 連線池失敗: %w", err)
	}

	// 3. 如果有 Slave，註冊 DBResolver
	if len(slaveCfgs) > 0 {
		// 建立 Replica (Slave) Dialectors
		replicas := make([]gorm.Dialector, 0, len(slaveCfgs))
		for _, slaveCfg := range slaveCfgs {
			slaveDSN := m.buildDSN(slaveCfg.Host, slaveCfg.Port, slaveCfg.User, slaveCfg.Password,
				slaveCfg.DBName, slaveCfg.Charset, slaveCfg.ParseTime, slaveCfg.Loc)
			replicas = append(replicas, mysql.Open(slaveDSN))
		}

		// 取得 Replica 連線池設定（使用第一個 Slave 的設定或全域設定）
		replicaMaxIdle := m.config.MaxIdleConns
		replicaMaxOpen := m.config.MaxOpenConns
		replicaMaxLifetime := m.config.ConnMaxLifetime
		if len(slaveCfgs) > 0 && slaveCfgs[0].MaxIdleConns != nil {
			replicaMaxIdle = *slaveCfgs[0].MaxIdleConns
		}
		if len(slaveCfgs) > 0 && slaveCfgs[0].MaxOpenConns != nil {
			replicaMaxOpen = *slaveCfgs[0].MaxOpenConns
		}
		if len(slaveCfgs) > 0 && slaveCfgs[0].ConnMaxLifetime != nil {
			replicaMaxLifetime = *slaveCfgs[0].ConnMaxLifetime
		}

		// 建立 DBResolver 設定
		resolverConfig := dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(masterDSN)}, // Master 作為 Source
			Replicas: replicas,                                 // Slaves 作為 Replicas
			Policy:   dbresolver.RandomPolicy{},               // 使用隨機策略
		}

		// 註冊 DBResolver 並設定 Replica 連線池
		if err := db.Use(dbresolver.Register(resolverConfig).
			SetMaxIdleConns(replicaMaxIdle).
			SetMaxOpenConns(replicaMaxOpen).
			SetConnMaxLifetime(time.Duration(replicaMaxLifetime) * time.Second)); err != nil {
			return nil, fmt.Errorf("註冊 DBResolver 失敗: %w", err)
		}
	}

	return db, nil
}

// configureConnectionPool 設定連線池參數
func (m *DBManager) configureConnectionPool(db *gorm.DB, maxIdle, maxOpen, maxLifetime *int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 使用個別設定或全域設定
	idle := m.config.MaxIdleConns
	if maxIdle != nil {
		idle = *maxIdle
	}

	open := m.config.MaxOpenConns
	if maxOpen != nil {
		open = *maxOpen
	}

	lifetime := m.config.ConnMaxLifetime
	if maxLifetime != nil {
		lifetime = *maxLifetime
	}

	sqlDB.SetMaxIdleConns(idle)
	sqlDB.SetMaxOpenConns(open)
	sqlDB.SetConnMaxLifetime(time.Duration(lifetime) * time.Second)

	return nil
}

// buildDSN 建立 MySQL DSN 字串
func (m *DBManager) buildDSN(host, port, user, password, dbName, charset string, parseTime bool, loc string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		user, password, host, port, dbName, charset, parseTime, loc)
}

// GetDB 取得資料庫連線（自動讀寫分離）
// 參數可以是 Master 名稱或 Schema 名稱
func (m *DBManager) GetDB(name ...string) *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 如果指定名稱，返回指定的資料庫
	if len(name) > 0 && name[0] != "" {
		if db, ok := m.databases[name[0]]; ok {
			return db
		}
		panic(fmt.Sprintf("找不到資料庫 [%s]", name[0]))
	}

	// 未指定名稱，返回預設資料庫（通常是 "main"）
	if len(m.databases) > 0 {
		// 如果有名為 "main" 的資料庫，優先返回
		if db, ok := m.databases["main"]; ok {
			return db
		}
		// 否則返回第一個
		for _, db := range m.databases {
			return db
		}
	}

	panic("沒有可用的資料庫")
}

// GetMaster 取得 Master 連線（明確指定使用 Master 進行寫入）
// 使用 Clauses(dbresolver.Write) 確保使用 Master
func (m *DBManager) GetMaster(name ...string) *gorm.DB {
	return m.GetDB(name...).Clauses(dbresolver.Write)
}

// GetSlave 取得 Slave 連線（明確指定使用 Slave 進行讀取）
// 使用 Clauses(dbresolver.Read) 確保使用 Slave（如無 Slave 則降級為 Master）
func (m *DBManager) GetSlave(name ...string) *gorm.DB {
	return m.GetDB(name...).Clauses(dbresolver.Read)
}

// GetBySchema 根據 Schema 取得對應的資料庫連線
func (m *DBManager) GetBySchema(schema string) *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 尋找符合 Schema 的資料庫
	if db, ok := m.databases[schema]; ok {
		return db
	}

	// 找不到符合的 Schema，返回預設資料庫
	return m.GetDB()
}

// Close 關閉所有資料庫連線
func (m *DBManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error

	// 關閉所有資料庫連線（包含 Master 和 Slave）
	for name, db := range m.databases {
		sqlDB, err := db.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("取得資料庫 [%s] SQL DB 失敗: %w", name, err))
			continue
		}
		if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("關閉資料庫 [%s] 失敗: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("關閉資料庫時發生錯誤: %v", errs)
	}

	return nil
}