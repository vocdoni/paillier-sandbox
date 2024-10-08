pragma circom 2.0.0;

template Num2Bits(n) {
    signal input in;
    signal output out[n];
    var lc1=0;

    var e2=1;
    for (var i = 0; i<n; i++) {
        out[i] <-- (in >> i) & 1;
        out[i] * (out[i] -1 ) === 0;
        lc1 += out[i] * e2;
        e2 = e2+e2;
    }

    lc1 === in;
}

// given an N-bit number,
// returns a list of N bits
template SigToBinary(N) {
    signal input val;
    signal output bits[N];
    var lc = 0;
    var e = 1;
    for (var i = 0; i < N; i++) {
        bits[i] <-- (val >> i) & 1;
        bits[i] * (1 - bits[i]) === 0;
        lc += bits[i] * e;
        e *= 2;
    }
    lc === val;
}

template IsNBits(N) {
    signal input val;
    signal bits[N];
    var lc = 0;
    var e = 1;
    for (var i = 0; i < N; i++) {
        bits[i] <-- (val >> i) & 1;
        bits[i] * (1 - bits[i]) === 0;
        lc += bits[i] * e;
        e *= 2;
    }
    lc === val;
}

// assert that a signal is an N bit signed integer (+ 1 bit for sign)
template IsSignedNBits(N) {
    signal input val;
    signal bits[N + 1];
    var lc = 0;
    var e = 1;
    var top = 1 << N;
    for (var i = 0; i <= N; i++) {
        bits[i] <-- ((val + top) >> i) & 1;
        bits[i] * (1 - bits[i]) === 0;
        lc += bits[i] * e;
        e *= 2;
    }
    lc === val + top;
}

