package cmd

type point struct {
	x int
	y int
}

func getLiveNeighborCount(x, y int, board [][]bool) int {
	sizeY := len(board)
	sizeX := len(board[0])
	count := 0
	neighbors := []point{
		// our grid is a torus. like video game "asteroids"
		// adding sizeX and sizeY whenever the neighbor index is subtracted from, to handle walking off
		// the grid in the negative direction. the mod below takes care of the rest.
		{sizeX + x - 1, sizeY + y - 1}, {x, sizeY + y - 1}, {x + 1, sizeY + y - 1},
		{sizeX + x - 1, y}, {x + 1, y},
		{sizeX + x - 1, y + 1}, {x, y + 1}, {x + 1, y + 1}}
	for i := range neighbors {
		// mod handles going off to the right or bottom.
		if board[neighbors[i].y%sizeY][neighbors[i].x%sizeX] {
			count++
		}
	}
	return count
}

func iterate(board [][]bool) bool {
	var switchStatePoints []point
	for y := range board {
		for x, isAlive := range board[y] {
			liveNeighborCount := getLiveNeighborCount(x, y, board)
			if (isAlive && (liveNeighborCount < 2 || liveNeighborCount > 3)) || (!isAlive && liveNeighborCount == 3) {
				switchStatePoints = append(switchStatePoints, point{x, y})
			}
		}
	}
	for i := range switchStatePoints {
		p := switchStatePoints[i]
		board[p.y][p.x] = !board[p.y][p.x]
	}
	return len(switchStatePoints) > 0
}
