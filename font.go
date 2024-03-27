package main

import (
	"image/color"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Font struct {
	face font.Face
}

func newFont(path string, size int) Font {
	font_data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	tt, err := opentype.Parse(font_data)
	if err != nil {
		log.Fatal(err)
	}
	font_face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	return Font{
		face: font_face,
	}
}

func (f *Font) drawFont(surface *ebiten.Image, str string, x int, y int, clr color.Color) {
	text.Draw(surface, str, f.face, x, y, clr)
}
