package interpolate

import (
	"fmt"

	"github.com/53jk1/pok2"
)

// Typ podstawowy zapewnia podstawową funkcjonalność dla dowolnego typu interpolacji
type Base struct {
	XYPairs []pok2.CoordinatePair
	X       []float64
	Y       []float64
}

// Fit otrzymuje dwa wycinki float64 - dla współrzędnych x i y, gdzie x[i] i y[i] reprezentują parę współrzędnych w siatce.
// Zwraca błąd, jeśli rozmiary X i Y nie są zgodne.
func (b *Base) Fit(x, y []float64) error {
	if len(x) != len(y) {
		return fmt.Errorf("Rozmiary X i Y nie pasują")
	}
	b.X = x
	b.Y = y
	b.XYPairs = pok2.SlicesToCoordinatePairs(x, y)
	pok2.SortCoordinatePairs(b.XYPairs)
	return nil
}
