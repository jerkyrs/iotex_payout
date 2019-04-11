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
	blockComm           uint64
	epochComm           uint64
	foundationComm      uint64
	outputFile          string
	epochToQuery        string
)

var PayoutCmd = &cobra.Command{
	Use:   "iotex_payout DELEGATE_NAME OPERATOR_ALIAS",
	Short: "Calculates voters' reward shares for IOTEX blockchain",
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
	PayoutCmd.Flags().Uint64VarP(&blockComm, "block-commission", "b", 10,
		"commission rate of block reward, 10 percent by default")
	PayoutCmd.Flags().Uint64VarP(&epochComm, "epoch-commission", "p", 10,
		"commission rate of epoch bonus, 10 percent by default")
	PayoutCmd.Flags().Uint64VarP(&foundationComm, "foundation-commission", "f", 10,
		"commission rate of foundation bonus, 10 percent by default")
	PayoutCmd.Flags().StringVarP(&outputFile, "output", "o", "",
		"file to output the result, output to stdout by default")
	PayoutCmd.Flags().StringVarP(&epochToQuery, "epoch", "e", "",
		"epoch(s) to calculate rewards, current epoch by default. " +
		"The input is in range format (e.g. 1-2,4,7-10)")

	if blockComm > 100 {
		fmt.Println("valid value for block reward commission rate is from 0 to 100")
		os.Exit(2)
	}
	if epochComm > 100 {
		fmt.Println("valid value for epoch reward commission rate is from 0 to 100")
		os.Exit(2)
	}
	if foundationComm > 100 {
		fmt.Println("valid value for extra reward commission rate is from 0 to 100")
		os.Exit(2)
	}
}
