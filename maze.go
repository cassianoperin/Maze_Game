package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"golang.org/x/image/colornames"
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Direction int

type player struct {
	direction				Direction
	sprites					map[Direction][]pixel.Rect
	currentSprite		pixel.Rect
	spriteMap				pixel.Picture
	grid_pos_X			int
	grid_pos_Y			int
}

const (
	// Window Size
	screen_height = 800
	screen_width  = 800

	// Directions
	up 		Direction = 0
	down	Direction = 1
	left	Direction = 2
	right	Direction = 3

	// Grid
	grid_size_x	int = 20
	grid_size_y	int = 20
)


// Load Picture File
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Cannot read file:", err)
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}


// Get the correct sprite in SpriteMap in setPlayerSprites()
func setSprite(spriteWidth float64, spriteHeight float64, posX int, posY int) pixel.Rect {
	// 4 points of the rect
	return pixel.R(
		float64(posX)*spriteWidth,
		float64(posY)*spriteHeight,
		float64(posX+1)*spriteWidth,
		float64(posY+1)*spriteHeight,
	)
}


// Set player sprites in a map based on its direction
func (p0 *player) setPlayerSprites(spriteMapImg pixel.Picture) {
	p0.spriteMap = spriteMapImg
	// X,Y(Size of each pixel), X, Y(Position in spriteMap)
	p0.sprites = make(map[Direction][]pixel.Rect)
	p0.sprites[up] = append(p0.sprites[up], setSprite(50, 70, 6, 0))
	p0.sprites[down] = append(p0.sprites[down], setSprite(50, 70, 6, 3))
	p0.sprites[left] = append(p0.sprites[left], setSprite(50, 70, 6, 2))
	p0.sprites[right] = append(p0.sprites[right], setSprite(50, 70, 6, 1))
}


// Get the coordinates to draw Player on screen
func getPlayerGridPosition(width float64, height float64, grid_x_size int, grid_y_size int, x int, y int) pixel.Rect {
	gridWidth := width / float64(grid_x_size)
	gridHeight := height / float64(grid_y_size)
	return pixel.R(float64(x)*gridWidth, float64(y)*gridHeight, float64((x+1))*gridWidth, float64((y+1))*gridHeight)
}


// Draw Player on screen
func (p0 *player) draw(win pixel.Target) {
	sprite := pixel.NewSprite(nil, pixel.Rect{})
	sprite.Set(p0.spriteMap, p0.currentSprite)
	pos := getPlayerGridPosition(screen_width, screen_height, grid_size_x, grid_size_y, p0.grid_pos_X, p0.grid_pos_Y)
	sprite.Draw(win, pixel.IM.ScaledXY(pixel.ZV, pixel.V(     pos.W()/sprite.Frame().W(),        pos.H()/sprite.Frame().H(),    ) ).Moved(pos.Center()),
	)
}


// Update the grid position accordingly to the direction of the next frame
func (p0 *player) getNewGridPos(direction Direction) (int, int) {
	if direction == right {
		return p0.grid_pos_X + 1, p0.grid_pos_Y
	}
	if direction == left {
		return p0.grid_pos_X - 1, p0.grid_pos_Y
	}
	if direction == up {
		return p0.grid_pos_X, p0.grid_pos_Y + 1
	}
	if direction == down {
		return p0.grid_pos_X, p0.grid_pos_Y - 1
	}
	return p0.grid_pos_X, p0.grid_pos_Y
}

// Update the direction, position on grid and the current sprite each frame
func (p0 *player) update(direction Direction) {
	p0.direction = direction
	p0.grid_pos_X, p0.grid_pos_Y = p0.getNewGridPos(direction)
	p0.currentSprite = p0.sprites[p0.direction][0]
}


// PixelGL Window
func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Maze Game",
		Bounds: pixel.R(0, 0, screen_width, screen_height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Load the PixelMap Image
	spriteMap, err := loadPicture("spritemap-rpg.png")

	// Initialize Player data
	p0 := &player{}
	// Initial Position
	p0.grid_pos_X = 1
	p0.grid_pos_Y = 5
	// Load the Player Sprites in a map
	p0.setPlayerSprites(spriteMap)
	// Initial Direction
	direction:=right

	// Infinite loop
	for !win.Closed() {
		// Clear Screen
		win.Clear(colornames.White)

		// Update player direction
		if win.Pressed(pixelgl.KeyLeft) {
			direction = left
		}
		if win.Pressed(pixelgl.KeyRight) {
			direction = right
		}
		if win.Pressed(pixelgl.KeyUp) {
			direction = up
		}
		if win.Pressed(pixelgl.KeyDown) {
			direction = down
		}
		p0.update(direction)

		// Draw Player on the screen
		p0.draw(win)

		// Update the screen
		win.Update()

		// Control speed of movement
		// TODO REPLACE BY A TIMER
		time.Sleep(200 * time.Millisecond)
	}
}

// Main function
func main() {
	pixelgl.Run(run)
}
