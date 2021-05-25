package xchain

import "github.com/xuperchain/xupercore/kernel/ledger"

type LedgerRely struct {
	reader ledger.XMSnapshotReader
}


type XMSnapshotReader struct {
}
func (XMSnapshotReader)Get(bucket string, key []byte) ([]byte, error){
	return []byte(""),nil
}

func (lr*LedgerRely)GetNewAccountGas() (int64, error){
	return 0,nil
}
func(lr*LedgerRely)	GetTipXMSnapshotReader() (ledger.XMSnapshotReader, error){
	return lr.reader,nil
}