CREATE TABLE push_subscriptions (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	endpoint TEXT NOT NULL,
	p256dh TEXT NOT NULL,
	auth TEXT NOT NULL,
	user_agent TEXT NOT NULL DEFAULT '',
	created_on TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX push_subscriptions_endpoint_idx ON push_subscriptions (endpoint);
CREATE INDEX push_subscriptions_user_id_idx ON push_subscriptions (user_id);

CREATE TABLE notification_preferences (
	user_id UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
	push_alerts BOOLEAN NOT NULL DEFAULT true,
	push_reports BOOLEAN NOT NULL DEFAULT true,
	push_act_approvals BOOLEAN NOT NULL DEFAULT true,
	updated_on TIMESTAMPTZ NOT NULL DEFAULT now()
);
