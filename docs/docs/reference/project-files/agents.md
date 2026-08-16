---
note: GENERATED. DO NOT EDIT.
title: Agent YAML
sidebar_position: 42
---

Agents are autonomous AI workers (Kairos Act). An agent couples an LLM, a system prompt (`instructions`), an allowlist of built-in read-only analytical tools, and optional outbound MCP connectors whose tools are the agent's actions on the outside world. Actions are governed by a per-connector approval posture: by default, every MCP tool call requires a human approval before it executes.

An agent runs when started manually or by an agent trigger. Triggers can be declared inline in the agent's `triggers:` list or as standalone `agent_trigger` resources; both forms behave identically.


## Properties

### `type`

_[string]_ - Refers to the resource type and must be `agent` _(required)_

### `display_name`

_[string]_ - Display name for the agent shown in the UI

### `description`

_[string]_ - Description for the agent

### `model`

_[object]_ - Declares the LLM intended to run the agent. NOTE: runs currently execute on the runtime's configured AI connector; these fields carry the declared intent for per-agent model routing.

  - **`connector`** - _[string]_ - Name of the AI connector for the agent's model. If it names a connector resource in the project, the agent depends on it; AI connectors are commonly configured in rill.yaml instead.

  - **`name`** - _[string]_ - Model identifier to request on the connector

### `instructions`

_[string]_ - The agent's system prompt. Required; an agent without instructions is a parse error. It is combined with the project's `ai_instructions` at runtime. _(required)_

### `prompt`

_[string]_ - The message sent to the agent when a run is started without one, so launching it by hand needs no typing. Optional. `instructions` say what the agent does and travel as the system prompt; this is the message that asks for it, the same role a trigger's `input.prompt` plays for an automatic run. When absent, a manual run falls back to the prompt of the agent's single trigger; an agent with several triggers, or with none, needs either this property or a prompt supplied at launch.

### `tools`

_[array of string]_ - Built-in tools the agent may call, drawn from the fixed allowlist of read-only analytical tools. This is the upper bound on the agent's authority over Rill's own tools; the model never sees a built-in tool outside this list. These tools read governed data and never require approval. MCP tools are NOT declared here; they come from `mcp` connectors and are discovered at run start.

### `mcp`

_[array of object]_ - Outbound MCP connectors; remote MCP servers whose tools become the agent's actions on the outside world. Each entry declares one server and the approval posture applied to its tools. Tools are offered to the model under namespaced names (`mcp.<connector>.<tool>`), but the approval globs below match the raw tool name WITHOUT the `mcp.<connector>.` prefix.

  - **`name`** - _[string]_ - Local connector name; required and unique among the agent's connectors. Must not contain dots; the dot is the namespace separator in effective tool names (`mcp.<connector>.<tool>`), so a dotted name would misroute tool calls and is rejected at parse time. _(required)_

  - **`url`** - _[string]_ - The remote MCP endpoint. Required; must be an absolute `https` URL (validated at reconcile time). Only the Streamable HTTP transport is supported. _(required)_

  - **`auth`** - _[object]_ - Credential presented to the remote server as a bearer token.

    - **`secret`** - _[string]_ - The NAME of a secret in the platform secret manager, never the secret value itself. It is resolved server-side on every execution, so rotating the credential does not invalidate a pending approval.

  - **`network`** - _[object]_ - Egress guard for the connector.

    - **`allowed_hosts`** - _[array of string]_ - Allowlist of hostnames the connector may reach. Every request host (including redirect targets) must match one entry exactly, case-insensitively. When empty, no host allowlist is enforced, but requests to loopback/private ranges are always denied.

  - **`approval`** - _[string]_ - Approval posture for the connector's tools. `manual` (the default) means every tool call requires a human approval, except tools matching `auto_approve`. `auto` means no tool call requires approval, except tools matching `require_approval`. An absent or unrecognized posture is treated as `manual`; the posture fails toward requiring approval.

  - **`require_approval`** - _[array of string]_ - Glob patterns (`*`, `?`, `[set]`) of tools that still require approval when `approval` is `auto`. Patterns match the RAW tool name as advertised by the MCP server, without the `mcp.<connector>.` prefix (e.g. `create_*`, not `mcp.crm.create_*`). A malformed pattern is a parse error, not a silent misconfiguration.

  - **`auto_approve`** - _[array of string]_ - Glob patterns (`*`, `?`, `[set]`) of tools that skip approval when `approval` is `manual` (the default). Patterns match the RAW tool name as advertised by the MCP server, without the `mcp.<connector>.` prefix (e.g. `search_*`, not `mcp.crm.search_*`). A malformed pattern is a parse error, not a silent misconfiguration.

  - **`trust_read_only_hint`** - _[boolean]_ - Off by default. When true, the connector trusts the `readOnlyHint` annotation the remote MCP server advertises on its tools, letting a read-only-hinted tool execute without approval. The hint is the server's own claim, not something the platform verifies; enabling this delegates that decision to the remote server.

