-- An organization's Kairos Style (kairos-cloud#242): the app's look its members get unless they pick their own.
-- A table of its own rather than a column on orgs: upstream rewrites every orgs column on each update, from over a
-- dozen call sites, so a Kairos column there would be reset by whichever one forgot to copy it.
-- No row means the default Style.
CREATE TABLE org_styles (
	org_id UUID PRIMARY KEY REFERENCES orgs (id) ON DELETE CASCADE,
	style TEXT NOT NULL,
	updated_on TIMESTAMPTZ NOT NULL DEFAULT now()
);
