package set

type Set interface {
	Intersect() Set
	Union() Set
	Complement() Set
}

type FuzzySet struct {
}

type Singleton struct {
	center float64
}
