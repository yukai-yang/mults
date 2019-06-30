package mults

import (
	"errors"
)

// MulTS is a struct for the multivariate time series
type MulTS struct {
	data      [][]float64
	freq      int
	start     [2]int
	end       [2]int
	vnames    []string
	laglength int
	dep       []int
	indep     []int
}

// SetData sets the data, start, end and frequency
// data is a slice of float64, row first sorted
// nvar is the number of variables
// vnames contains the names of the variables, nil for no names
func (ts *MulTS) SetData(data []float64, nvar int, vnames []string) error {
	var iTT = len(data) / nvar
	if iTT*nvar != len(data) {
		return errors.New("dimension does not fit")
	}
	ts.data = make([][]float64, nvar)
	for i := 0; i < nvar; i++ {
		ts.data[i] = make([]float64, iTT)
		copy(ts.data[i], data[i*iTT:(i*iTT+iTT)])
	}
	if vnames != nil && len(vnames) == nvar {
		ts.vnames = make([]string, nvar)
		copy(ts.vnames, vnames)
	}
	return nil
}

// SetFreq sets the frequency of the time series
// start and end are 2-slices of int representing the start and end time points
// if start and/or end are nil, they will be inferred
func (ts *MulTS) SetFreq(freq int, start, end [2]int) error {
	if freq < 1 {
		return errors.New("wrong frequency")
	}
	ts.freq = freq

	return nil
}
