package core

const NUMBER_OF_SPOTS = 1000000

type IndexRing struct {
	Ring []uint32
	Size uint32
}

func NewIndexRing(initialFilledIndex uint32, numberOfIndex uint32) *IndexRing {

	indexRing := make([]uint32, numberOfIndex)

	for id := range indexRing {
		indexRing[id] = initialFilledIndex
	}

	return &IndexRing{
		Ring: indexRing,
		Size: NUMBER_OF_SPOTS,
	}
}

func (indexRing *IndexRing) GetIndexId(key byte) uint32 {
	hashedKey := Get32MurmurHash([]byte{key})

	spotId := hashedKey % indexRing.Size

	return spotId
}

func (indexRing *IndexRing) GetNextFilledIndex(spotId uint32) uint32 {
	return indexRing.Ring[spotId]
}
