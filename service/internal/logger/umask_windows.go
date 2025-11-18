//go:build windows

package logger

// setUmask Windows 系統不支援 umask，這是一個空實作
func setUmask(mask int) int {
	return 0
}
