package main

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

type AstrSearch struct {
	Frontier PriorityQueueAstar
	Game     *Maze
}

func (d *AstrSearch) GetFrontier() []*Node {
	return d.Frontier
}

func (d *AstrSearch) Add(i *Node) {
	i.CostToGoal = i.ManhattanDistance(d.Game.Start)

	i.EstimatedCostToGoal = euclideanDist(i.State, d.Game.Goal) + float64(i.CostToGoal)

	d.Frontier.Push(i)

	heap.Init(&d.Frontier)

}

func (d *AstrSearch) ContainsState(i *Node) bool {

	for _, x := range d.Frontier {

		if x.State == i.State {
			return true

		}

	}

	return false

}

func (d *AstrSearch) Empty() bool {

	return len(d.Frontier) == 0

}

func (d *AstrSearch) Remove() (*Node, error) {

	if len(d.Frontier) > 0 {

		if d.Game.Debug {
			fmt.Println("Before removing...")
			for _, val := range d.Frontier {

				fmt.Println("Node: ", val.State)
			}

		}

		// using PriorityQueue
		// because AstrSearch use priorityQueue

		return heap.Pop(&d.Frontier).(*Node), nil
	}

	return nil, errors.New("Frontier is empty!")

}

func (d *AstrSearch) Solve() {

	fmt.Println("Starting to solve maze using AStar Search...")

	d.Game.NumExplored = 0

	start := Node{

		State:  d.Game.Start,
		Parent: nil,
		Action: "",
	}

	d.Add(&start)
	d.Game.CurrentNode = &start

	for {

		if d.Empty() {

			return
		}
		currentNode, err := d.Remove()

		if err != nil {

			log.Println(err)
			return
		}

		if d.Game.Debug {

			fmt.Println("Removed", currentNode.State)
			fmt.Println("-------")
			fmt.Println()
		}

		d.Game.CurrentNode = currentNode
		d.Game.NumExplored++
		// Have we found the solution?
		if d.Game.Goal == currentNode.State {
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

			d.Game.Solution = Solution{
				Actions: actions,
				Cells:   cells,
			}
			d.Game.Explored = append(d.Game.Explored, currentNode.State)
			break
		}

		d.Game.Explored = append(d.Game.Explored, currentNode.State)

		// Build animation frame if appropriate.
		if d.Game.Animate {
			d.Game.OutputImage(fmt.Sprintf("tmp/%06d.png", d.Game.NumExplored))
		}
		for _, x := range d.Neighbors(currentNode) {
			if !d.ContainsState(x) {
				if !inExplored(x.State, d.Game.Explored) {
					d.Add(&Node{
						State:  x.State,
						Parent: currentNode,
						Action: x.Action,
					})
				}
			}
		}

	}

}

func (d *AstrSearch) Neighbors(node *Node) []*Node {
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
		if 0 <= x.State.X && x.State.X < d.Game.Height {
			if 0 <= x.State.Y && x.State.Y < d.Game.Width {
				if !d.Game.Walls[x.State.X][x.State.Y].wall {
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
