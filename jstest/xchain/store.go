package xchain

// type mockStore is a implementation of contract.StateSandbox
//type mockStore struct {
//	db *leveldb.DB
//}
//
//func newMockStore() *mockStore {
//	db, err := leveldb.Open(storage.NewMemStorage(), nil)
//	if err != nil {
//		panic(err)
//	}
//	return &mockStore{
//		db: db,
//	}
//}
//
//func (m *mockStore) Get(bucket string, key []byte) ([]byte, error) {
//	return m.db.Get(makeRawKey(bucket, key), nil)
//}
//
//func (m *mockStore) Select(bucket string, startKey []byte, endKey []byte) (contract.Iterator, error) {
//	return &mockIterator{
//		Iterator: m.db.NewIterator(&util.Range{
//			Start: makeRawKey(bucket, startKey),
//			Limit: makeRawKey(bucket, endKey),
//		}, nil),
//	}, nil
//}
//
//func (m *mockStore) AddEvent(...*protos.ContractEvent) {
//	panic("not impl")
//}
//func (m *mockStore) Del(bucket string, key []byte) error {
//	return m.db.Delete(makeRawKey(bucket, key), nil)
//}
//func (m *mockStore) Flush() error {
//	return nil
//}
//func (m *mockStore) Put(bucket string, key []byte, value []byte) error {
//	return m.db.Put(makeRawKey(bucket, key), value, nil)
//}
//func (m *mockStore) RWSet() *contract.RWSet {
//	return &contract.RWSet{
//		RSet: []*ledger.VersionedData{},
//		WSet: []*ledger.PureData{},
//	}
//}
//
//func (m *mockStore) NewCache() ledger.XMReader {
//	return nil
//	//return &mockCache{store: m}
//}
//
//type mockIterator struct {
//	iterator.Iterator
//	err error
//}
//
//func (m *mockIterator) Key() []byte {
//	return bytes.Split(m.Iterator.Key(), []byte("/"))[1]
//}
//func (m *mockIterator) Value() []byte {
//	return m.Iterator.Value()
//}
//
//func (m *mockIterator) Next() bool {
//	if m.err != nil {
//		return false
//	}
//	return m.Iterator.Next()
//}
//
//func (m *mockIterator) Error() error {
//	return m.Iterator.Error()
//}
//
//func (m *mockIterator) Close() {}
//
//func makeRawKey(bucket string, key []byte) []byte {
//	buf := make([]byte, 0, len(bucket)+1+len(key))
//	buf = append(buf, bucket...)
//	buf = append(buf, '/')
//	return append(buf, key...)
//}
