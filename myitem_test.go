package objectpool_test

import (
	"fmt"
	"testing"

	"github.com/jiho-dev/objectpool"
)

/////////////////////////////

type MyItem struct {
	// common data
	objectpool.Element

	// private data
	Idx  int
	Name string
}

func (i *MyItem) GetElement() *objectpool.Element {
	return &i.Element
}

func (i *MyItem) GetItem() interface{} {
	return i
}

// Put the element back regardless its counter.
func (i *MyItem) Release() {
	i.GetElement().Release(i)
}

// Managing the ref count
func (i *MyItem) ReleaseRef() {
	i.GetElement().ReleaseRef(i)
}

func newMyItem() interface{} {
	return &MyItem{}
}

///////////////////////////

func TestUsageMyItem(t *testing.T) {
	pool := objectpool.New(newMyItem, 5)

	allItems := []*MyItem{}

	var useRefCount bool

	useRefCount = true

	// get a new item and hold it
	for i := 0; i < 10; i++ {
		item := pool.Get().(*MyItem)

		fmt.Printf("item=%p \n", item)

		if useRefCount {
			// managing the ref count
			item.HoldRef()
		}

		allItems = append(allItems, item)
	}
	a, c := pool.GetCount()
	fmt.Printf("1. pool Count: allocated=%d, cached=%d \n", a, c)

	// XXX: use all items

	// release them
	for _, item := range allItems {

		// 1) use the pool pointer
		//pool.Release(item)

		if useRefCount {
			// 2) use the ref count
			item.ReleaseRef()
		} else {
			// 3) use the interface
			// to release it
			item.Release()
		}
	}

	a, c = pool.GetCount()
	fmt.Printf("2. pool Count: allocated=%d, cached=%d \n", a, c)

	// reuse the item
	for i := 0; i < 7; i++ {
		item := pool.Get().(*MyItem)

		fmt.Printf("item=%p \n", item)

		// managing the ref count
		item.HoldRef()
		allItems = append(allItems, item)
	}

	a, c = pool.GetCount()
	fmt.Printf("3. pool Count: allocated=%d, cached=%d \n", a, c)
}

func BenchmarkPoolMyItem(t *testing.B) {
	t.ReportAllocs()

	pool := objectpool.New(newMyItem, 0)

	fmt.Printf("Test Count: %d \n", t.N)

	var cached int = 1024
	allItems := [10]*MyItem{}

	for i := 0; i < t.N; i++ {
		// free them
		if i > 0 && (i%cached) == 0 {
			for j := 0; j < cached; j++ {
				item1 := allItems[j]
				// reuse it
				pool.Release(item1)
				allItems[j] = nil
			}
		}

		item := pool.Get().(*MyItem)
		allItems[i%cached] = item
	}

	for j := 0; j < cached; j++ {
		item1 := allItems[j]
		if item1 != nil {
			pool.Release(item1)
		}

		allItems[j] = nil
	}

	a, c := pool.GetCount()
	fmt.Printf("4. pool Count: allocated=%d, cached=%d \n", a, c)

}

/*
allocs/op means how many distinct memory allocations occurred per op (single iteration).
B/op is how many bytes were allocated per op.
*/

func BenchmarkAllocMyItem(t *testing.B) {
	t.ReportAllocs()

	fmt.Printf("Test Count: %d \n", t.N)

	allItems := [10]*MyItem{}
	for i := 0; i < t.N; i++ {
		// free them
		if i > 0 && (i%10) == 0 {
			for j := 0; j < 10; j++ {
				allItems[j] = nil
			}
		}

		item := &MyItem{}
		_ = item
		item.Idx = i
		allItems[i%10] = item
	}

	for j := 0; j < 10; j++ {
		allItems[j] = nil
	}
}
