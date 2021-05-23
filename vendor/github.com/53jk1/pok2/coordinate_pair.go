package pok2

import (
	"sort"
)

// CoordinatePair to typ struktury z dwiema właściwościami, współrzędnymi X i Y.
type CoordinatePair struct {
	X float64
	Y float64
}

// SortCoordinatePairs to funkcja, która otrzymuje wycinek CoordinatePairs i sortuje go w porządku rosnącym.
func SortCoordinatePairs(cp []CoordinatePair) {
	sort.Slice(cp, func(i, j int) bool {
		return cp[i].X < cp[j].X
	})
}

// SlicesToCoordinatePairs to funkcja, która otrzymuje dwa wycinki liczb zmiennoprzecinkowych (x i y), zamienia je w wycinek CoordinatePairs i zwraca wynik.
func SlicesToCoordinatePairs(x, y []float64) []CoordinatePair {
	cp := make([]CoordinatePair, len(x))
	for i := 0; i < len(x); i++ {
		cp = append(cp, CoordinatePair{X: x[i], Y: y[i]})
	}
	return cp
}
