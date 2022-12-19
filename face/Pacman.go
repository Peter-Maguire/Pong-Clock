package face

import (
    "bytes"
    _ "embed"
    "github.com/peter-maguire/pong-clock/asset"
    "github.com/peter-maguire/pong-clock/clock"
    "github.com/peter-maguire/pong-clock/entity"
    rgbmatrix "github.com/peter-maguire/pong-clock/lib/go-rpi-rgb-led-matrix"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "math"
    "time"
)

var spritesheet image.Image

type Pacman struct {
    c         *rgbmatrix.Canvas
    t         *time.Ticker
    animTimer int
    Sprites   []*Character
}

type Facing int

const (
    FacingRight Facing = iota
    FacingUp
    FacingLeft
    FacingDown
)

type Sprite int

const (
    SpritePac1 Sprite = iota
    SpritePac2
    SpritePac3
    SpritePac4
    SpriteBlinky
    SpritePinky
    SpriteInky
    SpriteFruit
)

const wHeight = 12
const wLength = 10

type Character struct {
    X      int
    Y      int
    Step   int
    Goals  [][]int
    Facing Facing
    Sprite Sprite
}

func (c *Character) GetSprite() (image.Rectangle, image.Point) {
    return image.Rect(c.X, c.Y, c.X+8, c.Y+8), image.Pt(8*int(c.Sprite), 8*int(c.Facing))
}

func NewPacman(c *clock.Clock) *Pacman {
    return &Pacman{
        c: c.Canvas,
        Sprites: []*Character{{
            X:      0,
            Y:      3,
            Step:   0,
            Goals:  [][]int{},
            Facing: FacingRight,
            Sprite: SpritePac1,
        }},
    }
}

func (p *Pacman) Init() {
    spritesheet, _ = png.Decode(bytes.NewReader(asset.PacmanSpritesheet))
    go p.startTimer()
    p.setGoals(1, 1, numbers[0])
    p.setGoals(1+(3*8), 1, numbers[1])
    p.setGoals(1+(2*3*8), 1, numbers[2])
}

func (p *Pacman) Stop() {
    p.t.Stop()
    spritesheet = nil
}

func (p *Pacman) startTimer() {
    p.t = time.NewTicker(50 * time.Millisecond)
    for range p.t.C {
        p.animTimer++
        p.moveSprites()
    }
}

func (p *Pacman) moveSprites() {
    pacman := p.Sprites[0]
    pacman.Sprite = Sprite(p.animTimer % 4)
    for _, sprite := range p.Sprites {
        if sprite.Step >= len(sprite.Goals) {
            continue
        }
        goal := sprite.Goals[sprite.Step]
        goalX := goal[0]
        goalY := goal[1]
        if goalX > sprite.X {
            sprite.X++
            sprite.Facing = FacingRight
        } else if goalX < sprite.X {
            sprite.X--
            sprite.Facing = FacingLeft
        } else if goalY > sprite.Y {
            sprite.Y++
            sprite.Facing = FacingDown
        } else if goalY < sprite.Y {
            sprite.Y--
            sprite.Facing = FacingUp
        }
        if (goalX == sprite.X && math.Abs(float64(goalY-sprite.Y)) == 1) || (goalY == sprite.Y && math.Abs(float64(goalX-sprite.X)) == 1) {
            sprite.Step++
        }

    }
}

func (p *Pacman) drawWall(x int, y int, c color.Color, config []string) {
    for rowNum, row := range config {
        for colNum, chr := range row {
            if chr != '-' {
                for ax := 0; ax < wLength; ax++ {
                    // Top wall
                    if getUpperWall(rowNum, colNum, config) == '-' {
                        p.c.Set(x+ax+(wLength*colNum), y+(wHeight*rowNum), c)
                    }
                    // Bottom wall
                    if getLowerWall(rowNum, colNum, config) == '-' {
                        p.c.Set(x+ax+(wLength*colNum), wHeight-1+y+(wHeight*rowNum), c)
                    }
                }

                for ay := 0; ay < wHeight; ay++ {
                    // Left Wall
                    if getLeftWall(colNum, row) == '-' {
                        p.c.Set(x+(wLength*colNum), ay+y+(wHeight*rowNum), c)
                    }
                    // Right Wall
                    if getRightWall(colNum, row) == '-' {
                        p.c.Set(x+wLength-1+(wLength*colNum), ay+y+(wHeight*rowNum), c)
                    }
                }
            }
        }
    }
}

func (p *Pacman) setGoals(x int, y int, config []string) {
    pacman := p.Sprites[0]
    rowNum := 0
    colNum := 0
    //lastRowNum := 0
    lastColNum := 0
    for {
        rowDelta := 0
        colDelta := 0
        chr := config[rowNum][colNum]
        switch chr {
        case '>':
            colDelta = 1
        case '<':
            colDelta = -1
        case '=':
            if lastColNum < colNum {
                colDelta = 1
            } else {
                colDelta = -1
            }
        case '^':
            rowDelta = -1
        case 'v':
            rowDelta = 1
        default:
            rowDelta = 1
        }
        //lastRowNum = rowNum
        lastColNum = colNum
        rowNum += rowDelta
        colNum += colDelta
        pacman.Goals = append(pacman.Goals, []int{x + (colNum * wLength), y + (rowNum * wHeight)})
        if rowNum < 0 || colNum < 0 || rowNum > len(config)-1 || colNum > len(config[rowNum])-1 {
            break
        }
    }
}

func getLeftWall(colNum int, row string) rune {
    if colNum == 0 {
        return '-'
    }
    return rune(row[colNum-1])
}

func getRightWall(colNum int, row string) rune {
    if colNum >= len(row)-1 {
        return '-'
    }
    return rune(row[colNum+1])
}

func getLowerWall(rowNum int, colNum int, config []string) rune {
    if rowNum == len(config)-1 {
        return '-'
    }
    return rune(config[rowNum+1][colNum])
}

func getUpperWall(rowNum int, colNum int, config []string) rune {
    if rowNum == 0 {
        return '-'
    }
    return rune(config[rowNum-1][colNum])
}

// 0000
// LRUD

func (p *Pacman) Render() {
    for i := 0; i < 4; i++ {
        p.drawWall(0+(wLength*3*i)+(i*2), 1, color.RGBA{
            R: 0,
            G: 0,
            B: 255,
            A: 255,
        }, numbers[i])
    }

    for _, sprite := range p.Sprites {
        bounds, point := sprite.GetSprite()
        draw.Draw(p.c, bounds, spritesheet, point, draw.Over)
        for _, g := range sprite.Goals {
            draw.Draw(p.c, image.Rect(g[0], g[1], g[0]+2, g[1]+2), &image.Uniform{C: color.White}, image.Pt(0, 0), draw.Over)
        }
    }
}

func (p *Pacman) Logic() {

}

func (p *Pacman) HandleKey(key entity.Button) bool {
    return false
}
