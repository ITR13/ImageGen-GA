package main

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
)

func CircleGa(function func([]float64) float64) gago.GA {
	return gago.GA{
		NbrPopulations: 4,
		NbrIndividuals: 80,
		NbrGenes:       6,
		Ff: gago.Float64Function{
			Image: function,
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
