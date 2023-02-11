package main

import (
	"bytes"
	"image/color"
	"image/png"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// input byt: png bytes =>  output byt: png bytes
func AddTimeStamp(byt []byte, x, y, size int, color color.Color, label string) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(byt))
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(font, &truetype.Options{Size: float64(size)})

	w := img.Bounds().Size().X
	h := img.Bounds().Size().Y

	dc := gg.NewContext(w, h)
	dc.DrawImage(img, 0, 0)
	dc.SetFontFace(face)
	r, g, b, _ := color.RGBA()
	dc.SetRGB(
		float64(r)/float64(0xff),
		float64(g)/float64(0xff),
		float64(b)/float64(0xff),
	)
	dc.DrawStringAnchored(label, float64(x), float64(y), 0, 0.5)

	img = dc.Image()

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)

	return buf.Bytes(), nil
}
