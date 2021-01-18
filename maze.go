package main

import (
	"github.com/faiface/pixel/pixelgl"
	"Maze_Game/Maze"
)

// Main function
func main() {

	// Start Window system
	pixelgl.Run(Maze.Run)
}
