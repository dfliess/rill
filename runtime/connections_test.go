package runtime_test

import (
	"context"
	"testing"

	"github.com/rilldata/rill/runtime"
	_ "github.com/rilldata/rill/runtime/drivers/admin"
	_ "github.com/rilldata/rill/runtime/drivers/s3"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAcquireHandle(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			`rill.yaml`: `
display_name: Hello world
description: This project says hello to the world

connectors:
- name: my-s3
  type: s3
  defaults:
    AWS_ACCESS_KEY_ID: us-east-1
    aws_secret_access_key: xxxx

vars:
  foo: bar
  allow_host_access: false
`,
		},
		Variables: map[string]string{
			"aws_secret_access_key": "yyyy",
		},
	})
	ctx := context.Background()

	handle, _, err := rt.AcquireHandle(ctx, id, "my-s3")
	require.NoError(t, err)
	config := handle.Config()
	require.True(t, config["aws_access_key_id"].(string) == "us-east-1")
	require.True(t, config["allow_host_access"].(bool))
}

func TestNonSecretConnectorPropertiesSupportsManagedAdminAI(t *testing.T) {
	properties, err := runtime.NonSecretConnectorProperties("admin", map[string]any{
		"admin_url":    "https://admin.example.test",
		"access_token": "managed-secret",
		"project_id":   "project-1",
	})
	require.NoError(t, err)
	require.Equal(t, "https://admin.example.test", properties["admin_url"])
	require.Equal(t, "project-1", properties["project_id"])
	require.NotContains(t, properties, "access_token")

	checkpoint, err := structpb.NewStruct(properties)
	require.NoError(t, err)
	require.NotContains(t, checkpoint.String(), "managed-secret")
}
