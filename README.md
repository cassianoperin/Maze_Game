# Maze Game
Simple Maze game coded in Go with genetic algorithms option to search for the best solutions.

## Objective:
Cross the screen and reach the empty space on last column at the right.

**Human** | **Genetic Algorithms**
:-------------------------:|:-------------------------:
<img width="430" alt="horizontal" src="https://github.com/cassianoperin/Maze_Game/blob/main/Images/maze-human.png">  |  <img width="430" alt="vertical" src="https://github.com/cassianoperin/Maze_Game/blob/main/Images/maze-automations.gif">

## Usage
- Set variable "**automation**" on "game.go" file to enable or disable Genetic Algorithms.
- Set the genetic algorithms main options on file "**genetic_algorithms.go**":
  - Population size (population_size)
  - Number of genes (gene_number)
  - Number of participants of tournament for parents selection (k)
  - Crossover rate (crossover_rate)
  - Mutation rate (mutation_rate)
  - Number of generations (generations)
  - Elitism percentual (elitism_percentual)

## Missing:
- Improve score considering the individual that got the best result in less movements.

## Documentation:

- Pixel:

https://github.com/faiface/pixel/wiki

- Game programing:

https://www.codingdream.com/index.php/simple-pacman-in-using-go-and-pixelgl-part-1
