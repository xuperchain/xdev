module github.com/xuperchain/xdev

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/ddliu/motto v0.3.1
	github.com/golang/protobuf v1.4.3
	github.com/mitchellh/mapstructure v1.1.2
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/spf13/cobra v1.1.3
	github.com/xuperchain/log15 v0.0.0-20190620081506-bc88a9198230
	github.com/xuperchain/xupercore v0.0.0-20211223100656-b02fb7b21ce1
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)

replace github.com/hyperledger/burrow => github.com/xuperchain/burrow v0.30.6-0.20211229032028-fbee6a05ab0f

replace github.com/xuperchain/xvm => github.com/xuperchain/xvm v0.0.0-20220225110211-bd25eb4d8997

replace github.com/xuperchain/xupercore => github.com/xuperchain/xupercore v0.0.0-20220225071354-5439bf8c4bf5
