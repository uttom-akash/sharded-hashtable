package models

type Tag struct {
	Tag byte
}

func NewTag(tag byte) *Tag {
	return &Tag{
		Tag: tag,
	}
}
