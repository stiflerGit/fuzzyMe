package main

import "fuzzyMe/fuzzy/set"

func Fuzzify(u float64) set.Singleton {
	// u must belong to an universe
	return set.Singleton{u}
}

func Defuzzify(s set.Set) float64 {

}
