package dkg

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GeneratePolynomial generates a polynomial of degree k-1
func GeneratePolynomial(k int, secret *big.Int, n *big.Int) []*big.Int {
	coeffs := make([]*big.Int, k)
	coeffs[0] = secret // constant term (the secret)
	for i := 1; i < k; i++ {
		coeffs[i], _ = rand.Int(rand.Reader, n) // random coefficients
	}
	return coeffs
}

// GenerateCommitments generates commitments for the polynomial coefficients
func GenerateCommitments(coeffs []*big.Int, g *big.Int, n *big.Int) []*big.Int {
	commitments := make([]*big.Int, len(coeffs))
	for i, coeff := range coeffs {
		commitments[i] = new(big.Int).Exp(g, coeff, n) // g^a_i mod n
	}
	return commitments
}

// GenerateShare generates a share for participant i from the polynomial
func GenerateShare(i int, coeffs []*big.Int, n *big.Int) *big.Int {
	x := big.NewInt(int64(i))
	share := big.NewInt(0)
	for j := len(coeffs) - 1; j >= 0; j-- {
		term := new(big.Int).Exp(x, big.NewInt(int64(j)), n) // i^j
		term.Mul(term, coeffs[j])
		share.Add(share, term).Mod(share, n)
	}
	return share
}

// LagrangeInterpolation reconstructs the secret using Lagrange interpolation
func LagrangeInterpolation(shares map[int]*big.Int, n *big.Int) *big.Int {
	secret := big.NewInt(0)
	for i, share := range shares {
		li := big.NewInt(1)
		for j := range shares {
			if i != j {
				num := big.NewInt(int64(-j))
				den := big.NewInt(int64(i - j))
				den.ModInverse(den, n)
				li.Mul(li, num).Mul(li, den).Mod(li, n)
			}
		}
		term := new(big.Int).Mul(share, li)
		secret.Add(secret, term).Mod(secret, n)
	}
	return secret
}

// CombineShares combines the threshold shares to reconstruct the secret using Lagrange interpolation.
func CombineShares(shares []*big.Int, n *big.Int) *big.Int {
	secret := big.NewInt(0)

	// Perform Lagrange interpolation to compute the secret
	for i := 0; i < len(shares); i++ {
		// Lagrange coefficient for the i-th share
		lagrangeCoeff := big.NewInt(1)

		for j := 0; j < len(shares); j++ {
			if i != j {
				num := big.NewInt(int64(-j - 1))                              // -x_j
				denom := big.NewInt(int64(i - j))                             // x_i - x_j
				denom.ModInverse(denom, n)                                    // (x_i - x_j)^-1 mod n
				lagrangeCoeff.Mul(lagrangeCoeff, num).Mod(lagrangeCoeff, n)   // Multiply numerator
				lagrangeCoeff.Mul(lagrangeCoeff, denom).Mod(lagrangeCoeff, n) // Multiply denominator (inverted mod n)
			}
		}

		// Add the contribution of the i-th share (lagrangeCoeff * share_i mod n)
		term := new(big.Int).Mul(lagrangeCoeff, shares[i])
		term.Mod(term, n)
		secret.Add(secret, term).Mod(secret, n)
	}

	return secret
}

// VerifyShare verifies the correctness of a share using the public commitments.
func VerifyShare(share *big.Int, i int, commitments []*big.Int, g *big.Int, n *big.Int) bool {
	// lhs = g^share mod n
	lhs := new(big.Int).Exp(g, share, n)

	// Compute rhs = Product of (C_j^i^j mod n)
	rhs := big.NewInt(1)
	x := big.NewInt(int64(i))

	for j := 0; j < len(commitments); j++ {
		// Compute C_j^(i^j) mod n
		exponent := new(big.Int).Exp(x, big.NewInt(int64(j)), n)
		commitmentPower := new(big.Int).Exp(commitments[j], exponent, n)
		rhs.Mul(rhs, commitmentPower).Mod(rhs, n)
	}

	// Ensure both lhs and rhs are reduced mod n
	lhs.Mod(lhs, n)
	rhs.Mod(rhs, n)

	// Debug: Print lhs and rhs for troubleshooting
	fmt.Printf("lhs (g^share mod n) = %s\n", lhs.String())
	fmt.Printf("rhs (from commitments) = %s\n", rhs.String())

	// Check if lhs == rhs
	return lhs.Cmp(rhs) == 0
}
