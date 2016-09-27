package main

import (
	"fmt"
	"image"
	"math"
	"strconv"
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
	gas := make([]float64, 0)

	org, err := load("./input.png")
	if err != nil {
		panic(err)
	}
	bounds := org.Bounds()
	W, H = bounds.Dx(), bounds.Dy()
	if W > H {
		R = W
	} else {
		R = H
	}

	max := fit(nil, nil, nil, org, 100)
	previousFitness := math.Inf(1)
	img := image.NewRGBA(image.Rect(0, 0, W, H))

	for j := 0; j < 500; j++ {
		ga := CircleGa(img, org, max)
		ga.Initialize()
		fitness := ga.Best.Fitness
		count := 0
		i := 0
		prevI := 0
		increases := 0
		prevIncrease := 0

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

				if i-prevI > 25 && prevIncrease < increases {
					prevI = i
					prevIncrease = increases
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
					applyCircle(casted, img2, 0)
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

		applyCircle(casted, img, 0)
		go save("./output/gen"+strconv.Itoa(j+1)+"-0-0.png", img)
		fmt.Printf("Went from %f%% to %f%%\n", previousFitness, fitness)
		gas = append(gas, casted[0], casted[1], casted[2], casted[3],
			casted[4], casted[5])
		previousFitness = fitness
	}

}
