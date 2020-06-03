package profiling

import (
	"fmt"
	"time"
)

func MeasureTime(stage string, f func() error) (time.Duration, error) {
	fmt.Printf(">> start stage=%s\n", stage)
	start := time.Now()
	err := f()
	if err != nil {
		return 0, err
	}
	fmt.Printf(">> completed stage=%s duration=%s\n", stage, time.Since(start))
	return time.Since(start), nil
}
