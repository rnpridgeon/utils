package workerpool

import (
	"log"
	"sync"
	"time"
	"github.com/rnpridgeon/utils/collections"
)

// TODO: make configurable?
const maxIdleTime = 30 * time.Second
const cleanerCheckInterval = maxIdleTime + 1 * time.Second

func startCleaner(tickTime time.Duration, clean func(time.Time)) {
	for _ = range time.NewTicker(tickTime).C {
		log.Println("INFO: Reaping expired workers")
		clean(time.Now())
	}
}

type WorkerManager struct {
	free  *collections.DEQueue
	Ready chan *worker
	Close chan struct{}
	*sync.Mutex
	capacity int
}

// Returns reference to running worker manager
func NewWorkerManager(maxWorker int) (w *WorkerManager) {
	w = &WorkerManager{
		free:     collections.NewDEQueue(),
		Ready:    make(chan *worker),
		Close:    make(chan struct{}),
		Mutex:    &sync.Mutex{},
		capacity: maxWorker,
	}

	go startCleaner(cleanerCheckInterval, w.maybeClean)
	go w.run()

	return w
}

// Put first available worker on ready channel for consumption
func (w *WorkerManager) run() {

	for {
		if w.free.GetDepth() == 0 {
			maybeAdd(w)
		}
		w.Ready <- w.free.Watch().(*worker)
	}
}

// Add worker to pool assuming capacity has not yet been met
func maybeAdd(w *WorkerManager) {
	w.Lock()
	if w.capacity > 0 {
		w.free.Push(newWorker(w.free))
		log.Printf("INFO: Adding worker %d to the pool", w.capacity)
		w.capacity--
	}
	w.Unlock()
}

// shrink worker pool if workers are sitting around for more than max idle
func (w *WorkerManager) maybeClean(deadline time.Time) {
	for idle := w.free.PopBack(); idle != nil; idle = w.free.PopBack() {
		idle := idle.(*worker)
		if idle.expiry.Before(deadline) {
			log.Printf("INFO: Removing worker with expiry %v from the pool", idle.expiry)

			idle.stop <- struct{}{}
			idle = nil

			w.capacity++
			continue
		}
		w.free.Push(idle)
		// no need to continue, oldest eligible worker was current
		break;
	}
}

// TODO: Add maybeShrink otherwise underlying array remains in tact
//func (*workerManager) maybeShrink() (int) {
//	w.Lock()
//
//	w.Unlock()
//}

type worker struct {
	expiry     time.Time
	TaskChanel chan func()
	stop       chan struct{}
	registry   *collections.DEQueue
}

func newWorker(registry *collections.DEQueue) (w *worker) {
	w = &worker{
		expiry:     time.Now().Add(maxIdleTime),
		TaskChanel: make(chan func()),
		stop:       make(chan struct{}),
		registry:   registry,
	}

	go w.run()

	return w
}

func (w *worker) run() {
	for {
		select {
		case work := <-w.TaskChanel:
			work()
			w.registry.Push(w)
		case <-w.stop:
			close(w.TaskChanel)
			return
		}
	}
}

