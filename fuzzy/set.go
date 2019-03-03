package fuzzy

import (
	"errors"
	"fmt"
	. "github.com/stiflerGit/fuzzyMe/math"
	"math"
)

type Set interface {
	Universe() *Universe
	MembershipDegree(float64) (float64, error)
	Intersect(Set) (Set, error)
	Union(Set) (Set, error)
	Complement() Set
}

// membership function gives the degree to which an element belongs to a fuzzy
// set. It characterizes the fuzzy set itself.
type mf func(float64) float64

// a fuzzy set is defined by it's universe of application, and
// a membership function
type FuzzySet struct {
	universe *Universe // universe where the set is defined(ex: Temperature)
	mf       mf        // membership function
}

func NewFuzzySet(universe *Universe, membershipFunc mf) (FuzzySet, error) {
	if membershipFunc == nil {
		return FuzzySet{}, errors.New("membership function can't be nil")
	}
	return FuzzySet{universe: universe, mf: membershipFunc}, nil
}

func (s FuzzySet) Universe() *Universe {
	return s.universe
}

func (s FuzzySet) MembershipDegree(x float64) (float64, error) {
	if x < s.universe.min || x > s.universe.max {
		return 0, errors.New("value out of universe bounds")
	}
	return s.mf(x), nil
}

func intersect(a, b Set) mf {
	return func(x float64) float64 {
		aDegree, _ := a.MembershipDegree(x)
		bDegree, _ := b.MembershipDegree(x)
		return Minimum(aDegree, bDegree)
	}
}

func (s FuzzySet) Intersect(in Set) (Set, error) {
	//if s.universe != in.universe {
	//	return FuzzySet{}, fmt.Errorf("can't intersect Sets from different universes")
	//}
	if _, ok := in.(Singleton); ok {
		return in.Intersect(s)
	}
	return NewFuzzySet(s.universe, intersect(s, in))
}

func union(a, b Set) mf {
	return func(x float64) float64 {
		aDegree, _ := a.MembershipDegree(x)
		bDegree, _ := b.MembershipDegree(x)
		return Maximum(aDegree, bDegree)
	}
}

func (s FuzzySet) Union(in Set) (Set, error) {
	//if s.universe != in.universe {
	//	return FuzzySet{}, fmt.Errorf("can't union Sets from different universes")
	//}
	if _, ok := in.(Singleton); ok {
		return in.Intersect(s)
	}
	return NewFuzzySet(s.universe, union(s, in))
}

func (s FuzzySet) Complement() Set {
	return FuzzySet{
		universe: s.universe,
		mf: mf(func(u float64) float64 {
			return 1 - s.mf(u)
		})}
}

type MultiPointSet struct {
	universe *Universe
	points   Points
}

func newMultiPointSet(universe *Universe, points Points) (MultiPointSet, error) {
	if len(points) < 1 {
		return MultiPointSet{}, errors.New("no points")
	}
	pointsPlusInfinite := newPoints(points...)
	return MultiPointSet{universe, pointsPlusInfinite}, nil
}

func (s MultiPointSet) Universe() *Universe {
	return s.universe
}

// segment1 := p1 + (p2-p1) * r // segment2 := (u,0) + (0,1) * t
// intersection := (u,0) + (0,1)t = p1+(p2-p1)r
// u = p1.X + (p2.X-p1.X)*r
// t = p1.Y+(p2.Y-p1.Y)*r
// r = (u - p1.X) / (p2.X-p1.X)
// Y = p1.Y + (p2.Y-p1.Y)*s
func (s MultiPointSet) MembershipDegree(x float64) (float64, error) {
	if x < s.universe.min || x > s.universe.max {
		return 0, errors.New("value out of universe bounds")
	}
	if x == math.Inf(-1) {
		_, min := s.points.findMinY()
		return min.Y, nil
	}
	if x == math.Inf(1) {
		_, max := s.points.findMaxY()
		return max.Y, nil
	}

	var p1, p2 Point
	// find p1; p2 := next of p2
	i, _ := s.points.findLesserX(x)
	// last Point is impossible 'cause is inf
	if i == len(s.points)-2 {
		return s.points[i].Y, nil
	} else if i == 0 {
		return s.points[i+1].Y, nil
	}
	p1 = s.points[i]
	p2 = s.points[i+1]
	r := (x - p1.X) / (p2.X - p1.X)
	return p1.Y + (p2.Y-p1.Y)*r, nil
}

