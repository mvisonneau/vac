module github.com/mvisonneau/vac

go 1.15

require (
	github.com/go-ldap/ldap v3.0.2+incompatible // indirect
	github.com/hashicorp/vault/api v1.0.5-0.20200519221902-385fac77e20f
	github.com/hashicorp/vault/sdk v0.2.1
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/ktr0731/go-fuzzyfinder v0.2.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mvisonneau/go-helpers v0.0.1
	github.com/olekukonko/tablewriter v0.0.4
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli/v2 v2.3.0
	github.com/xeonx/timeago v1.0.0-rc4
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
)

replace github.com/ktr0731/go-fuzzyfinder => github.com/mvisonneau/go-fuzzyfinder v0.2.2-0.20200625134046-cc3ea9618b33
