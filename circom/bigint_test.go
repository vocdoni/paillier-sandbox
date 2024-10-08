package circom

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/witness"
)

func TestBigModExp(t *testing.T) {
	// read wasm file
	wasmFile := "./artifacts/bigint_test.wasm"
	bWasm, err := os.ReadFile(wasmFile)
	if err != nil {
		log.Fatalf("Error reading wasm file: %v\n", err)
		return
	}
	// read zkey file
	zkeyFile := "./artifacts/bigint_test_pkey.zkey"
	bZkey, err := os.ReadFile(zkeyFile)
	if err != nil {
		log.Fatalf("Error reading zkey file: %v\n", err)
		return
	}
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
	finalInputs, err := witness.ParseInputs(bInputs)
	if err != nil {
		t.Fatalf("Error parsing inputs: %v", err)
		return
	}
	// instance witness calculator
	calc, err := witness.NewCircom2WitnessCalculator(bWasm, true)
	if err != nil {
		t.Fatalf("Error creating witness calculator: %v", err)
		return
	}
	// calculate witness
	w, err := calc.CalculateWTNSBin(finalInputs, true)
	if err != nil {
		t.Fatalf("Error calculating witness: %v", err)
		return
	}
	proofData, pubSignals, err := prover.Groth16ProverRaw(bZkey, w)
	if err != nil {
		t.Fatalf("Error generating proof: %v", err)
		return
	}
	log.Println("Proof data:", proofData)
	log.Println("Public signals:", pubSignals)
}
