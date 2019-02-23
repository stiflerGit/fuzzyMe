package fuzzy

type Universe struct {
	name     string
	min, max float64 // minimum and maximum value a base variable belonging to it can assume
}

func NewUniverse(name string, min, max float64) Universe {
	return Universe{name, min, max}
}

func (u *Universe) NewFuzzySet(points points) (Set, error) {
	set, err := newFuzzySet(u, points)
	if err != nil {
		return Set{}, err
	}
	return set, nil
}

func (u *Universe) NewFuzzySingleton(v float64) (Set, error) {
	set, err := newFuzzySingleton(u, v)
	if err != nil {
		return Set{}, err
	}
	return set, nil
}
