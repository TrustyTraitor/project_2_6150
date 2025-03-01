package main

import (
	"errors"
	"math/rand"
)

// Decision -
// A generic function prototype used in the HillClimb function.
//
// Used when determining when the algorithm should exit.
//
// This is used to implement Sideways and non-sideways movement.
type Decision = func(...int) bool

// NoSolution -
// A custom error type returned if no solution can be found.
type NoSolution struct {
	msg string
}

// Error -
// Implementation of the error interface.
func (n NoSolution) Error() string {
	return n.msg
}

// findPiece - A helper function for finding a piece based on it's X coordinate
//
// Assumes that only a single piece can be on a particular X coordinate
func findPiece(pieces []Piece, x int) (int, error) {
	for idx, piece := range pieces {
		if piece.x == x {
			return idx, nil
		}
	}
	return 0, errors.New("no piece found")
}

// HillClimb -
// The HillClimb algorithm implementation for the 8 Queens problem.
//
// Allows for any size board to be input and will attempt to find the solution.
//
// Output -
//
// Board - The final state of the board. May or may not be solved.
//
// []Piece - A list of pieces, each containing their final locations. May or may not be solved.
//
// int - The number of steps it took to complete.
//
// error - A NoSolution error or Nil, depends on if a solution is found
func HillClimb(pieces []Piece, decision Decision, stepLimit int, showPath bool) (Board, []Piece, int, error) {
	count := 0
	board := NewBoard(len(pieces))

	if showPath {
		PrintPieces(pieces)
	}

	for {
		count++
		// SECTION: Computing the current H, stored in the hToBeat variable
		//          This also stores the current H in the board
		hToBeat := 0

		for i, p := range pieces {
			v := p.QueenHeuristic(pieces[i+1:])
			hToBeat += v
		}
		for _, p := range pieces {
			board.Set(p.x, p.y, hToBeat)
		}
		// SECTION: END

		// SECTION: Calculate all the possible moves and their H values
		//          Stores the calculated H in the board for referencing
		for x, piece := range pieces {
			for y := range len(pieces) {
				h := 0

				if y != piece.y { // This simply prevents a bit of recalculations
					tempPiece := piece
					tempPiece.y = y

					pieces[x] = tempPiece

					for i, p := range pieces {
						h += p.QueenHeuristic(pieces[i+1:])
					}
					board.Set(piece.x, y, h)

					pieces[x] = piece
				}
			}
		}
		// SECTION: END

		// SECTION: Finding the minimum move
		FindMinMoves := func(board *Board, h int) ([]Piece, int) {
			var minMoves []Piece

			minVal := h

			for x := range len(*board) {
				for y := range *board {
					v := board.Get(x, y)

					if v == minVal {
						minMoves = append(minMoves, Piece{x, y})
					} else if v < minVal {
						minMoves = nil

						minMoves = append(minMoves, Piece{x, y})
						minVal = v
					}
				}
			}

			return minMoves, minVal
		}
		minMoves, newH := FindMinMoves(&board, hToBeat)
		//fmt.Printf("Options: %v\nNewH: %v\n", minMoves, newH)

		randIdx := rand.Intn(len(minMoves)) // Picks random index from the minMoves slice
		chosenPiece := minMoves[randIdx]

		replaceIdx, _ := findPiece(pieces, chosenPiece.x) // Finds the index of the piece in []Pieces that matches the x of the new piece
		pieces[replaceIdx] = chosenPiece                  // Replaces old piece with new move
		// SECTION: END

		if showPath {
			PrintPieces(pieces)
		}

		// SECTION: Deciding if the algorithm should continue
		if !decision(hToBeat, newH, count, stepLimit) {
			if newH == 0 {
				return board, pieces, count, nil
			} else {
				return board, pieces, count, &NoSolution{"No Solution Found"}
			}
		}
		hToBeat = newH
		// SECTION: END
	}
}

// NoSideWays -
// A decision function used for when sideways move is not to be used.
//
// Arg 0: Current Value of H
//
// Arg 1: New Value of H
func NoSideWays(inputs ...int) bool {
	if inputs[1] < inputs[0] {
		return true
	}
	return false
}

// SideWays -
// A decision function used for when sideways move is to be used.
//
// Arg 0: Current Value of H
//
// Arg 1: New Value of H
//
// Arg 2: Step Count
//
// Arg 3: Step Limit
func SideWays(inputs ...int) bool {
	if inputs[0] == 0 || inputs[1] == 0 {
		return false
	}

	if inputs[1] < inputs[0] {
		return true
	}

	if inputs[2] == inputs[3] {
		return false
	} else {
		return true
	}

}
