package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/redis/go-redis/v9"
)

// RedisManager Redis 管理器（支援多個 Redis 實例）
type RedisManager struct {
	clients map[string]*redis.Client // 多個 Redis 實例，key 為實例名稱
	mu      sync.RWMutex             // 讀寫鎖

	config config.RedisConfig
}

var (
	redisInstance *RedisManager
	redisOnce     sync.Once
)

// InitRedis 初始化 Redis 連線（單例模式）
func InitRedis(cfg config.RedisConfig) error {
	var err error
	redisOnce.Do(func() {
		redisInstance = &RedisManager{
			clients: make(map[string]*redis.Client),
			config:  cfg,
		}
		err = redisInstance.initialize()
	})
	return err
}

// GetRedisInstance 取得 RedisManager 實例
func GetRedisInstance() *RedisManager {
	if redisInstance == nil {
		panic("redis not initialized, call InitRedis() first")
	}
	return redisInstance
}

// initialize 初始化所有 Redis 連線
func (m *RedisManager) initialize() error {
	if len(m.config.Instances) == 0 {
		logger.Warn("未設定任何 Redis 實例，跳過 Redis 初始化")
		return nil
	}

	// 為每個 Redis 實例建立客戶端
	for _, instanceCfg := range m.config.Instances {
		client, err := m.createRedisClient(instanceCfg)
		if err != nil {
			return fmt.Errorf("建立 Redis 實例 [%s] 失敗: %w", instanceCfg.Name, err)
		}

		// 測試連線
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			return fmt.Errorf("Redis 實例 [%s] 連線測試失敗: %w", instanceCfg.Name, err)
		}

		logger.Infof("Redis 實例 [%s] 連線成功: %s:%d (DB: %d)",
			instanceCfg.Name, instanceCfg.Host, instanceCfg.Port, instanceCfg.DB)

		// 使用實例名稱作為 key
		m.clients[instanceCfg.Name] = client
	}

	return nil
}

// createRedisClient 建立 Redis 客戶端
func (m *RedisManager) createRedisClient(cfg config.RedisInstanceConfig) (*redis.Client, error) {
	// 建立 Redis 選項
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	// 設定連線池參數（使用個別設定或全域設定）
	if cfg.PoolSize != nil {
		opts.PoolSize = *cfg.PoolSize
	} else {
		opts.PoolSize = m.config.PoolSize
	}

	if cfg.MinIdleConns != nil {
		opts.MinIdleConns = *cfg.MinIdleConns
	} else {
		opts.MinIdleConns = m.config.MinIdleConns
	}

	if cfg.MaxIdleConns != nil {
		opts.MaxIdleConns = *cfg.MaxIdleConns
	} else {
		opts.MaxIdleConns = m.config.MaxIdleConns
	}

	if cfg.ConnMaxIdleTime != nil {
		opts.ConnMaxIdleTime = *cfg.ConnMaxIdleTime
	} else {
		opts.ConnMaxIdleTime = m.config.ConnMaxIdleTime
	}

	if cfg.ConnMaxLifetime != nil {
		opts.ConnMaxLifetime = *cfg.ConnMaxLifetime
	} else {
		opts.ConnMaxLifetime = m.config.ConnMaxLifetime
	}

	// 設定超時參數
	if cfg.DialTimeout != nil {
		opts.DialTimeout = *cfg.DialTimeout
	} else {
		opts.DialTimeout = m.config.DialTimeout
	}

	if cfg.ReadTimeout != nil {
		opts.ReadTimeout = *cfg.ReadTimeout
	} else {
		opts.ReadTimeout = m.config.ReadTimeout
	}

	if cfg.WriteTimeout != nil {
		opts.WriteTimeout = *cfg.WriteTimeout
	} else {
		opts.WriteTimeout = m.config.WriteTimeout
	}

	// 建立客戶端
	client := redis.NewClient(opts)

	return client, nil
}

// GetClient 取得 Redis 客戶端
// 參數可以是實例名稱，如果不指定則返回預設實例
func (m *RedisManager) GetClient(name ...string) *redis.Client {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 如果指定名稱，返回指定的客戶端
	if len(name) > 0 && name[0] != "" {
		if client, ok := m.clients[name[0]]; ok {
			return client
		}
		panic(fmt.Sprintf("找不到 Redis 實例 [%s]", name[0]))
	}

	// 未指定名稱，返回預設實例（通常是 "main"）
	if len(m.clients) > 0 {
		// 如果有名為 "main" 的實例，優先返回
		if client, ok := m.clients["main"]; ok {
			return client
		}
		// 否則返回第一個
		for _, client := range m.clients {
			return client
		}
	}

	panic("沒有可用的 Redis 實例")
}

// GetClientByDB 根據資料庫名稱取得對應的 Redis 客戶端
// 這個方法會尋找設定中 Database 欄位符合的實例
func (m *RedisManager) GetClientByDB(database string) *redis.Client {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 尋找符合 Database 的實例
	for _, instanceCfg := range m.config.Instances {
		if instanceCfg.Database == database {
			if client, ok := m.clients[instanceCfg.Name]; ok {
				return client
			}
		}
	}

	// 找不到符合的 Database，返回預設實例
	logger.Warnf("找不到資料庫 [%s] 對應的 Redis 實例，使用預設實例", database)
	return m.GetClient()
}

// Ping 測試所有 Redis 連線
func (m *RedisManager) Ping() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var errs []error
	for name, client := range m.clients {
		if err := client.Ping(ctx).Err(); err != nil {
			errs = append(errs, fmt.Errorf("Redis 實例 [%s] Ping 失敗: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Redis 連線測試失敗: %v", errs)
	}

	return nil
}

// Close 關閉所有 Redis 連線
func (m *RedisManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error

	// 關閉所有 Redis 客戶端
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			errs = append(errs, fmt.Errorf("關閉 Redis 實例 [%s] 失敗: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("關閉 Redis 時發生錯誤: %v", errs)
	}

	return nil
}
