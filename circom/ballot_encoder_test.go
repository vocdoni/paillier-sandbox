package circom

import (
	"encoding/json"
	"testing"
)

func TestBallotEncoder(t *testing.T) {
	// circuit files
	wasmFile := "./artifacts/ballot_encoder.wasm"
	zkeyFile := "./artifacts/ballot_encoder_pkey.zkey"
	// init inputs
	inputs := map[string]any{
		"fields": []string{"5", "1", "4", "3", "0", "0", "0"},
		"mask":   []string{"1", "1", "1", "1", "0", "0", "0"},
		"base":   "100",
	}
	// compile and generate proof
	bInputs, _ := json.Marshal(inputs)
	proofData, pubSignals, err := CompileAndGenerateProof(bInputs, wasmFile, zkeyFile)
	if err != nil {
		t.Errorf("Error compiling and generating proof: %v\n", err)
		return
	}
	t.Log("Proof data:", proofData)
	t.Log("Public signals:", pubSignals)
}
