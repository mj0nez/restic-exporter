package internal

import promVersionInfo "github.com/prometheus/common/version"

// populated at build time
//
// WARN do not change to const !
var (
	Version  = "1.0.0"
	Revision = "abcdefghi123456789"
)

func init() {
	promVersionInfo.Version = Version
	promVersionInfo.Revision = Revision
}
