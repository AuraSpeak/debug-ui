package services

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNextID(t *testing.T) {
	// Reset the counter before test
	idCounter = ids{
		nextID: 0,
		mu:     sync.Mutex{},
	}

	// Test sequential IDs
	for i := 0; i < 100; i++ {
		id := GetNextID()
		assert.Equal(t, i, id, "ID should be sequential")
	}
}

func TestGetNextID_Concurrent(t *testing.T) {
	// Reset the counter before test
	idCounter = ids{
		nextID: 0,
		mu:     sync.Mutex{},
	}

	const numGoroutines = 100
	const idsPerGoroutine = 10

	idsChan := make(chan int, numGoroutines*idsPerGoroutine)
	var wg sync.WaitGroup

	// Spawn multiple goroutines to get IDs concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id := GetNextID()
				idsChan <- id
			}
		}()
	}

	wg.Wait()
	close(idsChan)

	// Collect all IDs
	ids := make(map[int]bool)
	for id := range idsChan {
		ids[id] = true
	}

	// Verify we got the expected number of unique IDs
	expectedCount := numGoroutines * idsPerGoroutine
	require.Equal(t, expectedCount, len(ids), "Should have unique IDs for all goroutines")

	// Verify IDs are in expected range
	for i := 0; i < expectedCount; i++ {
		assert.True(t, ids[i], "ID %d should exist", i)
	}
}
