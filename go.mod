module github.com/mvisonneau/vac

go 1.14

require (
	github.com/apex/log v1.4.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/ktr0731/go-fuzzyfinder v0.2.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mvisonneau/go-helpers v0.0.0-20200224131125-cb5cc4e6def9
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.4
)

replace github.com/ktr0731/go-fuzzyfinder => ../../os/go-fuzzyfinder
