package concurrent

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPool(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		wantSize int
	}{
		{
			name:     "Create pool with size 1",
			size:     1,
			wantSize: 1,
		},
		{
			name:     "Create pool with size 5",
			size:     5,
			wantSize: 5,
		},
		{
			name:     "Create pool with size 0",
			size:     0,
			wantSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewPool(tt.size)
			assert.NotNil(t, pool)
			assert.Equal(t, cap(pool.workers), tt.wantSize)
		})
	}
}

func TestPool_Submit(t *testing.T) {
	tests := []struct {
		name        string
		poolSize    int
		task        func() error
		wantErr     bool
		poolClosed  bool
		expectedErr error
	}{
		{
			name:     "Submit task to open pool",
			poolSize: 1,
			task: func() error {
				return nil
			},
			wantErr: false,
		},
		{
			name:     "Submit task that returns error",
			poolSize: 1,
			task: func() error {
				return errors.New("task error")
			},
			wantErr: false,
		},
		{
			name:       "Submit to closed pool",
			poolSize:   1,
			poolClosed: true,
			task: func() error {
				return nil
			},
			wantErr:     true,
			expectedErr: errors.New("pool is closed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewPool(tt.poolSize)
			if tt.poolClosed {
				pool.Close()
			}

			err := pool.Submit(tt.task)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPool_Close(t *testing.T) {
	t.Run("Close empty pool", func(t *testing.T) {
		pool := NewPool(1)
		pool.Close()
		assert.True(t, pool.closed)
	})

	t.Run("Close pool with running tasks", func(t *testing.T) {
		pool := NewPool(2)
		var wg sync.WaitGroup
		wg.Add(1)

		pool.Submit(func() error {
			time.Sleep(100 * time.Millisecond)
			wg.Done()
			return nil
		})

		pool.Close()
		wg.Wait()
		assert.True(t, pool.closed)
	})
}

func TestPool_ActiveTasks(t *testing.T) {
	t.Run("Check active tasks", func(t *testing.T) {
		pool := NewPool(2)
		assert.Equal(t, 0, pool.ActiveTasks())

		var wg sync.WaitGroup
		wg.Add(1)

		pool.Submit(func() error {
			time.Sleep(100 * time.Millisecond)
			wg.Done()
			return nil
		})

		time.Sleep(50 * time.Millisecond) // Wait for task to start
		assert.Equal(t, 1, pool.ActiveTasks())

		wg.Wait()
		time.Sleep(50 * time.Millisecond) // Wait for worker to be released
		assert.Equal(t, 0, pool.ActiveTasks())
	})
}
