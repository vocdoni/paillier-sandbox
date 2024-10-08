package circom

import (
	"encoding/json"
	"testing"
)

func TestUniqueNumbers(t *testing.T) {
	// circuit files
	wasmFile := "./artifacts/unique_test.wasm"
	zkeyFile := "./artifacts/unique_test_pkey.zkey"
	// init inpurs
	inputs := map[string]any{
		"arr": []string{"3", "2", "1", "5", "4"},
		"min": "0",
		"max": "6",
	}
	// compile and generate proof
	bInputs, _ := json.Marshal(inputs)
	proofData, pubSignals, err := compileAndGenerateProof(bInputs, wasmFile, zkeyFile)
	if err != nil {
		t.Errorf("Error compiling and generating proof: %v\n", err)
		return
	}
	t.Log("Proof data:", proofData)
	t.Log("Public signals:", pubSignals)
}