### `triggers`

_[array of object]_ - Inline declaration of when this agent runs. Each entry is an agent trigger body without `agent` (which is this agent); the parser desugars it into a standalone `agent_trigger` resource. The triggers live outside the agent's spec on purpose, so editing when the agent fires never changes its spec hash (what it does).

  - **`source`** - _[object]_ - What activates the trigger. `kind` selects the event source; which of the other fields apply depends on it. _(required)_

    - **`kind`** - _[string]_ - Kind of source the trigger subscribes to. Required. `alert` and `report` fire on events emitted by those resources; `schedule` fires on its own cron clock. _(required)_

    - **`name`** - _[string]_ - Name of the alert or report to subscribe to. Empty means ANY resource of that kind (wildcard); when set, the resource must exist, so a typo fails closed at reconcile instead of silently never firing. Must not be set for `kind: schedule`.

    - **`events`** - _[array of string]_ - Event short-names to subscribe to. Required (at least one) for `alert` and `report` sources; must not be set for `schedule`. Valid events are `entered_fail`, `recovered`, `entered_error`, `renotify_due` and `evaluated` for an alert source, and `completed` and `failed` for a report source.

    - **`cron`** - _[string]_ - Cron expression the trigger fires on (standard cron syntax, validated at reconcile time). Required for `kind: schedule`; must not be set for `alert` or `report` sources.

  - **`actor`** - _[object]_ - The identity a started run acts as; its security claims when reading governed data. Set at most ONE of `attributes` or `user_id`; combining them is a parse error. Naming the actor by email is deliberately not supported (an address can change or leave the org, silently breaking the trigger). With no actor, the run carries no claims and fails closed.

    - **`attributes`** - _[object]_ - Explicit user attributes the run acts as, mirroring an alert's `query_for_attributes`. Mutually exclusive with `user_id`.

    - **`user_id`** - _[string]_ - ID of the user whose attributes the run acts as; resolved via the admin service at reconcile time, exactly like an alert's `for.user_id`. Stable across email changes, unlike an email address. Mutually exclusive with `attributes`.

  - **`input`** - _[object]_ - The task handed to the started run.

    - **`prompt`** - _[string]_ - Prompt for the run. It specializes the task; when empty, the run starts from the agent's `instructions` with the triggering event as context.

    - **`context`** - _[object]_ - Key-value pairs injected as additional context for the run.

  - **`deduplication`** - _[object]_ - Collapses repeated activations of the trigger.

    - **`window`** - _[string]_ - Time window during which further matching events collapse onto the existing run instead of starting a new one. Accepts seconds, a Go duration (e.g. `24h`) or an ISO 8601 duration; must not be negative. Sub-second values round up, so `500ms` does not truncate to zero.

### `limits`

