package models

type Version struct {
	version uint32
}

func NewVersion() *Version {
	return &Version{
		version: 1,
	}
}

func (v *Version) Newer() {
	v.version++
}

func (v *Version) Get() uint32 {
	return v.version
}
