package mults

import "errors"

// MulTS is a struct for the multivariate time series
type MulTS struct {
	data   [][]float64 // SetData
	iTT    int         // SetData
	freq   int         // SetFreq
	rows   [][]int     // SetFreq
	vnames []string    // SetData, SetNames
	lag    int         // SetLag
	dep    []int       // SetDepByCol, SetDepByName
	indep  []int       // SetIndepByCol, SetIndepByName
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
	ts.indep = []int{}

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
func (ts *MulTS) SetDepByCol(deps []int, app bool) error {
	var nvar = len(ts.data)
	if !app {
		// not append
		ts.dep = []int{}
	}

	for dep := range deps {
		if dep < 0 || dep >= nvar {
			continue
		}

		if !containsint(ts.dep, dep) {
			ts.dep = append(ts.dep, dep)
		}
	}

	return nil
}

// SetDepByName appends dependent variables by variable names
func (ts *MulTS) SetDepByName(deps []string, app bool) error {
	if ts.vnames == nil {
		return errors.New("variables have no names")
	}

	if !app {
		// not append
		ts.dep = []int{}
	}

	for _, dep := range deps {
		for _, v := range ts.dep {
			if dep == ts.vnames[v] && !containsint(ts.dep, v) {
				ts.dep = append(ts.dep, v)
				break
			}
		}
	}

	return nil
}

// SetIndepByCol sets independent variables by column numbers
func (ts *MulTS) SetIndepByCol(indeps []int, app bool) error {
	var nvar = len(ts.data)
	if !app {
		// not append
		ts.indep = []int{}
	}

	for indep := range indeps {
		if indep < 0 || indep >= nvar {
			continue
		}

		if !containsint(ts.indep, indep) {
			ts.indep = append(ts.indep, indep)
		}
	}

	return nil
}

// SetIndepByName appends independent variables by variable names
func (ts *MulTS) SetIndepByName(indeps []string, app bool) error {
	if ts.vnames == nil {
		return errors.New("variables have no names")
	}

	if !app {
		// not append
		ts.indep = []int{}
	}

	for _, indep := range indeps {
		for _, v := range ts.indep {
			if indep == ts.vnames[v] && !containsint(ts.indep, v) {
				ts.indep = append(ts.indep, v)
				break
			}
		}
	}

	return nil
}

// DepVars returns a matrix containing the dependent variables
// subset from "from" to but without "to"
func (ts MulTS) DepVars(from, to int) ([][]float64, error) {
	if ts.data == nil {
		return nil, errors.New("no data")
	}
	if len(ts.dep) == 0 {
		return nil, errors.New("no dependent variable")
	}
	if from < 0 || from >= to || to > ts.iTT {
		return nil, errors.New("invalid from or to")
	}

	var dep = make([][]float64, len(ts.dep))
	for i, v := range ts.dep {
		dep[i] = make([]float64, to-from)
		copy(dep[i], ts.data[v][from:to])
	}

	return dep, nil
}

// IndepVars returns a matrix containing the independent variables
func (ts MulTS) IndepVars(from, to int) ([][]float64, error) {
	if ts.data == nil {
		return nil, errors.New("no data")
	}
	if len(ts.indep)+ts.lag*len(ts.dep) == 0 {
		return nil, errors.New("no independent variable")
	}
	if from-ts.lag < 0 || from > to || to >= ts.iTT {
		return nil, errors.New("invalid from or to")
	}

	var indep = [][]float64{}
	var tmp [][]float64
	for k := 1; k <= ts.lag; k++ {
		// DepVars copies the values
		tmp, _ = ts.DepVars(from-k, to-k)
		indep = append(indep, tmp...)
	}

	var tmpp []float64
	for _, v := range ts.indep {
		tmpp = make([]float64, to-from)
		copy(tmpp, ts.data[v][from:to])
		indep = append(indep, tmpp)
	}

	return indep, nil
}
