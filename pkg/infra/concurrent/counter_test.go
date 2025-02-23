package concurrent

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounter(t *testing.T) {
	tests := []struct {
		name     string
		initial  int64
		expected int64
	}{
		{
			name:     "zero initial value",
			initial:  0,
			expected: 0,
		},
		{
			name:     "positive initial value",
			initial:  42,
			expected: 42,
		},
		{
			name:     "negative initial value",
			initial:  -42,
			expected: -42,
		},
		{
			name:     "large positive value",
			initial:  1<<31 - 1,
			expected: 1<<31 - 1,
		},
		{
			name:     "large negative value",
			initial:  -(1 << 31),
			expected: -(1 << 31),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewCounter(tt.initial)
			assert.Equal(t, tt.expected, counter.Value())
		})
	}
}

func TestCounter_Increment(t *testing.T) {
	tests := []struct {
		name     string
		initial  int64
		times    int
		expected int64
	}{
		{
			name:     "increment from zero",
			initial:  0,
			times:    1,
			expected: 1,
		},
		{
			name:     "multiple increments",
			initial:  0,
			times:    5,
			expected: 5,
		},
		{
			name:     "increment from negative",
			initial:  -5,
			times:    5,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewCounter(tt.initial)
			for i := 0; i < tt.times; i++ {
				counter.Increment()
			}
			assert.Equal(t, tt.expected, counter.Value())
		})
	}
}

func TestCounter_Decrement(t *testing.T) {
	tests := []struct {
		name     string
		initial  int64
		times    int
		expected int64
	}{
		{
			name:     "decrement from zero",
			initial:  0,
			times:    1,
			expected: -1,
		},
		{
			name:     "multiple decrements",
			initial:  0,
			times:    5,
			expected: -5,
		},
		{
			name:     "decrement from positive",
			initial:  5,
			times:    5,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewCounter(tt.initial)
			for i := 0; i < tt.times; i++ {
				counter.Decrement()
			}
			assert.Equal(t, tt.expected, counter.Value())
		})
	}
}

func TestCounter_Add(t *testing.T) {
	tests := []struct {
		name     string
		initial  int64
		delta    int64
		expected int64
	}{
		{
			name:     "add positive to zero",
			initial:  0,
			delta:    5,
			expected: 5,
		},
		{
			name:     "add negative to zero",
			initial:  0,
			delta:    -5,
			expected: -5,
		},
		{
			name:     "add zero",
			initial:  5,
			delta:    0,
			expected: 5,
		},
		{
			name:     "add to positive",
			initial:  5,
			delta:    3,
			expected: 8,
		},
		{
			name:     "add to negative",
			initial:  -5,
			delta:    3,
			expected: -2,
		},
		{
			name:     "add negative to negative",
			initial:  -5,
			delta:    -3,
			expected: -8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewCounter(tt.initial)
			result := counter.Add(tt.delta)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expected, counter.Value())
		})
	}
}

func TestCounter_Reset(t *testing.T) {
	tests := []struct {
		name    string
		initial int64
	}{
		{
			name:    "reset from zero",
			initial: 0,
		},
		{
			name:    "reset from positive",
			initial: 42,
		},
		{
			name:    "reset from negative",
			initial: -42,
		},
		{
			name:    "reset from large value",
			initial: 1 << 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewCounter(tt.initial)
			counter.Reset()
			assert.Equal(t, int64(0), counter.Value())
		})
	}
}

func TestCounter_CompareAndSwap(t *testing.T) {
	tests := []struct {
		name        string
		initial     int64
		old         int64
		new         int64
		expectSwap  bool
		expectValue int64
	}{
		{
			name:        "successful swap",
			initial:     5,
			old:         5,
			new:         10,
			expectSwap:  true,
			expectValue: 10,
		},
		{
			name:        "failed swap - wrong old value",
			initial:     5,
			old:         6,
			new:         10,
			expectSwap:  false,
			expectValue: 5,
		},
		{
			name:        "swap to same value",
			initial:     5,
			old:         5,
			new:         5,
			expectSwap:  true,
			expectValue: 5,
		},
		{
			name:        "swap with zero",
			initial:     5,
			old:         5,
			new:         0,
			expectSwap:  true,
			expectValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := NewCounter(tt.initial)
			swapped := counter.CompareAndSwap(tt.old, tt.new)
			assert.Equal(t, tt.expectSwap, swapped)
			assert.Equal(t, tt.expectValue, counter.Value())
		})
	}
}

func TestCounter_ConcurrentOperations(t *testing.T) {
	counter := NewCounter(0)
	numGoroutines := 100
	numOperations := 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // For both increment and decrement operations

	// Concurrent increments
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				counter.Increment()
			}
		}()
	}

	// Concurrent decrements
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				counter.Decrement()
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(0), counter.Value())
}

func TestCounter_ConcurrentAdd(t *testing.T) {
	counter := NewCounter(0)
	numGoroutines := 100
	addValue := int64(5)

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // For both positive and negative adds

	// Concurrent positive adds
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Add(addValue)
		}()
	}

	// Concurrent negative adds
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Add(-addValue)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(0), counter.Value())
}

func TestCounter_ConcurrentCompareAndSwap(t *testing.T) {
	counter := NewCounter(0)
	numGoroutines := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Multiple goroutines trying to CAS from 0 to 1
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			if counter.CompareAndSwap(0, 1) {
				counter.Add(1)
			}
		}()
	}

	wg.Wait()
	// Only one goroutine should have succeeded in the CAS operation
	assert.Equal(t, int64(2), counter.Value())
}

func TestCounter_ResetUnderLoad(t *testing.T) {
	counter := NewCounter(0)
	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines + 1) // +1 for the reset goroutine

	// Start goroutines that continuously increment
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.Increment()
			}
		}()
	}

	// Start a goroutine that resets periodically
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			counter.Reset()
		}
	}()

	wg.Wait()
	// Final value doesn't matter, we just want to ensure no panics or race conditions
}
