package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	gocolor "github.com/gerow/go-color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func main() {
	path := "testdata"
	pathOut := "converted"
	files, _ := ioutil.ReadDir(path)

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, f := range files {
		go func(f os.FileInfo) {
			defer wg.Done()
			reader, err := os.Open(path + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			m, _, err := image.Decode(reader)
			if err != nil {
				log.Fatal(err)
			}

			bounds := m.Bounds()
			img := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := m.At(x, y).RGBA()

					rf := float64(r) / 65535
					gf := float64(g) / 65535
					bf := float64(b) / 65535

					hsl := gocolor.RGB{rf, gf, bf}.ToHSL()

					newLum := (1.0 - hsl.L) * 0.95

					newRgb := gocolor.HSL{hsl.H, hsl.S, newLum}.ToRGB()

					rn := uint8(newRgb.R * 255)
					gn := uint8(newRgb.G * 255)
					bn := uint8(newRgb.B * 255)

					img.Set(x, y, color.RGBA{rn, gn, bn, uint8(a / 256)})
				}
			}

			fout, err := os.Create(pathOut + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}

			if err := png.Encode(fout, img); err != nil {
				fout.Close()
				log.Fatal(err)
			}

			if err := fout.Close(); err != nil {
				log.Fatal(err)
			}

			fmt.Println("All done with ", f.Name())
		}(f)
	}

	wg.Wait()


}
