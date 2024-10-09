pragma circom 2.1.0;

include "./paillier_cipher.circom";

component main {public [ciphertext, n_plus_one, n_to_s_plus_one]} = EncryptWithPaillier(32, 16);