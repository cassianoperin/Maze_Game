package Maze

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	// --- Main variables --- //

	// Program Variables filled with INI information
	Population_size    int     // Default value = 100
	Gene_number        int     // Default value = 50
	K                  int     // Tournament size (number of participants) // Default value = 25
	Crossover_rate     float64 // Default value = 0.7
	Mutation_rate      float64 // I'm analyzing each gene so the mutation rate should be really small // Default value = 0.05
	Generations        int     // Default value = 100
	Elitism_percentual int     // Default value = 10 (10% of population size)

	// Other variables
	population          []string
	population_score    []int
	elitism_individuals int = (Elitism_percentual * Population_size) / 100

	// Counters
	mutation_count, mutation_ind_count int
	current_generation                 int = 0
	best_step                          int = 0

	// Print into screen variables
	print_best                    = ""
	print_score                   = 0
	print_current_generation      = 0
	print_crossover_count         = 0
	print_max_generation_position = 0
	print_average_score           = 0

	// Debug
	debug bool = false
)

func slice_average(slice []int, total int) (int, int) {
	var (
		sum        int = 0
		average    int = 0
		percentage int = 0
	)

	// Sum all values
	for i := 0; i < len(slice); i++ {
		sum += slice[i]
	}

	// Divide by the size of the slice
	average = sum / len(slice)
	// Percentage
	percentage = (average * 100) / total

	return average, percentage
}

// ------------------- Validate Parameters -------------------- //
func validate_parameters(pop_size int, competitors int) {
	// Minimal Population Size size accepted is 2
	if pop_size%2 == 1 {
		fmt.Printf("\nPopulation size should be ODD numbers. Exiting\n")
		os.Exit(0)
	}

	// Population Size should be positive
	if pop_size <= 0 {
		fmt.Printf("\nPopulation size should be Positive. Exiting\n")
		os.Exit(0)
	}

	// K (competitors) must be at least 2
	if competitors < 2 {
		fmt.Printf("\nNumber of competitors (k) must be at least 2. Exiting\n")
		os.Exit(0)
	}
}

// ------------------- Generate Individuals ------------------- //
func generate_individuals(gene_nr int) string {
	var individual string = ""

	// Initialize rand source
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < gene_nr; i++ {
		individual += strconv.Itoa(rand.Intn(2))
	}

	return individual
}

// --------- Generate the Evaluation of an Individual --------- //
func fitness_individual(individual string) int {
	ones := regexp.MustCompile("1")
	matches := ones.FindAllStringIndex(individual, -1)

	return len(matches)
}

// --------- Generate the Evaluation of a Population ---------- //
func fitness_population(pop []string) []int {

	var score []int

	for i := 0; i < len(pop); i++ {
		score = append(score, fitness_individual(pop[i]))
	}

	return score
}

// ------------------------- Elitism -------------------------- //
func elitism(pop []string, pop_score []int, pop_size int, elitism_number int) ([]string, []string) {
	var (
		elite, elite_score, tmp_slice []string
	)

	// Append score + individual in one slice
	for i := 0; i < pop_size; i++ {
		tmp_slice = append(tmp_slice, strconv.Itoa(pop_score[i])+","+pop[i])
	}

	// Sort slice
	sort.Strings(tmp_slice)

	// Insert individuals on Elite slice and score on elite_score
	for i := pop_size - 1; i > (pop_size-1)-elitism_number; i-- {
		tmp_slice := strings.Split(tmp_slice[i], ",")
		elite = append(elite, tmp_slice[1])             // Individual
		elite_score = append(elite_score, tmp_slice[0]) //Score
	}

	return elite, elite_score
}

// ---------------------- Define Parents ---------------------- //
func define_parents(pop []string, pop_size int, k int) []string {
	var parents []string

	// Quantity of tournaments is equal to the size of population
	for tournament := 0; tournament < pop_size; tournament++ {
		var (
			competitors []string
			score       []int
		)

		// Each tournament, K competitors
		for i := 0; i < k; i++ {
			competitors = append(competitors, pop[rand.Intn(pop_size)])
		}

		// Calculate the score of K competitors
		for i := 0; i < k; i++ {
			score = append(score, fitness_individual(competitors[i]))
		}

		bigger := score[0]
		winner := competitors[0]

		for i := 0; i < k; i++ {
			if score[i] > bigger {
				bigger = score[i]
				winner = competitors[i]
			}
		}

		parents = append(parents, winner)

		if debug {
			fmt.Printf("\tTournament: %d\t Competitors: %s\t Scores: %d\t Winner: %s (%d)\n", tournament, competitors, score, winner, bigger)
		}

	}

	return parents

}

