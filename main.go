package main

import (
	"fmt"
	"math/big"

	"github.com/niclabs/tcpaillier"
)

func main() {
	// Parameters
	bitSize := 512
	s := uint8(1)
	l := uint8(5) // Number of shares
	k := uint8(3) // Threshold

	// Step 1: Generate a new key with 5 shares
	shares, pk, err := tcpaillier.NewKey(bitSize, s, l, k)
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	// Step 2: Encrypt values 1 to 10
	values := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(4),
		big.NewInt(5),
		big.NewInt(6),
		big.NewInt(7),
		big.NewInt(8),
		big.NewInt(9),
		big.NewInt(10),
	}

	var encryptedValues []*big.Int

	// Encrypt the values
	for _, value := range values {
		encrypted, _, err := pk.Encrypt(value)
		if err != nil {
			fmt.Printf("Error encrypting value %d: %v\n", value, err)
			return
		}
		encryptedValues = append(encryptedValues, encrypted)
	}

	// Step 3: Homomorphically add the encrypted values using the standalone function
	encryptedSum := encryptedValues[0]
	for i := 1; i < len(encryptedValues); i++ {
		encryptedSum = homomorphicAdd(encryptedSum, encryptedValues[i], pk.N, s)
	}

	// Step 4: Partial decrypt with each share
	decryptionShares := make([]*tcpaillier.DecryptionShare, l)
	for i, share := range shares {
		decryptShare, err := share.PartialDecrypt(encryptedSum)
		if err != nil {
			fmt.Printf("Error decrypting share %d: %v\n", i+1, err)
			return
		}
		decryptionShares[i] = decryptShare
	}

	// Step 5: Combine the shares to get the final decrypted sum
	decryptedSum, err := pk.CombineShares(decryptionShares...)
	if err != nil {
		fmt.Printf("Error combining shares: %v\n", err)
		return
	}

	// Step 6: Print the result
	fmt.Printf("The decrypted sum of the values 1 to 10 is: %s\n", decryptedSum.String())
}
