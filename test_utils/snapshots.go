package test_utils

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/yosssi/gohtml"
)

// MakeSnapshot formats and creates a snapshot of the
// given string and compares it with the existing snapshot
func MakeSnapshot(t *testing.T, got string) {
	got = gohtml.Format(got)
	snaps.MatchSnapshot(t, got)
}
