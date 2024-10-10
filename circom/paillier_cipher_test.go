package circom

import (
	"encoding/json"
	"log"
	"math/big"
	"testing"
)

func TestPaillierCipher(t *testing.T) {
	var (
		// circuit parameters
		lSize  = 32
		nLimbs = 16
		// circuit assets
		wasmFile = "./artifacts/paillier_cipher_test.wasm"
		zkeyFile = "./artifacts/paillier_cipher_test_pkey.zkey"
	)
	// encrypt
	raw, _ := new(big.Int).SetString("102030405", 10)
	pk, rnd, c, err := EncryptWithPaillier(raw)
	if err != nil {
		log.Fatalf("Error encrypting: %v\n", err)
		return
	}
	// get the cached constant values
	cv := pk.Cache()
	rToNToS := new(big.Int).Exp(rnd, cv.NToS, cv.NToSPlusOne)
	// init inputs
	inputs := map[string]any{
		"m":               raw.String(),
		"n_plus_one":      BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, cv.NPlusOne)),
		"r_to_n_to_s":     BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, rToNToS)),
		"n_to_s_plus_one": BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, cv.NToSPlusOne)),
		"ciphertext":      BigIntArrayToStringArray(BigIntToArray(lSize, nLimbs, c)),
	}
	bInputs, _ := json.Marshal(inputs)
	log.Println("Inputs:", string(bInputs))
	proofData, pubSignals, err := CompileAndGenerateProof(bInputs, wasmFile, zkeyFile)
	if err != nil {
		t.Errorf("Error compiling and generating proof: %v\n", err)
		return
	}
	log.Println("Proof data:", proofData)
	log.Println("Public signals:", pubSignals)
}
