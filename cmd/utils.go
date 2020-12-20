package main

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
