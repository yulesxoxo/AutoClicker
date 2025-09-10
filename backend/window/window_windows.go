package window

import (
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

var (
	user32                  = windows.NewLazyDLL("user32.dll")
	procFindWindow          = user32.NewProc("FindWindowW")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowText       = user32.NewProc("GetWindowTextW")
	procGetClassName        = user32.NewProc("GetClassNameW")
)

const (
	WindowTextBufferSize = 256
	ClassNameBufferSize  = 256
)

func GetForegroundWindow() (hwnd uintptr) {
	hwnd, _, _ = procGetForegroundWindow.Call()
	return hwnd
}

func GetWindowTextW(hwnd uintptr) (title string, err error) {
	titleBuf := make([]uint16, WindowTextBufferSize)
	_, _, err = procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&titleBuf[0])), uintptr(len(titleBuf)))
	if err != nil {
		return title, err
	}
	return windows.UTF16ToString(titleBuf), nil
}

func GetClassNameW(hwnd uintptr) (string, error) {
	classBuf := make([]uint16, ClassNameBufferSize)
	_, _, err := procGetClassName.Call(hwnd, uintptr(unsafe.Pointer(&classBuf[0])), uintptr(len(classBuf)))
	if err != nil {
		return "", err
	}
	className := windows.UTF16ToString(classBuf)
	return className, err
}

func FindWindowW(className, windowTitle string) (uintptr, error) {
	var classNamePtr, titlePtr *uint16
	var err error

	if className != "" {
		classNamePtr, err = windows.UTF16PtrFromString(className)
		if err != nil {
			return 0, err
		}
	}

	if windowTitle != "" {
		titlePtr, err = windows.UTF16PtrFromString(windowTitle)
		if err != nil {
			return 0, err
		}
	}

	hwnd, _, _ := procFindWindow.Call(
		uintptr(unsafe.Pointer(classNamePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
	)

	if hwnd == 0 {
		return 0, fmt.Errorf("window not found with class: %s, title: %s", className, windowTitle)
	}

	return hwnd, nil
}
