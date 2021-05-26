package xchain

import "github.com/xuperchain/xupercore/kernel/ledger"

type MockLedgerRely struct {
	reader ledger.XMSnapshotReader
}

type XMSnapshotReader struct {
}

func (lr *MockLedgerRely) GetNewAccountGas() (int64, error) {
	return 0, nil
}
func (lr *MockLedgerRely) GetTipXMSnapshotReader() (ledger.XMSnapshotReader, error) {
	return lr.reader, nil
}

func (*XMSnapshotReader) Get(bucket string, key []byte) ([]byte, error) {
	return []byte(""), nil
}
