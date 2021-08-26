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
	github.com/xuperchain/xupercore v0.0.0-20210820071654-4a9d712921d2
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)

replace github.com/hyperledger/burrow => github.com/xuperchain/burrow v0.30.6-0.20210317023017-369050d94f4a
