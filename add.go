package main

import "math/big"

// homomorphicAdd adds two encrypted values mod n^(s+1).
// This assumes the values have been encrypted using a Paillier-like scheme.
func homomorphicAdd(c1, c2, n *big.Int, s uint8) *big.Int {
	// Calculate n^(s+1)
	sPlusOne := new(big.Int).Add(big.NewInt(int64(s)), big.NewInt(1))
	nToSPlusOne := new(big.Int).Exp(n, sPlusOne, nil)

	// Perform the multiplication (homomorphic addition in Paillier)
	sum := new(big.Int).Mul(c1, c2)

	// Apply mod n^(s+1)
	sum.Mod(sum, nToSPlusOne)

	return sum
}
