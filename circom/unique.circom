pragma circom 2.1.0;

include "./bits.circom";

template IsZero() {
    signal input in;
    signal output out;

    signal inv;
    inv <-- in!=0 ? 1/in : 0;
    out <== -in*inv +1;
    in*out === 0;
}

template IsEqual() {
    signal input in[2];
    signal output out;

    component isz = IsZero();
    in[1] - in[0] ==> isz.in;
    isz.out ==> out;
}

template LessThan(n) {
    assert(n <= 252);
    signal input in[2];
    signal output out;

    component n2b = Num2Bits(n+1);
    n2b.in <== in[0]+ (1<<n) - in[1];
    out <== 1-n2b.out[n];
}

template GreaterThan(n) {
    signal input in[2];
    signal output out;

    component lt = LessThan(n);
    lt.in[0] <== in[1];
    lt.in[1] <== in[0];
    lt.out ==> out;
}

template UniqueArrayInBounds(n) {
    signal input arr[n];
    signal input min;
    signal input max;
    // enforce each pair of elements is distinct
    component iseq[n][n];
    component lt[n];
    component gt[n];
    for (var i = 0; i < n; i++) {
        for (var j = i + 1; j < n; j++) {
            if (j != i) {
                iseq[i][j] = IsEqual();
                iseq[i][j].in[0] <== arr[i];
                iseq[i][j].in[1] <== arr[j];
                iseq[i][j].out === 0;
            }
        }
        // enforce each element is in bounds
        lt[i] = LessThan(252);
        lt[i].in[0] <== arr[i];
        lt[i].in[1] <== max;
        lt[i].out === 1;

        gt[i] = GreaterThan(252);
        gt[i].in[0] <== arr[i];
        gt[i].in[1] <== min;
        gt[i].out === 1;
    }
}