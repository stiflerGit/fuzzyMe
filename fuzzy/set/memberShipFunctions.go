package set

const (
	MembershipDegreeMin MembershipDegree = 0.0
	MembershipDegreeMax MembershipDegree = 1.0
)

type MembershipDegree float32

type MembershipFunc func(interface{}) MembershipDegree

func triangular(center, base float64) MembershipFunc {
	c := center
	b := base
	return func(elem interface{}) MembershipDegree {
		// TODO: IMPLEMENT
		return 0
	}
}

/*
func trapezoidal(points ...float64) MembershipFunc
*/
