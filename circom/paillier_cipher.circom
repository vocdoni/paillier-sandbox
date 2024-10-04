include "./bigint.circom";
include "./pow_mod.circom";

template EncryptWithPaillier(word_size, n_chunks, e_bits) {
    // PK {n, g}
    // E(msg, r) = g^msg * r^n^s mod n^s+1

    // input signals (precomputed public key parameters, message to cipher and fixed random value)
    signal input g[n_chunks];
    signal input n_to_s[n_chunks]; // precomputed n^s
    signal input n_to_s_plus_one[n_chunks]; // precomputed n^s+1
    signal input msg[n_chunks];
    signal input r[n_chunks];
    signal input ciphertext[n_chunks];
    // a = g^msg mod n^s+1
    component a = PowerMod(word_size, n_chunks, e_bits);
    for (var i  = 0; i < n_chunks; i++) {
        a.base[i] <== g[i];
        a.exp[i] <== msg[i];
        a.modulus[i] <== n_to_s_plus_one[i];
    }
    // b = r^n mod n^s+1
    component b = PowerMod(word_size, n_chunks, e_bits);
    for (var i  = 0; i < n_chunks; i++) {
        b.base[i] <== r[i];
        b.exp[i] <== n_to_s[i];
        b.modulus[i] <== n_to_s_plus_one[i];
    }
    // c = a * b mod n^s+1
    component c = BigMultModP(word_size, n_chunks);
    for (var i  = 0; i < n_chunks; i++) {
        c.a[i] <== a.out[i];
        c.b[i] <== b.out[i];
        c.p[i] <== n_to_s_plus_one[i];
    }
    // compare every c with ciphertext
    for (var i  = 0; i < n_chunks; i++) {
        c.out[i] === ciphertext[i];
    }
}