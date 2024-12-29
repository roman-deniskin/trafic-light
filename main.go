package main

import (
	"log"
	"runtime"
)

// Точка входа
func main() {
	runtime.LockOSThread() // Все вызовы WinAPI в одном потоке

	// Создаём приложение «Светофор»
	// Передаём все необходимые настройки:
	//   Size: 200x200
	//   FrameWidth: 6
	//   Colors: Red=0x000000FF, Green=0x0000FF00
	app := NewTrafficLightApp(AppConfig{
		Width:       200,
		Height:      200,
		FrameWidth:  6,
		RedColor:    0x000000FF,
		GreenColor:  0x0000FF00,
		X:           300,
		Y:           300,
		WindowTitle: "Светофор",
	})

	// Запускаем приложение
	if err := app.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
