-- What a person has been notified about, so a notification that appeared and went leaves a record they can go
-- back to (kairos-cloud#148). One row per recipient, not per device: the history is what concerns you, not what
-- your browser happened to receive, so it exists for someone with no browser subscribed and for someone who
-- turned the category off (they turned off the interruption, not the record).
CREATE TABLE notification_deliveries (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	-- The project the notification came from. Cascading is deliberate: a deleted project takes its notifications
	-- with it, like the rest of its data, rather than leaving rows pointing at something that no longer exists.
	project_id UUID NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
	category TEXT NOT NULL,
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	-- Path relative to the frontend, not an absolute URL: a deployment that changes domain keeps its history usable.
	link_path TEXT NOT NULL,
	sent_on TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- The only read is "this person's history, newest first", and the only delete is by age.
CREATE INDEX notification_deliveries_user_id_sent_on_idx ON notification_deliveries (user_id, sent_on DESC);
CREATE INDEX notification_deliveries_sent_on_idx ON notification_deliveries (sent_on);
