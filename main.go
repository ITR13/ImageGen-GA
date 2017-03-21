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
	"math"
	"strconv"

	"github.com/MaxHalford/gago"
)

var W, H, R int

func limit(v, max float64) float64 {
	for v < -max {
		v += max * 2
	}
	for v >= max*2 {
		v -= max * 2
	}
	if v > max {
		v = 2*max - v
	} else if v < 0 {
		v = -v
	}
	return v
}

func fit(X []float64, apply func([]float64, *image.RGBA, int),
	org *image.RGBA, diff *image.RGBA, max float64) float64 {
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	if org != nil {
		clone(org, img)
	}
	if X != nil && apply != nil {
		apply(X, img, 0)
	}
	sum := float64(0)
	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {
			c := img.RGBAAt(x, y)
			c2 := diff.RGBAAt(x, y)
			r, g, b := float64(c.R)-float64(c2.R), float64(c.G)-
				float64(c2.G), float64(c.B)-float64(c2.B)

			//Alternative way of doing it
			//sum += math.Sqrt(r*r + g*g + b*b)
			sum += math.Abs(r) + math.Abs(g) + math.Abs(b)
		}
	}
	return sum * 100 / max
}

func getFit(apply func([]float64, *image.RGBA, int), org, diff *image.RGBA, max float64) func([]float64) float64 {
	return func(X []float64) float64 {
		return fit(X, apply, org, diff, max)
	}
}

func main() {
	totalIncreases := 0

	org, err := load("./input.png")
	if err != nil {
		panic(err)
	}
	bounds := org.Bounds()
	W, H = bounds.Dx(), bounds.Dy()
	if W > H {
		R = W * 3 / 2
	} else {
		R = H * 3 / 2
	}

	previousFitness := math.Inf(1)
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	max := fit(nil, nil, img, org, 100)
	fmt.Println(max)

	for j := 0; previousFitness >= 0.001; j++ {
		usingCircle, ga := useCircle(img, org, max)
		if usingCircle {
			fmt.Println("Using Circle")
		} else {
			fmt.Println("Using Bezier")
		}

		ga.Initialize()
		fitness := ga.Best.Fitness
		lastFitness := fitness
		count := 0
		i := 0
		increases := 0

		for maxCount := 250; fitness > previousFitness || i < 200; maxCount += 250 {
			if maxCount != 250 {
				ga.Initialize()
				fmt.Printf("Increased MaxCount to %d\n", maxCount)
			}
			for count < maxCount {
				ga.Enhance()
				if fitness == ga.Best.Fitness {
					count++
				} else {
					fitness = ga.Best.Fitness
					count = 0
					totalIncreases++
					increases++
				}
				if count%10 == 0 {
					fmt.Printf("%d:\t%d:\t%d\t%d\tThe best obtained solution is %f%%\n",
						i, count, totalIncreases, j, ga.Best.Fitness)
				}
				i++

				if lastFitness-fitness >= 0.005 {
					lastFitness = fitness
					genome := ga.Best.Genome
					var casted = make([]float64, len(genome))
					for i := range genome {
						if genome[i] == nil {
							fmt.Printf("genome #%d is nil, not float\n", i)
						}
						casted[i] = genome[i].(float64)
					}

					img2 := image.NewRGBA(image.Rect(0, 0, W, H))
					clone(img, img2)
					if usingCircle {
						ApplyCircle(casted, img2, 0)
					} else {
						ApplyBezier(casted, img2, 0)
					}
					go save("./output/gen"+strconv.Itoa(j)+"-"+
						strconv.Itoa(maxCount/250)+"-"+strconv.Itoa(increases)+
						".png", img2)
				}
			}
		}

		genome := ga.Best.Genome
		var casted = make([]float64, len(genome))
		for i := range genome {
			casted[i] = genome[i].(float64)
		}

		if usingCircle {
			ApplyCircle(casted, img, 0)
		} else {
			ApplyBezier(casted, img, 0)
		}

		go save("./output/gen"+strconv.Itoa(j+1)+"-0-0.png", img)
		fmt.Printf("Went from %f%% to %f%%\n", previousFitness, fitness)
		previousFitness = fitness
	}

}

func useCircle(img, org *image.RGBA, max float64) (bool, gago.GA) {
	bezier := BezierGa(img, org, max)
	bezier.Initialize()
	circle := CircleGa(img, org, max)
	circle.Initialize()
	for i := 0; i < 40; i++ {
		bezier.Enhance()
		circle.Enhance()
	}
	fmt.Printf("%f, %f\n", bezier.Best.Fitness, circle.Best.Fitness)
	if bezier.Best.Fitness <= circle.Best.Fitness {
		return false, bezier
	}
	return true, circle
}
