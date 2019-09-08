package heap

import "sort"

type Interface interface {
	sort.Interface
}

// Init constructs a heap out of the provided unordered interface.
func Init(h Interface) {
	sort.Sort(h)
}

// Fix fixes the heap after the item at the index was changed.
// It will return the new index for the given item,
// or -1 if the index was not in range.
func Fix(h Interface, index int) int {
	index++
	newIndex := fix(h, index)
	if newIndex == -1 {
		return -1
	}
	newIndex--
	return newIndex
}

// Update fixes the heap after the first element has been updated.
func Update(h Interface) {
	if !inRange(h, 1) {
		return
	}
	fix(h, 1)
}

// Appended fixes the heap after a new element has been appended.
func Appended(h Interface) {
	if !inRange(h, h.Len()) {
		return
	}
	fix(h, h.Len())
}

// Remove moves the element at index to index h.Len()-1.
// This way it can be easily removed from the end of a slice.
func Remove(h Interface, index int) {
	index++
	if index == h.Len() {
		return
	}
	if !inRange(h, index) {
		return
	}
	doSwap(h, index, h.Len())
	fix(snip{Interface: h}, index)
}

type snip struct {
	Interface
}

func (s snip) Len() int {
	return s.Interface.Len() - 1
}

// fix uses a 1-based index
func fix(h Interface, index int) int {
	parentIndex := index / 2
	if inRange(h, parentIndex) && isLess(h, index, parentIndex) {
		doSwap(h, index, parentIndex)
		return fix(h, parentIndex)
	}
	leftIndex := index * 2
	rightIndex := index*2 + 1
	returnIndex := -1
	if inRange(h, leftIndex) && isLess(h, leftIndex, index) {
		doSwap(h, index, leftIndex)
		returnIndex = fix(h, leftIndex)
	}
	if inRange(h, rightIndex) && isLess(h, rightIndex, index) {
		doSwap(h, index, rightIndex)
		ri2 := fix(h, rightIndex)
		if returnIndex == -1 {
			returnIndex = ri2
		}
	}
	return returnIndex
}

// isLess uses a 1-based index
func isLess(h Interface, i, j int) bool {
	return h.Less(i-1, j-1)
}

// inRange uses a 1-based index
func inRange(h Interface, i int) bool {
	if i <= 0 {
		return false
	}
	return i <= h.Len()
}

// doSwap uses a 1-based index
func doSwap(h Interface, i, j int) {
	h.Swap(i-1, j-1)
}
