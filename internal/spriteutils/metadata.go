package spriteutils

import (
	"encoding/json"
)

type SpriteDefMeta struct {
	RawOriginX   int64   `json:"raw_origin_x"`
	RawOriginY   int64   `json:"raw_origin_y"`
	ExtraOriginX float64 `json:"extra_origin_x"`
	ExtraOriginY float64 `json:"extra_origin_y"`
}

func GetSpriteDef(rawjson []byte) (*SpriteDefMeta, error) {
	d := &SpriteDefMeta{}
	if err := json.Unmarshal(rawjson, d); err != nil {
		return nil, err
	}
	return d, nil
}

func MayGetSpriteDef(rawjson []byte) *SpriteDefMeta {
	d, err := GetSpriteDef(rawjson)
	if err != nil {
		return &SpriteDefMeta{}
	}
	return d
}

func (m *SpriteDefMeta) JSON() ([]byte, error) {
	return json.Marshal(m)
}

func (s *SpriteDefMeta) MustJSON() []byte {
	d, err := s.JSON()
	if err != nil {
		panic(err)
	}
	return d
}
