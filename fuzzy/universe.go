package fuzzy

import (
	"fmt"
)

type Universe struct {
	name     string
	min, max float64 // minimum and maximum value a base variable belonging to it can assume
}

func NewUniverse(name string, min, max float64) Universe {
	return Universe{name, min, max}
}

func (u *Universe) NewVar(v float64) base {
	return base{universe: u, val: v}
}

type base struct {
	universe *Universe
	val      float64
}

func (b *base) Set(v float64) (float64, error) {
	if v < b.universe.min {
		return 0, fmt.Errorf("minimum value can be: %f", b.universe.min)
	}
	if v > b.universe.max {
		return 0, fmt.Errorf("maximum value can be: %f", b.universe.max)
	}
	b.val = v
	return b.val, nil
}

func (b *base) IS(s Set) *Rule {
	return &Rule{
		props: []mf{
			s.mf,
		},
	}
}
