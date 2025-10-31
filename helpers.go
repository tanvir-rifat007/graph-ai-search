package main

import (
	"math"
	"os"
)

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

func euclideanDist(p, goal Point) float64 {
	return math.Sqrt(float64(p.X-goal.X)*float64(p.X-goal.X) + float64(p.Y-goal.Y)*float64(p.Y-goal.Y))

}
