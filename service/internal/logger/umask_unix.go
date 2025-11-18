// go:build !windows
//go:build !windows
// +build !windows

package logger

import "syscall"

// setUmask 設定檔案建立的 umask（僅限 Unix/Linux 系統）
func setUmask(mask int) int {
	return syscall.Umask(mask)
}
