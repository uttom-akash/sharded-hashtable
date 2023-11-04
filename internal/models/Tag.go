package models

import "scale.kv.store/internal/core"

type Tag struct {
	Tag byte
}

func NewTag(key byte) *Tag {
	return &Tag{
		Tag: core.Get8MurmurHash([]byte{key}),
	}
}

func (tag *Tag) IsEqual(pTag *Tag) bool {
	return tag.Tag == pTag.Tag
}
