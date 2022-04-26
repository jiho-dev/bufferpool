package bufferpool

type Element struct {
	pool *BufferPool
	ref  uint32
}

func (el *Element) GetElement() *Element {
	return el
}

type ElementInterface interface {
	GetElement() *Element
	//GetItem() interface{}
}
