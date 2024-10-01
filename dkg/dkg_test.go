// == dkg_test.go ==
package dkg

import (
	"fmt"
	"math/big"
	"testing"
)

func TestDKG(t *testing.T) {
	// Parameters
	bits := 512                     // Bit length for primes
	q, p := GenerateSafePrime(bits) // Generate safe prime p and corresponding q
	k := 3                          // Threshold
	nParties := 5                   // Number of parties

	// Verify that p is prime (should always be true here)
	if !p.ProbablyPrime(20) {
		t.Fatal("p is not prime")
	}

	// Find generator g of subgroup of order q
	g := FindGenerator(p, q)

	fmt.Printf("Prime p: %s\n", p.String())
	fmt.Printf("Prime q: %s\n", q.String())
	fmt.Printf("Generator g: %s\n", g.String())

	// Each party generates its polynomial and commitments
	polynomials := make([][]*big.Int, nParties)
	commitments := make([][]*big.Int, nParties)

	for i := 0; i < nParties; i++ {
		polynomials[i] = GeneratePolynomial(k, q)
		commitments[i] = GenerateCommitments(polynomials[i], g, p)
	}

	// Each party generates shares for all other parties
	shares := make([][]*big.Int, nParties) // shares[i][j]: share from party i to party j
	for i := 0; i < nParties; i++ {
		shares[i] = make([]*big.Int, nParties)
		for j := 0; j < nParties; j++ {
			shares[i][j] = GenerateShare(j+1, polynomials[i], q)
		}
	}

	// Each party verifies the shares received from others
	for i := 0; i < nParties; i++ {
		for j := 0; j < nParties; j++ {
			// Party i verifies share from party j
			valid := VerifyShare(shares[j][i], i+1, commitments[j], g, p)
			if !valid {
				t.Fatalf("Share verification failed for party %d's share from party %d", i+1, j+1)
			}
		}
	}

	// Each party computes its aggregate share
	aggregatedShares := make([]*big.Int, nParties)
	for i := 0; i < nParties; i++ {
		aggregatedShares[i] = big.NewInt(0)
		for j := 0; j < nParties; j++ {
			aggregatedShares[i].Add(aggregatedShares[i], shares[j][i]).Mod(aggregatedShares[i], q)
		}
	}

	// Any subset of k parties can reconstruct the secret
	indices := []int{1, 2, 3} // Indices of the parties (indices start from 1)
	subsetShares := []*big.Int{
		aggregatedShares[0],
		aggregatedShares[1],
		aggregatedShares[2],
	}

	// Reconstruct the secret
	secret := LagrangeInterpolation(subsetShares, indices, q)

	// Since all polynomials had zero constant term, the secret should be zero
	if secret.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("Reconstructed secret does not match expected value. Got: %s, expected: 0",
			secret.String())
	} else {
		fmt.Printf("Secret successfully reconstructed: %s\n", secret.String())
	}
}
