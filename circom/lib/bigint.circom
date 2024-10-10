pragma circom 2.0.0;

include "bits.circom";
include "comparators.circom";

function varmax(a, b) {
    if (a > b) {
        return a;
    } else {
        return b;
    }
}

function varmin(a, b) {
    if (a < b) {
        return a;
    } else {
        return b;
    }
}

// large integer addition
// given two large integers A and B
// where A is an large integer of Na limbs
// and B is an large integer of Nb limbs
// compute A + B as an large integer of max(Na, Nb) limbs of size max(na, nb) + 1
template BigAdd(Na, Nb) {
    signal input a[Na];
    signal input b[Nb];
    signal output out[varmax(Na, Nb)];
    for (var i = 0; i < varmin(Na, Nb); i++) {
        out[i] <== a[i] + b[i];
    }
    if (Na > Nb) {
        for (var i = Nb; i < Na; i++) {
            out[i] <== a[i];
        }
    } else {
        for (var i = Na; i < Nb; i++) {
            out[i] <== b[i];
        }
    }
}

// given large integer A with N limbs
// convert to a normalized large integer (N + 1) limbs with n bits per limb
// helper function, doesn't generate constraints
function BigVarNormalize(N, n, A) {
    var B[200]; // N + 1
    var r = 0;
    for (var i = 0; i < N; i++) {
        B[i] = ((A[i] + r) % (2 ** n));
        r = (A[i] + r) >> n;
    }
    B[N] = r;
    return B;
}

// multiply two large integers A and B and normalize the result (carries)
function BigVarMul(Na, Nb, a, b, n) {
    var out[200]; // Na + Nb - 1
    for (var i = 0; i < Na + Nb - 1; i++) {
        out[i] = 0;
    }
    for (var i = 0; i < Na; i++) {
        for (var j = 0; j < Nb; j++) {
            out[i + j] = out[i + j] + a[i] * b[j];
        }
    }
    return BigVarNormalize(Na + Nb - 1, n, out);
}

// large integer multiplication
// given two large integers A and B
// where A is an large integer of Na limbs
// and B is an large integer of Nb limbs
// compute A * B as an large integer of Na + Nb - 1 limbs
template BigMul(Na, Nb) {
    signal input a[Na];
    signal input b[Nb];
    signal output out[Na + Nb - 1];
    var ab[Na + Nb - 1];
    // first calculate C without constraints (polynomial multiplication)
    for (var i = 0; i < Na + Nb - 1; i++) {
        ab[i] = 0;
    }
    for (var i = 0; i < Na; i++) {
        for (var j = 0; j < Nb; j++) {
            ab[i + j] = ab[i + j] + a[i] * b[j];
        }
    }
    for (var i = 0; i < Na + Nb - 1; i++) {
        out[i] <-- ab[i];
    }
    // now add constraints to ensure that C is the correct result
    var s_a[Na + Nb - 1]; // A as a polynomial evaluated at 0, 1, 2, ...
    var s_b[Na + Nb - 1]; // B as a polynomial evaluated at 0, 1, 2, ...
    var s_c[Na + Nb - 1]; // C as a polynomial evaluated at 0, 1, 2, ...
    for (var i = 0; i < Na + Nb - 1; i++) {
        s_a[i] = 0; s_b[i] = 0; s_c[i] = 0;
        for (var j = 0; j < Na; j++) {
            s_a[i] = s_a[i] + a[j] * (i ** j);
        }
        for (var j = 0; j < Nb; j++) {
            s_b[i] = s_b[i] + b[j] * (i ** j);
        }
        for (var j = 0; j < Na + Nb - 1; j++) {
            s_c[i] = s_c[i] + out[j] * (i ** j);
        }
        s_c[i] === s_a[i] * s_b[i];
    }
}

