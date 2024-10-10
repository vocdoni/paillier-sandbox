pragma circom 2.1.0;

include "./ballot_protocol.circom";
include "./ballot_encoder.circom";
include "./paillier_cipher.circom";

// VocdoniZ is the circuit to prove a valid vote in the Vocdoni scheme. The 
// vote is valid if it meets the Ballot Protocol requirements, but also if the
// encrypted vote provided matches with the raw vote encrypted in this circuit.
// The circuit checks the the vote over the params provided using the 
// BallotProtocol template, encodes the vote using the BallotEncoder template
// and compares the result with the encrypted vote.
template VocdoniZ(n_fields, l_size, n_limbs) {
    // BallotProtocol inputs
    signal input fields[n_fields];  // private
    signal input max_count;         // public
    signal input force_uniqueness;  // public
    signal input max_value;         // public
    signal input min_value;         // public
    signal input max_total_cost;    // public
    signal input min_total_cost;    // public
    signal input cost_exp;          // public
    signal input cost_from_weight;  // public
    signal input weight;            // public
    // BallotEncoder inputs
    signal input base;              // public
    // EncryptWithPaillier inputs
    signal input n_plus_one[n_limbs];       // public
    signal input r_to_n_to_s[n_limbs];      // private
    signal input n_to_s_plus_one[n_limbs];  // public
    signal input ciphertext[n_limbs];       // public
    // 1. Check the vote meets the Ballot Protocol requirements
    component ballotProtocol = BallotProtocol(n_fields);
    ballotProtocol.fields <== fields;
    ballotProtocol.max_count <== max_count;
    ballotProtocol.force_uniqueness <== force_uniqueness;
    ballotProtocol.max_value <== max_value;
    ballotProtocol.min_value <== min_value;
    ballotProtocol.max_total_cost <== max_total_cost;
    ballotProtocol.min_total_cost <== min_total_cost;
    ballotProtocol.cost_exp <== cost_exp;
    ballotProtocol.cost_from_weight <== cost_from_weight;
    ballotProtocol.weight <== weight;
    // 2. Encode the vote
    component ballotEncoder = BallotEncoder(n_fields);
    ballotEncoder.fields <== fields;
    ballotEncoder.mask <== ballotProtocol.mask; // mask of valid fields
    ballotEncoder.base <== base; 
    // 3. Check the encrypted vote
    component encryptWithPaillier = EncryptWithPaillier(l_size, n_limbs);
    encryptWithPaillier.m <== ballotEncoder.out; // encoded vote from BallotEncoder
    encryptWithPaillier.n_plus_one <== n_plus_one;
    encryptWithPaillier.r_to_n_to_s <== r_to_n_to_s;
    encryptWithPaillier.n_to_s_plus_one <== n_to_s_plus_one;
    encryptWithPaillier.ciphertext <== ciphertext;
}

component main{public [max_count, force_uniqueness, max_value, min_value, max_total_cost, min_total_cost, cost_exp, cost_from_weight, weight, base, n_plus_one, n_to_s_plus_one, ciphertext]} = VocdoniZ(10, 32, 16);