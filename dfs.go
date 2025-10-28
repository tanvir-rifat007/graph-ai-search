package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

type DepthFirstSearch struct {
	Frontier []*Node
	Game     *Maze
}

func (dfs *DepthFirstSearch) GetFrontier() []*Node {
	return dfs.Frontier
}

func (dfs *DepthFirstSearch) Add(i *Node) {
	dfs.Frontier = append(dfs.Frontier, i)
}

func (dfs *DepthFirstSearch) ContainsState(i *Node) bool {
	for _, x := range dfs.Frontier {
		if x.State == i.State {
			return true
		}
	}
	return false
}

func (dfs *DepthFirstSearch) Empty() bool {
	return len(dfs.Frontier) == 0
}

func (dfs *DepthFirstSearch) Remove() (*Node, error) {
	if len(dfs.Frontier) > 0 {
		if dfs.Game.Debug {
			fmt.Println("Frontier before remove:")
			for _, x := range dfs.Frontier {
				fmt.Println("Node:", x.State)
			}
		}
		node := dfs.Frontier[len(dfs.Frontier)-1]
		dfs.Frontier = dfs.Frontier[:len(dfs.Frontier)-1]
		return node, nil
	}
	return nil, errors.New("frontier is empty")
}

func (dfs *DepthFirstSearch) Solve() {
	fmt.Println("Starting to solve maze using Depth First Search...")
	dfs.Game.NumExplored = 0

	start := Node{
		State:  dfs.Game.Start,
		Parent: nil,
		Action: "",
	}

	dfs.Add(&start)
	dfs.Game.CurrentNode = &start

	for {
		if dfs.Empty() {
			return
		}

		currentNode, err := dfs.Remove()
		if err != nil {
			log.Println(err)
			return
		}

		if dfs.Game.Debug {
			fmt.Println("Removed", currentNode.State)
			fmt.Println("-------")
			fmt.Println()
		}

		dfs.Game.CurrentNode = currentNode
		dfs.Game.NumExplored += 1

		// Have we found the solution?
		if dfs.Game.Goal == currentNode.State {
			var actions []string
			var cells []Point

			for {
				if currentNode.Parent != nil {
					// this is traversing child to parent(goal to start)
					actions = append(actions, currentNode.Action)
					cells = append(cells, currentNode.State)
					currentNode = currentNode.Parent
				} else {
					break
				}
			}

			// rever this(now it becomes start to goal)
			slices.Reverse(actions)
			slices.Reverse(cells)

			dfs.Game.Solution = Solution{
				Actions: actions,
				Cells:   cells,
			}
			dfs.Game.Explored = append(dfs.Game.Explored, currentNode.State)
			break
		}

		dfs.Game.Explored = append(dfs.Game.Explored, currentNode.State)

		for _, x := range dfs.Neighbors(currentNode) {
			if !dfs.ContainsState(x) {
				if !inExplored(x.State, dfs.Game.Explored) {
					dfs.Add(&Node{
						State:  x.State,
						Parent: currentNode,
						Action: x.Action,
					})
				}
			}
		}
	}
}

func (dfs *DepthFirstSearch) Neighbors(node *Node) []*Node {
	row := node.State.X
	col := node.State.Y

	// possible neighbors (that's why i named it candidates)
	candidates := []*Node{
		{State: Point{X: row - 1, Y: col}, Parent: node, Action: "up"},
		{State: Point{X: row + 1, Y: col}, Parent: node, Action: "down"},
		{State: Point{X: row, Y: col - 1}, Parent: node, Action: "left"},
		{State: Point{X: row, Y: col + 1}, Parent: node, Action: "right"},
	}

	var neighbors []*Node
	for _, x := range candidates {
		if 0 <= x.State.X && x.State.X < dfs.Game.Height {
			if 0 <= x.State.Y && x.State.Y < dfs.Game.Width {
				if !dfs.Game.Walls[x.State.X][x.State.Y].wall {
					neighbors = append(neighbors, x)
				}
			}
		}
	}

	// randomness of each node's neighbors each time

	for i := range neighbors {
		j := rand.Intn(i + 1)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}

	return neighbors
}
