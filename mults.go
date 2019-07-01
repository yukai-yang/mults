package mults

import (
	"errors"

	"gonum.org/v1/gonum/mat"
)

// MulTS is a struct for the multivariate time series
type MulTS struct {
	data   [][]float64 // SetData
	iTT    int         // SetData
	freq   int         // SetFreq
	rows   [][]int     // SetFreq
	vnames []string    // SetData, SetNames
	lag    int         // SetLag
	dep    []int       // SetData, SetDepByCol, SetDepByName
	indep  [][]int     // SetData, SetIndepByCol, SetIndepByName
}

// SetData sets the data, start, end and frequency
// data is a slice of float64, row first sorted
// nvar is the number of variables
// vnames contains the names of the variables, nil for no names
func (ts *MulTS) SetData(data []float64, nvar int, vnames []string) error {
	var iTT = len(data) / nvar
	if iTT*nvar != len(data) {
		return errors.New("data quantity does not fit")
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
	ts.dep = []int{}
	ts.indep = [][]int{}

	return nil
}

// SetNames sets the names of the variables
func (ts *MulTS) SetNames(vnames []string) error {
	if vnames == nil {
		return errors.New("names empty")
	}
	if ts.data == nil {
		return errors.New("no data")
	}
	if len(vnames) != len(ts.data) {
		return errors.New("variables number does not fit")
	}

	ts.vnames = make([]string, len(ts.data))
	copy(ts.vnames, vnames)

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

// SetLag sets the lag length
func (ts *MulTS) SetLag(k int) error {
	if k < 0 || k >= ts.iTT {
		return errors.New("invalid lag length")
	}
	ts.lag = k

	return nil
}

// SetDepByCol sets dependent variables by column numbers
func (ts *MulTS) SetDepByCol(app bool, deps ...int) error {
	if !app {
		// not append
		ts.dep = []int{}
	}

	for _, dep := range deps {
		if dep < 0 || dep >= len(ts.data) {
			continue
		}

		if !contains(ts.dep, dep) {
			ts.dep = append(ts.dep, dep)
		}
	}

	return nil
}

// SetDepByName appends dependent variables by variable names
func (ts *MulTS) SetDepByName(app bool, deps ...string) error {
	if ts.vnames == nil {
		return errors.New("variables have no names")
	}

	var tmp []int
	for _, dep := range deps {
		for i, v := range ts.vnames {
			if dep == v {
				tmp = append(tmp, i)
				break
			}
		}
	}

	ts.SetDepByCol(app, tmp...)

	return nil
}

// SetIndepByCol sets independent variables by column numbers
func (ts *MulTS) SetIndepByCol(app bool, indep int, lag int) error {
	if indep < 0 || indep >= len(ts.data) {
		return errors.New("invalid column number")
	}
	if lag < 0 || lag >= ts.iTT {
		return errors.New("invalid lag")
	}

	if !app {
		// not append
		ts.indep = [][]int{}
	}

	var tmp = []int{indep, lag}
	if !containpair(ts.indep, tmp) {
		ts.indep = append(ts.indep, tmp)
	}

	return nil
}

// SetIndepByName appends independent variables by variable names
func (ts *MulTS) SetIndepByName(app bool, indep string, lag int) error {
	if ts.vnames == nil {
		return errors.New("variables have no names")
	}
	if lag < 0 || lag >= ts.iTT {
		return errors.New("invalid lag")
	}

	for i, v := range ts.vnames {
		if indep == v {
			ts.SetIndepByCol(app, i, lag)
			break
		}
	}

	return nil
}

// DepVars returns (copies) a matrix containing the dependent variables
// subset from "from" to but without "to"
func (ts *MulTS) DepVars(from, to int) (mat.Matrix, error) {
	if ts.data == nil {
		return nil, errors.New("no data")
	}
	if len(ts.dep) == 0 {
		return nil, errors.New("no dependent variable")
	}
	mfrom, mto := ts.PossibleFrame()
	if from < mfrom || from >= to || to > mto {
		return nil, errors.New("invalid from or to")
	}

	den := mat.NewDense(len(ts.dep), to-from, depvars(ts, from, to))

	return den.T(), nil
}

// IndepVars returns a matrix containing the independent variables
func (ts *MulTS) IndepVars(from, to int) (mat.Matrix, error) {
	if ts.data == nil {
		return nil, errors.New("no data")
	}
	if len(ts.indep)+ts.lag*len(ts.dep) == 0 {
		return nil, errors.New("no independent variable")
	}
	mfrom, mto := ts.PossibleFrame()
	if from < mfrom || from >= to || to > mto {
		return nil, errors.New("invalid from or to")
	}

	var indep = []float64{}
	for k := 1; k <= ts.lag; k++ {
		// depvars copies the values
		indep = append(indep, depvars(ts, from-k, to-k)...)
	}

	var tmp []float64
	for _, v := range ts.indep {
		tmp = make([]float64, to-from)
		copy(tmp, ts.data[v[0]][(from-v[1]):(to-v[1])])
		indep = append(indep, tmp...)
	}

	den := mat.NewDense(len(ts.indep)+ts.lag*len(ts.dep), to-from, indep)

	return den.T(), nil
}

// PossibleFrame returns the possible time frame from and to with the largest sample size
func (ts *MulTS) PossibleFrame() (int, int) {
	var from = ts.lag
	for _, v := range ts.indep {
		if v[1] > from {
			from = v[1]
		}
	}
	return from, ts.iTT
}
