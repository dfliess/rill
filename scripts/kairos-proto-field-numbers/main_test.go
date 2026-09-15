package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// upstreamProto stands for an upstream .proto file. It packs the constructs the parser must get through: a service
// with a multi-line aggregate option, bracketed field options with braces inside a string, reserved statements, a
// block comment with braces, a map, a oneof, and an enum and a message nested in a message.
const upstreamProto = `syntax = "proto3";

package demo.v1;

import "google/api/annotations.proto";
import "validate/validate.proto";

option go_package = "example.com/demo/v1";

service DemoService {
  // Get returns a resource.
  rpc Get(GetRequest) returns (GetResponse) {
    option (google.api.http) = {
      post: "/v1/get/{name}",
      body: "*"
    };
  }
}

message GetRequest {
  string name = 1 [(validate.rules).string = {min_len: 1, in: ["a", "b}"]}];
  reserved 2, 3;
  reserved "old_name";
}

message GetResponse {
  Resource resource = 1;
}

/* A resource.
   Braces in a comment { do not open a block. */
message Resource {
  map<string, string> labels = 1;
  oneof kind {
    Model model = 2;
    Theme theme = 3;
  }
  enum State {
    STATE_UNSPECIFIED = 0;
    STATE_READY = 4;
  }
  State state = 5; // Not a field: other = 99;
  message Model {
    repeated string tags = 1;
  }
}

message Theme {
  optional string color = 1;
}
`

func parseProto(t *testing.T, src string) messages {
	t.Helper()
	res := make(messages)
	require.NoError(t, parseFile(res, "demo.proto", src))
	return res
}

func TestParseFile(t *testing.T) {
	got := parseProto(t, upstreamProto)

	numbers := make(map[string]map[string]int64)
	for msg, fields := range got {
		numbers[msg] = make(map[string]int64)
		for name, f := range fields {
			numbers[msg][name] = f.number
		}
	}
	require.Equal(t, map[string]map[string]int64{
		"demo.v1.GetRequest":     {"name": 1},
		"demo.v1.GetResponse":    {"resource": 1},
		"demo.v1.Resource":       {"labels": 1, "model": 2, "theme": 3, "state": 5},
		"demo.v1.Resource.Model": {"tags": 1},
		"demo.v1.Theme":          {"color": 1},
	}, numbers)

	labels := got["demo.v1.Resource"]["labels"]
	require.Equal(t, "demo.proto", labels.file)
	require.Equal(t, strings.Count(upstreamProto[:strings.Index(upstreamProto, "map<string")], "\n")+1, labels.line)
}

func TestCheck(t *testing.T) {
	tt := []struct {
		name           string
		after          string // Line of upstreamProto the head version inserts below; empty appends at the end.
		insert         string
		wantAdded      int
		wantViolations []string // Without the "file:line: " prefix.
	}{
		{
			name: "unchanged",
		},
		{
			name:           "low field in an upstream message",
			after:          "  Resource resource = 1;\n",
			insert:         "  string note = 2;\n",
			wantAdded:      1,
			wantViolations: []string{"demo.v1.GetResponse.note = 2"},
		},
		{
			name:      "field from 1000 in an upstream message",
			after:     "  Resource resource = 1;\n",
			insert:    "  string note = 1000;\n",
			wantAdded: 1,
		},
		{
			name:           "low field in a oneof of an upstream message",
			after:          "    Theme theme = 3;\n",
			insert:         "    Agent agent = 18;\n",
			wantAdded:      1,
			wantViolations: []string{"demo.v1.Resource.agent = 18"},
		},
		{
			name:           "low map field in an upstream message",
			after:          "  map<string, string> labels = 1;\n",
			insert:         "  map<string, int32> counts = 6;\n",
			wantAdded:      1,
			wantViolations: []string{"demo.v1.Resource.counts = 6"},
		},
		{
			name:           "low field in an upstream nested message",
			after:          "    repeated string tags = 1;\n",
			insert:         "    string owner = 2;\n",
			wantAdded:      1,
			wantViolations: []string{"demo.v1.Resource.Model.owner = 2"},
		},
		{
			name:           "low field with bracketed options",
			after:          "  reserved \"old_name\";\n",
			insert:         "  string mode = 4 [(validate.rules).string = {in: [\"x\", \"y\"]}, deprecated = true];\n",
			wantAdded:      1,
			wantViolations: []string{"demo.v1.GetRequest.mode = 4"},
		},
		{
			// If the parser lost track of the nested braces, the messages below the service would come out wrong.
			name:  "new rpc with a nested multi-line option",
			after: "      body: \"*\"\n    };\n  }\n",
			insert: `  rpc List(GetRequest) returns (GetResponse) {
    option (google.api.http) = {
      get: "/v1/list"
      additional_bindings {
        get: "/v1/items"
      }
    };
  }
`,
		},
		{
			name:           "low field with braces in a string of its options",
			after:          "  optional string color = 1;\n",
			insert:         "  string accent = 2 [(validate.rules).string = {pattern: \"^#[0-9a-f]{6}$\"}];\n",
			wantAdded:      1,
			wantViolations: []string{"demo.v1.Theme.accent = 2"},
		},
		{
			name:   "fields of a Kairos-only message",
			insert: "\nmessage Agent {\n  string name = 1;\n  AgentSpec spec = 2;\n}\n",
		},
		{
			name:   "fields of a Kairos-only message nested in an upstream message",
			after:  "    repeated string tags = 1;\n  }\n",
			insert: "  message Extra {\n    string x = 1;\n  }\n",
		},
		{
			name:   "value added to an upstream enum",
			after:  "    STATE_READY = 4;\n",
			insert: "    STATE_FAILED = 7;\n",
		},
		{
			name:   "reserved range added to an upstream message",
			after:  "  reserved 2, 3;\n",
			insert: "  reserved 10 to 20;\n",
		},
	}

	upstream := parseProto(t, upstreamProto)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			head := upstreamProto + tc.insert
			if tc.after != "" {
				require.Equal(t, 1, strings.Count(upstreamProto, tc.after), "anchor must match once")
				head = strings.Replace(upstreamProto, tc.after, tc.after+tc.insert, 1)
			}

			added, violations := check(upstream, parseProto(t, head))
			require.Len(t, added, tc.wantAdded)
			for i, v := range violations {
				_, desc, ok := strings.Cut(v, ": ")
				require.True(t, ok)
				violations[i] = desc
			}
			if tc.wantViolations == nil {
				require.Empty(t, violations)
			} else {
				require.Equal(t, tc.wantViolations, violations)
			}
		})
	}
}

func TestParseFileErrors(t *testing.T) {
	tt := map[string]string{
		"unterminated block":   "message A {\n  string x = 1;\n",
		"unbalanced brace":     "message A {\n}\n}\n",
		"unterminated comment": "/* message A {}",
		"unterminated string":  "message A {\n  string x = 1 [default = \"abc];\n}\n",
		"unterminated option":  "message A {\n  string x = 1 [deprecated = true;\n}\n",
	}
	for name, src := range tt {
		t.Run(name, func(t *testing.T) {
			require.Error(t, parseFile(make(messages), "bad.proto", src))
		})
	}
}
