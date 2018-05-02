// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gov

import (
	"regexp"

	"strings"

	"github.com/8thlight/sai_watcher/utils"
	"github.com/ethereum/go-ethereum/common"
)

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

// this is how web3.utils.hexToUtf8 makes a hex decodable without trailing characters
func RemoveZeroPadding(s string) string {
	var re = regexp.MustCompile(`(?:00)*`)
	s = re.ReplaceAllString(s, ``)
	s = re.ReplaceAllString(reverse(s), ``)
	return reverse(s)
}

func ConvertToModel(ge *GovEntity) *GovModel {
	var_ := string(common.FromHex(RemoveZeroPadding(ge.Var)))
	arg := utils.Convert("wad", ge.Arg, 16)
	guy := strings.ToLower(common.HexToAddress(ge.Guy).Hex())
	cap_ := utils.Convert("wad", ge.Cap, 16)
	mat := utils.Convert("ray", ge.Mat, 16)
	tax := utils.Convert("ray", ge.Tax, 16)
	fee := utils.Convert("ray", ge.Fee, 16)
	axe := utils.Convert("ray", ge.Axe, 16)
	gap := utils.Convert("wad", ge.Gap, 16)

	return &GovModel{
		Block: ge.Block,
		Tx:    ge.Tx,
		Var:   var_,
		Arg:   arg,
		Guy:   guy,
		Cap:   cap_,
		Mat:   mat,
		Tax:   tax,
		Fee:   fee,
		Axe:   axe,
		Gap:   gap,
	}
}
