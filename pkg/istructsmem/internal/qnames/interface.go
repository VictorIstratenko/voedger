/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package qnames

import "github.com/voedger/voedger/pkg/appdef"

// Identificator for QNames
type QNameID uint16

// QNames system view
//
//	Use GetID() to obtain QName ID.
//	Use GetQName() to obtain QName name by its ID.
//	Use Prepare() to load QNames IDs from storage.
type QNames struct {
	qNames  map[appdef.QName]QNameID
	ids     map[QNameID]appdef.QName
	lastID  QNameID
	changes uint
}
