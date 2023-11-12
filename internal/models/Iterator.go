package models

type Iterator struct {
	startIndex   int
	currentIndex int
	size         int
}

func NewIterator(startIndex int, size int) *Iterator {
	startIndex = startIndex % size

	return &Iterator{
		startIndex:   startIndex,
		currentIndex: startIndex,
		size:         size,
	}
}

func (iterator *Iterator) Next() int {

	nextIndex := (iterator.currentIndex + 1) % iterator.size

	if nextIndex == iterator.startIndex { //current index is at pervious index of start index
		return -1
	}

	iterator.currentIndex = nextIndex

	return iterator.currentIndex
}

func (iterator *Iterator) Current() int {
	return iterator.currentIndex
}
