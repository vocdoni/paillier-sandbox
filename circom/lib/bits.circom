pragma circom 2.1.0;

template Num2Bits(n) {
    signal input in;
    signal output out[n];
    var lc = 0;
    var e = 1;
    for (var i = 0; i < n; i++) {
        out[i] <-- (in >> i) & 1;
        out[i] * (1 - out[i]) === 0;
        lc += out[i] * e;
        e *= 2;
    }
    lc === in;
}

template IsNBits(n) {
    signal input in;
    signal bits[n];
    var lc = 0;
    var e = 1;
    for (var i = 0; i < n; i++) {
        bits[i] <-- (in >> i) & 1;
        bits[i] * (1 - bits[i]) === 0;
        lc += bits[i] * e;
        e *= 2;
    }
    lc === in;
}

// assert that a signal is an N bit signed integer (+ 1 bit for sign)
template IsSignedNBits(n) {
    signal input in;
    signal bits[n + 1];

    var lc = 0;
    var e = 1;
    var top = 1 << n;
    for (var i = 0; i <= n; i++) {
        bits[i] <-- ((in + top) >> i) & 1;
        bits[i] * (1 - bits[i]) === 0;
        lc += bits[i] * e;
        e *= 2;
    }
    lc === in + top;
}

