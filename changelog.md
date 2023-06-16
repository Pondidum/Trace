# Changelog

## [0.0.9] - 2023-06-16

### Added

* `group finish` and `finish` commands support the `--error` flag to set the span status to Error.  Optionally you can write `--error="some message"` to set the span status description.

## [0.0.8] - 2023-05-25

### Fixed

* Action failed as it was using a shell script, which isn't supported with a `composite` action type

## [0.0.7] - 2023-05-24

### Added

* `start` command's `--when` also supports `RFC3339` format dates

## [0.0.6] - 2023-05-10

### Fixed

* Improve `pondidum/trace` action's cache usage for GHE

## [0.0.5] - 2023-05-09

### Added

* Added github actions step to setup the tool: `pondidum/trace`

## [0.0.4] - 2023-05-09

### Added

* `start` command's `--when` supports both `iso-8601` and `epoch-seconds` time formats now.

### Fixed

* the `attr` command is now available!

## [0.0.3] - 2023-05-04

### Added

* Support configuring the OTLP endpoint via environment variables
* `start` now supports a `--when <date>` argument to allow using a different time for when a trace starts

## [0.0.2] - 2023-03-05

### Added

* all command support `--attr key=value` for setting additional span attributes

## [0.0.1] - 2023-02-27

### Added

- implemented most basic functionality
- added `attr` command to append attributes to spans
- rename `span` to `group`
- add `task` command

## [0.0.0] - 2023-02-13

### Added

- Initial Version
