# dhttp - StreamingFast HTTP Library

Go HTTP utilities for microservices with integrated logging, tracing, and error handling used in StreamingFast products.

## Key Features

- **Request handling**: Extract and validate URL parameters, query strings, and JSON payloads
- **Response writing**: JSON, text, HTML, and raw responses with proper headers and error logging
- **Error handling**: Structured HTTP error responses with tracing support
- **Handler wrappers**: Simplified handler patterns for common use cases
- **Utilities**: Real IP detection, response forwarding

## Core APIs

### Handlers
- [`JSONHandler`](https://pkg.go.dev/github.com/streamingfast/dhttp#JSONHandler) - Wrap functions that return JSON responses
- [`RawHandler`](https://pkg.go.dev/github.com/streamingfast/dhttp#RawHandler) - Stream raw content from io.ReadCloser
- [`DirectHandler`](https://pkg.go.dev/github.com/streamingfast/dhttp#DirectHandler) - Direct response writer control with error handling

### Request Processing
- [`ExtractRequest`](https://pkg.go.dev/github.com/streamingfast/dhttp#ExtractRequest) - Extract URL/query parameters into structs
- [`ExtractJSONRequest`](https://pkg.go.dev/github.com/streamingfast/dhttp#ExtractJSONRequest) - Parse and validate JSON request bodies

### Response Writing
- [`WriteJSON`](https://pkg.go.dev/github.com/streamingfast/dhttp#WriteJSON) - Write JSON responses
- [`WriteError`](https://pkg.go.dev/github.com/streamingfast/dhttp#WriteError) - Write structured error responses
- [`WriteText`](https://pkg.go.dev/github.com/streamingfast/dhttp#WriteText), [`WriteHTML`](https://pkg.go.dev/github.com/streamingfast/dhttp#WriteHTML) - Write text/HTML responses

### Error Types
- HTTP status error constructors: [`BadRequestError`](https://pkg.go.dev/github.com/streamingfast/dhttp#BadRequestError), [`NotFoundError`](https://pkg.go.dev/github.com/streamingfast/dhttp#NotFoundError), [`InternalServerError`](https://pkg.go.dev/github.com/streamingfast/dhttp#InternalServerError), etc.

See [full documentation](https://pkg.go.dev/github.com/streamingfast/dhttp) for complete API reference.
