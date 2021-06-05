// +build tools

package tools

import (
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/securego/gosec/v2/cmd/gosec"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
