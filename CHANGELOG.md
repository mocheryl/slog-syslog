# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2024-10-22

### Fixed

- Add option to enable writing source code origin of the logged event which is
  now be default disabled. To restore previous behaviour, set `AddSource`
  property in `Options` to true when passing it to the handler's constructor.

## [0.1.0] - 2024-07-20

### Added

- Structured log syslog handler that writes to a syslog server in a format as
  used by syslog package from the standard library.
