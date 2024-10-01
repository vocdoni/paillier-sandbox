# DKG (Distributed Key Generation) Package

## Introduction

This DKG (Distributed Key Generation) package provides an implementation of a cryptographic protocol that allows multiple parties to generate and share a private key in a decentralized manner. The key feature of DKG is that it eliminates the need for a trusted party to distribute the secret key. Instead, the secret key is distributed across multiple participants such that any subset of parties above a certain threshold can jointly perform cryptographic operations (like decryption or signing), but no subset of parties below the threshold can reconstruct the key or perform those operations.

### Key Features:
- **Threshold Cryptography**: A private key is distributed across `n` parties, and a threshold `k` is defined. Any group of `k` or more parties can jointly perform cryptographic operations (e.g., decryption or signing).
- **Distributed and Decentralized**: No trusted dealer is needed to distribute the key shares.
- **Secure**: The protocol is designed to prevent any coalition of less than `k` parties from reconstructing the secret key.

## How It Works

The DKG protocol works by each party independently generating a polynomial whose constant term is the party's secret contribution to the final key. The parties then exchange commitments to the coefficients of their polynomials, which allows each party to verify the correctness of the shares they receive from others. 

The protocol can be broken down into the following steps:
1. **Polynomial Generation**: Each party generates a random polynomial of degree `t-1` (where `t` is the threshold) with a secret as the constant term.
2. **Share Distribution**: Each party computes shares of its polynomial for every other party and sends them to the respective parties.
3. **Commitment Generation**: Each party generates public commitments to the coefficients of its polynomial to allow others to verify the correctness of the shares.
4. **Share Verification**: Each party verifies the shares it receives using the public commitments.
5. **Key Reconstruction**: Once enough shares have been collected, a subset of parties can jointly reconstruct the key using Lagrange interpolation.

### Example Usage
Here’s a basic flow of how this DKG implementation might be used:
```go
// Initialize DKG parameters
n := randomBigInt() // Modulus
g := big.NewInt(2)  // Generator
k := 3 // Threshold
nParties := 5 // Number of parties
secret := big.NewInt(123456789)

// Generate polynomials, shares, and commitments for each party
polynomials := make([][]*big.Int, nParties)
commitments := make([][]*big.Int, nParties)
shares := make([][]*big.Int, nParties)

for i := 0; i < nParties; i++ {
    poly := GeneratePolynomial(k, secret, n)
    polynomials[i] = poly
    commitments[i] = GenerateCommitments(poly, g, n)
    
    for j := 0; j < nParties; j++ {
        shares[j] = append(shares[j], GenerateShare(j+1, poly, n))
    }
}

// Each party verifies the received shares using the commitments
for i := 0; i < nParties; i++ {
    for j := 0; j < nParties; j++ {
        valid := VerifyShare(shares[j][i], j+1, commitments[i], g, n)
        if !valid {
            log.Fatalf("Share verification failed for party %d from party %d", j+1, i+1)
        }
    }
}
```

## The Mathematics Behind DKG

### Polynomial Secret Sharing

At the core of DKG is the use of **polynomial secret sharing** (based on Shamir's secret sharing scheme). Each party generates a random polynomial of degree `t-1` (where `t` is the threshold number of parties required to reconstruct the key). The polynomial is of the form:

```
f(x) = a_0 + a_1*x + a_2*x^2 + ... + a_{t-1}*x^{t-1} mod n
```

- `a_0` is the secret that the party contributes to the final shared key.
- `a_1, a_2, ..., a_{t-1}` are randomly generated coefficients.
- The modulus `n` is a large prime number to ensure that all operations are done within a finite field.

Each party generates `n` shares of its polynomial, one for each of the other participants, using the formula:

```
share_i = f(i) mod n
```

### Commitments

To prevent malicious parties from sending inconsistent or incorrect shares, each party generates **commitments** to its polynomial coefficients. The commitment to a coefficient `a_j` is:

```
C_j = g^{a_j} mod n
```

Where `g` is a publicly known generator. The commitments allow other parties to verify the correctness of the shares they receive.

### Share Verification

Once a party receives a share from another party, it can verify the share using the public commitments. To verify a share `s` from party `j` for party `i`, the following check is made:

```
g^s mod n == product(C_j^(i^j) mod n for j = 0 to t-1)
```

This ensures that the share corresponds to the committed polynomial.

### Reconstruction

When at least `t` shares are gathered, the secret can be reconstructed using **Lagrange interpolation**. The secret is the value of the polynomial at `x = 0`, which can be computed using the formula:

```
secret = sum(share_i * lagrange_coeff(i, S) mod n for i in S)
```

Where `S` is the set of indices of the parties whose shares are being used for reconstruction, and `lagrange_coeff(i, S)` is the Lagrange coefficient for index `i` with respect to set `S`.

## Assumptions and Security

### Assumptions

1. **Finite Field**: All operations are performed modulo a large prime `n` to ensure the computations are in a finite field.
2. **Threshold Security**: The protocol assumes that at least `t` parties must collaborate to reconstruct the secret. Any fewer than `t` parties do not have enough information to reconstruct the secret.
3. **Randomness**: The coefficients of the polynomials are chosen randomly to ensure that no information about the secret is leaked from the shares.

### Security Justification

1. **Threshold Security**: The security of the scheme is based on the properties of polynomial interpolation. Since the polynomial is of degree `t-1`, any group of fewer than `t` parties has no information about the constant term `a_0` (the secret). This ensures that an adversary cannot reconstruct the secret unless they have at least `t` shares.
   
2. **Zero Knowledge Proof**: The commitment scheme ensures that parties cannot forge or tamper with shares. Since the commitments are public and based on the polynomial coefficients, any incorrect share will fail the verification process.

3. **Distributed Trust**: No single party holds the secret key. Instead, the secret is distributed across the parties, and only a threshold number of parties can reconstruct the secret. This eliminates the need for a trusted dealer to distribute the key.

4. **Collusion Resistance**: The protocol is resistant to collusion as long as fewer than `t` parties collaborate. Even if up to `t-1` parties are compromised, they do not have enough information to reconstruct the secret.
