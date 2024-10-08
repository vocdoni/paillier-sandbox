pragma circom 2.1.0;

include "./bits.circom";
include "./unique.circom";

template Pow(n) {
    signal input base;
    signal input exp_bits[n];
    signal output out;
    // Initialize intermediate results
    signal intermediates[n+1];
    intermediates[0] <== 1; // Start with 1

    signal squares[n];
    signal multipliers[n];
    for (var i = 0; i < n; i++) {
        var bit_index = n - 1 - i; // Start from MSB
        squares[i] <== intermediates[i] * intermediates[i];
        multipliers[i] <== 1 + exp_bits[bit_index] * (base - 1);
        intermediates[i+1] <== squares[i] * multipliers[i];
    }
    out <== intermediates[n];
}

template Sum(n) {
    signal input inputs[n];
    signal output out;

    signal intermediates[n+1];
    intermediates[0] <== 0;
    for (var i = 0; i < n; i++) {
        intermediates[i+1] <== intermediates[i] + inputs[i];
    }
    out <== intermediates[n];
}

template SumPow(n, e_bits) {
    signal input inputs[n];
    signal input exp;
    signal output out;

    component n2b = Num2Bits(e_bits);
    n2b.in <== exp;

    signal powers[n];
    component pow[n];
    for (var i = 0; i < n; i++) {
        pow[i] = Pow(e_bits);
        pow[i].base <== inputs[i];
        pow[i].exp_bits <== n2b.out;
        powers[i] <== pow[i].out;
    }

    component sum = Sum(n);
    sum.inputs <== powers;
    out <== sum.out;
}

template BallotProtocol(max_count) {
    signal input fields[max_count];
    signal input max_value;
    signal input min_value;
    signal input max_total_cost;
    signal input min_total_cost;
    signal input cost_exp;
    // all fields must be different and every field must be between min_value and max_value
    component unique = UniqueArrayInBounds(max_count);
    unique.arr <== fields;
    unique.min <== min_value;
    unique.max <== max_value;
    // the sum of every field power cost_exp must be between min_total_cost and max_total_cost
    var sum = 0;
    for (var i = 0; i < max_count; i++) {
        sum += fields[i]**cost_exp;
    }
    // compute total cost
    signal total_cost;
    component sum_calc = SumPow(max_count, 252);
    sum_calc.inputs <== fields;
    sum_calc.exp <== cost_exp;
    total_cost <== sum_calc.out;
    // check bounds
    component lt = LessThan(252);
    lt.in[0] <== total_cost;
    lt.in[1] <== max_total_cost;
    lt.out === 1;
    component gt = GreaterThan(252);
    gt.in[0] <== total_cost;
    gt.in[1] <== min_total_cost;
    gt.out === 1;
}