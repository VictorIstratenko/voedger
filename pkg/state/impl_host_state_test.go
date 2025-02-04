/*
 * Copyright (c) 2022-present unTill Pro, Ltd.
 */

package state

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/voedger/voedger/pkg/appdef"
	amock "github.com/voedger/voedger/pkg/appdef/mock"
	"github.com/voedger/voedger/pkg/istructs"
)

func TestHostState_BasicUsage(t *testing.T) {
	require := require.New(t)

	factory := ProvideQueryProcessorStateFactory()
	hostState := factory(context.Background(), mockedHostStateStructs(), nil, SimpleWSIDFunc(istructs.WSID(1)), nil, nil, nil)

	// Declare simple extension
	extension := func(state istructs.IState) {
		//Create key
		key, err := state.KeyBuilder(ViewRecordsStorage, testViewRecordQName1)
		require.NoError(err)
		key.PutString("pkFld", "pkVal")

		// Call to storage
		require.NoError(state.MustNotExist(key))
	}

	// Run extension
	extension(hostState)

	require.NoError(hostState.ValidateIntents())
	require.NoError(hostState.ApplyIntents())
}

func mockedHostStateStructs() istructs.IAppStructs {
	mv := &mockValue{}
	mv.
		On("AsInt64", "vFld").Return(int64(10)).
		On("AsInt64", ColOffset).Return(int64(45))
	mvb1 := &mockValueBuilder{}
	mvb1.
		On("PutInt64", "vFld", int64(10)).
		On("PutInt64", ColOffset, int64(45)).
		On("Build").Return(mv)
	mvb2 := &mockValueBuilder{}
	mvb2.
		On("PutInt64", "vFld", int64(10)).Once().
		On("PutInt64", ColOffset, int64(45)).Once().
		On("PutInt64", "vFld", int64(17)).Once().
		On("PutInt64", ColOffset, int64(46)).Once()
	mkb := &mockKeyBuilder{}
	mkb.
		On("PutString", "pkFld", "pkVal")
	viewRecords := &mockViewRecords{}
	viewRecords.
		On("KeyBuilder", testViewRecordQName1).Return(mkb).
		On("NewValueBuilder", testViewRecordQName1).Return(mvb1).Once().
		On("NewValueBuilder", testViewRecordQName1).Return(mvb2).Once().
		On("GetBatch", istructs.WSID(1), mock.AnythingOfType("[]istructs.ViewRecordGetBatchItem")).
		Return(nil).
		Run(func(args mock.Arguments) {
			value := &mockValue{}
			value.On("AsString", "vk").Return("value")
			args.Get(1).([]istructs.ViewRecordGetBatchItem)[0].Value = value
		}).
		On("PutBatch", istructs.WSID(1), mock.AnythingOfType("[]istructs.ViewKV")).Return(nil)

	view := amock.NewView(testViewRecordQName1)
	view.
		AddPartField("pkFld", appdef.DataKind_string). //??? string in partition key
		AddClustColumn("ccFld", appdef.DataKind_string).
		AddValueField("vFld", appdef.DataKind_int64, false).
		AddValueField(ColOffset, appdef.DataKind_int64, false)

	appDef := amock.NewAppDef()
	appDef.AddView(view)

	appStructs := &mockAppStructs{}
	appStructs.
		On("AppDef").Return(appDef).
		On("ViewRecords").Return(viewRecords).
		On("Events").Return(&nilEvents{}).
		On("Records").Return(&nilRecords{})
	return appStructs
}
func TestHostState_KeyBuilder_Should_return_unknown_storage_ID_error(t *testing.T) {
	require := require.New(t)
	s := hostStateForTest(&mockStorage{})

	_, err := s.KeyBuilder(appdef.NullQName, appdef.NullQName)

	require.ErrorIs(err, ErrUnknownStorage)
}
func TestHostState_CanExist(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = &mockStateValue{}
			})
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, ok, err := s.CanExist(k)
		require.NoError(err)

		require.True(ok)
	})
	t.Run("Should return error when error occurred", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).Return(errTest)
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, _, err = s.CanExist(k)

		require.ErrorIs(err, errTest)
	})
	t.Run("Should return get batch not supported by storage error", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, _ := emptyHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, _, err = s.CanExist(kb)

		require.ErrorIs(err, ErrGetBatchNotSupportedByStorage)
	})
}
func TestHostState_CanExistAll(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		require := require.New(t)
		times := 0
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = &mockStateValue{}
			})
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.CanExistAll([]istructs.IStateKeyBuilder{k}, func(key istructs.IKeyBuilder, value istructs.IStateValue, ok bool) (err error) {
			times++
			require.Equal(k, key)
			require.True(ok)
			return
		})
		require.NoError(err)

		require.Equal(1, times)
	})
	t.Run("Should return error when error occurred", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).Return(errTest)
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.CanExistAll([]istructs.IStateKeyBuilder{k}, nil)

		require.ErrorIs(err, errTest)
	})
	t.Run("Should return get batch not supported by storage error", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, _ := emptyHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.CanExistAll([]istructs.IStateKeyBuilder{kb}, nil)

		require.ErrorIs(err, ErrGetBatchNotSupportedByStorage)
	})
}
func TestHostState_MustExist(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = &mockStateValue{}
			})
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, err = s.MustExist(k)

		require.NoError(err)
	})
	t.Run("Should return error when entity not exists", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = nil
			})
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, err = s.MustExist(k)

		require.ErrorIs(err, ErrNotExists)
	})
	t.Run("Should return error when error occurred on get batch", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).Return(errTest)
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, err = s.MustExist(k)

		require.ErrorIs(err, errTest)
	})
}
func TestHostState_MustExistAll(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = &mockStateValue{}
				args.Get(0).([]GetBatchItem)[1].value = &mockStateValue{}
			})
		s := hostStateForTest(ms)
		k1, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)
		k2, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)
		kk := make([]istructs.IKeyBuilder, 0, 2)

		err = s.MustExistAll([]istructs.IStateKeyBuilder{k1, k2}, func(key istructs.IKeyBuilder, value istructs.IStateValue, ok bool) (err error) {
			kk = append(kk, key)
			require.True(ok)
			return
		})
		require.NoError(err)

		require.Equal(k1, kk[0])
		require.Equal(k1, kk[1])
	})
	t.Run("Should return error on get batch", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).Return(errTest)
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustExistAll([]istructs.IStateKeyBuilder{k}, nil)

		require.ErrorIs(err, errTest)
	})
	t.Run("Should return error when entity not exists", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = &mockStateValue{}
				args.Get(0).([]GetBatchItem)[1].value = nil
			})
		s := hostStateForTest(ms)
		k1, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)
		k2, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustExistAll([]istructs.IStateKeyBuilder{k1, k2}, nil)

		require.ErrorIs(err, ErrNotExists)
	})
	t.Run("Should return get batch not supported by storage error", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, _ := emptyHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustExistAll([]istructs.IStateKeyBuilder{kb}, nil)

		require.ErrorIs(err, ErrGetBatchNotSupportedByStorage)
	})
}
func TestHostState_MustNotExist(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = nil
			})
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustNotExist(k)

		require.NoError(err)
	})
	t.Run("Should return error when entity exists", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = &mockStateValue{}
			})
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustNotExist(k)

		require.ErrorIs(err, ErrExists)
	})
	t.Run("Should return error when error occurred on get batch", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).Return(errTest)
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustNotExist(k)

		require.ErrorIs(err, errTest)
	})
}
func TestHostState_MustNotExistAll(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = nil
				args.Get(0).([]GetBatchItem)[1].value = nil
			})
		s := hostStateForTest(ms)
		k1, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)
		k2, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustNotExistAll([]istructs.IStateKeyBuilder{k1, k2})

		require.NoError(err)
	})
	t.Run("Should return error on get batch", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).Return(errTest)
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustNotExistAll([]istructs.IStateKeyBuilder{k})

		require.ErrorIs(err, errTest)
	})
	t.Run("Should return error when entity exists", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("GetBatch", mock.AnythingOfType("[]state.GetBatchItem")).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).([]GetBatchItem)[0].value = nil
				args.Get(0).([]GetBatchItem)[1].value = &mockStateValue{}
			})
		s := hostStateForTest(ms)
		k1, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)
		k2, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustNotExistAll([]istructs.IStateKeyBuilder{k1, k2})

		require.ErrorIs(err, ErrExists)
	})
	t.Run("Should return get batch not supported by storage error", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, _ := emptyHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.MustNotExistAll([]istructs.IStateKeyBuilder{kb})

		require.ErrorIs(err, ErrGetBatchNotSupportedByStorage)
	})
}
func TestHostState_Read(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("Read", mock.Anything, mock.AnythingOfType("istructs.ValueCallback")).Return(nil)
		s := hostStateForTest(ms)
		k, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(t, err)

		require.NoError(t, s.Read(k, nil))

		ms.AssertExpectations(t)
	})
	t.Run("Should return read not supported by storage error", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, _ := emptyHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		err = s.Read(kb, nil)

		require.ErrorIs(err, ErrReadNotSupportedByStorage)
	})
}
func TestHostState_NewValue(t *testing.T) {
	t.Run("Should return error when intents limit exceeded", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, i := limitedIntentsHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, err = i.NewValue(kb)

		require.ErrorIs(err, ErrIntentsLimitExceeded)
	})
	t.Run("Should return insert not supported by storage error", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, i := emptyHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, err = i.NewValue(kb)

		require.ErrorIs(err, ErrInsertNotSupportedByStorage)
	})
}
func TestHostState_UpdateValue(t *testing.T) {
	t.Run("Should return error when intents limit exceeded", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, i := limitedIntentsHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, err = i.UpdateValue(kb, nil)

		require.ErrorIs(err, ErrIntentsLimitExceeded)
	})
	t.Run("Should return update not supported by storage error", func(t *testing.T) {
		require := require.New(t)
		ms := &mockStorage{}
		ms.On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName))
		s, i := emptyHostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(err)

		_, err = i.UpdateValue(kb, nil)

		require.ErrorIs(err, ErrUpdateNotSupportedByStorage)
	})
}
func TestHostState_ValidateIntents(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("ProvideValueBuilder", mock.Anything, mock.Anything).Return(&viewRecordsValueBuilder{}).
			On("Validate", mock.Anything).Return(nil)
		s := hostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(t, err)
		_, err = s.NewValue(kb)
		require.NoError(t, err)

		err = s.ValidateIntents()

		require.NoError(t, err)
	})
	t.Run("Should return immediately when intents are empty", func(t *testing.T) {
		ms := &mockStorage{}
		s := hostStateForTest(&mockStorage{})

		require.NoError(t, s.ValidateIntents())

		ms.AssertNotCalled(t, "Validate", mock.Anything)
	})
	t.Run("Should return validation error", func(t *testing.T) {
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("ProvideValueBuilder", mock.Anything, mock.Anything).Return(&viewRecordsValueBuilder{}).
			On("Validate", mock.Anything).Return(errTest)
		s := hostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(t, err)
		_, err = s.NewValue(kb)
		require.NoError(t, err)

		err = s.ValidateIntents()

		require.ErrorIs(t, err, errTest)
	})
}
func TestHostState_ApplyIntents(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("ProvideValueBuilder", mock.Anything, mock.Anything).Return(&viewRecordsValueBuilder{}).
			On("ApplyBatch", mock.Anything).Return(nil)
		s := hostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(t, err)
		_, err = s.NewValue(kb)
		require.NoError(t, err)

		require.NoError(t, s.ApplyIntents())

		ms.AssertExpectations(t)
	})
	t.Run("Should return apply batch error", func(t *testing.T) {
		ms := &mockStorage{}
		ms.
			On("NewKeyBuilder", appdef.NullQName, nil).Return(newKeyBuilder(testStorage, appdef.NullQName)).
			On("ProvideValueBuilder", mock.Anything, mock.Anything).Return(&viewRecordsValueBuilder{}).
			On("ApplyBatch", mock.Anything).Return(errTest)
		s := hostStateForTest(ms)
		kb, err := s.KeyBuilder(testStorage, appdef.NullQName)
		require.NoError(t, err)
		_, err = s.NewValue(kb)
		require.NoError(t, err)

		err = s.ApplyIntents()

		require.ErrorIs(t, err, errTest)
	})
}
func hostStateForTest(s IStateStorage) IHostState {
	hs := newHostState("ForTest", 10)
	hs.addStorage(testStorage, s, S_GET_BATCH|S_READ|S_INSERT|S_UPDATE)
	return hs
}
func emptyHostStateForTest(s IStateStorage) (istructs.IState, istructs.IIntents) {
	bs := ProvideQueryProcessorStateFactory()(context.Background(), &nilAppStructs{}, nil, nil, nil, nil, nil).(*hostState)
	bs.addStorage(testStorage, s, math.MinInt)
	return bs, bs
}
func limitedIntentsHostStateForTest(s IStateStorage) (istructs.IState, istructs.IIntents) {
	hs := newHostState("LimitedIntentsForTest", 0)
	hs.addStorage(testStorage, s, S_GET_BATCH|S_READ|S_INSERT|S_UPDATE)
	return hs, hs
}
