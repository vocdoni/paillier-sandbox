package dkg

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

// Utility function to generate a random big integer
func randomBigInt() *big.Int {
	// Generate a large prime number for modulus (e.g., 512 bits)
	p, _ := rand.Prime(rand.Reader, 512)
	q, _ := rand.Prime(rand.Reader, 512)
	n := new(big.Int).Mul(p, q) // n = p * q
	return n
}

func TestGeneratePolynomial(t *testing.T) {
	n := randomBigInt() // Modulus
	secret := big.NewInt(123456789)
	k := 3 // Threshold

	// Generate a polynomial of degree k-1
	poly := GeneratePolynomial(k, secret, n)

	// Check if the first coefficient is the secret
	if poly[0].Cmp(secret) != 0 {
		t.Fatalf("Expected first coefficient to be the secret, got %s", poly[0].String())
	}

	// Check if the polynomial has the correct number of coefficients
	if len(poly) != k {
		t.Fatalf("Expected %d coefficients, but got %d", k, len(poly))
	}
}

func TestGenerateCommitments(t *testing.T) {
	n := randomBigInt() // Modulus
	g := big.NewInt(2)  // Generator
	secret := big.NewInt(123456789)
	k := 3 // Threshold

	// Generate a polynomial and its commitments
	poly := GeneratePolynomial(k, secret, n)
	commitments := GenerateCommitments(poly, g, n)

	// Check if the commitments are the correct size
	if len(commitments) != k {
		t.Fatalf("Expected %d commitments, but got %d", k, len(commitments))
	}

	// Verify that the first commitment is g^secret mod n
	expectedCommitment := new(big.Int).Exp(g, secret, n)
	if commitments[0].Cmp(expectedCommitment) != 0 {
		t.Fatalf("First commitment doesn't match expected value: got %s, expected %s",
			commitments[0].String(), expectedCommitment.String())
	}
}

func TestGenerateShare(t *testing.T) {
	n := randomBigInt() // Modulus
	secret := big.NewInt(123456789)
	k := 3 // Threshold

	// Generate a polynomial
	poly := GeneratePolynomial(k, secret, n)

	// Generate a share for party 1
	share := GenerateShare(1, poly, n)

	// Manually calculate the expected share for party 1 (f(1))
	expectedShare := big.NewInt(0)
	x := big.NewInt(1) // i = 1
	for j := len(poly) - 1; j >= 0; j-- {
		term := new(big.Int).Exp(x, big.NewInt(int64(j)), n)
		term.Mul(term, poly[j])
		expectedShare.Add(expectedShare, term).Mod(expectedShare, n)
	}

	if share.Cmp(expectedShare) != 0 {
		t.Fatalf("Generated share doesn't match expected value: got %s, expected %s",
			share.String(), expectedShare.String())
	}
}

func TestVerifyShare(t *testing.T) {
	n := randomBigInt() // Modulus
	g := big.NewInt(2)  // Generator
	secret := big.NewInt(123456789)
	k := 3 // Threshold

	// Generate a polynomial and commitments
	poly := GeneratePolynomial(k, secret, n)
	commitments := GenerateCommitments(poly, g, n)

	// Generate a share for party 1
	share := GenerateShare(1, poly, n)

	// Verify the share using the commitments
	valid := VerifyShare(share, 1, commitments, g, n)
	if !valid {
		t.Fatalf("Share verification failed for party 1")
	}

	// Tamper with the share and verify again (should fail)
	tamperedShare := new(big.Int).Add(share, big.NewInt(1))
	valid = VerifyShare(tamperedShare, 1, commitments, g, n)
	if valid {
		t.Fatalf("Tampered share verification succeeded when it should have failed")
	}
}

func TestEndToEnd(t *testing.T) {
	// Step 1: Generate parameters
	n := randomBigInt()             // Modulus
	g := big.NewInt(2)              // Generator
	k := 3                          // Threshold
	nParties := 5                   // Number of parties
	secret := big.NewInt(123456789) // Arbitrary secret for testing

	// Step 2: Each party generates its own secret polynomial and commitments
	polynomials := make([][]*big.Int, nParties)
	commitments := make([][]*big.Int, nParties)

	for i := 0; i < nParties; i++ {
		// Each party generates its own polynomial
		polynomials[i] = GeneratePolynomial(k, secret, n)

		// Each party generates its commitments
		commitments[i] = GenerateCommitments(polynomials[i], g, n)

		fmt.Printf("Party %d's commitments: %v\n", i+1, commitments[i])
	}

	// Step 3: Each party generates shares for all other parties
	shares := make([][]*big.Int, nParties) // shares[i][j] is party i's share for party j
	for i := 0; i < nParties; i++ {
		shares[i] = make([]*big.Int, nParties)
		for j := 1; j <= nParties; j++ {
			shares[i][j-1] = GenerateShare(j, polynomials[i], n)
			fmt.Printf("Party %d's share for party %d: %s\n", i+1, j, shares[i][j-1].String())
		}
	}

	// Step 4: Verify shares using commitments
	for i := 0; i < nParties; i++ {
		for j := 0; j < nParties; j++ {
			valid := VerifyShare(shares[i][j], j+1, commitments[i], g, n)
			if !valid {
				t.Fatalf("Share verification failed for party %d's share for party %d", i+1, j+1)
			} else {
				fmt.Printf("Party %d's share for party %d verified successfully\n", i+1, j+1)
			}
		}
	}

	// Step 5: Combine the shares to reconstruct the secret (using interpolation)

	// Use the shares of the first threshold number of parties (k = 3 here)
	thresholdShares := make([]*big.Int, k)
	for i := 0; i < k; i++ {
		thresholdShares[i] = shares[i][i] // Use the i-th party's share for itself
	}

	// Lagrange interpolation to reconstruct the secret
	combinedSecret := CombineShares(thresholdShares, n)

	// Step 6: Verify if the reconstructed secret matches the original
	if combinedSecret.Cmp(secret) != 0 {
		t.Fatalf("Reconstructed secret does not match the original. Got: %s, expected: %s",
			combinedSecret.String(), secret.String())
	} else {
		fmt.Printf("Secret successfully reconstructed: %s\n", combinedSecret.String())
	}
}
