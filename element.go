package bufferpool

type Element struct {
	ref int
}

type ElementInterface interface {
	GetElement() Element
	GetItem() interfce{}
}