_[object]_ - Execution limits for a run of the agent.

  - **`max_steps`** - _[integer]_ - Maximum number of model/tool iterations in the agent loop. Zero or absent means the default limit.

  - **`timeout`** - _[string]_ - Wall-clock timeout for a run. Accepts seconds, a Go duration (e.g. `10m`) or an ISO 8601 duration. Sub-second values round up, so `500ms` does not truncate to zero. Zero or absent means no agent-level timeout.

### `security`

_[object]_ - Access control for the agent, resolved by Rill's security engine (ADR-0018). Only `access`, `launch` and `approve` apply; the dashboard policy's `row_filter`/`include`/`exclude`/`rules` are rejected on agents (an agent has no rows or fields), and `execute` is reserved and rejected. Expressions are templated booleans with the same grammar and user attributes as a dashboard's `security.access`. Without a `security` block, only admins can see the agent; declaring `security` without `access` denies it to everyone except admins, who always retain access.

  - **`access`** - _[oneOf]_ - Expression indicating whether the user can see the agent, including all of the agent's runs. If `security` is defined but `access` is not, it resolves to false and the agent is visible only to admins, who are never excluded by an agent's policy (unlike a metrics view's, which applies to them too). It cannot reference the `.action`/`.run` context (that exists only while deciding an approval).

    - **option 1** - _[string]_ - SQL expression that evaluates to a boolean to determine access

    - **option 2** - _[boolean]_ - Direct boolean value to allow or deny access

  - **`launch`** - _[oneOf]_ - Expression indicating whether the user can start a run of the agent. When omitted, it inherits `access`. Launching is always a subset of `access`. It cannot reference the `.action`/`.run` context (that exists only while deciding an approval).

    - **option 1** - _[string]_ - SQL expression that evaluates to a boolean to determine who can launch

    - **option 2** - _[boolean]_ - Direct boolean value to allow or deny launching

  - **`approve`** - _[oneOf]_ - Expression indicating whether the user can approve or reject an action that the connector's approval posture flagged for a signature. When omitted, only admins can approve (deliberately stricter than `launch`; launching is inert, approving touches the outside world). Beyond `.user`, the expression may reference only `.action.tool`, `.action.connector` and `.run.actor`; the action's arguments are model-chosen and can never select the approver.

    - **option 1** - _[string]_ - SQL expression that evaluates to a boolean to determine who can approve

    - **option 2** - _[boolean]_ - Direct boolean value to allow or deny approving

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
# Example: A collections agent that investigates a delinquency alert and proposes a CRM task
type: agent
display_name: Collections Management Agent
description: Investigates delinquency spikes and proposes opening a task in the CRM.
model:
    connector: deepseek
    name: deepseek-v4-flash
instructions: |
    You are a credit and collections analyst. Investigate using the governed
    metrics tools, and only at the end propose the task with create_crm_task.
tools:
    - list_metrics_views
    - get_metrics_view
    - query_metrics_view
    - query_metrics_view_summary
mcp:
    - name: crm
      url: https://crm.example.com/mcp
      auth:
        secret: CRM_TOKEN
      network:
        allowed_hosts:
            - crm.example.com
      approval: manual
      auto_approve:
        - "search_*"
        - "get_*"
triggers:
    - source:
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
limits:
    max_steps: 20
    timeout: 10m
```

```yaml
# Example: An agent launched by hand, with no trigger. `prompt` is the request a trigger would
# otherwise carry, so starting a run needs nothing typed.
type: agent
display_name: Returns Reviewer
description: Reviews last month's returns on demand and proposes a quality ticket.
model:
    connector: deepseek
    name: deepseek-v4-flash
instructions: |
    You are a returns analyst. Use the governed metrics tools to find which
    category and channel drive returns, and only at the end propose the ticket.
prompt: Review last month's returns and propose the quality ticket.
tools:
    - query_metrics_view
    - get_metrics_view
mcp:
    - name: quality
      url: https://quality.example.com/mcp
      auth:
        secret: QUALITY_TOKEN
      network:
        allowed_hosts:
            - quality.example.com
      approval: manual
limits:
    max_steps: 20
    timeout: 10m
```
