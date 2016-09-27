package main

import (
	"fmt"
	"image"
	"image/color"
)

//taken from: https://rosettacode.org/wiki/Bitmap/Bresenham%27s_line_algorithm#Go
func drawLine(x0, y0, x1, y1 int, c color.RGBA, img *image.RGBA, verbose int) {
	if verbose > 1 {
		fmt.Println(c)
	}
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		img.SetRGBA(x0, y0, c)
		if verbose > 0 {
			if img.RGBAAt(x0, y0) != c {
				fmt.Printf("%v != %v", img.RGBAAt(x0, y0), c)
				//panic(";`;")
			}
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

//taken from: https://rosettacode.org/wiki/Bitmap/Midpoint_circle_algorithm#Go
func drawCircle(x, y, r int, c color.RGBA, img *image.RGBA, verbose int) {
	if verbose > 0 {
		fmt.Println(c)
	}
	if r < 0 {
		return
	}
	// Bresenham algorithm
	x1, y1, err := -r, 0, 2-2*r
	for {
		drawLine(x+x1, y+y1, x-x1, y+y1, c, img, verbose)
		drawLine(x-x1, y-y1, x+x1, y-y1, c, img, 0)
		drawLine(x+y1, y-x1, x-y1, y-x1, c, img, 0)
		drawLine(x-y1, y+x1, x+y1, y+x1, c, img, 0)
		//img.SetRGBA(x-x1, y+y1, c)
		//img.SetRGBA(x-y1, y-x1, c)
		//img.SetRGBA(x+x1, y-y1, c)
		//img.SetRGBA(x+y1, y+x1, c)
		r = err
		if r > x1 {
			x1++
			err += x1*2 + 1
		}
		if r <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

type safeImg struct {
}
