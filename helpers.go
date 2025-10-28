package main

func inExplored(needle Point, items []Point) bool {

	for _, val := range items {

		if val.X == needle.X && val.Y == needle.Y {

			return true
		}

	}

	return false
}
