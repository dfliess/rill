---
note: GENERATED. DO NOT EDIT.
title: Agent Trigger YAML
sidebar_position: 43
---

Agent triggers declare when an agent runs on its own; on an alert transition, on a report completion, or on a cron schedule. A trigger only starts runs; what the agent does is defined on the agent itself. The same body can be declared inline in the agent's `triggers:` list, where `agent` is implicit; both forms behave identically.

## Properties

### `type`

_[string]_ - Refers to the resource type and must be `agent_trigger` _(required)_

### `agent`

_[string]_ - Name of the agent resource this trigger starts. Required; the trigger reconciles after the agent and fails closed if the agent is missing or invalid. _(required)_

### `source`

_[object]_ - What activates the trigger. `kind` selects the event source; which of the other fields apply depends on it. _(required)_

  - **`kind`** - _[string]_ - Kind of source the trigger subscribes to. Required. `alert` and `report` fire on events emitted by those resources; `schedule` fires on its own cron clock. _(required)_

  - **`name`** - _[string]_ - Name of the alert or report to subscribe to. Empty means ANY resource of that kind (wildcard); when set, the resource must exist, so a typo fails closed at reconcile instead of silently never firing. Must not be set for `kind: schedule`.

  - **`events`** - _[array of string]_ - Event short-names to subscribe to. Required (at least one) for `alert` and `report` sources; must not be set for `schedule`. Valid events are `entered_fail`, `recovered`, `entered_error`, `renotify_due` and `evaluated` for an alert source, and `completed` and `failed` for a report source.

  - **`cron`** - _[string]_ - Cron expression the trigger fires on (standard cron syntax, validated at reconcile time). Required for `kind: schedule`; must not be set for `alert` or `report` sources.

### `actor`

_[object]_ - The identity a started run acts as; its security claims when reading governed data. Set at most ONE of `attributes` or `user_id`; combining them is a parse error. Naming the actor by email is deliberately not supported (an address can change or leave the org, silently breaking the trigger). With no actor, the run carries no claims and fails closed.

  - **`attributes`** - _[object]_ - Explicit user attributes the run acts as, mirroring an alert's `query_for_attributes`. Mutually exclusive with `user_id`.

  - **`user_id`** - _[string]_ - ID of the user whose attributes the run acts as; resolved via the admin service at reconcile time, exactly like an alert's `for.user_id`. Stable across email changes, unlike an email address. Mutually exclusive with `attributes`.

### `input`

_[object]_ - The task handed to the started run.

  - **`prompt`** - _[string]_ - Prompt for the run. It specializes the task; when empty, the run starts from the agent's `instructions` with the triggering event as context.

  - **`context`** - _[object]_ - Key-value pairs injected as additional context for the run.

### `deduplication`

_[object]_ - Collapses repeated activations of the trigger.

  - **`window`** - _[string]_ - Time window during which further matching events collapse onto the existing run instead of starting a new one. Accepts seconds, a Go duration (e.g. `24h`) or an ISO 8601 duration; must not be negative. Sub-second values round up, so `500ms` does not truncate to zero.

## Common Properties

### `name`

_[string]_ - Name is usually inferred from the filename, but can be specified manually.

### `refs`

_[array of string]_ - List of resource references

### `tags`

_[array of string]_ - Tags for organizing and filtering the resource (e.g. on the project dashboards list).

### `dev`

_[object]_ - Overrides any properties in development environment.

### `prod`

_[object]_ - Overrides any properties in production environment.

## Examples

```yaml
# Example: Start the collections agent whenever the delinquency alert starts failing
type: agent_trigger
agent: collections_agent
source:
    kind: alert
    name: high-delinquency
    events:
        - entered_fail
actor:
    attributes:
        email: act@example.com
input:
    prompt: Investigate the delinquency spike and propose the collections task.
deduplication:
    window: 24h
```

```yaml
# Example: Run the weekly review agent every Monday at 08:00
type: agent_trigger
agent: weekly_review_agent
source:
    kind: schedule
    cron: 0 8 * * 1
actor:
    attributes:
        email: act@example.com
```
