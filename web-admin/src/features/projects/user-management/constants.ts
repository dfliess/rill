import * as m from "@rilldata/web-common/paraglide/messages.js";
import { ProjectUserRoles } from "@rilldata/web-common/features/users/roles.ts";

export function getProjectRolesOptions() {
  return [
    {
      value: ProjectUserRoles.Admin,
      label: m.role_admin(),
      description: m.role_project_admin_desc(),
    },
    {
      value: ProjectUserRoles.Editor,
      label: m.role_editor(),
      description: m.role_project_editor_desc(),
    },
    {
      value: ProjectUserRoles.Viewer,
      label: m.role_viewer(),
      description: m.role_project_viewer_desc(),
    },
  ];
}

export function getProjectRoleDescription(role: string): string {
  switch (role) {
    case "admin": return m.role_project_admin_desc();
    case "editor": return m.role_project_editor_desc();
    case "viewer": return m.role_project_viewer_desc();
    case "guest": return m.role_guest_desc();
    default: return "";
  }
}
