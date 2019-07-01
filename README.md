# mults [![GoDoc](https://godoc.org/github.com/yukai-yang/mults?status.svg)](https://godoc.org/github.com/yukai-yang/mults)

Package mults implements the container for multivariate time series data.

## How to use the package

```go
package main

import (
	"fmt"
	"math/rand"
	"github.com/yukai-yang/mults"
)

func main() {
	data := make([]float64, 40)
	for i := range data {
		data[i] = rand.NormFloat64()
	}

	ts := &mults.MulTS{}
	ts.SetData(data, 4, nil)
	ts.SetFreq(4, nil, nil)
	ts.SetNames([]string{"V0", "V1", "V2", "V3"})
	ts.SetLag(2)

	ts.SetDepByCol(false, 0)
	ts.SetDepByName(true, "V1")

	ts.SetIndepByCol(false, 2, 0)
	ts.SetIndepByName(true, "V3", 2)

	var dep, _ = ts.DepVars(2, 10)
	fmt.Println(mults.ViewMatrix(dep))
	var indep, _ = ts.IndepVars(2, 10)
	fmt.Println(mults.ViewMatrix(indep))
}
```
