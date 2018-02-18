package workerpool

import (
	"sync"
	"sync/atomic"
	"github.com/rnpridgeon/utils/collections"
)

type Task interface {
	Process() error
	onSuccess()
	onFailure(error)
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


func (w *WorkerManager) run() {
	for {
		t := w.runnable.WatchQueue().(Task)

		w.L.Lock()
		if atomic.LoadInt64(&w.len) > w.capacity {
			w.Cond.Wait()
		}
		w.L.Unlock()
		w.decrement()

		go func(t Task){
			if err := t.Process(); err  == nil {
				t.onSuccess()
			} else {
				t.onFailure(err)
			}
			w.increment()
		}(t)

	}
}

// Executes task if we aren't at capacity, otherwise dumps on the runnable queue
func (w *WorkerManager) Execute(t Task) {
	w.runnable.Enqueue(t)
}

func (w *WorkerManager) SetCapacity(n int64) {
	w.capacity = atomic.SwapInt64(&w.capacity, n)
}

func (w *WorkerManager) increment() {
	atomic.AddInt64(&w.len, - 1)
	w.Cond.Signal()
}

func (w *WorkerManager) decrement() {
	atomic.AddInt64(&w.len, 1)
}



