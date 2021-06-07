package xchain

import (
	"errors"
	"github.com/xuperchain/xupercore/bcs/ledger/xledger/xldgpb"
	"github.com/xuperchain/xupercore/kernel/contract/bridge/pb"

	"math/big"
)

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
func (c *chainCore) QueryTransaction(txid []byte) (*pb.Transaction, error) {
	return new(pb.Transaction), nil
}

// QueryBlock query block
func (c *chainCore) QueryBlock(blockid []byte) (*xldgpb.InternalBlock, error) {
	return new(xldgpb.InternalBlock), nil
}

// QueryBlockByHeight query block by height
func (c *chainCore) QueryBlockByHeight(height int64) (*xldgpb.InternalBlock, error) {
	return new(xldgpb.InternalBlock), nil
}

// QueryLastBlock query last block
func (c *chainCore) QueryLastBlock() (*xldgpb.InternalBlock, error) {
	return new(xldgpb.InternalBlock), nil
}

func (c *chainCore) Transfer(from string, to string, amount *big.Int) error {
	return nil
}

// CrossQuery query contract from otherchain
//func (c *chainCore) ResolveChain(chainName string) (*pb.CrossQueryMeta, error) {
//	return new(pb.CrossQueryMeta), nil
//}
