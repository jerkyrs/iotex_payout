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
	"testing"
	"math/big"
)

func TestCalculateReward(t *testing.T) {
	var r Reward

	votes := big.NewInt(500)
	total := big.NewInt(1000)

	r = calculateReward(/*blks=*/8, /*elected=*/true,
			/*votes=*/votes, /*total=*/total)

	if r.Block != 8 * 16 {
		t.Error("Expect 16 IOTX per block")
	}

	if r.FoundationBonus != 80 {
		t.Error("Expect 80 IOTX as foundation bonus if elected")
	}

	if r.EpochBonus != uint64(6250) {
		t.Error("Expect 12500 IOTX as epoch bonus")
	}

	r = calculateReward(/*blks=*/8, /*elected=*/false,
			/*votes=*/votes, /*total=*/total)

	if r.FoundationBonus != 0 {
		t.Error("Expect 0 IOTX as foundation bonus if not elected")
	}
}

// Test reward shares are caculated correctly
func TestCalculateRewardShares(t *testing.T) {

	blockCommOrig := blockComm
	foundationCommOrig := foundationComm
	epochCommOrig := epochComm
	simpleJsonOrig := simpleJson

	// clear commission rate
	blockComm = 0
	foundationComm = 0
	epochComm = 0
	simpleJson = false

	bps := make(map[string]*big.Int)
	vote1 := big.NewInt(4001)
	vote2 := big.NewInt(5999)
	bps["45831656370acf0b345cc25558dc9b3b1424ddc3"] = vote1
	bps["7d5943b6ee8093be7dc75cd095ee12bfd6c8bf6c"] = vote2

	total := big.NewInt(10000)

	reward := Reward{
		/*Block=*/100,
		/*FoundationBonus=*/1000,
		/*EpochBonus=*/10000,
	}

	rs := NewRewardShares().SetReward(reward)
	rs = rs.CalculateShares(/*bps=*/bps, /*total=*/total, /*epoch=*/10)

	if rs.Shares == nil {
		t.Error("Failed to calculate reward shares")
	}
	for _, share := range rs.Shares {
		if share.ETHAddr == "45831656370acf0b345cc25558dc9b3b1424ddc3" {
			if share.Votes[0] != "4001" {
				t.Fatalf("Expect 40 votes, but actual obtained %v", share.Votes[0])
			}
			if share.Share[0] != 400 {
				t.Fatalf("Expect 400 shares, but actual obtained %v", share.Share[0])
			}
			if share.VotedPeriod[0] != 10 {
				t.Fatalf("Expect epoch 10, but actual obtained %v", share.VotedPeriod[0])
			}
			if share.Reward.Block != 40 ||
			   share.Reward.FoundationBonus != 400 ||
			   share.Reward.EpochBonus != 4000 {
				t.Fatalf("Expect reward {40, 400, 4000}, but actual obtained {%v, %v, %v}",
					 share.Reward.Block, share.Reward.FoundationBonus,
					 share.Reward.EpochBonus)
			}
		} else if share.ETHAddr == "7d5943b6ee8093be7dc75cd095ee12bfd6c8bf6c" {
			if share.Votes[0] != "5999" {
				t.Fatalf("Expect 60 votes, but actual obtained %v", share.Votes[0])
			}
			if share.Share[0] != 599 {
				t.Fatalf("Expect 600 shares, but actual obtained %v", share.Share[0])
			}
			if share.VotedPeriod[0] != 10 {
				t.Fatalf("Expect epoch 10, but actual obtained %v", share.VotedPeriod[0])
			}
			if share.Reward.Block != 59 ||
			   share.Reward.FoundationBonus != 599 ||
			   share.Reward.EpochBonus != 5990 {
				t.Fatalf("Expect reward {59, 599, 5990}, but actual obtained {%v, %v, %v}",
					 share.Reward.Block, share.Reward.FoundationBonus,
					 share.Reward.EpochBonus)
			}
		} else {
			t.Fatalf("Unexpected ETH address %v", share.ETHAddr)
		}
	}

	// test commission fee
	blockComm = 10
	foundationComm = 10
	epochComm = 10

	rs = rs.CalculateShares(/*bps=*/bps, /*total=*/total, /*epoch=*/10)
	for _, share := range rs.Shares {
		if share.ETHAddr == "45831656370acf0b345cc25558dc9b3b1424ddc3" {
			if share.Reward.Block != 36 ||
			   share.Reward.FoundationBonus != 360 ||
			   share.Reward.EpochBonus != 3600 {
				t.Fatalf("Expect reward {36, 360, 3600}, but actual obtained {%v, %v, %v}",
					 share.Reward.Block, share.Reward.FoundationBonus,
					 share.Reward.EpochBonus)
			}
		} else if share.ETHAddr == "7d5943b6ee8093be7dc75cd095ee12bfd6c8bf6c" {
			if share.Reward.Block != 53 ||
			   share.Reward.FoundationBonus != 539 ||
			   share.Reward.EpochBonus != 5391 {
				t.Fatalf("Expect reward {54, 540, 5391}, but actual obtained {%v, %v, %v}",
					 share.Reward.Block, share.Reward.FoundationBonus,
					 share.Reward.EpochBonus)
			}
		} else {
			t.Fatalf("Unexpected ETH address %v", share.ETHAddr)
		}
	}

	// test simpleJson
	simpleJson = true

	rs = rs.CalculateShares(/*bps=*/bps, /*total=*/total, /*epoch=*/10)
	for _, share := range rs.Shares {
		if share.Votes != nil {
			t.Fatalf("Expect (nil) votes, but actual obtained %v", share.Votes)
		}
		if share.Share != nil {
			t.Fatalf("Expect (nil) shares, but actual obtained %v", share.Share)
		}
		if share.VotedPeriod != nil {
			t.Fatalf("Expect epoch (nil), but actual obtained %v", share.VotedPeriod)
		}
	}

	// restore config
	blockComm = blockCommOrig
	foundationComm = foundationCommOrig
	epochComm = epochCommOrig
	simpleJson = simpleJsonOrig
}

