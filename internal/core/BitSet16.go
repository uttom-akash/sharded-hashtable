package core

import "math/bits"

const LENGTH = 16

type BitSet16 struct {
	bitset uint16
}

func NewBitSet16() *BitSet16 {
	return &BitSet16{
		bitset: 0,
	}
}

func (b *BitSet16) Set(index uint8) {
	if index < LENGTH {
		b.bitset = b.bitset | (1 << index)
	}
}

func (b *BitSet16) Unset(index uint8) {
	if index < LENGTH {
		b.bitset = b.bitset & ^(1 << index)
	}
}

func (b *BitSet16) IsSet(index uint8) bool {
	if index < LENGTH {
		return b.bitset&(1<<index) == 1
	}
	return false
}

func (b *BitSet16) GetSetBitCount() uint8 {
	return uint8(bits.OnesCount16(b.bitset))
}

func (b *BitSet16) GetSetBitIndex() int8 {
	for index := int8(0); index < LENGTH; index++ {
		if b.bitset&(1<<index) != 0 {
			return index
		}
	}
	return int8(-1)
}

func (b *BitSet16) GetUnsetBitIndex() int8 {
	for index := int8(0); index < LENGTH; index++ {
		if b.bitset&(1<<index) == 0 {
			return index
		}
	}
	return -1
}
