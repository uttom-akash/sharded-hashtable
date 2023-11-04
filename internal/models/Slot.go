package models

type Key struct {
	Key byte
	Tag *Tag
}

func NewKey(key byte) *Key {
	return &Key{
		Key: key,
		Tag: NewTag(key),
	}
}

func (key *Key) IsEqual(pKey *Key) bool {
	return key.Key == pKey.Key
}

type Value struct {
	Value byte
}

func NewValue(value byte) *Value {
	return &Value{Value: value}
}

type Slot struct {
	Id    uint8
	key   *Key
	value *Value
}

func NewSlot(key *Key, value *Value, id uint8) *Slot {
	return &Slot{
		id,
		key,
		value,
	}
}

func NewEmptySlot(id uint8) *Slot {
	return &Slot{
		id,
		&Key{Key: 0},
		&Value{Value: 0},
	}
}
