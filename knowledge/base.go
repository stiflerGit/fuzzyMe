package knowledge

import "fuzzyMe/fuzzy/set"

//Database contains the membership functions
type Database struct {
	Sets []set.MembershipFunc
}

//Rule base contains the fuzzy rules
type Rulebase struct {
	Rules []Rules
}
