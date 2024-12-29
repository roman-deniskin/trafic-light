package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// ----------------------------------------------------------------------------
// Глобальные переменные WinAPI
// ----------------------------------------------------------------------------
var (
	user32 = windows.NewLazySystemDLL("user32.dll")
	gdi32  = windows.NewLazySystemDLL("gdi32.dll")

	// Функции
	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procUpdateWindow     = user32.NewProc("UpdateWindow")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")

	procBeginPaint    = user32.NewProc("BeginPaint")
	procEndPaint      = user32.NewProc("EndPaint")
	procGetClientRect = user32.NewProc("GetClientRect")

	procCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	procCreatePen        = gdi32.NewProc("CreatePen")
	procSelectObject     = gdi32.NewProc("SelectObject")
	procEllipse          = gdi32.NewProc("Ellipse")
	procDeleteObject     = gdi32.NewProc("DeleteObject")

	procReleaseCapture = user32.NewProc("ReleaseCapture")
	procSendMessageW   = user32.NewProc("SendMessageW")

	procSetWindowsHookExW   = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState")

	procCreateEllipticRgn = gdi32.NewProc("CreateEllipticRgn")
	procSetWindowRgn      = user32.NewProc("SetWindowRgn")
	procInvalidateRect    = user32.NewProc("InvalidateRect")
)

// ----------------------------------------------------------------------------
// Типы
// ----------------------------------------------------------------------------
type (
	HANDLE    uintptr
	HINSTANCE HANDLE
	HBRUSH    HANDLE
	HICON     HANDLE
	HCURSOR   HANDLE
	HWND      HANDLE
	HRGN      HANDLE
)

// WNDCLASSEXW
type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       HICON
}

// MSG
type MSG struct {
	Hwnd     HWND
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Point    POINT
	LPrivate uint32
}

// PAINTSTRUCT
type PAINTSTRUCT struct {
	Hdc         HANDLE
	FErase      bool
	RcPaint     RECT
	FRestore    bool
	FIncUpdate  bool
	RgbReserved [32]byte
}

// RECT
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// POINT
type POINT struct {
	X int32
	Y int32
}

// KBDLLHOOKSTRUCT
type KBDLLHOOKSTRUCT struct {
	VkCode    uint32
	ScanCode  uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

// ----------------------------------------------------------------------------
// Константы WinAPI
// ----------------------------------------------------------------------------
const (
	WS_POPUP         = 0x80000000
	WS_VISIBLE       = 0x10000000
	WS_EX_TOPMOST    = 0x00000008
	WS_EX_TOOLWINDOW = 0x00000080

	SW_SHOW = 5

	WM_PAINT         = 0x000F
	WM_DESTROY       = 0x0002
	WM_LBUTTONDOWN   = 0x0201
	WM_NCLBUTTONDOWN = 0x00A1
	WM_QUIT          = 0x0012

	HTCAPTION = 2

	WH_KEYBOARD_LL = 13
	HC_ACTION      = 0
	WM_KEYDOWN     = 0x0100
	WM_SYSKEYDOWN  = 0x0104

	VK_LCONTROL = 0xA2
	VK_RCONTROL = 0xA3
	VK_LSHIFT   = 0xA0
	VK_RSHIFT   = 0xA1
	VK_SPACE    = 0x20
)

// ----------------------------------------------------------------------------
// Глобальная ссылка на приложение (для wndProc/keyboardProc)
// ----------------------------------------------------------------------------
var gApp *TrafficLightApp

// ----------------------------------------------------------------------------
// RegisterClassEx
// ----------------------------------------------------------------------------
func registerWindowClass(className *uint16, wndProcCallback uintptr) error {
	var wcx WNDCLASSEXW
	wcx.CbSize = uint32(unsafe.Sizeof(wcx))
	wcx.LpfnWndProc = wndProcCallback
	// остальное оставим по умолчанию
	wcx.LpszClassName = className

	r, _, e := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wcx)))
	atom := uint16(r)
	if atom == 0 {
		return e
	}
	return nil
}

// ----------------------------------------------------------------------------
// Create Window
// ----------------------------------------------------------------------------
func createMainWindow(className, windowName *uint16, x, y, width, height int32) (HWND, error) {
	dwExStyle := WS_EX_TOPMOST | WS_EX_TOOLWINDOW
	dwStyle := WS_POPUP | WS_VISIBLE

	r, _, e := procCreateWindowExW.Call(
		uintptr(dwExStyle),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		uintptr(dwStyle),
		uintptr(x), uintptr(y),
		uintptr(width), uintptr(height),
		0, 0,
		0, 0,
	)
	if r == 0 {
		return 0, e
	}
	return HWND(r), nil
}

// ----------------------------------------------------------------------------
// makeWindowRound
// ----------------------------------------------------------------------------
func makeWindowRound(hwnd HWND, width, height int32) {
	hrgn := createEllipticRgn(0, 0, width, height)
	if hrgn != 0 {
		setWindowRgn(hwnd, hrgn, true)
		deleteObject(HANDLE(hrgn))
	}
}

// ----------------------------------------------------------------------------
// Show/Update Window
// ----------------------------------------------------------------------------
func showWindow(hwnd HWND, cmdShow int32) {
	procShowWindow.Call(uintptr(hwnd), uintptr(cmdShow))
}

