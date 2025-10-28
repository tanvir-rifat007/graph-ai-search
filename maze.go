package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// maze's point on x and y axis
type Point struct {
	X int
	Y int
}

// where the maze will be blocked
// in maze.txt the "#" is the wall
type Wall struct {
	State Point
	wall  bool
}

type Maze struct {
	Height int
	Width  int
	Start  Point
	Goal   Point
	Walls  [][]Wall
}

func (app *application) loadMaze(filename string) error {

	f, err := os.Open(filename)

	if err != nil {

		return err
	}

	defer f.Close()

	var fileContents []string

	reader := bufio.NewReader(f)

	for {

		line, err := reader.ReadString('\n')

		if err == io.EOF {

			break
		} else if err != nil {

			return err
		}

		fileContents = append(fileContents, line)

	}

	foundStart, foundEnd := false, false

	for _, line := range fileContents {

		if strings.Contains(line, "A") {

			foundStart = true
		}

		if strings.Contains(line, "B") {

			foundEnd = true
		}

	}

	if !foundStart {
		app.logger.Error("starting point 'A' not found")

	}

	if !foundEnd {

		app.logger.Error("ending point 'B' not found")
	}

	app.maze.Height = len(fileContents)
	app.maze.Width = len(fileContents[0])

	var rows [][]Wall

	for i, row := range fileContents {

		var cols []Wall

		for j, col := range row {
			currLetter := fmt.Sprintf("%c", col)
			var wall Wall

			switch currLetter {

			case "A":
				app.maze.Start = Point{X: i, Y: j}
				wall.State.X = i
				wall.State.Y = j
				wall.wall = false

			case "B":
				app.maze.Goal = Point{X: i, Y: j}
				wall.State.X = i
				wall.State.Y = j
				wall.wall = false

			case " ":
				wall.State.X = i
				wall.State.Y = j
				wall.wall = false

			case "#":
				wall.State.X = i
				wall.State.Y = j
				wall.wall = true

			default:
				continue
			}

			cols = append(cols, wall)

		}
		rows = append(rows, cols)

	}

	app.maze.Walls = rows

	return nil

}
