package circom

import "math/big"

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
