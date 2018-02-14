package collections

import (
	"fmt"
	"sync"
	"testing"
)

var (
	wg     sync.WaitGroup
	cycles int = 1000000
)

func TestStackDrain(t *testing.T) {
	stack := NewDEQueue()

	wg.Add(cycles)

	var i int
	for i = i; i < cycles; i++ {
		go func(i int) {
			stack.Push(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println(stack.GetDepth())
	wg.Add(stack.GetDepth())

	for ; cycles > 0; cycles-- {
		go func() {
			stack.Pop()
			wg.Done()
		}()
	}

	wg.Wait()

	if stack.GetDepth() != 0 {
		t.Errorf("FAILED: Failed to drain stack %d frames left", stack.GetDepth())
	}
}

// TODO: turn into an actual test
func TestStackWait(t *testing.T) {
	stack := NewDEQueue()

	var i int
	var waitVal interface{}
	go func(){ waitVal = stack.Watch()}()

	for i = 1; i < cycles; i++ {
		wg.Add(1)
		go func(i int) {
			stack.Push(i)
			wg.Done()
		}(i)
	}

	for ; cycles-1 > 0; cycles-- {
		wg.Add(1)
		go func() {
			stack.Pop()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("Val: %v\n", waitVal.(int))
}

func TestStackLIFO(t *testing.T) {
	stack := NewDEQueue()

	wg.Add(cycles)

	var i int
	for i = i; i < cycles; i++ {
		go func(i int) {
			stack.Push(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	testVal := stack.Pop()

	stack.Pop()
	stack.Push(testVal)

	testVal2 := stack.Pop()

	if testVal.(int) != testVal2.(int) {
		t.Errorf("FAILED: TestStackLIFO got %d expected %d", testVal, testVal2)
	}
}

func TestStackEmpty(t *testing.T) {
	stack := NewDEQueue()
	stack.Pop()
	stack.Pop()
}

