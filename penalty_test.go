package main

import (
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	miner18 "github.com/filecoin-project/go-state-types/builtin/v18/miner"
	smoothing18 "github.com/filecoin-project/go-state-types/builtin/v18/util/smoothing"
	miner19 "github.com/filecoin-project/go-state-types/builtin/v19/miner"
	"github.com/filecoin-project/go-state-types/builtin/v19/util/smoothing"
)

func TestFaultFeeSectorQuality(t *testing.T) {
	size := abi.SectorSize(32 << 30)
	reward := smoothing.NewEstimate(big.NewInt(1_000_000_000_000_000_000), big.Zero())
	power := smoothing.NewEstimate(big.NewInt(1<<60), big.Zero())
	sector := miner19.SectorOnChainInfo{
		PowerBaseEpoch:     100,
		Expiration:         100 + 180*2880,
		VerifiedDealWeight: big.Zero(),
		Flags:              miner19.SIMPLE_QA_POWER,
	}
	legacyFee := FaultFee(size, &sector, reward, power)
	oldSector := miner18.SectorOnChainInfo{
		PowerBaseEpoch:     sector.PowerBaseEpoch,
		Expiration:         sector.Expiration,
		VerifiedDealWeight: sector.VerifiedDealWeight,
		Flags:              miner18.SIMPLE_QA_POWER,
	}
	oldFee := miner18.ExpectedRewardForPower(
		smoothing18.FilterEstimate(reward), smoothing18.FilterEstimate(power),
		miner18.QAPowerForSector(size, &oldSector), 10108,
	)
	if !legacyFee.Equals(oldFee) {
		t.Fatalf("legacy fee changed: got %s, want %s", legacyFee, oldFee)
	}

	// nv29 full-QA sectors get 10x power even without verified deals.
	sector.Flags |= miner19.FULL_QA_POWER
	fullFee := FaultFee(size, &sector, reward, power)
	difference := big.Sub(fullFee, big.Mul(legacyFee, big.NewInt(10)))
	// The reward calculation rounds down to whole attoFIL.
	if difference.LessThan(big.Zero()) || !difference.LessThan(big.NewInt(10)) {
		t.Fatalf("full-QA fee should be 10x legacy fee within rounding: legacy=%s full=%s", legacyFee, fullFee)
	}
}
