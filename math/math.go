package math

import "math"

type Tnorm func(float64, float64) float64

type Snorm func(float64, float64) float64

// Triangular Norms
var Minimum = Tnorm(math.Min)

var Product = Tnorm(
	func(a, b float64) float64 {
		return a * b
	},
)

var BoundedProduct = Tnorm(
	func(a, b float64) float64 {
		return math.Max(0, a+b-1)
	},
)

// Triangular CoNorms
var Maximum = Snorm(
	func(a, b float64) float64 {
		return math.Max(a, b)
	},
)

var ProbabilisticSum = Snorm(
	func(a, b float64) float64 {
		return a + b - a*b
	},
)

var BoundedSum = Snorm(
	func(a, b float64) float64 {
		return math.Min(1, a+b)
	},
)
