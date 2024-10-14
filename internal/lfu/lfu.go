package lfu

import (
	"errors"
	"iter"
	"lfucache/internal/linkedlist"
)

// Main algorithm idea: usage of:
// 1) linked linkedlist for storing all cache elements.
// linkedlist structure: blocks of elements of each frequency.
// The blocks are ordered in descending frequency.
// Elements in the block are ordered by recency of use.
// 2) Frequency to first Element of block map.
// 3) Key to Element map.
// Pic notation: ()_i - block of i frequency, + - Element in block. Example: (+, +)_3 - (+)_1.

var ErrKeyNotFound = errors.New("key not found")

const DefaultCapacity = 5

// Cache
// O(capacity) memory
type Cache[K comparable, V any] interface {
	// Get returns the value of the key if the key exists in the cache,
	// otherwise, returns ErrKeyNotFound.
	//
	// O(1)
	Get(key K) (V, error)

	// Put updates the value of the key if present, or inserts the key if not already present.
	//
	// When the cache reaches its capacity, it should invalidate and remove the least frequently used key
	// before inserting a new item. For this problem, when there is a tie
	// (i.e., two or more keys with the same frequency), the least recently used key would be invalidated.
	//
	// O(1)
	Put(key K, value V)

	// All returns the iterator in descending order of frequency.
	// If two or more keys have the same frequency, the most recently used key will be listed first.
	//
	// O(capacity)
	All() iter.Seq2[K, V]

	// Size returns the cache size.
	//
	// O(1)
	Size() int

	// Capacity returns the cache capacity.
	//
	// O(1)
	Capacity() int

	// GetKeyFrequency returns the element's frequency if the key exists in the cache,
	// otherwise, returns ErrKeyNotFound.
	//
	// O(1)
	GetKeyFrequency(key K) (int, error)
}

// CacheImpl represents LFU cache implementation
type CacheImpl[K comparable, V any] struct {
	listOfData               linkedlist.ListInterface[*node[K, V]]
	keyToElement             map[K]*linkedlist.Element[*node[K, V]]
	frequencyToRecentElement map[int]*linkedlist.Element[*node[K, V]]
	capacity                 int
	defaultValue             V
}

// New initializes the cache with the given capacity.
// If no capacity is provided, the cache will use DefaultCapacity.
func New[K comparable, V any](capacity ...int) *CacheImpl[K, V] {
	var cacheCapacity int
	if len(capacity) > 0 {
		if capacity[0] < 0 {
			panic("Negative capacity does not supported.")
		}
		cacheCapacity = capacity[0]
	} else {
		cacheCapacity = DefaultCapacity
	}
	return &CacheImpl[K, V]{
		listOfData:               linkedlist.NewList[*node[K, V]](),
		keyToElement:             make(map[K]*linkedlist.Element[*node[K, V]], cacheCapacity),
		frequencyToRecentElement: make(map[int]*linkedlist.Element[*node[K, V]], cacheCapacity),
		capacity:                 cacheCapacity,
	}
}

// Get removes Element from linkedlist, increase its frequency and puts it back in linkedlist.
func (l *CacheImpl[K, V]) Get(key K) (V, error) {
	if link, ok := l.keyToElement[key]; ok {
		l.removeFromList(link)
		l.removeFreqLevel(link)
		link.Value.freq++
		l.addToList(link)
		return link.Value.value, nil
	}
	return l.defaultValue, ErrKeyNotFound
}

// node: struct for storing cache elements.
type node[K comparable, V any] struct {
	key   K
	value V
	freq  int
}

// addToList: method to add new Element to cache if cache was empty.
// func (l *CacheImpl[K, V]) addToList(link *element[*node[K, V]], previousLinkOfElement *element[*node[K, V]]) {
func (l *CacheImpl[K, V]) addToList(link *linkedlist.Element[*node[K, V]]) {
	// Empty cache case.
	if l.listOfData.Len() == 0 {
		l.listOfData.Add(link)
	} else {
		l.addToNotEmptyList(link)
	}
	n := link.Value
	l.frequencyToRecentElement[n.freq] = link
	l.keyToElement[n.key] = link
}

