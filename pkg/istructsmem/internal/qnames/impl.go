/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package qnames

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/istorage"
	"github.com/voedger/voedger/pkg/istructs"
	"github.com/voedger/voedger/pkg/istructsmem/internal/consts"
	"github.com/voedger/voedger/pkg/istructsmem/internal/utils"
	"github.com/voedger/voedger/pkg/istructsmem/internal/vers"
)

func newQNames() *QNames {
	return &QNames{
		qNames: make(map[appdef.QName]QNameID),
		ids:    make(map[QNameID]appdef.QName),
		lastID: QNameIDSysLast,
	}
}

// Returns ID for specified QName
func (names *QNames) GetID(qName appdef.QName) (QNameID, error) {
	if id, ok := names.qNames[qName]; ok {
		return id, nil
	}
	return 0, fmt.Errorf("unknown QName «%v»: %w", qName, ErrNameNotFound)
}

// Retrieve QName for specified ID
func (names *QNames) GetQName(id QNameID) (qName appdef.QName, err error) {
	qName, ok := names.ids[id]
	if ok {
		return qName, nil
	}

	return appdef.NullQName, fmt.Errorf("unknown QName ID «%v»: %w", id, ErrIDNotFound)
}

// Reads all application QNames from storage, add all system and application QNames and write result to storage if some changes. Must be called at application starts
func (names *QNames) Prepare(storage istorage.IAppStorage, versions *vers.Versions, appDef appdef.IAppDef, resources istructs.IResources) error {
	if err := names.load(storage, versions); err != nil {
		return err
	}

	if err := names.collectAllQNames(appDef, resources); err != nil {
		return err
	}

	if names.changes > 0 {
		if err := names.store(storage, versions); err != nil {
			return err
		}
	}

	return nil
}

// Collect all system and application QName IDs
func (names *QNames) collectAllQNames(appDef appdef.IAppDef, r istructs.IResources) (err error) {

	// system QNames
	names.
		collectSysQName(appdef.NullQName, NullQNameID).
		collectSysQName(istructs.QNameForError, QNameIDForError).
		collectSysQName(istructs.QNameCommandCUD, QNameIDCommandCUD)

	if appDef != nil {
		appDef.Defs(
			func(d appdef.IDef) {
				err = errors.Join(err,
					names.collectAppQName(d.QName()))
			})
	}

	if r != nil {
		r.Resources(
			func(q appdef.QName) {
				err = errors.Join(err,
					names.collectAppQName(q))
			})
	}

	return err
}

// Checks is exists ID for application QName in cache. If not then adds it with new ID
func (names *QNames) collectAppQName(qName appdef.QName) error {
	if _, ok := names.qNames[qName]; ok {
		return nil // already known QName
	}

	for id := names.lastID + 1; id < MaxAvailableQNameID; id++ {
		if _, ok := names.ids[id]; !ok {
			names.qNames[qName] = id
			names.ids[id] = qName
			names.lastID = id
			names.changes++
			return nil
		}
	}

	return ErrQNameIDsExceeds
}

// Adds system QName to cache
func (names *QNames) collectSysQName(qName appdef.QName, id QNameID) *QNames {
	names.qNames[qName] = id
	names.ids[id] = qName
	return names
}

// loads all stored QNames from storage
func (names *QNames) load(storage istorage.IAppStorage, versions *vers.Versions) (err error) {

	ver := versions.Get(vers.SysQNamesVersion)
	switch ver {
	case vers.UnknownVersion: // no sys.QName storage exists
		return nil
	case ver01:
		return names.load01(storage)
	}

	return fmt.Errorf("unknown version of system QNames view (%v): %w", ver, vers.ErrorInvalidVersion)
}

// loads all stored QNames from storage version ver01
func (names *QNames) load01(storage istorage.IAppStorage) error {

	readQName := func(cCols, value []byte) error {
		qName, err := appdef.ParseQName(string(cCols))
		if err != nil {
			return err
		}
		id := QNameID(binary.BigEndian.Uint16(value))
		if id == NullQNameID {
			return nil // deleted QName
		}

		if id <= QNameIDSysLast {
			return fmt.Errorf("unexpected ID (%v) is readed from system QNames view: %w", id, ErrWrongQNameID)
		}

		names.qNames[qName] = id
		names.ids[id] = qName

		if names.lastID < id {
			names.lastID = id
		}

		return nil
	}
	pKey := utils.ToBytes(consts.SysView_QNames, ver01)
	return storage.Read(context.Background(), pKey, nil, nil, readQName)
}

// Stores all known QNames to storage
func (names *QNames) store(storage istorage.IAppStorage, versions *vers.Versions) (err error) {
	pKey := utils.ToBytes(consts.SysView_QNames, ver01)

	batch := make([]istorage.BatchItem, 0)
	for qName, id := range names.qNames {
		if (id > QNameIDSysLast) ||
			(qName != appdef.NullQName) && (id == NullQNameID) { // deleted QName
			item := istorage.BatchItem{
				PKey:  pKey,
				CCols: []byte(qName.String()),
				Value: utils.ToBytes(id),
			}
			batch = append(batch, item)
		}
	}

	if err = storage.PutBatch(batch); err != nil {
		return fmt.Errorf("error store application QName IDs to storage: %w", err)
	}

	if ver := versions.Get(vers.SysQNamesVersion); ver != lastestVersion {
		if err = versions.Put(vers.SysQNamesVersion, lastestVersion); err != nil {
			return fmt.Errorf("error store system QNames view version: %w", err)
		}
	}

	names.changes = 0
	return nil
}
