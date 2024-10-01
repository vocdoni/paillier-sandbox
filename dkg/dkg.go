// == dkg.go ==
package dkg

import (
	"crypto/rand"
	"math/big"
)

// GeneratePolynomial generates a random polynomial of degree k-1 with zero constant term.
func GeneratePolynomial(k int, q *big.Int) []*big.Int {
	coeffs := make([]*big.Int, k)
	coeffs[0] = big.NewInt(0) // Zero constant term
	for i := 1; i < k; i++ {
		coeffs[i], _ = rand.Int(rand.Reader, q) // Random coefficients modulo q
	}
	return coeffs
}

// GenerateCommitments generates commitments for the polynomial coefficients.
func GenerateCommitments(coeffs []*big.Int, g, p *big.Int) []*big.Int {
	commitments := make([]*big.Int, len(coeffs))
	for i, coeff := range coeffs {
		commitments[i] = new(big.Int).Exp(g, coeff, p) // C_i = g^{a_i} mod p
	}
	return commitments
}

// GenerateShare computes the share for participant i using the polynomial.
func GenerateShare(i int, coeffs []*big.Int, q *big.Int) *big.Int {
	x := big.NewInt(int64(i))
	share := big.NewInt(0)
	for j := 0; j < len(coeffs); j++ {
		// term = a_j * x^{j} mod q
		exp := big.NewInt(int64(j))
		xExpJ := new(big.Int).Exp(x, exp, q)
		term := new(big.Int).Mul(coeffs[j], xExpJ)
		share.Add(share, term).Mod(share, q)
	}
	return share
}

// VerifyShare verifies a share using the public commitments.
func VerifyShare(share *big.Int, i int, commitments []*big.Int, g, p *big.Int) bool {
	lhs := new(big.Int).Exp(g, share, p) // lhs = g^{s_i} mod p

	rhs := big.NewInt(1)
	x := big.NewInt(int64(i))
	for j := 0; j < len(commitments); j++ {
		exp := new(big.Int).Exp(x, big.NewInt(int64(j)), p)
		commitmentExp := new(big.Int).Exp(commitments[j], exp, p)
		rhs.Mul(rhs, commitmentExp).Mod(rhs, p)
	}

	// Check if lhs == rhs
	return lhs.Cmp(rhs) == 0
}

// LagrangeInterpolation reconstructs the secret using the provided shares and indices.
func LagrangeInterpolation(shares []*big.Int, indices []int, q *big.Int) *big.Int {
	secret := big.NewInt(0)
	for i := 0; i < len(shares); i++ {
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)
		xi := big.NewInt(int64(indices[i]))
		for j := 0; j < len(shares); j++ {
			if i != j {
				xj := big.NewInt(int64(indices[j]))
				numerator.Mul(numerator, new(big.Int).Neg(xj)).Mod(numerator, q) // numerator *= -xj
				diff := new(big.Int).Sub(xi, xj)                                 // xi - xj
				denominator.Mul(denominator, diff).Mod(denominator, q)           // denominator *= xi - xj
			}
		}
		// Compute inverse of denominator modulo q
		invDenominator := new(big.Int).ModInverse(denominator, q)
		if invDenominator == nil {
			panic("Denominator has no inverse modulo q")
		}
		lagrangeCoeff := new(big.Int).Mul(numerator, invDenominator)
		lagrangeCoeff.Mod(lagrangeCoeff, q)

		term := new(big.Int).Mul(shares[i], lagrangeCoeff)
		secret.Add(secret, term).Mod(secret, q)
	}
	return secret
}

// GenerateSafePrime generates a safe prime p and its corresponding q such that p = 2q + 1.
func GenerateSafePrime(bits int) (q, p *big.Int) {
	one := big.NewInt(1)
	two := big.NewInt(2)
	for {
		qCandidate, err := rand.Prime(rand.Reader, bits-1) // q has bits-1 bits
		if err != nil {
			panic("Failed to generate random prime q")
		}
		pCandidate := new(big.Int).Mul(qCandidate, two)
		pCandidate.Add(pCandidate, one) // p = 2q + 1

		if pCandidate.ProbablyPrime(20) {
			return qCandidate, pCandidate
		}
		// Else, try again
	}
}

// FindGenerator finds a generator g of the subgroup of order q in Z_p^*.
func FindGenerator(p, q *big.Int) *big.Int {
	one := big.NewInt(1)
	for {
		// Randomly pick h in [2, p-2]
		h, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(3)))
		if err != nil {
			panic("Failed to generate random h")
		}
		h.Add(h, big.NewInt(2)) // Ensure h >= 2

		// Compute g = h^2 mod p
		g := new(big.Int).Exp(h, big.NewInt(2), p)

		// Check if g^q mod p == 1 (i.e., g ∈ subgroup of order q)
		if new(big.Int).Exp(g, q, p).Cmp(one) == 0 && g.Cmp(one) != 0 {
			return g
		}
	}
}
