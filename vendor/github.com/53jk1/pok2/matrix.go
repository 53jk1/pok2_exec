package pok2

import (
	"fmt"
	"math"
)

// Typ macierzy to wycinek wektorów z niestandardowymi metodami potrzebnymi do operacji na macierzach.
type Matrix []Vector

// Dim zwraca wymiary macierzy w postaci (wiersze, kolumny).
func (m Matrix) Dim() (int, int) {
	if m.isNil() {
		return 0, 0
	}
	return len(m), len(m[0])
}

// Invert zwraca odwróconą macierz przy użyciu eliminacji Gaussa-Jordana
func (m Matrix) Invert() (Matrix, error) {

	if !m.isSquare() {
		return nil, fmt.Errorf("Nie można odwrócić macierzy niekwadratowej")
	}

	var rows, _ = m.Dim()

	vec := make(Vector, rows)

	// 1. Redukcja do formy tożsamości
	for currentRow := 1; currentRow <= rows; currentRow++ {

		// Obrotowe
		p := currentRow
		for i := currentRow + 1; i <= rows; i++ {
			if math.Abs(m[i-1][currentRow-1]) > math.Abs(m[p-1][currentRow-1]) {
				p = i
			}
		}

		// Jeśli nie ma elementu a (k, i) różnego od zera, macierz jest pojedyncza i nie ma żadnego lub więcej niż jednego rozwiązania
		if math.Abs(m[p-1][currentRow-1]) < 1e-10 {
			return nil, fmt.Errorf("Macierz jest pojedyncza")
		}

		// Jeśli znajdziemy pivot, który jest największym a (i, currentRow), zamieniamy wiersze
		if p != currentRow {
			tmp := m[currentRow-1]
			m[currentRow-1] = m[p-1]
			m[p-1] = tmp
		}

		vec[currentRow-1] = float64(p)

		mi := m[currentRow-1][currentRow-1]
		m[currentRow-1][currentRow-1] = 1.0

		// Dzielenie przez mi
		div, err := m[currentRow-1].DivideByScalar(mi)

		if err != nil {
			return nil, err
		}

		m[currentRow-1] = div

		for i := 1; i <= rows; i++ {
			if i != currentRow {
				mi = m[i-1][currentRow-1]
				m[i-1][currentRow-1] = 0.0
				for j := 1; j <= rows; j++ {
					m[i-1][j-1] -= mi * m[currentRow-1][j-1]
				}
			}
		}
	}

	// Odwrotna zamiana
	for j := rows; j >= 1; j-- {
		p := vec[j-1]
		if p != float64(j) {
			for i := 1; i <= rows; i++ {
				tmp := m[i-1][int64(p)-1]
				m[i-1][int64(p)-1] = m[i-1][j-1]
				m[i-1][j-1] = tmp
			}
		}
	}
	return m, nil
}

// Log stosuje logarytm naturalny do wszystkich elementów macierzy i zwraca wynikową macierz.
func (m Matrix) Log() Matrix {
	row, col := m.Dim()
	result := make(Matrix, row)
	for i := range m {
		result[i] = make(Vector, col)
		for j := range m[i] {
			result[i][j] = math.Log(m[i][j])
		}
	}
	return result
}

// Exp stosuje e^x do wszystkich elementów macierzy i zwraca wynikową macierz.
func (m Matrix) Exp() Matrix {
	row, col := m.Dim()
	result := make(Matrix, row)
	for i := range m {
		result[i] = make(Vector, col)
		for j := range m[i] {
			result[i][j] = math.Exp(m[i][j])
		}
	}
	return result
}

// LeftDivide otrzymuje inną macierz jako parametr.
// Metoda rozwiązuje symboliczny układ równań liniowych w postaci macierzowej, A * X = B dla X.
// Zwraca wyniki w postaci macierzowej i błędu (jeśli istnieje).
func (m Matrix) LeftDivide(m2 Matrix) (Matrix, error) {
	var r Matrix
	mT, err := m.Transpose()

	if err != nil {
		return r, err
	}

	mtm, err := mT.MultiplyBy(m)

	if err != nil {
		return r, err
	}

	pInv, err := mtm.Invert()
	if err != nil {
		return r, err
	}

	pmt, err := pInv.MultiplyBy(mT)

	if err != nil {
		return r, err
	}

	r, err = pmt.MultiplyBy(m2)

	if err != nil {
		return r, err
	}

	return r, nil
}

func (m Matrix) sumAbs() float64 {
	var sum float64
	rows, cols := m.Dim()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			sum += math.Abs(m[i][j])
		}
	}
	return sum
}

func (m Matrix) isSquare() bool {
	rows, cols := m.Dim()
	return rows == cols
}

// MultiplyBy otrzymuje inną macierz jako parametr.
// Mnoży macierze i zwraca wynikową macierz i błąd.
func (m Matrix) MultiplyBy(m2 Matrix) (Matrix, error) {
	var r Matrix

	_, cols1 := m.Dim()
	rows2, _ := m2.Dim()

	if cols1 != rows2 {
		return r, fmt.Errorf("Liczba kolumn pierwszej macierzy musi być równa liczbie wierszy drugiej macierzy")
	}

	for currentRowIndex := range m {
		r = append(r, Vector{})
		for currentColumnIndex := range m2[0] {
			m2Col, err := m2.Col(currentColumnIndex)

			if err != nil {
				return r, err
			}

			dot, err := m[currentRowIndex].Dot(m2Col)
			if err != nil {
				return r, err
			}

			r[currentRowIndex] = append(r[currentRowIndex], dot)
		}
	}

	return r, nil
}

