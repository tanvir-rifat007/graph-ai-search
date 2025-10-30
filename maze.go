package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	DFS = iota
	BFS
	GBFS
	ASTAR
	DIJKSTRA
)

type Node struct {
	index  int
	State  Point
	Parent *Node
	Action string
	// from starting to the current node's cost for DIJKSTRA and other like A*,GBFS algorithms
	CostToGoal int
}

// calculate the cost from current node to the starting point
func (n *Node) ManhattanDistance(goal Point) int {

	return abs(n.State.X-goal.X) + abs(n.State.Y-goal.Y)

}

type Solution struct {
	Actions []string
	Cells   []Point
}

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
	Height      int
	Width       int
	Start       Point
	Goal        Point
	Walls       [][]Wall
	CurrentNode *Node
	Solution    Solution
	Explored    []Point
	Steps       int
	NumExplored int
	Debug       bool
	SearchType  int
	Animate     bool
}

func (m *Maze) loadMaze(filename string) error {

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
		fmt.Errorf("Starting point ('A') not found : %w", err.Error())

	}

	if !foundEnd {

		fmt.Errorf("Ending point ('B') not found : %w", err.Error())
	}

	m.Height = len(fileContents)
	m.Width = len(fileContents[0])

	var rows [][]Wall

	for i, row := range fileContents {

		var cols []Wall

		for j, col := range row {
			currLetter := fmt.Sprintf("%c", col)
			var wall Wall

			switch currLetter {

			case "A":
				m.Start = Point{X: i, Y: j}
				wall.State.X = i
				wall.State.Y = j
				wall.wall = false

			case "B":
				m.Goal = Point{X: i, Y: j}
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

	m.Walls = rows

	return nil

}

func (g *Maze) printMaze() {
	for r, row := range g.Walls {
		for c, col := range row {
			if col.wall {
				fmt.Print("█")
			} else if g.Start.X == col.State.X && g.Start.Y == col.State.Y {
				fmt.Print("A")
			} else if g.Goal.X == col.State.X && g.Goal.Y == col.State.Y {
				fmt.Print("B")
			} else if g.inSolution(Point{r, c}) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (g *Maze) inSolution(x Point) bool {
	for _, step := range g.Solution.Cells {
		if step.X == x.X && step.Y == x.Y {
			return true
		}
	}
	return false
}
