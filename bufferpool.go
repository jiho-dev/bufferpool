package bufferpool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Package bufferpool is a simple wrapper around sync.Pool that is specific
// to bytes.Buffer.

/*
var global = New()

// Get returns a bytes.Buffer from a global pool
func Get() *bytes.Buffer {
	return global.Get()
}

// Release puts the given bytes.Buffer instance back in the global pool.
func Release(buf *bytes.Buffer) {
	global.Release(buf)
}
*/

// BufferPool is a sync.Pool for bytes.Buffer objects
type BufferPool struct {
	pool  sync.Pool
	count uint32
	alloc uint32
}

// New creates a new BufferPool instance
func New(new func() interface{}) *BufferPool {
	var bp BufferPool
	bp.pool.New = func() interface{} {
		newItem := new()

		atomic.AddUint32(&bp.alloc, 1)
		atomic.AddUint32(&bp.count, 1)
		/*
			runtime.SetFinalizer(newItem, func(newItem interface{}) {
				atomic.AddUint32(&bp.count, ^uint32(0))
			})
		*/

		return newItem
	}

	return &bp
}

func (bp *BufferPool) GetCount() (uint32, uint32) {
	return bp.alloc, bp.count
}

// Get returns a bytes.Buffer from the specified pool
func (bp *BufferPool) Get() interface{} {
	d := bp.pool.Get()
	e := d.(ElementInterface)
	el := e.GetElement()
	atomic.AddUint32(&el.ref, 1)
	el.pool = bp

	atomic.AddUint32(&bp.count, ^uint32(0))

	//return e.GetItem()
	return d
}

// Release puts the given bytes.Buffer back in the specified pool after
// resetting it
func (bp *BufferPool) Release(e ElementInterface) {
	el := e.GetElement()
	atomic.AddUint32(&el.ref, ^uint32(0))
	atomic.AddUint32(&bp.count, 1)

	if el.ref != 0 {
		fmt.Printf("Ref is not zero: %d \n", el.ref)
		el.ref = 0
	}

	//i := e.GetItem()
	bp.pool.Put(e)
}

func Release(e ElementInterface) {
	el := e.GetElement()
	el.pool.Release(el)
}