func TestRewardShareString(t *testing.T) {
	rs := RewardShares{
		/*EpochNum=*/"20",
		/*Productivity=*/10,
		/*TotalVotes=*/[]string{"10"},
		/*Reward=*/Reward{/*Block=*/10,
			/*FoundationBonus=*/10,
			/*EpochBonus=*/10,
		},
		/*Shares=*/[]Share{Share{/*IOAddr=*/"io1",
			/*ETHAddr=*/"xxx",
			/*Votes=*/[]string{"5"},
			/*Share=*/[]uint64{500},
			/*VotedPeriod=*/[]uint64{0},
			/*Reward=*/Reward{5, 5, 5},
		}, Share{/*IOAddr=*/"io2",
			/*ETHAddr=*/"yyy",
			/*Votes=*/[]string{"5"},
			/*Share=*/[]uint64{500},
			/*VotedPeriod=*/[]uint64{0},
			/*Reward=*/Reward{5, 5, 5},
		}},
	}

	const expected = `{"epochnum":"20","productivity":10,"votes":["10"],` +
		`"reward":{"block":10,"foundation":10,"epoch":10},` +
		`"shares":[{"ioaddr":"io1","ethaddr":"xxx","votes":["5"],"Share":[500],` +
		`"voteperiod":[0],"reward":{"block":5,"foundation":5,"epoch":5}},` +
		`{"ioaddr":"io2","ethaddr":"yyy","votes":["5"],"Share":[500],` +
		`"voteperiod":[0],"reward":{"block":5,"foundation":5,"epoch":5}}]}`

	if rs.String() != expected {
		t.Fatalf("Expected :\n" +
			"%v\n" +
			"But actually:\n" +
			"%v", expected, rs.String())
	}
}

func TestCombineReward(t *testing.T) {
	rs1 := RewardShares{
		/*EpochNum=*/"",
		/*Productivity=*/10,
		/*TotalVotes=*/[]string{"10"},
		/*Reward=*/Reward{/*Block=*/10,
			/*FoundationBonus=*/10,
			/*EpochBonus=*/10,
		},
		/*Shares=*/[]Share{Share{/*IOAddr=*/"io1",
			/*ETHAddr=*/"xxx",
			/*Votes=*/[]string{"5"},
			/*Share=*/[]uint64{500},
			/*VotedPeriod=*/[]uint64{0},
			/*Reward=*/Reward{5, 5, 5},
		}, Share{/*IOAddr=*/"io2",
			/*ETHAddr=*/"yyy",
			/*Votes=*/[]string{"5"},
			/*Share=*/[]uint64{500},
			/*VotedPeriod=*/[]uint64{0},
			/*Reward=*/Reward{5, 5, 5},
		}},
	}
	rs2 := RewardShares{
		/*EpochNum=*/"",
		/*Productivity=*/20,
		/*TotalVotes=*/[]string{"20"},
		/*Reward=*/Reward{/*Block=*/20,
			/*FoundationBonus=*/20,
			/*EpochBonus=*/20,
		},
		/*Shares=*/[]Share{Share{/*IOAddr=*/"io1",
			/*ETHAddr=*/"xxx",
			/*Votes=*/[]string{"10"},
			/*Share=*/[]uint64{500},
			/*VotedPeriod=*/[]uint64{1},
			/*Reward=*/Reward{10, 10, 10},
		}, Share{/*IOAddr=*/"io3",
			/*ETHAddr=*/"zzz",
			/*Votes=*/[]string{"10"},
			/*Share=*/[]uint64{500},
			/*VotedPeriod=*/[]uint64{1},
			/*Reward=*/Reward{10, 10, 10},
		}},
	}
	expected := RewardShares{
		/*EpochNum=*/"",
		/*Productivity=*/30,
		/*TotalVotes=*/[]string{"10", "20"},
		/*Reward=*/Reward{/*Block=*/30,
			/*FoundationBonus=*/30,
			/*EpochBonus=*/30,
		},
		/*Shares=*/[]Share{Share{/*IOAddr=*/"io1",
			/*ETHAddr=*/"xxx",
			/*Votes=*/[]string{"5", "10"},
			/*Share=*/[]uint64{500, 500},
			/*VotedPeriod=*/[]uint64{0, 1},
			/*Reward=*/Reward{15, 15, 15},
		}, Share{/*IOAddr=*/"io2",
			/*ETHAddr=*/"yyy",
			/*Votes=*/[]string{"5"},
			/*Share=*/[]uint64{500},
			/*VotedPeriod=*/[]uint64{0},
			/*Reward=*/Reward{5, 5, 5},
		}, Share{/*IOAddr=*/"io3",
			/*ETHAddr=*/"zzz",
			/*Votes=*/[]string{"10"},
			/*Share=*/[]uint64{500},
			/*VotedPeriod=*/[]uint64{1},
			/*Reward=*/Reward{10, 10, 10},
		}},
	}

	rs1.Combine(&rs2)
	if rs1.String() != expected.String() {
		t.Fatalf("Expected :\n" +
			"%v\n" +
			"But actually:\n" +
			"%v", expected.String(), rs1.String())
	}
}
