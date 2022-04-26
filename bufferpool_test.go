package bufferpool_test

import (
	"fmt"
	"testing"

	"github.com/lestrrat-go/bufferpool"
)

type Item struct {
	bufferpool.Element
	name string
}

func (i *Item) GetElement() *bufferpool.Element {
	return &i.Element
}

func (i *Item) GetItem() interface{} {
	return i
}

func allocItem() interface{} {
	return &Item{}
}

func TestUsage(t *testing.T) {
	pool := bufferpool.New(allocItem)

	all := []*Item{}

	for i := 0; i < 10; i++ {
		el := pool.Get().(*Item)

		/**
		if !assert.Equal(t, &bytes.Buffer{}, buf, `should be an empty buffer`) {
			return
		}
		*/
		fmt.Printf("el=%p \n", el)
		//pool.Release(el)

		all = append(all, el)
	}

	a, c := pool.GetCount()
	fmt.Printf("pool.Count: alloc=%d, count=%d \n", a, c)

	for _, e := range all {
		//pool.Release(e)

		// XXX: don't need pool pointer
		bufferpool.Release(e)
	}

	a, c = pool.GetCount()
	fmt.Printf("pool.Count: alloc=%d, count=%d \n", a, c)
}
