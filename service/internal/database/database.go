package database

import (
	"fmt"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/logger"
)

// Init 初始化所有資料庫連線（MySQL 和 Redis）
func Init(cfg *config.Config) error {
	// 1. 初始化 MySQL
	if err := InitMySQL(cfg.Database); err != nil {
		return fmt.Errorf("初始化 MySQL 失敗: %w", err)
	}
	logger.Info("MySQL 初始化成功")

	// 2. 初始化 Redis（如果有設定）
	if len(cfg.Redis.Instances) > 0 {
		if err := InitRedis(cfg.Redis); err != nil {
			return fmt.Errorf("初始化 Redis 失敗: %w", err)
		}
		logger.Info("Redis 初始化成功")
	} else {
		logger.Warn("未設定 Redis，跳過 Redis 初始化")
	}

	return nil
}

// Close 關閉所有資料庫連線
func Close() error {
	var errs []error

	// 關閉 MySQL
	if mysqlInstance != nil {
		if err := mysqlInstance.Close(); err != nil {
			errs = append(errs, fmt.Errorf("關閉 MySQL 失敗: %w", err))
		} else {
			logger.Info("MySQL 連線已關閉")
		}
	}

	// 關閉 Redis
	if redisInstance != nil {
		if err := redisInstance.Close(); err != nil {
			errs = append(errs, fmt.Errorf("關閉 Redis 失敗: %w", err))
		} else {
			logger.Info("Redis 連線已關閉")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("關閉資料庫時發生錯誤: %v", errs)
	}

	return nil
}

// GetMySQL 快速取得 MySQL 管理器（convenience function）
func GetMySQL() *MySQLManager {
	return GetMySQLInstance()
}

// GetRedis 快速取得 Redis 管理器（convenience function）
func GetRedis() *RedisManager {
	return GetRedisInstance()
}
