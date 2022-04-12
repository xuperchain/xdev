module github.com/xuperchain/xdev

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/ddliu/motto v0.3.1
	github.com/golang/protobuf v1.4.3
	github.com/mitchellh/mapstructure v1.1.2
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/spf13/cobra v1.1.3
	github.com/xuperchain/xupercore v0.0.0-20211223100656-b02fb7b21ce1
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)

replace github.com/hyperledger/burrow => github.com/xuperchain/burrow v0.30.6-0.20211229032028-fbee6a05ab0f

replace github.com/xuperchain/xvm => github.com/xuperchain/xvm v0.0.0-20220412083010-7737da252598

replace github.com/xuperchain/xupercore => github.com/xuperchain/xupercore v0.0.0-20220323054134-beddc96c027d
