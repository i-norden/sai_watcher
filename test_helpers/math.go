package test_helpers

import "math/big"

func GetFloat(n string) *big.Float {
	strAsFloat := big.NewFloat(0)
	strAsFloat.SetString(n)
	return strAsFloat
}
