package act

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// defaultListLimit caps a listing that did not request a limit. It bounds a single page so an unbounded catalog of
// runs or events can never be pulled in one query.
const defaultListLimit = 100

// schemaNamePattern constrains the store's schema name to a bare SQL identifier. The schema is interpolated into
// DDL and DML (pgx cannot parameterize identifiers), so restricting it to this shape is what keeps that
// interpolation injection-safe.
var schemaNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// StoreConfig configures a PostgresRunStore.
type StoreConfig struct {
	// DatabaseURL is the Postgres DSN for the Act product tables. For Act this is the Platform Postgres instance,
	// logical database "act". It is the same instance the DBOS tables live in, but the product tables sit in their
	// own schema (Schema), separate from the DBOS system schema (§15.2).
	DatabaseURL string
	// Schema is the schema the product tables live in. It must be a bare SQL identifier. Defaults to "act_product".
	Schema string
}

// PostgresRunStore is the Postgres-backed RunStore. It owns a pgxpool and its own schema; it does not share tables
// with the workflow engine. Reads are plain SELECTs so the API's query path never depends on DBOS being healthy.
type PostgresRunStore struct {
	pool       *pgxpool.Pool
	schema     string
	sharedPool bool // when true the pool is owned by the caller and Close leaves it open
}

var _ RunStore = (*PostgresRunStore)(nil)

// NewPostgresRunStore opens a pool against cfg.DatabaseURL and returns a store for cfg.Schema. It does not create
// the schema; call Migrate for that. The caller owns the returned store's lifetime and must Close it.
func NewPostgresRunStore(ctx context.Context, cfg StoreConfig) (*PostgresRunStore, error) {
	schema := cfg.Schema
	if schema == "" {
		schema = "act_product"
	}
	if !schemaNamePattern.MatchString(schema) {
		return nil, fmt.Errorf("act: invalid store schema %q", schema)
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("act: StoreConfig.DatabaseURL is required")
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("act: connect run store: %w", err)
	}
	return &PostgresRunStore{pool: pool, schema: schema}, nil
}

// NewPostgresRunStoreWithPool returns a store over an already-open pool, for callers that share one pool across the
// Act components. The store does not take ownership: Close is a no-op, and the caller closes the pool.
func NewPostgresRunStoreWithPool(pool *pgxpool.Pool, schema string) (*PostgresRunStore, error) {
	if schema == "" {
		schema = "act_product"
	}
	if !schemaNamePattern.MatchString(schema) {
		return nil, fmt.Errorf("act: invalid store schema %q", schema)
	}
	if pool == nil {
		return nil, errors.New("act: NewPostgresRunStoreWithPool requires a non-nil pool")
	}
	// pool is left owned by the caller: shared stores do not close it (Close is a no-op unless this store opened it).
	return &PostgresRunStore{pool: pool, schema: schema, sharedPool: true}, nil
}

// Close releases the pool if this store owns it. A store built over a shared pool leaves it open for the owner.
func (s *PostgresRunStore) Close() {
	if s.pool != nil && !s.sharedPool {
		s.pool.Close()
	}
}

// t qualifies a table name with the store's schema. schema is validated at construction, so this is injection-safe.
func (s *PostgresRunStore) t(name string) string { return s.schema + "." + name }

// DropSchemaForTest drops the store's schema and everything in it. It exists only so a test that migrates a
// throwaway schema can leave the database clean; production never drops the product tables. schema is validated at
// construction, so the interpolation is injection-safe.
func (s *PostgresRunStore) DropSchemaForTest(ctx context.Context) {
	_, _ = s.pool.Exec(ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %s CASCADE`, s.schema))
}

// ClobberCanonicalArgsForTest overwrites an approval's stored preimage, which CreateApproval refuses to do because
// it enforces the binding to args_hash. It exists so a test can reproduce the two states no legitimate write path
// can produce any more: a row from before the column existed (nil) and a row whose bytes were corrupted after the
// fact. Both must be refused at the approve endpoint rather than signed blind.
func (s *PostgresRunStore) ClobberCanonicalArgsForTest(ctx context.Context, instanceID, approvalID string, canonicalArgs *string) error {
	_, err := s.pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET canonical_args=$3 WHERE approval_id=$1 AND instance_id=$2`,
		s.t("agent_approvals")), approvalID, instanceID, canonicalArgs)
	return err
}

// BackdateRunForTest moves a run's updated_on into the past. It exists so a test can reproduce a run that has sat
// untouched for long enough to be considered stranded, which the startup repair requires before it closes anything;
// no legitimate write path moves updated_on backwards.
func (s *PostgresRunStore) BackdateRunForTest(ctx context.Context, instanceID, runID string, age time.Duration) error {
	_, err := s.pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET updated_on = now() - $3::interval WHERE run_id=$1 AND instance_id=$2`,
		s.t("agent_runs")), runID, instanceID, age.String())
	return err
}

// BackdateApprovalDecisionsForTest moves a run's approval decisions into the past. The startup repair measures
// staleness from decided_on, not from the run row, so a test that wants to look stranded has to age the decisions
// themselves; ageing only the run is precisely the mistake that made an earlier version of that predicate able to
// cancel a live run.
func (s *PostgresRunStore) BackdateApprovalDecisionsForTest(ctx context.Context, instanceID, runID string, age time.Duration) error {
	_, err := s.pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET decided_on = now() - $3::interval WHERE run_id=$1 AND instance_id=$2 AND decided_on IS NOT NULL`,
		s.t("agent_approvals")), runID, instanceID, age.String())
	return err
}

