package core

type BitSet struct {
	bits []bool
}

func NewBitSet(size int) *BitSet {
	return &BitSet{
		bits: make([]bool, size),
	}
}

func (b *BitSet) Set(index int) {
	if index < len(b.bits) {
		b.bits[index] = true
	}
}

func (b *BitSet) Clear(index int) {
	if index < len(b.bits) {
		b.bits[index] = false
	}
}

func (b *BitSet) Get(index int) bool {
	if index < len(b.bits) {
		return b.bits[index]
	}
	return false
}
