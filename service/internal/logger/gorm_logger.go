package logger

import (
	"context"
	"fmt"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

// GormLogger 將 GORM 的 log 輸出到 Logrus
type GormLogger struct {
	LogLevel           gormLogger.LogLevel
	SlowThreshold      time.Duration // 慢查詢門檻
	LogSlowQuery       bool          // 是否記錄慢查詢
	LogAllQueries      bool          // 是否記錄所有查詢
}

// NewGormLogger 建立 GORM Logger 實例
func NewGormLogger(level gormLogger.LogLevel, slowThreshold time.Duration, logSlowQuery, logAllQueries bool) *GormLogger {
	return &GormLogger{
		LogLevel:      level,
		SlowThreshold: slowThreshold,
		LogSlowQuery:  logSlowQuery,
		LogAllQueries: logAllQueries,
	}
}

// LogMode 設定日誌級別
func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 輸出 Info 級別日誌
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		Infof(msg, data...)
	}
}

// Warn 輸出 Warn 級別日誌
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		Warnf(msg, data...)
	}
}

// Error 輸出 Error 級別日誌
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		Errorf(msg, data...)
	}
}

// Trace 追蹤 SQL 執行
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := map[string]interface{}{
		"elapsed": fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
		"rows":    rows,
		"sql":     sql,
	}

	switch {
	// 1. 有錯誤時一定記錄
	case err != nil && l.LogLevel >= gormLogger.Error:
		fields["error"] = err
		WithFields(fields).Error("SQL 執行錯誤")

	// 2. 慢查詢記錄（根據設定）
	case l.LogSlowQuery && elapsed >= l.SlowThreshold && l.LogLevel >= gormLogger.Warn:
		WithFields(fields).Warn("SQL 慢查詢")

	// 3. 記錄所有查詢（開發模式）
	case l.LogAllQueries && l.LogLevel >= gormLogger.Info:
		WithFields(fields).Info("SQL 執行")
	}
}
