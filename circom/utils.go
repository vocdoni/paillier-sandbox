package circom

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
	"github.com/iden3/go-rapidsnark/witness"
	"github.com/niclabs/tcpaillier"
)

type ProofData struct {
	A []string   `json:"pi_a"`
	B [][]string `json:"pi_b"`
	C []string   `json:"pi_c"`
}

const (
	// paillier parameters
	bitSize = 64
	s       = uint8(1)
	l       = uint8(5) // number of shares
	k       = uint8(3) // threshold
)

// BigIntToArray converts a big.Int into an array of k big.Int elements, it is
// the go implementation of the javascript function:
//
//	function bigint_to_array(n: number, k: number, x: bigint) {
//	    let mod: bigint = 1n;
//	    for (var idx = 0; idx < n; idx++) {
//	        mod = mod * 2n;
//	    }
//	    let ret: bigint[] = [];
//	    var x_temp: bigint = x;
//	    for (var idx = 0; idx < k; idx++) {
//	        ret.push(x_temp % mod);
//	        x_temp = x_temp / mod;
//	    }
//	    return ret;
//	}
func BigIntToArray(n int, k int, x *big.Int) []*big.Int {
	// precompute 2^n as mod, left shift is equivalent to multiplying by 2^n
	mod := new(big.Int).Lsh(big.NewInt(1), uint(n))
	// initialize the result array
	ret := make([]*big.Int, k)
	// temporary variable to avoid re-allocating big.Int in each iteration
	xTemp := new(big.Int).Set(x)
	// loop to compute each element in the result array
	for i := 0; i < k; i++ {
		// ret[i] = xTemp % mod
		ret[i] = new(big.Int).Mod(xTemp, mod)
		// xTemp = xTemp / mod
		xTemp.Div(xTemp, mod)
	}
	return ret
}

// BigIntArrayToStringArray converts an array of big.Int into an array of strings
func BigIntArrayToStringArray(arr []*big.Int) []string {
	ret := make([]string, len(arr))
	for i, v := range arr {
		ret[i] = v.String()
	}
	return ret
}

func IntArrayToStringArray(arr []int, n int) []string {
	strArr := make([]string, n)
	for i := 0; i < n; i++ {
		if i < len(arr) {
			strArr[i] = fmt.Sprint(arr[i])
		} else {
			strArr[i] = "0"
		}
	}
	return strArr
}

func CompileAndGenerateProof(inputs []byte, wasmFile, zkeyFile string) (string, string, error) {
	finalInputs, err := witness.ParseInputs(inputs)
	if err != nil {
		return "", "", err
	}
	// read wasm file
	bWasm, err := os.ReadFile(wasmFile)
	if err != nil {
		return "", "", err
	}
	// read zkey file
	bZkey, err := os.ReadFile(zkeyFile)
	if err != nil {
		return "", "", err
	}
	// instance witness calculator
	calc, err := witness.NewCircom2WitnessCalculator(bWasm, true)
	if err != nil {
		return "", "", err
	}
	// calculate witness
	w, err := calc.CalculateWTNSBin(finalInputs, true)
	if err != nil {
		return "", "", err
	}
	// generate proof
	return prover.Groth16ProverRaw(bZkey, w)
}

func VerifyProof(proofData, pubSignals string, vkey []byte) error {
	data := ProofData{}
	if err := json.Unmarshal([]byte(proofData), &data); err != nil {
		return err
	}
	signals := []string{}
	if err := json.Unmarshal([]byte(pubSignals), &signals); err != nil {
		return err
	}
	proof := types.ZKProof{
		Proof: &types.ProofData{
			A: data.A,
			B: data.B,
			C: data.C,
		},
		PubSignals: signals,
	}
	return verifier.VerifyGroth16(proof, vkey)
}

func EncryptWithPaillier(raw *big.Int) (*tcpaillier.PubKey, *big.Int, *big.Int, error) {
	// generate the public key
	_, pk, err := tcpaillier.NewKey(bitSize, s, l, k)
	if err != nil {
		return nil, nil, nil, err
	}
	// get a random mod
	rnd, err := pk.RandomModNToSPlusOneStar()
	if err != nil {
		return nil, nil, nil, err
	}
	// encrypt with rnd
	c, err := pk.EncryptFixed(raw, rnd)
	if err != nil {
		return nil, nil, nil, err
	}
	return pk, rnd, c, nil
}

// BallotConfig holds the configuration for the ballot protocol
type BallotConfig struct {
	MaxCount int
	Base     int
}

// powBigInt computes base^exp using big.Int
func powBigInt(base, exp int) *big.Int {
	result := big.NewInt(1)
	bBase := big.NewInt(int64(base))

	for i := 0; i < exp; i++ {
		result.Mul(result, bBase)
	}

	return result
}

// EncodeBallot encodes the ballot into a single big.Int number
func EncodeBallot(ballot []int, config BallotConfig) *big.Int {
	encoded := big.NewInt(0)
	for i, value := range ballot {
		positionValue := new(big.Int).Mul(big.NewInt(int64(value)), powBigInt(config.Base, config.MaxCount-i-1))
		encoded.Add(encoded, positionValue)
	}
	return encoded
}
