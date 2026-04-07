// named/pending_bb_test.go v1
package named_test

import (
	"os"
	"testing"
)

func skipIfPending(t *testing.T, reason string) {
	t.Helper()
	if os.Getenv("RUN_PENDING_TESTS") == "" {
		t.Skip("PENDING: " + reason + " — set RUN_PENDING_TESTS=1 to run")
	}
}
