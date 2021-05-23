package main

import (
	"fmt"

	"github.com/53jk1/pok2/interpolate"
	"github.com/53jk1/pok2/interpolate/linear"
)

var x float64

func main() {
	x := []float64{1.3, 1.8, 2.5, 3.1, 3.8, 4.4, 4.9, 5.5, 6.2}
	y := []float64{3.37, 4.45, 4.81, 3.96, 3.31, 2.72, 3.02, 3.43, 4.07}
	valToInterp := 5.1

	li := linear.New()
	li.Fit(x, y)

	estimate, err := interpolate.WithSingle(li, valToInterp)
	fmt.Println(estimate)
	fmt.Println(err)
}