// given large integers A and M (A potentially not normalized)
// compute R = A % M and Q = A / M
// where R is an large integer of Nm limbs
// and Q is an large integer of Na - Nm + 1 limbs
// does not generate constraints, just computes the result
// TO DO, maybe turn this into a function instead of a template to avoid warnings
template BigVarDiv(Na, Nm, n) {
    signal input a[Na];
    signal input mod[Nm];
    signal output out[Nm];
    signal output q[Na - Nm + 1];
    var anorm[200] = BigVarNormalize(Na, n, a); // Na + 1
    var r = 0;
    for (var ii = 0; ii <= Na - Nm; ii++) {
        var i = Na - Nm - ii;
        // determine the ith limb of Q by binary search
        // this could be swapped for something faster but its not the bottleneck
        var l = 0;
        var h = 2 ** n;
        while (l + 1 < h) {
            var m = (l + h) / 2;
            var tmp[200] = BigVarMul(Nm, 1, mod, [m], n); // Nm + 1
            // check if tmp <= anorm[i:]
            var larger = 0;
            var done = 0;
            for (var jj = 0; jj <= Nm && done == 0; jj++) {
                var j = Nm - jj;
                if ((i + j > Na && tmp[j] > 0) ||
                    (i + j <= Na && tmp[j] > anorm[i + j])) {
                    larger = 1;
                    done = 1;
                } else if (i + j <= Na && tmp[j] < anorm[i + j]) {
                    done = 1;
                }
            }
            if (larger == 0) {
                l = m;
            } else {
                h = m;
            }
        }
        q[i] <-- l;
        // subtract M * Q[i] from a slice of anorm
        var tmp[200] = BigVarMul(Nm, 1, mod, [q[i]], n); // Nm + 1
        for (var j = 0; j < Nm + 1 && i + j < Na + 1; j++) {
            if (anorm[i + j] < tmp[j]) {
                anorm[i + j] = anorm[i + j] + (2 ** n) - tmp[j];
                tmp[j + 1] = tmp[j + 1] + 1;
            } else {
                anorm[i + j] = anorm[i + j] - tmp[j];
            }
        }
    }
    // copy the lower Nm limbs of anorm into R
    for (var i = 0; i < Nm; i++) {
        out[i] <-- anorm[i];
    }
}

template BigLimbCheck(N, n) {
    signal input a[N];
    component isbits[N];
    for (var i = 0; i < N; i++) {
        isbits[i] = IsNBits(n);
        isbits[i].in <== a[i];
    }
}

// only use after a single multiplication, not generic!
template BigEq(N, n) {
    signal input a[N];
    signal input b[N];
    // get ceil log2 N
    var logN = 0;
    var tmp = N;
    while (tmp > 0) {
        logN++;
        tmp >>= 1;
    }
    var k = ((253 - logN) \ n) - 1;
    var l = ((N - 1) \ k);
    signal Zero[l];
    for (var i = 0; k * i < N; i++) {
        var lc = 0;
        if (i != 0) {
            lc = Zero[i - 1];
        }
        for (var j = 0; j < k && k * i + j < N; j++) {
            lc += (a[k * i + j] - b[k * i + j]) * (2 ** (n * j));
        }
        if (k * (i + 1) < N) {
            Zero[i] <== lc / (2 ** (n * k));
        } else {
            lc === 0;
        }
    }
    component safe[l];
    for (var i = 0; i < l; i++) {
        safe[i] = IsSignedNBits(n + logN);
        safe[i].in <== Zero[i];
    }
}

