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

//Taken from: https://rosettacode.org/wiki/Bitmap/B%C3%A9zier_curves/Cubic#Go
func drawBézier3(x1, y1, x2, y2, x3, y3, x4, y4, b3Seg int,
	c color.RGBA, img *image.RGBA) {
	b3Segf := float64(b3Seg)

	px, py := make([]int, b3Seg+1), make([]int, b3Seg+1)
	fx1, fy1 := float64(x1), float64(y1)
	fx2, fy2 := float64(x2), float64(y2)
	fx3, fy3 := float64(x3), float64(y3)
	fx4, fy4 := float64(x4), float64(y4)
	for i := range px {
		d := float64(i) / b3Segf
		a := 1 - d
		b, c := a*a, d*d
		a, b, c, d = a*b, 3*b*d, 3*a*c, c*d
		px[i] = int(a*fx1 + b*fx2 + c*fx3 + d*fx4)
		py[i] = int(a*fy1 + b*fy2 + c*fy3 + d*fy4)
	}
	x0, y0 := px[0], py[0]
	for i := 1; i <= b3Seg; i++ {
		x1, y1 := px[i], py[i]
		drawLine(x0, y0, x1, y1, c, img, 0)
		x0, y0 = x1, y1
	}
}
