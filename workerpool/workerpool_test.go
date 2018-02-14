package workerpool

import(
	"testing"
	"time"
	"fmt"
	"sync/atomic"
)

const cycles = 1000

func TestCleaner(t *testing.T) {
	w := NewWorkerManager(10)
	var inFlight int64 = 0
	for i := 0; i <= cycles; i++ {
		taskCh := <-w.Ready
			atomic.AddInt64(&inFlight, 1)
			taskCh.TaskChanel<-func(){
						time.Sleep(1 * time.Second)
						atomic.AddInt64(&inFlight, -1)
				}
	}
	fmt.Println(inFlight)
	time.Sleep(cleanerCheckInterval + 1 * time.Second)
	// There should always be 1 worker in the chamber
	fmt.Println(inFlight)
	if w.capacity != 9 {
		t.Error("Wrokerpool capacity failed to grow after ")
	}
}
