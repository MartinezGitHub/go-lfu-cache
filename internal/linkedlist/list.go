package linkedlist

import "iter"

// ListInterface interface to hide realization from user
type ListInterface[E any] interface {
	Add(element *Element[E]) *Element[E]
	Remove(e *Element[E])
	Len() int
	PopBack() *Element[E]
	Front() *Element[E]
	Back() *Element[E]
	AddBefore(newElement, existingElement *Element[E])
	ReplaceDeletedElement(newElement, deletedElement *Element[E])
	Iterator() iter.Seq[E]
	Root() *Element[E]
}

// ElementInterface interface to hide realization from user
type ElementInterface[E any] interface {
	GetNext() *Element[E]
}

// Element struct of List elements.
// It contains pointers to the next and previous elements
// and the value it holds.
type Element[E any] struct {
	next, prev *Element[E] // Pointers to the next and previous elements in the List.
	Value      E           // The value stored in this Element.
}

// List struct of List.
// It contains a sentinel root Element and the length of the List.
type List[E any] struct {
	root Element[E] // Sentinel node to simplify List operations, does not hold data.
	len  int        // Current number of elements in the List.
}

// init: initialize and return empty List.
// The root's next and prev pointers reference the root itself, forming a circular List.
func (l *List[E]) init() *List[E] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// link: links two List elements.
func (l *List[E]) link(firstNode *Element[E], secondNode *Element[E]) {
	firstNode.next = secondNode
	secondNode.prev = firstNode
}

// NewList return new initialized List.
func NewList[E any]() *List[E] {
	return new(List[E]).init() // Allocate a new List and initialize it.
}

// NewElement return new Element with given value.
func NewElement[E any](value E) *Element[E] {
	return &Element[E]{Value: value} // Create and return a new Element with the given value.
}

// Add adds Element to the end of List and return pointer on it.
// If the Element is already part of a List, it returns nil.
func (l *List[E]) Add(elementToAdd *Element[E]) *Element[E] {
	last := l.root.prev
	elementToAdd.prev = last
	elementToAdd.next = &l.root

	last.next = elementToAdd
	l.root.prev = elementToAdd

	l.len++ // Increment the List length.
	return elementToAdd
}

// Remove removes Element with pointer from List.
// Only removes the Element if it's part of the List.
func (l *List[E]) Remove(e *Element[E]) {
	if l.len == 0 { // Check if List is not empty.
		return
	}
	l.link(e.prev, e.next)
	l.len-- // Decrement List length.
}

// Len returns List length.
func (l *List[E]) Len() int {
	return l.len
}

// PopBack deletes last Element of List and returns pointer on it.
// If the List is empty, returns nil.
func (l *List[E]) PopBack() *Element[E] {
	if l.len == 0 { // Check if List is empty.
		return nil
	}
	last := l.root.prev

	last.prev.next = &l.root
	l.root.prev = last.prev

	last.next = nil
	last.prev = nil

	l.len-- // Decrement List length.
	return last
}

// Front returns pointer on first Element.
// Returns nil if the List is empty.
func (l *List[E]) Front() *Element[E] {
	return l.root.next
}

// GetNext returns pointer on the next Element of argument pointer Element.
// If at the end, returns reference to the root.
func (e Element[E]) GetNext() *Element[E] {
	return e.next
}

// AddBefore add Element with given pointer before other Element with given pointer.
// Does nothing if the existing Element is the root or if new Element is already linked.
func (l *List[E]) AddBefore(elementToAdd *Element[E], elementOld *Element[E]) {
	if elementOld == &l.root {
		return
	}
	elementToAdd.next = elementOld
	elementToAdd.prev = elementOld.prev

	elementOld.prev.next = elementToAdd
	elementOld.prev = elementToAdd

	l.len++ // Increment the List length.
}

// Back returns pointer on last Element of the List.
// Returns nil if the List is empty.
func (l *List[E]) Back() *Element[E] {
	if l.len == 0 { // Check if List is empty.
		return nil
	}
	return l.root.prev
}

// ReplaceDeletedElement pick the node between parent and child
// of the deleted node and increase List length.
func (l *List[E]) ReplaceDeletedElement(elementToAdd *Element[E], deletedElement *Element[E]) {
	if l.len == 0 {
		return
	}
	elementToAdd.next = deletedElement.next
	elementToAdd.prev = deletedElement.prev
	deletedElement.prev.next = elementToAdd
	deletedElement.next.prev = elementToAdd
	l.len++ // Increment the List length.
}

// Iterator returns iterator function witch iterates all List elements in the order they are stored.
func (l *List[E]) Iterator() iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := l.Front(); e != &l.root; e = e.next {
			if !yield(e.Value) { // yield function from value
				return
			}
		}
	}
}

// Root returns root
func (l *List[E]) Root() *Element[E] {
	return &l.root
}
