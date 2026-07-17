package act

import (
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// requirePostgresInternal is the package-internal counterpart of the external requirePostgres helper: it returns the
// spike DSN and skips the test if Postgres is unreachable, so an internal test that needs the pool can run under the
// same local container.
func requirePostgresInternal(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("ACT_TEST_DBOS_URL")
	if dsn == "" {
		dsn = "postgres://postgres:dbos@localhost:55432/dbos_spike?sslmode=disable"
	}
	u, err := url.Parse(dsn)
	require.NoError(t, err)
	conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
	if err != nil {
		t.Skipf("act internal tests need Postgres at %s (docker start dbos-spike-pg): %v", u.Host, err)
	}
	_ = conn.Close()
	return dsn
}

// TestMigrateBackfillsLeaseColumns proves Migrate is idempotent and repairs a pre-existing agent_actions table (A1). A
// database created by an earlier build lacks attempt_id and lease_expiry, and since CREATE TABLE IF NOT EXISTS never
// alters an existing table, every projection that selects attempt_id would fail until Migrate backfills the columns.
func TestMigrateBackfillsLeaseColumns(t *testing.T) {
	dsn := requirePostgresInternal(t)
	ctx := t.Context()
	schema := "act_int_" + uuid.New().String()[:8]
	store, err := NewPostgresRunStore(ctx, StoreConfig{DatabaseURL: dsn, Schema: schema})
	require.NoError(t, err)
	t.Cleanup(func() { store.DropSchemaForTest(ctx); store.Close() })

	require.NoError(t, store.Migrate(ctx))

	// Simulate an older database: drop the lease columns the earlier build never created.
	_, err = store.pool.Exec(ctx, `ALTER TABLE `+store.t("agent_actions")+` DROP COLUMN attempt_id, DROP COLUMN lease_expiry`)
	require.NoError(t, err)

	// A projection that selects attempt_id must now fail: the table is back in the pre-migration shape.
	_, err = store.GetAction(ctx, "i", "r", "c")
	require.Error(t, err, "reading a table missing attempt_id fails until the backfill runs")

	// Re-migrate: the backfill ALTERs restore both columns; running Migrate a further time proves idempotency.
	require.NoError(t, store.Migrate(ctx))
	require.NoError(t, store.Migrate(ctx))

	// A full propose -> approve -> claim cycle now works, proving the restored columns are present and usable.
	require.NoError(t, store.ProposeAction(ctx, NewAction{
		ToolCallID: "c", RunID: "r", InstanceID: "i", Tool: "t", Connector: "k",
		ArgsHash: "sha256:x", IdempotencyKey: "act-x", PolicyDecision: PolicyApprovalRequired,
	}))
	_, err = store.TransitionAction(ctx, ActionTransition{InstanceID: "i", RunID: "r", ToolCallID: "c", Status: ActionApproved})
	require.NoError(t, err)
	claimed, action, err := store.ClaimExecuting(ctx, "i", "r", "c", "attempt-1", 90*time.Second)
	require.NoError(t, err)
	require.True(t, claimed)
	require.Equal(t, "attempt-1", action.AttemptID)
	require.Equal(t, ActionExecuting, action.Status)
}
