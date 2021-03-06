// Copyright 2019 Infinity Stones
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

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Flags
var (
	blockComm           int64
	epochComm           int64
	foundationComm      int64
	outputFile          string
	epochToQuery        string
	simpleJson          bool
)

var PayoutCmd = &cobra.Command{
	Use:   "iotex_payout DELEGATE_NAME OPERATOR_[ALIAS|ADDRESS]",
	Short: "Calculates voters' reward shares for IOTEX blockchain, output the input for iotex multisend",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		output := payout(args[0], args[1])
		if outputFile == "" {
			fmt.Println(output)
			return
		}
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		if _, err := f.Write([]byte(output + "\n")); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
	},
}

func main() {
	if err := PayoutCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	PayoutCmd.Flags().Int64VarP(&blockComm, "block-commission", "b", 100,
		"commission rate of block reward, 100% by default")
	PayoutCmd.Flags().Int64VarP(&epochComm, "epoch-commission", "p", 100,
		"commission rate of epoch bonus, 100% by default")
	PayoutCmd.Flags().Int64VarP(&foundationComm, "foundation-commission", "f", 100,
		"commission rate of foundation bonus, 100% by default")
	PayoutCmd.Flags().StringVarP(&outputFile, "output", "o", "",
		"file to output the result, output to stdout by default")
	PayoutCmd.Flags().StringVarP(&epochToQuery, "epoch", "e", "",
		"epoch(s) to calculate rewards, current epoch by default. " +
		"The input is in range format (e.g. 1-2,4,7-10)")
	PayoutCmd.Flags().BoolVarP(&simpleJson, "simple", "s", false,
		"also print out votes information, print rewards only by default")

	if blockComm > 100 {
		fmt.Println("valid value for block reward commission rate is up to 100")
		os.Exit(2)
	}
	if epochComm > 100 {
		fmt.Println("valid value for epoch reward commission rate is up to 100")
		os.Exit(2)
	}
	if foundationComm > 100 {
		fmt.Println("valid value for extra reward commission rate is up to 100")
		os.Exit(2)
	}
}
