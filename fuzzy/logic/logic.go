package logic

import "math"

type Implication func(a, b float64) float64

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
