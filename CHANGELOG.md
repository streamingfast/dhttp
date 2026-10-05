# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

* Add `middleware.NewAddTraceIDHeaderMiddleware()` to return the request's trace ID in the `X-Trace-ID` response header.

### Changed

* Require Go 1.26 or newer.

### Security

* Bump `google.golang.org/grpc` to v1.83.2, `golang.org/x/crypto` to v0.57.0, `golang.org/x/net` to v0.59.0, OpenTelemetry to v1.47.0 and `github.com/gorilla/schema` to v1.4.1 to fix known vulnerabilities.

## v0.1.2

* Added license file so [pkg.go.dev/github.com/streamingfast/dhttp](https://pkg.go.dev/github.com/streamingfast/dhttp@latest) renders correctly.

## v0.1.1

* Added more documentation.

## v0.1.0

* Added `dhttp.<>Error` to represents all HTTP errors.

* Added `dhttp.WriteErrorResponse(...)` in addition to `dhttp.WriteError(...)` to control the message logged.

* Initial release.