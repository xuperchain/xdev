package xchain

import (
	"fmt"
	"github.com/xuperchain/xupercore/bcs/contract/evm/abi"
	"github.com/xuperchain/xupercore/kernel/contract/bridge"
	"io/ioutil"
	"math/big"
	"testing"

	"encoding/hex"
	"errors"
	"github.com/xuperchain/xdev/jstest"
)

type xchainAdapter struct {
}

// NewAdapter is the xchain adapter
func NewAdapter() jstest.Adapter {
	return new(xchainAdapter)
}

func (x *xchainAdapter) OnSetup(r *jstest.Runner) {
	r.GlobalObject().Set("Xchain", func() *xchainObject {
		x, err := newXchainObject()
		if err != nil {
			jstest.Throw(err)
		}
		return x
	})
}

func (x *xchainAdapter) OnTeardown(r *jstest.Runner) {
}

func (x *xchainAdapter) OnTestCase(r *jstest.Runner, test jstest.TestCase) jstest.TestCase {
	body := func(t *testing.T) {
		xctx, err := newXchainObject()
		if err != nil {
			t.Fatal(err)
		}
		defer xctx.env.Close()

		if !r.Option.Quiet {
			// TODO: add log output
		}
		// reset xchain environment
		r.GlobalObject().Set("xchain", xctx)

		test.F(t)
	}
	return jstest.TestCase{
		Name: test.Name,
		F:    body,
	}
}

var (
	errUnimplemented = errors.New("unimplemented")
)

type chainCore struct {
}

// GetAccountAddress get addresses associated with account name
func (c *chainCore) GetAccountAddresses(accountName string) ([]string, error) {
	return []string{}, nil
}

// GetBalance get balance from utxo
func (c *chainCore) GetBalance(addr string) (*big.Int, error) {
	return big.NewInt(0), nil
}

// VerifyContractPermission verify permission of calling contract
func (c *chainCore) VerifyContractPermission(initiator string, authRequire []string, contractName, methodName string) (bool, error) {
	return true, nil
}

// VerifyContractOwnerPermission verify contract ownership permisson
func (c *chainCore) VerifyContractOwnerPermission(contractName string, authRequire []string) error {
	return nil
}

// QueryTransaction query confirmed tx
//func (c *chainCore) QueryTransaction(txid []byte) (*pb.Transaction, error) {
//	return new(pb.Transaction), nil
//}

// QueryBlock query block
//func (c *chainCore) QueryBlock(blockid []byte) (*pb.InternalBlock, error) {
//	return new(pb.InternalBlock), nil
//}

// QueryBlockByHeight query block by height
//func (c *chainCore) QueryBlockByHeight(height int64) (*pb.InternalBlock, error) {
//	return new(pb.InternalBlock), nil
//}

// QueryLastBlock query last block
//func (c *chainCore) QueryLastBlock() (*pb.InternalBlock, error) {
//	return new(xledgerpb.InternalBlock), nil
//}

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
