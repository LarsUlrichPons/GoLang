package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 500
	sampleRate   = 44100
)

var audioContext *audio.Context

// AUDIO STREAM (Same as before)
type stream struct {
	freq float64
	pos  int64
}

func (s *stream) Read(b []byte) (int, error) {
	for i := 0; i < len(b)/4; i++ {
		v := math.Sin(2 * math.Pi * float64(s.pos) * s.freq / sampleRate)
		s.pos++
		v16 := int16(v * 32767 * 0.1)
		b[4*i] = byte(v16)
		b[4*i+1] = byte(v16 >> 8)
		b[4*i+2] = byte(v16)
		b[4*i+3] = byte(v16 >> 8)
	}
	return len(b), nil
}

// Added 'label' to store the text (e.g., "C4")
type KeyInfo struct {
	ebitenKey ebiten.Key
	freq      float64
	isBlack   bool
	drawIndex int
	label     string
}

type Game struct {
	keys       []KeyInfo
	players    map[ebiten.Key]*audio.Player
	activeNote string
	activeFreq float64
}

func NewGame() *Game {
	audioContext = audio.NewContext(sampleRate)

	return &Game{
		players: make(map[ebiten.Key]*audio.Player),
		keys: []KeyInfo{
			// WHITE KEYS (C4 to B4)
			{ebiten.KeyA, 261.63, false, 0, "C4"},
			{ebiten.KeyS, 293.66, false, 1, "D4"},
			{ebiten.KeyD, 329.63, false, 2, "E4"},
			{ebiten.KeyF, 349.23, false, 3, "F4"},
			{ebiten.KeyG, 392.00, false, 4, "G4"},
			{ebiten.KeyH, 440.00, false, 5, "A4"},
			{ebiten.KeyJ, 493.88, false, 6, "B4"},

			// BLACK KEYS
			{ebiten.KeyW, 277.18, true, 1, "C#4"},
			{ebiten.KeyE, 311.13, true, 2, "D#4"},
			{ebiten.KeyT, 369.99, true, 4, "F#4"},
			{ebiten.KeyY, 415.30, true, 5, "G#4"},
			{ebiten.KeyU, 466.16, true, 6, "A#4"},
		},
		activeNote: "---",
		activeFreq: 0.0,
	}
}

func (g *Game) Update() error {
	pressedAny := false

	for _, k := range g.keys {
		if ebiten.IsKeyPressed(k.ebitenKey) {
			pressedAny = true

			// Update the LCD display data
			g.activeNote = k.label
			g.activeFreq = k.freq

			if g.players[k.ebitenKey] == nil {
				p, _ := audioContext.NewPlayer(&stream{freq: k.freq})
				g.players[k.ebitenKey] = p
				p.Play()
			}
		} else {
			if p, ok := g.players[k.ebitenKey]; ok && p != nil {
				p.Close()
				g.players[k.ebitenKey] = nil
			}
		}
	}

	// Reset LCD if no keys are pressed
	if !pressedAny {
		g.activeNote = "---"
		g.activeFreq = 0.0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 1. Background (Dark Carbon look)
	screen.Fill(color.RGBA{25, 25, 25, 255})

	// 2. Draw LCD Screen at Top Center
	lcdWidth, lcdHeight := 300.0, 50.0
	lcdX := (screenWidth - lcdWidth) / 2
	lcdY := 20.0

	// Green LCD Background
	ebitenutil.DrawRect(screen, lcdX, lcdY, lcdWidth, lcdHeight, color.RGBA{50, 255, 50, 255})

	// LCD Text
	lcdText := fmt.Sprintf("Note: %s  |  %.2f Hz", g.activeNote, g.activeFreq)
	if g.activeFreq == 0 {
		lcdText = "Ready to Play"
	}
	text.Draw(screen, lcdText, basicfont.Face7x13, int(lcdX+50), int(lcdY+30), color.Black)

	// Calculate Key Layout constraints
	pianoY := 120.0
	whiteKeyWidth := float64(screenWidth) / 7.0
	blackKeyWidth := whiteKeyWidth * 0.6

	// 3. Draw White Keys
	for _, k := range g.keys {
		if k.isBlack {
			continue
		}
		x := float64(k.drawIndex) * whiteKeyWidth

		// Draw Black Border underneath
		ebitenutil.DrawRect(screen, x, pianoY, whiteKeyWidth, screenHeight-pianoY, color.Black)

		// Draw Inner White Key (slightly smaller to reveal the border)
		drawColor := color.RGBA{250, 250, 250, 255}
		if ebiten.IsKeyPressed(k.ebitenKey) {
			drawColor = color.RGBA{200, 200, 200, 255} // Grey when pressed
		}
		ebitenutil.DrawRect(screen, x+2, pianoY, whiteKeyWidth-4, screenHeight-pianoY-2, drawColor)

		// Draw Red Label at Bottom
		text.Draw(screen, k.label, basicfont.Face7x13, int(x+whiteKeyWidth/2-10), int(screenHeight-20), color.RGBA{200, 0, 0, 255})
	}

	// 4. Draw Black Keys
	for _, k := range g.keys {
		if !k.isBlack {
			continue
		}
		x := float64(k.drawIndex)*whiteKeyWidth - (blackKeyWidth / 2)

		drawColor := color.RGBA{30, 30, 30, 255}
		if ebiten.IsKeyPressed(k.ebitenKey) {
			drawColor = color.RGBA{80, 80, 80, 255} // Lighter grey when pressed
		}

		keyHeight := 220.0
		ebitenutil.DrawRect(screen, x, pianoY, blackKeyWidth, keyHeight, drawColor)

		// Draw White Label on Black Key
		text.Draw(screen, k.label, basicfont.Face7x13, int(x+blackKeyWidth/2-12), int(pianoY+keyHeight-20), color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Realistic Virtual Piano")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
