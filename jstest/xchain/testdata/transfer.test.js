Test("transfer", function (t) {
        c1 = xchain.Deploy({
            name: "contract1",
            code: "/Users/chenfengjin/baidu/contract-sdk-go/example/call/c1/c1",
            lang: "go",
            type:"native",
            init_args: {},
            options: {"account": "XC1111111111111111@xuper"}
        })
        c2 = xchain.Deploy({
                name: "contract2",
                code: "/Users/chenfengjin/baidu/contract-sdk-go/example/call/c2/c2",
                lang: "go",
                init_args: {},
            type:"native",
                options: {"account": "XC1111111111111111@xuper"}
            }
        )
    resp =c2.Invoke("Invoke",{},{})
    console.log(resp.Message)
    }

)