// InsertCol otrzymuje indeks i wektor.
// Dodaje podany wektor jako kolumnę o indeksie k i zwraca wynikową macierz oraz błąd (jeśli istnieje).
func (m Matrix) InsertCol(k int, c Vector) (Matrix, error) {
	var r Matrix

	if k < 0 {
		return r, fmt.Errorf("Indeks nie może być mniejszy niż 0")
	} else if _, width := m.Dim(); k > width {
		return r, fmt.Errorf("Indeks nie może być większy niż liczba kolumn + 1")
	} else if len(c) != len(m) {
		return r, fmt.Errorf("Wymiary kolumn muszą być zgodne")
	}

	for i := 0; i < len(m); i++ {
		row := m[i]
		expRow := append(row[:k], append(Vector{c[i]}, row[k:]...)...)
		r = append(r, expRow)
	}

	return r, nil
}

// Row otrzymuje indeks jako parametr.
// Zwraca wiersz wektora pod podanym indeksem i błąd (jeśli istnieje).
func (m Matrix) Row(i int) (Vector, error) {
	if i < 0 {
		return nil, fmt.Errorf("Indeks nie może być ujemny")
	} else if i > len(m) {
		return nil, fmt.Errorf("Indeks nie może być większy niż długość")
	}
	return m[i], nil
}

// Col otrzymuje indeks jako parametr.
// Zwraca kolumnę wektora pod podanym indeksem i błąd (jeśli istnieje).
func (m Matrix) Col(i int) (Vector, error) {
	var r Vector

	if i < 0 {
		return nil, fmt.Errorf("Indeks nie może być ujemny")
	} else if i > len(m[0]) {
		return nil, fmt.Errorf("Indeks nie może być większy niż długość")
	}

	for row := range m {
		r = append(r, m[row][i])
	}

	return r, nil
}

// Transpozycja zwraca transponowaną macierz i błąd.
func (m Matrix) Transpose() (Matrix, error) {
	var t Matrix

	for columnIndex := range m[0] {
		column, err := m.Col(columnIndex)
		if err != nil {
			return t, err
		}
		t = append(t, column)
	}

	return t, nil
}

// IsSimilar otrzymuje inną macierz i tolerancję jako parametry.
// Sprawdza, czy dwie macierze są podobne w ramach podanej tolerancji.
func (m Matrix) IsSimilar(m2 Matrix, tol float64) bool {

	if m.IsEqual(m2) {
		return true
	}

	if !m.areDimsEqual(m2) {
		return false
	}

	for col := range m {
		for row := range m[col] {
			if math.Abs(m[col][row]-m2[col][row]) > tol {
				return false
			}
		}
	}

	return true
}

// IsEqual otrzymuje inną macierz jako parametr.
// Zwraca prawdę, jeśli wartości dwóch macierzy są równe, lub fałsz w przeciwnym razie.
func (m Matrix) IsEqual(m2 Matrix) bool {
	if m == nil && m2 == nil {
		return true
	} else if m == nil || m2 == nil {
		return false
	} else if !m.areDimsEqual(m2) {
		return false
	}

	for row := range m {
		for col := range m[row] {
			if m[row][col] != m2[row][col] {
				return false
			}
		}
	}
	return true
}

// Add otrzymuje inną macierz jako parametr.
// Dodaje dwie macierze i zwraca macierz wyników oraz błąd (jeśli taki istnieje).
func (m Matrix) Add(m2 Matrix) (Matrix, error) {
	rows, cols := m.Dim()
	var r = make(Matrix, rows, cols)
	if ok, err := m.canPerformOperationsWith(m2); !ok {
		return nil, err
	}

	for row := range m {
		for col := range m[row] {
			r[row] = append(r[row], m[row][col]+m2[row][col])
		}
	}

	return r, nil
}

// Subtract otrzymuje inną macierz jako parametr.
// Odejmuje dwie macierze i zwraca macierz wyników oraz błąd (jeśli istnieje).
func (m Matrix) Subtract(m2 Matrix) (Matrix, error) {
	rows, cols := m.Dim()
	var r = make(Matrix, rows, cols)
	if ok, err := m.canPerformOperationsWith(m2); !ok {
		return nil, err
	}

	for row := range m {
		for col := range m[row] {
			r[row] = append(r[row], m[row][col]-m2[row][col])
		}
	}

	return r, nil
}

func (m Matrix) areDimsEqual(m2 Matrix) bool {
	mRows, mCols := m.Dim()
	m2Rows, m2Cols := m2.Dim()

	if mRows != m2Rows || mCols != m2Cols {
		return false
	}
	return true
}

func (m Matrix) isNil() bool {
	if m == nil {
		return true
	}
	return false
}

func (m Matrix) canPerformOperationsWith(m2 Matrix) (bool, error) {
	if m == nil || m2 == nil {
		return false, fmt.Errorf("Macierze nie mogą być <nil>")
	} else if !m.areDimsEqual(m2) {
		return false, fmt.Errorf("Wymiary macierzy muszą być zgodne")
	}
	return true, nil
}
