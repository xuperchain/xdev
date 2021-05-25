module github.com/xuperchain/xdev

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/ddliu/motto v0.3.1
	github.com/golang/protobuf v1.4.2
	github.com/mitchellh/mapstructure v1.1.2
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/spf13/cobra v1.1.3
	github.com/xuperchain/xuperchain v0.0.0-20210511082518-b2d6bd248cc3
	github.com/xuperchain/xupercore v0.0.0-20210427110058-ad573e2286cc
)

replace github.com/hyperledger/burrow => github.com/xuperchain/burrow v0.30.6-0.20210317023017-369050d94f4a

replace github.com/xuperchain/xupercore => ../xupercore