// -------------------- Generate Children --------------------- //
func generate_children(parents []string, pop_size int, elitism_number int, elite []string) ([]string, int) {
	var (
		father1, father2, child1, child2 string
		pop_new                          []string
		cross_count                      int = 0
	)

	if debug {
		fmt.Printf("\n\tSelected parents:\n")
	}

	for i := 0; i < pop_size/2; i++ {
		// Define the couples
		randomIndex := rand.Intn(len(parents))
		father1 = parents[randomIndex]

		randomIndex = rand.Intn(len(parents))
		father2 = parents[randomIndex]

		if debug {
			fmt.Printf("\t%d) %s with %s\n", i, father1, father2)
		}

		// Define if will have crossover (the parents will be copied to next generation)
		if rand.Float64() < Crossover_rate {

			// Define the cut-point
			cut_point := rand.Intn(Gene_number-1) + 1
			if debug {
				fmt.Printf("\t\tCut-point: %d\n", cut_point)
			}

			// Split father's values
			// Father1
			father1_split := strings.Split(father1, "")
			father1_split_p1 := father1_split[0:cut_point]
			father1_split_p2 := father1_split[cut_point:]
			// Father2
			father2_split := strings.Split(father2, "")
			father2_split_p1 := father2_split[0:cut_point]
			father2_split_p2 := father2_split[cut_point:]

			// Child1
			child1_p1 := strings.Join(father1_split_p1, "")
			child1_p2 := strings.Join(father2_split_p2, "")
			child1 = child1_p1 + child1_p2
			if debug {
				fmt.Printf("\t\tChild1: %s + %s: %s\n", child1_p1, child1_p2, child1)
			}

			// Child2
			child2_p1 := strings.Join(father2_split_p1, "")
			child2_p2 := strings.Join(father1_split_p2, "")
			child2 = child2_p1 + child2_p2
			if debug {
				fmt.Printf("\t\tChild2: %s + %s: %s\n", child2_p1, child2_p2, child2)
			}

			// Put the childs in the new generation
			pop_new = append(pop_new, child1)
			pop_new = append(pop_new, child2)

		} else {
			if debug {
				fmt.Printf("\t\tCrossover:\n")
			}
			pop_new = append(pop_new, father1)
			pop_new = append(pop_new, father2)
			if debug {
				fmt.Printf("\t\tChild1 (Father1): %s\n", father1)
				fmt.Printf("\t\tChild2 (Father2): %s\n", father2)
			}
			cross_count++
		}

	}

	// Ensure place of elite members on next generation
	if elitism_number > 0 {
		if debug {
			fmt.Printf("\n\tElitism: Regular individual removal:\n")
		}

		// Remove randomically the number os elite elements
		for i := 0; i < elitism_number; i++ {
			random := rand.Intn(len(pop_new))
			if debug {
				fmt.Printf("\t\tIndividual %d:\t%s removed randomically from new population\n", i, pop_new[random])
			}

			// Remove the element at index 'random' from pop_new
			pop_new[random] = pop_new[len(pop_new)-1] // Copy last element to index 'random'.
			pop_new[len(pop_new)-1] = ""              // Erase last element (write zero value).
			pop_new = pop_new[:len(pop_new)-1]        // Truncate slice.
		}

		// Insert Elite Members on next generation
		if debug {
			fmt.Printf("\n\tElitism: Elite individual insertion:\n")
		}
		for i := 0; i < elitism_number; i++ {
			pop_new = append(pop_new, elite[i])
			if debug {
				fmt.Printf("\t\tIndividual %d\t%s inserted to new population\n", i, elite[i])
			}
		}
	}

	return pop_new, cross_count
}

// ------------------------- Mutation ------------------------- //
func generate_mutation(new_pop []string, pop_size int, gene_nr int, Mutation_rate float64) ([]string, int, int) {

	var (
		new_pop_mutated   []string
		count_genes       int = 0
		count_individuals int = 0
	)

	// For all individuals in population
	for i := 0; i < pop_size; i++ {

		var (
			individual              string = ""
			individual_mutated_flag bool
		)

		individual = new_pop[i]

		// For each gene, check for mutations
		for gene := 0; gene < gene_nr; gene++ {

			// Check if there is a mutation
			if Mutation_rate >= rand.Float64() {

				individual_split := strings.Split(individual, "")

				// Invert the mutated gene
				if individual_split[gene] == "0" {
					individual_split[gene] = "1"

				} else {
					individual_split[gene] = "0"
				}

				// Update the mutated individual
				individual = strings.Join(individual_split, "")

				if debug {
					fmt.Printf("\tIndividual #%d (%s) mutated on gene %d. New Individual: %s \n", i, new_pop[i], gene, individual)
				}

				count_genes++ // Generation genes mutated count
				individual_mutated_flag = true

			}

		}

		// Generation individuals mutated count
		if individual_mutated_flag {
			count_individuals++
			individual_mutated_flag = false
		}

		// Add mutated individuals to a new generation
		new_pop_mutated = append(new_pop_mutated, individual)
	}

	return new_pop_mutated, count_genes, count_individuals
}

