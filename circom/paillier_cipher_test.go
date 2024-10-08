package circom

import (
	"encoding/json"
	"log"
	"math/big"
	"testing"

	"github.com/niclabs/tcpaillier"
)

func TestPaillierCipher(t *testing.T) {
	var (
		bitSize = 64
		s       = uint8(1)
		l       = uint8(5) // number of shares
		k       = uint8(3) // threshold

		wasmFile = "./artifacts/paillier_cipher_test.wasm"
		zkeyFile = "./artifacts/paillier_cipher_test_pkey.zkey"
	)
	// generate the public key
	_, pk, err := tcpaillier.NewKey(bitSize, s, l, k)
	if err != nil {
		log.Fatalf("Error generating key: %v\n", err)
		return
	}
	// get a random mod
	randMod, err := pk.RandomModNToSPlusOneStar()
	if err != nil {
		log.Fatalf("Error generating random mod: %v\n", err)
		return
	}
	// encrypt with r
	raw := big.NewInt(255)
	c, err := pk.EncryptFixed(raw, randMod)
	if err != nil {
		log.Fatalf("Error encrypting: %v\n", err)
		return
	}
	// get the cached constant values
	cv := pk.Cache()
	rToNToS := new(big.Int).Exp(randMod, cv.NToS, cv.NToSPlusOne)
	// init inputs
	inputs := map[string]any{
		"m":               raw.String(),
		"n_plus_one":      bigIntArrayToStringArray(bigIntToArray(32, 16, cv.NPlusOne)),
		"r_to_n_to_s":     bigIntArrayToStringArray(bigIntToArray(32, 16, rToNToS)),
		"n_to_s_plus_one": bigIntArrayToStringArray(bigIntToArray(32, 16, cv.NToSPlusOne)),
		"ciphertext":      bigIntArrayToStringArray(bigIntToArray(32, 16, c)),
	}
	bInputs, _ := json.Marshal(inputs)
	log.Println("Inputs:", string(bInputs))
	proofData, pubSignals, err := compileAndGenerateProof(bInputs, wasmFile, zkeyFile)
	if err != nil {
		t.Errorf("Error compiling and generating proof: %v\n", err)
		return
	}
	log.Println("Proof data:", proofData)
	log.Println("Public signals:", pubSignals)
}
