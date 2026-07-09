package klotski

import (
	"fmt"
	"strings"
)

type Board struct {
	Size    Vec2[int]
	Pieces  []Piece
	HasPins bool
}

type Piece struct {
	Size Vec2[int]
	Pos  Vec2[int]
}

type Vec2[T any] struct {
	X, Y T
}

func (v Vec2[any]) String() string {
	return fmt.Sprintf("(%v, %v)", v.X, v.Y)
}

func (b Board) String() string {
	var sb strings.Builder

	var curLetter byte = 'A'
	var letters = make(map[int]byte)

	grid := b.grid()
	if grid == nil {
		return "invalid board"
	}

	for i, row := range grid {
		for _, pieceId := range row {
			if pieceId == -1 {
				sb.WriteByte(' ')
				continue
			}

			if letter, ok := letters[pieceId]; ok {
				sb.WriteByte(letter)
			} else {
				letters[pieceId] = curLetter
				sb.WriteByte(curLetter)
				curLetter++
			}
		}

		if i < len(grid)-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
