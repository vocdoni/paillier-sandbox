# Paillier Circom circuits PoC

## BigModExp

> Result = Base ^ Exp % Modulus

```
template instances: 12
non-linear constraints: 1643
linear constraints: 0
public inputs: 4
private inputs: 9
public outputs: 0
wires: 1592
labels: 4216
```

### Inputs

| Name | Pub/Priv | Type 
|:---:|:---:|:---:|
| result | `Pub` | `[]int` |
| exponent | `Priv` | `int` |
| base | `Priv` | `[]int` |
| modulus | `Priv` | `[]int` |
| result | `Priv` | `[]int` |

### Parameters

* `l_size`: Size of each limb (bigints chunks).
* `n_limbs`: Number of limbs.

### Test

```bash
sh prepare-circuit.sh bigint_test.circom
go test -timeout 30s -run ^TestBigModExp$ github.com/vocdoni/paillier-sandbox/circom -v -count=1
```

## Paillier encrypt

> E(m, r) = g^m * r^n^s mod n^s+1

```
template instances: 13
non-linear constraints: 45190
linear constraints: 0
public inputs: 16
private inputs: 49
public outputs: 0
wires: 44757
labels: 63974
```

### Inputs

| Name | Pub/Priv | Type | Description |
|:---:|:---:|:---:|:---:|
| ciphertext | `Pub` | `[]int` | The result of cipher `m` |
| m | `Priv` | `int` |  The raw message |
| n_plus_one | `Priv` | `[]int` | `g` component |
| r_to_n_to_s | `Priv` | `[]int` | `n^r^s mod n^s+1` component precalculated |
| n_to_s_plus_one | `Priv` | `[]int` | `n^s+1` component precalculated |

### Parameters

* `l_size`: Size of each limb (bigints chunks).
* `n_limbs`: Number of limbs.

### Test

```bash
sh prepare-circuit.sh paillier_cipher_test.circom
go test -timeout 30s -run ^TestPaillierCipher$ github.com/vocdoni/paillier-sandbox/circom -v -count=1
```