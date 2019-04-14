package fuzzy

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

var negInf = math.Inf(-1)
var posInf = math.Inf(1)

type System struct {
	Universes []Universe
}

func (s *System) Fit(inData ...[]float64) error {
	// check the arrays have all the same length
	if len(inData) < 2 {
		return errors.New("expect two slice")
	}
	l := len(inData[0])
	for _, slice := range inData {
		if len(slice) != l {
			return errors.New("slice must be of same length")
		}
	}

	for i, slice := range inData {
		sort.Float64s(slice)
		min := slice[0]
		max := slice[len(slice)-1]
		s.Universes[i] = NewUniverse(fmt.Sprintf("inData%d", i), min, max)
	}
	// Divide the input and output spaces into fuzzy regions
	for i, slice := range inData {
		if err := divideUniverseInRegions(&s.Universes[i], slice); err != nil {
			return err
		}
	}
	// TODO: Generate fuzzy rules from given data pairs

	// TODO: Assign a degree to each rule

	// TODO: Create a combined fuzzy rule base

	// TODO: Determine a mapping based on the combined fuzzy rule base

	return nil
}

func divideUniverseInRegions(u *Universe, slice []float64) error {
	previous := Point{negInf, MaxMembershipDegree}
	for j := 0; j < len(slice); j++ {
		var next Point
		if j < len(slice)-1 {
			next = Point{slice[j] + 1, MinMembershipDegree}
		} else {
			next = Point{posInf, MaxMembershipDegree}
		}
		this := Point{slice[j], MaxMembershipDegree}
		points := newPoints(previous, this, next)
		if _, err := u.NewFuzzyMultipointSet(points); err != nil {
			return err
		}
		previous.X = slice[j]
		previous.Y = MinMembershipDegree
	}
	return nil
}
