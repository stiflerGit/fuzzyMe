package knowledge

import (
	"fuzzyMe/fuzzy"
)

//Database contains the membership functions
type Database struct {
	Sets []fuzzy.Set
}

//Rule base contains the fuzzy rules
type Rulebase struct {
	Rules []Rules
}
