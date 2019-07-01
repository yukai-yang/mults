package mults

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

/* utility functions */

func makerows(rows [][]int, start []int, iTT, freq int) {
	rows[0] = []int{start[0], start[1]}
	for i := 1; i < iTT; i++ {
		rows[i] = make([]int, 2)
		rows[i][0] = rows[i-1][0]
		rows[i][1] = rows[i-1][1] + 1
		if rows[i][1] == freq {
			rows[i][1]--
			rows[i][0]++
		}
	}
}

func makerowsrev(rows [][]int, end []int, iTT, freq int) {
	rows[iTT-1] = []int{end[0], end[1]}
	for i := iTT - 2; i >= 0; i-- {
		rows[i] = make([]int, 2)
		rows[i][0] = rows[i+1][0]
		rows[i][1] = rows[i+1][1] - 1
		if rows[i][1] < 0 {
			rows[i][1]++
			rows[i][0]--
		}
	}

}

func containsint(xs []int, x int) bool {
	for _, v := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// make a copy to the origin
func depvars(ts *MulTS, from, to int) []float64 {
	var dep = []float64{}
	var tmp []float64
	for _, v := range ts.dep {
		tmp = make([]float64, to-from)
		copy(tmp, ts.data[v][from:to])
		dep = append(dep, tmp...)
	}
	return dep
}

// ViewMatrix prints the matrix
func ViewMatrix(m mat.Matrix) string {
	var r, c = m.Dims()
	var str string
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			str = str + fmt.Sprintf("%v ", m.At(i, j))
		}
		str = str + "\n"
	}

	return str
}
