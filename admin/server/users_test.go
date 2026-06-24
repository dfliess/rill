package server_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/testadmin"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUser(t *testing.T) {
	ctx := context.Background()
	fix := testadmin.New(t)

	t.Run("Deleting a user", func(t *testing.T) {
		// Create a superuser and two normal users
		_, sc1 := fix.NewSuperuser(t)
		u2, c2 := fix.NewUser(t)
		u3, c3 := fix.NewUser(t)

		// A normal user can't delete another normal user
		_, err := c2.DeleteUser(ctx, &adminv1.DeleteUserRequest{
			Email: u3.Email,
		})
		require.Error(t, err)
		require.Equal(t, codes.PermissionDenied, grpc.Code(err))

		// A normal user can delete themselves
		_, err = c2.DeleteUser(ctx, &adminv1.DeleteUserRequest{
			Email: u2.Email,
		})
		require.NoError(t, err)
		fix.Admin.PurgeAuthTokenCache()

		_, err = c2.GetCurrentUser(ctx, &adminv1.GetCurrentUserRequest{})
		require.Error(t, err)
		require.Equal(t, codes.Unauthenticated, grpc.Code(err))

		// A superuser can delete any user
		_, err = sc1.DeleteUser(ctx, &adminv1.DeleteUserRequest{
			Email:                u3.Email,
			SuperuserForceAccess: true,
		})
		require.NoError(t, err)
		fix.Admin.PurgeAuthTokenCache()

		_, err = c3.GetCurrentUser(ctx, &adminv1.GetCurrentUserRequest{})
		require.Error(t, err)
		require.Equal(t, codes.Unauthenticated, grpc.Code(err))
	})

	t.Run("Single-user orgs quota", func(t *testing.T) {
		u1, c1 := fix.NewUser(t)

		_, err := fix.Admin.DB.UpdateUser(ctx, u1.ID, &database.UpdateUserOptions{
			QuotaSingleuserOrgs: 3,
		})
		require.NoError(t, err)

		for i := 0; i < 4; i++ {
			orgName := "org" + strconv.Itoa(i)
			org, err := c1.CreateOrganization(ctx, &adminv1.CreateOrganizationRequest{
				Name: orgName,
			})
			if err != nil {
				require.Equal(t, codes.FailedPrecondition, status.Code(err), "error is: %v", err)
				require.ErrorContains(t, err, "quota exceeded")
				break
			}
			require.NoError(t, err)
			require.Equal(t, org.Organization.Name, orgName)
		}
		resp, err := c1.ListOrganizations(ctx, &adminv1.ListOrganizationsRequest{})
		require.NoError(t, err)
		require.Equal(t, 3, len(resp.Organizations))
	})

	t.Run("Preference language valid BCP-47", func(t *testing.T) {
		_, c1 := fix.NewUser(t)

		// Set a valid BCP-47 language tag
		resp, err := c1.UpdateUserPreferences(ctx, &adminv1.UpdateUserPreferencesRequest{
			Preferences: &adminv1.UserPreferences{
				PreferredLocale: strPtr("es"),
			},
		})
		require.NoError(t, err)
		require.Equal(t, "es", *resp.Preferences.PreferredLocale)

		// Verify it persists via GetCurrentUser
		cur, err := c1.GetCurrentUser(ctx, &adminv1.GetCurrentUserRequest{})
		require.NoError(t, err)
		require.Equal(t, "es", *cur.Preferences.PreferredLocale)

		// Set a more complex BCP-47 tag
		resp, err = c1.UpdateUserPreferences(ctx, &adminv1.UpdateUserPreferencesRequest{
			Preferences: &adminv1.UserPreferences{
				PreferredLocale: strPtr("pt-BR"),
			},
		})
		require.NoError(t, err)
		require.Equal(t, "pt-BR", *resp.Preferences.PreferredLocale)

		// Set language to empty string (reset)
		resp, err = c1.UpdateUserPreferences(ctx, &adminv1.UpdateUserPreferencesRequest{
			Preferences: &adminv1.UserPreferences{
				PreferredLocale: strPtr(""),
			},
		})
		require.NoError(t, err)
		require.Equal(t, "", *resp.Preferences.PreferredLocale)
	})

	t.Run("Preference language invalid tag", func(t *testing.T) {
		_, c1 := fix.NewUser(t)

		// Try a malformed language tag (single char is not valid BCP-47)
		_, err := c1.UpdateUserPreferences(ctx, &adminv1.UpdateUserPreferencesRequest{
			Preferences: &adminv1.UserPreferences{
				PreferredLocale: strPtr("a"),
			},
		})
		require.Error(t, err)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("Preference language returned in GetCurrentUser", func(t *testing.T) {
		_, c1 := fix.NewUser(t)

		// Initially empty
		cur, err := c1.GetCurrentUser(ctx, &adminv1.GetCurrentUserRequest{})
		require.NoError(t, err)
		require.NotNil(t, cur.Preferences)
		require.NotNil(t, cur.Preferences.PreferredLocale)
		require.Equal(t, "", *cur.Preferences.PreferredLocale)

		// Set a language
		_, err = c1.UpdateUserPreferences(ctx, &adminv1.UpdateUserPreferencesRequest{
			Preferences: &adminv1.UserPreferences{
				PreferredLocale: strPtr("en-US"),
			},
		})
		require.NoError(t, err)

		// Read it back
		cur, err = c1.GetCurrentUser(ctx, &adminv1.GetCurrentUserRequest{})
		require.NoError(t, err)
		require.Equal(t, "en-US", *cur.Preferences.PreferredLocale)
	})

	t.Run("Token basics", func(t *testing.T) {
		u1, c1 := fix.NewUser(t)

		// Issue a plain token
		res, err := c1.IssueUserAuthToken(ctx, &adminv1.IssueUserAuthTokenRequest{
			UserId:      "current",
			ClientId:    database.AuthClientIDRillManual,
			DisplayName: "Foo",
		})
		require.NoError(t, err)
		require.NotEmpty(t, res.Token)

		// Check the token works
		uTmp := fix.NewClient(t, res.Token)
		res2, err := uTmp.GetCurrentUser(ctx, &adminv1.GetCurrentUserRequest{})
		require.NoError(t, err)
		require.Equal(t, res2.User.Email, u1.Email)

		// Issue a token with an expiration
		res3, err := c1.IssueUserAuthToken(ctx, &adminv1.IssueUserAuthTokenRequest{
			UserId:     "current",
			ClientId:   database.AuthClientIDRillManual,
			TtlMinutes: 10,
		})
		require.NoError(t, err)
		require.NotEmpty(t, res3.Token)

		// Check the token were created
		res4, err := c1.ListUserAuthTokens(ctx, &adminv1.ListUserAuthTokensRequest{
			UserId: "current",
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(res4.Tokens), 3) // 2 created above and 1 from fix.NewUser

		// One should have description "Foo" and one should have an expiration
		var foundFoo, foundExpiration bool
		for _, token := range res4.Tokens {
			if token.DisplayName == "Foo" {
				foundFoo = true
			}
			if token.ExpiresOn != nil {
				foundExpiration = true
			}
		}
		require.True(t, foundFoo)
		require.True(t, foundExpiration)

		// Find an ID for the "Foo" token
		var tokenID string
		for _, token := range res4.Tokens {
			if token.DisplayName == "Foo" {
				tokenID = token.Id
				break
			}
		}
		require.NotEmpty(t, tokenID)

		// Revoke the token
		_, err = c1.RevokeUserAuthToken(ctx, &adminv1.RevokeUserAuthTokenRequest{
			TokenId: tokenID,
		})
		require.NoError(t, err)

		// Check the token is revoked
		res5, err := c1.ListUserAuthTokens(ctx, &adminv1.ListUserAuthTokensRequest{
			UserId: "current",
		})
		require.NoError(t, err)
		require.Equal(t, 2, len(res5.Tokens))
		for _, token := range res5.Tokens {
			require.NotEqual(t, token.Id, tokenID)
		}

	})
}

func TestCreateOrUpdateUserLocale(t *testing.T) {
	ctx := context.Background()
	fix := testadmin.New(t)

	t.Run("Locale seeds preference_language for new user", func(t *testing.T) {
		email := "locale-new@test.io"
		u, err := fix.Admin.CreateOrUpdateUser(ctx, email, "Locale Test", "", "es")
		require.NoError(t, err)
		require.Equal(t, "es", u.PreferenceLanguage)
	})

	t.Run("Locale does not overwrite existing preference_language", func(t *testing.T) {
		email := "locale-nooverwrite@test.io"

		// Create user with locale "fr"
		u, err := fix.Admin.CreateOrUpdateUser(ctx, email, "Test", "", "fr")
		require.NoError(t, err)
		require.Equal(t, "fr", u.PreferenceLanguage)

		// Call again with different locale "de" — should NOT overwrite
		u, err = fix.Admin.CreateOrUpdateUser(ctx, email, "Test", "", "de")
		require.NoError(t, err)
		require.Equal(t, "fr", u.PreferenceLanguage)
	})

	t.Run("Locale seeds empty preference_language on update", func(t *testing.T) {
		email := "locale-seedupdate@test.io"

		// Create user with no locale
		u, err := fix.Admin.CreateOrUpdateUser(ctx, email, "Test", "", "")
		require.NoError(t, err)
		require.Equal(t, "", u.PreferenceLanguage)

		// Update with locale — should seed because currently empty
		u, err = fix.Admin.CreateOrUpdateUser(ctx, email, "Test", "", "ja")
		require.NoError(t, err)
		require.Equal(t, "ja", u.PreferenceLanguage)
	})

	t.Run("Invalid locale is ignored", func(t *testing.T) {
		email := "locale-invalid@test.io"

		// Create user with malformed locale tag — should be silently ignored
		u, err := fix.Admin.CreateOrUpdateUser(ctx, email, "Test", "", "123")
		require.NoError(t, err)
		require.Equal(t, "", u.PreferenceLanguage)
	})
}

func strPtr(s string) *string {
	return &s
}
