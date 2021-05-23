package xchain

//
//import (
//	"github.com/syndtr/goleveldb/leveldb/iterator"
//	"github.com/syndtr/goleveldb/leveldb/util"
//	"github.com/xuperchain/xupercore/kernel/ledger"
//)
//
//type mockCacheIterator struct {
//	iterator.Iterator
//
//	data ledger.VersionedData
//	err  error
//}
//
//func (m *mockCacheIterator) Next() bool {
//	if m.err != nil {
//		return false
//	}
//	return m.Iterator.Next()
//}
//
//func (m *mockCacheIterator) Close() {
//
//}
//func (m *mockCacheIterator) Error() error {
//	if m.err != nil {
//		return m.err
//	}
//	return m.Iterator.Error()
//}
//
//func (m *mockCacheIterator) Value() *ledger.VersionedData {
//	return &ledger.VersionedData{}
//}
//
//// mockCache is a implementation of contract.XMReader
//type mockCache struct {
//	store *mockStore
//}
//
//func (m *mockCache) Get(string, []byte) (*ledger.VersionedData, error) {
//	return &ledger.VersionedData{}, nil
//}
//
//func (m *mockCache) Select(bucket string, startKey []byte, endKey []byte) (ledger.XMIterator, error) {
//	return &mockCacheIterator{
//		Iterator: m.store.db.NewIterator(&util.Range{makeRawKey(bucket, startKey), makeRawKey(bucket, endKey)}, nil),
//		data: ledger.VersionedData{
//			PureData: &ledger.PureData{
//				Bucket: "",
//				Key:    nil,
//				Value:  nil,
//			},
//			RefTxid:   []byte(""),
//			RefOffset: 0,
//		},
//		err: nil,
//	}, nil
//}
