package objectpool_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/jiho-dev/objectpool"
)

/////////////////////////////

type MyString struct {
	// common data
	objectpool.Element

	Buffer bytes.Buffer
}

func (i *MyString) GetElement() *objectpool.Element {
	return &i.Element
}

func (i *MyString) GetItem() interface{} {
	return i
}

// Put the element back regardless its counter.
func (i *MyString) Release() {
	i.GetElement().Release(i)
}

// Managing the ref count
func (i *MyString) ReleaseRef() {
	i.GetElement().ReleaseRef(i)
}

func newMyString() interface{} {
	return &MyString{}
}

///////////////////////////

func TestUsageMyString(t *testing.T) {
	pool := objectpool.New(newMyString, 5)

	allItems := []*MyString{}

	var useRefCount bool

	useRefCount = true

	// get a new item and hold it
	for i := 0; i < 10; i++ {
		item := pool.Get().(*MyString)
		item.Buffer.Reset()
		item.Buffer.WriteString(fmt.Sprintf("1st use, my index: %d", i))

		fmt.Printf("item=%p, %s \n", item, item.Buffer.String())

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
		item := pool.Get().(*MyString)
		item.Buffer.Reset()
		item.Buffer.WriteString(fmt.Sprintf("2nd use, my index: %d", i))

		fmt.Printf("item=%p, %s \n", item, item.Buffer.String())

		// managing the ref count
		item.HoldRef()
		allItems = append(allItems, item)
	}

	a, c = pool.GetCount()
	fmt.Printf("3. pool Count: allocated=%d, cached=%d \n", a, c)
}

func BenchmarkPoolMyString(t *testing.B) {
	t.ReportAllocs()

	pool := objectpool.New(newMyString, 0)

	fmt.Printf("Test Count: %d \n", t.N)

	var cached int = 1024
	allItems := make([]*MyString, cached)

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

		item := pool.Get().(*MyString)
		item.Buffer.Reset()
		item.Buffer.WriteString("1st use, my index: ")

		// XXX: occured the memory allocation
		// 1 allocation is needed at least
		// https://gist.github.com/evalphobia/caee1602969a640a4530
		item.Buffer.WriteString(strconv.Itoa(i))

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
