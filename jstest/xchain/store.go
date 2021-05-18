package xchain

import (
	"bytes"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/kernel/ledger"
	"github.com/xuperchain/xupercore/protos"
)

type mockStore struct {
	db *leveldb.DB
}

func newMockStore() *mockStore {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		panic(err)
	}
	return &mockStore{
		db: db,
	}
}

func makeRawKey(bucket string, key []byte) []byte {
	buf := make([]byte, 0, len(bucket)+1+len(key))
	buf = append(buf, bucket...)
	buf = append(buf, '/')
	return append(buf, key...)
}

func (m *mockStore) Get(bucket string, key []byte) ([]byte, error) {
	return m.db.Get(makeRawKey(bucket, key), nil)
}

func (m *mockStore) Select(bucket string, startKey []byte, endKey []byte) (contract.Iterator, error) {
	return &mockLedgerIterator{
		Iterator: m.db.NewIterator(&util.Range{
			Start: makeRawKey(bucket, startKey),
			Limit: makeRawKey(bucket, endKey),
		}, nil),
	}, nil
}

type mockLedgerIterator struct {
	iterator.Iterator
	err error
}

func (m *mockLedgerIterator) Key() []byte {
	//TODO @fengjin
	return bytes.Split(m.Iterator.Key(), []byte("/"))[1]
}
func (m *mockLedgerIterator) Value() []byte {
	return m.Iterator.Value()
}

func (m *mockLedgerIterator) Next() bool {
	if m.err != nil {
		return false
	}
	return m.Iterator.Next()
}

func (m *mockLedgerIterator) Error() error {
	return m.Iterator.Error()
}

func (m *mockLedgerIterator) Close() {}

type mockIterator struct {
	iterator.Iterator

	data ledger.VersionedData
	err  error
}

func newMockIterator(iter iterator.Iterator) ledger.XMIterator {
	return &mockIterator{
		Iterator: iter,
	}
}

func (m *mockIterator) Next() bool {
	if m.err != nil {
		return false
	}
	return m.Iterator.Next()

}

func (m *mockIterator) Close() {

}
func (m *mockIterator) Error() error {
	if m.err != nil {
		return m.err
	}
	return m.Iterator.Error()
}

func (m *mockIterator) Value() *ledger.VersionedData {
	return &ledger.VersionedData{}
}
func (m *mockStore) AddEvent(...*protos.ContractEvent) {
	panic("not impl")
}
func (m *mockStore) Del(string, []byte) error {
	return nil
}
func (m *mockStore) Flush() error {
	return nil
}
func (m *mockStore) Put(bucket string, key []byte, value []byte) error {
	return m.db.Put(makeRawKey(bucket, key), value, nil)
}
func (m *mockStore) RWSet() *contract.RWSet {
	return &contract.RWSet{
		RSet: []*ledger.VersionedData{},
		WSet: []*ledger.PureData{},
	}
}

type mockCache struct {
	store *mockStore
}

func (m *mockStore) NewCache() ledger.XMReader {
	return &mockCache{store: m}
}
func (m *mockCache) Get(string, []byte) (*ledger.VersionedData, error) {
	return &ledger.VersionedData{}, nil
}

func (m *mockCache) Select(bucket string, startKey []byte, endKey []byte) (ledger.XMIterator, error) {
	return &mockIterator{
		Iterator: m.store.db.NewIterator(&util.Range{makeRawKey(bucket, startKey), makeRawKey(bucket, endKey)}, nil),
		data: ledger.VersionedData{
			PureData: &ledger.PureData{
				Bucket: "",
				Key:    nil,
				Value:  nil,
			},
			RefTxid:   []byte(""),
			RefOffset: 0,
		},
		err: nil,
	}, nil
}
