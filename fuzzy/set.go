package fuzzy

import (
	"errors"
	"fmt"
	. "fuzzyMe/math"
	"math"
)

//
//type Set interface {
//	Intersect() Set
//	Union() Set
//	Complement() Set
//}

// a fuzzy set is defined by it's universe of application, and
// a membership function
type Set struct {
	universe *Universe // universe where the set is defined(ex: Temperature)
	mf       mf        // membership function
}

func newFuzzySingleton(universe *Universe, x float64) (Set, error) {
	if x < universe.min || x > universe.max {
		return Set{}, errors.New("value out of universe bounds")
	}
	return Set{universe: universe, mf: mf(func(u float64) float64 {
		if u != x {
			return 0
		}
		return 1
	})}, nil
}

func newFuzzySet(universe *Universe, points points) (Set, error) {
	if len(points) < 2 {
		return Set{}, errors.New("too few points")
	}
	mf, err := multiPoint(points...)
	if err != nil {
		return Set{}, err
	}
	return Set{universe: universe, mf: mf}, nil
}

func (f Set) Intersect(in Set) (Set, error) {
	if f.universe != in.universe {
		return Set{}, fmt.Errorf("can't intersect Sets from different universes")
	}
	return Set{
		universe: f.universe, // universe same as originals
		mf: mf(func(u float64) float64 {
			return Minimum(f.mf(u), in.mf(u))
		})}, nil
}

func (f Set) Union(in Set) (Set, error) {
	if f.universe != in.universe {
		return Set{}, fmt.Errorf("can't union Sets from different universes")
	}
	return Set{
		universe: f.universe,
		mf: mf(func(u float64) float64 {
			return Maximum(f.mf(u), in.mf(u))
		})}, nil
}

func (f Set) Complement() Set {
	return Set{
		universe: f.universe,
		mf: mf(func(u float64) float64 {
			return 1 - f.mf(u)
		})}
}

type Singleton struct {
	center float64
}

const (
	MembershipDegreeMin membershipDegree = 0.0
	MembershipDegreeMax membershipDegree = 1.0
)

type membershipDegree float64

type mf func(float64) float64

func triangular(center, base float64) mf {
	mf, _ := multiPoint(
		point{center - base/2, 0},
		point{center, 1},
		point{center + base/2, 0})
	return mf
}

func trapezioid(center, majorBase, minorBase float64) mf {
	mf, _ := multiPoint(
		point{center - majorBase/2, 0},
		point{center - minorBase/2, 1},
		point{center + minorBase/2, 1},
		point{center + majorBase/2, 0})
	return mf
}

//
// segment1 := p1 + (p2-p1) * s
// segment2 := (u,0) + (0,1) * t
// intersection := (u,0) + (0,1)t = p1+(p2-p1)s
//
// u = p1.x + (p2.x-p1.x)*s
// t = p1.y+(p2.y-p1.y)*s
//
// s = (u - p1.x) / (p2.x-p1.x)
//
// y = p1.y + (p2.y-p1.y)*s
func multiPoint(points ...point) (mf, error) {
	myPoints := newPoints(points...)
	return func(u float64) float64 {
		if u == math.Inf(-1) {
			_, min := myPoints.findMinY()
			return min.y
		}
		if u == math.Inf(1) {
			_, max := myPoints.findMaxY()
			return max.y
		}

		var p1, p2 point
		// find p1; p2 := next of p2
		i, _ := myPoints.findLesserX(u)
		// last point is impossible 'cause is inf
		if i == len(myPoints)-2 {
			return myPoints[i].y
		} else if i == 0 {
			return myPoints[i+1].y
		}
		p1 = myPoints[i]
		p2 = myPoints[i+1]
		s := (u - p1.x) / (p2.x - p1.x)
		return p1.y + (p2.y-p1.y)*s
	}, nil
}

func (mf mf) Min() float64 {
	return mf(math.Inf(-1))
}

func (mf mf) Max() float64 {
	return mf(math.Inf(1))
}

//
// Points
//
type point struct {
	x, y float64
}

type points []point

func newPoints(p ...point) points {
	// TODO: sort points for x
	firstPoint := p[0]
	lastPoint := p[len(p)-1]
	P := points{}
	P = append(points{point{math.Inf(-1), firstPoint.y}}, p...)
	P = append(P, point{math.Inf(1), lastPoint.y})
	return P
}

func (p points) firstPoint() point {
	return p[1]
}

func (p points) lastPoint() point {
	return p[len(p)-1]
}

func (p points) findLesserX(x float64) (int, point) {
	var i int
	for i = range p {
		if p[i].x > x {
			break
		}
	}
	return i - 1, p[i-1]
}

func (p points) findMinY() (int, point) {
	iMin, min := 0, p.firstPoint().y
	for i := range p {
		if p[i].y < min {
			iMin = i
			min = p[i].y
		}
	}
	return iMin, p[iMin]
}

func (p points) findMaxY() (int, point) {
	iMax, max := 0, p.firstPoint().y
	for i := range p {
		if p[i].y > max {
			iMax = i
			max = p[i].y
		}
	}
	return iMax, p[iMax]
}
