package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
)

// AppConfig хранит настройки приложения
type AppConfig struct {
	Width      int32
	Height     int32
	FrameWidth int32

	RedColor   uint32
	GreenColor uint32

	X int32
	Y int32

	WindowTitle string
}

// TrafficLightApp — «Светофорное» приложение
type TrafficLightApp struct {
	config  AppConfig
	hWnd    HWND           // главное окно
	hHook   windows.Handle // хук клавиатуры
	isGreen bool           // состояние: false=красный, true=зелёный

	drawer TrafficLightDrawer // отвечает за отрисовку
}

// NewTrafficLightApp конструирует приложение
func NewTrafficLightApp(cfg AppConfig) *TrafficLightApp {
	app := &TrafficLightApp{
		config:  cfg,
		isGreen: false,
	}
	// Инициализируем «художника» (drawer), который будет рисовать окно
	app.drawer = TrafficLightDrawer{
		FrameWidth:    cfg.FrameWidth,
		GetIsGreen:    func() bool { return app.isGreen },
		GetRedColor:   func() uint32 { return app.config.RedColor },
		GetGreenColor: func() uint32 { return app.config.GreenColor },
	}
	return app
}

// Run — основной метод запуска
func (app *TrafficLightApp) Run() error {
	// 1) Регистрируем класс окна (через WinAPI)
	className, _ := syscall.UTF16PtrFromString("TrafficLightClassSOLID")
	title, _ := syscall.UTF16PtrFromString(app.config.WindowTitle)

	// Передаём ссылку на app в глобальную переменную, чтобы wndProc мог к нему обратиться
	gApp = app

	// Регистрируем класс
	if err := registerWindowClass(className, syscall.NewCallback(wndProcCallback)); err != nil {
		return fmt.Errorf("registerWindowClass: %v", err)
	}

	// 2) Создаём окно
	hwnd, err := createMainWindow(className, title, app.config.X, app.config.Y, app.config.Width, app.config.Height)
	if err != nil {
		return err
	}
	app.hWnd = hwnd

	// Делаем окно круглым
	makeWindowRound(hwnd, app.config.Width, app.config.Height)

	// Показываем окно
	showWindow(hwnd, SW_SHOW)
	updateWindow(hwnd)

	// 3) Ставим хук клавиатуры
	if err := app.setKeyboardHook(); err != nil {
		return fmt.Errorf("setKeyboardHook: %v", err)
	}
	defer app.unhookKeyboard()

	// 4) Запускаем цикл сообщений
	return messageLoop()
}

// OnPaint вызывается из wndProc, чтобы отрисовать окно
func (app *TrafficLightApp) OnPaint(hdc HANDLE, width, height int32) {
	// Передаём управление «художнику»
	app.drawer.Draw(hdc, width, height)
}

// OnKeyPressed вызывается из keyboardProc, когда нажаты клавиши
func (app *TrafficLightApp) OnKeyPressed(kb *KBDLLHOOKSTRUCT) {
	// Проверяем, надо ли переключать цвет
	if isToggleHotkey(kb) {
		// Меняем состояние
		app.isGreen = !app.isGreen
		// Перерисуем окно
		invalidateRect(app.hWnd, nil, true)
	}
}

// Снятие/установка хука
func (app *TrafficLightApp) setKeyboardHook() error {
	cb := syscall.NewCallback(keyboardProcCallback)
	r1, _, err := procSetWindowsHookExW.Call(
		WH_KEYBOARD_LL,
		cb,
		0,
		0,
	)
	if r1 == 0 {
		return err
	}
	app.hHook = windows.Handle(r1)
	return nil
}

func (app *TrafficLightApp) unhookKeyboard() {
	if app.hHook != 0 {
		procUnhookWindowsHookEx.Call(uintptr(app.hHook))
		app.hHook = 0
	}
}
