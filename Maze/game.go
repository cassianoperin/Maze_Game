package Maze

import (
	"os"
	"fmt"
  "strings"
	"math"
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
	score int
	max_ind_position int
	max_ind_position_cycle int
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
	// IA
	automation	bool = true
	commands_matrix	[][]Direction

	// Slice of players
	player_list []*player
	direction	Direction

	// Background
	// 0 = path
	// 1 = light green tree		2 = pink tree
	// 3 = dark green tree		4 = middle green tree
	// 10 x 10
	// backgroundMap [][]uint8 = [][]uint8{
	// 			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	// 			{1, 0, 2, 0, 0, 0, 0, 2, 0, 1},
	// 			{0, 0, 0, 0, 1, 0, 0, 3, 0, 1},
	// 			{1, 0, 0, 0, 3, 0, 0, 0, 0, 1},
	// 			{1, 0, 0, 2, 4, 2, 0, 0, 0, 1},
	// 			{1, 0, 0, 0, 2, 0, 0, 0, 0, 0},
	// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 			{1, 0, 1, 0, 0, 0, 0, 3, 2, 1},
	// 			{1, 0, 1, 0, 0, 0, 0, 3, 2, 1},
	// 			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	// }

	// 15 x 10
	backgroundMap [][]uint8 = [][]uint8{
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 0, 3, 0, 0, 0, 0, 0, 2, 3, 1, 0, 0, 0, 1},
				{0, 0, 0, 0, 0, 0, 0, 0, 4, 3, 2, 0, 0, 0, 1},
				{1, 0, 0, 0, 4, 1, 3, 0, 0, 1, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 2, 0, 1},
				{1, 1, 0, 0, 1, 3, 2, 0, 0, 0, 0, 0, 4, 0, 1},
				{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 2, 3, 0, 1},
				{1, 0, 1, 3, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0},
				{1, 0, 2, 4, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	// // 15 x 10 Empty
	// backgroundMap [][]uint8 = [][]uint8{
	// 			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	// }
	// Grid size is defined by : X = number of objects per line	Y = number of objects in the slice
	grid_size_x	int = len(backgroundMap[0])
	grid_size_y	int = len(backgroundMap)

	// Score
	max_generation_position int = 0

	// Keyboard
	keyboard_human		map[Direction][]bool	// Human Keyboard
	keyboard_automations		map[Direction][]bool	// Automation Virtual Keyboards

	// Cycle counter
	cycle int = 0
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
func (object *player) setPlayerSprites(spriteMapImg pixel.Picture) {
	object.spriteMap = spriteMapImg
	// X,Y(Size of each pixel), X, Y(Position in spriteMap)
	object.sprites = make(map[Direction][]pixel.Rect)
	object.sprites[up] = append(object.sprites[up], setSprite(50, 70, 6, 0))
	object.sprites[down] = append(object.sprites[down], setSprite(50, 70, 6, 3))
	object.sprites[left] = append(object.sprites[left], setSprite(50, 70, 6, 2))
	object.sprites[right] = append(object.sprites[right], setSprite(50, 70, 6, 1))
}

// Draw Player on screen
// func (p0 *player) draw(win pixel.Target) {
func (object *player) draw(win pixel.Target) {
	sprite := pixel.NewSprite(nil, pixel.Rect{})
	sprite.Set(object.spriteMap, object.currentSprite)
	pos := getObjectGridPosition(screen_width, screen_height, grid_size_x, grid_size_y, object.grid_pos_X, object.grid_pos_Y)
	sprite.Draw(win, pixel.IM.ScaledXY(pixel.ZV, pixel.V(     pos.W()/sprite.Frame().W(),        pos.H()/sprite.Frame().H(),    ) ).Moved(pos.Center()),
	)
}

// Update the grid position accordingly to the direction of the next frame
// Collision Detection
func (object *player) getNewGridPos(direction Direction) (int, int) {
	if direction == right {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if object.grid_pos_X + 1 < grid_size_x  && backgroundMap[len(backgroundMap) - 1 - object.grid_pos_Y][object.grid_pos_X + 1] == 0 {
			object.grid_pos_X += 1
		}
		return object.grid_pos_X, object.grid_pos_Y
	}
	// backgroundMap[line][column]
	if direction == left {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if object.grid_pos_X - 1 >= 0  && backgroundMap[len(backgroundMap) - 1 - object.grid_pos_Y][object.grid_pos_X - 1] == 0{
			object.grid_pos_X -= 1
		}
		return object.grid_pos_X, object.grid_pos_Y
	}
	if direction == up {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if object.grid_pos_Y + 1 < grid_size_y  && backgroundMap[len(backgroundMap) - 1 - (object.grid_pos_Y + 1)][object.grid_pos_X] == 0{
			object.grid_pos_Y += 1
		}
		return object.grid_pos_X, object.grid_pos_Y
	}
	if direction == down {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if object.grid_pos_Y - 1 >= 0  && backgroundMap[len(backgroundMap) - 1 - (object.grid_pos_Y - 1)][object.grid_pos_X] == 0{
			object.grid_pos_Y -= 1
		}
		return object.grid_pos_X, object.grid_pos_Y
	}
	return object.grid_pos_X, object.grid_pos_Y
}

// Update the direction, position on grid and the current sprite each frame
func (object *player) update(direction Direction, player_index int) {
	// Update grid positiom
	object.grid_pos_X, object.grid_pos_Y = object.getNewGridPos(direction)

	// Update current sprite based on direction
	object.currentSprite = object.sprites[direction][0]

	// Test if its new generation record:
	if object.grid_pos_X > max_generation_position {
		max_generation_position =  object.grid_pos_X
	}

	// Update the Individual Maximum score and check the objective!
	if object.grid_pos_X > object.max_ind_position {

		// Update the new maximum posistion of player
		object.max_ind_position = object.grid_pos_X

		// Add score points
		player_score(object.max_ind_position, cycle, player_index)

		// Objective reached!!
		if object.grid_pos_X == len(backgroundMap[0]) - 1 {
			fmt.Printf("\n\n\n\t\tObjective accomplished!\n\t\tIndividual: %s\tPosition: %d\tMovements: %d\n\n\n", population[player_index], len(backgroundMap[0]) - 1, cycle)
		}
	}
	// Punishment
	// } else {
	// 	object.score -= 10
	// }

}


// Calculate the player's score
func player_score (max_pos int, cycle int, plr_index int) {
	var (
		tmp_score float64
		cycles_needed int
	)

	cycles_needed = cycle - player_list[plr_index].max_ind_position_cycle
	// fmt.Printf("Player %d\tCycle: %d\tMaxPos: %d\tcycles needed: %d\t\n",plr_index, cycle, max_pos, cycles_needed)

	tmp_score = (float64(max_pos) / float64(cycles_needed)) * 100
	// fmt.Println(tmp_score)
	player_list[plr_index].score += int(math.Round(tmp_score))

	// fmt.Printf("Individual: %d (%s)\tGeneration:%d\tNew max_pos: %d\tSteps: %d\tNew Score: %d\n",individual_number, population[individual_number], current_generation, max_pos, cycle, score)

	// Update the cycle of last jump (for next avaliation)
	player_list[plr_index].max_ind_position_cycle = cycle

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


func (*player) restart_player(sprMap pixel.Picture, object *player) {
	// Initial Position
	object.grid_pos_X = 0
	object.grid_pos_Y = 7
	// Load the Player Sprites in a map
	object.setPlayerSprites(sprMap)
	// Initial Direction
	direction = right	// To identify the initial sprite
	object.currentSprite = object.sprites[direction][0]
	// Score
	object.score = 0
	object.max_ind_position = 0
	object.max_ind_position_cycle = 0
}




// Convert the binary string of individuals to commands
func individualtoCommands(pop[] string, gene_nr int) [][]Direction {

	// -------------- Prepare the Multidimensional Slice -------------- //
	// Declaring a slice of slices with a length of POPULATION
	commands := make([][]Direction , len(pop))

	// looping through the slice to declare a slice of each slice size
	for i := 0; i < len(pop); i++ {

		// Length of each slice shoulg be a half of gene_number (2 digits = 1 command)
		new_length := gene_nr / 2

    commands[i] = make([]Direction, new_length)
	}

	// ----------------------- Process the data ----------------------- //

	// Decode each individual into commands
	for i := 0 ; i < len(pop) ; i++ {

		individual_split := strings.Split(pop[i], "")

		// Read individual and transform it to a slice with commands
		index := 0
		for j := 0 ; j < len(individual_split) / 2 ; j ++ {
			code := fmt.Sprintf("%s%s", individual_split[index], individual_split[index+1])

			if code == "00" {
				commands[i][j] = 0
			} else if code == "01" {
				commands[i][j] = 1
			} else if code == "10" {
				commands[i][j] = 2
			} else if code == "11" {
				commands[i][j] = 3
			} else {
				fmt.Printf("Value unexpected on binary to command conversion, exiting.\n")
				os.Exit(2)
			}

			index+=2
		}

	}
	return commands
}





// ------------------------ PixelGL Window ------------------------ //
func Run() {

	// ----------------------- Config ----------------------- //

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

	// ------------------------- IA ------------------------- //

	// Validate parameters
	validate_parameters(population_size, k)

	// 0 - Generate the population
	// Generate each individual for population
	for i := 0 ; i < population_size ; i++ {
	  population = append( population, generate_individuals(gene_number) )
	}

	// ---------------------- Keyboard ---------------------- //

	// Keyboard used by human user
	keyboard_human = make(map[Direction][]bool)
	keyboard_human[up]		= append(keyboard_human[up], false)
	keyboard_human[down]	= append(keyboard_human[down], false)
	keyboard_human[left]	= append(keyboard_human[left], false)
	keyboard_human[right]	= append(keyboard_human[right], false)


	// Keyboard used by automations
	keyboard_automations = make(map[Direction][]bool)
	for i := 0 ; i < len(population) ; i ++ {
		keyboard_automations[up]		= append(keyboard_automations[up], false)
		keyboard_automations[down]	= append(keyboard_automations[down], false)
		keyboard_automations[left]	= append(keyboard_automations[left], false)
		keyboard_automations[right]	= append(keyboard_automations[right], false)
	}


	// ---------------- Player and background --------------- //

	// Load the PixelMap Image
	spriteMap, err := loadPicture("Maze/spritemap-rpg.png")




	// // Initialize Player data
	if automation == false {	// Draw just player0
		player_list = append(player_list, &player{})
		player_list[0].restart_player(spriteMap, player_list[0])
	} else {	// draw all population
			// Add players accordingly to population
			for i := 0 ; i < population_size ; i++ {
				player_list = append(player_list, &player{})
				player_list[i].restart_player(spriteMap, player_list[i])
			}
	}

	// Initialize the background
	bgd := &background{}
	bgd.setPlayerSprites(spriteMap)





	// Infinite loop
	for !win.Closed() {

		// Draw all background objects first to this object and just draw to window one time later
		imd := imdraw.New(spriteMap)

		// Esc to quit program
    if win.JustPressed(pixelgl.KeyEscape) {
      break
    }

		// Clear Screen
		win.Clear(colornames.Lightgreen)

		// ---------------------- Keyboard ---------------------- //

		// Update player direction and keys pressed
		if win.JustPressed(pixelgl.KeyUp) {
			keyboard_human[up][0] = true
		}
		if win.JustPressed(pixelgl.KeyDown) {
			keyboard_human[down][0] = true
		}
		if win.JustPressed(pixelgl.KeyLeft) {
			keyboard_human[left][0] = true
		}
		if win.JustPressed(pixelgl.KeyRight) {
			keyboard_human[right][0] = true
		}


		// ---------- Read and execute commands from IA ---------- //
		if automation {

			// Decode all individuals into commands and save it to a Matrix
			if cycle == 0 {
				commands_matrix = individualtoCommands(population, gene_number)
			}

			// Loop for all commands available
			if cycle < len(commands_matrix[0]) {

				// Fill the commands in all virtual keyboards
				for i := 0 ; i < len(population) ; i ++ {
					// Execute the command on keyboard

					// UP[0] first player, UP[1] second player...
					keyboard_automations[commands_matrix[i][cycle]][i] = true
				}

				// Update cycle
				cycle ++

			// Finished all commands for this generation, reset and start again
			} else {

				// If there are more generations to run
				if current_generation < generations {

					// Update the Score slice
					for i := 0 ; i < population_size ; i ++ {
						population_score = append(population_score, player_list[i].score)
					}

					// Clean variables for the next generation
					cycle = 0
					// // Restart game for next individual
					for i := 0 ; i < population_size ; i++ {
						player_list[i].restart_player(spriteMap, player_list[i])
					}

					genetic_algorithm()
					current_generation ++
					max_generation_position = 0
				} else {
					fmt.Println("\nSimulation Ended\n\n")
					automation = false
				}
			}

		}

		// ---------------------- Keyboard ---------------------- //

		// Move Player - Necessary for the automation of player execution
		// Update player direction and keys pressed
		if win.JustPressed(pixelgl.KeyUp) {
			keyboard_human[up][0] = true
		}
		if win.JustPressed(pixelgl.KeyDown) {
			keyboard_human[down][0] = true
		}
		if win.JustPressed(pixelgl.KeyLeft) {
			keyboard_human[left][0] = true
		}
		if win.JustPressed(pixelgl.KeyRight) {
			keyboard_human[right][0] = true
		}


		// Move Player - Necessary for the automation of player execution
		if keyboard_human[up][0] == true {
			direction = up
			player_list[0].update(up, 0)
		}
		if keyboard_human[down][0] == true {
			direction = down
			player_list[0].update(down, 0)
		}
		if keyboard_human[left][0] == true {
			direction = left
			player_list[0].update(left, 0)
		}
		if keyboard_human[right][0] == true {
			direction = right
			player_list[0].update(right, 0)
		}

		// Virtual Keyboard for automation
		// Move Automated Players - Necessary for the automation of player execution
		for i := 0 ; i < len(population) ; i ++ {
			if keyboard_automations[up][i] == true {
				direction = up
				player_list[i].update(up, i)
			}
			if keyboard_automations[down][i] == true {
				direction = down
					player_list[i].update(down, i)
			}
			if keyboard_automations[left][i] == true {
				direction = left
				player_list[i].update(left, i)
			}
			if keyboard_automations[right][i] == true {
				direction = right
				player_list[i].update(right, i)
			}
		}

		// Clean key pressed for the next cycle
		keyboard_human[up][0] = false
		keyboard_human[down][0] = false
		keyboard_human[left][0] = false
		keyboard_human[right][0] = false

		for i := 0 ; i < len(population) ; i ++ {
			keyboard_automations[up][i] = false
			keyboard_automations[down][i] = false
			keyboard_automations[left][i] = false
			keyboard_automations[right][i] = false
		}

		// -------------------- Draw Objects -------------------- //

		// Draw the entire background
		bgd.draw(imd)

		// Draw Players on the screen
		for j := 0 ; j < len(player_list) ; j++ {
			player_list[j].draw(imd)
		}

		// Draw with just one draw() call to screen
		imd.Draw(win)

		// Update the screen
		win.Update()

		// time.Sleep(100 * time.Millisecond)
	}
}
