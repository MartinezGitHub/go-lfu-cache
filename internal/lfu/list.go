package lfu

// element: struct of List elements.
// It contains pointers to the next and previous elements,
// a reference to the list it belongs to, and the value it holds.
type element[E any] struct {
	next, prev *element[E] // Pointers to the next and previous elements in the list.
	list       *list[E]    // Pointer to the list this element belongs to.
	value      E           // The value stored in this element.
}

// list: struct of List.
// It contains a sentinel root element and the length of the list.
type list[E any] struct {
	root element[E] // Sentinel node to simplify list operations, does not hold data.
	len  int        // Current number of elements in the list.
}

// init: initialize and return empty List.
// The root's next and prev pointers reference the root itself, forming a circular list.
func (l *list[E]) init() *list[E] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// newList return new initialized List.
func newList[E any]() *list[E] {
	return new(list[E]).init() // Allocate a new list and initialize it.
}

// newElement return new Element with given value.
func newElement[E any](value E) *element[E] {
	return &element[E]{value: value} // Create and return a new element with the given value.
}

// Add adds element to the end of List and return pointer on it.
// If the element is already part of a list, it returns nil.
func (l *list[E]) Add(elementToAdd *element[E]) *element[E] {
	if elementToAdd.list != nil { // Check if the element is already in a list.
		return nil
	}

	last := l.root.prev

	elementToAdd.list = l
	elementToAdd.prev = last
	elementToAdd.next = &l.root

	last.next = elementToAdd
	l.root.prev = elementToAdd

	l.len++ // Increment the list length.
	return elementToAdd
}

// Remove removes element with pointer from List.
// Only removes the element if it's part of the list.
func (l *list[E]) Remove(e *element[E]) {
	if e.list != l || l.len == 0 { // Check if element belongs to list and list is not empty.
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev
	l.len-- // Decrement list length.
}

// Len returns List length.
func (l *list[E]) Len() int {
	return l.len
}

// PopBack deletes last element of list and returns pointer on it.
// If the list is empty, returns nil.
func (l *list[E]) PopBack() *element[E] {
	if l.len == 0 { // Check if list is empty.
		return nil
	}
	last := l.root.prev

	last.prev.next = &l.root
	l.root.prev = last.prev

	last.next = nil
	last.prev = nil
	last.list = nil

	l.len-- // Decrement list length.
	return last
}

// Front returns pointer on first element.
// Returns nil if the list is empty.
func (l *list[E]) Front() *element[E] {
	return l.root.next
}

// Next returns pointer on the next element of argument pointer element.
// If at the end, returns reference to the root.
func (e element[E]) Next() *element[E] {
	return e.next
}

// AddBefore add element with given pointer before other element with given pointer.
// Does nothing if the existing element is the root or if new element is already linked.
func (l *list[E]) AddBefore(elementToAdd *element[E], elementOld *element[E]) {
	if elementOld == &l.root || elementOld.list != l || elementToAdd.list != nil {
		return
	}

	elementToAdd.list = l
	elementToAdd.next = elementOld
	elementToAdd.prev = elementOld.prev

	elementOld.prev.next = elementToAdd
	elementOld.prev = elementToAdd

	l.len++ // Increment the list length.
}

// Back returns pointer on last element of the List.
// Returns nil if the list is empty.
func (l *list[E]) Back() *element[E] {
	if l.len == 0 { // Check if list is empty.
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
	l.len++ // Increment the list length.
}
