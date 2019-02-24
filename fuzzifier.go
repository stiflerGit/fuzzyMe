package main

import (
	. "fuzzyMe/fuzzy"
)

const (
	resolution = 0.01
)

type defuzzType string

const (
	COG = defuzzType("COG") // Center Of Gravity
	MOM = defuzzType("MOM") // Mean Of Maxima
)

func Defuzzify(s Set, defuzzType defuzzType) float64 {
	switch defuzzType {
	case COG:
		return centerOfGravity(s)
	case MOM:
		return meanOfMaxima(s)
	default:
		panic("unimplemented defuzzification type")
	}
}

func centerOfGravity(s Set) float64 {
	interval := s.Universe().GetMax() - s.Universe().GetMin()
	nSamples := int(interval / resolution)
	muSum := 0.0
	nominator := 0.0
	for i := 0; i < nSamples; i++ {
		xi := float64(i) * resolution
		mu, _ := s.MembershipDegree(xi)
		nominator += xi * mu
		muSum += mu
	}
	return nominator / muSum
}

func meanOfMaxima(s Set) float64 {
	// TODO: IMPLEMENT ME
	panic("Implement me")
	return 0
}
