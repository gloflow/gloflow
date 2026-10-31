// Package gloflow holds the version of the whole repo. VERSION (repo root) is embedded here so the Go code
// (gf_core.GetSysReleseInfo) reports the same version as the file that AGENTS.md tells us to bump.
package gloflow

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var versionFile string

// Version returns the content of VERSION, e.g. "0.11.0".
func Version() string {
	return strings.TrimSpace(versionFile)
}
