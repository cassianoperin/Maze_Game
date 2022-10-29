package Maze

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// ----------- Player ----------- //
type Direction int

type player struct {
	sprites                map[Direction][]pixel.Rect
	currentSprite          pixel.Rect
	spriteMap              pixel.Picture
	grid_pos_X             int
	grid_pos_Y             int
	score                  int
	max_ind_position       int
	max_ind_position_cycle int
}

// --------- Background --------- //
type background struct {
	sprites map[int][]pixel.Rect
}

type block struct {
	currentSprite pixel.Rect
	spriteMap     pixel.Picture
	gridX         int
	gridY         int
}

// ---------- Objective --------- //
type objective_reached struct {
	generation int
	individual string
	score      int
	steps      int
}

// ---------- Constants --------- //
const (
	// Window Size
	screen_height = 800
	screen_width  = 800

	// Directions
	up    Direction = 0
	down  Direction = 1
	left  Direction = 2
	right Direction = 3
)

// --------- Variables ---------- //
var (
	// IA
	Automation      bool = false
	commands_matrix [][]Direction

	// Slice of players
	player_list []*player
	direction   Direction

	// Background
	backgroundMap [][]uint8
	// Grid size is defined by : X = number of objects per line	Y = number of objects in the slice
	grid_size_x       int
	grid_size_y       int
	map_best_solution int

	// Score
	max_generation_position int = 0

	// Keyboard
	keyboard_human       map[Direction][]bool // Human Keyboard
	keyboard_automations map[Direction][]bool // Automation Virtual Keyboards

	// Cycle counter
	cycle int = 0

	// Objective slice
	objective []objective_reached

	// Simulation finish and show results
	simlation_finished bool = false

	// Fonts
	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// Screen messages
	textMessage *text.Text // On screen Message content
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
	// gridWidth := width / float64(grid_x_size)
	// gridHeight := height / float64(grid_y_size)
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
	sprite.Draw(win, pixel.IM.ScaledXY(pixel.ZV, pixel.V(pos.W()/sprite.Frame().W(), pos.H()/sprite.Frame().H())).Moved(pos.Center()))
}

