package face

import (
    _ "embed"
    "fmt"
    "github.com/fogleman/gg"
    "github.com/golang/freetype/truetype"
    "github.com/peter-maguire/pong-clock/asset"
    "github.com/peter-maguire/pong-clock/clock"
    "github.com/peter-maguire/pong-clock/entity"
    rgbmatrix "github.com/peter-maguire/pong-clock/lib/go-rpi-rgb-led-matrix"
    "golang.org/x/image/font"
    "image"
    "image/color"
    "image/draw"
    "math"
)

var menuFont font.Face

var menuItems = map[string][]string{
    "":        {"Exit", "Faces", "Display", "Time", "WiFi"},
    "Exit":    {},
    "Faces":   {"Default", "Cycle: On", "Back"},
    "Display": {"Brightness", "Auto-Off: Off", "Back"},
    "Time":    {"Timezone", "Manual Set", "Back"},
    "WiFi":    {"Network", "HTTP Server", "Back"},
}

type Menu struct {
    clk         *clock.Clock
    c           *rgbmatrix.Canvas
    ctx         *gg.Context
    currentMenu string
    currentItem int
    animTime    float64
}

func NewMenu(c *clock.Clock) *Menu {
    return &Menu{
        clk:         c,
        c:           c.Canvas,
        currentItem: 1,
    }
}

func (m *Menu) Init() {
    if menuFont == nil {
        menuFontObj, _ := truetype.Parse(asset.MenuFontData)

        menuFont = truetype.NewFace(menuFontObj, &truetype.Options{
            Size: 12,
        })
    }

    m.ctx = gg.NewContext(128, 64)
    m.ctx.SetFontFace(menuFont)
}

func (m *Menu) Stop() {

}

func (m *Menu) Render() {

    m.ctx.SetColor(color.Transparent)
    m.ctx.Clear()
    m.ctx.SetColor(color.White)

    currentMenu := menuItems[m.currentMenu]

    for i := -2; i < 5; i++ {
        x, y := m.DrawMenuItem(230 + m.animTime + float64(i*35))
        if i == 1 && m.animTime < 5 {
            m.ctx.SetColor(color.RGBA{R: 255, A: 255})
        }
        index := (len(currentMenu) + m.currentItem - 1 + i) % len(currentMenu)
        if index < 0 {
            continue
        }
        m.ctx.DrawStringAnchored(currentMenu[index], x, y, 0, 1)
        m.ctx.SetColor(color.White)
    }

    x, y := m.DrawMenuItem(230 + (35 * 3))
    m.ctx.DrawStringAnchored("> ", x+10, y-32, 0, 1)

    m.ctx.Stroke()
    draw.Draw(m.c, m.c.Bounds(), m.ctx.Image(), image.Point{X: 0, Y: 0}, draw.Over)
}

func (m *Menu) DrawMenuItem(rad float64) (float64, float64) {
    r := float64(32)

    cx := float64(-16)
    cy := float64(24)

    x := cx + r*math.Cos((rad+90)*math.Pi/180)
    y := cy + r*math.Sin((rad+90)*math.Pi/180)
    return x, y
}

func (m *Menu) Logic() {
    if m.animTime != 0 {
        m.animTime = math.Round(m.animTime / 3)
    }
}

func (m *Menu) HandleKey(key entity.Button) bool {
    switch key {
    case entity.ButtonUp:
        if m.currentItem > 0 {
            m.currentItem--
        } else {
            m.currentItem = len(menuItems[m.currentMenu]) - 1
        }
        m.animTime = -35

    case entity.ButtonDown:
        if m.currentItem < len(menuItems[m.currentMenu])-1 {
            m.currentItem++
        } else {
            m.currentItem = 0
        }
        m.animTime = 35
    case entity.ButtonSelect:
        fmt.Println(m.currentMenu, m.currentItem)
        selectedItem := m.currentItem
        m.currentItem = 0
        m.currentMenu = menuItems[m.currentMenu][selectedItem]
        fmt.Println("Switched to ", m.currentMenu, m.currentItem)
        if m.currentMenu == "Exit" {
            m.currentMenu = ""
            fmt.Println("Exiting")
            m.clk.ChangeFace(m.clk.Faces[1])
        } else if m.currentMenu == "Back" {
            m.currentMenu = ""
            fmt.Println("Going back")
        } else {
            face := NewMenuSettingSlider(m.clk, m.clk.Settings.Values[0].(*clock.IntRangeSetting))
            m.clk.ChangeFace(face)

        }
    }
    return true
}
