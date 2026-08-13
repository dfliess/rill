package adminenv

import (
	"fmt"
)

var EnvURLs = map[string]string{
	"prod":  "https://admin.rilldata.com",
	"stage": "https://admin.rilldata.io",
	"test":  "https://admin.rilldata.in",
	"dev":   "http://localhost:8080",

	// Kairos self-hosted environments (ADR-0012). `rill devtool switch-env` parks the active token under
	// tokens.<env> and swaps in the target's, which is the only mechanism that keeps a dev login from
	// overwriting the prod one; Infer matches the admin URL exactly, so a self-hosted deployment needs its
	// URLs listed here or the command fails before it can park anything.
	//
	// kairos-prod is the local end of the SSH tunnel, not a public host: the admin is not published to the
	// internet (the public host answers 405 for gRPC), so the CLI always reaches prod through
	// `ssh -L 28080:<admin-ip>:20080 kairos`.
	"kairos-dev":  "http://localhost:20080",
	"kairos-prod": "http://localhost:28080",
}

func Infer(adminURL string) (string, error) {
	for env, url := range EnvURLs {
		if url == adminURL {
			return env, nil
		}
	}
	return "", fmt.Errorf("could not infer env from admin URL %q", adminURL)
}

func AdminURL(env string) string {
	u, ok := EnvURLs[env]
	if !ok {
		panic(fmt.Errorf("invalid environment %q", env))
	}
	return u
}
