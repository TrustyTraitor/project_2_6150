package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

// RandomPlace -
// Partially random. Places pieces in random row, but ensures there is not more than 1 in any column.
func RandomPlace(pieces []Piece) ([]Piece, error) {

	size := len(pieces)

	x := rand.Intn(size)

	for i := range pieces {
		y := rand.Intn(size)

		pieces[i].x = x
		pieces[i].y = y

		x++
		if x >= len(pieces) {
			x = 0
		}
	}

	return pieces, nil
}

// SetManualPieces -
// Allows the user to manually enter a board instead of randomly generating it.
//
// Does not prevent invalid boards from being entered.
func SetManualPieces(pieces *[]Piece) []Piece {
	var x int
	var y int

	fmt.Println("Manually enter pieces by typing the x and y coordinates.")
	fmt.Println("The coordinates are space separated")
	fmt.Println("Each piece must be on a separate line. Hit enter after entering each coordinate.")
	fmt.Print("Example for a Queen that is placed at 3,4 (0 indexed) then another at 5,5:\n\n")
	fmt.Print("3 4\n5 5\n\n")
	fmt.Println("There is limited error checking. So ensure your entries are within [0, n)")
	fmt.Println()

	for i := range len(*pieces) {
		_, err := fmt.Scanf("%d %d\n", &x, &y)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		(*pieces)[i] = Piece{x, y}
	}

	return *pieces
}

// main -
// The main function handles program execution
//
// 1. Reads user input/command line args
//
// 2. Executes HillClimb and collects statistics
//
// 3. Prints stats to stdout.
func main() {
	// Get board dimensions
	var boardSize int
	if len(os.Args) == 1 {
		fmt.Print("What is the dimensions of the board? : ")

		_, err := fmt.Scanln(&boardSize)
		if err != nil {
			fmt.Println("Error reading input")
			return
		}
	} else {
		boardSize, _ = strconv.Atoi(os.Args[1])

	}

	if boardSize < 4 {
		fmt.Println("Board size must be at least 4")
		return
	}

	// Determine if the board needs to be randomly generated
	var useManual bool
	if len(os.Args) >= 3 {
		useManual, _ = strconv.ParseBool(os.Args[2])
	} else {
		fmt.Print("Manually enter pieces? : ")
		_, _ = fmt.Scanf("%t\n", &useManual)
	}

	// Determine how many times the program must be run
	var rerunCount int
	if len(os.Args) >= 4 {
		rerunCount, _ = strconv.Atoi(os.Args[3])
	} else {
		fmt.Print("How many runs? : ")
		_, _ = fmt.Scanf("%d\n", &rerunCount)
	}

	var randomRestartDef bool
	if len(os.Args) >= 5 {
		randomRestartDef, _ = strconv.ParseBool(os.Args[4])
	} else {
		fmt.Print("Use Random Restarts?? : ")
		_, _ = fmt.Scanf("%t\n", &randomRestartDef)
	}

	var useSidewaysMove bool
	if len(os.Args) >= 6 {
		useSidewaysMove, _ = strconv.ParseBool(os.Args[5])
	} else {
		fmt.Print("Use Sideways move? : ")
		_, _ = fmt.Scanf("%t\n", &useSidewaysMove)
	}

	pieces := make([]Piece, boardSize)

	// Random restart statistics and control
	randomRestart := randomRestartDef
	randRestartCount := 0

	// Controls some program functionality
	printPath := true
	stepLimit := 1000

	// Used for keeping track of statistics
	succSteps := 0
	failSteps := 0

	succPuzzle := 0
	failPuzzle := 0

	var decisionFunc Decision
	if useSidewaysMove {
		decisionFunc = SideWays
	} else {
		decisionFunc = NoSideWays
	}

	// A nested function definition. Created for clarity of code. However little it helps...
	RunHillClimb := func(idx int) {
		var steps int
		var err error

		totalSteps := 0

		for {
			if !useManual {
				pieces, _ = RandomPlace(pieces)
			} else {
				pieces = SetManualPieces(&pieces)
			}

			if idx > 3 {
				printPath = false
			} else if printPath {
				fmt.Printf("### Run: %v\n", idx)
			}

			_, _, steps, err = HillClimb(pieces, decisionFunc, stepLimit, printPath)

			totalSteps += steps

			if err == nil {
				randomRestart = false
			}

			if !randomRestart {
				break
			} else {
				randRestartCount++
			}
		}
		randomRestart = randomRestartDef

		if err != nil {
			failSteps += steps
			failPuzzle++

			if idx < 4 && printPath {
				fmt.Println(err.Error())
				fmt.Println()
			}
		} else {
			succPuzzle++

			if randomRestartDef {
				succSteps += totalSteps
			} else {
				succSteps += steps
			}

			if idx < 4 && printPath {
				fmt.Print("Solution Found\n\n")
			}
		}
	}

	// Runs the HillClimb algorithm some number of times based on rerunCount
	for idx := range rerunCount {
		RunHillClimb(idx)
	}

	// Calculates and displays statistics
	avgSuccSteps := float64(succSteps) / float64(succPuzzle)
	avgFailSteps := float64(failSteps) / float64(failPuzzle)

	fmt.Printf("Runs: %d   Sideways Move: %t   Random Restarts: %t   Step Limit: %d\n\n",
		rerunCount, useSidewaysMove, randomRestartDef, stepLimit)

	fmt.Printf("Succeeded Steps Avg: %.2f  Failed Steps Avg: %v\n",
		avgSuccSteps, avgFailSteps)

	totalPuzzles := succPuzzle + failPuzzle
	succPuzzlePercent := float64(succPuzzle) / float64(totalPuzzles)
	fmt.Printf(
		"Succeeded Puzzles: %v  Failed Puzzles: %v  Total Puzzles: %v  Percent: %.2f\n",
		succPuzzle, failPuzzle, totalPuzzles, succPuzzlePercent*100)

	if randomRestartDef {
		fmt.Printf("Random Restart Average: %.2f\n", float64(randRestartCount)/float64(totalPuzzles))
	}
}
