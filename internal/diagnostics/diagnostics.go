package diagnostics

import (
	"os"
	"strings"
)

// OSReleasePrettyName returns PRETTY_NAME from /etc/os-release (os-release(5)),
// e.g. `Ubuntu 24.04.1 LTS`, or "" when unavailable. Exported because
// internal/analytics reports the distro string.
//
// os-release(5): "/etc/os-release takes precedence over /usr/lib/os-release.
// Applications should check for the former, and exclusively use its data if
// it exists, and only fall back to /usr/lib/os-release if /etc/os-release
// does not exist." Most distros symlink the former to the latter, but some
// minimal/container images ship only /usr/lib/os-release.
func OSReleasePrettyName() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		data, err = os.ReadFile("/usr/lib/os-release")
		if err != nil {
			return ""
		}
	}
	return parseOSReleasePrettyName(string(data))
}

// parseOSReleasePrettyName is the pure parser behind OSReleasePrettyName.
// os-release(5) values are shell-compatible: unquoted, double- or
// single-quoted all occur in the wild (Alpine single-quotes).
func parseOSReleasePrettyName(data string) string {
	for line := range strings.SplitSeq(data, "\n") {
		v, ok := strings.CutPrefix(line, "PRETTY_NAME=")
		if !ok {
			continue
		}
		v = strings.TrimSpace(v)
		if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
			v = v[1 : len(v)-1]
		}
		return v
	}
	return ""
}
