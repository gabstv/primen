package spr

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type SpriteDefVersion string

const (
	SpriteV1 SpriteDefVersion = "1"
)

type SpriteDef struct {
	Version  string          `json:"version"`
	Source   *SourceData     `json:"source,omitempty"`
	Position Vec2            `json:"position"`
	Size     Vec2            `json:"size"`
	Origin   Vec2            `json:"origin"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
	filename string          `json:"-"`
}

func (sd *SpriteDef) Filename() string {
	return sd.filename
}

func NewSpriteDef(name string, meta interface{}) *SpriteDef {
	d := &SpriteDef{
		filename: name,
		Version:  "1",
	}
	if meta == nil {
		return d
	}
	if db, ok := meta.([]byte); ok {
		d.Metadata = json.RawMessage(db)
	}
	if db, ok := meta.(string); ok {
		d.Metadata = json.RawMessage([]byte(db))
	}
	if bb, err := json.Marshal(meta); err == nil {
		d.Metadata = json.RawMessage(bb)
	}
	return d
}

func ReadSpriteDefFile(name string) (*SpriteDef, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := &SpriteDef{}
	if err = json.NewDecoder(f).Decode(d); err != nil {
		return nil, err
	}
	d.filename = name
	return d, nil
}

func (sd *SpriteDef) WriteToFile(name string, perm os.FileMode) error {
	data, err := json.Marshal(sd)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, data, perm)
}
