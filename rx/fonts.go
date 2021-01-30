package rx

import (
	"io/ioutil"
	"sync"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

var (
	fontsOnce            sync.Once
	sourceSansProRegular *truetype.Font
	defaultFont          *truetype.Font
	defaultFontFace      font.Face
)

// SourceSansProRegular Truetype Font
func SourceSansProRegular() *truetype.Font {
	fontsOnce.Do(fontsSetup)
	return sourceSansProRegular
}

// DefaultFont Truetype Font
func DefaultFont() *truetype.Font {
	return SourceSansProRegular()
}

// DefaultFontFace is the 20pt 72dpi*DeviceScaleFactor() font face
//
// Warning: do not use this inside an init() function, specially on Android
func DefaultFontFace() font.Face {
	fontsOnce.Do(fontsSetup)
	return defaultFontFace
}

func fontsSetup() {
	f, err := res.Open("public/fonts/SourceSansPro-Regular.ttf")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	tt, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	sourceSansProRegular = tt
	defaultFont = tt
	defaultFontFace = truetype.NewFace(defaultFont, &truetype.Options{
		DPI:     72 * ebiten.DeviceScaleFactor(),
		Size:    20,
		Hinting: font.HintingFull,
	})
}