// ForceTerminalWithoutCloseForTest marks a run terminal by writing the status column directly, bypassing CloseRun and
// so leaving any pending approvals attached. It reproduces what the previous binary did, which a replayed DBOS step
// can still reproduce during the deploy that ships CloseRun: the step's output is checkpointed, so the new body that
// withdraws approvals never runs. No production path may do this, which is why it lives behind a test-only name.
func (s *PostgresRunStore) ForceTerminalWithoutCloseForTest(ctx context.Context, instanceID, runID string, status RunStatus) error {
	if !status.IsTerminal() {
		return fmt.Errorf("act: ForceTerminalWithoutCloseForTest requires a terminal status, got %q", status)
	}
	_, err := s.pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET status=$3, updated_on=now(), finished_on=COALESCE(finished_on, now()) WHERE run_id=$1 AND instance_id=$2`,
		s.t("agent_runs")), runID, instanceID, string(status))
	return err
}

func (s *PostgresRunStore) Migrate(ctx context.Context) error {
	// One statement per Exec: pgx's simple-protocol batching aside, keeping them separate makes a failure point to
	// the exact DDL. All are IF NOT EXISTS so Migrate is safe to run on every startup.
	stmts := []string{
		fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, s.schema),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			run_id text PRIMARY KEY,
			instance_id text NOT NULL,
			organization_id text,
			project_id text,
			agent_name text NOT NULL,
			spec_hash text NOT NULL DEFAULT '',
			trigger text NOT NULL DEFAULT '',
			trigger_ref text NOT NULL DEFAULT '',
			idempotency_key text NOT NULL DEFAULT '',
			conversation_id text NOT NULL DEFAULT '',
			actor_subject text NOT NULL DEFAULT '',
			actor_service_principal boolean NOT NULL DEFAULT false,
			status text NOT NULL,
			error text NOT NULL DEFAULT '',
			created_on timestamptz NOT NULL DEFAULT now(),
			updated_on timestamptz NOT NULL DEFAULT now(),
			started_on timestamptz,
			finished_on timestamptz
		)`, s.t("agent_runs")),

		// Listings are always instance-scoped and newest-first, and the inbox/status views filter by agent/status;
		// these composite indexes lead with instance_id so a tenant's slice is contiguous.
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_runs_instance_created_idx ON %s (instance_id, created_on DESC)`, s.t("agent_runs")),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_runs_instance_agent_idx ON %s (instance_id, agent_name, created_on DESC)`, s.t("agent_runs")),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_runs_instance_status_idx ON %s (instance_id, status, created_on DESC)`, s.t("agent_runs")),

		// Backfill trigger_ref onto an agent_runs table created by an earlier build: CREATE TABLE IF NOT EXISTS never
		// alters an existing table, so a database migrated before this column existed would lack it and every run
		// projection (which now selects trigger_ref) would fail with a missing-column error. ADD COLUMN IF NOT EXISTS is
		// idempotent and safe on both a fresh and a pre-existing table.
		fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS trigger_ref text NOT NULL DEFAULT ''`, s.t("agent_runs")),

		// dedupe_key is the per-run idempotency key of an event. For the once-per-run lifecycle edges it equals
		// event_type; for the approval edges the executor scopes it per segment ("run.waiting_approval:1"), so a run
		// that pauses on several governed actions records every pause while a replayed edge still conflicts onto its
		// own row. The uniqueness lives on (run_id, dedupe_key) — created as a separate index below so the fresh and
		// the migrated table converge on the same shape.
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			id bigserial PRIMARY KEY,
			run_id text NOT NULL,
			instance_id text NOT NULL,
			seq bigint NOT NULL,
			event_type text NOT NULL,
			dedupe_key text NOT NULL DEFAULT '',
			status text NOT NULL DEFAULT '',
			payload jsonb,
			visibility text NOT NULL DEFAULT 'user',
			created_on timestamptz NOT NULL DEFAULT now(),
			UNIQUE (run_id, seq)
		)`, s.t("agent_run_events")),

		// The stream/poll path reads events for one run with id greater than a cursor, in id order; id is the global
		// monotonic sequence (bigserial), so this index serves both the scoping and the cursor.
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_run_events_instance_run_id_idx ON %s (instance_id, run_id, id)`, s.t("agent_run_events")),

		// Migrate an agent_run_events table created by an earlier build, which keyed idempotency on
		// UNIQUE (run_id, event_type): that uniqueness silently swallowed the second waiting_approval/resumed of a run
		// with several governed actions, leaving the status stuck on running (kairos-cloud#129). Add dedupe_key, backfill
		// it from event_type — suffixing the approval edges with ":0", the segment every pre-migration run implicitly
		// was on, so a replayed old edge still deduplicates against the new keys — then move the uniqueness over and
		// drop the old constraint. Every statement is idempotent: the backfills match no rows once dedupe_key is
		// populated (new inserts always carry it), and the index/constraint statements are IF (NOT) EXISTS.
		fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS dedupe_key text NOT NULL DEFAULT ''`, s.t("agent_run_events")),
		fmt.Sprintf(`UPDATE %s SET dedupe_key = event_type || ':0' WHERE dedupe_key = '' AND event_type IN ('%s','%s','%s')`,
			s.t("agent_run_events"), EventTypeWaitingApproval, EventTypeResumed, EventTypeRejected),
		fmt.Sprintf(`UPDATE %s SET dedupe_key = event_type WHERE dedupe_key = ''`, s.t("agent_run_events")),
		fmt.Sprintf(`CREATE UNIQUE INDEX IF NOT EXISTS agent_run_events_run_dedupe_idx ON %s (run_id, dedupe_key)`, s.t("agent_run_events")),
		fmt.Sprintf(`ALTER TABLE %s DROP CONSTRAINT IF EXISTS agent_run_events_run_id_event_type_key`, s.t("agent_run_events")),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			approval_id text PRIMARY KEY,
			run_id text NOT NULL,
			instance_id text NOT NULL,
			tool_name text NOT NULL DEFAULT '',
			connector text NOT NULL DEFAULT '',
			tool_call_id text NOT NULL DEFAULT '',
			args_hash text NOT NULL,
			canonical_args text,
			proposal text NOT NULL DEFAULT '',
			policy text NOT NULL DEFAULT '',
			status text NOT NULL,
			requested_by text NOT NULL DEFAULT '',
			decided_by text NOT NULL DEFAULT '',
			created_on timestamptz NOT NULL DEFAULT now(),
			decided_on timestamptz
		)`, s.t("agent_approvals")),

		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_approvals_instance_status_idx ON %s (instance_id, status, created_on DESC)`, s.t("agent_approvals")),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_approvals_instance_run_idx ON %s (instance_id, run_id)`, s.t("agent_approvals")),

		// The action ledger (§15.1): one row per proposed external write, keyed by (instance, run, tool call) so a
		// replayed proposal conflicts onto the same row. Secret values never land here — only redacted args/results and
		// secret references (§16.2). It is a peer table of the run store, in the same product schema, distinct from the
		// DBOS system tables.
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			instance_id text NOT NULL,
			run_id text NOT NULL,
			tool_call_id text NOT NULL,
			agent_name text NOT NULL DEFAULT '',
			tool text NOT NULL DEFAULT '',
			connector text NOT NULL DEFAULT '',
			tool_version text NOT NULL DEFAULT '',
			class text NOT NULL DEFAULT '',
			args_hash text NOT NULL,
			idempotency_key text NOT NULL,
			redacted_args jsonb,
			policy_decision text NOT NULL DEFAULT '',
			policy_reason text NOT NULL DEFAULT '',
			proposal text NOT NULL DEFAULT '',
			status text NOT NULL,
			attempt_id text NOT NULL DEFAULT '',
			lease_expiry timestamptz,
			external_reference text NOT NULL DEFAULT '',
			redacted_result jsonb,
			redacted_message text NOT NULL DEFAULT '',
			redacted_verification jsonb,
			error text NOT NULL DEFAULT '',
			requested_by text NOT NULL DEFAULT '',
			decided_by text NOT NULL DEFAULT '',
			created_on timestamptz NOT NULL DEFAULT now(),
			updated_on timestamptz NOT NULL DEFAULT now(),
			executed_on timestamptz,
			verified_on timestamptz,
			PRIMARY KEY (instance_id, run_id, tool_call_id)
		)`, s.t("agent_actions")),

		// Audit and inbox views list a run's actions; the idempotency-key index backs a lookup by dedup key.
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_actions_instance_run_idx ON %s (instance_id, run_id, created_on DESC)`, s.t("agent_actions")),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_actions_idem_idx ON %s (instance_id, idempotency_key)`, s.t("agent_actions")),

		// Backfill the execution-lease columns onto an agent_actions table created by an earlier build: CREATE TABLE IF
		// NOT EXISTS never alters an existing table, so a database migrated before these columns existed would lack them
		// and every projection that selects attempt_id (and the recovery path that reads lease_expiry) would fail with a
		// missing-column error. ADD COLUMN IF NOT EXISTS is idempotent and safe on both a fresh and a pre-existing table.
		fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS attempt_id text NOT NULL DEFAULT ''`, s.t("agent_actions")),
		fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS lease_expiry timestamptz`, s.t("agent_actions")),
		// Backfill the action's redacted text result onto a table created before this column existed (same rationale as
		// the lease columns above): a tool that returns only text lands here, distinct from the structured redacted_result.
		fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS redacted_message text NOT NULL DEFAULT ''`, s.t("agent_actions")),

		// Backfill position/total columns for per-tool-call approvals: position is the 1-based index of the action within
		// its batch, total is the batch size. Zero defaults are fine for old single-action approvals.
		fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS position int NOT NULL DEFAULT 0`, s.t("agent_approvals")),
		fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS total int NOT NULL DEFAULT 0`, s.t("agent_approvals")),

		// Add canonical_args: the exact byte preimage of args_hash, persisted so an approver can be shown the very
		// bytes their decision binds to. There is nothing to backfill INTO old rows: the bytes those hashes were
		// computed over were never stored, and recomputing them from the truncated proposal would fabricate a
		// preimage. Old approvals stay NULL, which is why the column is nullable rather than NOT NULL DEFAULT '':
		// refusing to sign what cannot be shown is now an authorization decision, and it must not confuse "no
		// preimage was ever stored" with "the preimage is the empty string".
		fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS canonical_args text`, s.t("agent_approvals")),
		// A database that already has the column from an earlier build of this same change carries it as
		// NOT NULL DEFAULT '': ADD COLUMN IF NOT EXISTS would leave that in place, and its rows would then read
		// back as an empty preimage (a hash contradiction) instead of as the absence they are. Converge the shape
		// explicitly, then restore the absence those rows should have had.
		//
		// Guarded on attnotnull rather than run unconditionally: Migrate runs on every startup, and an unconditional
		// UPDATE ... WHERE canonical_args='' would scan the whole table each time forever. Once the column is
		// nullable the conversion has already happened, so the block does nothing and costs one catalog lookup.
		fmt.Sprintf(`DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM pg_attribute a
				WHERE a.attrelid = '%s'::regclass AND a.attname = 'canonical_args' AND a.attnotnull
			) THEN
				ALTER TABLE %s ALTER COLUMN canonical_args DROP NOT NULL;
				ALTER TABLE %s ALTER COLUMN canonical_args DROP DEFAULT;
				UPDATE %s SET canonical_args=NULL WHERE canonical_args='';
			END IF;
		END $$`, s.t("agent_approvals"), s.t("agent_approvals"), s.t("agent_approvals"), s.t("agent_approvals")),

		// Drop the approval expiry column: approval deadlines have been removed. The column is no longer written or
		// read, and historical rows that carried an expiry date lose nothing of value (the deadline was never enforced
		// in production). The companion statements that normalized historical "expired" statuses to "cancelled" are
		// gone: both deployments ran them and hold zero expired rows, so they only survived as a false lead — reading
		// them suggests an expiry mechanism that has not existed since the column was dropped.
		fmt.Sprintf(`ALTER TABLE %s DROP COLUMN IF EXISTS expires_on`, s.t("agent_approvals")),

		// Referential integrity: every child table's run_id must point to an existing agent_runs row. Each FK is added
		// NOT VALID (short AccessExclusive lock, no full-table scan), orphan rows are cleaned up (they are unreachable
		// through any API path since no parent run exists — today there is no delete path for runs, so orphans should
		// only come from manual database cleanup or test harnesses), and then the constraint is validated
		// (ShareUpdateExclusive, concurrent DML proceeds). Every step is idempotent: the DO block checks
		// pg_constraint so ADD CONSTRAINT runs at most once, the DELETE is a no-op when no orphans remain, and
		// VALIDATE on an already-valid constraint is a no-op.

		// -- agent_run_events.run_id → agent_runs.run_id
		fmt.Sprintf(`DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint c JOIN pg_namespace n ON n.oid = c.connamespace
				WHERE c.conname = 'agent_run_events_run_id_fkey' AND n.nspname = '%s'
			) THEN
				ALTER TABLE %s ADD CONSTRAINT agent_run_events_run_id_fkey
					FOREIGN KEY (run_id) REFERENCES %s (run_id) ON DELETE CASCADE NOT VALID;
			END IF;
		END $$`, s.schema, s.t("agent_run_events"), s.t("agent_runs")),
		fmt.Sprintf(`DELETE FROM %s WHERE NOT EXISTS (SELECT 1 FROM %s r WHERE r.run_id = %s.run_id)`,
			s.t("agent_run_events"), s.t("agent_runs"), s.t("agent_run_events")),
		fmt.Sprintf(`ALTER TABLE %s VALIDATE CONSTRAINT agent_run_events_run_id_fkey`, s.t("agent_run_events")),

		// -- agent_approvals.run_id → agent_runs.run_id
		fmt.Sprintf(`DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint c JOIN pg_namespace n ON n.oid = c.connamespace
				WHERE c.conname = 'agent_approvals_run_id_fkey' AND n.nspname = '%s'
			) THEN
				ALTER TABLE %s ADD CONSTRAINT agent_approvals_run_id_fkey
					FOREIGN KEY (run_id) REFERENCES %s (run_id) ON DELETE CASCADE NOT VALID;
			END IF;
		END $$`, s.schema, s.t("agent_approvals"), s.t("agent_runs")),
		fmt.Sprintf(`DELETE FROM %s WHERE NOT EXISTS (SELECT 1 FROM %s r WHERE r.run_id = %s.run_id)`,
			s.t("agent_approvals"), s.t("agent_runs"), s.t("agent_approvals")),
		fmt.Sprintf(`ALTER TABLE %s VALIDATE CONSTRAINT agent_approvals_run_id_fkey`, s.t("agent_approvals")),

		// -- agent_actions.run_id → agent_runs.run_id
		fmt.Sprintf(`DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint c JOIN pg_namespace n ON n.oid = c.connamespace
				WHERE c.conname = 'agent_actions_run_id_fkey' AND n.nspname = '%s'
			) THEN
				ALTER TABLE %s ADD CONSTRAINT agent_actions_run_id_fkey
					FOREIGN KEY (run_id) REFERENCES %s (run_id) ON DELETE CASCADE NOT VALID;
			END IF;
		END $$`, s.schema, s.t("agent_actions"), s.t("agent_runs")),
		fmt.Sprintf(`DELETE FROM %s WHERE NOT EXISTS (SELECT 1 FROM %s r WHERE r.run_id = %s.run_id)`,
			s.t("agent_actions"), s.t("agent_runs"), s.t("agent_actions")),
		fmt.Sprintf(`ALTER TABLE %s VALIDATE CONSTRAINT agent_actions_run_id_fkey`, s.t("agent_actions")),

		// Repair runs stranded by the write asymmetry this release fixes: a shutdown cut short the approval wait, the
		// approval sweep ran on a context that outlived it while the run's terminal transition did not, and the run was
		// left reading waiting_approval with every approval cancelled and no decider. Such a run is not decidable (there
		// is nothing pending to approve) and not resumable (its workflow is gone), so it shows in the inbox forever
		// offering nothing (kairos-cloud#135). Both deployments carry these; the fix above stops new ones, and this
		// clears the ones already written, on startup, so nobody runs UPDATEs against a production database by hand.
		//
		// The predicate is deliberately narrow, and each clause earns its place by excluding a run that is alive:
		//
		//   - no approval still pending: there is genuinely nothing for a human to decide.
		//   - at least one approval cancelled with NO decider: the specific fingerprint of this bug. A withdrawal
		//     names no decider, so this separates "the process withdrew them" from "a human decided them all".
		//   - the most recent decision is over an hour old. This is the clause that must NOT be read off the run's
		//     updated_on: that column records when the run entered waiting_approval, so a run parked for hours and
		//     approved one second before a replica boots would satisfy an age check on it and be cancelled while
		//     live. Reading the age off the approvals' decided_on measures the thing that actually has to be stale.
		//     A null MAX (rows predating the column) compares false, which leaves the run alone.
		// It is serialized with an advisory lock because two replicas booting together would otherwise each compute
		// MAX(seq)+1 for the same run and collide on the (run_id, seq) uniqueness, failing a migration over a repair.
		// The lock is transaction-scoped, so it releases on commit and a crashed replica cannot hold it. The rolling
		// deploy is safe for a second reason too: an old replica can still strand a run while the new one boots, but
		// the hour of grace on the most recent decision keeps this from touching anything that recent.
		fmt.Sprintf(`
			WITH locked AS (SELECT pg_advisory_xact_lock(hashtext('%s.act_repair_stranded'))),
			stranded AS (
				SELECT r.run_id, r.instance_id
				FROM %s r CROSS JOIN locked
				WHERE r.status = '%s'
				  AND NOT EXISTS (SELECT 1 FROM %s a WHERE a.run_id = r.run_id AND a.instance_id = r.instance_id AND a.status = '%s')
				  AND EXISTS (SELECT 1 FROM %s a WHERE a.run_id = r.run_id AND a.instance_id = r.instance_id
				              AND a.status = '%s' AND a.decided_by = '')
				  AND (SELECT MAX(a.decided_on) FROM %s a WHERE a.run_id = r.run_id AND a.instance_id = r.instance_id)
				      < now() - interval '1 hour'
			), explained AS (
				INSERT INTO %s (run_id, instance_id, seq, event_type, dedupe_key, status, payload, visibility)
				SELECT s.run_id, s.instance_id,
					COALESCE((SELECT MAX(e.seq) FROM %s e WHERE e.run_id = s.run_id), 0) + 1,
					'%s', '%s', '%s',
					jsonb_build_object('reason', 'The run was interrupted while it waited for approval and could not be resumed. Its pending approvals were withdrawn at the time; this closes the run itself.'),
					'user'
				FROM stranded s
				ON CONFLICT (run_id, dedupe_key) DO NOTHING
				RETURNING run_id
			)
			UPDATE %s r SET status = '%s', updated_on = now(), finished_on = COALESCE(r.finished_on, now())
			FROM stranded s WHERE r.run_id = s.run_id AND r.instance_id = s.instance_id`,
			s.schema,
			s.t("agent_runs"), RunStatusWaitingApproval,
			s.t("agent_approvals"), ApprovalStatusPending,
			s.t("agent_approvals"), ApprovalStatusCancelled,
			s.t("agent_approvals"),
			s.t("agent_run_events"),
			s.t("agent_run_events"),
			EventTypeCancelled, EventTypeCancelled, RunStatusCancelled,
			s.t("agent_runs"), RunStatusCancelled),

		// The same invariant from the other side: an approval left pending on a run that already reached a terminal
		// state. Nothing can come of deciding it (the run is over), but the inbox offers the decision anyway. A
		// replayed DBOS step is one way to produce it — a run that recorded its terminal transition under the previous
		// binary replays that step's checkpointed output and so skips the new body that withdraws approvals — and any
		// future write path that marks a run terminal without going through CloseRun is another. The hour of grace
		// keeps it off runs that are being closed right now, where the two writes are simply in flight.
		fmt.Sprintf(`
			UPDATE %s a SET status = '%s', decided_on = now()
			WHERE a.status = '%s'
			  AND EXISTS (
				SELECT 1 FROM %s r
				WHERE r.run_id = a.run_id AND r.instance_id = a.instance_id
				  AND r.status IN ('%s','%s','%s','%s')
				  AND r.updated_on < now() - interval '1 hour'
			  )`,
			s.t("agent_approvals"), ApprovalStatusCancelled, ApprovalStatusPending,
			s.t("agent_runs"),
			RunStatusSucceeded, RunStatusFailed, RunStatusRejected, RunStatusCancelled),
	}

	for _, stmt := range stmts {
		if _, err := s.pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("act: migrate run store: %w", err)
		}
	}
	return nil
}

func (s *PostgresRunStore) CreateRun(ctx context.Context, run NewRun) error {
	if run.RunID == "" || run.InstanceID == "" || run.AgentName == "" {
		return errors.New("act: CreateRun requires RunID, InstanceID and AgentName")
	}
	return s.inTx(ctx, func(tx pgx.Tx) error {
		// Insert the run in the queued state. ON CONFLICT DO NOTHING makes a replayed Start a no-op rather than a
		// duplicate-key error: idempotency is enforced here, not only inside the workflow engine.
		ct, err := tx.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (run_id, instance_id, organization_id, project_id, agent_name, spec_hash, trigger,
				trigger_ref, idempotency_key, conversation_id, actor_subject, actor_service_principal, status)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
			ON CONFLICT (run_id) DO NOTHING`, s.t("agent_runs")),
			run.RunID, run.InstanceID, nullIfEmpty(run.OrganizationID), nullIfEmpty(run.ProjectID), run.AgentName,
			run.SpecHash, run.Trigger, run.TriggerRef, run.IdempotencyKey, run.ConversationID, run.Actor.Subject,
			run.Actor.ServicePrincipal, string(RunStatusQueued))
		if err != nil {
			return fmt.Errorf("insert run: %w", err)
		}
		if ct.RowsAffected() == 0 {
			// The run already existed: nothing more to do (its queued event was appended on the first create).
			return nil
		}
		// The run was just inserted, so its queued event is always new; the inserted flag cannot be false here.
		_, err = s.appendEventTx(ctx, tx, run.RunID, run.InstanceID, EventTypeQueued, "", RunStatusQueued, nil)
		return err
	})
}

