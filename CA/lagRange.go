package main

import "math/big"

func LagRange(a []*LagPoint, x *big.Int) *big.Int {
    result := big.NewInt(0)
    for i, lP := range a {
        up := big.NewInt(1)
        down := big.NewInt(1)
        for j, lp := range a {
            if j != i {
                up.Mul(up, new(big.Int).Sub(x, lp.X))
                down.Mul(down, new(big.Int).Sub(lP.X, lp.X))
            }
        }
        LxResult := new(big.Int).Div(up, down)
        result.Add(result, new(big.Int).Mul(lP.Y, LxResult))
    }
    return result
}
