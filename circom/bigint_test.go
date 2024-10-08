package circom

import (
	"encoding/json"
	"log"
	"math/big"
	"testing"
)

func TestBigModExp(t *testing.T) {
	// circuit files
	wasmFile := "./artifacts/bigint_test.wasm"
	zkeyFile := "./artifacts/bigint_test_pkey.zkey"
	// init the inputs and calculate the result
	exponent := big.NewInt(3)
	base := big.NewInt(2)
	modulus := big.NewInt(3)
	result := new(big.Int).Exp(base, exponent, modulus)
	// debug print
	log.Printf("%v^%v mod %v = %v\n", base, exponent, modulus, result)
	// create circuit inputs
	inputs := map[string]any{
		"exponent": exponent.String(),
		"base":     bigIntArrayToStringArray(bigIntToArray(8, 4, base)),
		"modulus":  bigIntArrayToStringArray(bigIntToArray(8, 4, modulus)),
		"result":   bigIntArrayToStringArray(bigIntToArray(8, 4, result)),
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
