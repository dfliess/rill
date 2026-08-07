package act

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// PostgresRunStore also implements ActionLedger: the ledger's agent_actions table lives in the same product schema as
// the run tables (created by Migrate), so one store and one pool serve both. Keeping them together lets a run's
// lifecycle transition and its action-ledger writes share the same Postgres without a second connection or schema.
var _ ActionLedger = (*PostgresRunStore)(nil)

func (s *PostgresRunStore) ProposeAction(ctx context.Context, a NewAction) error {
	if a.ToolCallID == "" || a.RunID == "" || a.InstanceID == "" {
		return errors.New("act: ProposeAction requires ToolCallID, RunID and InstanceID")
	}
	if a.ArgsHash == "" || a.IdempotencyKey == "" {
		return errors.New("act: ProposeAction requires ArgsHash and IdempotencyKey")
	}
	argsJSON, err := jsonbParam(a.RedactedArgs)
	if err != nil {
		return fmt.Errorf("act: marshal redacted args: %w", err)
	}
	// ON CONFLICT DO NOTHING on the (instance, run, tool call) key makes a replayed proposal a no-op rather than a
	// duplicate-key error: an action is proposed exactly once even under at-least-once step delivery (§12).
	tag, err := s.pool.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (instance_id, run_id, tool_call_id, agent_name, tool, connector, tool_version, class, args_hash,
			idempotency_key, redacted_args, policy_decision, policy_reason, proposal, status, requested_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11::jsonb,$12,$13,$14,$15,$16)
		ON CONFLICT (instance_id, run_id, tool_call_id) DO NOTHING`, s.t("agent_actions")),
		a.InstanceID, a.RunID, a.ToolCallID, a.AgentName, a.Tool, a.Connector, a.ToolVersion, string(a.Class),
		a.ArgsHash, a.IdempotencyKey, argsJSON, string(a.PolicyDecision), a.PolicyReason, a.Proposal,
		string(ActionProposed), a.RequestedBy)
	if err != nil {
		return fmt.Errorf("act: propose action: %w", err)
	}
	if tag.RowsAffected() == 0 {
		// The row already existed. A replay of the same proposal is benign, but two DIFFERENT proposals racing onto one
		// tool call must not pass silently: reload the winner and reject if its immutable identity differs from ours. A
		// Proposer that mutated the tool, connector, version or arguments cannot ride the earlier row (§12, §17.1).
		existing, gErr := s.GetAction(ctx, a.InstanceID, a.RunID, a.ToolCallID)
		if gErr != nil {
			return fmt.Errorf("act: reload conflicting action: %w", gErr)
		}
		if existing.Tool != a.Tool || existing.Connector != a.Connector || existing.ToolVersion != a.ToolVersion ||
			existing.Class != a.Class || existing.ArgsHash != a.ArgsHash || existing.IdempotencyKey != a.IdempotencyKey {
			return fmt.Errorf("act: proposal for %q conflicts with the recorded action: %w", a.ToolCallID, ErrActionTransitionInvalid)
		}
	}
	return nil
}

// ClaimExecuting flips an approved action to executing atomically, so exactly one caller ever owns the write (§12).
// The UPDATE is conditional on the row still being approved: its affected-row count is the claim result. A caller that
// wins (one row) may execute; one that loses (zero rows) must not, and reads the current state to report the recorded
// outcome. It stamps attempt_id (the owner) and lease_expiry = now()+lease on the database clock so recovery can tell
// a live attempt from a dead one. executed_on latches on the first move into executing so an interrupted attempt stays
// visible on recovery.
func (s *PostgresRunStore) ClaimExecuting(ctx context.Context, instanceID, runID, toolCallID, attemptID string, lease time.Duration) (bool, *Action, error) {
	if toolCallID == "" || runID == "" || instanceID == "" {
		return false, nil, errors.New("act: ClaimExecuting requires ToolCallID, RunID and InstanceID")
	}
	var (
		claimed bool
		action  *Action
	)
	err := s.inTx(ctx, func(tx pgx.Tx) error {
		// lease_expiry is computed from the database now() so it is comparable to the now() ReclaimExpiredLease reads; a
		// non-positive lease (used by tests to force an already-expired claim) yields an expiry in the past.
		tag, err := tx.Exec(ctx, fmt.Sprintf(`
			UPDATE %s SET
				status=$4,
				attempt_id=$5,
				lease_expiry=now() + make_interval(secs => $6::double precision),
				executed_on=COALESCE(executed_on, now()),
				updated_on=now()
			WHERE instance_id=$1 AND run_id=$2 AND tool_call_id=$3 AND status=$7`, s.t("agent_actions")),
			instanceID, runID, toolCallID, string(ActionExecuting), attemptID, lease.Seconds(), string(ActionApproved))
		if err != nil {
			return fmt.Errorf("claim executing: %w", err)
		}
		claimed = tag.RowsAffected() == 1
		action, err = scanAction(tx.QueryRow(ctx, s.selectActionSQL()+` WHERE instance_id=$1 AND run_id=$2 AND tool_call_id=$3`,
			instanceID, runID, toolCallID))
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrActionNotFound
		}
		if err != nil {
			return fmt.Errorf("reload action: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, nil, err
	}
	return claimed, action, nil
}

// ReclaimExpiredLease moves an executing action to indeterminate iff its lease has expired or was never set (§12). The
// WHERE clause is the whole guard: it matches only a row still executing whose lease_expiry is null (a crash that
// reached executing without a claim) or already in the past, so a row under a LIVE lease affects zero rows and is left
// for its owner to finalize. reclaimed is the affected-row count; the returned action is the current row either way.
func (s *PostgresRunStore) ReclaimExpiredLease(ctx context.Context, instanceID, runID, toolCallID, reason string) (bool, *Action, error) {
	if toolCallID == "" || runID == "" || instanceID == "" {
		return false, nil, errors.New("act: ReclaimExpiredLease requires ToolCallID, RunID and InstanceID")
	}
	var (
		reclaimed bool
		action    *Action
	)
	err := s.inTx(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, fmt.Sprintf(`
			UPDATE %s SET
				status=$4,
				error = CASE WHEN $5<>'' THEN $5 ELSE error END,
				updated_on=now()
			WHERE instance_id=$1 AND run_id=$2 AND tool_call_id=$3
				AND status=$6
				AND (lease_expiry IS NULL OR lease_expiry < now())`, s.t("agent_actions")),
			instanceID, runID, toolCallID, string(ActionIndeterminate), reason, string(ActionExecuting))
		if err != nil {
			return fmt.Errorf("reclaim expired lease: %w", err)
		}
		reclaimed = tag.RowsAffected() == 1
		action, err = scanAction(tx.QueryRow(ctx, s.selectActionSQL()+` WHERE instance_id=$1 AND run_id=$2 AND tool_call_id=$3`,
			instanceID, runID, toolCallID))
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrActionNotFound
		}
		if err != nil {
			return fmt.Errorf("reload action: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, nil, err
	}
	return reclaimed, action, nil
}

func (s *PostgresRunStore) TransitionAction(ctx context.Context, tr ActionTransition) (*Action, error) {
	if tr.ToolCallID == "" || tr.RunID == "" || tr.InstanceID == "" || tr.Status == "" {
		return nil, errors.New("act: TransitionAction requires ToolCallID, RunID, InstanceID and Status")
	}
	resultJSON, err := jsonbParam(tr.RedactedResult)
	if err != nil {
		return nil, fmt.Errorf("act: marshal redacted result: %w", err)
	}
	verificationJSON, err := jsonbParam(tr.RedactedVerification)
	if err != nil {
		return nil, fmt.Errorf("act: marshal redacted verification: %w", err)
	}

	var action *Action
	err = s.inTx(ctx, func(tx pgx.Tx) error {
		// Lock the action row: validates existence and instance/run scoping, and serializes concurrent transitions on
		// one action so the state machine is evaluated against a stable current status.
		var current string
		lockErr := tx.QueryRow(ctx, fmt.Sprintf(`SELECT status FROM %s WHERE instance_id=$1 AND run_id=$2 AND tool_call_id=$3 FOR UPDATE`, s.t("agent_actions")),
			tr.InstanceID, tr.RunID, tr.ToolCallID).Scan(&current)
		if errors.Is(lockErr, pgx.ErrNoRows) {
			return ErrActionNotFound
		}
		if lockErr != nil {
			return fmt.Errorf("lock action: %w", lockErr)
		}

		// Same-status transition is an idempotent no-op: a replayed step that re-applies the state the action already
		// holds must not error or double-write. A transition out of a terminal state is likewise a no-op: a terminal
		// action's outcome is final, so a stray replayed edge after it must not regress or overwrite it (mirrors the
		// run store's terminal guard). Only a genuine forward edge from a non-terminal state is applied.
		cur := ActionStatus(current)
		if cur != tr.Status && !cur.IsTerminal() {
			if !actionEdges[cur][tr.Status] {
				return fmt.Errorf("%w: %s -> %s", ErrActionTransitionInvalid, current, tr.Status)
			}
			// Apply the edge. Optional columns are written only when supplied: an empty external reference or error, or
			// a nil redacted payload, leaves the stored value untouched (COALESCE / CASE guards). executed_on latches on
			// the first move into executing; verified_on is set on verified. When AttemptID is set the update is
			// ownership-guarded: it applies only while the row is still executing under that exact attempt, so a finalize
			// whose attempt was reclaimed by recovery affects no row and the reloaded action reveals the recovered outcome
			// instead of overwriting it (§12).
			_, execErr := tx.Exec(ctx, fmt.Sprintf(`
				UPDATE %s SET
					status=$4,
					external_reference = CASE WHEN $5<>'' THEN $5 ELSE external_reference END,
					redacted_result = COALESCE($6::jsonb, redacted_result),
					redacted_verification = COALESCE($7::jsonb, redacted_verification),
					error = CASE WHEN $8<>'' THEN $8 ELSE error END,
					decided_by = CASE WHEN $9<>'' THEN $9 ELSE decided_by END,
					redacted_message = CASE WHEN $13<>'' THEN $13 ELSE redacted_message END,
					updated_on = now(),
					executed_on = COALESCE(executed_on, CASE WHEN $4=$10 THEN now() END),
					verified_on = CASE WHEN $4=$11 THEN now() ELSE verified_on END
				WHERE instance_id=$1 AND run_id=$2 AND tool_call_id=$3
					AND ($12='' OR attempt_id=$12)`, s.t("agent_actions")),
				tr.InstanceID, tr.RunID, tr.ToolCallID, string(tr.Status), tr.ExternalReference, resultJSON,
				verificationJSON, tr.Error, tr.DecidedBy, string(ActionExecuting), string(ActionVerified), tr.AttemptID,
				tr.RedactedMessage)
			if execErr != nil {
				return fmt.Errorf("update action: %w", execErr)
			}
		}

		action, err = scanAction(tx.QueryRow(ctx, s.selectActionSQL()+` WHERE instance_id=$1 AND run_id=$2 AND tool_call_id=$3`,
			tr.InstanceID, tr.RunID, tr.ToolCallID))
		if err != nil {
			return fmt.Errorf("reload action: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (s *PostgresRunStore) GetAction(ctx context.Context, instanceID, runID, toolCallID string) (*Action, error) {
	action, err := scanAction(s.pool.QueryRow(ctx, s.selectActionSQL()+` WHERE instance_id=$1 AND run_id=$2 AND tool_call_id=$3`,
		instanceID, runID, toolCallID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrActionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("act: get action: %w", err)
	}
	return action, nil
}

func (s *PostgresRunStore) ListActions(ctx context.Context, instanceID, runID string) ([]*Action, error) {
	if instanceID == "" || runID == "" {
		return nil, errors.New("act: ListActions requires InstanceID and RunID")
	}
	rows, err := s.pool.Query(ctx, s.selectActionSQL()+` WHERE instance_id=$1 AND run_id=$2 ORDER BY created_on DESC`, instanceID, runID)
	if err != nil {
		return nil, fmt.Errorf("act: list actions: %w", err)
	}
	defer rows.Close()

	var out []*Action
	for rows.Next() {
		action, scanErr := scanAction(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("act: scan action: %w", scanErr)
		}
		out = append(out, action)
	}
	return out, rows.Err()
}

func (s *PostgresRunStore) WithdrawPendingActions(ctx context.Context, instanceID, runID string) (int64, error) {
	tag, err := s.pool.Exec(ctx, fmt.Sprintf(
		`UPDATE %s SET status=$3, updated_on=now() WHERE instance_id=$1 AND run_id=$2 AND status=$4`,
		s.t("agent_actions")),
		instanceID, runID, string(ActionWithdrawn), string(ActionApprovalPending))
	if err != nil {
		return 0, fmt.Errorf("act: withdraw pending actions: %w", err)
	}
	return tag.RowsAffected(), nil
}

// selectActionSQL is the column list and source for an action projection, in the order scanAction expects.
func (s *PostgresRunStore) selectActionSQL() string {
	return fmt.Sprintf(`
		SELECT instance_id, run_id, tool_call_id, agent_name, tool, connector, tool_version, class, args_hash,
			idempotency_key, redacted_args, policy_decision, policy_reason, proposal, status, attempt_id,
			external_reference, redacted_result, redacted_message, redacted_verification, error, requested_by, decided_by,
			created_on, updated_on, executed_on, verified_on
		FROM %s`, s.t("agent_actions"))
}

func scanAction(row rowScanner) (*Action, error) {
	var (
		a                                      Action
		class, decision, status                string
		argsJSON, resultJSON, verificationJSON []byte
	)
	err := row.Scan(&a.InstanceID, &a.RunID, &a.ToolCallID, &a.AgentName, &a.Tool, &a.Connector, &a.ToolVersion, &class,
		&a.ArgsHash, &a.IdempotencyKey, &argsJSON, &decision, &a.PolicyReason, &a.Proposal, &status, &a.AttemptID,
		&a.ExternalReference, &resultJSON, &a.RedactedMessage, &verificationJSON, &a.Error, &a.RequestedBy, &a.DecidedBy,
		&a.CreatedOn, &a.UpdatedOn, &a.ExecutedOn, &a.VerifiedOn)
	if err != nil {
		return nil, err
	}
	a.Class = ToolClass(class)
	a.PolicyDecision = PolicyDecision(decision)
	a.Status = ActionStatus(status)
	if a.RedactedArgs, err = jsonbMap(argsJSON); err != nil {
		return nil, fmt.Errorf("unmarshal redacted args: %w", err)
	}
	if a.RedactedResult, err = jsonbMap(resultJSON); err != nil {
		return nil, fmt.Errorf("unmarshal redacted result: %w", err)
	}
	if a.RedactedVerification, err = jsonbMap(verificationJSON); err != nil {
		return nil, fmt.Errorf("unmarshal redacted verification: %w", err)
	}
	return &a, nil
}

// jsonbParam renders a map for a jsonb parameter: a nil map becomes a SQL NULL (so COALESCE keeps an existing value),
// otherwise the JSON text (passed as a string, which pgx encodes to jsonb). It never carries resolved secrets: the
// gateway redacts before persisting (§16.2).
func jsonbParam(m map[string]any) (any, error) {
	if m == nil {
		return nil, nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// jsonbMap unmarshals a jsonb column into a map, mapping SQL NULL / empty to a nil map.
func jsonbMap(b []byte) (map[string]any, error) {
	if len(b) == 0 {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}
