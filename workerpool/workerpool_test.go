package workerpool

import(
	"testing"
	"time"
	"sync/atomic"
)

const cycles = 1000

type TestTask struct {
	id int
}

var processed, successful, failure  int64

func (t *TestTask) Process() {
	atomic.AddInt64(&processed, 1)
}

func (t *TestTask) onSuccess() {
	atomic.AddInt64(&successful, 1)
}

func (t *TestTask) onFailure(err error) {
	atomic.AddInt64(&successful, 1)
}

//TODO: Add actual tests
func TestExecutor(t *testing.T) {
	w := NewWorkerManager(10)
	processed = 0
	successful = 0
	failure = 0
	//wg := sync.WaitGroup{}
	for i := 0; i < cycles; i++ {
		t := &TestTask{ id : i}
		w.Execute(t.Process)
		w.Execute(t.onSuccess)
	}

	//wg.Wait()
	time.Sleep(1 * time.Second)
	if successful + failure != processed {
		t.Errorf("ERROR: TestExecutor : %d successful + %d failure != %d processed\n", successful, failure, processed)
	}
	}