func updateWindow(hwnd HWND) {
	procUpdateWindow.Call(uintptr(hwnd))
}

// ----------------------------------------------------------------------------
// Message loop
// ----------------------------------------------------------------------------
func messageLoop() error {
	var msg MSG
	for {
		r, _, _ := procGetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0, 0, 0,
		)
		if r == 0 {
			// WM_QUIT
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
	return nil
}

// ----------------------------------------------------------------------------
// wndProcCallback
// ----------------------------------------------------------------------------
func wndProcCallback(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_PAINT:
		var ps PAINTSTRUCT
		hdc := beginPaint(HWND(hwnd), &ps)
		defer endPaint(HWND(hwnd), &ps)

		// Определим размеры
		var rc RECT
		getClientRect(HWND(hwnd), &rc)
		w := rc.Right - rc.Left
		h := rc.Bottom - rc.Top

		if gApp != nil {
			gApp.OnPaint(hdc, w, h)
		}

		return 0

	case WM_LBUTTONDOWN:
		// Перемещение окна
		releaseCapture()
		sendMessageW(HWND(hwnd), WM_NCLBUTTONDOWN, HTCAPTION, 0)
		return 0

	case WM_DESTROY:
		if gApp != nil {
			gApp.unhookKeyboard()
		}
		postQuitMessage(0)
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(
		hwnd,
		uintptr(msg),
		wparam,
		lparam,
	)
	return ret
}

// ----------------------------------------------------------------------------
// keyboardProcCallback
// ----------------------------------------------------------------------------
func keyboardProcCallback(code int, wparam uintptr, lparam uintptr) uintptr {
	if code == HC_ACTION && gApp != nil {
		// Разбираем структуру KBDLLHOOKSTRUCT
		kb := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		// WM_KEYDOWN / WM_SYSKEYDOWN
		if wparam == WM_KEYDOWN || wparam == WM_SYSKEYDOWN {
			gApp.OnKeyPressed(kb)
		}
	}
	ret, _, _ := procCallNextHookEx.Call(
		0,
		uintptr(code),
		wparam,
		lparam,
	)
	return ret
}

// ----------------------------------------------------------------------------
// GDI Helpers
// ----------------------------------------------------------------------------
func beginPaint(hwnd HWND, ps *PAINTSTRUCT) HANDLE {
	r, _, _ := procBeginPaint.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(ps)),
	)
	return HANDLE(r)
}

func endPaint(hwnd HWND, ps *PAINTSTRUCT) {
	procEndPaint.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(ps)),
	)
}

func getClientRect(hwnd HWND, rect *RECT) {
	procGetClientRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(rect)),
	)
}

func releaseCapture() {
	procReleaseCapture.Call()
}

func sendMessageW(hwnd HWND, msg uint32, wParam, lParam uintptr) {
	procSendMessageW.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
}

func createEllipticRgn(left, top, right, bottom int32) HRGN {
	r, _, _ := procCreateEllipticRgn.Call(
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
	)
	return HRGN(r)
}

func setWindowRgn(hwnd HWND, hrgn HRGN, redraw bool) int32 {
	rd := 0
	if redraw {
		rd = 1
	}
	r, _, _ := procSetWindowRgn.Call(
		uintptr(hwnd),
		uintptr(hrgn),
		uintptr(rd),
	)
	return int32(r)
}

func createSolidBrush(color uint32) HBRUSH {
	r, _, _ := procCreateSolidBrush.Call(uintptr(color))
	return HBRUSH(r)
}

func createPen(style, width int32, color uint32) HANDLE {
	r, _, _ := procCreatePen.Call(
		uintptr(style),
		uintptr(width),
		uintptr(color),
	)
	return HANDLE(r)
}

func selectObject(hdc HANDLE, obj HANDLE) HANDLE {
	r, _, _ := procSelectObject.Call(uintptr(hdc), uintptr(obj))
	return HANDLE(r)
}

func ellipse(hdc HANDLE, left, top, right, bottom int32) bool {
	rr, _, _ := procEllipse.Call(
		uintptr(hdc),
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
	)
	return rr != 0
}

func deleteObject(obj HANDLE) {
	procDeleteObject.Call(uintptr(obj))
}

func postQuitMessage(exitCode int32) {
	procPostQuitMessage.Call(uintptr(exitCode))
}

func invalidateRect(hwnd HWND, rect *RECT, erase bool) {
	e := 0
	if erase {
		e = 1
	}
	procInvalidateRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(rect)),
		uintptr(e),
	)
}

// ----------------------------------------------------------------------------
// Горячие клавиши
// ----------------------------------------------------------------------------
func isToggleHotkey(kb *KBDLLHOOKSTRUCT) bool {
	ctrl := getAsyncKeyState(VK_LCONTROL) || getAsyncKeyState(VK_RCONTROL)
	shift := getAsyncKeyState(VK_LSHIFT) || getAsyncKeyState(VK_RSHIFT)
	return (kb.VkCode == VK_SPACE) && ctrl && shift
}

func getAsyncKeyState(vkCode uint32) bool {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(vkCode))
	return (r>>15)&1 == 1
}
