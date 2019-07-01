package mults

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestMulTS(t *testing.T) {
	data := make([]float64, 40)
	for i := range data {
		data[i] = rand.NormFloat64()
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 4; j++ {
			fmt.Print(data[i+j*10], "\t")
		}
		fmt.Println()
	}
	fmt.Println()

	mults := &MulTS{}
	mults.SetData(data, 4, nil)
	mults.SetFreq(4, nil, nil)
	mults.SetNames([]string{"V0", "V1", "V2", "V3"})
	mults.SetLag(2)

	mults.SetDepByCol(false, 0)
	mults.SetDepByName(true, "V1")

	mults.SetIndepByCol(false, 2)
	mults.SetIndepByName(true, "V3")

	var dep, _ = mults.DepVars(2, 10)
	fmt.Println(ViewMatrix(dep))
	var indep, _ = mults.IndepVars(2, 10)
	fmt.Println(ViewMatrix(indep))
}
