package models

type Coordinator struct {
	IndexRing *IndexRing
}

func NewCoordinator() *Coordinator {

	return &Coordinator{
		IndexRing: NewIndexRing(),
	}
}

func (coordinator *Coordinator) Get(key byte) *Result {

	spotId := coordinator.IndexRing.getSpotId(key)

	spot := coordinator.IndexRing.getNextSpot(spotId)

	var result *Result

	for _, shard := range spot.Shards {

		result := shard.Get(key)

		if result.Status == Found || result.Status == StopSearch {
			return result
		}
	}

	return result
}

func (coordinator *Coordinator) Put(key byte, value byte) *Result {

	spotId := coordinator.IndexRing.getSpotId(key)

	spot := coordinator.IndexRing.getNextSpot(spotId)

	spot.writeLock.Lock()
	defer spot.writeLock.Unlock()

	var bucketId *int
	var shardd *Shard
	var result *Result

	for _, shard := range spot.Shards {

		result := shard.SearchForWrite(key)

		if result.Status == Found || result.Status == StopSearch {

			result = shard.Put(key, value, *result.BucketId)

			return result

		} else {
			shardd = shard
			bucketId = result.BucketId
		}
	}

	if bucketId != nil && shardd != nil {
		result = shardd.Put(key, value, *bucketId)
	}

	return result
}

func (coordinator *Coordinator) Delete(key byte) *Result {

	spotId := coordinator.IndexRing.getSpotId(key)

	spot := coordinator.IndexRing.getNextSpot(spotId)

	var result *Result

	for _, shard := range spot.Shards {

		result := shard.Delete(key)

		if result.Status == Deleted || result.Status == StopSearch {
			return result
		}
	}

	return result
}
