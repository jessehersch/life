package cmd

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func getBoard(sizeX, sizeY int) [][]bool {
	// there are sizeY rows and sizeX columns per row. first index is row because I think of rows first when indexing a matrix
	// that's maybe confusing because to get a point you index by y first, then x. Which is opposite to the normal ordering of
	// points in cartesian coordinates (x,y). I considered flipping it... But I have usually stored matrices as row-major in
	// the past so that's what I did!
	// https://en.wikipedia.org/wiki/Row-_and_column-major_order
	board := make([][]bool, sizeY)
	for y := range board {
		board[y] = make([]bool, sizeX)
	}
	return board
}

func drawBoard(screen tcell.Screen, board [][]bool) {
	for y := range board {
		for x, isAlive := range board[y] {
			// +1 because of the border
			xx, yy := x+1, y+1
			if isAlive {
				screen.SetContent(xx, yy, tcell.RuneCkBoard, nil, tcell.StyleDefault)
			} else {
				screen.SetContent(xx, yy, ' ', nil, tcell.StyleDefault)
			}
		}
	}
	screen.Show()
}

func drawBorder(s tcell.Screen, board [][]bool) {
	sizeY := len(board)
	sizeX := len(board[0])
	style := tcell.StyleDefault
	s.SetContent(0, 0, tcell.RuneULCorner, nil, style)
	s.SetContent(sizeX+1, 0, tcell.RuneURCorner, nil, style)
	s.SetContent(0, sizeY+1, tcell.RuneLLCorner, nil, style)
	s.SetContent(sizeX+1, sizeY+1, tcell.RuneLRCorner, nil, style)
	for x := 1; x <= sizeX; x++ {
		s.SetContent(x, 0, tcell.RuneHLine, nil, style)
		s.SetContent(x, sizeY+1, tcell.RuneHLine, nil, style)
	}
	for y := 1; y <= sizeY; y++ {
		s.SetContent(0, y, tcell.RuneVLine, nil, style)
		s.SetContent(sizeX+1, y, tcell.RuneVLine, nil, style)
	}
}

func populateBoard(board [][]bool, filename string) {
	sizeY := len(board)
	sizeX := len(board[0])
	if filename == "" {
		// if caller doesn't specify an initial condition, use acorn as that's a small but rich starting point
		board[20%sizeY][50%sizeX] = true
		board[20%sizeY][51%sizeX] = true
		board[18%sizeY][51%sizeX] = true
		board[19%sizeY][53%sizeX] = true
		board[20%sizeY][54%sizeX] = true
		board[20%sizeY][55%sizeX] = true
		board[20%sizeY][56%sizeX] = true
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		var p [2]int
		if err := json.Unmarshal(s.Bytes(), &p); err != nil {
			log.Fatalf("%+v", err)
		}
		board[p[0]%sizeY][p[1]%sizeX] = true
	}
}
