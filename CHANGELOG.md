# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [0ver](https://0ver.org) (more or less).

## [Unreleased]

### Changed

- Refactored codebase following golang standard structure
- Bumped all dependencies to their latest version
- Bumped to go `1.15`

## [0.0.3] - 2020-09-03

### Added

- securego/gosec tests

### Changed

- Bumped golang to 1.15
- Bumped goreleaser to 0.142.0
- Bumped urfave/cli to v2

### Removed

- Dropped support for darwin/386

## [0.0.2] - 2020-06-29

### Added

- New `ttl`, `min-ttl` and `force-generate` flags on the **get** function to manipulate credentials lengths
- New `status` function to disclose some info about the current context, cached credentials and Vault server connectivity

### Changed

- Removed some typos in the CLI flags definition
- Removed unused parameter RenewBefore on the AWSCredential objects
- Added some tests

## [0.0.1] - 2020-06-26

### Added

- Working state of the app
- Makefile
- LICENSE
- README

[Unreleased]: https://github.com/mvisonneau/vac/compare/0.0.2...HEAD
[0.0.2]: https://github.com/mvisonneau/vac/tree/0.0.2
[0.0.1]: https://github.com/mvisonneau/vac/tree/0.0.1
