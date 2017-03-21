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

func ApplyBezier(X []float64, img *image.RGBA, verbose int) {
	for i := 0; i+11 < len(X); i += 12 {
		x1 := int(X[i+0] * float64(W))
		y1 := int(X[i+1] * float64(H))
		x2 := int(X[i+2] * float64(W))
		y2 := int(X[i+3] * float64(H))
		x3 := int(X[i+4] * float64(W))
		y3 := int(X[i+5] * float64(H))
		x4 := int(X[i+6] * float64(W))
		y4 := int(X[i+7] * float64(H))
		r := int(X[i+8]*20 + 1)
		b3Seg := int(X[i+9]*60 + 2)
		red := X[i+10] * 255
		green := X[i+11] * 255
		blue := X[i+12] * 255
		c := color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
		if verbose > 0 {
			fmt.Println(c)
		}
		drawBézier3(x1, y1, x2, y2, x3, y3, x4, y4, r, b3Seg, c, img)
	}
}

func BezierGa(org, diff *image.RGBA, max float64) gago.GA {
	return gago.GA{
		NbrPopulations: 6,
		NbrIndividuals: 80,
		NbrGenes:       13,
		Ff: gago.Float64Function{
			Image: getFit(ApplyBezier, org, diff, max),
		},
		Initializer: gago.InitUniformF{
			Lower: 1,
			Upper: 1,
		},
		Model: gago.ModSteadyState{
			Selector: gago.SelTournament{
				NbParticipants: 3,
			},
			Crossover: BezierCross{},
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

type BezierCross struct{}

func (cross BezierCross) Apply(p1 gago.Individual, p2 gago.Individual,
	rng *rand.Rand) (gago.Individual, gago.Individual) {
	var (
		o1 = gago.Individual{
			Genome:    make([]interface{}, 13),
			Fitness:   math.Inf(1),
			Evaluated: false,
			Name:      "-",
		}
		o2 = gago.Individual{
			Genome:    make([]interface{}, 13),
			Fitness:   math.Inf(1),
			Evaluated: false,
			Name:      "-",
		}
	)
	for i := 0; i < 8; i += 2 {
		p := rng.Float64()
		for j := i; j < i+2; j++ {
			o1.Genome[j] = p*p1.Genome[j].(float64) + (1-p)*p2.Genome[j].(float64)
			o2.Genome[j] = (1-p)*p1.Genome[j].(float64) + p*p2.Genome[j].(float64)
		}
	}
	for i := 8; i < 10; i++ {
		p := rng.Float64()
		o1.Genome[i] = p*p1.Genome[i].(float64) + (1-p)*p2.Genome[i].(float64)
		o2.Genome[i] = (1-p)*p1.Genome[i].(float64) + p*p2.Genome[i].(float64)
	}
	p := rng.Float64()
	for i := 10; i < 13; i++ {
		o1.Genome[i] = p*p1.Genome[i].(float64) + (1-p)*p2.Genome[i].(float64)
		o2.Genome[i] = (1-p)*p1.Genome[i].(float64) + p*p2.Genome[i].(float64)
	}
	return o1, o2
}
