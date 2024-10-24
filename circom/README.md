# Paillier Circom circuits PoC

## Paillier encrypt

> E(m, r) = g^m * r^n^s mod n^s+1

The public key used in a proof can be verified checking the public inputs `g = n + 1` and `n^s+1` with `s = 1`:

```
g - 1 === √(n^s+1) === expected_n

with s = 1
```

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
For `l_size = 32` and `n_limbs = 16`.


### Inputs

| Name | Pub/Priv | Type | Description |
|:---:|:---:|:---:|:---|
| ciphertext | `Pub` | `[]int` | The result of cipher `m` |
| m | `Priv` | `int` |  The raw message |
| n_plus_one | `Pub` | `[]int` | `g` component |
| r_to_n_to_s | `Priv` | `[]int` | `n^r^s mod n^s+1` component precalculated |
| n_to_s_plus_one | `Pub` | `[]int` | `n^s+1` component precalculated |

### Parameters

* `l_size`: Size of each limb (bigints chunks).
* `n_limbs`: Number of limbs.

### Test

```bash
sh prepare-circuit.sh paillier_cipher_test.circom
go test -timeout 30s -run ^TestPaillierCipher$ github.com/vocdoni/paillier-sandbox/circom -v -count=1
```

## Ballot Protocol

![Ballot protocol example](https://blog.aragon.org/content/images/2021/04/ballot-variables-1.png)

Read more in [Ballot Protocol documentation](https://docs.vocdoni.io/architecture/data-schemes/ballot-protocol.html).


```
template instances: 13
non-linear constraints: 8639
linear constraints: 0
public inputs: 0
private inputs: 12
public outputs: 0
wires: 8613
labels: 10144
```
For `n_fields = 5`.

### Inputs

| Name | Pub/Priv | Type | Description |
|:---:|:---:|:---:|:---|
| fields | `Priv` | `[]int` | Each position of the array contains an answer to one of the process' fields. |
| max_count | `Priv` | `int` | The number of valid values of *fields*. Must be lower or equal to `n_fields` parameter. |
| force_uniqueness | `Priv` | `int` | Choices for a question cannot appear twice or more |
| max_value | `Priv` | `int` |  Determines the acceptable maximum value for all fields. |
| min_value | `Priv` | `int` | Determines the acceptable minimum value for all fields. |
| max_total_cost | `Priv` | `int` | Maximum limit on the total sum of all ballot fields' values. |
| min_total_cost | `Priv` | `int` | Minimum limit on the total sum of all ballot fields' values. |
| cost_exp | `Priv` | `int` | The exponent that will be used to compute the "cost" of the field values. |
| cost_from_weight | `Priv` | `int` | Bit to select if the circuit should use *max_total_cost* or *weight* as *total_cost* bound. |
| weight | `Priv` | `int` | The voter assigned ballot weight |

### Parameters

* `n_fields`: The number of `fields` items.

### Test

```bash
sh prepare-circuit.sh ballot_protocol_test.circom
go test -timeout 30s -run ^TestBallotProtocol$ github.com/vocdoni/paillier-sandbox/circom -v -count=1
```

## Ballot encoder

The encoding is done by the sum of every field multiplied by the power of the base and the position of the field. 

Ex.:
```
fields   = [5, 1, 4, 3, 0, 0, 0]
mask     = [1, 1, 1, 1, 0, 0, 0]
base     = 100
out      = 5 * 100^3 + 1 * 100^2 + 4 * 100^1 + 3 * 100^0 = 5010403
```

```
template instances: 6
non-linear constraints: 7070
linear constraints: 0
public inputs: 0
private inputs: 15
public outputs: 1
wires: 7086
labels: 8945
```
For `n_fields = 7`.

### Inputs

| Name | Pub/Priv | Type | Description |
|:---:|:---:|:---:|:---|
| fields | `Priv` | `[]int` | Ballot protocol fields |
| mask | `Priv` | `[]int` | A list of bits indicating which positions contain a valid field. |
| base | `Priv` | `int` | To be used as the base of the power to encode each field. |

### Parameters

* `n_fields`: The number of `fields` items.

### Test

```bash
sh prepare-circuit.sh ballot_encoder_test.circom
go test -timeout 30s -run ^TestBallotEncoder$ github.com/vocdoni/paillier-sandbox/circom -v -count=1
```

## VocdoniZ circuit

VocdoniZ is the circuit to prove a valid vote in the Vocdoni scheme. The vote is valid if it meets the Ballot Protocol requirements, but also if the encrypted vote provided matches with the raw vote encrypted in this circuit.
The circuit checks the the vote over the params provided using the BallotProtocol template, encodes the vote using the BallotEncoder template and compares the result with the encrypted vote.

```
template instances: 104
non-linear constraints: 134160
linear constraints: 0
public inputs: 35
private inputs: 15
public outputs: 0
wires: 132758
labels: 189042
```
For `n_fields = 5`, `l_size = 32` and `n_limbs = 8`.

### Inputs

| Name | Pub/Priv | Type | Description |
|:---:|:---:|:---:|:---|
| fields | `Priv` | `[]int` | Each position of the array contains an answer to one of the process' fields. |
| max_count | `Pub` | `int` | The number of valid values of *fields*. Must be lower or equal to `n_fields` parameter. |
| force_uniqueness | `Pub` | `int` | Choices for a question cannot appear twice or more |
| max_value | `Pub` | `int` |  Determines the acceptable maximum value for all fields. |
| min_value | `Pub` | `int` | Determines the acceptable minimum value for all fields. |
| max_total_cost | `Pub` | `int` | Maximum limit on the total sum of all ballot fields' values. |
| min_total_cost | `Pub` | `int` | Minimum limit on the total sum of all ballot fields' values. |
| cost_exp | `Pub` | `int` | The exponent that will be used to compute the "cost" of the field values. |
| cost_from_weight | `Pub` | `int` | Bit to select if the circuit should use *max_total_cost* or *weight* as *total_cost* bound. |
| weight | `Pub` | `int` | The voter assigned ballot weight |
| base | `Pub` | `int` | To be used as the base of the power to encode each field. |
| ciphertext | `Pub` | `[]int` | The result of cipher the encoded ballot. |
| n_plus_one | `Pub` | `[]int` | `g` component |
| r_to_n_to_s | `Priv` | `[]int` | `n^r^s mod n^s+1` component precalculated |
| n_to_s_plus_one | `Pub` | `[]int` | `n^s+1` component precalculated |
| nullifier | `Pub` | `int` | The voter nullifier: *Hash(commitment, secret)* |
| commitment | `Priv` | `int` | The first part of the pre-image of the nullifier |
| secret | `Priv` | `int` | The second part of the pre-image of the nullifier |

### Parameters

* `n_fields`: The number of `fields` items.
* `l_size`: Size of each limb (bigints chunks).
* `n_limbs`: Number of limbs.

### Test

```bash
sh prepare-circuit.sh vocdoni_z.circom
go test -timeout 3m -run ^TestVocdoniZ$ github.com/vocdoni/paillier-sandbox/circom -v -count=1
```