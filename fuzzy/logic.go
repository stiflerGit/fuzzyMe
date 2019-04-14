package fuzzy

import (
	"fmt"
	"math"
)

type Set interface {
	Universe() *Universe
	MembershipDegree(float64) (float64, error)
	Intersect(Set) (Set, error)
	Union(Set) (Set, error)
	Complement() Set
}

type Implication func(a, b float64) float64

// ruleBase is a Set of rules
type RuleBase []rule

// create a new ruleBase
func NewRuleBase() RuleBase {
	return RuleBase{}
}

// TODO: document
// see https://github.com/golang/lint/issues/210
//type (
//	RuleIn interface {
//		IF(Set) RuleConjunction
//	}
//	RuleIs interface {
//		AND(Set) RuleConjunction
//		OR(Set) RuleConjunction
//		THEN(Set) RuleThen
//	}
//	RuleConjunction interface {
//		IS(Set) RuleIs
//	}
//	RuleThen interface {
//		IS(Set)
//	}
//)

// NewRule add a rule to the ruleBase
func (r *RuleBase) NewRule() RuleIn {
	*r = append(*r, rule{})
	ruleIndex := len(*r) - 1
	return RuleIn{
		baseElem: &(*r)[ruleIndex],
		_func:    nil,
	}
}

// Exec execute all the rules in the ruleBase and return the resulting FuzzySet
func (r *RuleBase) Exec() (Set, error) {
	var err error
	// to allow rule chaining rule can't return error, that's why rules
	// panic. Exec capture panic and returns it
	defer func(err *error) {
		if r := recover(); r != nil {
			*err = fmt.Errorf("error in execution %v", r)
		}
	}(&err)
	// one output set for each rule
	sets := make([]Set, len(*r))
	for i, rule := range *r {
		rule._func(&sets[i])
	}
	// union of all resulting set
	res := sets[0]
	for i := 1; i < len(sets); i++ {
		res, err = res.Union(sets[i])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// rule represent an operation on a FuzzySet
type rule struct {
	baseElem *rule      // pointer to the rule in the ruleBase
	_func    func(*Set) // function representing the logic of the rule
}

// a type for each rule, avoid to chain function call wrongly
// (ex of wrong chain: rule.IF(a).IF(b))
type (
	RuleIn rule
	RuleIf rule
	RuleIs rule
	RuleThen rule
	RuleAnd rule
	RuleOr rule
)

// RuleIn has only IF function. So the first operation MUST be an IF
func (r RuleIn) IF(set Set) RuleIf {
	return RuleIf{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			var err error
			*s, err = (*s).Intersect(set)
			if err != nil {
				panic(err)
			}
		}}
}

// RuleIf has only IS function. Hence, an IF operation MUST be followed by an IS
func (r RuleIf) IS(set Set) RuleIs {
	return RuleIs{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			r._func(&set)
			*s = set
		}}
}

// RuleIs can be followed either by an OR operation or an AND operation or a
// THEN operation.
// OR operation results in selecting the maximum among all propositions.
// Proposition is the smallest piece of rule: i.e. setA.is(setB)
func (r RuleIs) OR(set Set) RuleOr {
	return RuleOr{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			var err error
			*s, err = (*s).Intersect(set)
			if err != nil {
				panic(err)
			}
			var leftSet Set
			r._func(&leftSet)
			leftSingleton, leftOk := leftSet.(Singleton)
			rightSingleton, rightOk := (*s).(Singleton)
			if !leftOk || !rightOk {
				panic("unimplemented AND without singleton")
			}
			if leftSingleton.y > rightSingleton.y {
				*s = leftSingleton
			}
		}}
}

// RuleIs can be followed either by an OR operation or an AND operation or a
// THEN operation.
// OR operation results in selecting the minimum among all propositions.
// Proposition is the smallest piece of rule: i.e. setA.is(setB)
func (r RuleIs) AND(set Set) RuleAnd {

	return RuleAnd{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			var err error
			*s, err = (*s).Intersect(set)
			if err != nil {
				panic(err)
			}
			var leftSet Set
			r._func(&leftSet)
			leftSingleton, leftOk := leftSet.(Singleton)
			rightSingleton, rightOk := (*s).(Singleton)
			if !leftOk || !rightOk {
				panic("unimplemented AND without singleton")
			}
			if leftSingleton.y < rightSingleton.y {
				*s = leftSingleton
			}
		}}
}

// RuleIs can be followed either by an OR operation or an AND operation or a
// THEN operation.
// THEN operation results in a projection of the input set on the output
// co-domain
func (r RuleIs) THEN(Set) RuleThen {

	return RuleThen{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			// TODO: Improve for every type of Sets
			var domainSet Set
			r._func(&domainSet)

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
		},
	}
}

// RuleThen has only IS function. Then operation represent the logic
// implication
func (r RuleThen) IS(set Set) {
	*r.baseElem = rule{
		baseElem: nil,
		_func: func(s *Set) {
			r._func(&set)
			*s = set
		}}
}

// RuleAnd has only IS function. Hence, an IF operation MUST be followed by an IS
func (r RuleAnd) IS(set Set) RuleIs {
	return RuleIs{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			r._func(&set)
			*s = set
		}}
}

// RuleOr has only IS function. Hence, an IF operation MUST be followed by an IS
func (r RuleOr) IS(set Set) RuleIs {
	return RuleIs{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			r._func(&set)
			*s = set
		}}
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
