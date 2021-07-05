module github.com/mvisonneau/vac

go 1.15

require (
	github.com/hashicorp/vault/api v1.0.4
	github.com/hashicorp/vault/sdk v0.1.13
	github.com/ktr0731/go-fuzzyfinder v0.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mvisonneau/go-helpers v0.0.1
	github.com/olekukonko/tablewriter v0.0.4
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli/v2 v2.3.0
	github.com/xeonx/timeago v1.0.0-rc4
)

replace github.com/ktr0731/go-fuzzyfinder => github.com/mvisonneau/go-fuzzyfinder v0.2.2-0.20200625134046-cc3ea9618b33
