package act

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestApplyActionWithoutSinkFailsClosed pins what a misconfigured runner does at the moment it would write: it
// refuses. Before the fallback existed this dereferenced nil, so the failure mode of forgetting to wire a sink was
// a panic inside the run step rather than an error the run could report.
func TestApplyActionWithoutSinkFailsClosed(t *testing.T) {
	r := &SessionRunner{}
	_, err := r.ApplyAction(t.Context(), ActionRequest{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no action sink configured")
}
