import inputs from "./inputs.js";

const startBtn = document.getElementById("start");
const infoElem = document.getElementById("info");
const proofElem = document.getElementById("proof");
const publicSignalsElem = document.getElementById("public-signals");

startBtn.addEventListener("click", async function() {
    const start = Date.now();
    infoElem.textContent = "Generating proof...";
    const {proof, publicSignals} = await snarkjs.groth16.fullProve(
        inputs,
        "../circom/artifacts/vocdoni_z.wasm",
        "../circom/artifacts/vocdoni_z_pkey.zkey",
    );
    infoElem.textContent += `\nProof generation took ${Date.now() - start}ms`;
    proofElem.textContent = JSON.stringify(proof, null, 2);
    publicSignalsElem.textContent = JSON.stringify(publicSignals, null, 2);
});