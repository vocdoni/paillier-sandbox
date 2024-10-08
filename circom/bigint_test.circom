pragma circom 2.1.0;

include "./bigint.circom";

template BigModExpTest(l_size, n_limbs) {
    signal input exponent;
    signal input base[n_limbs];
    signal input modulus[n_limbs];
    signal input result[n_limbs];

    component powMod = BigModExp(n_limbs, l_size, 10);
    powMod.exp <== exponent;
    powMod.base <== base;
    powMod.mod <== modulus;
    for (var i = 0; i < n_limbs; i++) {
        powMod.out[i] === result[i];
    }
}

component main {public [result]} = BigModExpTest(8, 4);