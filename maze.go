package main

import (
	"os"
	"fmt"
	"image"
	_ "image/png"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

// ----------- Player ----------- //
type Direction int

type player struct {
	sprites		map[Direction][]pixel.Rect
	currentSprite	pixel.Rect
	spriteMap	pixel.Picture
	grid_pos_X	int
	grid_pos_Y	int
}

// --------- Background --------- //
type background struct {
	sprites		map[int][]pixel.Rect
}

type block struct {
  currentSprite	pixel.Rect
  spriteMap	pixel.Picture
  gridX		int
  gridY		int
}

// ---------- Constants --------- //
const (
	// Window Size
	screen_height = 800
	screen_width  = 800

	// Directions
	up	Direction = 0
	down	Direction = 1
	left	Direction = 2
	right	Direction = 3

)

// --------- Variables ---------- //
var (

	// Background
	// 0 = path
	// 1 = light green tree		2 = pink tree
	// 3 = dark green tree		4 = middle green tree
	// 10 x 10
	backgroundMap [][]uint8 = [][]uint8{
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1, 2, 0, 0, 0, 1},
				{1, 0, 0, 1, 3, 1, 3, 0, 0, 1},
				{1, 0, 0, 2, 4, 2, 4, 0, 0, 1},
				{1, 0, 0, 0, 2, 4, 0, 0, 0, 0},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	// Grid size is defined by : X = number of objects per line	Y = number of objects in the slice
	grid_size_x	int = len(backgroundMap[0])
	grid_size_y	int = len(backgroundMap)

	// Score
	score int = 0

	// Keyboard
	keys		map[Direction][]bool
)

// ----------------------- Common Functions ----------------------- //

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

// Get the correct sprite in SpriteMap, based on it size and coordinates
func setSprite(spriteWidth float64, spriteHeight float64, posX int, posY int) pixel.Rect {
	// 4 points of the rect
	return pixel.R(
		float64(posX)*spriteWidth,
		float64(posY)*spriteHeight,
		float64(posX+1)*spriteWidth,
		float64(posY+1)*spriteHeight,
	)
}

// Get the coordinates to draw objects on screen
func getObjectGridPosition(width float64, height float64, grid_x_size int, grid_y_size int, x int, y int) pixel.Rect {
	gridWidth := width / float64(grid_x_size)
	gridHeight := height / float64(grid_y_size)
	return pixel.R(float64(x)*gridWidth, float64(y)*gridHeight, float64((x+1))*gridWidth, float64((y+1))*gridHeight)
}


// ---------------------------- Player ---------------------------- //

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

// Draw Player on screen
func (p0 *player) draw(win pixel.Target) {
	sprite := pixel.NewSprite(nil, pixel.Rect{})
	sprite.Set(p0.spriteMap, p0.currentSprite)
	pos := getObjectGridPosition(screen_width, screen_height, grid_size_x, grid_size_y, p0.grid_pos_X, p0.grid_pos_Y)
	sprite.Draw(win, pixel.IM.ScaledXY(pixel.ZV, pixel.V(     pos.W()/sprite.Frame().W(),        pos.H()/sprite.Frame().H(),    ) ).Moved(pos.Center()),
	)
}

// Update the grid position accordingly to the direction of the next frame
func (p0 *player) getNewGridPos(direction Direction) (int, int) {
	if direction == right {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if p0.grid_pos_X + 1 < grid_size_x && backgroundMap[len(backgroundMap[0]) - 1 - p0.grid_pos_Y][p0.grid_pos_X + 1] == 0 {
			p0.grid_pos_X += 1
		}
		return p0.grid_pos_X, p0.grid_pos_Y
	}
	if direction == left {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if p0.grid_pos_X - 1 >= 0 && backgroundMap[len(backgroundMap[0]) - 1 - p0.grid_pos_Y][p0.grid_pos_X - 1] == 0 {
			p0.grid_pos_X -= 1
		}
		return p0.grid_pos_X, p0.grid_pos_Y
	}
	if direction == up {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if p0.grid_pos_Y + 1 < grid_size_x && backgroundMap[len(backgroundMap) - 1 - (p0.grid_pos_Y + 1)][p0.grid_pos_X] == 0 {
			p0.grid_pos_Y += 1
		}
		return p0.grid_pos_X, p0.grid_pos_Y
	}
	if direction == down {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if p0.grid_pos_Y - 1 >= 0 && backgroundMap[len(backgroundMap) - 1 - (p0.grid_pos_Y - 1)][p0.grid_pos_X] == 0 {
			p0.grid_pos_Y -= 1
		}
		return p0.grid_pos_X, p0.grid_pos_Y
	}
	return p0.grid_pos_X, p0.grid_pos_Y
}

// Update the direction, position on grid and the current sprite each frame
func (p0 *player) update(direction Direction) {
	p0.grid_pos_X, p0.grid_pos_Y = p0.getNewGridPos(direction)
	p0.currentSprite = p0.sprites[direction][0]
	score = p0.grid_pos_X
}


// -------------------------- Background -------------------------- //

// Set board sprites in a map based on its direction
func (bgd *background) setPlayerSprites(spriteMapImg pixel.Picture) {

	// X,Y(Size of each pixel), X, Y(Position in spriteMap)
	bgd.sprites = make(map[int][]pixel.Rect)
	bgd.sprites[0] = append(bgd.sprites[0], setSprite(130, 150, 0, 2))
	bgd.sprites[1] = append(bgd.sprites[1], setSprite(130, 150, 1, 2))
	bgd.sprites[2] = append(bgd.sprites[2], setSprite(130, 150, 2, 2))
	bgd.sprites[3] = append(bgd.sprites[3], setSprite(130, 150, 3, 2))
	bgd.sprites[4] = append(bgd.sprites[4], setSprite(50, 70, 0, 3))
}

// Draw a single block of the background
func (blk block) draw(t pixel.Target) {
    sprite := pixel.NewSprite(nil, pixel.Rect{})
    sprite.Set(blk.spriteMap, blk.currentSprite)
    pos := getObjectGridPosition(screen_width, screen_height, len(backgroundMap[0]), len(backgroundMap), blk.gridY, blk.gridX)

    sprite.Draw(t, pixel.IM.
        ScaledXY(pixel.ZV, pixel.V(
            pos.W()/sprite.Frame().W(),
            pos.H()/sprite.Frame().H(),
        )).
        Moved(pos.Center()),
    )
}

// Draw blocks into the background
func (bgd *background) draw(t pixel.Target) error {
	for i := 0; i < len(backgroundMap); i++ {				// Lines
		for j := 0; j < len(backgroundMap[0]); j++ {	// Columns
			if backgroundMap[i][j] == 0 {
				// Don't draw anything, its the path
			} else if backgroundMap[i][j] == 1 {
				b:=block{currentSprite: bgd.sprites[0][0], gridX:(len(backgroundMap) -1 ) -i, gridY:j}
				b.draw(t)
			} else if backgroundMap[i][j] == 2 {
				b:=block{currentSprite: bgd.sprites[1][0], gridX:(len(backgroundMap) -1 ) -i, gridY:j}
				b.draw(t)
			} else if backgroundMap[i][j] == 3 {
				b:=block{currentSprite: bgd.sprites[2][0], gridX:(len(backgroundMap) -1 ) -i, gridY:j}
				b.draw(t)
			} else if backgroundMap[i][j] == 4 {
				b:=block{currentSprite: bgd.sprites[3][0], gridX:(len(backgroundMap) -1 ) -i, gridY:j}
				b.draw(t)
			}
		}
	}
	return nil
}



// ------------------------ PixelGL Window ------------------------ //
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

	// Disable on screen mouse cursor
	win.SetCursorVisible(false)

	// Initialize Keyboard

	// X,Y(Size of each pixel), X, Y(Position in spriteMap)
	keys = make(map[Direction][]bool)
	keys[up]		= append(keys[up], false)
	keys[down]	= append(keys[down], false)
	keys[left]	= append(keys[left], false)
	keys[right]	= append(keys[right], false)

	// Load the PixelMap Image
	spriteMap, err := loadPicture("spritemap-rpg.png")

	// Initialize Player data
	p0 := &player{}
	// Initial Position
	p0.grid_pos_X = 0
	p0.grid_pos_Y = 7
	// Load the Player Sprites in a map
	p0.setPlayerSprites(spriteMap)
	// Initial Direction
	direction := right	// To identify the initial sprite
	p0.currentSprite = p0.sprites[direction][0]

	// Initialize the background
	bgd := &background{}
	bgd.setPlayerSprites(spriteMap)

	// Draw all background objects first to this object and just draw to window one time later
	imd := imdraw.New(spriteMap)


	// Infinite loop
	for !win.Closed() {

		// Esc to quit program
    if win.JustPressed(pixelgl.KeyEscape) {
      break
    }

		// Clear Screen
		win.Clear(colornames.Lightgreen)

		// ---------------------- Keyboard ---------------------- //

		// Update player direction and keys pressed
		if win.JustPressed(pixelgl.KeyLeft) {
			direction = left
			keys[direction][0] = true

		}
		if win.JustPressed(pixelgl.KeyRight) {
			direction = right
			keys[direction][0] = true
		}
		if win.JustPressed(pixelgl.KeyUp) {
			direction = up
			keys[direction][0] = true

		}
		if win.JustPressed(pixelgl.KeyDown) {
			direction = down
			keys[direction][0] = true
		}

		// Move Player - Necessary for the automation of player execution
		if keys[up][0] == true {
			p0.update(up)
		}
		if keys[down][0] == true {
			p0.update(down)
		}
		if keys[left][0] == true {
			p0.update(left)
		}
		if keys[right][0] == true {
			p0.update(right)
		}

		// Clean key pressed for the next cycle
		keys[up][0] = false
		keys[down][0] = false
		keys[left][0] = false
		keys[right][0] = false

		// -------------------- Draw Objects -------------------- //

		// Draw the entire background
		bgd.draw(imd)
		// Draw with just one draw() call to screen
		imd.Draw(win)

		// Draw Player on the screen
		p0.draw(win)

		// Update the screen
		win.Update()

	}
}

// Main function
func main() {
	pixelgl.Run(run)
}