func (s MultiPointSet) Intersect(r Set) (Set, error) {
	if s.universe != r.Universe() {
		return MultiPointSet{}, fmt.Errorf("can't intersect Sets from different universes")
	}
	if _, ok := r.(Singleton); ok {
		return r.Intersect(s)
	}
	return NewFuzzySet(s.universe, intersect(s, r))
}

func (s MultiPointSet) Union(r Set) (Set, error) {
	if s.universe != r.Universe() {
		return MultiPointSet{}, fmt.Errorf("can't union Sets from different universes")
	}
	if _, ok := r.(Singleton); ok {
		return r.Intersect(s)
	}
	return NewFuzzySet(s.universe, union(s, r))
}

func (s MultiPointSet) Complement() Set {
	var complementedSet MultiPointSet
	complementedSet.points = append(Points{}, s.points...)
	for i := range complementedSet.points {
		complementedSet.points[i].Y = 1 - s.points[i].Y
	}
	return complementedSet
}

// A fuzzy singleton is a set with only one value whose degree of membership is
// equal to 1 while all the other values have a degree of membership = 0
type Singleton struct {
	universe *Universe
	x, y     float64
}

func newFuzzySingleton(universe *Universe, x float64) (Singleton, error) {
	if x < universe.min || x > universe.max {
		return Singleton{}, errors.New("value out of universe bounds")
	}
	return Singleton{universe, x, 1}, nil
}

func (s Singleton) Universe() *Universe {
	return s.universe
}

func (s Singleton) MembershipDegree(x float64) (float64, error) {
	if x != s.x {
		return 0, nil
	}
	return 1, nil
}

func (s Singleton) Intersect(r Set) (Set, error) {
	md, _ := r.MembershipDegree(s.x)
	return Singleton{s.universe, s.x, md}, nil
}

func (s Singleton) Union(e Set) (Set, error) {
	// TODO: IMPLEMENT ME
	panic("IMPLEMENT ME")
	return e, nil
}

func (s Singleton) Complement() Set {
	set, _ := NewFuzzySet(s.universe,
		func(x float64) float64 {
			if x != s.x {
				return 1
			}
			return 0
		})
	return set
}

const (
	MembershipDegreeMin membershipDegree = 0.0
	MembershipDegreeMax membershipDegree = 1.0
)

type membershipDegree float64

//
// Points
//
type Point struct {
	X, Y float64
}

type Points []Point

func newPoints(p ...Point) Points {
	// TODO: sort Points for X
	firstPoint := p[0]
	lastPoint := p[len(p)-1]
	P := Points{}
	P = append(Points{Point{math.Inf(-1), firstPoint.Y}}, p...)
	P = append(P, Point{math.Inf(1), lastPoint.Y})
	return P
}

func (p Points) firstPoint() Point {
	return p[1]
}

func (p Points) lastPoint() Point {
	return p[len(p)-1]
}

func (p Points) findLesserX(x float64) (int, Point) {
	var i int
	for i = range p {
		if p[i].X > x {
			break
		}
	}
	return i - 1, p[i-1]
}

func (p Points) findMinY() (int, Point) {
	iMin, min := 0, p.firstPoint().Y
	for i := range p {
		if p[i].Y < min {
			iMin = i
			min = p[i].Y
		}
	}
	return iMin, p[iMin]
}

func (p Points) findMaxY() (int, Point) {
	iMax, max := 0, p.firstPoint().Y
	for i := range p {
		if p[i].Y > max {
			iMax = i
			max = p[i].Y
		}
	}
	return iMax, p[iMax]
}
