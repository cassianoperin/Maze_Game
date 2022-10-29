# Maze Game
Simple Maze game coded in Go with genetic algorithms option to search for the best solutions.

## Objective:
Cross the screen and reach the empty space on last column at the right.

**Human** | **Genetic Algorithms**
:-------------------------:|:-------------------------:
<img width="430" alt="horizontal" src="https://github.com/cassianoperin/Maze_Game/blob/main/Images/maze-human.png">  |  <img width="430" alt="vertical" src="https://github.com/cassianoperin/Maze_Game/blob/main/Images/maze-automations.gif">

## Compile

### MAC
`export GO111MODULE=off`
`go get github.com/faiface/pixel`
`go get github.com/faiface/glhf`
`go get github.com/go-gl/glfw/v3.3/glfw`
`go get gopkg.in/ini.v1`
`go build`

#### Instructions to pack into Mac executable:
1) Baixar o binário
https://github.com/kindlychung/genicon

2) Install dependency:
brew install imagemagick

3) Create the icon based on a PNG image
`./genicon icon.png tmp_folder`

4) Rename the folder mv AppIcon.appiconset AppIcon.iconset

5) Create the icon in the format needed by Mac executable
iconutil -c icns -o icon.icns AppIcon.iconset


### WINDOWS
TO DO

## Usage
1)  After the first execution, the program will create an ini file named '.maze.ini' into user home folder
  - To execute the game, set the value 'Automation' to false, otherwise, it will start in simulation mode
  - Select the map from 0 to 3
2) Define the genetic altorithm configuration:
  - Number of generations (Generations)
  - Population size (Population_size)
  - Number of genes (Gene_number)
  - Number of participants of tournament for parents selection (K)
  - Crossover rate (Crossover_rate)
  - Mutation rate (Mutation_rate)
  - Elitism percentual (Elitism_percentual)
3) Run the program

## Next steps:
- Improve score considering the individual that got the best result in less movements.
- After finish, show the path of winner
- Key to reset
- Binary for Windows
- Clean code
- Show the time spent after the execution
- Translate the individual into arrows
- Reactivate the background drawing (game.go) it efficiently
- Put the debug into the down side of screen


## Documentation:

- Pixel:

https://github.com/faiface/pixel/wiki

- Game programing:

https://www.codingdream.com/index.php/simple-pacman-in-using-go-and-pixelgl-part-1

- Mac APPs

https://medium.com/@mattholt/packaging-a-go-application-for-macos-f7084b00f6b5
