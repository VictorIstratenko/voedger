/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package teststore

import (
	"bytes"
	"context"

	"github.com/voedger/voedger/pkg/istorage"
	"github.com/voedger/voedger/pkg/istorageimpl"
	"github.com/voedger/voedger/pkg/istructs"
)

// Test storage. Trained to return the specified error
type (
	scheduleStorageError struct {
		err         error
		pKey, cCols []byte
	}

	damagedStorageFunc    func(*[]byte)
	scheduleStorageDamage struct {
		dam damagedStorageFunc
		scheduleStorageError
	}

	TestMemStorage struct {
		storage  istorage.IAppStorage
		get, put scheduleStorageError
		damage   scheduleStorageDamage
	}

	testStorageProvider struct {
		testStorage *TestMemStorage
	}
)

func (tsp *testStorageProvider) AppStorage(appName istructs.AppQName) (structs istorage.IAppStorage, err error) {
	return tsp.testStorage, nil
}

func NewStorageProvider(ts *TestMemStorage) istorage.IAppStorageProvider {
	return &testStorageProvider{testStorage: ts}
}

func NewStorage() *TestMemStorage {
	s := TestMemStorage{get: scheduleStorageError{}, put: scheduleStorageError{}}
	asf := istorage.ProvideMem()
	sp := istorageimpl.Provide(asf)
	var err error
	if s.storage, err = sp.AppStorage(istructs.AppQName_test1_app1); err != nil {
		panic(err)
	}

	return &s
}

// Returns is partition key and clustering columns matches the scheduled error
func (e *scheduleStorageError) match(pKey, cCols []byte) bool {
	return ((len(e.pKey) == 0) || bytes.Equal(e.pKey, pKey)) &&
		((len(e.cCols) == 0) || bytes.Equal(e.cCols, cCols))
}

// Clear all scheduled errors
func (s *TestMemStorage) Reset() {
	s.get = scheduleStorageError{}
	s.put = scheduleStorageError{}
	s.damage = scheduleStorageDamage{}
}

// Schedule Get() to return error
func (s *TestMemStorage) ScheduleGetError(err error, pKey, cCols []byte) {
	s.get.err = err
	s.get.pKey = make([]byte, len(pKey))
	copy(s.get.pKey, pKey)
	s.get.cCols = make([]byte, len(cCols))
	copy(s.get.cCols, cCols)
}

// Schedule Get() to return damaged data
func (s *TestMemStorage) ScheduleGetDamage(dam damagedStorageFunc, pKey, cCols []byte) {
	s.damage.dam = dam
	s.damage.pKey = make([]byte, len(pKey))
	copy(s.damage.pKey, pKey)
	s.damage.cCols = make([]byte, len(cCols))
	copy(s.damage.cCols, cCols)
}

// Schedule Put() to return error
func (s *TestMemStorage) SchedulePutError(err error, pKey, cCols []byte) {
	s.put.err = err
	s.put.pKey = make([]byte, len(pKey))
	copy(s.put.pKey, pKey)
	s.put.cCols = make([]byte, len(cCols))
	copy(s.put.cCols, cCols)
}

func (s *TestMemStorage) Get(pKey []byte, cCols []byte, data *[]byte) (ok bool, err error) {
	if s.get.err != nil {
		if s.get.match(pKey, cCols) {
			err = s.get.err
			s.get.err = nil
			return false, err
		}
	}

	ok, err = s.storage.Get(pKey, cCols, data)

	if ok && (s.damage.dam != nil) {
		if s.damage.match(pKey, cCols) {
			s.damage.dam(data)
			s.damage.dam = nil
			return ok, err
		}
	}

	return ok, err
}

func (s *TestMemStorage) GetBatch(pKey []byte, items []istorage.GetBatchItem) (err error) {
	if s.get.err != nil {
		for _, item := range items {
			if s.get.match(pKey, item.CCols) {
				err = s.get.err
				s.get.err = nil
				return err
			}
		}
	}

	err = s.storage.GetBatch(pKey, items)

	if s.damage.dam != nil {
		for i := 0; i < len(items); i++ {
			if s.damage.match(pKey, items[i].CCols) {
				if items[i].Ok {
					s.damage.dam(items[i].Data)
					s.damage.dam = nil
				}
			}
		}
	}

	return err
}

func (s *TestMemStorage) Put(pKey []byte, cCols []byte, value []byte) (err error) {
	if s.put.err != nil {
		if s.put.match(pKey, cCols) {
			err = s.put.err
			s.put.err = nil
			return err
		}
	}
	return s.storage.Put(pKey, cCols, value)
}

func (s *TestMemStorage) PutBatch(items []istorage.BatchItem) (err error) {
	for _, p := range items {
		if err = s.Put(p.PKey, p.CCols, p.Value); err != nil {
			return err
		}
	}
	return nil
}

func (s *TestMemStorage) Read(ctx context.Context, pKey []byte, startCCols, finishCCols []byte, cb istorage.ReadCallback) (err error) {
	cbWrap := func(cCols []byte, data []byte) (err error) {
		if s.get.err != nil {
			if s.get.match(pKey, cCols) {
				err = s.get.err
				s.get.err = nil
				return err
			}
		}

		if s.damage.dam != nil {
			if s.damage.match(pKey, cCols) {
				s.damage.dam(&data)
				s.damage.dam = nil
			}
		}

		return cb(cCols, data)
	}

	return s.storage.Read(ctx, pKey, startCCols, finishCCols, cbWrap)
}
