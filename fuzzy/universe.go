package fuzzy

const (
	MaxMembershipDegree = float64(1.0)
	MinMembershipDegree = float64(0.0)
)

type Universe struct {
	name     string
	min, max float64 // minimum and maximum value a base variable belonging to it can assume
	// Logic* ??? // TODO: add logic to a universe (How to Intersect, Union, ... sets belonging to the universe)
}

func NewUniverse(name string, min, max float64) Universe {
	return Universe{name, min, max}
}

func (u Universe) GetName() string {
	return u.name
}

func (u Universe) GetMin() float64 {
	return u.min
}

func (u Universe) GetMax() float64 {
	return u.max
}

func (u *Universe) NewFuzzyMultipointSet(points Points) (MultiPointSet, error) {
	set, err := newMultiPointSet(u, points)
	if err != nil {
		return MultiPointSet{}, err
	}
	return set, nil
}

func (u *Universe) NewFuzzySingleton(v float64) (Singleton, error) {
	set, err := newFuzzySingleton(u, v)
	if err != nil {
		return Singleton{}, err
	}
	return set, nil
}
