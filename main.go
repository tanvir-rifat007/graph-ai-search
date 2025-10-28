package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var m Maze

	var maze, searchType string

	flag.StringVar(&maze, "file", "maze.txt", "maze file")
	flag.StringVar(&searchType, "search", "DFS", "search type")

	flag.Parse()

	err := m.loadMaze(maze)

	if err != nil {

		fmt.Errorf("Error : %w", err.Error())
		os.Exit(1)
	}

	fmt.Printf("height: %d and width : %d\n", m.Height, m.Width)

	startTime := time.Now()

	switch searchType {

	case "DFS":
		m.SearchType = DFS
		solveDFS(&m)
	default:
		fmt.Println("Invalid search Type")
		os.Exit(1)

	}

	if len(m.Solution.Actions) > 0 {

		fmt.Println("Solution : ")
		m.printMaze()
		fmt.Printf("Solution is : %d steps\n", len(m.Solution.Cells))
		fmt.Println("Total time taken : ", time.Since(startTime))

	} else {

		fmt.Println("No solution")
	}

	fmt.Printf("Explored %d nodes \n", len(m.Explored))

}

func solveDFS(m *Maze) {

	var s DepthFirstSearch
	s.Game = m
	fmt.Println("Goal is : ", s.Game.Goal)
	s.Solve()

}
