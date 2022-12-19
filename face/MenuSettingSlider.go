package face

import (
    "github.com/fogleman/gg"
    "github.com/peter-maguire/pong-clock/clock"
    "github.com/peter-maguire/pong-clock/entity"
    rgbmatrix "github.com/peter-maguire/pong-clock/lib/go-rpi-rgb-led-matrix"
    "image"
    "image/color"
    "image/draw"
)

type MenuSettingSlider struct {
    clk        *clock.Clock
    c          *rgbmatrix.Canvas
    ctx        *gg.Context
    setting    *clock.IntRangeSetting
    draftValue int
}

func (m *MenuSettingSlider) Init() {
    m.ctx = gg.NewContext(128, 64)
    m.ctx.SetFontFace(menuFont)
    m.draftValue = m.setting.GetValue()
}

func (m *MenuSettingSlider) Stop() {

}

func (m *MenuSettingSlider) Render() {
    m.ctx.SetColor(color.Transparent)
    m.ctx.Clear()
    m.ctx.SetColor(color.White)

    m.ctx.DrawStringAnchored(m.setting.GetName(), 64, 1, 0.5, 1)

    m.ctx.DrawRectangle(8, 32-8, 128-16, 16)
    m.ctx.Stroke()

    m.ctx.DrawRectangle(9, 32-8+1, (float64(m.draftValue)/float64(m.setting.GetMaxValue()))*(128-16-1), 15)

    m.ctx.Fill()

    draw.Draw(m.c, m.c.Bounds(), m.ctx.Image(), image.Point{X: 0, Y: 0}, draw.Over)
}

func (m *MenuSettingSlider) Logic() {

}

func (m *MenuSettingSlider) HandleKey(key entity.Button) bool {
    switch key {
    case entity.ButtonUp:
        if m.draftValue < m.setting.GetMaxValue() {
            m.draftValue++
        }
    case entity.ButtonDown:
        if m.draftValue > m.setting.GetMinValue() {
            m.draftValue--
        }
    }
    return true
}

func NewMenuSettingSlider(c *clock.Clock, setting *clock.IntRangeSetting) *MenuSettingSlider {
    return &MenuSettingSlider{
        clk:     c,
        c:       c.Canvas,
        setting: setting,
    }
}
