package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type config struct {
	maze       string
	searchType string
}

type application struct {
	cfg    config
	logger *slog.Logger
	maze   *Maze
}

func main() {
	var cfg config

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	flag.StringVar(&cfg.maze, "file", "maze.txt", "maze file")
	flag.StringVar(&cfg.searchType, "search", "DFS", "search type")

	flag.Parse()
	app := application{

		cfg:    cfg,
		logger: logger,
		maze:   &Maze{},
	}

	err := app.loadMaze(cfg.maze)

	if err != nil {

		app.logger.Error("Error : %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("height: %d and width : %d\n", app.maze.Height, app.maze.Width)

}
