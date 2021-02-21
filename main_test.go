package main

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

// Test for MapErrorToResponseError functionality
func TestErrorMapper(t *testing.T) {
	for _, test := range []struct {
		label              string
		inputError         error
		expectedHTTPStatus int
		expectedOp         string
		expectedCode       string
	}{
		{
			label:              "timeout-error-received",
			inputError:         fmt.Errorf("%w: %v", ErrVerbSvcTimeout, "ya504"),
			expectedHTTPStatus: http.StatusGatewayTimeout,
			expectedOp:         "verbSvc",
			expectedCode:       "errTimeout",
		},
		{
			label:              "invalid-input-received",
			inputError:         fmt.Errorf("%w: %v", ErrVerbSvcInvalidInput, "ya bad input received"),
			expectedHTTPStatus: http.StatusBadRequest,
			expectedOp:         "verbSvc",
			expectedCode:       "errInvalidInput",
		},
		{
			label:              "invalid-session-received",
			inputError:         fmt.Errorf("%w: %v", ErrVerbSvcInvalidSession, "this user closed its session"),
			expectedHTTPStatus: http.StatusForbidden,
			expectedOp:         "verbSvc",
			expectedCode:       "errInvalidSession",
		},
		{
			label:              "internal-error-received",
			inputError:         fmt.Errorf("%w: %v", ErrVerbSvcInternal, "Traditional Java null pointer exception"),
			expectedHTTPStatus: http.StatusInternalServerError,
			expectedOp:         "verbSvc",
			expectedCode:       "errInternal",
		},
		{
			label:              "unmmaped-error-received",
			inputError:         fmt.Errorf("%w: %v", ErrVerbSvcUnauthenticated, "who is this user?"),
			expectedHTTPStatus: http.StatusInternalServerError,
			expectedOp:         "verbSvc",
			expectedCode:       "errInternal",
		},
	} {
		e := MapErrorToResponseError(verbSvcRegisteredErrors, test.inputError, "verbSvc")
		responseError, ok := e.(ResponseError)
		if !ok {
			t.Errorf("For %q was expected a `ResponseError` type, but got %T", test.label, e)
		}

		wrappedE := errors.Unwrap(responseError)
		if !errors.Is(wrappedE, test.inputError) {
			t.Errorf("For %q was expected a wrapped error: %v, but got: %v", test.label, test.inputError, wrappedE)
		}

		contextualErrorMessage := fmt.Sprintf("%s: %s", responseError.Op, test.inputError.Error())
		if responseError.Error() != contextualErrorMessage {
			t.Errorf("For %q was expected an exact match for contextual error message: %v, but got: %v", test.label, test.inputError.Error(), responseError.Error())
		}

		if responseError.httpStatus != test.expectedHTTPStatus {
			t.Errorf("For %q was expected http status: %d, but got: %d", test.label, test.expectedHTTPStatus, responseError.httpStatus)
		}

		if responseError.Op != test.expectedOp {
			t.Errorf("For %q was expected operation: %s, but got: %s", test.label, test.expectedOp, responseError.Op)
		}

		if responseError.Code != test.expectedCode {
			t.Errorf("For %q was expected error code: %s, but got: %s", test.label, test.expectedCode, responseError.Code)
		}
	}
}
