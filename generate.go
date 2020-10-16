package main

import "github.com/rubenvp8510/jaeger-storage-perf/helpers"

func main()  {
	err := helpers.GenerateFixtures("traces_fixtures")
	if err != nil {
		panic(err)
	}
}
