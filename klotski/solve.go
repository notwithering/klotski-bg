package klotski

import (
	"errors"
	"fmt"
	"strings"
)

func Solve(start Board) (Graph, error) {
	if start.grid() == nil {
		return Graph{}, errors.New("invalid starting state")
	}

	index := map[string]int{}
	var g Graph
	seenEdges := map[[2]int]bool{}

	startID := len(g.Nodes)
	index[start.key()] = startID
	g.Nodes = append(g.Nodes, start)

	queue := []int{startID}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		curBoard := g.Nodes[cur]

		for _, neighbor := range curBoard.neighbors() {
			k := neighbor.key()
			id, known := index[k]
			if !known {
				if len(g.Nodes) >= 5000 {
					continue
				}

				id = len(g.Nodes)
				index[k] = id
				g.Nodes = append(g.Nodes, neighbor)
				queue = append(queue, id)
			}

			e := [2]int{cur, id}
			if e[0] > e[1] {
				e[0], e[1] = e[1], e[0]
			}
			if e[0] != e[1] && !seenEdges[e] {
				seenEdges[e] = true
				g.Edges = append(g.Edges, e)
			}
		}
	}

	return g, nil
}

func (b Board) grid() [][]int {
	grid := make([][]int, b.Size.Y)
	for y := range grid {
		row := make([]int, b.Size.X)
		for x := range row {
			row[x] = -1
		}
		grid[y] = row
	}

	for i, piece := range b.Pieces {
		if piece.Pos.X < 0 ||
			piece.Pos.Y < 0 ||
			piece.Pos.X+piece.Size.X > b.Size.X ||
			piece.Pos.Y+piece.Size.Y > b.Size.Y {
			return nil
		}

		for dy := 0; dy < piece.Size.Y; dy++ {
			for dx := 0; dx < piece.Size.X; dx++ {
				cell := &grid[piece.Pos.Y+dy][piece.Pos.X+dx]
				if *cell != -1 {
					// overlappsing
					return nil
				}
				*cell = i
			}
		}
	}
	return grid
}

func (b Board) key() string {
	var sb strings.Builder
	for _, piece := range b.Pieces {
		fmt.Fprintf(&sb, "%d,%d;", piece.Pos.X, piece.Pos.Y) // this could be better
	}
	return sb.String()
}

func (b Board) neighbors() []Board {
	var out []Board

	for i, piece := range b.Pieces {
		for _, direction := range []Vec2[int]{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}} {
			if b.HasPins {
				if piece.Size.X > 1 && direction.Y != 0 {
					continue
				}

				if piece.Size.Y > 1 && direction.X != 0 {
					continue
				}
			}

			neighbor := b.Clone()

			neighbor.Pieces[i].Pos = Vec2[int]{
				X: piece.Pos.X + direction.X,
				Y: piece.Pos.Y + direction.Y,
			}

			if neighbor.grid() != nil {
				out = append(out, neighbor)
			}
		}
	}

	return out
}

func (b Board) Clone() Board {
	pieces := make([]Piece, len(b.Pieces))
	copy(pieces, b.Pieces)

	return Board{
		Size:    b.Size,
		Pieces:  pieces,
		HasPins: b.HasPins,
	}
}
