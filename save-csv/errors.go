package main

func PanicError(e error) {
	if e != nil {
		panic(e)
	}
}
