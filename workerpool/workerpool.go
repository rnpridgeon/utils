package workerpool

import (
	"sync"
	"github.com/rnpridgeon/utils/collections"
)

type Task interface {
	Process()
}

type WorkerManager struct {
	runnable  *collections.DEQueue
	*sync.Cond
	len int64
	capacity int64
}

// Returns reference to running worker manager
func NewWorkerManager(maxWorker int64) (w *WorkerManager) {
	w = &WorkerManager{
		runnable:     collections.NewDEQueue(),
		Cond:    &sync.Cond{L:&sync.Mutex{}},
		len: 0,
		capacity: maxWorker,
	}

	go w.run()

	return w
}

// Polls run queue for tasks, executing at most w.capacity tasks concurrently
func (w *WorkerManager) run() {
	for {
		t := w.runnable.WatchQueue().(func())

		w.L.Lock()
		if w.len > w.capacity {
			w.Cond.Wait()
		}
		w.L.Unlock()

		w.decrementCapacity()
		go func(t func()){
			t()
			w.incrementCapcity()
		}(t)

	}
}

// Adds task to the runnable queue
func (w *WorkerManager) Execute(t func()) {
	w.runnable.Enqueue(t)
}

func (w *WorkerManager) SetCapacity(n int64) {
	w.L.Lock()
		w.capacity = n
	w.L.Unlock()
}

func (w *WorkerManager) decrementCapacity() {
	w.L.Lock()
		w.len++
	w.L.Unlock()
}

func (w *WorkerManager) incrementCapcity() {
	w.L.Lock()
		w.len--
		w.Cond.Signal()
	w.L.Unlock()
}



