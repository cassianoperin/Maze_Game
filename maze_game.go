package main

import (
	"Maze_Game/Maze"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/faiface/pixel/pixelgl"
	"gopkg.in/ini.v1"
)

var (
	// Configuration file (ini)
	maze_ini string = ""
)

// Main function
func main() {

	// Load INI Variables
	load_INI()

	// Start Window system
	pixelgl.Run(Maze.Run)
}

func load_INI() {

	// Check the Operational System to save the ini file
	myos := runtime.GOOS

	// Get home dir
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error retriving home dir: %s. Exiting.\n", err)
		os.Exit(2)
	}

	if myos == "darwin" {
		maze_ini = home + "/.maze.ini"

		// Check if the ini file already exist
		if _, err := os.Stat(maze_ini); err != nil || os.IsNotExist(err) {
			// File not found, create
			f, err := os.Create(maze_ini)
			if err != nil {
				fmt.Printf("Error creating ini file: %s. Exiting.", err)
				os.Exit(2)
			}
			defer f.Close()

			// Write initial INI Values
			_, err2 := f.WriteString("[Maps]\nmap=1\t\t\t; 0 || 1\n\n[Mode]\nAutomation=true\t\t; true || false\n\n[Settings]\nGenerations=100\nPopulation_size=100\nGene_number=50\nK=25\nCrossover_rate=0.7\nMutation_rate=0.05\nElitism_percentual=10\n")
			if err2 != nil {
				fmt.Printf("Error writing to ini file: %s. Exiting.", err)
				os.Exit(2)
			}
		}

	} else if myos == "windows" {
		// Windows
		maze_ini = home + "\\.maze.ini"

		// Check if the ini file already exist
		if _, err := os.Stat(maze_ini); err != nil || os.IsNotExist(err) {
			// File not found, create
			f, err := os.Create(maze_ini)
			if err != nil {
				fmt.Printf("Error creating ini file: %s. Exiting.", err)
				os.Exit(2)
			}
			defer f.Close()

			// Write initial INI Values
			_, err2 := f.WriteString("[Maps]\nmap=1\t\t\t; 0 || 1\n\n[Mode]\nAutomation=true\t\t; true || false\n\n[Settings]\nGenerations=100\nPopulation_size=100\nGene_number=50\nK=25\nCrossover_rate=0.7\nMutation_rate=0.05\nElitism_percentual=10\n")
			if err2 != nil {
				fmt.Printf("Error writing to ini file: %s. Exiting.", err)
				os.Exit(2)
			}
		}

	} else if myos == "linux" {
		// Linux
		maze_ini = home + "/.maze.ini"

		// Check if the ini file already exist
		if _, err := os.Stat(maze_ini); err != nil || os.IsNotExist(err) {
			// File not found, create
			f, err := os.Create(maze_ini)
			if err != nil {
				fmt.Printf("Error creating ini file: %s. Exiting.", err)
				os.Exit(2)
			}
			defer f.Close()

			// Write initial INI Values
			_, err2 := f.WriteString("[Maps]\nmap=1\t\t\t; 0 || 1\n\n[Mode]\nAutomation=true\t\t; true || false\n\n[Settings]\nGenerations=100\nPopulation_size=100\nGene_number=50\nK=25\nCrossover_rate=0.7\nMutation_rate=0.05\nElitism_percentual=10\n")
			if err2 != nil {
				fmt.Printf("Error writing to ini file: %s. Exiting.", err)
				os.Exit(2)
			}
		}
	} else {
		fmt.Printf("Operational system not supported: %s. Exiting\n\n", myos)
		os.Exit(2)
	}

	// Load INI information:
	cfg_ini, err := ini.Load(maze_ini)
	if err != nil {
		fmt.Printf("Fail to read file: %s", err)
		os.Exit(1)
	}

	// ------------ Read ini options into program variables ------------ //

	// [Maps]
	tmp_value, err := strconv.ParseInt(cfg_ini.Section("Maps").Key("map").String(), 0, 32)
	Maze.Maze_map = int(tmp_value)
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'map': %s", err)
		os.Exit(2)
	}

	// [Mode] - Automation
	Maze.Automation, err = strconv.ParseBool(cfg_ini.Section("Mode").Key("Automation").String())
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'automation': %s", err)
		os.Exit(2)
	}

	// [Settings] - Population_size
	tmp_value, err = strconv.ParseInt(cfg_ini.Section("Settings").Key("Population_size").String(), 0, 32)
	Maze.Population_size = int(tmp_value)
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'population_size': %s", err)
		os.Exit(2)
	}

	// [Settings] - Gene_number
	tmp_value, err = strconv.ParseInt(cfg_ini.Section("Settings").Key("Gene_number").String(), 0, 32)
	Maze.Gene_number = int(tmp_value)
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'Gene_number': %s", err)
		os.Exit(2)
	}

	// [Settings] - K
	tmp_value, err = strconv.ParseInt(cfg_ini.Section("Settings").Key("K").String(), 0, 8)
	Maze.K = int(tmp_value)
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'K': %s", err)
		os.Exit(2)
	}

	// [Settings] - Crossover_rate
	Maze.Crossover_rate, err = strconv.ParseFloat(cfg_ini.Section("Settings").Key("Crossover_rate").String(), 0)
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'Crossover_rate': %s", err)

	}

	// [Settings] - Mutation_rate
	Maze.Mutation_rate, err = strconv.ParseFloat(cfg_ini.Section("Settings").Key("Mutation_rate").String(), 0)
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'Mutation_rate': %s", err)

	}

	// [Settings] - Generations
	tmp_value, err = strconv.ParseInt(cfg_ini.Section("Settings").Key("Generations").String(), 0, 32)
	Maze.Generations = int(tmp_value)
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'Generations': %s", err)
		os.Exit(2)
	}

	// [Settings] - Elitism_percentual
	tmp_value, err = strconv.ParseInt(cfg_ini.Section("Settings").Key("Elitism_percentual").String(), 0, 32)
	Maze.Elitism_percentual = int(tmp_value)
	if err != nil {
		fmt.Printf("Fail to read ini attribute 'Elitism_percentual': %s", err)
		os.Exit(2)
	}

}