func (s *PostgresRunStore) RecordTransition(ctx context.Context, tr RunTransition) error {
	if tr.RunID == "" || tr.InstanceID == "" || tr.Status == "" || tr.EventType == "" {
		return errors.New("act: RecordTransition requires RunID, InstanceID, Status and EventType")
	}
	// A terminal status must go through CloseRun, which withdraws the run's approvals in the same transaction.
	// Refusing it here is what makes "a terminal run has no pending approvals" an invariant rather than a habit: with
	// this door open, any caller could mark a run terminal and leave its approvals pending just by picking the more
	// obvious method name (kairos-cloud#135).
	if tr.Status.IsTerminal() {
		return fmt.Errorf("act: RecordTransition refuses the terminal status %q: use CloseRun, which also withdraws the run's approvals", tr.Status)
	}
	return s.inTx(ctx, func(tx pgx.Tx) error {
		return s.recordTransitionTx(ctx, tx, tr)
	})
}

// CloseRun records a run's terminal transition and withdraws its still-pending approvals in ONE transaction, so no
// reader can ever observe one without the other. Two writes done separately is what stranded runs showing
// waiting_approval with every approval cancelled and no decider: the sweep landed on a context that outlived the
// shutdown, the transition did not (kairos-cloud#135). Callers closing a run must use this, not RecordTransition
// followed by CancelPendingApprovals.
//
// The approvals are cancelled even when the transition is a no-op (the run was already terminal, or this edge was
// already recorded): the invariant being enforced is "a terminal run has no pending approvals", and an already-terminal
// run holding one is exactly the inconsistency to clear. Returns how many approvals were cancelled.
func (s *PostgresRunStore) CloseRun(ctx context.Context, tr RunTransition) (int64, error) {
	if tr.RunID == "" || tr.InstanceID == "" || tr.Status == "" || tr.EventType == "" {
		return 0, errors.New("act: CloseRun requires RunID, InstanceID, Status and EventType")
	}
	if !tr.Status.IsTerminal() {
		return 0, fmt.Errorf("act: CloseRun requires a terminal status, got %q", tr.Status)
	}
	var cancelled int64
	err := s.inTx(ctx, func(tx pgx.Tx) error {
		if err := s.recordTransitionTx(ctx, tx, tr); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, fmt.Sprintf(`UPDATE %s SET status=$3, decided_on=now() WHERE run_id=$1 AND instance_id=$2 AND status=$4`, s.t("agent_approvals")),
			tr.RunID, tr.InstanceID, ApprovalStatusCancelled, ApprovalStatusPending)
		if err != nil {
			return fmt.Errorf("cancel pending approvals: %w", err)
		}
		cancelled = tag.RowsAffected()

		// The ledger's counterpart, in the same transaction. This store IS the ledger in production, so an action left
		// in approval_pending beside a withdrawn approval is the same half-applied state one table over. The executor
		// also withdraws actions best-effort afterwards, which covers the deployment where the ledger sits behind the
		// gateway instead; doing it here as well makes the common case atomic rather than merely likely.
		if _, err := tx.Exec(ctx, fmt.Sprintf(`UPDATE %s SET status=$3, updated_on=now() WHERE run_id=$1 AND instance_id=$2 AND status=$4`, s.t("agent_actions")),
			tr.RunID, tr.InstanceID, ActionWithdrawn, ActionApprovalPending); err != nil {
			return fmt.Errorf("withdraw pending actions: %w", err)
		}
		return nil
	})
	return cancelled, err
}

