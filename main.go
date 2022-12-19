package main

import (
    "github.com/peter-maguire/pong-clock/clock"
    "github.com/peter-maguire/pong-clock/entity"
    "github.com/peter-maguire/pong-clock/face"
    rgbmatrix "github.com/peter-maguire/pong-clock/lib/go-rpi-rgb-led-matrix"
)

func main() {
    m, _ := rgbmatrix.NewRGBLedMatrix(&rgbmatrix.HardwareConfig{
        Rows:                   64,
        Cols:                   64,
        ChainLength:            2,
        Parallel:               1,
        PWMBits:                10,
        PWMLSBNanoseconds:      115,
        Brightness:             60,
        ScanMode:               rgbmatrix.Progressive,
        DisableHardwarePulsing: false,
        ShowRefreshRate:        false,
        HardwareMapping:        "adafruit-hat-pwm",
    })

    c := rgbmatrix.NewCanvas(m)
    defer c.Close()

    clk := clock.NewClock(&m, c)

    clk.Faces = []entity.ClockFace{
        face.NewMenu(clk),
        face.NewPongGame(clk),
        face.NewPacman(clk),
    }
    clk.ChangeFace(clk.Faces[1])

    clk.Start()
}
