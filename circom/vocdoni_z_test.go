package circom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"testing"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"go.vocdoni.io/dvote/util"
)

func TestVocdoniZ(t *testing.T) {
	var (
		fields       = []int{3, 5, 2, 4, 1}
		n_fields     = 5
		maxCount     = 5
		maxValue     = 16 + 1
		minValue     = 0
		base         = 16000001
		costExp      = 2
		address, _   = hex.DecodeString("0x6Db989fbe7b1308cc59A27f021e2E3de9422CF0A")
		processID, _ = hex.DecodeString("0xf16236a51F11c0Bf97180eB16694e3A345E42506")
		secret, _    = hex.DecodeString("super-secret-mnemonic-phrase")
		// circuit parameters
		lSize  = 32
		nLimbs = 8
		// circuit assets
		wasmFile = "./artifacts/vocdoni_z.wasm"
		zkeyFile = "./artifacts/vocdoni_z_pkey.zkey"
		vkeyFile = "./artifacts/vocdoni_z_vkey.json"
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
	// generate the nullifier
	commitment, err := poseidon.Hash([]*big.Int{
		util.BigToFF(new(big.Int).SetBytes(address)),
		util.BigToFF(new(big.Int).SetBytes(processID)),
		util.BigToFF(new(big.Int).SetBytes(secret)),
	})
	if err != nil {
		log.Fatalf("Error hashing: %v\n", err)
		return
	}
	nullifier, err := poseidon.Hash([]*big.Int{
		commitment,
		util.BigToFF(new(big.Int).SetBytes(secret)),
	})
	log.Println("Commitment:", commitment)
	log.Println("Nullifier:", nullifier)
	// circuit inputs
	inputs := map[string]any{
		"fields":           IntArrayToStringArray(fields, n_fields),
		"max_count":        fmt.Sprint(maxCount),
		"force_uniqueness": "0",
		"max_value":        fmt.Sprint(maxValue),
		"min_value":        fmt.Sprint(minValue),
		"cost_exp":         fmt.Sprint(costExp),
		"max_total_cost":   fmt.Sprint(int(math.Pow(float64(maxValue-1), float64(costExp))) * maxCount), // (maxValue-1)^costExp * maxCount
		"min_total_cost":   fmt.Sprint(maxCount),
		"cost_from_weight": "0",
		"weight":           "1",
		"base":             fmt.Sprint(base),
		"n_plus_one":       BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, cv.NPlusOne)),
		"r_to_n_to_s":      BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, rToNToS)),
		"n_to_s_plus_one":  BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, cv.NToSPlusOne)),
		"ciphertext":       BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, c)),
		"nullifier":        nullifier.String(),
		"commitment":       commitment.String(),
		"secret":           util.BigToFF(new(big.Int).SetBytes(secret)).String(),
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
	// read zkey file
	vkey, err := os.ReadFile(vkeyFile)
	if err != nil {
		t.Errorf("Error reading zkey file: %v\n", err)
		return
	}
	if err := VerifyProof(proofData, pubSignals, vkey); err != nil {
		t.Errorf("Error verifying proof: %v\n", err)
		return
	}
	log.Println("Proof verified")
}
