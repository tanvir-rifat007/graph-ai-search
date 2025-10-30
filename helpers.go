package main

import "os"

func inExplored(needle Point, items []Point) bool {

	for _, val := range items {

		if val.X == needle.X && val.Y == needle.Y {

			return true
		}

	}

	return false
}

func emptyTmp() {
	directory := "./tmp/"
	dir, _ := os.Open(directory)
	filesToDelete, _ := dir.Readdir(0)

	for index := range filesToDelete {
		f := filesToDelete[index]
		fullPath := directory + f.Name()
		_ = os.Remove(fullPath)
	}
}

func abs(x int) int {

	if x < 0 {

		return -x
	}

	return x
}
