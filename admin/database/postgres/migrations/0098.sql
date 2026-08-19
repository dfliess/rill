-- Notification preferences become per-organization: a user may want a category from one org and not from another.
-- The table has never been deployed, so its rows are dropped instead of migrated.
DROP TABLE notification_preferences;

CREATE TABLE notification_preferences (
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	org_id UUID NOT NULL REFERENCES orgs (id) ON DELETE CASCADE,
	push_alerts BOOLEAN NOT NULL DEFAULT true,
	push_reports BOOLEAN NOT NULL DEFAULT true,
	push_act_approvals BOOLEAN NOT NULL DEFAULT true,
	updated_on TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, org_id)
);

CREATE INDEX notification_preferences_org_id_idx ON notification_preferences (org_id);
