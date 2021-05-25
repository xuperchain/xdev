package xchain

import (
	"crypto/rand"
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/kernel/contract/sandbox"
	"github.com/xuperchain/xupercore/kernel/ledger"
	"github.com/xuperchain/xupercore/kernel/permission/acl/utils"
	"io/ioutil"
)

const (
	ContractAccount = "XC1111111111111111@xuper"
)

type mockStore struct {
	basedir string

	state *sandbox.MemXModel
}

func NewmockStore() *mockStore {
	basedir, err := ioutil.TempDir("", "xdev-test")
	if err != nil {
		panic(err)
	}
	state := sandbox.NewMemXModel()
	store := &mockStore{
		basedir: basedir,
		state:   state,
	}
	store.initAccount()
	return store
}

func (t *mockStore) State() ledger.XMReader {
	return t.state
}

func (t *mockStore) initAccount() {
	t.state.Put(utils.GetAccountBucket(), []byte(ContractAccount), &ledger.VersionedData{
		RefTxid: []byte("txid"),
	})
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
