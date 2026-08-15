import type { PartialMessage } from "@bufbuild/protobuf";
import type {
  AgentApproval,
  AgentRun,
} from "@rilldata/web-common/proto/gen/rill/runtime/v1/agents_pb";

// The AgentService hooks return proto messages serialized with `toJson`, so the
// data these components consume is the JSON projection of each message. We model
// that as `PartialMessage<T>` to match what the generated hooks hand back.
export type AgentRunData = PartialMessage<AgentRun>;
export type AgentApprovalData = PartialMessage<AgentApproval>;

/**
 * A run's event normalized to a plain shape for display. Populated from the
 * StreamAgentRunEvents server stream (see agent-run-events.ts); int64 cursors
 * are surfaced as strings/numbers and the google.protobuf.Struct payload as a
 * plain object.
 */
export interface AgentRunTimelineEvent {
  id: string;
  seq: number;
  eventType: string;
  status: string;
  visibility: string;
  payload: unknown;
  createdOn?: string;
}
