import BigNumber from "bignumber.js";

// Encryption key params (given)
const r  = new BigNumber("14972380668335538207658515373379008223", 10);
const n_plus_one = new BigNumber("12339172240886319110", 10);
const n_to_s_plus_one = new BigNumber("152255171590259505891894229882978553881", 10);

// Message input and expected ciphertext (to test)
const m = new BigNumber("102030405", 10);
const expected_c = new BigNumber("113892038734439757816648928291563035903", 10);

function paillier_encryption(m: BigNumber): BigNumber {
    console.log(m)
    // const g = (n_plus_one ** m) % n_to_s_plus_one;
    const g = n_plus_one.pow(m).mod(n_to_s_plus_one);
    console.log(g)
    // const r_to_n_to_s = (r ** n_to_s_plus_one) % n_plus_one;
    const r_to_n_to_s = r.pow(n_to_s_plus_one).mod(n_plus_one);
    console.log(r_to_n_to_s)
    // return (g * r_to_n_to_s) % n_to_s_plus_one;
    return g.multipliedBy(r_to_n_to_s).mod(n_to_s_plus_one);
}
const c = paillier_encryption(m)
console.log(c)
console.log(c === expected_c);