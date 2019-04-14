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
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/iotex-address/address"
	"math/big"
)

// reward type:
//   - block reward: reward for each block mined
//   - epoch reward: epoch bonus reward
//   - extra reward: extra reward during first year
//
// see https://medium.com/iotex/iotex-delegates-program-application-voting-and-rewards-5cab7e87bd20
type Reward struct {
	Block           uint64 `json:"block"`
	FoundationBonus uint64 `json:"foundation"`
	EpochBonus      uint64 `json:"epoch"`
}

type Share struct {
	IOAddr  string `json:"ioaddr"`
	ETHAddr string `json:"ethaddr"`
	Votes   []string `json:"votes"`
	Share   []uint64 `json:share`
	VotedPeriod []uint64 `json:"voteperiod"`
	Reward  Reward `json:"reward"`
}

type RewardShares struct {
	EpochNum     string   `json:"epochnum"`
	Productivity uint64   `json:"productivity"`
	TotalVotes   []string `json:"votes"`
	Reward       Reward   `json:"reward"`
	Shares       []Share  `json:"shares"`
}

func addReward(self Reward, other Reward) Reward {
	return Reward{
		self.Block + other.Block,
		self.FoundationBonus + other.FoundationBonus,
		self.EpochBonus + other.EpochBonus}
}

// Set total reward
func (rs *RewardShares) SetReward(total Reward) *RewardShares {
	rs.Reward = total
	return rs
}

// Set total votes
func (rs *RewardShares) SetTotalVotes(votes *big.Int) *RewardShares {
	rs.TotalVotes = []string{votes.String()}
	return rs
}

// Set epoch number
func (rs *RewardShares) SetEpochNum(epoch string) *RewardShares {
	rs.EpochNum = epoch
	return rs
}

// Set number of produced blocks
func (rs *RewardShares) SetProductivity(prod uint64) *RewardShares {
	rs.Productivity = prod
	return rs
}

// Combine two epochs' rewardshares
func (rs *RewardShares) Combine(other *RewardShares) *RewardShares {
	rs.Productivity += other.Productivity
	rs.TotalVotes = append(rs.TotalVotes, other.TotalVotes...)

	rs.Reward = addReward(rs.Reward, other.Reward)

	for _, right := range other.Shares {
		existing := false
		for i, left := range rs.Shares {
			// update the voters that exist in previous epochs.
			if right.IOAddr == left.IOAddr {
				existing = true
				rs.Shares[i].Reward = addReward(left.Reward, right.Reward)
				if !simpleJson {
					rs.Shares[i].Votes = append(left.Votes, right.Votes...)
					rs.Shares[i].Share = append(left.Share, right.Share...)
					rs.Shares[i].VotedPeriod = append(left.VotedPeriod, right.VotedPeriod...)
				}
			}
		}
		if !existing {
			rs.Shares = append(rs.Shares, right)
		}
	}
	return rs
}

// Debug string
func (rs *RewardShares) String() string {
	rs_str, _ := json.Marshal(rs)
	return string(rs_str)
}

// Based on the obtained votes, calculate voter's shares
func (rs *RewardShares) CalculateShares(bps map[string]*big.Int, total *big.Int, epoch uint64) *RewardShares {
	// calculate each voter's meta info
	rs.Shares = nil
	for addr, vote := range bps {
		var share Share

		hex_addr, _ := address.FromBytes(common.HexToAddress(addr).Bytes())
		share.IOAddr = hex_addr.String()
		share.ETHAddr = addr

		base := big.NewInt(1000)
		base = base.Mul(base, vote)
		base = base.Div(base, total)
		percentMille := base.Uint64()

		if !simpleJson {
			share.Votes = []string{vote.Text(10)}
			share.Share = []uint64{percentMille}
			share.VotedPeriod = []uint64{epoch}
		}

		share.Reward = Reward{
			percentMille * rs.Reward.Block * (100 - blockComm) / (1000 * 100),
			percentMille * rs.Reward.FoundationBonus * (100 - foundationComm) / (1000 * 100),
			percentMille * rs.Reward.EpochBonus * (100 - epochComm) / (1000 * 100)}

		rs.Shares = append(rs.Shares, share)
	}

	return rs
}

// Allocate new RewardShares
func NewRewardShares() *RewardShares {
	return new(RewardShares)
}
