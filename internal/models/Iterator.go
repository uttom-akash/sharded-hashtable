package models

type Iterator struct {
	startIndex   int
	currentIndex int
	size         int
}

func NewIterator(startIndex int, size int) *Iterator
{
	return &Iterator{
		startIndex: startIndex,
		currentIndex: -1,
		size: size,
	}
}

func (iterator *Iterator) Next() int {

	if iterator.currentIndex == (iterator.startIndex - 1+iterator.size)%iterator.size{//current index is at pervious index of start index
		return -1;
	}

	iterator.currentIndex= (iterator.currentIndex + 1)%iterator.size

	return iterator.currentIndex
}
