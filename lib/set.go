package lib

type NumberSet struct {
	intervals []Interval
}

type Interval struct {
	Start int
	End   int
}
