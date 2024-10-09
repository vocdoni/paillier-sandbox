pragma circom 2.1.0;

include "./lib/bits.circom";
include "./lib/comparators.circom";
include "./lib/math.circom";

// BallotEncoder is a template that encodes a ballot with n_fields fields
// into a single integer. The encoding is done by multiplying each field
// by a power of base and summing the results.
// Ex.:
//  fields   = [5, 1, 4, 3, 0, 0, 0]
//  n_fields = 7
//  mask     = [1, 1, 1, 1, 0, 0, 0]
//  base     = 100
//  out      = 5 * 100^3 + 1 * 100^2 + 4 * 100^1 + 3 * 100^0 = 5010403
template BallotEncoder(n_fields) {
    signal input fields[n_fields];
    signal input mask[n_fields];
    signal input base;
    signal output out;
    // calculate the difference between the position of the valid fields and the
    // exponent of the base
    signal exp_diff[n_fields+1];
    exp_diff[0] <== 0;
    component iseq[n_fields];
    for (var i = 0; i < n_fields; i++) {
        iseq[i] = IsEqual();
        iseq[i].in[0] <== mask[i];
        iseq[i].in[1] <== 1;
        exp_diff[i+1] <== iseq[i].out + exp_diff[i];
    }
    assert(exp_diff[n_fields] >= 0);
    var exp = n_fields - exp_diff[n_fields];

    signal powers[n_fields];
    component pow[n_fields];
    component n2b[n_fields];
    for (var i = 0; i < n_fields; i++) {
        n2b[i] = Num2Bits_unsafe(252);
        n2b[i].in <== exp;

        pow[i] = Pow(252);
        pow[i].base <== base;
        pow[i].exp_bits <== n2b[i].out;

        powers[i] <== pow[i].out * fields[i];
        exp -= mask[i];
    }

    component sum = Sum(n_fields);
    sum.inputs <== powers;
    sum.mask <== mask;
    out <== sum.out;
}

component main = BallotEncoder(7);