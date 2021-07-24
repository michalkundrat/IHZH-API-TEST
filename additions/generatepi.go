package additions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big" // high-precision math
)

func arccot(x int64, unity *big.Int) *big.Int {
	bigx := big.NewInt(x)
	xsquared := big.NewInt(x * x)
	sum := big.NewInt(0)
	sum.Div(unity, bigx)
	xpower := big.NewInt(0)
	xpower.Set(sum)
	n := int64(3)
	zero := big.NewInt(0)
	sign := false

	term := big.NewInt(0)
	for {
		xpower.Div(xpower, xsquared)
		term.Div(xpower, big.NewInt(n))
		if term.Cmp(zero) == 0 {
			break
		}
		if sign {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}
		sign = !sign
		n += 2
	}
	return sum
}

func generatePiNumber() {
	ndigits := int64(5000)
	digits := big.NewInt(ndigits + 10)
	unity := big.NewInt(0)
	unity.Exp(big.NewInt(10), digits, nil)
	pi := big.NewInt(0)
	four := big.NewInt(4)
	pi.Mul(four, pi.Sub(pi.Mul(four, arccot(5, unity)), arccot(239, unity)))
	//val := big.Mul(4, big.Sub(big.Mul(4, arccot(5, unity)), arccot(239, unity)))
	pistring := pi.String()[0:ndigits]
	fmt.Println("Computed pi: ", pistring)
	digitcount := make([]int, 10)
	for _, digit := range pistring {
		val := digit - '0'
		digitcount[val]++
	}
	fmt.Printf("Digit\tCount\n")
	for i, digit := range digitcount {
		fmt.Printf("%d\t%d\n", i, digit)
	}
	file, _ := json.MarshalIndent(pistring, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)
}
