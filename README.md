# [POC] Contextual Error Information Handling for Catched Error in Transport Layer

The intend of the **POC** is to validate some about and to settled an approach for handling error contextual information at Transport Layer (http transport in this particular case) in a `Go` application.

So for a **Business Logic** error catched a *Transport Layer* we can be able to index a richer error interface valuable to this specific transport (`http - application/json`).

## POC Result

```bash
Printing mapped errors...

Mapping error errVerbSvcTimeout
Error from Business Logic⤵️
        yasvc internal call got timeout: ya504
Printing Response json...
{
  "httpStatus": 504,
  "data": null,
  "error": {
    "op": "verbSvc",
    "code": "errTimeout",
    "message": "internal timeout",
    "meta": null
  },
  "extra": null
}
Mapping error errVerbSvcInvalidInput
Error from Business Logic⤵️
        yasvc input validation error: ya bad input received
Printing Response json...
{
  "httpStatus": 400,
  "data": null,
  "error": {
    "op": "verbSvc",
    "code": "errInvalidInput",
    "message": "input validation error",
    "meta": null
  },
  "extra": null
}
Mapping error errVerbSvcInvalidSession
Error from Business Logic⤵️
        yasvc invalid session: this user closed its session
Printing Response json...
{
  "httpStatus": 403,
  "data": null,
  "error": {
    "op": "verbSvc",
    "code": "errInvalidSession",
    "message": "invalid session",
    "meta": null
  },
  "extra": null
}
Mapping error errVerbSvcInternal
Error from Business Logic⤵️
        yasvc internal error: Traditional Java null pointer exception
Printing Response json...
{
  "httpStatus": 500,
  "data": null,
  "error": {
    "op": "verbSvc",
    "code": "errInternal",
    "message": "internal error",
    "meta": null
  },
  "extra": null
}
Handling an unmapped error with an Internal Error
Error from Business Logic⤵️
        yasvc not an authenticated user: who is this user?
Printing Response json...
{
  "httpStatus": 500,
  "data": null,
  "error": {
    "op": "verbSvc",
    "code": "errInternal",
    "message": "internal error",
    "meta": null
  },
  "extra": null
}
```
## Test Coverage Result

```bash
main.go:29:   Error                   100.0%
main.go:33:   Unwrap                  100.0%
main.go:57:   justPrint               0.0%
main.go:111:  MapErrorToResponseError 100.0%
main.go:135:  main                    0.0%
total:        (statements)            15.9%
```
