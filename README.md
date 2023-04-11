# vac | AWS credentials management leveraging Vault

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mvisonneau/vac)](https://pkg.go.dev/mod/github.com/mvisonneau/vac)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/vac)](https://goreportcard.com/report/github.com/mvisonneau/vac)
[![test](https://github.com/mvisonneau/vac/actions/workflows/test.yml/badge.svg)](https://github.com/mvisonneau/vac/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/vac/badge.svg?branch=main)](https://coveralls.io/github/mvisonneau/vac?branch=main)
[![release](https://github.com/mvisonneau/vac/actions/workflows/release.yml/badge.svg)](https://github.com/mvisonneau/vac/actions/workflows/release.yml)

`vac` is a wrapper to manage AWS credentials dynamically using [Hashicorp Vault](https://www.vaultproject.io/).

It is heavily inspired from [jantman/vault-aws-creds](https://github.com/jantman/vault-aws-creds) and [ahmetb/kubectx](https://github.com/ahmetb/kubectx).

Written in golang, it can work on most common platforms (Linux, MacOS, Windows).

It leverages the [external process sourcing](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html) capabilities of the AWS CLI config definition.

[![asciicast](https://asciinema.org/a/343653.svg)](https://asciinema.org/a/343653?t=60)

## Install

Have a look onto the [latest release page](https://github.com/mvisonneau/vac/releases/latest) and pick your flavor.

Checksums are signed with the [following GPG key](https://keybase.io/mvisonneau/pgp_keys.asc): `C09C A9F7 1C5C 988E 65E3  E5FC ADEA 38ED C46F 25BE`

### Go

```bash
~$ go install github.com/mvisonneau/vac/cmd/vac@latest
~$ sudo setcap cap_ipc_lock=ep ${GOPATH:-~/go}/bin/vac
```

### Homebrew

```bash
~$ brew install mvisonneau/tap/vac
```

### Docker

```bash
~$ docker run -it --rm docker.io/mvisonneau/vac
~$ docker run -it --rm ghcr.io/mvisonneau/vac
~$ docker run -it --rm quay.io/mvisonneau/vac
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

- Once you have [installed it](#install), create a new profile in your `~/.aws/credentials` file:

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

- Finally assuming that you have sorted out your Vault accesses already, you need to chose which engine/role to use:

```bash
~$ vac
[follow prompt]
```

## Advanced usage

```bash
~$ vac --help
NAME:
   vac - Manage AWS credentials dynamically using Vault

USAGE:
   vac [global options] command [command options] [arguments...]

COMMANDS:
   get      get the creds in credential_process format (json)
   status   returns some info about the current context, cached credentials and Vault server connectivity
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --engine path, -e path  engine path [$VAC_ENGINE]
   --role name, -r name    role name [$VAC_ROLE]
   --state path, -s path   state path (default: "~/.vac_state") [$VAC_STATE_PATH]
   --log-level level       log level (debug,info,warn,fatal,panic) (default: "info") [$VAC_LOG_LEVEL]
   --log-format format     log format (json,text) (default: "text") [$VAC_LOG_FORMAT]
   --auth value            auth method (token, kubernetes) (default: "token") [$VAC_AUTH]
   --auth-k8s-role value   Kubernetes role to authenticate to (for --auth kubernetes) [$VAC_AUTH_K8S_ROLE]
   --auth-k8s-mount value  Kubernetes auth mount path (for --auth kubernetes) (default: "kubernetes") [$VAC_AUTH_K8S_MOUNT]
   --help, -h              show help
```

### Static configuration

You are forced to use the fuzzyfinding capabilities. This is particularily useful in a non-TTY usage scenario. eg:

```toml
# ~/.aws/credentials
[default]
credential_process = /usr/local/bin/vac get
[dev-admin]
credential_process = /usr/local/bin/vac -e dev -r admin get
[staging-admin]
credential_process = /usr/local/bin/vac -e staging -r admin get
```

You can also dynamically switch your context without a prompt by doing the following:

```bash
# only prompt for chosing a role in the "dev" engine
~$ vac -e dev

# no prompt, automatically switch to "admin" role in the "dev" engine
~$ vac -e dev -r admin
```

### TTL and cache bypass

The `get` command can take various flags in order to manage the credentials TTLs but also when to refresh them:

```bash
~$ vac get --help
NAME:
   vac get - get the creds in credential_process format (json)

USAGE:
   vac get [command options] [arguments...]

OPTIONS:
   --min-ttl duration           min-ttl duration (default: 0s) [$VAC_MIN_TTL]
   --ttl duration, -t duration  ttl duration (default: 0s) [$VAC_TTL]
   --force-generate, -f         bypass currently cached creds and generate new ones [$VAC_FORCE_GENERATE]
```

#### Examples

```bash
# Generate credentials valid for 1h
~$ vac get --ttl 1h

# Generate credentials valid for 1h but replace them if existing ones expire in less than 30m
~$ vac get --ttl 1h --min-ttl 30m

# Generate credentials valid for 2h, indepently if some valid ones are still present in the cache
~$ vac get --ttl 2 -f
```

you can of course define them in your `~/.aws/credentials` profiles as well

```bash
~$ cat - <<EOF >> ~/.aws/credentials
[vac-4h]
credential_process = $(which vac) get --ttl 4h
[vac-no-cache]
credential_process = $(which vac) get -f
EOF
```

### Get information about current configuration

You can use the `status` command in order to retrieve some info about:

- Current context (selected engine/role)
- Cached credentials
- Vault server connectivity details

```
~$ vac status
+----------------+---------------------+
|  LOCAL STATE   |                     |
+----------------+---------------------+
| Current Engine | dev                 |
| Current Role   | admin               |
+----------------+---------------------+
+-----------+--------+---------------+
|  ENGINE   |  ROLE  |  EXPIRATION   |
+-----------+--------+---------------+
| dev       | admin  | in 2 hours    |
| prod      | admin  | 2 days ago    |
| staging   | admin  | 2 days ago    |
+-----------+--------+---------------+
+-------------+--------------------------------------+
|    VAULT    |                                      |
+-------------+--------------------------------------+
| ClusterID   | 0e6b2fcd-e84b-a7cd-f84d-6b31947a8d73 |
| ClusterName | vault-cluster-90f72c95               |
| Initialized | true                                 |
| Sealed      | false                                |
| Version     | 1.5.3                                |
+-------------+--------------------------------------+
```

## Required Vault policies

To be able to use all the lookup features, you will need some 

```hcl
# List available AWS engines
path "sys/mounts" {
  capabilities = ["read"]
}

# List all the roles for each <engine_path>
# path "<engine_path>/roles" {
#  capabilities = ["list"]
# }
# eg:
path "dev/roles" {
 capabilities = ["list"]
}

path "staging/roles" {
 capabilities = ["list"]
}

# Assume the role
# path "<engine_path>/sts/<role>" {
#  capabilities = ["update"]
# }
# eg:

path "dev/sts/admin" {
  capabilities = ["update"]
}

path "staging/sts/admin" {
  capabilities = ["update"]
}
```

## Limitations

- It currently **only supports** authenticating using STS [assumed_role](https://www.vaultproject.io/docs/secrets/aws#sts-assumerole) or [federation_tokens](https://www.vaultproject.io/docs/secrets/aws#sts-federation-tokens) methods.
- It will list all available engines and roles according to the defined policies. This result may not be relevant to what an user can actually assume.

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/vac/pulls).