// addToNotEmptyList: method to add new Element to not empty cache.
func (l *CacheImpl[K, V]) addToNotEmptyList(link *linkedlist.Element[*node[K, V]]) {
	n := link.Value
	lastFreqRoot, freqExists := l.frequencyToRecentElement[n.freq]
	if freqExists {
		// Case when we just need to add Element before the most recent used Element of such frequency.
		// pic: (...)_n+1 - (+, +)_n - (x, ...)_n-1 --> (...)_n+1 - (x, +, +) - (...)_n-1
		l.listOfData.AddBefore(link, lastFreqRoot)
	} else if n.freq < l.listOfData.Back().Value.freq {
		// Case when we add Element of minimal frequency. We need to add it in the end of linkedlist.
		// pic: (...)_n+1 - (...)_n  --> (...)_n+1 - (...)_n - (x)_m, m < n
		l.listOfData.Add(link)
	} else if n.freq != 1 {

		// Case when last frequency level might exist.
		lastPastFreqRoot, pastFreqExists := l.frequencyToRecentElement[n.freq-1]
		if pastFreqExists {
			// Case when frequency - 1 level exists. It means that we need to add our Element before
			// the most recent used Element of frequency - 1 level.
			// pic: (...)_n+2 - (x, +)_n --> (...)_n+2 - (x)_n+1 - (+)_n
			l.listOfData.AddBefore(link, lastPastFreqRoot)
		} else {
			// Case when frequency - 1 level does not exist. It means, that Element to add in linkedlist was the only one
			// on frequency - 1 level before frequency increase. We can replace old Element with new element.
			// pic: (+)_n+1 - (x)_n-1 - (+)_n-2 --> (+)_n+1 - (x)_n - (+)_n-2
			l.listOfData.ReplaceDeletedElement(link, link)
		}
	}
}

// removeFromList removes Element from linkedlist and removes frequency level if Element was only one on it.
func (l *CacheImpl[K, V]) removeFromList(link *linkedlist.Element[*node[K, V]]) {
	l.listOfData.Remove(link)
	l.removeFreqLevel(link)
}

// removeFreqLevel removes frequency level if argument was single on it or makes other element
// of this frequency level representative.
func (l *CacheImpl[K, V]) removeFreqLevel(link *linkedlist.Element[*node[K, V]]) {
	if link == nil || link.GetNext() == nil {
		return
	}
	n := link.Value
	previousEl := link.GetNext().Value
	if l.frequencyToRecentElement[n.freq] == link {
		if previousEl != nil && previousEl.freq == n.freq {
			l.frequencyToRecentElement[n.freq] = link.GetNext()
		} else {
			delete(l.frequencyToRecentElement, n.freq)
		}
	}
}

// extractLatest extracts the least recently used Element of all least frequently used elements.
// If this Element was frequency level least frequently used, removes this level.
func (l *CacheImpl[K, V]) extractLatest() {
	del := l.listOfData.PopBack()
	n := del.Value
	if l.frequencyToRecentElement[n.freq] == del {
		delete(l.frequencyToRecentElement, n.freq)
	}
	delete(l.keyToElement, n.key)
}

// Put puts new node to cache.
func (l *CacheImpl[K, V]) Put(key K, value V) {
	if link, ok := l.keyToElement[key]; ok {
		// Case when cache contains Element with such key. Put removes this Element from linkedlist
		// and adds it with new frequency.
		l.removeFromList(link)
		l.removeFreqLevel(link)
		link.Value.value = value
		link.Value.freq++
		l.addToList(link)
		return
	}

	if l.Size() == l.capacity {
		// Case when adding occurs to a full cache. Extract latest Element and then add new.
		l.extractLatest()
	}

	// New cache Element(node) creation.
	n := &node[K, V]{
		key:   key,
		value: value,
		freq:  1,
	}

	l.keyToElement[key] = linkedlist.NewElement(n)
	// Addition to the linkedlist of elements.
	l.addToList(l.keyToElement[key])
}

// All returns iterator function witch iterates all cache elements in the order they are stored,
// yielding each key-value pair to the provided yield function.
func (l *CacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for element := range l.listOfData.Iterator() {
			if !(yield(element.key, element.value)) {
				return
			}
		}
	}
}

// Size returns the cache size (value of elements in linkedlist).
func (l *CacheImpl[K, V]) Size() int {
	return l.listOfData.Len()
}

// Capacity returns the cache capacity.
func (l *CacheImpl[K, V]) Capacity() int {
	return l.capacity
}

// GetKeyFrequency returns the frequency of given key Element if such key exists.
func (l *CacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	if link, ok := l.keyToElement[key]; ok {
		n := link.Value
		return n.freq, nil
	}
	return 0, ErrKeyNotFound
}
