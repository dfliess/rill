package postgres

import (
	"context"
	"errors"

	"github.com/rilldata/rill/admin/database"
)

// Kairos Styles (kairos-cloud#242); the table is migrations_kairos/0004.sql.

func (c *connection) FindOrganizationStyle(ctx context.Context, orgID string) (string, error) {
	var style string
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT style FROM org_styles WHERE org_id=$1", orgID).Scan(&style)
	if err != nil {
		if errors.Is(parseErr("organization style", err), database.ErrNotFound) {
			return "", nil
		}
		return "", parseErr("organization style", err)
	}
	return style, nil
}

func (c *connection) UpdateOrganizationStyle(ctx context.Context, orgID, style string) error {
	if style == "" {
		_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM org_styles WHERE org_id=$1", orgID)
		return parseErr("organization style", err)
	}
	_, err := c.getDB(ctx).ExecContext(ctx, `
		INSERT INTO org_styles (org_id, style) VALUES ($1, $2)
		ON CONFLICT (org_id) DO UPDATE SET style=excluded.style, updated_on=now()`,
		orgID, style,
	)
	return parseErr("organization style", err)
}
