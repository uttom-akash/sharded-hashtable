package storageengine

import (
	"scale.kv.store/internal/core"
)

type Bucket struct {
	version      *Version
	validSlots   *core.BitSet16
	deletedSlots *core.BitSet16
	tags         [16]*Tag
	slots        [16]*Slot
}

func NewBucket() *Bucket {
	return &Bucket{
		version:      NewVersion(),
		validSlots:   core.NewBitSet16(), // set false to make this slot empty and search path end here.
		deletedSlots: core.NewBitSet16(), // set while deleting slot so that search path doesn't end here but this slot is not readable
	}
}

func (bucket *Bucket) Get(key *Key) *Slot {

	slot := bucket.find(key)

	return slot
}

func (bucket *Bucket) Put(key *Key, value *Value) *Slot {

	slot := bucket.find(key)

	if slot != nil {
		return bucket.update(*key, *value, slot)
	}

	return bucket.PutNewKey(key, value)
}

func (bucket *Bucket) PutNewKey(key *Key, value *Value) *Slot {

	if !bucket.anySlotToInsert() {
		return nil
	}

	return bucket.insert(*key, *value)
}

func (bucket *Bucket) Delete(key *Key) bool {

	slot := bucket.find(key)
	if slot == nil {
		return false
	}
	//step1
	bucket.deletedSlots.Set(slot.Id)
	bucket.version.Newer()

	return true
}

// func (bucket *Bucket) findP(key *Key) *Slot {
// 	foundChannel := make(chan int)
// 	// defer close(foundChannel)

// 	for index := range bucket.tags {
// 		id := index
// 		go func() {
// 			tag := bucket.tags[id]

// 			isMatch := tag != nil &&
// 				tag.IsEqual(key.Tag) &&
// 				bucket.isSlotReadable(uint8(id)) &&
// 				bucket.slots[id].key.IsEqual(key)

// 			if isMatch {
// 				foundChannel <- id
// 			} else {
// 				foundChannel <- -1
// 			}
// 		}()
// 	}

// 	index := 0

// 	for id := range foundChannel {
// 		if id != -1 {
// 			return bucket.slots[id]
// 		}

// 		if index == len(bucket.tags)-1 {
// 			break
// 		}
// 		index++
// 	}

// 	return nil
// }

func (bucket *Bucket) find(key *Key) *Slot {

	for id, tag := range bucket.tags {

		if tag == nil {
			continue
		}

		isMatch := bucket.isSlotVisible(uint8(id)) &&
			tag.IsEqual(key.Tag) &&
			bucket.slots[id].key.IsEqual(key)

		if isMatch {
			return bucket.slots[id]
		}
	}

	return nil
}

func (bucket *Bucket) insert(key Key, value Value) *Slot {
	slot := bucket.locateInsertableSlot()

	//step 1
	slot.key = &key
	slot.value = &value
	bucket.tags[slot.Id] = slot.key.Tag

	//step 2
	bucket.validSlots.Set(slot.Id)
	bucket.deletedSlots.Unset(slot.Id)
	bucket.version.Newer()

	return slot
}

func (bucket *Bucket) update(key Key, value Value, oldSlot *Slot) *Slot {

	reserveSlotId := bucket.validSlots.GetUnsetBitIndex()

	//step 1
	updatedSlot := NewSlot(&key, &value, uint8(reserveSlotId))
	bucket.slots[reserveSlotId] = updatedSlot
	bucket.tags[reserveSlotId] = key.Tag

	//step 2
	// make new slot visible
	bucket.validSlots.Set(updatedSlot.Id)
	bucket.deletedSlots.Unset(updatedSlot.Id)

	// make old slot invisible
	bucket.validSlots.Unset(oldSlot.Id)
	bucket.slots[oldSlot.Id] = nil
	bucket.tags[oldSlot.Id] = nil

	bucket.version.Newer()

	return updatedSlot
}

func (bucket *Bucket) locateInsertableSlot() *Slot {

	// first try to fill deleted slot
	// then try to fill empty slot

	slotId := bucket.deletedSlots.GetSetBitIndex()

	if slotId == -1 {
		slotId = bucket.validSlots.GetUnsetBitIndex()
	}

	bucket.slots[slotId] = NewEmptySlot(uint8(slotId))

	return bucket.slots[slotId]
}

func (bucket *Bucket) isSlotVisible(id uint8) bool {
	// slot is visible if it is valid or not deleted

	return bucket.validSlots.IsSet(id) &&
		!bucket.deletedSlots.IsSet(id)
}

func (bucket *Bucket) HasEmptySlots() bool {
	// total slots 16, 1 slot should be reserved

	return bucket.validSlots.GetSetBitCount() < 15
}

func (bucket *Bucket) HasDeletedSlots() bool {

	return bucket.deletedSlots.GetSetBitCount() > 0
}

func (bucket *Bucket) anySlotToInsert() bool {

	// there is empty slot or deleted slot
	return bucket.HasEmptySlots() || bucket.HasDeletedSlots()
}
