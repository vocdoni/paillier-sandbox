package circom

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"testing"
)

func TestVocdoniZ(t *testing.T) {
	var (
		fields   = []int{6, 8, 13, 0, 2, 15, 10, 8, 0, 16}
		n_fields = 10
		maxCount = 10
		maxValue = 16 + 1
		minValue = 0
		base     = 160000001
		// circuit parameters
		lSize  = 32
		nLimbs = 16
		// circuit assets
		wasmFile = "./artifacts/vocdoni_z.wasm"
		zkeyFile = "./artifacts/vocdoni_z_pkey.zkey"
	)
	encodedBallot := EncodeBallot(fields, BallotConfig{
		MaxCount: maxCount,
		Base:     base,
	})
	// encrypt with r
	pk, rnd, c, err := EncryptWithPaillier(encodedBallot)
	if err != nil {
		log.Fatalf("Error encrypting: %v\n", err)
		return
	}
	// get the cached constant values
	cv := pk.Cache()
	rToNToS := new(big.Int).Exp(rnd, cv.NToS, cv.NToSPlusOne)
	// circuit inputs
	inputs := map[string]any{
		"fields":           IntArrayToStringArray(fields, n_fields),
		"max_count":        fmt.Sprint(maxCount),
		"force_uniqueness": "0",
		"max_value":        fmt.Sprint(maxValue),
		"min_value":        fmt.Sprint(minValue),
		"cost_exp":         "2",
		"max_total_cost":   "919",
		"min_total_cost":   "8",
		"cost_from_weight": "0",
		"weight":           "1",
		"base":             fmt.Sprint(base),
		"n_plus_one":       BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, cv.NPlusOne)),
		"r_to_n_to_s":      BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, rToNToS)),
		"n_to_s_plus_one":  BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, cv.NToSPlusOne)),
		"ciphertext":       BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, c)),
	}
	bInputs, _ := json.Marshal(inputs)
	t.Log("Inputs:", string(bInputs))
	proofData, pubSignals, err := CompileAndGenerateProof(bInputs, wasmFile, zkeyFile)
	if err != nil {
		t.Errorf("Error compiling and generating proof: %v\n", err)
		return
	}
	log.Println("Proof data:", proofData)
	log.Println("Public signals:", pubSignals)
}
