package fuzzy

import (
	"fmt"
	"math"
)

type Implication func(a, b float64) float64

type RuleBase []rule

func NewRuleBase() RuleBase {
	return RuleBase{}
}

func (r *RuleBase) NewRule() rule {
	*r = append(*r, nil)
	return (*r)[len(*r)-1]
}

func (r *RuleBase) Exec() (Set, error) {
	var res, set Set
	var err error
	if len(*r) < 1 {
		// is it an error ???
		return nil, nil
	}
	res, err = (*r)[0].EXEC()
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(r); i++ {
		set, err = r[i].EXEC()
		if err != nil {
			return nil, err
		}
		res, err = set.Union(set)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

type rule func(*Set) rule

func (r rule) IF(set Set) rule {

	return func(s *Set) rule {
		var err error
		if s == nil || (*s) == nil {
			panic("IF must be followed by IS")
		}
		*s, err = (*s).Intersect(set)
		if err != nil {
			panic(err)
		}
		return r
	}
}

func (r rule) IS(set Set) rule {

	return func(s *Set) rule {
		r(&set)
		*s = set
		return r
	}
}

func (r rule) THEN(set Set) rule {

	return func(s *Set) rule {
		if s == nil || (*s) == nil {
			panic("THEN must be followed by IS")
		}
		// TODO: Improve for every type of Sets
		var domainSet Set
		r(&domainSet)
		min := math.Inf(1)
		if singleton, ok := domainSet.(Singleton); ok {
			min = singleton.y
		} else {
			panic("not a singleton: find a solution")
		}
		projectionSet, err := newMultiPointSet((*s).Universe(), Points{{0, min}})
		if err != nil {
			panic(err)
		}
		*s, err = (*s).Intersect(projectionSet)
		if err != nil {
			panic(err)
		}
		return r
	}
}

func (r rule) EXEC() (set Set, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	r(&set)
	return set, err
}

func (r rule) AND(set Set) rule {
	return func(s *Set) rule {
		if s == nil || *s == nil {
			panic("AND must be followed by IS")
		}
		var err error
		*s, err = (*s).Intersect(set)
		if err != nil {
			panic(err)
		}
		var leftSet Set
		r(&leftSet)
		leftSingleton, leftOk := leftSet.(Singleton)
		rightSingleton, rightOk := (*s).(Singleton)
		if !leftOk || !rightOk {
			panic("unimplemented AND without singleton")
		}
		if leftSingleton.y < rightSingleton.y {
			*s = leftSingleton
		}
		return r
	}
}

func (r rule) OR(set Set) rule {
	return func(s *Set) rule {
		if s == nil || *s == nil {
			panic("OR must be followed by IS")
		}
		var err error
		*s, err = (*s).Intersect(set)
		if err != nil {
			panic(err)
		}
		var leftSet Set
		r(&leftSet)
		leftSingleton, leftOk := leftSet.(Singleton)
		rightSingleton, rightOk := (*s).(Singleton)
		if !leftOk || !rightOk {
			panic("unimplemented AND without singleton")
		}
		if leftSingleton.y > rightSingleton.y {
			*s = leftSingleton
		}
		return r
	}
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
