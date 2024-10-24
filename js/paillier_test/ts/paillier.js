// Encryption key params (given)
var r = globalThis.BigInt("14972380668335538207658515373379008223");
var n_plus_one = globalThis.BigInt("12339172240886319110");
var n_to_s_plus_one = globalThis.BigInt("152255171590259505891894229882978553881");
// Message input and expected ciphertext (to test)
var m = globalThis.BigInt("102030405");
var expected_c = globalThis.BigInt("113892038734439757816648928291563035903");
function paillier_encryption(m) {
    var g = (Math.pow(n_plus_one, m)) % n_to_s_plus_one;
    var r_to_n_to_s = (Math.pow(r, n_to_s_plus_one)) % n_plus_one;
    return (g * r_to_n_to_s) % n_to_s_plus_one;
}
console.log(paillier_encryption(m) === expected_c);
