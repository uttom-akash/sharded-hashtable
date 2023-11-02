package models

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
		validSlots:   core.NewBitSet16(),
		deletedSlots: core.NewBitSet16(),
	}
}

func (bucket *Bucket) Get(key *Key) *Slot {

	slot := bucket.find(key)

	return slot
}

func (bucket *Bucket) Put(key *Key, value *Value) *Slot {

	slot := bucket.find(key)

	if slot == nil && !bucket.hasEmptySlots() && !bucket.hasDeletedSlots() {
		return nil
	}

	if slot == nil {
		return bucket.insert(*key, *value)
	}

	return bucket.update(*key, *value, slot)
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

func (bucket *Bucket) hasEmptySlots() bool {
	return bucket.validSlots.GetSetBitCount() <= 14
}

func (bucket *Bucket) hasDeletedSlots() bool {
	return bucket.deletedSlots.GetSetBitCount() > 0
}

func (bucket *Bucket) find(key *Key) *Slot {
	for id, tag := range bucket.tags {
		if tag == nil {
			continue
		}
		
		if tag.Tag == key.Key &&
			bucket.validSlots.IsSet(uint8(id)) &&
			!bucket.deletedSlots.IsSet(uint8(id)) &&
			bucket.slots[id].key.Key == key.Key {
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

func (bucket *Bucket) locateInsertableSlot() *Slot {
	slotId := bucket.deletedSlots.GetSetBitIndex()

	if slotId == -1 {
		slotId = bucket.validSlots.GetUnsetBitIndex()
	}

	bucket.slots[slotId] = NewEmptySlot(uint8(slotId))

	return bucket.slots[slotId]
}

func (bucket *Bucket) update(key Key, value Value, oldSlot *Slot) *Slot {
	reserveSlotId := bucket.validSlots.GetUnsetBitIndex()

	//step 1
	updatedSlot := NewSlot(&key, &value, uint8(reserveSlotId))
	bucket.slots[reserveSlotId] = updatedSlot
	bucket.tags[reserveSlotId] = key.Tag

	//step 2
	bucket.validSlots.Set(updatedSlot.Id)
	bucket.deletedSlots.Unset(updatedSlot.Id)

	bucket.validSlots.Unset(oldSlot.Id)
	bucket.slots[oldSlot.Id] = nil
	bucket.tags[oldSlot.Id] = nil

	bucket.version.Newer()

	return updatedSlot
}
