package circom

import (
	"math/big"
	"os"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/witness"
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
	return prover.Groth16ProverRaw(bZkey, w)
}
