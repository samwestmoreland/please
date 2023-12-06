// Package heap provides a pool to acquire memory arenas to use to allocate objects when parsing and interpreting asp.
// This package pools the areas to improve efficiency, as short-lived arenas, that allocate a handful of small objects
// don't perform well. By re-using arenas between package parses, we can have them live longer and allocate more
// objects.
package heap

import (
	"arena"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thought-machine/please/src/cli/logging"
)

var log = logging.Log

// Heap represents a memory arena we can use for asp heap allocated objects
type Heap struct {
	usages    int
	lastUsage time.Time
	Arena     *arena.Arena
	mux       sync.Mutex
}

func (h *Heap) free() {
	if h.Arena == nil {
		return
	}
	h.Arena.Free()
	h.Arena = nil
	h.usages = 0
	h.lastUsage = time.Time{}
}

// Pool is a struct for managing a pool of heaps, freeing them after a set number of usages, or if they've not been used
// for a set duration. This allows the heaps to be freed once we finish parsing.
type Pool struct {
	heaps     []*Heap
	available chan *Heap
	heapsMux  sync.Mutex
	// Ideally this would be based on bytes allocated, but we don't have access to those stats.
	UsagesBeforeFree  int
	IdleTimeUntilFree time.Duration
	Stats Stats
}

type Stats struct {
	Frees atomic.Uint64
	NewArena atomic.Uint64
}

// NewPool creates a new dynamically sized pool of heaps that can be used to allocated memory during parsing and
// interpreting asp code.
//
// usagesBeforeFree:  The number of times a heap will be used before the pool frees the underlying arena. This can be
//
//	negative, in which case the arena will never be freed by this heuristic.
//
// idleTimeUntilFree: The duration in which the underlying arena will be freed if this heap is not used. Idle time is
//
//	calculated from when the heap was returned to the pool.
func NewPool(size, usagesBeforeFree int, idleTimeUntilFree time.Duration) *Pool {
	pool := &Pool{
		heaps:             make([]*Heap, size),
		UsagesBeforeFree:  usagesBeforeFree,
		IdleTimeUntilFree: idleTimeUntilFree,
		available:         make(chan *Heap, size),
	}
	for i := 0; i < size; i++ {
		pool.heaps[i] = new(Heap)
		pool.available <- pool.heaps[i]
	}

	go pool.freeIdleHeaps()

	return pool
}

func (p *Pool) freeIdleHeaps() {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			for _, h := range p.heaps {
				if time.Since(h.lastUsage) < p.IdleTimeUntilFree {
					continue
				}
				if !h.mux.TryLock() {
					continue // Something must be using it so it's not idle
				}
				h.free()
				h.mux.Unlock()
			}
		}
	}
}

func (p *Pool) Get() *Heap {
	var heap = <-p.available
	heap.mux.Lock()
	if heap.Arena == nil {
		heap.Arena = arena.NewArena()
		p.Stats.NewArena.Add(1)
	}
	return heap
}

func (p *Pool) Return(heap *Heap) {
	heap.usages++
	if heap.usages > p.UsagesBeforeFree {
		heap.free()
		p.Stats.Frees.Add(1)
	}
	heap.mux.Unlock()
	p.available <- heap
}

func MakeSlice[T any](a *arena.Arena, len, cap int) []T {
	if a == nil {
		return make([]T, len, cap)
	}
	return arena.MakeSlice[T](a, len, cap)
}

func Append[T any](a *arena.Arena, slice []T, values ...T) []T {
	targetSize := len(slice) + len(values)
	if targetSize <= cap(slice) {
		return append(slice, values...)
	}
	var newSlice []T
	if a == nil {
		newSlice = make([]T, 0, targetSize*2)
	} else {
		newSlice = arena.MakeSlice[T](a, 0, targetSize*2)
	}
	return append(append(newSlice, slice...), values...)
}

func New[T any](a *arena.Arena) *T {
	if a == nil {
		return new(T)
	}
	return arena.New[T](a)
}
