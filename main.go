package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ResponseEnvelope struct {
	HTTPStatus int           `json:"httpStatus"`
	Data       interface{}   `json:"data"`
	Error      ResponseError `json:"error"`
	Meta       interface{}   `json:"extra"`
}

// ResponseError represents the JSON tagged error response object and also it is
// itself an `error` wrapping the **Business Logic** `error` at which it mapped
// to.
type ResponseError struct {
	httpStatus int
	Op         string      `json:"op"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Extra      interface{} `json:"meta"`
	Err        error       `json:"-"`
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("%s%s: %s", r.Op, r.Code, r.Err.Error())
}

func (r *ResponseError) Unwrap() error {
	return r.Err
}

type ResponseErrorsRegistry []ResponseError

// Business Logic expected errors
var (
	ErrVerbSvcInvalidSession  = errors.New("yasvc invalid session")
	ErrVerbSvcTimeout         = errors.New("yasvc internal call got timeout")
	ErrVerbSvcInvalidInput    = errors.New("yasvc input validation error")
	ErrVerbSvcInternal        = errors.New("yasvc internal error")
	ErrVerbSvcUnauthenticated = errors.New("yasvc not an authenticated user")
)

var dummyRes = ResponseError{
	httpStatus: http.StatusInternalServerError,
	Op:         "verbSvc",
	Code:       "errInternal",
	Message:    "internal error",
	Extra:      nil,
	Err:        ErrVerbSvcInternal,
}

func justPrint(payload ResponseError) {
	res := ResponseEnvelope{
		HTTPStatus: payload.httpStatus,
		Data:       nil,
		Meta:       nil,
		Error:      payload,
	}
	j, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Printing Response json...\n%v\n", string(j))
}

var verbSvcRegisteredErrors = []ResponseError{
	{
		httpStatus: http.StatusGatewayTimeout,
		Op:         "verbSvc",
		Code:       "errTimeout",
		Message:    "internal timeout",
		Extra:      nil,
		Err:        ErrVerbSvcTimeout,
	},
	{
		httpStatus: http.StatusBadRequest,
		Op:         "verbSvc",
		Code:       "errInvalidInput",
		Message:    "input validation error",
		Extra:      nil,
		Err:        ErrVerbSvcInvalidInput,
	},
	{
		httpStatus: http.StatusForbidden,
		Op:         "verbSvc",
		Code:       "errInvalidSession",
		Message:    "invalid session",
		Extra:      nil,
		Err:        ErrVerbSvcInvalidSession,
	},
	{
		httpStatus: http.StatusInternalServerError,
		Op:         "verbSvc",
		Code:       "errInternal",
		Message:    "internal error",
		Extra:      nil,
		Err:        ErrVerbSvcInternal,
	},
}

func main() {
	fmt.Println("Error info handling POC")
	fmt.Println("Refference Response...")
	justPrint(dummyRes)

	fmt.Printf("\nPrinting mapped errors...\n\n")

	fmt.Println("Mapping error errVerbSvcTimeout")
	e := fmt.Errorf("%w: %v", ErrVerbSvcTimeout, "ya504")
	fmt.Printf("Error from Business Logic⤵️\n\t%v\n", e)
	errorInfo := MapErrorToResponseError(verbSvcRegisteredErrors, e, "verbSvc")
	justPrint(errorInfo)

	fmt.Println("Mapping error errVerbSvcInvalidInput")
	e = fmt.Errorf("%w: %v", ErrVerbSvcInvalidInput, "ya bad input received")
	fmt.Printf("Error from Business Logic⤵️\n\t%v\n", e)
	errorInfo = MapErrorToResponseError(verbSvcRegisteredErrors, e, "verbSvc")
	justPrint(errorInfo)

	fmt.Println("Mapping error errVerbSvcInvalidSession")
	e = fmt.Errorf("%w: %v", ErrVerbSvcInvalidSession, "this user closed its session")
	fmt.Printf("Error from Business Logic⤵️\n\t%v\n", e)
	errorInfo = MapErrorToResponseError(verbSvcRegisteredErrors, e, "verbSvc")
	justPrint(errorInfo)

	fmt.Println("Mapping error errVerbSvcInternal")
	e = fmt.Errorf("%w: %v", ErrVerbSvcInternal, "Traditional Java null pointer exception")
	fmt.Printf("Error from Business Logic⤵️\n\t%v\n", e)
	errorInfo = MapErrorToResponseError(verbSvcRegisteredErrors, e, "verbSvc")
	justPrint(errorInfo)

	fmt.Println("Handling an unmapped error with an Internal Error")
	e = fmt.Errorf("%w: %v", ErrVerbSvcUnauthenticated, "who is this user?")
	fmt.Printf("Error from Business Logic⤵️\n\t%v\n", e)
	errorInfo = MapErrorToResponseError(verbSvcRegisteredErrors, e, "verbSvc")
	justPrint(errorInfo)
}

func MapErrorToResponseError(errorRegistry ResponseErrorsRegistry, err error, defaultOp string) ResponseError {
	for _, r := range errorRegistry {
		wrappedErr := errors.Unwrap(&r)
		if errors.Is(err, wrappedErr) {
			return r
		}
	}
	return ResponseError{
		httpStatus: http.StatusInternalServerError,
		Op:         defaultOp,
		Code:       "errInternal",
		Message:    "internal error",
		Extra:      nil,
		Err:        ErrVerbSvcInternal,
	}
}
