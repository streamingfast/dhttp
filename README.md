## dfuse HTTP Library

This repository contains all common stuff around HTTP handling across our
various services.

### Philosophy

The package contains wide adoption common method re-used over and over again
across our micro-services. The methods are around the following subjects:

- Utilities
- Middlewares
- Requests
- Responses

They usually perform the most standard operation handling everything related
to logging, tracing and error handling.

### Reference

- [WriteJSON](#writejson)

#### Utilities

#### Requests

##### `ExtractRequest`

##### `ExtractJSONRequest`

#### Responses

##### `WriteError`

##### `WriteJSON`

Writes a struct as a JSON body for a particular handler correctly handling logging
of errors and correctly sets all headers
