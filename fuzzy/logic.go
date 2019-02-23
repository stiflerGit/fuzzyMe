package fuzzy

import (
	. "fuzzyMe/math"
	"math"
)

type Implication func(a, b float64) float64

type Rule struct {
	sets []Set // propositions
}

func (r *Rule) IF(in *Set) *Rule {
	r.sets = append(r.sets, *in)
	return r
}

func (r *Rule) THEN(in *Rule) Set {
	// TODO: Improve for every type of Sets
	min := r.sets[0].mf.Min()
	mf, err := multiPoint(
		point{math.Inf(-1), min},
		point{math.Inf(1), min},
	)
	if err != nil {
		panic("err")
	}
	domainSet := Set{
		mf: mf,
	}
	return Set{mf:mf}
	}
}

/* // TODO: MISO Systems
func (r *Rule) AND() *Rule {
	//TODO: IMPLEMENT ME
	return nil
}
*/

func (r *Rule) OR() *Rule {
	//TODO: IMPLEMENT ME
	return nil
}

// and??

// or

// Implications
var KleeneDienes = Implication(
	func(a, b float64) float64 {
		return math.Max(1-a, b)
	})

var Lukasiewicz = Implication(
	func(a, b float64) float64 {
		return math.Min(1, 1-a+b)
	})

var Zadeh = Implication(
	func(a, b float64) float64 {
		return math.Max(math.Min(a, b), 1-a)
	})

var Mamdani = Implication(math.Min)

var Larsen = Implication(
	func(a, b float64) float64 {
		return a * b
	})

var _if = mf(func(float64) float64 {
	return 0
})
