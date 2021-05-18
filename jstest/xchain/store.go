package xchain

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
	xmodel "github.com/xuperchain/xupercore/kernel/contract/sandbox"
	"github.com/xuperchain/xupercore/kernel/ledger"
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

func (m *mockStore) Get(bucket string, key []byte) (*ledger.VersionedData, error) {
	//value, err := m.db.Get(makeRawKey(bucket, key), nil)
	//if err != nil {
	//	return nil, err
	//}
	//return ledger.VersionedData{RefTxid:}
	//data := new(ledger.VersionedData)
	//err = proto.Unmarshal(value, data)
	//if err != nil {
	//	return nil, err
	//}
	//return data, nil
	return nil,nil
}

func (m *mockStore) Select(bucket string, startKey []byte, endKey []byte) (ledger.XMIterator, error) {
	start, end := makeRawKey(bucket, startKey), makeRawKey(bucket, endKey)
	iter := m.db.NewIterator(&util.Range{
		Start: start,
		Limit: end,
	}, nil)
	//newMockIterator(iter).Value()
	return newMockIterator(iter), nil
}

func (m *mockStore) Commit(cache *xmodel.XMCache) error {
	//txid := make([]byte, 32)
	//rand.Read(txid)
	//
	//batch := new(leveldb.Batch)
	//wset := cache.RWSet().WSet
	//for i, w := range wset {
	//	rawKey := makeRawKey(w.GetBucket(), w.GetKey())
	//	value, _ := proto.Marshal(&ledger.VersionedData{
	//		RefTxid:   txid,
	//		RefOffset: int32(i),
	//		PureData:  w,
	//	})
	//	batch.Put(rawKey, value)
	//}
return nil
	//return m.db.Write(batch, nil)
}

func (m *mockStore) NewCache() *xmodel.XMCache {
	cache := xmodel.NewXModelCache(m)
	return cache
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

func(m*mockIterator)Close(){

}
func (m *mockIterator) Error() error {
	if m.err != nil {
		return m.err
	}
	return m.Iterator.Error()
}

func (m *mockIterator) Data() *ledger.VersionedData {
	return &m.data
}

func(m *mockIterator) Value()* ledger.VersionedData{
	//TODO
	return nil
}