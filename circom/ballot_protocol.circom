pragma circom 2.1.0;

include "./lib/bits.circom";
include "./lib/math.circom";
include "./lib/comparators.circom";

template BallotProtocol(n_fields) {
    signal input fields[n_fields];
    signal input max_count;
    signal input force_uniqueness;
    signal input max_value;
    signal input min_value;
    signal input max_total_cost;
    signal input min_total_cost;
    signal input cost_exp;
    // generate a mask of valid fields
    signal mask[n_fields];
    component mask_gen = MaskGenerator(n_fields);
    mask_gen.in <== max_count;
    mask <== mask_gen.out;
    // all fields must be different
    component unique = UniqueArray(n_fields);
    unique.arr <== fields;
    unique.mask <== mask;
    unique.sel <== force_uniqueness;
    // every field must be between min_value and max_value
    component inBounds = ArrayInBounds(n_fields);
    inBounds.arr <== fields;
    inBounds.mask <== mask;
    inBounds.min <== min_value;
    inBounds.max <== max_value;
    // compute total cost: sum of all fields to the power of cost_exp
    signal total_cost;
    component sum_calc = SumPow(n_fields, 252);
    sum_calc.inputs <== fields;
    sum_calc.mask <== mask;
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