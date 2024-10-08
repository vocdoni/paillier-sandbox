pragma circom 2.1.0;

include "./paillier_cipher.circom";

component main {public [ciphertext]} = EncryptWithPaillier(32, 16);