// given large integers A and M
// where A is an large integer of Na limbs
// and M is an large integer of Nm limbs
// compute R = A % M with some flexibility in the size of C
// namely R can be larger than M, but it must fit within Nm limbs of size n
template BigRelaxMod(Na, Nm, n) {
    signal input a[Na];
    signal input mod[Nm];
    signal output out[Nm];
    signal q[Na - Nm + 1];
    // A = Q * M + R
    // first calculate Q and R without constraints via long division
    component longDiv = BigVarDiv(Na, Nm, n);
    longDiv.a <== a;
    longDiv.mod <== mod;
    q <== longDiv.q;
    out <== longDiv.out;
    // check that A = Q * M + R
    component mul = BigMul(Na - Nm + 1, Nm);
    mul.a <== q;
    mul.b <== mod;
    component add = BigAdd(Na, Nm);
    add.a <== mul.out;
    add.b <== out;
    // check equality
    component eq = BigEq(Na, n);
    eq.a <== add.out;
    eq.b <== a;
    // check that R fits in Nm limbs of with n bits each
    component limbCheckR = BigLimbCheck(Nm, n);
    limbCheckR.a <== out;
    // check that Q fits in Na - Nm + 1 limbs of with n bits each
    component limbCheckQ = BigLimbCheck(Na - Nm + 1, n);
    limbCheckQ.a <== q;
}

template BigModMul(N, n) {
    signal input a[N];
    signal input b[N];
    signal input mod[N];
    signal output out[N];
    component mul = BigMul(N, N);
    mul.a <== a;
    mul.b <== b;
    component bigMod = BigRelaxMod(N + N - 1, N, n);
    bigMod.a <== mul.out;
    bigMod.mod <== mod;
    out <== bigMod.out;
}

// Large integer modular exponentiation
// Given a large integer A and modulus M, both as arrays of N limbs of size n
// Compute A^e mod M for some integer e (input as a signal), where e is up to K bits
template BigModExp(N, n, K) {
    // Inputs
    signal input base[N];          // Base (N limbs)
    signal input mod[N];          // Modulus (N limbs)
    signal input exp;             // Exponent (single limb for simplicity)
    signal output out[N];         // Result (N limbs)

    // Step 1: Convert the exponent to bits
    component bitsizer = Num2Bits(K);
    bitsizer.in <== exp;
    signal e_bits[K];
    for (var i = 0; i < K; i++) {
        e_bits[i] <== bitsizer.out[i];
    }

    // Step 2: Initialize R and currentA at iteration 0
    signal R[K+1][N];
    signal currentA[K+1][N];

    for (var j = 0; j < N; j++) {
        R[0][j] <== (j == 0) ? 1 : 0;    // R[0] = 1 represented in limbs
        currentA[0][j] <== base[j];      // currentA[0] = base
    }

    // Step 3: Pre-declare all required components
    component modmul[2 * K];
    component muxes[K][N];    // One Mux1 per bit per limb

    // Step 4: Iterate over each bit of the exponent
    for (var i = 0; i < K; i++) {
        // 4.1: Conditional Multiplication: R[i+1] = E_bits[i] ? (R[i] * currentA[i] mod M) : R[i]
        modmul[2 * i] = BigModMul(N, n);
        modmul[2 * i].a <== R[i];
        modmul[2 * i].b <== currentA[i];
        modmul[2 * i].mod <== mod;

        // 4.2: Use Mux1 to select between R[i] and modmul[2*i].C for each limb
        for (var j = 0; j < N; j++) {
            muxes[i][j] = Mux();
            muxes[i][j].a <== R[i][j];
            muxes[i][j].b <== modmul[2 * i].out[j];
            muxes[i][j].sel <== e_bits[i];
            R[i + 1][j] <== muxes[i][j].out;
        }

        // 4.3: Square the base: currentA[i+1] = (currentA[i] * currentA[i]) mod M
        modmul[2 * i + 1] = BigModMul(N, n);
        modmul[2 * i + 1].a <== currentA[i];
        modmul[2 * i + 1].b <== currentA[i];
        modmul[2 * i + 1].mod <== mod;

        for (var j = 0; j < N; j++) {
            currentA[i + 1][j] <== modmul[2 * i + 1].out[j];
        }
    }

    // Step 5: Assign the final result
    for (var j = 0; j < N; j++) {
        out[j] <== R[K][j];
    }
}