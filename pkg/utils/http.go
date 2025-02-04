/*
 * Copyright (c) 2020-present unTill Pro, Ltd.
 * @author Denis Gribanov
 */

package coreutils

import (
	"fmt"
	"net/http"

	ibus "github.com/untillpro/airs-ibus"
)

func NewHTTPErrorf(httpStatus int, args ...interface{}) SysError {
	return SysError{
		HTTPStatus: httpStatus,
		Message:    fmt.Sprint(args...),
	}
}

func NewHTTPError(httpStatus int, err error) SysError {
	return NewHTTPErrorf(httpStatus, err.Error())
}

func ReplyErrf(bus ibus.IBus, sender interface{}, status int, args ...interface{}) {
	ReplyErrDef(bus, sender, NewHTTPErrorf(status, args...), http.StatusInternalServerError)
}

func ReplyErrDef(bus ibus.IBus, sender interface{}, err error, defaultStatusCode int) {
	res := WrapSysError(err, defaultStatusCode).(SysError)
	ReplyJSON(bus, sender, res.HTTPStatus, res.ToJSON())
}

func ReplyErr(bus ibus.IBus, sender interface{}, err error) {
	ReplyErrDef(bus, sender, err, http.StatusInternalServerError)
}

func ReplyJSON(bus ibus.IBus, sender interface{}, httpCode int, body string) {
	bus.SendResponse(sender, ibus.Response{
		ContentType: ApplicationJSON,
		StatusCode:  httpCode,
		Data:        []byte(body),
	})
}

func ReplyBadRequest(bus ibus.IBus, sender interface{}, message string) {
	ReplyErrf(bus, sender, http.StatusBadRequest, message)
}

func replyAccessDenied(bus ibus.IBus, sender interface{}, code int, message string) {
	msg := "access denied"
	if len(message) > 0 {
		msg += ": " + message
	}
	ReplyErrf(bus, sender, code, msg)
}

func ReplyAccessDeniedUnauthorized(bus ibus.IBus, sender interface{}, message string) {
	replyAccessDenied(bus, sender, http.StatusUnauthorized, message)
}

func ReplyAccessDeniedForbidden(bus ibus.IBus, sender interface{}, message string) {
	replyAccessDenied(bus, sender, http.StatusForbidden, message)
}

func ReplyUnauthorized(bus ibus.IBus, sender interface{}, message string) {
	ReplyErrf(bus, sender, http.StatusUnauthorized, message)
}

func ReplyInternalServerError(bus ibus.IBus, sender interface{}, message string, err error) {
	ReplyErrf(bus, sender, http.StatusInternalServerError, message, ": ", err)
}
