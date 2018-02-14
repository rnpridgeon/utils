package collections

import (
	"sync"
)

// Double-ended queue
type DEQueue struct {
	locker *sync.Cond
	elements []interface{}
	depth    int
}

func NewDEQueue() *DEQueue{
	return &DEQueue{
		locker :sync.NewCond(&sync.Mutex{}),
		elements: nil,
		depth:    0,
	}
}

func pop(deq *DEQueue) (element interface{}) {
	if deq.depth > 0 {
		element = deq.elements[deq.depth-1]
		deq.elements = deq.elements[:deq.depth-1]
		deq.depth--
	}
	return element
}

func popBack(deq *DEQueue) (element interface{}) {
	if deq.depth > 0 {
		element = deq.elements[0]
		deq.elements = deq.elements[1:]
		deq.depth--
	}
	return element
}

func push(deq *DEQueue, element interface{}) {
	deq.elements = append(deq.elements[:deq.depth], element)
	deq.depth++
}

// Return current depth of stack
func (deq *DEQueue) GetDepth() (depth int) {
	deq.locker.L.Lock()
	depth = deq.depth
	deq.locker.L.Unlock()
	return depth
}

// Return first element on the stack
func (deq *DEQueue) Pop() (element interface{}) {
	deq.locker.L.Lock()
		element = pop(deq)
	deq.locker.L.Unlock()
	return element
}

// Maybe I should have called it a Squeue?
func (deq *DEQueue) PopBack() (element interface{}) {
	deq.locker.L.Lock()
		element = popBack(deq)
	deq.locker.L.Unlock()
	return element
}

func (deq *DEQueue) Push(element interface{}) {
	deq.locker.L.Lock()
		push(deq, element)
		deq.locker.Signal()
	deq.locker.L.Unlock()
}

// Mimic channel behavior
func (deq *DEQueue) Watch() (element interface{}) {
	deq.locker.L.Lock()
		if len(deq.elements) == 0 {
			deq.locker.Wait()
		}
		element = pop(deq)
	deq.locker.L.Unlock()

	return element
}
