package core

import (
	"github.com/spaolacci/murmur3"
)

func Get32MurmurHash(data []byte) uint32 {
	seed := uint32(42)

	hash := murmur3.Sum32WithSeed(data, seed)

	return hash
}

func Get16MurmurHash(data []byte) uint16 {
	hash32 := Get32MurmurHash(data)

	hash16 := uint16(hash32 & 0xffff)

	return hash16
}

func Get8MurmurHash(data []byte) uint8 {
	hash32 := Get32MurmurHash(data)

	hash8 := uint8(hash32 & 0xff)

	return hash8
}
