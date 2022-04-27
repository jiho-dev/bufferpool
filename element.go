package objectpool

import "sync/atomic"

type Element struct {
	pool *ObjectPool
	ref  int64
}

/*
func (el *Element) GetElement() *Element {
	return el
}

func (el *Element) GetItem() interface{} {
	return el
}
*/

// Put the element back regardless its counter.
func (el *Element) Release(e ElementInterface) {
	el.pool.Release(e)
}

func (el *Element) HoldRef() {
	atomic.AddInt64(&el.ref, 1)
}

// Decrease its counter
// and release it if the counter is zero
func (el *Element) ReleaseRef(e ElementInterface) {
	if atomic.AddInt64(&el.ref, -1) > 0 {
		return
	}

	// release itself
	el.Release(e)
}

type ElementInterface interface {
	GetElement() *Element
	GetItem() interface{}

	/*
		// use reference count
		HoldRef()
		ReleaseRef(e ElementInterface)
		Release(e ElementInterface)
	*/
}