// recordTransitionTx applies a status change and appends its event inside an open transaction. It is the shared body
// of RecordTransition and CloseRun; the caller owns the transaction so a close can bind the transition and the
// approval sweep into one atomic unit.
func (s *PostgresRunStore) recordTransitionTx(ctx context.Context, tx pgx.Tx, tr RunTransition) error {
	// Lock the run row: this validates existence and instance scoping, and serializes transitions for the run so
	// the per-run event sequence stays monotonic under concurrency.
	var current string
	err := tx.QueryRow(ctx, fmt.Sprintf(`SELECT status FROM %s WHERE run_id=$1 AND instance_id=$2 FOR UPDATE`, s.t("agent_runs")),
		tr.RunID, tr.InstanceID).Scan(&current)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrRunNotFound
	}
	if err != nil {
		return fmt.Errorf("lock run: %w", err)
	}

	// Guard terminal states: a run that has reached a terminal status never transitions again. Without this, a
	// replayed transition of an earlier edge (a DBOS step re-running an old run.running after the run already
	// succeeded) would regress the status back out of its terminal state. Terminal is final, so this is a no-op.
	if RunStatus(current).IsTerminal() {
		return nil
	}

	// Append the event first, and let its (run_id, dedupe_key) uniqueness decide idempotency: if this edge was
	// already recorded, no row is inserted and we must NOT touch the status. This is what stops a re-applied
	// transition from moving the status a second time (e.g. regressing a run whose later edges already ran). The
	// key defaults to the event type (strict once-per-run); the executor scopes the repeatable approval edges per
	// segment so a second governed pause is a NEW edge — it inserts and moves the status — while a replay of any
	// recorded edge still carries the same key and stays a no-op.
	inserted, err := s.appendEventTx(ctx, tx, tr.RunID, tr.InstanceID, tr.EventType, tr.DedupeKey, tr.Status, tr.Payload)
	if err != nil {
		return err
	}
	if !inserted {
		// Duplicate edge: the status already reflects this transition. Leave the row untouched.
		return nil
	}

	// Update lifecycle state. started_on latches on the first move into running; finished_on latches on the
	// first move into a terminal status; error and spec_hash are written only when supplied.
	_, err = tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s SET
			status=$3,
			updated_on=now(),
			started_on = COALESCE(started_on, CASE WHEN $3=$4 THEN now() END),
			finished_on = CASE WHEN $5 THEN COALESCE(finished_on, now()) ELSE finished_on END,
			error = CASE WHEN $6<>'' THEN $6 ELSE error END,
			spec_hash = CASE WHEN $7<>'' THEN $7 ELSE spec_hash END
		WHERE run_id=$1 AND instance_id=$2`, s.t("agent_runs")),
		tr.RunID, tr.InstanceID, string(tr.Status), string(RunStatusRunning), tr.Status.IsTerminal(), tr.Error, tr.SpecHash)
	if err != nil {
		return fmt.Errorf("update run: %w", err)
	}
	return nil
}

// appendEventTx appends one lifecycle event inside an open transaction and reports whether a new row was inserted.
// The next per-run sequence is computed under the run-row lock the caller holds, and the insert is idempotent on
// (run_id, dedupe_key) — dedupeKey defaults to eventType when empty: re-appending the same edge after a crash
// inserts nothing and returns false, which is what lets RecordTransition detect a duplicate edge and leave the
// run's status untouched.
func (s *PostgresRunStore) appendEventTx(ctx context.Context, tx pgx.Tx, runID, instanceID, eventType, dedupeKey string, status RunStatus, payload map[string]any) (bool, error) {
	if dedupeKey == "" {
		dedupeKey = eventType
	}
	var payloadJSON []byte
	if payload != nil {
		var err error
		payloadJSON, err = json.Marshal(payload)
		if err != nil {
			return false, fmt.Errorf("marshal event payload: %w", err)
		}
	}
	ct, err := tx.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (run_id, instance_id, seq, event_type, dedupe_key, status, payload, visibility)
		VALUES ($1, $2, (SELECT COALESCE(MAX(seq),0)+1 FROM %s WHERE run_id=$1), $3, $4, $5, $6, $7)
		ON CONFLICT (run_id, dedupe_key) DO NOTHING`, s.t("agent_run_events"), s.t("agent_run_events")),
		runID, instanceID, eventType, dedupeKey, string(status), payloadJSON, VisibilityUser)
	if err != nil {
		return false, fmt.Errorf("append event: %w", err)
	}
	return ct.RowsAffected() > 0, nil
}

