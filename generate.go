package main

import (
	"fmt"
	"github.com/rubenvp8510/jaeger-storage-perf/fixtures"
)

func main()  {
	nSpans, err := fixtures.SaveSpans("data/traces",3500, 100, 100)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Spans writen %d\n", nSpans)
}
