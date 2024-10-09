pragma circom 2.1.0;

include "./lib/bigint.circom";

template EncryptWithPaillier(l_size, n_limbs) {
    // E(m, r) = g^m * r^n^s mod n^s+1

    // private inputs
    signal input m; // plain text
    signal input n_plus_one[n_limbs]; // g = n + 1
    signal input r_to_n_to_s[n_limbs]; // r^n^s mod n^s+1
    signal input n_to_s_plus_one[n_limbs]; // n^s+1
    // public inputs
    signal input ciphertext[n_limbs];

    // compute g^m mod n^s+1 
    component powMod = BigModExp(n_limbs, l_size, 17);
    powMod.base <== n_plus_one;
    powMod.exp <== m;
    powMod.mod <== n_to_s_plus_one;

    // compute g^m * r^n^s mod n^s+1
    component mulMod = BigModMul(n_limbs, l_size);
    mulMod.a <== powMod.out;
    mulMod.b <== r_to_n_to_s;
    mulMod.mod <== n_to_s_plus_one;

    // check result with public input
    for (var i = 0; i < n_limbs; i++) {
        mulMod.out[i] === ciphertext[i];
    }
}

component main {public [ciphertext, n_plus_one, n_to_s_plus_one]} = EncryptWithPaillier(32, 16);