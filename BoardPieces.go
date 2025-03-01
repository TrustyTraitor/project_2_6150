package main

import "fmt"

/*****************************************************************\
|                              Board                              |
\*****************************************************************/

// Board -
// A Simple type that is a slice of int slices.
type Board []BoardRow

// Set -
// Sets the value at (x,y)
// Values are stored mirrored over the X axis to align better with typical graphs.
func (b *Board) Set(x, y, v int) {
	(*b)[len(*b)-1-y][x] = v
}

// Get -
// Gets the value at (x,y)
// Values are stored mirrored over the X axis to align better with typical graphs.
func (b *Board) Get(x, y int) int {
	return (*b)[len(*b)-1-y][x]
}

// Print -
// Used to print the board in a more readable way
func (b *Board) Print() {
	for _, line := range *b {
		line.Print()
	}
}

type BoardRow []int

// Print -
// Prints the row of a board with some formatting for readability.
func (br BoardRow) Print() {
	for _, i := range br {
		if i != 0 {
			fmt.Printf("%-3v", i)
		} else {
			fmt.Print(".  ")
		}
	}
	fmt.Println()
}

// NewBoard -
// Constructor for the Board type.
//
// Initializes the required memory.
func NewBoard(size int) Board {
	board := make(Board, size)

	for idx := range board {
		board[idx] = make([]int, size)
	}

	return board
}

/*****************************************************************\
|                              Piece                              |
\*****************************************************************/

// Piece -
// The piece type is used to hold the x,y coordinate for a given queen.
type Piece struct {
	x int
	y int
}

// QueenHeuristic
// Counts how many pieces are causing threat on this piece.
func (p *Piece) QueenHeuristic(pieces []Piece) int {

	// CheckDiagonal
	// Checks if a piece is on either diagonal of another
	//
	CheckDiagonal := func(piece Piece) int {
		intercept1 := p.y - p.x
		check1 := (piece.x + intercept1 - piece.y) == 0

		intercept2 := p.y + p.x
		check2 := (-piece.x + intercept2 - piece.y) == 0

		if check1 || check2 {
			return 1
		}

		return 0
	}

	h := 0

	for _, piece := range pieces {
		// Checks if the piece is on the horizontal or vertical
		if p.x == piece.x || p.y == piece.y {
			h++
		} else {
			h += CheckDiagonal(piece)
		}
	}

	return h
}

// PrintPieces -
// Prints the pieces onto a board. Only shows queen locations, not the H values for moves.
func PrintPieces(pieces []Piece) {
	h := 0

	board := NewBoard(len(pieces))

	for idx, piece := range pieces {
		board[len(pieces)-1-piece.y][piece.x] = idx + 1
	}

	board.Print()

	for idx, piece := range pieces {
		h += piece.QueenHeuristic(pieces[idx+1:])
		//fmt.Printf("%v: %v\n", idx+1, piece)
	}
	fmt.Printf("H: %v\n\n", h)
}
