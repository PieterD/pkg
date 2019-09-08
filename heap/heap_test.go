package heap

import (
	"reflect"
	"sort"
	"testing"
)

func TestRemove(t *testing.T) {
	h := sort.IntSlice{
		9, 4, 6, 1, 3, 67, 1, 43, 7, 3, 5, 2, 5,
	}
	expected := make(sort.IntSlice, len(h))
	copy(expected, h)
	sort.Sort(expected)
	orig := h
	Init(h)
	check(t, h)
	for i, want := range expected {
		if h[0] != want {
			t.Logf(" got: %v", h[0])
			t.Logf("want: %v", want)
			t.Fatalf("unordered heap min index %d", i)
		}
		check(t, h)
		Remove(h, 0)
		h = h[:len(h)-1]
		check(t, h)
	}
	reverse(orig)
	if !reflect.DeepEqual(orig, expected) {
		t.Logf(" got: %v", orig)
		t.Logf("want: %v", expected)
		t.Fatalf("expected sorted h after reversing")
	}
}

func TestAppended(t *testing.T) {
	var h sort.IntSlice
	data := sort.IntSlice{6, 3, 5, 9, 1, 3, 1, 65, 43, 7, 1, 9, 6, 2, 5, 21, 44, 7, 9, 11}
	for _, insert := range data {
		check(t, h)
		h = append(h, insert)
		Appended(h)
		check(t, h)
	}
}

func check(t *testing.T, h sort.IntSlice) {
	t.Helper()
	for index := 1; index <= h.Len(); index++ {
		parentIndex := index / 2
		leftIndex := index * 2
		rightIndex := index*2 + 1
		if inRange(h, parentIndex) && isLess(h, index, parentIndex) {
			t.Logf("heap: %v", h)
			t.Errorf("bad parent at index %d (value %d) for item at index %d (value %d)", parentIndex-1, h[parentIndex-1], index-1, h[index-1])
		}
		if inRange(h, leftIndex) && isLess(h, leftIndex, index) {
			t.Logf("heap: %v", h)
			t.Errorf("bad left at index %d (value %d) for item at index %d (value %d)", leftIndex-1, h[leftIndex-1], index-1, h[index-1])
		}
		if inRange(h, rightIndex) && isLess(h, rightIndex, index) {
			t.Logf("heap: %v", h)
			t.Errorf("bad right at index %d (value %d) for item at index %d (value %d)", rightIndex-1, h[rightIndex-1], index-1, h[index-1])
		}
	}
}

func reverse(h sort.IntSlice) {
	for i, j := 0, len(h)-1; i < j; {
		h[i], h[j] = h[j], h[i]
		i++
		j--
	}
}
