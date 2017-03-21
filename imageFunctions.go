/*
    This file is part of ImageGenGA.

    ImageGenGA is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    ImageGenGA is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with ImageGenGA.  If not, see <http://www.gnu.org/licenses/>.
*/
	
package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func clone(src, dst *image.RGBA) {
	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {
			dst.SetRGBA(x, y, src.RGBAAt(x, y))
		}
	}
}

func load(path string) (*image.RGBA, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	orgImg, err := png.Decode(in)
	if err != nil {
		return nil, err
	}
	in.Close()
	org := image.NewRGBA(orgImg.Bounds())
	bounds := orgImg.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			org.Set(x, y, orgImg.At(x, y))
		}
	}
	return org, nil
}

func save(path string, img *image.RGBA) {
	out, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
	}
	out.Close()

	stored, err := load(path)
	if isEqual(img, stored) {
		fmt.Println("The image was stored without any changes")
		return
	}
	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {
			c1 := img.RGBAAt(x, y)
			c2 := stored.RGBAAt(x, y)
			if c1 != c2 {
				if c1.R > 250 {
					c1.R -= 2
				} else if c1.R < 5 {
					c1.R += 2
				}
				if c1.G > 250 {
					c1.G -= 2
				} else if c1.G < 5 {
					c1.G += 2
				}
				if c1.B > 250 {
					c1.B -= 2
				} else if c1.B < 5 {
					c1.B += 2
				}
				stored.SetRGBA(x, y, c2)
			}
		}
	}

	out, err = os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	err = png.Encode(out, stored)
	if err != nil {
		fmt.Println(err)
	}
	out.Close()

	fmt.Println("Modifed picture and saved")
}

func isEqual(a, b *image.RGBA) bool {
	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {
			c1 := a.RGBAAt(x, y)
			c2 := b.RGBAAt(x, y)
			if c1 != c2 {
				return false
			}
		}
	}
	return true
}
