package main

// our container/heap package that we used in the dijkstra.go file
// requires these method's that we create on t he PriorityQueueAstar

type PriorityQueueAstar []*Node

func (pq PriorityQueueAstar) Len() int {

	return len(pq)
}

func (pq PriorityQueueAstar) Less(i, j int) bool {

	return int(pq[i].EstimatedCostToGoal) < int(pq[j].EstimatedCostToGoal)
}

func (pq PriorityQueueAstar) Swap(i, j int) {

	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueueAstar) Push(x any) {

	n := pq.Len()

	// convert the x to Node type
	item := x.(*Node)

	item.index = n

	*pq = append(*pq, item)

}

func (pq *PriorityQueueAstar) Pop() any {
	old := *pq
	n := pq.Len()

	item := old[n-1]

	old[n-1] = nil

	item.index = -1

	*pq = old[:n-1]

	return item

}
