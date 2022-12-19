package face

import (
    _ "embed"
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
    "math/rand"
    "strconv"
    "time"
)

const pPad = 1
const pHeight = 15
const pWidth = 5
const pMax = 64 - pHeight + (pPad * 2)

const ballSize = 4

type Ball struct {
    X  int
    Y  int
    TX int
    TY int
}

type Paddle struct {
    Pos   int
    Goal  int
    Score int
    Lose  bool
}

type PongGame struct {
    Left  *Paddle
    Right *Paddle
    Ball  *Ball
    c     *rgbmatrix.Canvas
    ctx   *gg.Context
}

var scoreFont font.Face

func NewPongGame(c *clock.Clock) *PongGame {
    now := time.Now()

    return &PongGame{
        Left: &Paddle{
            Pos:   0,
            Goal:  64,
            Score: now.Hour(),
            Lose:  false,
        },
        Right: &Paddle{
            Pos:   0,
            Goal:  64,
            Score: now.Minute(),
            Lose:  false,
        },
        Ball: &Ball{X: 63, Y: 31, TX: 1, TY: 1},
        c:    c.Canvas,
    }
}

func (pg *PongGame) Init() {
    scoreFontObj, _ := truetype.Parse(asset.PongFontData)

    scoreFont = truetype.NewFace(scoreFontObj, &truetype.Options{
        Size: 12,
    })

    pg.ctx = gg.NewContext(128, 64)
    pg.ctx.SetFontFace(scoreFont)
}

func (pg *PongGame) Stop() {
    //scoreFont = nil
    //pg.ctx = nil
}

func (pg *PongGame) Render() {
    zp := image.Point{X: 0, Y: 0}
    colour := image.Uniform{C: color.White}
    // Left paddle
    draw.Draw(pg.c, image.Rect(pPad, pPad+pg.Left.Pos, pPad+pWidth, pPad+pg.Left.Pos+pHeight), &colour, zp, draw.Over)
    // Right Paddle
    draw.Draw(pg.c, image.Rect(128-pPad, pPad+pg.Right.Pos, 128-pPad-pWidth, pPad+pg.Right.Pos+pHeight), &colour, zp, draw.Over)
    // Ball
    draw.Draw(pg.c, image.Rect(pg.Ball.X, pg.Ball.Y, pg.Ball.X+ballSize, pg.Ball.Y+ballSize), &colour, zp, draw.Over)

    // Text

    scoreText := strconv.Itoa(pg.Left.Score)
    if pg.Left.Score < 10 {
        scoreText = "0" + scoreText
    }
    if pg.Right.Score < 10 {
        scoreText += "0"
    }
    scoreText += strconv.Itoa(pg.Right.Score)
    pg.ctx.SetColor(color.Transparent)
    pg.ctx.Clear()
    pg.ctx.SetColor(color.White)
    pg.ctx.DrawStringAnchored(scoreText, 64, 2, 0.5, 1)
    pg.ctx.Stroke()
    draw.Draw(pg.c, pg.c.Bounds(), pg.ctx.Image(), zp, draw.Over)
}

func (pg *PongGame) Logic() {
    if pg.Ball.TX == 0 && pg.Ball.TY == 0 {
        return
    }

    now := time.Now()

    // Left gains a point when the hour changes, right gains a point when the minute changes
    pg.Right.Lose = pg.Left.Score != now.Hour()
    pg.Left.Lose = !pg.Right.Lose && pg.Right.Score != now.Minute()

    // Move the ball forward by its trajectory
    pg.Ball.X += pg.Ball.TX
    pg.Ball.Y += pg.Ball.TY

    // Move the paddles towards their goal
    pg.movePaddle(pg.Left)
    pg.movePaddle(pg.Right)

    if pg.Ball.X < 64 {
        pg.setPaddleGoals(pg.Left, pg.Right)
    } else {
        pg.setPaddleGoals(pg.Right, pg.Left)
    }

    // Deflect the ball of the top and bottom of the display
    if !pg.valueWithin(pg.Ball.Y, 0, 64-ballSize) {
        pg.Ball.TY *= -1
    }

    ballRect := image.Rect(pg.Ball.X, pg.Ball.Y, pg.Ball.X+ballSize, pg.Ball.Y+ballSize)
    // Left paddle deflect
    if pg.Ball.X <= pPad+pWidth && !ballRect.Intersect(image.Rect(pPad, pPad+pg.Left.Pos, pPad+pWidth, pPad+pg.Left.Pos+pHeight)).Empty() {
        pg.Ball.TX = 1
        if rand.Intn(2) == 1 {
            pg.Ball.TY = -1
        } else {
            pg.Ball.TY = 1
        }
    }

    // Right paddle deflect
    if pg.Ball.X >= 128-pWidth-ballSize && !ballRect.Intersect(image.Rect(128-pPad, pPad+pg.Right.Pos, 128-pPad-pWidth, pPad+pg.Right.Pos+pHeight)).Empty() {
        pg.Ball.TX = -1
        if rand.Intn(2) == 1 {
            pg.Ball.TY = -1
        } else {
            pg.Ball.TY = 1
        }
    }

    // If the ball hits the left and right walls, reset the ball
    if !pg.valueWithin(pg.Ball.X, 1, 128-ballSize-1) {
        pg.Ball.TX = 0
        pg.Ball.TY = 0
        time.Sleep(1 * time.Second)
        pg.Left.Score = now.Hour()
        pg.Right.Score = now.Minute()
        pg.Ball.X = 63
        pg.Ball.Y = 31
        pg.Ball.TX = 1
    }
}

func (pg *PongGame) HandleKey(key entity.Button) bool {
    return false
}

func (pg *PongGame) movePaddle(p *Paddle) {
    if p.Pos > p.Goal && p.Pos > 0 {
        p.Pos--
    } else if p.Pos < p.Goal && p.Pos < pMax {
        p.Pos++
    }
}

func (pg *PongGame) setPaddleGoals(relevant *Paddle, other *Paddle) {
    newGoal := pg.Ball.Y - pHeight/4 - rand.Intn(pHeight/2)
    if math.Abs(float64(newGoal-relevant.Goal)) > ballSize {
        relevant.Goal = newGoal
        // If the relevant paddle is supposed to lose, move the paddle just out of reach
        if relevant.Lose {
            if relevant.Goal > 64 {
                relevant.Goal -= pHeight
            } else {
                relevant.Goal += pHeight
            }
        }
    }

    if math.Abs(float64(other.Goal-other.Pos)) < ballSize && rand.Intn(20) < 2 {
        other.Goal = rand.Intn(64 - pHeight)
    }
}

func (pg *PongGame) valueWithin(value int, min int, max int) bool {
    return value < max && value > min
}
