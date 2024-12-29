# Traffic Light (Светофор)

## 1. Назначение программы
Эта программа создаёт круглое окно «светофор», которое может переключаться между красным и зелёным цветом при нажатии **Ctrl + Shift + Space**. Окно можно перемещать (Drag & Drop), «схватив» его левой кнопкой мыши. Программа предназначена для работы только в среде Windows.

---

## 2. Как правильно редактировать параметры программы
В коде программы (в файле `main.go` или любом другом месте, где вы инициализируете приложение) создаётся экземпляр приложения:

```go
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
