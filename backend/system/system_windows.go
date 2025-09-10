package system

import (
	"fmt"
	"golang.org/x/sys/windows"
)

var (
	user32               = windows.NewLazyDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

func GetScreenDimensions() (int, int, error) {
	width, _, _ := procGetSystemMetrics.Call(0)  // SM_CXSCREEN
	height, _, _ := procGetSystemMetrics.Call(1) // SM_CYSCREEN

	if width == 0 || height == 0 {
		return 0, 0, fmt.Errorf("invalid screen dimensions: width=%d, height=%d", width, height)
	}

	return int(width), int(height), nil
}
