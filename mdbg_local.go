package main

import "github.com/jcramb/cedict"

type MdgbLocal struct {
	d *cedict.Dict
}

func NewMdgbLocal() *MdgbLocal {
	return &MdgbLocal{
		d: cedict.New(),
	}
}

func (m *MdgbLocal) Get(word string) ResultObj {
	r := ResultObj{}

	entry := m.d.GetByHanzi(word)
	r.English = entry.Meanings
	r.Reading = []string{entry.Pinyin}

	return r
}
