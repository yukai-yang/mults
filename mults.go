package mults

import (
	"errors"
)

// MulTS is a struct for the multivariate time series
type MulTS struct {
	data      [][]float64
	iTT       int
	freq      int
	rows      [][]int
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
	ts.iTT = iTT
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
func (ts *MulTS) SetFreq(freq int, start, end []int) error {
	if ts.data == nil {
		return errors.New("no data")
	}
	if freq < 1 {
		return errors.New("wrong frequency")
	}

	if start != nil {
		if len(start) < 2 {
			return errors.New("start has less than two values")
		}
		if start[1] < 0 || start[1] >= freq {
			return errors.New("start's freq is invalid")
		}

		ts.rows = make([][]int, ts.iTT)
		makerows(ts.rows, start, ts.iTT, freq)

	} else {
		if end != nil {
			if len(end) < 2 {
				return errors.New("end has less than two values")
			}
			if end[1] < 0 || end[1] >= freq {
				return errors.New("end's freq is invalid")
			}

			ts.rows = make([][]int, ts.iTT)
			makerowsrev(ts.rows, end, ts.iTT, freq)
		} else {
			ts.rows = make([][]int, ts.iTT)
			makerows(ts.rows, []int{0, 0}, ts.iTT, freq)
		}
	}

	ts.freq = freq

	return nil
}
