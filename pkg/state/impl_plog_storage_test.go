/*
 * Copyright (c) 2022-present unTill Pro, Ltd.
 */

package state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/istructs"
)

func TestPLogStorage_Read(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		require := require.New(t)
		touched := false
		events := &mockEvents{}
		events.On("ReadPLog", context.Background(), istructs.PartitionID(1), istructs.FirstOffset, 1, mock.AnythingOfType("istructs.PLogEventsReaderCallback")).
			Return(nil).
			Run(func(args mock.Arguments) {
				require.NoError(args.Get(4).(istructs.PLogEventsReaderCallback)(istructs.FirstOffset, nil))
			})
		appStructs := &mockAppStructs{}
		appStructs.On("AppDef").Return(&nilAppDef{})
		appStructs.On("Events").Return(events)
		appStructs.On("Records").Return(&nilRecords{})
		appStructs.On("ViewRecords").Return(&nilViewRecords{})
		s := ProvideQueryProcessorStateFactory()(context.Background(), appStructs, SimplePartitionIDFunc(istructs.PartitionID(1)), nil, nil, nil, nil)
		kb, err := s.KeyBuilder(PLogStorage, appdef.NullQName)
		require.NoError(err)
		kb.PutInt64(Field_Offset, 1)
		kb.PutInt64(Field_Count, 1)

		err = s.Read(kb, func(key istructs.IKey, _ istructs.IStateValue) (err error) {
			touched = true
			require.Equal(int64(1), key.AsInt64(Field_Offset))
			return
		})
		require.NoError(err)

		require.True(touched)
	})
	t.Run("Should return error on read plog", func(t *testing.T) {
		require := require.New(t)
		events := &mockEvents{}
		events.On("ReadPLog", context.Background(), istructs.PartitionID(1), istructs.FirstOffset, 1, mock.AnythingOfType("istructs.PLogEventsReaderCallback")).Return(errTest)
		appStructs := &mockAppStructs{}
		appStructs.On("AppDef").Return(&nilAppDef{})
		appStructs.On("Events").Return(events)
		appStructs.On("Records").Return(&nilRecords{})
		appStructs.On("ViewRecords").Return(&nilViewRecords{})
		s := ProvideQueryProcessorStateFactory()(context.Background(), appStructs, SimplePartitionIDFunc(istructs.PartitionID(1)), nil, nil, nil, nil)
		k, err := s.KeyBuilder(PLogStorage, appdef.NullQName)
		require.NoError(err)
		k.PutInt64(Field_Offset, 1)
		k.PutInt64(Field_Count, 1)

		err = s.Read(k, func(istructs.IKey, istructs.IStateValue) error { return nil })

		require.ErrorIs(err, errTest)
	})
}
func TestPLogStorage_GetBatch(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		require := require.New(t)
		events := &mockEvents{}
		events.On("ReadPLog", context.Background(), istructs.PartitionID(1), istructs.FirstOffset, 1, mock.AnythingOfType("istructs.PLogEventsReaderCallback")).
			Return(nil).
			Run(func(args mock.Arguments) {
				cb := args.Get(4).(istructs.PLogEventsReaderCallback)
				require.NoError(cb(istructs.FirstOffset, nil))
				require.NoError(cb(istructs.Offset(2), nil))
				require.NoError(cb(istructs.Offset(3), nil))
			})
		appStructs := &mockAppStructs{}
		appStructs.On("AppDef").Return(&nilAppDef{})
		appStructs.On("Events").Return(events)
		appStructs.On("Records").Return(&nilRecords{})
		appStructs.On("ViewRecords").Return(&nilViewRecords{})
		s := ProvideCommandProcessorStateFactory()(context.Background(), func() istructs.IAppStructs { return appStructs }, SimplePartitionIDFunc(istructs.PartitionID(1)), nil, nil, nil, nil, nil, 0)
		kb, err := s.KeyBuilder(PLogStorage, appdef.NullQName)
		require.NoError(err)
		kb.PutInt64(Field_Offset, 1)
		kb.PutInt64(Field_Count, 1)

		sv, ok, err := s.CanExist(kb)
		require.NoError(err)

		require.True(ok)
		require.Equal(int64(1), sv.AsInt64(Field_Offset))
	})
	t.Run("Should return error when error occurred on read plog", func(t *testing.T) {
		require := require.New(t)
		events := &mockEvents{}
		events.
			On("ReadPLog", context.Background(), istructs.PartitionID(1), istructs.FirstOffset, 1, mock.AnythingOfType("istructs.PLogEventsReaderCallback")).
			Return(nil).
			Run(func(args mock.Arguments) {
				require.NoError(args.Get(4).(istructs.PLogEventsReaderCallback)(istructs.FirstOffset, nil))
			}).
			On("ReadPLog", context.Background(), istructs.PartitionID(1), istructs.Offset(2), 1, mock.AnythingOfType("istructs.PLogEventsReaderCallback")).
			Return(errTest)
		appStructs := &mockAppStructs{}
		appStructs.On("AppDef").Return(&nilAppDef{})
		appStructs.On("Events").Return(events)
		appStructs.On("Records").Return(&nilRecords{})
		appStructs.On("ViewRecords").Return(&nilViewRecords{})
		s := ProvideCommandProcessorStateFactory()(context.Background(), func() istructs.IAppStructs { return appStructs }, SimplePartitionIDFunc(istructs.PartitionID(1)), nil, nil, nil, nil, nil, 0)
		kb1, err := s.KeyBuilder(PLogStorage, appdef.NullQName)
		require.NoError(err)
		kb1.PutInt64(Field_Offset, 1)
		kb1.PutInt64(Field_Count, 1)
		kb2, err := s.KeyBuilder(PLogStorage, appdef.NullQName)
		require.NoError(err)
		kb2.PutInt64(Field_Offset, 2)
		kb2.PutInt64(Field_Count, 1)

		err = s.CanExistAll([]istructs.IStateKeyBuilder{kb1, kb2}, nil)

		require.ErrorIs(err, errTest)
	})
}
func TestPLogStorage_ToJSON(t *testing.T) {
	s := &pLogStorage{appDefFunc: func() appdef.IAppDef { return nil }}
	require := require.New(t)
	eventError := &mockEventError{}
	eventError.
		On("ErrStr").Return("Error string").
		On("QNameFromParams").Return(istructs.QNameForError).
		On("ValidEvent").Return(false).
		On("OriginalEventBytes").Return([]byte("Love bites"))
	event := &mockPLogEvent{}
	event.
		On("QName").Return(testEvent).
		On("ArgumentObject").Return(istructs.NewNullObject()).
		On("CUDs", mock.Anything).Return(nil).
		On("RegisteredAt").Return(istructs.UnixMilli(1662972220000)).
		On("Synced").Return(true).
		On("DeviceID").Return(istructs.ConnectedDeviceID(7)).
		On("SyncedAt").Return(istructs.UnixMilli(1662972220001)).
		On("Workspace").Return(istructs.WSID(1001)).
		On("WLogOffset").Return(istructs.Offset(56)).
		On("Error").Return(eventError)

	sv := &pLogStorageValue{
		event:      event,
		offset:     123,
		toJSONFunc: s.toJSON,
	}

	json, err := sv.ToJSON()
	require.NoError(err)

	require.JSONEq(`
						{
						  "ArgumentObject": {},
						  "CUDs":[],
						  "DeviceID": 7,
						  "Error": {
							"ErrStr": "Error string",
							"OriginalEventBytes": "TG92ZSBiaXRlcw==",
							"QNameFromParams": "sys.Error",
							"ValidEvent": false
						  },
						  "Offset": 123,
						  "QName": "test.event",
						  "RegisteredAt": 1662972220000,
						  "Synced": true,
						  "SyncedAt": 1662972220001,
						  "WLogOffset": 56,
						  "Workspace": 1001
						}`, json)
}
