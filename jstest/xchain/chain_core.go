package xchain

import (
	"errors"
	"github.com/xuperchain/xupercore/bcs/ledger/xledger/state"
	pb "github.com/xuperchain/xupercore/bcs/ledger/xledger/xldgpb"
	pb2 "github.com/xuperchain/xupercore/kernel/contract/bridge/pb"
	"github.com/xuperchain/xupercore/kernel/ledger"
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
func (c *chainCore) QueryTransaction(txid []byte) (*pb2.Transaction, error) {
	return &pb2.Transaction{}, nil
}

// QueryBlock query block
func (c *chainCore) QueryBlock(blockid []byte) (ledger.BlockHandle, error) {
	return &state.BlockAgent{}, nil
}

// QueryBlockByHeight query block by height
func (c *chainCore) QueryBlockByHeight(height int64) (*pb.InternalBlock, error) {
	return new(pb.InternalBlock), nil
}

// QueryLastBlock query last block
func (c *chainCore) QueryLastBlock() (*pb.InternalBlock, error) {
	return new(pb.InternalBlock), nil
}

// CrossQuery query contract from otherchain
//func (c *chainCore) ResolveChain(chainName string) (*pb.CrossQueryMeta, error) {
//	return new(pb.CrossQueryMeta), nil
//}
