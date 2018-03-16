package utils

import (
	"math/big"
	"path"
	"runtime"
)

var Ether *big.Float
var Ray *big.Float

func init() {
	Ether = big.NewFloat(1e18)
	Ray = big.NewFloat(1e27)
}

func ProjectRoot() string {
	var _, filename, _, _ = runtime.Caller(0)
	return path.Join(path.Dir(filename), "..")
}

func Convert(conversion string, value string, prec int) string {
	var bgflt = big.NewFloat(0.0)
	bgflt.SetString(value)
	switch conversion {
	case "ray":
		bgflt.Quo(bgflt, Ray)
	case "wad":
		bgflt.Quo(bgflt, Ether)
	}
	return bgflt.Text('g', prec)
}

func Arg(s string) string {
	if len(s) < 27 {
		return ""
	}
	val := "0x" + s[26:]
	switch {
	case val == "0x0000000000000000000000000000000000000000":
		return ""
	default:
		return Convert("wad", val, 16)
	}
}
