package xchain

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/storage"
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
	//start, end := makeRawKey(bucket, startKey), makeRawKey(bucket, endKey)
	//iter := m.db.NewIterator(&util.Range{
	//	Start: start,
	//	Limit: end,
	//}, nil)
	//newMockIterator(iter).Value()
	//return newMockIterator(iter), nil
	return nil, nil
}

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

func (m *mockIterator) First() bool {
	if m.err != nil {
		return false
	}
	ok := m.Iterator.First()
	if !ok {
		return false
	}
	//err := proto.Unmarshal(m.Iterator.Value(), &m.data)
	//if err != nil {
	//	m.err = err
	//	return false
	//}
	return true
}

func (m *mockIterator) Next() bool {
	if m.err != nil {
		return false
	}
	ok := m.Iterator.Next()
	if !ok {
		return false
	}

	//err := proto.Unmarshal(m.Iterator.Value(), &m.data)
	//if err != nil {
	//	m.err = err
	//	return false
	//}
	return true
}

func (m *mockIterator) Close() {

}
func (m *mockIterator) Error() error {
	if m.err != nil {
		return m.err
	}
	return m.Iterator.Error()
}

func (m *mockIterator) Data() *ledger.VersionedData {
	return &ledger.VersionedData{}
}

func (m *mockIterator) Value() *ledger.VersionedData {
	return &ledger.VersionedData{}
}
func (m *mockStore) AddEvent(...*protos.ContractEvent) {

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
}

//
//type mockIterator struct {
//}

func (m *mockCache) Get(string, []byte) (*ledger.VersionedData, error) {
	return &ledger.VersionedData{}, nil
}
func (m *mockCache) Select(bucket string, startKey []byte, endKey []byte) (ledger.XMIterator, error) {
	return &mockIterator{}, nil
}

func (m *mockStore) NewCache() ledger.XMReader {
	return &mockCache{}
}
