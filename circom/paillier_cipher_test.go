package circom

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/witness"
	"github.com/niclabs/tcpaillier"
)

const (
	wordSize = 32
	nChunks  = 16
	bitSize  = 512
	s        = uint8(1)
	l        = uint8(5) // number of shares
	k        = uint8(3) // threshold

	wasmFile = "./artifacts/paillier_cipher_test_js/paillier_cipher_test.wasm"
	zkeyFile = "./artifacts/proving_key.zkey"
)

var (
	bInputs []byte
	bWasm   []byte
	bZkey   []byte
)

// bigIntToArray converts a big.Int into an array of k big.Int elements, it is
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
func bigIntToArray(n int, k int, x *big.Int) []*big.Int {
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

// bigIntArrayToStringArray converts an array of big.Int into an array of strings
func bigIntArrayToStringArray(arr []*big.Int) []string {
	ret := make([]string, len(arr))
	for i, v := range arr {
		ret[i] = v.String()
	}
	return ret
}

func TestMain(m *testing.M) {
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
	raw := big.NewInt(1234567890)
	c, err := pk.EncryptFixed(raw, randMod)
	if err != nil {
		log.Fatalf("Error encrypting: %v\n", err)
		return
	}
	// get the cached constant values
	cv := pk.Cache()
	g := bigIntToArray(wordSize, nChunks, cv.NPlusOne)
	nToS := bigIntToArray(wordSize, nChunks, cv.NToS)
	nToSPlusOne := bigIntToArray(wordSize, nChunks, cv.NToSPlusOne)
	msg := bigIntToArray(wordSize, nChunks, raw)
	r := bigIntToArray(wordSize, nChunks, randMod)
	ciphertext := bigIntToArray(wordSize, nChunks, c)
	// init inputs
	inputs := map[string][]string{
		"g":               bigIntArrayToStringArray(g),
		"n_to_s":          bigIntArrayToStringArray(nToS),
		"n_to_s_plus_one": bigIntArrayToStringArray(nToSPlusOne),
		"msg":             bigIntArrayToStringArray(msg),
		"r":               bigIntArrayToStringArray(r),
		"ciphertext":      bigIntArrayToStringArray(ciphertext),
	}
	// write inputs in json file
	if bInputs, err = json.Marshal(inputs); err != nil {
		log.Fatalf("Error marshalling inputs: %v\n", err)
		return
	}
	log.Println("Inputs:", string(bInputs))
	// read wasm file
	if bWasm, err = os.ReadFile(wasmFile); err != nil {
		log.Fatalf("Error reading wasm file: %v\n", err)
		return
	}
	// read zkey file
	if bZkey, err = os.ReadFile(zkeyFile); err != nil {
		log.Fatalf("Error reading zkey file: %v\n", err)
		return
	}
	os.Exit(m.Run())
}

func TestPaillierCipher(t *testing.T) {
	// parse inputs
	inputs, err := witness.ParseInputs(bInputs)
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
	w, err := calc.CalculateWTNSBin(inputs, true)
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
