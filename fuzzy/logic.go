package fuzzy

import (
	"math"
)

type Implication func(a, b float64) float64

type RuleBase []rule

func NewRuleBase() RuleBase {
	return make(RuleBase, 1)
}

func (r RuleBase) NewRule() *rule {
	r = append(r, rule{})
	return &r[len(r)-1]
}

func (r RuleBase) Exec() Set {
	for _, rule := range r {
		rule.EXEC()
	}
	// TODO: IMPLEMENT ME
	panic("IMPLEMENT ME")
	return nil
}

type rule struct {
	sets []Set // propositions
}

func (r *rule) IF(s Set) *rule {
	r.sets = append(r.sets, s)
	return r
}

func (r *rule) IS(s Set) *rule {
	if r.sets == nil {
		panic("first call must be to IF")
	}
	set, err := s.Intersect(r.sets[len(r.sets)-1])
	if err != nil {
		panic(err)
	}
	r.sets[len(r.sets)-1] = set
	return r
}

func (r *rule) THEN(s Set) *rule {
	if r.sets == nil {
		panic("first call must be to IF")
	}
	// TODO: Improve for every type of Sets
	min := math.Inf(1)
	for _, set := range r.sets {
		if singleton, ok := set.(Singleton); ok {
			if singleton.y < min {
				min = singleton.y
			}
		} else {
			panic("Not a singleton: find a solution")
		}
	}
	projectionSet, err := newMultiPointSet(s.Universe(), Points{{0, min}})
	if err != nil {
		panic(err)
	}
	r.sets = []Set{projectionSet}
	return r
}

// TODO: MISO Systems
func (r *rule) AND(s Set) *rule {
	r.sets = append(r.sets, s)
	return r
	//min := 0.0
	//for _, set := range r.sets {
	//	if singleton, ok := set.(Singleton); ok {
	//		if singleton.y < min {
	//			min = singleton.y
	//		}
	//	} else {
	//		panic("Not a singleton: find a solution")
	//	}
	//}
}

func (r *rule) OR(s Set) *rule {
	//TODO: IMPLEMENT ME
	panic("implement me")
	return nil
	//max := 0.0
	//for _, set := range r.sets {
	//	if singleton, ok := set.(Singleton); ok {
	//		if singleton.y > max {
	//			max = singleton.y
	//		}
	//	} else {
	//		panic("Not a singleton: find a solution")
	//	}
	//}
}

func (r *rule) EXEC() Set {
	return r.sets[0]
}

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
