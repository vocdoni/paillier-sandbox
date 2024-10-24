const g_input = document.getElementById('g');
const r_input = document.getElementById('r');
const n_to_s_input = document.getElementById('n_to_s');
const n_to_s_plus_one_input = document.getElementById('n_to_s_plus_one');
const msg_input = document.getElementById('msg');
const encrypt_button = document.getElementById('encrypt');
const mock_button = document.getElementById('mock_inputs');
const result_elem = document.getElementById('result');

const paillierWasm = "./paillier.wasm";

async function initWasm() {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch(paillierWasm), go.importObject);
    go.run(result.instance);
}

encrypt_button.addEventListener('click', async () => {
    await initWasm();
    const result = Paillier.encrypt(JSON.stringify({
        g: g_input.value,
        r: r_input.value,
        n_to_s: n_to_s_input.value,
        n_to_s_plus_one: n_to_s_plus_one_input.value,
        msg: msg_input.value
    }))
    result_elem.textContent = result;
});

mock_button.addEventListener('click', () => {
    g_input.value = "14783057525512922414";
    r_input.value = "83418981857195657416018508207837029889";
    n_to_s_input.value = "14783057525512922413";
    n_to_s_plus_one_input.value = "218538789802624248699744705051757742569";
    msg_input.value = "102030405";
});