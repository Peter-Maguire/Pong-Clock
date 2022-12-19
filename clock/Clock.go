package clock

import (
    "fmt"
    "github.com/eiannone/keyboard"
    "github.com/peter-maguire/pong-clock/entity"
    rgbmatrix "github.com/peter-maguire/pong-clock/lib/go-rpi-rgb-led-matrix"
    "time"
)

const renderInterval = 15 * time.Millisecond
const logicInterval = 20 * time.Millisecond

type Clock struct {
    matrix      *rgbmatrix.Matrix
    Canvas      *rgbmatrix.Canvas
    Faces       []entity.ClockFace
    currentFace entity.ClockFace
    Settings    Settings
}

func NewClock(m *rgbmatrix.Matrix, c *rgbmatrix.Canvas) *Clock {
    return &Clock{
        matrix: m,
        Canvas: c,
        Settings: Settings{Values: []SettingType{
            NewIntRangeSetting("Brightness", 100, 0, 100),
        }},
    }
}

func (c *Clock) Start() {
    c.currentFace.Init()

    go c.renderLoop()
    go c.gameLogic()
    go c.readInput()
    forever := make(chan bool)
    <-forever
}

func (c *Clock) CurrentFace() entity.ClockFace {
    return c.currentFace
}

func (c *Clock) ChangeFace(face entity.ClockFace) {
    face.Init()
    if c.currentFace != nil {
        c.currentFace.Stop()
    }
    c.currentFace = face
}

func (c *Clock) renderLoop() {
    t := time.NewTicker(renderInterval)
    for range t.C {
        c.CurrentFace().Render()
        c.Canvas.Render()
    }
}

func (c *Clock) gameLogic() {
    t := time.NewTicker(logicInterval)
    for range t.C {
        c.CurrentFace().Logic()
    }
}

func (c *Clock) readInput() {
    keysEvents, err := keyboard.GetKeys(10)
    if err != nil {
        fmt.Println("Failed to open keyboard")
        return
    }
    defer keyboard.Close()
    for {
        event := <-keysEvents
        if event.Err != nil {
            panic(event.Err)
        }
        var key entity.Button
        switch event.Key {
        case keyboard.KeyArrowUp:
            key = entity.ButtonUp
        case keyboard.KeyArrowDown:
            key = entity.ButtonDown
        case keyboard.KeyArrowLeft:
            key = entity.ButtonLeft
        case keyboard.KeyArrowRight:
            key = entity.ButtonRight
        case keyboard.KeyEnter:
            key = entity.ButtonSelect
        default:
            continue
        }
        fmt.Println("Key press ", key)
        handled := c.CurrentFace().HandleKey(key)
        if handled {
            continue
        }
        // Change to menu if current face doesn't handle the key
        if key == entity.ButtonSelect {
            c.ChangeFace(c.Faces[0])
        }
    }
}
