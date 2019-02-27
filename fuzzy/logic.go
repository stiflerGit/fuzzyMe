package fuzzy

import (
	"fmt"
	"math"
)

type Implication func(a, b float64) float64

// ruleBase is a Set of rules
type RuleBase []rule

// create a new ruleBase
func NewRuleBase() RuleBase {
	return RuleBase{}
}

// NewRule add a rule to the ruleBase
func (r *RuleBase) NewRule() ruleIn {
	*r = append(*r, rule{})
	ruleIndex := len(*r) - 1
	return ruleIn{
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
	ruleIn   rule
	ruleIf   rule
	ruleIs   rule
	ruleThen rule
	ruleAnd  rule
	ruleOr   rule
)

// ruleIn has only IF function. So the first operation MUST be an IF
func (r ruleIn) IF(set Set) ruleIf {
	return ruleIf{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			var err error
			*s, err = (*s).Intersect(set)
			if err != nil {
				panic(err)
			}
		}}
}

// ruleIf has only IS function. Hence, an IF operation MUST be followed by an IS
func (r ruleIf) IS(set Set) ruleIs {
	return ruleIs{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			r._func(&set)
			*s = set
		}}
}

// ruleIs can be followed either by an OR operation or an AND operation or a
// THEN operation.
// OR operation results in selecting the maximum among all propositions.
// Proposition is the smallest piece of rule: i.e. setA.is(setB)
func (r ruleIs) OR(set Set) ruleOr {
	return ruleOr{
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

// ruleIs can be followed either by an OR operation or an AND operation or a
// THEN operation.
// OR operation results in selecting the minimum among all propositions.
// Proposition is the smallest piece of rule: i.e. setA.is(setB)
func (r ruleIs) AND(set Set) ruleAnd {

	return ruleAnd{
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

// ruleIs can be followed either by an OR operation or an AND operation or a
// THEN operation.
// THEN operation results in a projection of the input set on the output
// co-domain
func (r ruleIs) THEN(Set) ruleThen {

	return ruleThen{
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

// ruleThen has only IS function. Then operation represent the logic
// implication
func (r ruleThen) IS(set Set) {
	*r.baseElem = rule{
		baseElem: nil,
		_func: func(s *Set) {
			r._func(&set)
			*s = set
		}}
}

// ruleAnd has only IS function. Hence, an IF operation MUST be followed by an IS
func (r ruleAnd) IS(set Set) ruleIs {
	return ruleIs{
		baseElem: r.baseElem,
		_func: func(s *Set) {
			r._func(&set)
			*s = set
		}}
}

// ruleOr has only IS function. Hence, an IF operation MUST be followed by an IS
func (r ruleOr) IS(set Set) ruleIs {
	return ruleIs{
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
