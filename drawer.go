package main

// TrafficLightDrawer отвечает за рисование окна (круг + рамка).
type TrafficLightDrawer struct {
	// Ширина рамки (px)
	FrameWidth int32

	// Функции, которые говорят, какой цвет сейчас
	GetIsGreen    func() bool
	GetRedColor   func() uint32
	GetGreenColor func() uint32
}

// Draw — рисуем светофор (круг) в клиентской области
func (d TrafficLightDrawer) Draw(hdc HANDLE, width, height int32) {
	// Определим заливочный цвет (BGR)
	var fillColor uint32
	if d.GetIsGreen() {
		fillColor = d.GetGreenColor()
	} else {
		fillColor = d.GetRedColor()
	}

	// Кисть
	brush := createSolidBrush(fillColor)
	// Перо (рамка) — ширина d.FrameWidth, чёрный
	pen := createPen(0, d.FrameWidth, 0x00000000) // PS_SOLID=0, чёрный=0x00000000

	oldPen := selectObject(hdc, HANDLE(pen))
	oldBrush := selectObject(hdc, HANDLE(brush))

	ellipse(hdc, 0, 0, width, height)

	// Возвращаем старые объекты
	selectObject(hdc, oldPen)
	selectObject(hdc, oldBrush)
	// Удаляем GDI-объекты
	deleteObject(HANDLE(pen))
	deleteObject(HANDLE(brush))
}
