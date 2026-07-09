package main

import (
	"fmt"
	"image/png"
	"log"
	"main/klotski"
	"main/layout"
	"main/render"
	"os"
)

func main() {
	board := klotski.Board{
		Size: klotski.Vec2[int]{X: 4, Y: 5},
		Pieces: []klotski.Piece{
			{Pos: klotski.Vec2[int]{X: 0, Y: 0}, Size: klotski.Vec2[int]{X: 1, Y: 2}},
			{Pos: klotski.Vec2[int]{X: 1, Y: 0}, Size: klotski.Vec2[int]{X: 2, Y: 2}},
			{Pos: klotski.Vec2[int]{X: 3, Y: 0}, Size: klotski.Vec2[int]{X: 1, Y: 2}},
			{Pos: klotski.Vec2[int]{X: 0, Y: 2}, Size: klotski.Vec2[int]{X: 1, Y: 2}},
			{Pos: klotski.Vec2[int]{X: 1, Y: 2}, Size: klotski.Vec2[int]{X: 2, Y: 1}},
			{Pos: klotski.Vec2[int]{X: 3, Y: 2}, Size: klotski.Vec2[int]{X: 1, Y: 2}},
			{Pos: klotski.Vec2[int]{X: 1, Y: 3}, Size: klotski.Vec2[int]{X: 1, Y: 1}},
			{Pos: klotski.Vec2[int]{X: 2, Y: 3}, Size: klotski.Vec2[int]{X: 1, Y: 1}},
			{Pos: klotski.Vec2[int]{X: 0, Y: 4}, Size: klotski.Vec2[int]{X: 1, Y: 1}},
			{Pos: klotski.Vec2[int]{X: 3, Y: 4}, Size: klotski.Vec2[int]{X: 1, Y: 1}},
		},
	}
	fmt.Println(board)

	var graph klotski.Graph
	if persisted("graph") {
		fmt.Print("loading graph...")
		// im aware that it should try to load the layout and skip the graph but...

		var err error
		graph, err = getPersistence[klotski.Graph]("graph")
		if err != nil {
			panic(err)
		}

		fmt.Println("done")
	} else {
		fmt.Print("solving board...")

		var err error
		graph, err = klotski.Solve(board)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%d states, %d edges\n", len(graph.Nodes), len(graph.Edges))

		fmt.Print("saving graph...")

		if err := persist("graph", graph); err != nil {
			panic(err)
		}

		fmt.Println("done")
	}

	var positions []layout.Vec3
	if persisted("layout") {
		fmt.Print("loading layout...")

		var err error
		positions, err = getPersistence[[]layout.Vec3]("layout")
		if err != nil {
			panic(err)
		}

		fmt.Println("done")
	} else {
		fmt.Print("laying out...")
		positions = layout.ForceDirected3D(len(graph.Nodes), graph.Edges)
		fmt.Println("done")

		fmt.Print("saving layout...")
		if err := persist("layout", positions); err != nil {
			panic(err)
		}
		fmt.Println("done")
	}

	fmt.Print("rasterizing...")
	points := render.Rasterize(positions)
	fmt.Println("done")

	fmt.Print("rendering...")
	img := render.Render(points, graph.Edges)
	fmt.Println("done")

	f, err := os.Create("background.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
