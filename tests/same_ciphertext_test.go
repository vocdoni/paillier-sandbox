package tests

import (
	"math/big"
	"testing"

	"github.com/niclabs/tcpaillier"
)

func TestSameCipherWithSameR(t *testing.T) {
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
	c1, proof, err := pk.EncryptFixedWithProof(raw, r)
	if err != nil {
		t.Errorf("Error encrypting: %v\n", err)
		return
	}
	// encrypt again with r
	c2, err := pk.EncryptFixed(raw, r)
	if err != nil {
		t.Errorf("Error encrypting: %v\n", err)
		return
	}
	// compare the two ciphertexts
	if c1.Cmp(c2) != 0 {
		t.Error("Ciphertexts are different")
	}
	// verify the proof
	if err := proof.Verify(pk, c1); err != nil {
		t.Errorf("Error verifying proof: %v\n", err)
	}
}
