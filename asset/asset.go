package asset

import _ "embed"

//go:embed arial.ttf
var MenuFontData []byte

//go:embed pong_font.ttf
var PongFontData []byte

//go:embed pacman.png
var PacmanSpritesheet []byte
