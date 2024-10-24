//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"math/big"
	"syscall/js"
)

// paillierEncrypt is a basic implementation of the Paillier encryption
// algorithm. It receives the public key n+1, n^s, n^(s+1), the message to
// encrypt and a random number r. It returns the resulting ciphertext.
func paillierEncrypt(nPlusOne, nToS, nToSPlusOne, msg, r *big.Int) *big.Int {
	// msg mod n^s+1
	m := new(big.Int).Mod(msg, nToSPlusOne)
	// g^m mod n^s+1
	nPlusOneToM := new(big.Int).Exp(nPlusOne, m, nToSPlusOne)
	// g^m * r^n^s mod n^s+1
	rToNToS := new(big.Int).Exp(r, nToS, nToSPlusOne)
	c := new(big.Int).Mul(nPlusOneToM, rToNToS)
	return new(big.Int).Mod(c, nToSPlusOne)
}

type paillierInputs struct {
	G           string `json:"g"`
	NToS        string `json:"n_to_s"`
	NToSPlusOne string `json:"n_to_s_plus_one"`
	Msg         string `json:"msg"`
	R           string `json:"r"`
}

func main() {
	jsClass := js.ValueOf(map[string]interface{}{})
	jsClass.Set("encrypt", js.FuncOf(func(this js.Value, args []js.Value) any {
		inputs := paillierInputs{}
		if err := json.Unmarshal([]byte(args[0].String()), &inputs); err != nil {
			return ""
		}
		g, _ := new(big.Int).SetString(inputs.G, 10)
		nToS, _ := new(big.Int).SetString(inputs.NToS, 10)
		nToSPlusOne, _ := new(big.Int).SetString(inputs.NToSPlusOne, 10)
		msg, _ := new(big.Int).SetString(inputs.Msg, 10)
		r, _ := new(big.Int).SetString(inputs.R, 10)
		c := paillierEncrypt(g, nToS, nToSPlusOne, msg, r)
		return c.String()
	}))
	js.Global().Set("Paillier", jsClass)
	select {}
}
