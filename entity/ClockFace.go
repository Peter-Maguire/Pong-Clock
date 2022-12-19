package entity

type ClockFace interface {
    Init()
    Stop()
    Render()
    Logic()
    HandleKey(key Button) bool
}