// SetRunConversationID links a run to the AI session it executed in. It overwrites conversation_id unconditionally
// (the create path seeds it, usually empty) and is scoped by instance so it is tenant-safe; a run that does not
// exist for the instance yields ErrRunNotFound rather than silently affecting no rows.
func (s *PostgresRunStore) SetRunConversationID(ctx context.Context, instanceID, runID, conversationID string) error {
	if runID == "" || instanceID == "" {
		return errors.New("act: SetRunConversationID requires RunID and InstanceID")
	}
	ct, err := s.pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET conversation_id=$3, updated_on=now() WHERE run_id=$1 AND instance_id=$2`, s.t("agent_runs")),
		runID, instanceID, conversationID)
	if err != nil {
		return fmt.Errorf("act: set run conversation id: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrRunNotFound
	}
	return nil
}

func (s *PostgresRunStore) GetRun(ctx context.Context, instanceID, runID string) (*Run, error) {
	row := s.pool.QueryRow(ctx, s.selectRunSQL()+` WHERE run_id=$1 AND instance_id=$2`, runID, instanceID)
	run, err := scanRun(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRunNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("act: get run: %w", err)
	}
	return run, nil
}

func (s *PostgresRunStore) ListRuns(ctx context.Context, f ListRunsFilter) ([]*Run, error) {
	if f.InstanceID == "" {
		return nil, errors.New("act: ListRuns requires InstanceID")
	}
	limit := f.Limit
	if limit <= 0 {
		limit = defaultListLimit
	}
	// Optional filters are folded into the WHERE via NULL-guards so the query text is fixed: an empty AgentName or
	// Status matches everything, which keeps the prepared statement stable and the index usable. The access
	// allow-list rides the same way: a nil slice arrives as SQL NULL (no restriction), a non-nil one restricts in
	// the WHERE itself so the page is full of rows the caller may actually see.
	rows, err := s.pool.Query(ctx, s.selectRunSQL()+`
		WHERE instance_id=$1
		  AND ($2='' OR agent_name=$2)
		  AND ($3='' OR status=$3)
		  AND ($4::text[] IS NULL OR agent_name = ANY($4))
		ORDER BY created_on DESC
		LIMIT $5`, f.InstanceID, f.AgentName, string(f.Status), f.AccessibleAgents, limit)
	if err != nil {
		return nil, fmt.Errorf("act: list runs: %w", err)
	}
	defer rows.Close()

	var out []*Run
	for rows.Next() {
		run, err := scanRun(rows)
		if err != nil {
			return nil, fmt.Errorf("act: scan run: %w", err)
		}
		out = append(out, run)
	}
	return out, rows.Err()
}

func (s *PostgresRunStore) ListRunEvents(ctx context.Context, instanceID, runID string, afterID int64, limit int) ([]*RunEvent, error) {
	if instanceID == "" || runID == "" {
		return nil, errors.New("act: ListRunEvents requires InstanceID and RunID")
	}
	if limit <= 0 {
		limit = defaultListLimit
	}
	rows, err := s.pool.Query(ctx, fmt.Sprintf(`
		SELECT id, run_id, instance_id, seq, event_type, status, payload, visibility, created_on
		FROM %s
		WHERE instance_id=$1 AND run_id=$2 AND id>$3
		ORDER BY id ASC
		LIMIT $4`, s.t("agent_run_events")), instanceID, runID, afterID, limit)
	if err != nil {
		return nil, fmt.Errorf("act: list run events: %w", err)
	}
	defer rows.Close()

	var out []*RunEvent
	for rows.Next() {
		var (
			e           RunEvent
			status      string
			payloadJSON []byte
		)
		if err := rows.Scan(&e.ID, &e.RunID, &e.InstanceID, &e.Seq, &e.EventType, &status, &payloadJSON, &e.Visibility, &e.CreatedOn); err != nil {
			return nil, fmt.Errorf("act: scan run event: %w", err)
		}
		e.Status = RunStatus(status)
		if len(payloadJSON) > 0 {
			if err := json.Unmarshal(payloadJSON, &e.Payload); err != nil {
				return nil, fmt.Errorf("act: unmarshal event payload: %w", err)
			}
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}

func (s *PostgresRunStore) CreateApproval(ctx context.Context, a NewApproval) error {
	if a.ApprovalID == "" || a.RunID == "" || a.InstanceID == "" || a.ArgsHash == "" {
		return errors.New("act: CreateApproval requires ApprovalID, RunID, InstanceID and ArgsHash")
	}
	// Enforce the preimage binding where it is written, not only where it is read. An approval that cannot show an
	// approver what they are signing cannot be signed at all (the API refuses it), so one created without a
	// coherent preimage would be dead on arrival: better to fail the caller here, loudly, than to persist a row
	// that only reveals its defect when a human is waiting on it.
	//
	// An empty preimage is rejected even when its hash matches. It would be a legitimate value in the abstract —
	// the SHA-256 of no bytes is a real hash — but nothing describes an action with it: the gateway path always
	// canonicalizes to at least "{}", and a simulated proposal with no text describes nothing to approve. Allowing
	// it would only create a state that survives to the wire as an empty string and there becomes indistinguishable
	// from "no preimage recorded", which is exactly the ambiguity the nullable column exists to remove.
	if a.CanonicalArgs == "" {
		return fmt.Errorf("act: CreateApproval requires CanonicalArgs (approval %q)", a.ApprovalID)
	}
	if HashArgs(a.CanonicalArgs) != a.ArgsHash {
		return fmt.Errorf("act: CreateApproval requires CanonicalArgs to be the preimage of ArgsHash (approval %q)", a.ApprovalID)
	}
	// Lock the run and refuse to attach a pending approval to one that is already over. Without the lock this races
	// CloseRun in the losing direction: CloseRun withdraws what is pending and commits, this insert lands a moment
	// later, and the run ends up terminal with a live approval nobody's decision can affect. Taking the same row lock
	// CloseRun takes serializes the two, and the terminal check makes the loser fail instead of writing that state.
	return s.inTx(ctx, func(tx pgx.Tx) error {
		var current string
		err := tx.QueryRow(ctx, fmt.Sprintf(`SELECT status FROM %s WHERE run_id=$1 AND instance_id=$2 FOR UPDATE`, s.t("agent_runs")),
			a.RunID, a.InstanceID).Scan(&current)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRunNotFound
		}
		if err != nil {
			return fmt.Errorf("act: create approval: lock run: %w", err)
		}
		if RunStatus(current).IsTerminal() {
			return fmt.Errorf("act: CreateApproval refuses to attach approval %q to run %q, which is already %s",
				a.ApprovalID, a.RunID, current)
		}
		_, err = tx.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (approval_id, run_id, instance_id, tool_name, connector, tool_call_id, args_hash,
				canonical_args, proposal, policy, status, requested_by, position, total)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
			ON CONFLICT (approval_id) DO NOTHING`, s.t("agent_approvals")),
			a.ApprovalID, a.RunID, a.InstanceID, a.ToolName, a.Connector, a.ToolCallID, a.ArgsHash,
			a.CanonicalArgs, a.Proposal, a.Policy, ApprovalStatusPending, a.RequestedBy, a.Position, a.Total)
		if err != nil {
			return fmt.Errorf("act: create approval: %w", err)
		}
		return nil
	})
}

