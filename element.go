package bufferpool

type Element struct {
	ref uint32
}

type ElementInterface interface {
	GetElement() *Element
	GetItem() interface{}
}
