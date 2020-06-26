# vac | AWS credentials management leveraging Vault

[![GoDoc](https://godoc.org/github.com/mvisonneau/vac?status.svg)](https://godoc.org/github.com/mvisonneau/vac/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/vac)](https://goreportcard.com/report/github.com/mvisonneau/vac)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/vac.svg)](https://hub.docker.com/r/mvisonneau/vac/)
[![Build Status](https://cloud.drone.io/api/badges/mvisonneau/vac/status.svg)](https://cloud.drone.io/mvisonneau/vac)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/vac/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/vac?branch=master)

`vac` is a wrapper to manage AWS credentials dynamically using [Hashicorp Vault](https://www.vaultproject.io/). It is heavily inspired from [jantman/vault-aws-creds](https://github.com/jantman/vault-aws-creds) and [ahmetb/kubectx](https://github.com/ahmetb/kubectx). Written in golang, it can work on most common platforms (Linux, MacOS, Windows).

## TL:DR

There will be a nice GIF or asciicinema here soon!

## Install

### Go

```bash
~$ go get -u github.com/mvisonneau/vac
```

### Homebrew

```bash
~$ brew install mvisonneau/tap/vac
```

### Docker

```bash
~$ docker run -it --rm mvisonneau/vac
```

### Scoop

```bash
~$ scoop bucket add https://github.com/mvisonneau/scoops
~$ scoop install vac
```

### Binaries, DEB and RPM packages

Have a look onto the [latest release page](https://github.com/mvisonneau/vac/releases/latest) to pick your flavor and version. Here is an helper to fetch the most recent one:

```bash
~$ export VAC_VERSION=$(curl -s "https://api.github.com/repos/mvisonneau/vac/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
```

```bash
# Binary (eg: linux/amd64)
~$ wget https://github.com/mvisonneau/vac/releases/download/${VAC_VERSION}/vac_${VAC_VERSION}_linux_amd64.tar.gz
~$ tar zxvf vac_${VAC_VERSION}_linux_amd64.tar.gz -C /usr/local/bin

# DEB package (eg: linux/386)
~$ wget https://github.com/mvisonneau/vac/releases/download/${VAC_VERSION}/vac_${VAC_VERSION}_linux_386.deb
~$ dpkg -i vac_${VAC_VERSION}_linux_386.deb

# RPM package (eg: linux/arm64)
~$ wget https://github.com/mvisonneau/vac/releases/download/${VAC_VERSION}/vac_${VAC_VERSION}_linux_arm64.rpm
~$ rpm -ivh vac_${VAC_VERSION}_linux_arm64.rpm
```

## Quickstart

- Once you have installed it, create a new profile in your `~/.aws/credentials` file:

```bash
~$ cat - <<EOF >> ~/.aws/credentials

[vac]
credential_process = $(which vac) get
EOF
```

- You will need to set the following env variable to use this profile

```bash
~$ export AWS_PROFILE=vac
```

(you can omit this part by using it as your default profile instead)

- Finally asusming that you have sorted out your Vault accesses already, you need to chose which engine/role to use:

```bash
~$ vac
[follow prompt]
```

## Usage

```bash
~$ vac --help
NAME:
   vac - Manage AWS credentials dynamically using Vault

USAGE:
   vac [global options] command [command options] [arguments...]

COMMANDS:
   get      get the creds in credential_process format (json)
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --engine path, -e path  engine path [$VAC_ENGINE]
   --role name, -r name    role name [$VAC_ROLE]
   --state path, -s path   state path (default: "~/.vac_state") [$VAC_STATE_PATH]
   --log-level level       log level (debug,info,warn,fatal,panic) (default: "info") [$VAC_LOG_LEVEL]
   --log-format format     log format (json,text) (default: "text") [$VAC_LOG_FORMAT]
   --help, -h              show help
```

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/vac/pulls).
