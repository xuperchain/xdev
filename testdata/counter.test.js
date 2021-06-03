var assert = require("assert");

var codePath = "counter-c.wasm";

var lang = "c"
var type = "wasm"
function deploy() {
    return xchain.Deploy({
        name: "counter",
        code: codePath,
        lang: lang,
        type: type,
        init_args: { "creator": "xchain" },
        options: { "account": "XC1111111111111111@xuper" }

    });
}

Test("Increase", function (t) {
    var c = deploy();
    var resp = c.Invoke("increase", { "key": "xchain" }, { "name": "11111" });
    assert.equal(resp.Body, "1");
})

Test("Get", function (t) {
    var c = deploy()
    c.Invoke("increase", { "key": "xchain" });
    var resp = c.Invoke("get", { "key": "xchain" })
    assert.equal(resp.Body, "1")
})
