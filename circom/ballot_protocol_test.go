package circom

import (
	"encoding/json"
	"testing"
)

func TestBallotProtocol(t *testing.T) {
	// circuit files
	wasmFile := "./artifacts/ballot_protocol.wasm"
	zkeyFile := "./artifacts/ballot_protocol_pkey.zkey"
	// init inputs
	inputs := map[string]any{
		"fields":           []string{"1", "2", "3", "0", "0"}, // total_cost = 1^2 + 2^2 + 3^2 = 14
		"max_count":        "3",                               // number of valid values in fields
		"force_uniqueness": "1",                               // no boolean type in circom
		"max_value":        "4",
		"min_value":        "0",
		"max_total_cost":   "15",
		"min_total_cost":   "13",
		"cost_exp":         "2",
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
