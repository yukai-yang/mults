package mults

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

func contains(xs []int, x int) bool {
	for v := range xs {
		if x == v {
			return true
		}
	}
	return false
}
