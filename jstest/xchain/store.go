package xchain

import (
	"crypto/rand"
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/kernel/contract/sandbox"
	"github.com/xuperchain/xupercore/kernel/ledger"
)

const (
	ContractAccount = "XC1111111111111111@xuper"
)

type mockStore struct {
	state *sandbox.MemXModel
}

func NewmockStore() *mockStore {
	state := sandbox.NewMemXModel()
	store := &mockStore{
		state: state,
	}
	return store
}

func (t *mockStore) State() ledger.XMReader {
	return t.state
}

func (t *mockStore) Commit(state contract.StateSandbox) {
	rwset := state.RWSet()
	txbuf := make([]byte, 32)
	rand.Read(txbuf)
	for i, w := range rwset.WSet {
		t.state.Put(w.Bucket, w.Key, &ledger.VersionedData{
			RefTxid:   txbuf,
			RefOffset: int32(i),
			PureData: &ledger.PureData{
				Bucket: w.Bucket,
				Key:    w.Key,
				Value:  w.Value,
			},
		})
	}
}
