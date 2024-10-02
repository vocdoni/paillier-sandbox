package tests

import (
	"math/big"
	"testing"

	"github.com/niclabs/tcpaillier"
)

var one = big.NewInt(1)

func dummyEncrypt(g, nToS, nToSPlusOne, msg, r *big.Int) *big.Int {
	// (n+1)^m % n^(s+1)
	m := new(big.Int).Mod(msg, nToSPlusOne)
	nPlusOneToM := new(big.Int).Exp(g, m, nToSPlusOne)
	// r^(n^s) % n^(s+1)
	rToNToS := new(big.Int).Exp(r, nToS, nToSPlusOne)
	// (n+1)^m * r^(n^s) % n^(s+1)
	c := new(big.Int).Mul(nPlusOneToM, rToNToS)
	c.Mod(c, nToSPlusOne)
	return c
}

func TestDummy(t *testing.T) {
	// set parameters
	bitSize := 254
	s := uint8(1)
	l := uint8(5) // number of shares
	k := uint8(3) // threshold
	// generate the public key
	_, pk, err := tcpaillier.NewKey(bitSize, s, l, k)
	if err != nil {
		t.Errorf("Error generating key: %v\n", err)
		return
	}
	// get a random mod
	r, err := pk.RandomModNToSPlusOneStar()
	if err != nil {
		t.Errorf("Error generating random mod: %v\n", err)
		return
	}
	// encrypt with r
	raw := big.NewInt(1234567890)
	c1, err := pk.EncryptFixed(raw, r)
	if err != nil {
		t.Errorf("Error encrypting: %v\n", err)
		return
	}
	// s
	bigS := big.NewInt(int64(pk.S))
	// n+1
	nPlusOne := new(big.Int).Add(pk.N, one)
	// (s+1)
	sPlusOne := new(big.Int).Add(bigS, one)
	// n^s
	nToS := new(big.Int).Exp(pk.N, bigS, nil)
	// n^(s+1)
	nToSPlusOne := new(big.Int).Exp(pk.N, sPlusOne, nil)

	// cipher raw with dummy encryption and compare with c1
	if c2 := dummyEncrypt(nPlusOne, nToS, nToSPlusOne, raw, r); c1.Cmp(c2) != 0 {
		t.Error("Ciphertexts are different")
	}
}
