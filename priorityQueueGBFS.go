package main

// our container/heap package that we used in the dijkstra.go file
// requires these method's that we create on t he PriorityQueueDijkstra

type PriorityQueueGBFS []*Node

func (pq PriorityQueueGBFS) Len() int {

	return len(pq)
}

func (pq PriorityQueueGBFS) Less(i, j int) bool {

	return pq[i].CostToGoal < pq[j].CostToGoal
}

func (pq PriorityQueueGBFS) Swap(i, j int) {

	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueueGBFS) Push(x any) {

	n := pq.Len()

	// convert the x to Node type
	item := x.(*Node)

	item.index = n

	*pq = append(*pq, item)

}

func (pq *PriorityQueueGBFS) Pop() any {
	old := *pq
	n := pq.Len()

	item := old[n-1]

	old[n-1] = nil

	item.index = -1

	*pq = old[:n-1]

	return item

}
