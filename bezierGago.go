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
		b3Seg := int(X[i+8]*60 + 2)
		red := X[i+9] * 255
		green := X[i+10] * 255
		blue := X[i+11] * 255
		c := color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
		if verbose > 0 {
			fmt.Println(c)
		}
		drawBézier3(x1, y1, x2, y2, x3, y3, x4, y4, b3Seg, c, img)
	}
}

func BezierGa(org, diff *image.RGBA, max float64) gago.GA {
	return gago.GA{
		NbrPopulations: 6,
		NbrIndividuals: 80,
		NbrGenes:       12,
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
			Genome:    make([]interface{}, 12),
			Fitness:   math.Inf(1),
			Evaluated: false,
			Name:      "-",
		}
		o2 = gago.Individual{
			Genome:    make([]interface{}, 12),
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
	p := rng.Float64()
	o1.Genome[8] = p*p1.Genome[8].(float64) + (1-p)*p2.Genome[8].(float64)
	o2.Genome[8] = (1-p)*p1.Genome[8].(float64) + p*p2.Genome[8].(float64)
	p = rng.Float64()
	for i := 9; i < 12; i++ {
		o1.Genome[i] = p*p1.Genome[i].(float64) + (1-p)*p2.Genome[i].(float64)
		o2.Genome[i] = (1-p)*p1.Genome[i].(float64) + p*p2.Genome[i].(float64)
	}
	return o1, o2
}
