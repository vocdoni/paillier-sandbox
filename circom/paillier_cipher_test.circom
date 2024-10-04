pragma circom 2.0.0;

include "./paillier_cipher.circom";

component main {public [g, n_to_s, n_to_s_plus_one, msg, r]} = EncryptWithPaillier(32, 16, 64);