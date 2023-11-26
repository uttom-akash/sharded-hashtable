package coordinationservice

import (
	"scale.kv.store/internal/core"
	"scale.kv.store/internal/metadata"
	"scale.kv.store/internal/storageengine"
)

const NUMBER_OF_SPOTS = 1000000

type Coordinator struct {
	IndexRing                 *core.IndexRing
	ShardGroupMetadataManager *ShardGroupMetadataManager
	ShardManager              *ShardManager
}

func NewCoordinator() *Coordinator {

	return &Coordinator{
		IndexRing:                 core.NewIndexRing(0, NUMBER_OF_SPOTS),
		ShardManager:              NewShardManager(),
		ShardGroupMetadataManager: NewShardGroupMetadataManager(),
	}
}

func (coordinator *Coordinator) Get(key byte) *storageengine.Result {

	shardGroup := coordinator.getShardGroup(key)

	var result *storageengine.Result

	for _, shardMetadata := range shardGroup.Shards {

		// replacce with rpc
		shard := coordinator.ShardManager.GetShard(shardMetadata.ShardId)

		result := shard.Get(key)
		//end replacce with rpc

		if result.Status == storageengine.Found || result.Status == storageengine.StopSearch {
			return result
		}
	}

	return result
}

func (coordinator *Coordinator) Put(key byte, value byte) *storageengine.Result {

	shardGroup := coordinator.getShardGroup(key)

	shardGroup.WriteLock.Lock()
	defer shardGroup.WriteLock.Unlock()

	var bucketId *int
	var shardd *storageengine.Shard
	var result *storageengine.Result

	for _, shardMetadata := range shardGroup.Shards {

		// replacce with rpc
		shard := coordinator.ShardManager.GetShard(shardMetadata.ShardId)

		result := shard.SearchForWrite(key)
		//end replacce with rpc

		if result.Status == storageengine.Found || result.Status == storageengine.StopSearch {

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

func (coordinator *Coordinator) Delete(key byte) *storageengine.Result {

	shardGroup := coordinator.getShardGroup(key)

	var result *storageengine.Result

	for _, shardMetadata := range shardGroup.Shards {

		// replacce with rpc
		shard := coordinator.ShardManager.GetShard(shardMetadata.ShardId)

		result := shard.Delete(key)
		//end replacce with rpc

		if result.Status == storageengine.Deleted || result.Status == storageengine.StopSearch {
			return result
		}
	}

	return result
}

func (coordinator *Coordinator) getShardGroup(key byte) *metadata.ShardGroupMetadata {
	spotId := coordinator.IndexRing.GetIndexId(key)

	filledSpotId := coordinator.IndexRing.GetNextFilledIndex(spotId)

	shardGroup := coordinator.ShardGroupMetadataManager.GetShardGroup(filledSpotId)

	return shardGroup
}
