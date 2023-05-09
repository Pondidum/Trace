# Changelog

## [0.0.4] - 2023-05-09

### Added

* `start` command's `--when` supports both `iso-8601` and `epoch-seconds` time formats now.

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
