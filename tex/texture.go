package tex

import (
	"image/png"
	"os"
)

type Texture struct {
	Pixels []byte
	W, H   int
	Pitch  int
}

func NewTexture(w, h int) *Texture {
	tex := new(Texture)
	tex.Pixels = make([]byte, w*h*4)
	tex.W = w
	tex.H = h
	tex.Pitch = w * 4
	return tex
}

func LoadFromFile(fileName string) *Texture {
	inFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	img, err := png.Decode(inFile)
	if err != nil {
		panic(err)
	}
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	tex := NewTexture(w, h)
	bIndex := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			tex.Pixels[bIndex] = byte(r / 256)
			bIndex++
			tex.Pixels[bIndex] = byte(g / 256)
			bIndex++
			tex.Pixels[bIndex] = byte(b / 256)
			bIndex++
			tex.Pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}
	return tex
}
