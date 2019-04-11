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
	"context"
	"encoding/hex"
	"strings"
	"strconv"
	"fmt"
	"math/big"

	"github.com/iotexproject/iotex-core/cli/ioctl/cmd/alias"
	"github.com/iotexproject/iotex-core/cli/ioctl/cmd/bc"
	"github.com/iotexproject/iotex-core/cli/ioctl/util"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
	"github.com/iotexproject/iotex-election/committee"
)

// Config needed by committee
//
// TODO: read config from yaml file
var config committee.Config = committee.Config{
	8, // NumOfRetries   uint8
	[]string{"https://mainnet.infura.io/v3/b355cae6fafc4302b106b937ee6c15af"},
	// gravityChainAPIs []string
	100,     // gravityChainHeightInterval  uint64
	7368630, // gravityChainStartHeight     uint64
	"0x95724986563028deb58f15c5fac19fa09304f32d", // RegisterContractAddress   string
	"0x87c9dbff0016af23f5b1ab9b8e072124ab729193", // StakingContractAddress    string
	100, // PaginationSize            uint8
	"0", // VoteThreshold             string
	"0", // ScoreThreshold            string
	"0", // SelfStakingThreshold      string
	100, // CacheSize                 uint32
	1}   // NumOfFetchInParallel      uint8

var epochResponse *iotexapi.GetEpochMetaResponse

// get voter's votes
func getVotes(delegate []byte, height uint64) (map[string]*big.Int, *big.Int, *big.Int) {
	comm, err := committee.NewCommittee(nil, config)
	if err != nil {
		panic(err)
	}
	result, err := comm.FetchResultByHeight(height)
	if err != nil {
		panic(err)
	}

	delegateVotes := new(big.Int)
	bps := make(map[string]*big.Int)

	for _, vote := range result.VotesByDelegate(delegate) {
		ethAddr := hex.EncodeToString(vote.Voter()) // []byte
		votes := vote.WeightedAmount()              // *big.Int
		_, ok := bps[ethAddr]
		if ok {
			// Already have this eth addr, need to combine the votes
			bps[ethAddr] = new(big.Int).Add(bps[ethAddr], votes)
		} else {
			bps[ethAddr] = votes
		}
		delegateVotes = delegateVotes.Add(delegateVotes, votes)
	}
	return bps, delegateVotes, result.TotalVotes()
}

// get current epoch
func currentEpochNum() uint64 {
	chainMeta, err := bc.GetChainMeta()
	if err != nil {
		panic(err)
	}
	return chainMeta.Epoch.Num
}

// calculate rewards
func calculateReward(blks uint64, elected bool, votes *big.Int, total *big.Int) Reward {
	var reward Reward

	// block reward
	reward.Block = 16 * blks

	// epoch reward
	if elected {
		reward.FoundationBonus = 1920
	}

	// bonus reward
	extra, _ := new(big.Int).SetString("300000", 10)
	extra = extra.Mul(extra, votes)
	extra = extra.Div(extra, total)
	reward.EpochBonus = extra.Uint64()

	return reward
}

func getEpochResponse(epoch_num uint64) {
	conn, err := util.ConnectToEndpoint()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cli := iotexapi.NewAPIServiceClient(conn)
	request := &iotexapi.GetEpochMetaRequest{EpochNumber: epoch_num}
	ctx := context.Background()
	epochResponse, err = cli.GetEpochMeta(ctx, request)
	if err != nil {
		panic(err)
	}
}

// get epoch data
func epochNum() uint64 {
	return epochResponse.GetEpochData().GetNum()
}

func epochHeight() uint64 {
	return epochResponse.GetEpochData().GetHeight()
}

func epochGravityHeight() uint64 {
	return epochResponse.GetEpochData().GetGravityChainStartHeight()
}

// get number of produced blocks
func delegateProductivity(operator string) uint64 {
	for _, bp := range epochResponse.GetBlockProducersInfo() {
		if operator == bp.GetAddress() {
			return bp.GetProduction()
		}
	}
	return 0
}

// whether the delegate was elected
func isDelegateElected(operator string) bool {
	for _, bp := range epochResponse.GetBlockProducersInfo() {
		if operator == bp.GetAddress() {
			return true
		}
	}
	return false
}

// get delegate's name as byte array
func delegateName(delegate string) []byte {
	length := len(delegate)
	if length > 12 {
		return []byte(delegate[length-12:])
	} else {
		return append(make([]byte, 12-length), []byte(delegate)...)
	}
}

// populate reward shares for a single epoch
func calculateEpochRewardShares(operator string, delegate []byte, epoch_num uint64) *RewardShares {
	// get epoch response
	getEpochResponse(epoch_num)

	// get gravity height
	gravity_height := epochGravityHeight()

	// get number of produced blocks
	blocks := delegateProductivity(operator)

	// delegate is not in the BP list, no epoch reward
	elected := isDelegateElected(operator)

	// get delegate's votes
	votes_distribution, delegate_votes, total_votes := getVotes(delegate, gravity_height)

	// calculate reward
	reward := calculateReward(blocks, elected, delegate_votes, total_votes)

	// populate rewardshare structure
	return NewRewardShares().
		SetEpochNum(epoch_num).
		SetProductivity(blocks).
		SetTotalVotes(total_votes).
		SetReward(reward).
		CalculateShares(votes_distribution, delegate_votes)
}

// Generator for epoch range
func epochRangeGen(epochs string) chan uint64 {
	c := make(chan uint64)

	// string parser for extracting epoch range
	go func(input string) {
		for _, part := range strings.Split(input, ",") {
			if i := strings.Index(part[1:], "-"); i == -1 {
				n, err := strconv.ParseUint(part, 10, 64)
				if err != nil {
					fmt.Println(err)
					break
				}
				c <- n
			} else {
				n1, err := strconv.ParseUint(part[:i+1], 10, 64)
				if err != nil {
					fmt.Println(err)
					break
				}
				n2, err := strconv.ParseUint(part[i+2:], 10, 64)
				if err != nil {
					fmt.Println(err)
					break
				}
				if n2 < n1 {
					fmt.Printf("Invalid range %d-%d\n", n1, n2)
					break
				}
				for ii := n1; ii <= n2; ii++ {
					c <- ii
				}
			}
		}
		close(c)
	}(epochs)
	return c
}

// populate reward shares for a range of epoches
func calculateRewardShares(operator string, delegate []byte, epoches string) string {
	if epoches == "" {
		return calculateEpochRewardShares(
			operator, delegate, currentEpochNum()).String()
	}

	// parse a range of epoches
	var result string
	for e := range epochRangeGen(epoches) {
		result += calculateEpochRewardShares(operator, delegate, e).String() + "\n"
	}
	return result
}

// payout pays tokens out to delegates on IoTeX blockchain
func payout(delegate string, operator string) string {
	// get operator's address
	operator_addr, err := alias.Address(operator)
	if err != nil {
		panic(err)
	}

	// get delegate's name to 12-byte array
	delegate_name := delegateName(delegate)

	return calculateRewardShares(operator_addr, delegate_name, epochToQuery)
}
