package aseprite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const smap = `{ "frames": {
	"player (person) 0.aseprite": {
	 "frame": { "x": 55, "y": 0, "w": 8, "h": 11 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 17, "y": 8, "w": 8, "h": 11 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (person) 1.aseprite": {
	 "frame": { "x": 55, "y": 0, "w": 8, "h": 11 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 17, "y": 8, "w": 8, "h": 11 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (person) 2.aseprite": {
	 "frame": { "x": 55, "y": 0, "w": 8, "h": 11 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 17, "y": 8, "w": 8, "h": 11 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (bumper) 0.aseprite": {
	 "frame": { "x": 46, "y": 0, "w": 9, "h": 24 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 27, "y": 4, "w": 9, "h": 24 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (bumper) 1.aseprite": {
	 "frame": { "x": 46, "y": 0, "w": 9, "h": 24 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 31, "y": 4, "w": 9, "h": 24 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (bumper) 2.aseprite": {
	 "frame": { "x": 46, "y": 0, "w": 9, "h": 24 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 29, "y": 4, "w": 9, "h": 24 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (ball_front) 0.aseprite": {
	 "frame": { "x": 0, "y": 0, "w": 32, "h": 32 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 32, "h": 32 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (ball_front) 1.aseprite": {
	 "frame": { "x": 0, "y": 0, "w": 32, "h": 32 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 32, "h": 32 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (ball_front) 2.aseprite": {
	 "frame": { "x": 0, "y": 0, "w": 32, "h": 32 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 32, "h": 32 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (glass) 0.aseprite": {
	 "frame": { "x": 32, "y": 0, "w": 14, "h": 17 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 15, "y": 0, "w": 14, "h": 17 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (glass) 1.aseprite": {
	 "frame": { "x": 32, "y": 0, "w": 14, "h": 17 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 15, "y": 0, "w": 14, "h": 17 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (glass) 2.aseprite": {
	 "frame": { "x": 32, "y": 0, "w": 14, "h": 17 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 15, "y": 0, "w": 14, "h": 17 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (propulsion) 0.aseprite": {
	 "frame": { "x": 32, "y": 17, "w": 10, "h": 12 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 11, "y": 13, "w": 10, "h": 12 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (propulsion) 1.aseprite": {
	 "frame": { "x": 32, "y": 17, "w": 10, "h": 12 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 11, "y": 13, "w": 10, "h": 12 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (propulsion) 2.aseprite": {
	 "frame": { "x": 32, "y": 17, "w": 10, "h": 12 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 11, "y": 13, "w": 10, "h": 12 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (plasma) 0.aseprite": {
	 "frame": { "x": 42, "y": 24, "w": 6, "h": 8 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 13, "y": 25, "w": 6, "h": 8 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (plasma) 1.aseprite": {
	 "frame": { "x": 55, "y": 11, "w": 6, "h": 9 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 13, "y": 25, "w": 6, "h": 9 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	},
	"player (plasma) 2.aseprite": {
	 "frame": { "x": 55, "y": 20, "w": 7, "h": 7 },
	 "rotated": false,
	 "trimmed": true,
	 "spriteSourceSize": { "x": 13, "y": 25, "w": 7, "h": 7 },
	 "sourceSize": { "w": 40, "h": 40 },
	 "duration": 100
	}
  },
  "meta": {
   "app": "http://www.aseprite.org/",
   "version": "1.2.18",
   "image": "player.png",
   "format": "RGBA8888",
   "size": { "w": 63, "h": 32 },
   "scale": "1",
   "frameTags": [
   ],
   "layers": [
	{ "name": "person", "opacity": 255, "blendMode": "normal" },
	{ "name": "bumper", "opacity": 255, "blendMode": "normal" },
	{ "name": "ball_front", "opacity": 255, "blendMode": "normal" },
	{ "name": "glass", "opacity": 255, "blendMode": "normal" },
	{ "name": "propulsion", "opacity": 255, "blendMode": "normal" },
	{ "name": "plasma", "opacity": 255, "blendMode": "normal" }
   ],
   "slices": [
   ]
  }
 }`

const lmap = `{ "frames": [
	{
	 "filename": "bg1",
	 "frame": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg2",
	 "frame": { "x": 64, "y": 0, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg3",
	 "frame": { "x": 128, "y": 0, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg4",
	 "frame": { "x": 192, "y": 0, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg5",
	 "frame": { "x": 0, "y": 64, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg6",
	 "frame": { "x": 64, "y": 64, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg7",
	 "frame": { "x": 128, "y": 64, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg8",
	 "frame": { "x": 192, "y": 64, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg9",
	 "frame": { "x": 0, "y": 128, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg10",
	 "frame": { "x": 64, "y": 128, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg11",
	 "frame": { "x": 128, "y": 128, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	},
	{
	 "filename": "bg12",
	 "frame": { "x": 192, "y": 128, "w": 64, "h": 64 },
	 "rotated": false,
	 "trimmed": false,
	 "spriteSourceSize": { "x": 0, "y": 0, "w": 64, "h": 64 },
	 "sourceSize": { "w": 64, "h": 64 },
	 "duration": 100
	}
  ],
  "meta": {
   "app": "http://www.aseprite.org/",
   "version": "1.2.18",
   "image": "background.png",
   "format": "RGBA8888",
   "size": { "w": 256, "h": 192 },
   "scale": "1"
  }
 }
 `

func TestParse(t *testing.T) {
	x, err := Parse([]byte(smap))
	assert.NoError(t, err)
	assert.Equal(t, FileTypeMap, x.Type())
	assert.Equal(t, "1", x.GetMetadata().Scale)
	assert.Equal(t, "RGBA8888", x.GetMetadata().Format)

	x, err = Parse([]byte(lmap))
	assert.NoError(t, err)
	assert.Equal(t, FileTypeSlice, x.Type())
}
