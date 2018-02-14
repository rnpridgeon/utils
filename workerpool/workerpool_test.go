package workerpool

import(
	"testing"
	//"log"
	"sync"
	"time"
)

const cycles = 1000

func TestCleaner(t *testing.T) {
	w := NewWorkerManager(10)
	wg := sync.WaitGroup{}
	//done := make(chan struct{})
	for i := 0; i <= cycles; i++ {
		wg.Add(1)
		taskCh := <-w.Ready
			taskCh.TaskChanel<-func(){
				wg.Done()
				}
	}

	wg.Wait()
	time.Sleep(cleanerCheckInterval)

	// There should always be 1 worker in the chamber
	if w.capacity != 9 {
		t.Error("Wrokerpool capacity failed to grow after ")
	}
}
