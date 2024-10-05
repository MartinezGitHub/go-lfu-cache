package lfu

// element: struct of List elements.
type element[E any] struct {
	next, prev *element[E]
	list       *list[E]
	value      E
}

// list: struct of List.
type list[E any] struct {
	root element[E]
	len  int
}

// init: initialize and return empty List.
func (l *list[E]) init() *list[E] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// newList return new initialized List.
func newList[E any]() *list[E] {
	return new(list[E]).init()
}

// newElement return new Element with given value.
func newElement[E any](value E) *element[E] {
	return &element[E]{value: value}
}

// Add adds element to the end of List and return pointer on it.
func (l *list[E]) Add(elementToAdd *element[E]) *element[E] {
	if elementToAdd.list != nil {
		return nil
	}

	last := l.root.prev

	elementToAdd.list = l
	elementToAdd.prev = last
	elementToAdd.next = &l.root

	last.next = elementToAdd
	l.root.prev = elementToAdd

	l.len++
	return elementToAdd
}

// Remove removes element with pointer from List.
func (l *list[E]) Remove(e *element[E]) {
	if e.list != l || l.len == 0 {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev
	l.len--
}

// Len returns List length.
func (l *list[E]) Len() int {
	return l.len
}

// PopBack deletes last element of list and returns pointer on it.
func (l *list[E]) PopBack() *element[E] {
	if l.len == 0 {
		return nil
	}
	last := l.root.prev

	last.prev.next = &l.root
	l.root.prev = last.prev

	last.next = nil
	last.prev = nil
	last.list = nil

	l.len--
	return last
}

// Front returns pointer on first element.
func (l *list[E]) Front() *element[E] {
	return l.root.next
}

// Next returns pointer on the next element of argument pointer element.
func (e element[E]) Next() *element[E] {
	return e.next
}

// AddBefore add element with given pointer before other element with given pointer.
func (l *list[E]) AddBefore(elementToAdd *element[E], elementOld *element[E]) {
	if elementOld == &l.root || elementOld.list != l || elementToAdd.list != nil {
		return
	}

	elementToAdd.list = l
	elementToAdd.next = elementOld
	elementToAdd.prev = elementOld.prev

	elementOld.prev.next = elementToAdd
	elementOld.prev = elementToAdd

	l.len++
}

// Back returns pointer on last element of the List.
func (l *list[E]) Back() *element[E] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// ReplaceDeletedElement pick the node between parent and child
// of the deleted node and increase List length.
func (l *list[E]) ReplaceDeletedElement(elementToAdd *element[E], deletedElement *element[E]) {
	if deletedElement.list != l || l.len == 0 || elementToAdd.list != nil {
		return
	}
	elementToAdd.next = deletedElement.next
	elementToAdd.prev = deletedElement.prev
	deletedElement.prev.next = elementToAdd
	deletedElement.next.prev = elementToAdd
	elementToAdd.list = l
	l.len++
}
