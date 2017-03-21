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
	"image/color"
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
)

func ApplyCircle(X []float64, img *image.RGBA, verbose int) {
	for i := 0; i+5 < len(X); i += 6 {
		x := X[i+0] * float64(W)
		y := X[i+1] * float64(H)
		r := X[i+2]*float64(R) + 10
		red := X[i+3] * 255
		green := X[i+4] * 255
		blue := X[i+5] * 255
		c := color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
		if verbose > 0 {
			fmt.Printf("%f, %f, %f, %f, %f, %f\n", x, y, r, red, green, blue)
			fmt.Println(c)
		}
		drawCircle(int(x), int(y), int(r),
			c, img, verbose-1)
	}
}

func CircleGa(org, diff *image.RGBA, max float64) gago.GA {
	return gago.GA{
		NbrPopulations: 6,
		NbrIndividuals: 80,
		NbrGenes:       6,
		Ff: gago.Float64Function{
			Image: getFit(ApplyCircle, org, diff, max),
		},
		Initializer: gago.InitUniformF{
			Lower: 1,
			Upper: 1,
		},
		Model: gago.ModSteadyState{
			Selector: gago.SelTournament{
				NbParticipants: 3,
			},
			Crossover: CircleCross{},
			KeepBest:  true,
			Mutator: MutClamped{
				Rate: 0.25,
				Div:  3,
			},
			MutRate: 0.5,
		},
		Migrator:     gago.MigShuffle{},
		MigFrequency: 10,
	}
}

type CircleCross struct{}

func (cross CircleCross) Apply(p1 gago.Individual, p2 gago.Individual,
	rng *rand.Rand) (gago.Individual, gago.Individual) {
	var (
		o1 = gago.Individual{
			Genome:    make([]interface{}, 6),
			Fitness:   math.Inf(1),
			Evaluated: false,
			Name:      "-",
		}
		o2 = gago.Individual{
			Genome:    make([]interface{}, 6),
			Fitness:   math.Inf(1),
			Evaluated: false,
			Name:      "-",
		}
	)
	p := rng.Float64()
	for i := 0; i < 2; i++ {
		o1.Genome[i] = p*p1.Genome[i].(float64) + (1-p)*p2.Genome[i].(float64)
		o2.Genome[i] = (1-p)*p1.Genome[i].(float64) + p*p2.Genome[i].(float64)
	}
	p = rng.Float64()
	o1.Genome[2] = p*p1.Genome[2].(float64) + (1-p)*p2.Genome[2].(float64)
	o2.Genome[2] = (1-p)*p1.Genome[2].(float64) + p*p2.Genome[2].(float64)
	p = rng.Float64()
	for i := 3; i < 6; i++ {
		o1.Genome[i] = p*p1.Genome[i].(float64) + (1-p)*p2.Genome[i].(float64)
		o2.Genome[i] = (1-p)*p1.Genome[i].(float64) + p*p2.Genome[i].(float64)
	}
	return o1, o2
}

type MutClamped struct {
	Rate float64
	Div  float64
}

func (mut MutClamped) Apply(indi *gago.Individual, rng *rand.Rand) {
	for i := range indi.Genome {
		if rng.Float64() <= mut.Rate {
			r := rng.Float64()
			r -= 0.5
			r /= 2
			indi.Genome[i] = limit(indi.Genome[i].(float64)+r, 1)
		}
	}
}
