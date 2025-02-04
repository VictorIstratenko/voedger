/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 */

package istructsmem

import (
	"errors"
	"fmt"
)

var ErrorInvalidVersion = errors.New("invalid version")

var ErrorEventNotValid = errors.New("event is not valid")

var ErrNameMissed = errors.New("name is empty")

var ErrNameNotFound = errors.New("name not found")

var ErrInvalidName = errors.New("name not valid")

var ErrIDNotFound = errors.New("ID not found")

var ErrorRecordIDNotFound = fmt.Errorf("recordID cannot be found: %w", ErrIDNotFound)

var ErrRecordNotFound = errors.New("record cannot be found")

var ErrMinOccursViolation = errors.New("minimum occurs violated")

var ErrMaxOccursViolation = errors.New("maximum occurs violated")

var ErrFieldIsEmpty = errors.New("field is empty")

var ErrInvalidVerificationKind = errors.New("invalid verification kind")

var ErrCUDsMissed = errors.New("CUDs are missed")

var ErrRawRecordIDExpected = errors.New("raw record ID expected")

var ErrRecordIDUniqueViolation = errors.New("record ID duplicates")

var ErrWrongRecordID = errors.New("wrong record ID")

var ErrUnableToUpdateSystemField = errors.New("unable to update system field")

var ErrWrongDefinition = errors.New("wrong definition")

var ErrUnexpectedDefKind = errors.New("unexpected definition kind")

var ErrUnknownCodec = errors.New("unknown codec")

var ErrSingletonIDsExceeds = errors.New("the maximum number of singleton document identifiers has been exceeded")

var ErrMaxGetBatchRecordCountExceeds = errors.New("the maximum count of records to batch is exceeded")

var ErrEmptySetOfKeyFields = errors.New("empty set of key fields")

var ErrWrongFieldType = errors.New("wrong field type")

var ErrKeyMustHaveNotMoreThanOneVarSizeField = errors.New("key must have not more than one variable size field")

var ErrKeyFieldMustBeRequired = errors.New("key field must be required")

var ErrUnknownDefinitionQName = errors.New("unknown definition QName")

var ErrDefKindMayNotHaveUniques = errors.New("definition kind may not have uniques")

var ErrUnknownKeyField = errors.New("unknown key field")

var ErrUniquesHaveSameFields = errors.New("uniques have same fields")

var ErrKeyFieldIsUsedMoreThanOnce = errors.New("key field is used more than once")

var ErrDefChanged = errors.New("definition has been changed")

var ErrReferentialIntegrityViolation = errors.New("referencial integrity violation")

const errFieldNotFoundWrap = "%s-type field «%s» is not found in definition «%v»: %w" // int32-type field «myField» is not found …

const errFieldValueTypeMismatchWrap = "value type «%s» is not applicable for %s-type field «%s»: %w" // value type «float64» is not applicable for int32-type field «myField»: …

const errFieldMustBeVerificated = "field «%s» must be verificated, token expected, but value «%T» passed: %w"

const errFieldConvertErrorWrap = "field «%s» value type «%T» can not to be converted to «%s»: %w"

const errCantGetFieldQNameIDWrap = "QName field «%s» can not get ID for value «%v»: %w"

const errDefNotFoundWrap = "definition «%v» not found: %w"

// ValidateError: an interface for describing errors that occurred during validation
//   - methods:
//     — Code(): returns error code, see ECode_××× constants
type ValidateError interface {
	error
	Code() int
}
