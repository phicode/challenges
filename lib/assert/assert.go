package assert

func True(c bool) {
	if !c {
		panic("assert failed")
	}
}

func False(c bool) {
	if c {
		panic("assert failed")
	}
}

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}
