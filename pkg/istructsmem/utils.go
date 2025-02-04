/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package istructsmem

import (
	"encoding/binary"
	"fmt"

	"github.com/untillpro/dynobuffers"
	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/irates"
	"github.com/voedger/voedger/pkg/istorage"
	"github.com/voedger/voedger/pkg/istorageimpl"
	"github.com/voedger/voedger/pkg/istructs"
	payloads "github.com/voedger/voedger/pkg/itokens-payloads"
	"github.com/voedger/voedger/pkg/itokensjwt"
)

var NullAppConfig = newAppConfig(istructs.AppQName_null, appdef.New())

var (
	nullDynoBuffer = dynobuffers.NewBuffer(dynobuffers.NewScheme())
	// not a func -> golang tokensjwt.TimeFunc will be initialized on process init forever
	testTokensFactory    = func() payloads.IAppTokensFactory { return payloads.TestAppTokensFactory(itokensjwt.TestTokensJWT()) }
	simpleStorageProvder = func() istorage.IAppStorageProvider {
		asf := istorage.ProvideMem()
		return istorageimpl.Provide(asf)
	}
)

// crackID splits ID to two-parts key — partition key (hi) and clustering columns (lo)
func crackID(id uint64) (hi uint64, low uint16) {
	return uint64(id >> partitionBits), uint16(id) & lowMask
}

// crackRecordID splits record ID to two-parts key — partition key (hi) and clustering columns (lo)
func crackRecordID(id istructs.RecordID) (hi uint64, low uint16) {
	return crackID(uint64(id))
}

// crackLogOffset splits log offset to two-parts key — partition key (hi) and clustering columns (lo)
func crackLogOffset(ofs istructs.Offset) (hi uint64, low uint16) {
	return crackID(uint64(ofs))
}

// uncrackLogOffset calculate log offset from two-parts key — partition key (hi) and clustering columns (lo)
func uncrackLogOffset(hi uint64, low uint16) istructs.Offset {
	return istructs.Offset(hi<<partitionBits | uint64(low))
}

const uint64bits, uint16bits = 8, 2

// splitRecordID splits record ID to two-parts key — partition key and clustering columns
func splitRecordID(id istructs.RecordID) (pk, cc []byte) {
	hi, lo := crackRecordID(id)
	pkBuf := make([]byte, uint64bits)
	binary.BigEndian.PutUint64(pkBuf, hi)
	ccBuf := make([]byte, uint16bits)
	binary.BigEndian.PutUint16(ccBuf, lo)
	return pkBuf, ccBuf
}

// splitLogOffset splits offset to two-parts key — partition key and clustering columns
func splitLogOffset(offset istructs.Offset) (pk, cc []byte) {
	hi, lo := crackLogOffset(offset)
	pkBuf := make([]byte, uint64bits)
	binary.BigEndian.PutUint64(pkBuf, hi)
	ccBuf := make([]byte, uint16bits)
	binary.BigEndian.PutUint16(ccBuf, lo)
	return pkBuf, ccBuf
}

// calcLogOffset calculate log offset from two-parts key — partition key and clustering columns
func calcLogOffset(pk, cc []byte) istructs.Offset {
	hi := binary.BigEndian.Uint64(pk)
	low := binary.BigEndian.Uint16(cc)
	return uncrackLogOffset(hi, low)
}

// used in tests only
func IBucketsFromIAppStructs(as istructs.IAppStructs) irates.IBuckets {
	// appStructs implementation has method Buckets()
	return as.(interface{ Buckets() irates.IBuckets }).Buckets()
}

func FillElementFromJSON(data map[string]interface{}, def appdef.IDef, b istructs.IElementBuilder) error {
	for fieldName, fieldValue := range data {
		switch fv := fieldValue.(type) {
		case float64:
			b.PutNumber(fieldName, fv)
		case string:
			b.PutChars(fieldName, fv)
		case bool:
			b.PutBool(fieldName, fv)
		case []interface{}:
			// e.g. TestBasicUsage_Dashboard(), "order_item": [<2 elements>]
			containerName := fieldName
			containerDef := def.ContainerDef(containerName)
			if containerDef.Kind() == appdef.DefKind_null {
				return fmt.Errorf("container with name %s is not found", containerName)
			}
			for i, intf := range fv {
				objContainerElem, ok := intf.(map[string]interface{})
				if !ok {
					return fmt.Errorf("element #%d of %s is not an object", i, fieldName)
				}
				containerElemBuilder := b.ElementBuilder(fieldName)
				if err := FillElementFromJSON(objContainerElem, containerDef, containerElemBuilder); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func NewIObjectBuilder(cfg *AppConfigType, qName appdef.QName) istructs.IObjectBuilder {
	obj := newObject(cfg, qName)
	return &obj
}

func CheckRefIntegrity(obj istructs.IRowReader, appStructs istructs.IAppStructs, wsid istructs.WSID) (err error) {
	appDef := appStructs.AppDef()
	def := appDef.Def(obj.AsQName(appdef.SystemField_QName))
	def.Fields(
		func(f appdef.IField) {
			if err != nil || f.DataKind() != appdef.DataKind_RecordID {
				return
			}
			recID := obj.AsRecordID(f.Name())
			if recID.IsRaw() || recID == istructs.NullRecordID {
				return
			}
			var rec istructs.IRecord
			rec, err = appStructs.Records().Get(wsid, true, recID)
			if err != nil {
				return
			}
			if rec.QName() == appdef.NullQName {
				err = fmt.Errorf("%w: record ID %d referenced by %s.%s does not exist", ErrReferentialIntegrityViolation, recID,
					obj.AsQName(appdef.SystemField_QName), f.Name())
			}
		})
	return err
}
