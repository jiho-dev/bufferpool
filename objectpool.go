package objectpool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ObjectPool is a sync.Pool for objects
type ObjectPool struct {
	pool           sync.Pool
	max            uint32
	cachedCount    uint64
	allocatedCount uint64
}

// New creates a new ObjectPool instance
func New(new func() interface{}, max uint32) *ObjectPool {
	var bp ObjectPool
	bp.pool.New = func() interface{} {
		newItem := new()

		atomic.AddUint64(&bp.allocatedCount, 1)
		atomic.AddUint64(&bp.cachedCount, 1)

		/*
			runtime.SetFinalizer(newItem, func(newItem interface{}) {
				atomic.AddUint64(&bp.cachedCount, ^uint32(0))
				fmt.Printf("Delete item memory: %p \n", newItem)
			})
		*/

		return newItem
	}

	bp.max = max

	return &bp
}

func (bp *ObjectPool) GetCount() (uint64, uint64) {
	a, c := atomic.LoadUint64(&bp.allocatedCount), atomic.LoadUint64(&bp.cachedCount)

	return a, c
}

// Get returns a object from the specified pool
func (bp *ObjectPool) Get() interface{} {
	item := bp.pool.Get()

	e := item.(ElementInterface)
	el := e.GetElement()
	//el.HoldRef()
	el.ref = 0
	el.pool = bp

	atomic.AddUint64(&bp.cachedCount, ^uint64(0))

	return item
}

// Release puts the given object back in the specified pool after resetting it
func (bp *ObjectPool) Release(e ElementInterface) {
	el := e.GetElement()
	//el.ReleaseRef()

	if el.ref != 0 {
		fmt.Printf("Ref is not zero: %d \n", el.ref)
		el.ref = 0
	}

	// full
	if bp.max > 0 && atomic.LoadUint64(&bp.cachedCount) >= uint64(bp.max) {
		return
	}

	atomic.AddUint64(&bp.cachedCount, 1)

	item := e.GetItem()
	bp.pool.Put(item)
}