func (s *PostgresRunStore) GetApproval(ctx context.Context, instanceID, approvalID string) (*Approval, error) {
	row := s.pool.QueryRow(ctx, s.selectApprovalSQL()+` WHERE a.approval_id=$1 AND a.instance_id=$2`, approvalID, instanceID)
	a, err := scanApproval(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrApprovalNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("act: get approval: %w", err)
	}
	return a, nil
}

func (s *PostgresRunStore) ListApprovals(ctx context.Context, f ListApprovalsFilter) ([]*Approval, error) {
	if f.InstanceID == "" {
		return nil, errors.New("act: ListApprovals requires InstanceID")
	}
	limit := f.Limit
	if limit <= 0 {
		limit = defaultListLimit
	}
	// The access allow-list resolves through the approval's run (an approval row does not carry its agent, but
	// the projection already joins it), still inside the WHERE so pagination stays truthful; see ListRuns.
	rows, err := s.pool.Query(ctx, s.selectApprovalSQL()+`
		WHERE a.instance_id=$1
		  AND ($2='' OR a.run_id=$2)
		  AND ($3='' OR a.status=$3)
		  AND ($4::text[] IS NULL OR r.agent_name = ANY($4))
		ORDER BY a.created_on DESC
		LIMIT $5`, f.InstanceID, f.RunID, f.Status, f.AccessibleAgents, limit)
	if err != nil {
		return nil, fmt.Errorf("act: list approvals: %w", err)
	}
	defer rows.Close()

	var out []*Approval
	for rows.Next() {
		a, err := scanApproval(rows)
		if err != nil {
			return nil, fmt.Errorf("act: scan approval: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *PostgresRunStore) ResolveApproval(ctx context.Context, instanceID, approvalID, status, decidedBy string) (*Approval, error) {
	if status != ApprovalStatusApproved && status != ApprovalStatusDenied {
		return nil, fmt.Errorf("act: invalid approval resolution %q", status)
	}
	var a *Approval
	err := s.inTx(ctx, func(tx pgx.Tx) error {
		// Claim the approval under a row lock: only a pending approval can be resolved, so two concurrent decisions
		// serialize and the second sees a non-pending status. This is the guard behind double-submit safety.
		var current string
		err := tx.QueryRow(ctx, fmt.Sprintf(`SELECT status FROM %s WHERE approval_id=$1 AND instance_id=$2 FOR UPDATE`, s.t("agent_approvals")),
			approvalID, instanceID).Scan(&current)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrApprovalNotFound
		}
		if err != nil {
			return fmt.Errorf("lock approval: %w", err)
		}
		if current != ApprovalStatusPending {
			return ErrApprovalNotResolvable
		}
		_, err = tx.Exec(ctx, fmt.Sprintf(`UPDATE %s SET status=$3, decided_by=$4, decided_on=now() WHERE approval_id=$1 AND instance_id=$2`, s.t("agent_approvals")),
			approvalID, instanceID, status, decidedBy)
		if err != nil {
			return fmt.Errorf("update approval: %w", err)
		}
		a, err = scanApproval(tx.QueryRow(ctx, s.selectApprovalSQL()+` WHERE a.approval_id=$1 AND a.instance_id=$2`, approvalID, instanceID))
		if err != nil {
			return fmt.Errorf("reload approval: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// CancelPendingApprovals moves every pending approval of a run to cancelled in one statement, stamping decided_on but
// no decider (the run was cancelled, not decided). It is idempotent: a run with no pending approval cancels zero rows.
func (s *PostgresRunStore) CancelPendingApprovals(ctx context.Context, instanceID, runID string) (int64, error) {
	tag, err := s.pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET status=$3, decided_on=now() WHERE run_id=$1 AND instance_id=$2 AND status=$4`, s.t("agent_approvals")),
		runID, instanceID, ApprovalStatusCancelled, ApprovalStatusPending)
	if err != nil {
		return 0, fmt.Errorf("act: cancel pending approvals: %w", err)
	}
	return tag.RowsAffected(), nil
}

// selectRunSQL is the column list and source for a run projection. COALESCE folds the two nullable tenancy columns
// to empty strings so every consumer scans into a plain Run without null handling.
func (s *PostgresRunStore) selectRunSQL() string {
	return fmt.Sprintf(`
		SELECT run_id, instance_id, COALESCE(organization_id,''), COALESCE(project_id,''), agent_name, spec_hash,
			trigger, trigger_ref, idempotency_key, conversation_id, actor_subject, actor_service_principal, status, error,
			created_on, updated_on, started_on, finished_on
		FROM %s`, s.t("agent_runs"))
}

// selectApprovalSQL is the column list and source for an approval projection. It joins the approval's run
// (every approval has one, enforced by the FK) to carry the agent name and the run's actor: the pair a reader
// needs to resolve the approve policy for this concrete approval. Consumers must qualify shared column names
// (a.run_id, a.instance_id, a.status) in their WHERE clauses.
func (s *PostgresRunStore) selectApprovalSQL() string {
	return fmt.Sprintf(`
		SELECT a.approval_id, a.run_id, a.instance_id, r.agent_name, r.actor_subject, a.tool_name, a.connector,
			a.tool_call_id, a.args_hash, a.canonical_args, a.proposal, a.policy, a.status, a.requested_by, a.decided_by,
			a.created_on, a.decided_on, a.position, a.total
		FROM %s a JOIN %s r ON r.run_id = a.run_id AND r.instance_id = a.instance_id`,
		s.t("agent_approvals"), s.t("agent_runs"))
}

// rowScanner is the read surface shared by pgx.Row (single-row) and pgx.Rows (iterated), so one scan helper serves
// both Get and List.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanRun(row rowScanner) (*Run, error) {
	var (
		run    Run
		status string
	)
	err := row.Scan(&run.RunID, &run.InstanceID, &run.OrganizationID, &run.ProjectID, &run.AgentName, &run.SpecHash,
		&run.Trigger, &run.TriggerRef, &run.IdempotencyKey, &run.ConversationID, &run.Actor.Subject, &run.Actor.ServicePrincipal,
		&status, &run.Error, &run.CreatedOn, &run.UpdatedOn, &run.StartedOn, &run.FinishedOn)
	if err != nil {
		return nil, err
	}
	run.Status = RunStatus(status)
	return &run, nil
}

func scanApproval(row rowScanner) (*Approval, error) {
	var a Approval
	// canonical_args scans into a pointer so a NULL (a row recorded before the column existed) stays
	// distinguishable from a stored empty preimage; VerifiedCanonicalArgs relies on that difference.
	err := row.Scan(&a.ApprovalID, &a.RunID, &a.InstanceID, &a.AgentName, &a.RunActorSubject, &a.ToolName,
		&a.Connector, &a.ToolCallID, &a.ArgsHash, &a.CanonicalArgs, &a.Proposal, &a.Policy, &a.Status, &a.RequestedBy,
		&a.DecidedBy, &a.CreatedOn, &a.DecidedOn, &a.Position, &a.Total)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// inTx runs fn in a transaction, committing on success and rolling back on error or panic.
func (s *PostgresRunStore) inTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("act: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op after a successful Commit
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// nullIfEmpty maps "" to a SQL NULL so the nullable tenancy columns stay null rather than empty-string until an
// identity model populates them (§15.2).
func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
