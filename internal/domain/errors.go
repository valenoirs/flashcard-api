package domain

import "errors"

// domain can have a shared error because at the end of the day
// we'll write a domain error map in the delivery layer,
// so we can handle the "user not found" or entity related error there
var (
	ErrNotFound            = errors.New("record not found")           // 404
	ErrConflict            = errors.New("record already exist")       // 409
	ErrUnprocessableEntity = errors.New("unprocessable entity")       // 422
	ErrPreconditionFailed  = errors.New("precondition not meet")      // 412
	ErrGone                = errors.New("record permanently deleted") // 410
	ErrResourceLocked      = errors.New("record locked")              // 423
	ErrUnsuportedMediaType = errors.New("unsupported media type")     // 415
)
