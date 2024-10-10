package circom

import (
	"encoding/json"
	"testing"
)

func TestBallotEncoder(t *testing.T) {
	// circuit files
	wasmFile := "./artifacts/ballot_encoder_test.wasm"
	zkeyFile := "./artifacts/ballot_encoder_test_pkey.zkey"
	// init inputs
	inputs := map[string]any{
		"fields": IntArrayToStringArray([]int{1, 2, 3, 4, 5}, 7),
		"mask":   IntArrayToStringArray([]int{1, 1, 1, 1, 1}, 7),
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

	var decPubSignals []string
	if err := json.Unmarshal([]byte(pubSignals), &decPubSignals); err != nil {
		t.Errorf("Error unmarshalling public signals: %v\n", err)
		return
	}
	expected := EncodeBallot([]int{1, 2, 3, 4, 5}, BallotConfig{
		MaxCount: 5,
		Base:     100,
	})
	if decPubSignals[0] != expected.String() {
		t.Errorf("Incorrect public signal: expected %s, got %s\n", expected.String(), decPubSignals[0])
		return
	}
}