// Update the grid position accordingly to the direction of the next frame
// Collision Detection
func (object *player) getNewGridPos(direction Direction) (int, int) {
	if direction == right {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if object.grid_pos_X+1 < grid_size_x && backgroundMap[len(backgroundMap)-1-object.grid_pos_Y][object.grid_pos_X+1] == 0 {
			object.grid_pos_X += 1
		}
		return object.grid_pos_X, object.grid_pos_Y
	}
	// backgroundMap[line][column]
	if direction == left {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if object.grid_pos_X-1 >= 0 && backgroundMap[len(backgroundMap)-1-object.grid_pos_Y][object.grid_pos_X-1] == 0 {
			object.grid_pos_X -= 1
		}
		return object.grid_pos_X, object.grid_pos_Y
	}
	if direction == up {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if object.grid_pos_Y+1 < grid_size_y && backgroundMap[len(backgroundMap)-1-(object.grid_pos_Y+1)][object.grid_pos_X] == 0 {
			object.grid_pos_Y += 1
		}
		return object.grid_pos_X, object.grid_pos_Y
	}
	if direction == down {
		// Keep the player inside the window && just update if there isn't an object on the next move position
		if object.grid_pos_Y-1 >= 0 && backgroundMap[len(backgroundMap)-1-(object.grid_pos_Y-1)][object.grid_pos_X] == 0 {
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
		max_generation_position = object.grid_pos_X
	}

	// Update the Individual Maximum score and check the objective!
	if object.grid_pos_X > object.max_ind_position {

		// Update the new maximum posistion of player
		object.max_ind_position = object.grid_pos_X

		// Add score points
		player_score(object.max_ind_position, cycle, player_index)

		// Objective reached!!
		if object.grid_pos_X == len(backgroundMap[0])-1 {

			objective = append(objective, objective_reached{generation: current_generation, individual: population[player_index], score: object.grid_pos_X, steps: cycle})

			// fmt.Printf("\n\n\n\t\tObjective accomplished!\n\t\tIndividual: %s\tPosition: %d\tMovements: %d\n\n\n", population[player_index], len(backgroundMap[0]) - 1, cycle)
		}
	}
	// Punishment
	// } else {
	// 	object.score -= 10
	// }

}

// Calculate the player's score
func player_score(max_pos int, cycle int, plr_index int) {
	var (
		tmp_score     float64
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
	for i := 0; i < len(backgroundMap); i++ { // Lines
		for j := 0; j < len(backgroundMap[0]); j++ { // Columns
			if backgroundMap[i][j] == 0 {
				// Don't draw anything, its the path
			} else if backgroundMap[i][j] == 1 {
				b := block{currentSprite: bgd.sprites[0][0], gridX: (len(backgroundMap) - 1) - i, gridY: j}
				b.draw(t)
			} else if backgroundMap[i][j] == 2 {
				b := block{currentSprite: bgd.sprites[1][0], gridX: (len(backgroundMap) - 1) - i, gridY: j}
				b.draw(t)
			} else if backgroundMap[i][j] == 3 {
				b := block{currentSprite: bgd.sprites[2][0], gridX: (len(backgroundMap) - 1) - i, gridY: j}
				b.draw(t)
			} else if backgroundMap[i][j] == 4 {
				b := block{currentSprite: bgd.sprites[3][0], gridX: (len(backgroundMap) - 1) - i, gridY: j}
				b.draw(t)
			} else if backgroundMap[i][j] == 5 {
				b := block{currentSprite: bgd.sprites[4][0], gridX: (len(backgroundMap) - 1) - i, gridY: j}
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
	direction = right // To identify the initial sprite
	object.currentSprite = object.sprites[direction][0]
	// Score
	object.score = 0
	object.max_ind_position = 0
	object.max_ind_position_cycle = 0
}

// Convert the binary string of individuals to commands
func individualtoCommands(pop []string, gene_nr int) [][]Direction {

	// -------------- Prepare the Multidimensional Slice -------------- //
	// Declaring a slice of slices with a length of POPULATION
	commands := make([][]Direction, len(pop))

	// looping through the slice to declare a slice of each slice size
	for i := 0; i < len(pop); i++ {

		// Length of each slice shoulg be a half of Gene_number (2 digits = 1 command)
		new_length := gene_nr / 2

		commands[i] = make([]Direction, new_length)
	}

	// ----------------------- Process the data ----------------------- //

	// Decode each individual into commands
	for i := 0; i < len(pop); i++ {

		individual_split := strings.Split(pop[i], "")

		// Read individual and transform it to a slice with commands
		index := 0
		for j := 0; j < len(individual_split)/2; j++ {
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

			index += 2
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
	validate_parameters(Population_size, K)

	// 0 - Generate the population
	// Generate each individual for population
	for i := 0; i < Population_size; i++ {
		population = append(population, generate_individuals(Gene_number))
	}

	// ---------------------- Keyboard ---------------------- //

	// Keyboard used by human user
	keyboard_human = make(map[Direction][]bool)
	keyboard_human[up] = append(keyboard_human[up], false)
	keyboard_human[down] = append(keyboard_human[down], false)
	keyboard_human[left] = append(keyboard_human[left], false)
	keyboard_human[right] = append(keyboard_human[right], false)

	// Keyboard used by automations
	keyboard_automations = make(map[Direction][]bool)
	for i := 0; i < len(population); i++ {
		keyboard_automations[up] = append(keyboard_automations[up], false)
		keyboard_automations[down] = append(keyboard_automations[down], false)
		keyboard_automations[left] = append(keyboard_automations[left], false)
		keyboard_automations[right] = append(keyboard_automations[right], false)
	}

	// ------------------- Define the map ------------------- //

	if Automation {
		if Maze_map == 0 {
			backgroundMap = backgroundMap_0_automate
			map_best_solution = backgroundMap_0_best_solution
		} else if Maze_map == 1 {
			backgroundMap = backgroundMap_1_automate
			map_best_solution = backgroundMap_1_best_solution
		} else if Maze_map == 2 {
			backgroundMap = backgroundMap_2_automate
			map_best_solution = backgroundMap_2_best_solution
		} else if Maze_map == 3 {
			backgroundMap = backgroundMap_3_automate
			map_best_solution = backgroundMap_3_best_solution
		} else {
			fmt.Printf("Map %d not found! Exiting.\n", Maze_map)
			os.Exit(2)
		}

	} else {
		if Maze_map == 0 {
			backgroundMap = backgroundMap_0
			map_best_solution = backgroundMap_0_best_solution
		} else if Maze_map == 1 {
			backgroundMap = backgroundMap_1
			map_best_solution = backgroundMap_1_best_solution
		} else if Maze_map == 2 {
			backgroundMap = backgroundMap_2
			map_best_solution = backgroundMap_2_best_solution
		} else if Maze_map == 3 {
			backgroundMap = backgroundMap_3
			map_best_solution = backgroundMap_3_best_solution
		} else {
			fmt.Printf("Map %d not found! Exiting.\n", Maze_map)
			os.Exit(2)
		}
	}

	// Calculate the size of the grid according to map selected
	grid_size_x = len(backgroundMap[0])
	grid_size_y = len(backgroundMap)

	// ---------------- Player and background --------------- //

	// Load the PixelMap Image
	spriteMap, err := loadPicture("Images/spritemap-rpg.png")

	// // Initialize Player data
	if Automation == false { // Draw just player0
		player_list = append(player_list, &player{})
		player_list[0].restart_player(spriteMap, player_list[0])
	} else { // draw all population
		// Add players accordingly to population
		for i := 0; i < Population_size; i++ {
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

		// Clear Screen
		win.Clear(colornames.Lightgreen)

		// Esc to quit program
		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}

		if !simlation_finished {

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
			if Automation {

				// Decode all individuals into commands and save it to a Matrix
				if cycle == 0 {
					commands_matrix = individualtoCommands(population, Gene_number)
				}

				// Loop for all commands available
				if cycle < len(commands_matrix[0]) {

					// Fill the commands in all virtual keyboards
					for i := 0; i < len(population); i++ {
						// Execute the command on keyboard

						// UP[0] first player, UP[1] second player...
						keyboard_automations[commands_matrix[i][cycle]][i] = true
					}

					// Update cycle
					cycle++

					// Finished all commands for this generation, reset and start again
				} else {

					// If there are more generations to run
					if current_generation < Generations {

						// Update the Score slice
						for i := 0; i < Population_size; i++ {
							population_score = append(population_score, player_list[i].score)
						}

						// Clean variables for the next generation
						cycle = 0
						// // Restart game for next individual
						for i := 0; i < Population_size; i++ {
							player_list[i].restart_player(spriteMap, player_list[i])
						}

						genetic_algorithm()
						current_generation++
						max_generation_position = 0
					} else {
						fmt.Printf("\n\n\n|| ---------------------------------- Simulation Ended ---------------------------------- ||\n\nWinners:\n")
						for i := 0; i < len(objective); i++ {
							fmt.Printf("%d\tGen: %d\tIndividual: %s\tScore: %d\tSteps: %d\n", i+1, objective[i].generation, objective[i].individual, objective[i].score, objective[i].steps)
						}

						// Calculate the best one (less steps)
						quickest := Gene_number / 2
						for i := 0; i < len(objective); i++ {
							if objective[i].steps < quickest {
								quickest = objective[i].steps
							}
						}

						fmt.Printf("\nBest performances:\n")
						for i := 0; i < len(objective); i++ {
							if objective[i].steps == quickest {
								fmt.Printf("Gen: %d\tIndividual: %s\tScore: %d\tSteps: %d\n", objective[i].generation, objective[i].individual, objective[i].score, objective[i].steps)
							}
						}
						fmt.Println()

						// Disable automation
						Automation = false

						// Show Results
						simlation_finished = true
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
			for i := 0; i < len(population); i++ {
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

			for i := 0; i < len(population); i++ {
				keyboard_automations[up][i] = false
				keyboard_automations[down][i] = false
				keyboard_automations[left][i] = false
				keyboard_automations[right][i] = false
			}

			// // ------------------- Draw Background ------------------ //
			// pic, err := loadPicture("Images/background.png")
			// if err != nil {
			// 	panic(err)
			// }

			// sprite := pixel.NewSprite(pic, pic.Bounds())

			// // Initial X and Y values
			// var X float64 = 110
			// var Y float64 = 70

			// for i := 0; i < 6; i++ {
			// 	for j := 0; j < 5; j++ {
			// 		// sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.1).Moved(pixel.V(110, 70)))
			// 		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.1).Moved(pixel.V(X, Y)))
			// 		X += 220
			// 	}
			// 	X = 110
			// 	Y += 140
			// }

			// -------------------- Draw Objects -------------------- //

			// Just draw Degug screen if Automation is enabled
			if Automation {
				imd.Color = colornames.Gray
				imd.Push(pixel.V(0, 630))
				imd.Push(pixel.V(800, 800))
				imd.Rectangle(0)

				imd.Color = colornames.Whitesmoke
				imd.Push(pixel.V(5, 635))
				imd.Push(pixel.V(795, 795))
				imd.Rectangle(0)
			}

			// Draw the entire background
			bgd.draw(imd)

			// Draw Players on the screen
			for j := 0; j < len(player_list); j++ {
				player_list[j].draw(imd)
			}

			// Draw with just one draw() call to screen
			imd.Draw(win)

			// Just draw Degug Text Information if Automation is enabled
			if Automation {
				// Generation
				textMessage = text.New(pixel.V(20, 780), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "GENERATION: %d of %d", print_current_generation+1, Generations)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

				// Mutated individuals
				textMessage = text.New(pixel.V(20, 760), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "Mutated individuals: %d", mutation_ind_count)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

				// Mutated genes
				textMessage = text.New(pixel.V(260, 760), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "Mutated genes: %d", mutation_count)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

				// Crossovers
				textMessage = text.New(pixel.V(20, 740), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "Crossovers: %d", print_crossover_count)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

				// Best Individual
				textMessage = text.New(pixel.V(20, 720), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "Best Individual: %s", print_best)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

				// Fitness Average
				textMessage = text.New(pixel.V(20, 700), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "Fitness Average: %d", print_average_score)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

				// Maximum position
				textMessage = text.New(pixel.V(20, 660), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "Maximum position: %d of %d", print_max_generation_position+1, grid_size_x)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

				// Fitness
				textMessage = text.New(pixel.V(260, 660), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "Fitness: %d", print_score)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

				// Number of Winners
				if len(objective) > 0 {
					textMessage = text.New(pixel.V(20, 640), atlas)
					textMessage.Clear()
					textMessage.Color = colornames.Black
					fmt.Fprintf(textMessage, "Number of Winners: %d", len(objective))
					textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))
				}
			}

		} else {

			imd.Color = colornames.Gray
			imd.Push(pixel.V(0, 630))
			imd.Push(pixel.V(800, 800))
			imd.Rectangle(0)

			imd.Color = colornames.Whitesmoke
			imd.Push(pixel.V(5, 635))
			imd.Push(pixel.V(795, 795))
			imd.Rectangle(0)

			// Draw the entire background
			bgd.draw(imd)

			// // Draw Players on the screen
			// for j := 0; j < len(player_list); j++ {
			// 	player_list[j].draw(imd)
			// }

			// Draw with just one draw() call to screen
			imd.Draw(win)

			// Banner
			textMessage = text.New(pixel.V(20, 780), atlas)
			textMessage.Clear()
			textMessage.Color = colornames.Black
			fmt.Fprintf(textMessage, "|| ---------------------- Simulation Ended ---------------------- ||")
			textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

			// Generation
			textMessage = text.New(pixel.V(20, 760), atlas)
			textMessage.Clear()
			textMessage.Color = colornames.Black
			fmt.Fprintf(textMessage, "|| GENERATIONS: %d", print_current_generation+1)
			textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

			// Number of Winners
			textMessage = text.New(pixel.V(20, 740), atlas)
			textMessage.Clear()
			textMessage.Color = colornames.Black
			fmt.Fprintf(textMessage, "|| Number of Winners: %d", len(objective))
			textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

			// Maximum position
			textMessage = text.New(pixel.V(260, 740), atlas)
			textMessage.Clear()
			textMessage.Color = colornames.Black
			fmt.Fprintf(textMessage, "Maximum position reached: %d of %d", best_step, grid_size_x)
			textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

			if len(objective) > 0 {

				// Calculate the best individual (less steps)
				quickest := Gene_number / 2
				for i := 0; i < len(objective); i++ {
					if objective[i].steps < quickest {
						quickest = objective[i].steps
					}
				}

				best_performer := 0 // less steps and the first one to reach the objective
				for i := len(objective) - 1; i >= 0; i-- {
					if objective[i].steps == quickest {
						// fmt.Printf("Gen: %d\tIndividual: %s\tScore: %d\tSteps: %d\n", objective[i].generation, objective[i].individual, objective[i].score, objective[i].steps)
						best_performer = i
					}
				}

				// Best Individual
				textMessage = text.New(pixel.V(20, 680), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "Best Individual: %s\n\nGeneration: %d, with %d steps (Best solution: %d)", objective[best_performer].individual, objective[best_performer].generation, objective[best_performer].steps, map_best_solution)
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))

			} else {
				// Best Individual
				textMessage = text.New(pixel.V(20, 700), atlas)
				textMessage.Clear()
				textMessage.Color = colornames.Black
				fmt.Fprintf(textMessage, "No Winner")
				textMessage.Draw(win, pixel.IM.Scaled(textMessage.Orig, 1))
			}

		}

		// Update the screen
		win.Update()

	}
}
