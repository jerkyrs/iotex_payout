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
	Votes   string `json:"votes"`
	Share   uint64 `json:"share"`
	Reward  Reward `json:"reward"`
}

type RewardShares struct {
	EpochNum     uint64  `json:"epochnum"`
	Productivity uint64  `json:"productivity"`
	TotalVotes   string  `json:"votes"`
	Reward       Reward  `json:"reward"`
	Shares       []Share `json:"shares"`
}

// Set total reward
func (rs *RewardShares) SetReward(total Reward) *RewardShares {
	rs.Reward = total
	return rs
}

// Set total votes
func (rs *RewardShares) SetTotalVotes(votes *big.Int) *RewardShares {
	rs.TotalVotes = votes.String()
	return rs
}

// Set epoch number
func (rs *RewardShares) SetEpochNum(epoch uint64) *RewardShares {
	rs.EpochNum = epoch
	return rs
}

// Set number of produced blocks
func (rs *RewardShares) SetProductivity(prod uint64) *RewardShares {
	rs.Productivity = prod
	return rs
}

// Debug string
func (rs *RewardShares) String() string {
	rs_str, _ := json.Marshal(rs)
	return string(rs_str)
}

// Based on the obtained votes, calculate voter's shares
func (rs *RewardShares) CalculateShares(bps map[string]*big.Int, total *big.Int) *RewardShares {
	// calculate each voter's meta info
	rs.Shares = nil
	for addr, vote := range bps {
		var share Share

		hex_addr, _ := address.FromBytes(common.HexToAddress(addr).Bytes())
		share.IOAddr = hex_addr.String()
		share.ETHAddr = addr

		share.Votes = vote.Text(10)

		base := big.NewInt(1000)
		base = base.Mul(base, vote)
		base = base.Div(base, total)
		share.Share = base.Uint64()

		share.Reward = Reward{
			share.Share * rs.Reward.Block * (100 - blockComm) / (1000 * 100),
			share.Share * rs.Reward.FoundationBonus * (100 - foundationComm) / (1000 * 100),
			share.Share * rs.Reward.EpochBonus * (100 - epochComm) / (1000 * 100)}

		rs.Shares = append(rs.Shares, share)
	}

	return rs
}

// Allocate new RewardShares
func NewRewardShares() *RewardShares {
	return new(RewardShares)
}
