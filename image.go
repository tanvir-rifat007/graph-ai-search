package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/StephaneBunel/bresenham"
	"github.com/kettek/apng"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Constant.
const cellSize = 60

// Neon cyberpunk color palette
var (
	// Background - deep space blue/black
	bgColor = color.RGBA{R: 10, G: 10, B: 25, A: 255}

	// Walls - dark purple
	wallColor = color.RGBA{R: 25, G: 15, B: 45, A: 255}

	// Solution path - electric cyan
	solutionColor = color.RGBA{R: 0, G: 255, B: 255, A: 255}

	// Start point - neon green
	startColor = color.RGBA{R: 57, G: 255, B: 20, A: 255}

	// Goal point - hot pink/magenta
	goalColor = color.RGBA{R: 255, G: 0, B: 128, A: 255}

	// Current node - electric purple
	currentColor = color.RGBA{R: 191, G: 64, B: 191, A: 255}

	// Explored cells - deep blue with purple tint
	exploredColor = color.RGBA{R: 75, G: 50, B: 150, A: 255}

	// Empty cells - darker blue
	emptyColor = color.RGBA{R: 20, G: 25, B: 45, A: 255}

	// Grid lines - bright cyan with glow effect
	gridColor = color.RGBA{R: 0, G: 200, B: 255, A: 255}

	// Text color - bright cyan
	textColor = color.RGBA{R: 0, G: 255, B: 255, A: 255}
)

// OutputImage draw the maze as png file with neon theme.
func (g *Maze) OutputImage(fileName ...string) {
	fmt.Printf("🎨 Generating neon cyberpunk maze image %s...\n", fileName)
	width := cellSize * (g.Width - 1)
	height := cellSize * g.Height

	var outFile = "image.png"
	if len(fileName) > 0 {
		outFile = fileName[0]
	}

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// Dark cyberpunk background
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	// draw squares on the image with neon colors
	for i, row := range g.Walls {
		for j, col := range row {
			p := Point{
				X: i,
				Y: j,
			}
			if col.wall {
				// draw dark purple square for wall
				g.drawSquare(col, p, img, wallColor, cellSize, j*cellSize, i*cellSize)
			} else if g.inSolution(p) {
				// part of solution, so draw electric cyan square
				g.drawSquare(col, p, img, solutionColor, cellSize, j*cellSize, i*cellSize)
			} else if col.State.X == g.Start.X && col.State.Y == g.Start.Y {
				// Starting point, so draw neon green square
				g.drawSquare(col, p, img, startColor, cellSize, j*cellSize, i*cellSize)
			} else if col.State.X == g.Goal.X && col.State.Y == g.Goal.Y {
				// Ending point. Draw hot pink square
				g.drawSquare(col, p, img, goalColor, cellSize, j*cellSize, i*cellSize)
			} else if col.State == g.CurrentNode.State {
				// Current location. Draw in electric purple
				g.drawSquare(col, p, img, currentColor, cellSize, j*cellSize, i*cellSize)
			} else if inExplored(Point{i, j}, g.Explored) {
				// An explored cell - deep blue purple
				g.drawSquare(col, p, img, exploredColor, cellSize, j*cellSize, i*cellSize)
			} else {
				// empty, unexplored. Draw in dark blue
				g.drawSquare(col, p, img, emptyColor, cellSize, j*cellSize, i*cellSize)
			}
		}
	}

	// draw a glowing grid with neon cyan lines
	for i, _ := range g.Walls {
		bresenham.DrawLine(img, 0, i*cellSize, g.Width*cellSize, i*cellSize, gridColor)
	}

	for i := 0; i <= g.Width; i++ {
		bresenham.DrawLine(img, i*cellSize, 0, i*cellSize, g.Height*cellSize, gridColor)
	}

	f, _ := os.Create(outFile)
	_ = png.Encode(f, img)
}

// drawSquare with neon styling
func (g *Maze) drawSquare(col Wall, p Point, img *image.RGBA, c color.Color, size, x, y int) {
	patch := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(patch, patch.Bounds(), &image.Uniform{
		C: c,
	}, image.Point{}, draw.Src)

	if !col.wall {
		// Choose text color based on background brightness
		var txtColor color.Color
		// Use dark text for bright backgrounds, bright text for dark backgrounds
		if isBrightColor(c) {
			txtColor = color.RGBA{R: 10, G: 10, B: 25, A: 255} // Dark blue text
		} else {
			txtColor = textColor // Cyan text
		}
		g.printLocation(p, txtColor, patch)
	}

	draw.Draw(img, image.Rect(x, y, x+size, y+size), patch, image.Point{}, draw.Src)
}

// printLocation with cyberpunk styling
func (g *Maze) printLocation(p Point, c color.Color, patch *image.RGBA) {
	point := fixed.Point26_6{X: fixed.I(6), Y: fixed.I(40)}
	d := &font.Drawer{
		Dst:  patch,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(fmt.Sprintf("[%d %d]", p.X, p.Y))
}

// isBrightColor determines if a color is bright (needs dark text)
func isBrightColor(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	// Convert to 0-255 range and calculate perceived brightness
	brightness := (r>>8)*299 + (g>>8)*587 + (b>>8)*114
	return brightness > 128000 // Threshold for bright colors
}

// OutputAnimatedImage creates slower animated maze visualization
func (g *Maze) OutputAnimatedImage() {
	g.Animate = true
	output := "./animation.png"
	fmt.Println("🎬 Creating cyberpunk animated maze...")

	files, _ := os.ReadDir("./tmp")

	var images []string
	var delays []int

	for _, file := range files {
		images = append(images, fmt.Sprintf("./tmp/%s", file.Name()))
		// Slower animation: 100ms per frame (was 30ms)
		delays = append(delays, 100)
	}
	images = append(images, "./image.png")

	a := apng.APNG{
		Frames: make([]apng.Frame, len(images)),
	}
	out, _ := os.Create(output)
	defer out.Close()

	for i, s := range images {
		in, err := os.Open(s)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		defer in.Close()

		m, err := png.Decode(in)
		if err != nil {
			continue
		}
		a.Frames[i].Image = m
	}

	err := apng.Encode(out, a)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("✨ Neon maze animation complete!")
}