// --------------------- Best Individual ---------------------- //
func best_individual() (string, int) {
	// var score []int
	//
	// // Calculate the score of the latest population
	// score = fitness_population(population)

	bigger := population_score[0]
	winner := population[0]

	for i := 0; i < len(population_score); i++ {
		if population_score[i] > bigger {
			bigger = population_score[i]
			winner = population[i]
		}
	}

	return winner, bigger
}

// ------------------------- MAIN FUNCTION ------------------------- //
func genetic_algorithm() {

	// // --------------------- Validate parameters --------------------- //
	// validate_parameters(Population_size, k)
	//
	//
	// // ----------------- 0 - Generate the population ----------------- //
	// // Generate each individual for population
	// for i := 0 ; i < Population_size ; i++ {
	//   population = append( population, generate_individuals(Gene_number) )
	// }

	// ----------------------- Generations Loop ---------------------- //
	// for i := 0 ; i < Generations ; i ++ {

	if debug {
		fmt.Printf("\n// ---------------------------------- GENERATION: %d ---------------------------------- //\n\n", current_generation)
	}

	// ----------------------- 1 - Evaluation ------------------------ //
	if debug {
		fmt.Printf("1 - Evaluation:\n\n")
	}
	// population_score = fitness_population(population)
	// Evaluation

	// Show the evaluation of each individual
	if debug {
		for i := 0; i < Population_size; i++ {
			fmt.Printf("\tIndividual %d:\t%s\tEvaluation %d\n", i, population[i], population_score[i])
		}
	}

	// ---------------------- 2 - Define Parents --------------------- //
	if debug {
		fmt.Printf("\n2 - Define Parents:\n\n")
	}

	parents := define_parents(population, Population_size, K)

	if debug {
		fmt.Printf("\n\tParents: %s\n\n", parents)
	}

	// ------------------------- 3 - Elitism ------------------------- //
	elite, elite_score := elitism(population, population_score, Population_size, elitism_individuals)
	if debug {
		fmt.Printf("\n3 - Elitism:\n\n\tNumber of elite members: %d\n\n", elitism_individuals)
		for i := 0; i < elitism_individuals; i++ {
			fmt.Printf("\tIndividual %d:\t%s set for elite with score: %s\n", i, elite[i], elite_score[i])
		}
	}

	// -------------------- 4 - Generate Children -------------------- //
	new_population, crossover_count := generate_children(parents, Population_size, elitism_individuals, elite)
	if debug {
		fmt.Printf("\n4 - Generate Chindren:\n\n\tNew population: %s\n", new_population)
	}

	// ------------------------ 5 - Mutation ------------------------- //
	new_population, mutation_count, mutation_ind_count = generate_mutation(new_population, Population_size, Gene_number, Mutation_rate)
	if debug {
		fmt.Printf("\n5 - Mutation:\n\tMutated Generation: %s\n\n", new_population)
	}

	// ---- 6 - Replace population vector with new population one ---- //
	population = nil // Clean ond population
	for i := 0; i < len(new_population); i++ {
		population = append(population, new_population[i])
	}

	average_score := 0
	for i := 0; i < len(population_score); i++ {
		average_score += population_score[i]
	}

	average_score = average_score / len(population_score)

	// -------------------- 7 - Best individual ---------------------- //

	// Print debug to console
	best, score := best_individual()
	fmt.Printf("\nGENERATION: %d\n", current_generation)
	fmt.Printf("Mutated individuals: %d\t\tMutated Genes: %d\n", mutation_ind_count, mutation_count)
	fmt.Printf("Crossovers: %d\n", crossover_count)
	fmt.Printf("Best Individual: %s\n", best)
	fmt.Printf("Fitness Average: %d\n\n", average_score)
	fmt.Printf("Maximum position: %d\tFitness: %d\n\n", max_generation_position+1, score)

	// Keep the max number of steps reached
	if max_generation_position+1 > best_step {
		best_step = max_generation_position + 1
	}

	// Now set the variables to be printed on screen
	print_best, print_score = best_individual()
	print_current_generation = current_generation
	print_crossover_count = crossover_count
	print_max_generation_position = max_generation_position
	print_average_score = average_score

	// Restart Variables
	population_score = nil

	// }

}
