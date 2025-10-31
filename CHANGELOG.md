# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [0ver](https://0ver.org) (more or less).

## [Unreleased]

### Changed

- Golang updated to `1.20`
- Bumped all dependencies
- internal/cli: enhanced flags implementation

## [v0.0.8] - 2021-11-15

### Added

- Support for tokens containing a trailing `\n`

### Changed

- Golang updated to `1.19`
- Bumped all dependencies
  
## [v0.0.7] - 2021-02-11

### Added

- Released container images over `quay.io`
### Changed

- Bumped all dependencies

## [v0.0.6] - 2021-08-19

### Changed

- Generate pre-releases on default branch pushes
- Bumped go from **1.15** to **1.17**
- Updated all dependencies to their latest versions
- Release npm/deb packages correctly
- Add support for Apple M1 silicon

## [v0.0.5] - 2020-12-17

### Added

- Release GitHub container registry based images: [ghcr.io/mvisonneau/vac](https://github.com/users/mvisonneau/packages/container/package/vac)
- Release `arm64v8` based container images as part of docker manifests in both **docker.io** and **ghcr.io**
- GPG sign released artifacts checksums

### Changed

- Prefix new releases with `^v` to make `pkg.go.dev` happy
- Updated all dependencies
- Migrated CI from Drone to GitHub actions

## [0.0.4] - 2020-10-22

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

[Unreleased]: https://github.com/mvisonneau/vac/compare/v0.0.8...HEAD
[v0.0.8]: https://github.com/mvisonneau/vac/tree/v0.0.8
[v0.0.7]: https://github.com/mvisonneau/vac/tree/v0.0.7
[v0.0.6]: https://github.com/mvisonneau/vac/tree/v0.0.6
[v0.0.5]: https://github.com/mvisonneau/vac/tree/v0.0.5
[0.0.4]: https://github.com/mvisonneau/vac/tree/0.0.4
[0.0.3]: https://github.com/mvisonneau/vac/tree/0.0.3
[0.0.2]: https://github.com/mvisonneau/vac/tree/0.0.2
[0.0.1]: https://github.com/mvisonneau/vac/tree/0.0.1
