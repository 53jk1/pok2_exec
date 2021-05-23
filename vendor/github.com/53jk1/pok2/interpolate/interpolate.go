package interpolate

type interpolator interface {
	Interpolate(float64) float64
}

type validator interface {
	Validate(float64) error
}

type validateInterpolator interface {
	interpolator
	validator
}

// WithMulti akceptuje wycinek float64 i zwraca interpolowane wartości dla przekazanych wartości wycinka, a błąd
func WithMulti(vi validateInterpolator, vals []float64) ([]float64, error) {
	var r []float64
	for _, val := range vals {
		est, err := WithSingle(vi, val)
		if err != nil {
			return r, err
		}
		r = append(r, est)
	}
	return r, nil
}

// WithSingle akceptuje pojedynczą wartość float64 i zwraca dla niej interpolowaną wartość oraz błąd
func WithSingle(vi validateInterpolator, val float64) (float64, error) {
	var est float64

	err := vi.Validate(val)
	if err != nil {
		return est, err
	}

	est = vi.Interpolate(val)
	return est, nil
}
