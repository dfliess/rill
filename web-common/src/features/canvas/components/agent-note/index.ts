import { BaseCanvasComponent } from "@rilldata/web-common/features/canvas/components/BaseCanvasComponent";
import type { InputParams } from "@rilldata/web-common/features/canvas/inspector/types";
import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
import type { V1Resource } from "@rilldata/web-common/runtime-client";
import type { CanvasEntity, ComponentPath } from "../../stores/canvas-entity";
import type { CanvasComponentType, ComponentCommonProperties } from "../types";
import AgentNote from "./AgentNote.svelte";

export { default as AgentNote } from "./AgentNote.svelte";

export interface AgentNoteSpec extends ComponentCommonProperties {
  /** Name of a `type: agent` resource in the project. It must be read-only (no MCP block). */
  agent: string;
  /** What the agent is asked to do with the dashboard it is placed on. */
  prompt: string;
  /**
   * Generate on mount. Repeat visits with unchanged filters attach to the existing run, so leaving this
   * on costs one completion per context rather than one per page load. Turn it off to require a click.
   */
  auto_run?: boolean;
  /**
   * Lines shown before the note is clamped (default 3). The canvas cannot size a row to its content, so
   * the row that hosts the note declares how many lines it has room for.
   */
  lines?: number;
}

export class AgentNoteCanvasComponent extends BaseCanvasComponent<AgentNoteSpec> {
  minSize = { width: 2, height: 1 };
  defaultSize = { width: 6, height: 2 };
  resetParams = [];
  type: CanvasComponentType = "agent_note";
  component = AgentNote;

  constructor(resource: V1Resource, parent: CanvasEntity, path: ComponentPath) {
    const defaultSpec: AgentNoteSpec = {
      title: "",
      description: "",
      agent: "",
      prompt: "",
      auto_run: true,
    };
    super(resource, parent, path, defaultSpec);
  }

  isValid(spec: AgentNoteSpec): boolean {
    return (
      typeof spec.agent === "string" &&
      spec.agent.trim().length > 0 &&
      typeof spec.prompt === "string" &&
      spec.prompt.trim().length > 0
    );
  }

  inputParams(): InputParams<AgentNoteSpec> {
    return {
      options: {
        agent: {
          type: "text",
          label: m.canvas_agent_note_agent_label(),
          description: m.canvas_agent_note_agent_description(),
        },
        prompt: {
          type: "textarea",
          label: m.canvas_agent_note_prompt_label(),
          description: m.canvas_agent_note_prompt_description(),
        },
        auto_run: {
          type: "boolean",
          optional: true,
          showInUI: true,
          label: m.canvas_agent_note_auto_run_label(),
          description: m.canvas_agent_note_auto_run_description(),
        },
        lines: {
          type: "number",
          optional: true,
          showInUI: true,
          label: m.canvas_agent_note_lines_label(),
          description: m.canvas_agent_note_lines_description(),
        },
      },
      filter: {},
    };
  }

  static newComponentSpec(): AgentNoteSpec {
    return {
      agent: "",
      prompt: "Summarise what this dashboard shows and why it moved.",
      auto_run: true,
    };
  }
}
