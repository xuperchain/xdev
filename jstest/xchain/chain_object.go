package xchain

import (
	"encoding/hex"
	"fmt"
	"github.com/xuperchain/xdev/jstest"
	"github.com/xuperchain/xupercore/bcs/contract/evm/abi"
	"github.com/xuperchain/xupercore/kernel/contract/bridge"
	"io/ioutil"
)

type xchainObject struct {
	env *environment
}

func (*xchainObject) GetAccountAddresses(accountName string) ([]string, error) {
	return []string{}, nil
}

func (*xchainObject) VerifyContractPermission(initiator string, authRequire []string, contractName, methodName string) (bool, error) {
	return true, nil
}
func (*xchainObject) VerifyContractOwnerPermission(contractName string, authRequire []string) error {
	return nil
}
func newXchainObject() (*xchainObject, error) {
	env, err := newEnvironment()
	if err != nil {
		return nil, err
	}
	return &xchainObject{
		env: env,
	}, nil
}

func (x *xchainObject) Contract(name string) *contractObject {
	if !x.env.ContractExists(name) {
		jstest.Throw(fmt.Errorf("contract %s not found", name))
	}
	return &contractObject{
		Name: name,
		env:  x.env,
	}
}

func (x *xchainObject) Deploy(args deployArgs) *contractObject {
	codeBuf, err := ioutil.ReadFile(args.Code)
	if err != nil {
		jstest.Throw(err)
	}

	if args.Type == string(bridge.TypeEvm) {
		dst, err := hex.DecodeString(string(codeBuf))
		if err != nil {
			jstest.Throw(err)
		}
		codeBuf = dst
	}

	args.codeBuf = codeBuf
	if args.Type == string(bridge.TypeEvm) && args.ABIFile == "" {
		jstest.Throws("missing abi")
	}
	var enc *abi.ABI
	if args.ABIFile != "" {
		buf, err := ioutil.ReadFile(args.ABIFile)
		if err != nil {
			jstest.Throw(err)
		}
		enc, err = abi.New(buf)
		if err != nil {
			jstest.Throw(err)
		}
	}
	var trueArgs map[string][]byte
	if args.ABIFile != "" {
		input, err := enc.Encode("", args.InitArgs)
		if err != nil {
			jstest.Throw(err)
		}
		codeBuf = append(codeBuf, input...)
	} else {
		trueArgs = convertArgs(args.InitArgs)
	}
	args.trueArgs = trueArgs
	args.codeBuf = codeBuf

	_, err = x.env.Deploy(args)
	if err != nil {
		jstest.Throw(err)
	}

	return &contractObject{
		env:  x.env,
		abi:  enc,
		Name: args.Name,
		Type: args.Type,
	}
}
