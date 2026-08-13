package adminenv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The Kairos deployment is self-hosted, so its admin URLs have to be listed for Infer to resolve them: without an
// entry, `rill devtool switch-env` fails before it can park the active token and a dev login overwrites the prod one.
func TestInferResolvesKairosEnvironments(t *testing.T) {
	env, err := Infer("http://localhost:20080")
	require.NoError(t, err)
	require.Equal(t, "kairos-dev", env)

	// Prod is reached through the SSH tunnel; the admin is not published to the internet.
	env, err = Infer("http://localhost:28080")
	require.NoError(t, err)
	require.Equal(t, "kairos-prod", env)
}

func TestInferRejectsUnknownURL(t *testing.T) {
	_, err := Infer("https://example.com")
	require.Error(t, err)
}

// Infer and AdminURL read the same map from opposite ends, so a duplicated or mistyped URL would make an environment
// unreachable in one direction only. Round-tripping every entry catches that at the point the map is edited.
func TestEnvURLsRoundTrip(t *testing.T) {
	for env, url := range EnvURLs {
		require.Equal(t, url, AdminURL(env), "AdminURL(%q)", env)

		inferred, err := Infer(url)
		require.NoError(t, err, "Infer(%q)", url)
		require.Equal(t, env, inferred, "%q and %q share the URL %q", env, inferred, url)
	}
}
