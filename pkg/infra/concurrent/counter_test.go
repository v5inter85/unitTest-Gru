package concurrent_test

import (
	"order-system/pkg/infra/concurrent"
	"sync"
	"testing"
)

func TestNewCounter(t *testing.T) {
	counter := concurrent.NewCounter(5)
	if counter.Value() != 5 {
		t.Errorf("Expected value 5, got %d", counter.Value())
	}
}

func TestIncrement(t *testing.T) {
	counter := concurrent.NewCounter(0)
	if counter.Increment() != 1 {
		t.Errorf("Expected value 1 after increment")
	}
	if counter.Value() != 1 {
		t.Errorf("Expected value 1, got %d", counter.Value())
	}
}

func TestDecrement(t *testing.T) {
	counter := concurrent.NewCounter(1)
	if counter.Decrement() != 0 {
		t.Errorf("Expected value 0 after decrement")
	}
	if counter.Value() != 0 {
		t.Errorf("Expected value 0, got %d", counter.Value())
	}
}

func TestAdd(t *testing.T) {
	counter := concurrent.NewCounter(5)
	if counter.Add(3) != 8 {
		t.Errorf("Expected value 8 after adding 3")
	}
	if counter.Value() != 8 {
		t.Errorf("Expected value 8, got %d", counter.Value())
	}

	if counter.Add(-5) != 3 {
		t.Errorf("Expected value 3 after adding -5")
	}
	if counter.Value() != 3 {
		t.Errorf("Expected value 3, got %d", counter.Value())
	}
}

func TestValue(t *testing.T) {
	counter := concurrent.NewCounter(10)
	if counter.Value() != 10 {
		t.Errorf("Expected value 10, got %d", counter.Value())
	}
}

func TestReset(t *testing.T) {
	counter := concurrent.NewCounter(5)
	counter.Reset()
	if counter.Value() != 0 {
		t.Errorf("Expected value 0 after reset, got %d", counter.Value())
	}
}

func TestCompareAndSwap(t *testing.T) {
	counter := concurrent.NewCounter(5)

	// Successful swap
	if !counter.CompareAndSwap(5, 10) {
		t.Error("Expected successful swap")
	}
	if counter.Value() != 10 {
		t.Errorf("Expected value 10, got %d", counter.Value())
	}

	// Failed swap
	if counter.CompareAndSwap(5, 15) {
		t.Error("Expected failed swap")
	}
	if counter.Value() != 10 {
		t.Errorf("Expected value 10, got %d", counter.Value())
	}
}

func TestConcurrentOperations(t *testing.T) {
	counter := concurrent.NewCounter(0)
	var wg sync.WaitGroup
	numGoroutines := 100

	// Test concurrent increments
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	if counter.Value() != int64(numGoroutines) {
		t.Errorf("Expected value %d, got %d", numGoroutines, counter.Value())
	}

	// Test concurrent decrements
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Decrement()
		}()
	}
	wg.Wait()
	if counter.Value() != 0 {
		t.Errorf("Expected value 0, got %d", counter.Value())
	}

	// Test concurrent adds
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Add(2)
		}()
	}
	wg.Wait()
	if counter.Value() != int64(numGoroutines*2) {
		t.Errorf("Expected value %d, got %d", numGoroutines*2, counter.Value())
	}
}
