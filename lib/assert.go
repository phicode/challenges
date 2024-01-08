package lib

func Assert(c bool) {
	if !c {
		panic("assert failed")
	}
}
