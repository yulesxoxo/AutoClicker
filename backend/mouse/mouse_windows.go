package mouse

import (
	"autoclicker/backend/system"
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

var (
	user32           = windows.NewLazyDLL("user32.dll")
	procSendInput    = user32.NewProc("SendInput")
	procPostMessage  = user32.NewProc("PostMessageW")
	procSendMessage  = user32.NewProc("SendMessageW")
	procSetCursorPos = user32.NewProc("SetCursorPos")
)

// Windows API constants
const (
	INPUT_MOUSE            = 0
	MOUSEEVENTF_LEFTDOWN   = 0x0002
	MOUSEEVENTF_LEFTUP     = 0x0004
	MOUSEEVENTF_RIGHTDOWN  = 0x0008
	MOUSEEVENTF_RIGHTUP    = 0x0010
	MOUSEEVENTF_MIDDLEDOWN = 0x0020
	MOUSEEVENTF_MIDDLEUP   = 0x0040
	MOUSEEVENTF_ABSOLUTE   = 0x8000
)

// SendMessage/PostMessage constants
const (
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
)

const (
	WindowsAbsoluteCoordMax = 65535 // Windows uses 0-65535 for absolute coordinates
)

type POINT struct {
	X, Y int32
}

type MOUSEINPUT struct {
	DX          int32
	DY          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	Mi   MOUSEINPUT
}

// Absolute coordinate range 0-65535
func convertScreenCoordinatesToAbsoluteCoordinates(x, y int) (int32, int32, error) {
	width, height, err := system.GetScreenDimensions()
	if err != nil {
		return 0, 0, err
	}
	absX := int32((x * WindowsAbsoluteCoordMax) / width)
	absY := int32((y * WindowsAbsoluteCoordMax) / height)

	return absX, absY, nil
}

// SendInput method
func sendInput(x, y int, button MouseButton, moveCursor bool) error {
	downFlag, upFlag := getInputFlags(button)
	var inputs [2]INPUT

	if moveCursor {
		_, _, err := procSetCursorPos.Call(uintptr(x), uintptr(y))
		if err != nil {
			return err
		}
		inputs[0] = INPUT{Type: INPUT_MOUSE, Mi: MOUSEINPUT{DwFlags: downFlag}}
		inputs[1] = INPUT{Type: INPUT_MOUSE, Mi: MOUSEINPUT{DwFlags: upFlag}}
	} else {
		absX, absY, err := convertScreenCoordinatesToAbsoluteCoordinates(x, y)
		if err != nil {
			return err
		}
		inputs[0] = INPUT{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				DX:      absX,
				DY:      absY,
				DwFlags: downFlag | MOUSEEVENTF_ABSOLUTE,
			},
		}
		inputs[1] = INPUT{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				DX:      absX,
				DY:      absY,
				DwFlags: upFlag | MOUSEEVENTF_ABSOLUTE,
			},
		}
	}

	ret, _, err := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		return err
	}
	return nil
}

// PostMessage method (async, sends to specific window, doesn't move cursor)
func postMessage(x, y int, button MouseButton, hwnd uintptr) error {
	if hwnd == 0 {
		return fmt.Errorf("no target window set for PostMessage (use SetTargetWindow or FindAndSetWindow)")
	}

	downMsg, upMsg := getMessageCodes(button)
	lParam := uintptr(y<<16 | x)

	ret1, _, err := procPostMessage.Call(hwnd, uintptr(downMsg), 1, lParam)
	if err != nil {
		return err
	}
	ret2, _, err := procPostMessage.Call(hwnd, uintptr(upMsg), 0, lParam)
	if err != nil {
		return err
	}
	if ret1 == 0 || ret2 == 0 {
		return fmt.Errorf("PostMessage failed")
	}

	return nil
}

// SendMessage method (sync, sends to specific window, doesn't move cursor)
func sendMessage(x, y int, button MouseButton, hwnd uintptr) error {
	if hwnd == 0 {
		return fmt.Errorf("no target window set for SendMessage (use SetTargetWindow or FindAndSetWindow)")
	}

	downMsg, upMsg := getMessageCodes(button)
	lParam := uintptr(y<<16 | x)

	_, _, err := procSendMessage.Call(hwnd, uintptr(downMsg), 1, lParam)
	if err != nil {
		return err
	}
	_, _, err = procSendMessage.Call(hwnd, uintptr(upMsg), 0, lParam)
	if err != nil {
		return err
	}

	return nil
}

func getInputFlags(button MouseButton) (down, up uint32) {
	switch button {
	case LeftButton:
		return MOUSEEVENTF_LEFTDOWN, MOUSEEVENTF_LEFTUP
	case RightButton:
		return MOUSEEVENTF_RIGHTDOWN, MOUSEEVENTF_RIGHTUP
	case MiddleButton:
		return MOUSEEVENTF_MIDDLEDOWN, MOUSEEVENTF_MIDDLEUP
	default:
		return MOUSEEVENTF_LEFTDOWN, MOUSEEVENTF_LEFTUP
	}
}

func getMessageCodes(button MouseButton) (down, up uint32) {
	switch button {
	case LeftButton:
		return WM_LBUTTONDOWN, WM_LBUTTONUP
	case RightButton:
		return WM_RBUTTONDOWN, WM_RBUTTONUP
	case MiddleButton:
		return WM_MBUTTONDOWN, WM_MBUTTONUP
	default:
		return WM_LBUTTONDOWN, WM_LBUTTONUP
	}
}
