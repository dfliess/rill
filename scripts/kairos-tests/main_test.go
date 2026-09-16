package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOwnTests(t *testing.T) {
	src := `package server_test

import "testing"

// func TestCommentedOut(t *testing.T) is not a declaration
func TestOne(t *testing.T) {}

func TestTwo (t *testing.T) {
	t.Run("subtest", func(t *testing.T) {})
}

func testHelper(t *testing.T)   {}
func TestingSomething(t *T)     {}
func BenchmarkThree(b *testing.B) {}

// The test binary's entry point reports no result of its own, so it is not a test to account for.
func TestMain(m *testing.M) {}
`
	own := ownTests(map[string]string{"runtime/server/agents_test.go": src})
	require.Equal(t, map[string]map[string]bool{
		"runtime/server": {"TestOne": true, "TestTwo": true},
	}, own)
}

func TestOwnTestsGroupsFilesByPackage(t *testing.T) {
	own := ownTests(map[string]string{
		"runtime/server/a_test.go": "func TestA(t *testing.T) {}\n",
		"runtime/server/b_test.go": "func TestB(t *testing.T) {}\n",
		"runtime/c_test.go":        "func TestC(t *testing.T) {}\n",
	})
	require.Equal(t, map[string]map[string]bool{
		"runtime/server": {"TestA": true, "TestB": true},
		"runtime":        {"TestC": true},
	}, own)
}

func TestCheckPackages(t *testing.T) {
	own := map[string]map[string]bool{
		"runtime/server": {"TestA": true, "TestB": true},
		"runtime/act":    {"TestC": true},
	}

	require.NoError(t, checkPackages(own, []string{"runtime/act", "runtime/server"}))

	err := checkPackages(own, []string{"runtime/act"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime/server (2 tests)")
	require.Contains(t, err.Error(), "nothing runs them")

	err = checkPackages(own, []string{"runtime/act", "runtime/server", "runtime/gone"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime/gone")
	require.Contains(t, err.Error(), "hold no Kairos test any more")
}

func TestParseSkips(t *testing.T) {
	skips, err := parseSkips(`
# a comment

runtime/drivers/mssql TestOLAP # needs a real SQL Server
runtime/drivers/odoo  TestAgainstLocalOdoo   # needs ODOO_TEST_URL
`)
	require.NoError(t, err)
	require.Equal(t, map[testRef]string{
		{Pkg: "runtime/drivers/mssql", Test: "TestOLAP"}:            "needs a real SQL Server",
		{Pkg: "runtime/drivers/odoo", Test: "TestAgainstLocalOdoo"}: "needs ODOO_TEST_URL",
	}, skips)

	_, err = parseSkips("runtime/drivers/mssql TestOLAP\n")
	require.ErrorContains(t, err, "every skip needs a reason")

	_, err = parseSkips("runtime/drivers/mssql TestOLAP #   \n")
	require.ErrorContains(t, err, "every skip needs a reason")

	_, err = parseSkips("TestOLAP # no package\n")
	require.ErrorContains(t, err, "expected \"<package> <Test> # <reason>\"")
}

func TestParseList(t *testing.T) {
	require.Equal(t, []string{"runtime/act", "runtime/server"}, parseList(`
# Packages
runtime/act

runtime/server  # the agent API
`))
	require.Nil(t, parseList("# only comments\n\n"))
}

func TestCheckStaleSkips(t *testing.T) {
	own := map[string]map[string]bool{"runtime/drivers/mssql": {"TestOLAP": true}}

	require.NoError(t, checkStaleSkips(map[testRef]string{
		{Pkg: "runtime/drivers/mssql", Test: "TestOLAP"}: "needs SQL Server",
	}, own))

	err := checkStaleSkips(map[testRef]string{
		{Pkg: "runtime/drivers/mssql", Test: "TestRenamedAway"}: "needs SQL Server",
	}, own)
	require.ErrorContains(t, err, "runtime/drivers/mssql TestRenamedAway")
}

func TestScanEvents(t *testing.T) {
	// A run holding: a pass, a skip, a fail, a subtest (ignored, its parent counts), a package-level event (ignored)
	// and a plain build error line (echoed, not parsed).
	in := strings.Join([]string{
		`{"Action":"run","Package":"github.com/rilldata/rill/runtime/server","Test":"TestOne"}`,
		`{"Action":"output","Package":"github.com/rilldata/rill/runtime/server","Test":"TestOne","Output":"=== RUN   TestOne\n"}`,
		`{"Action":"pass","Package":"github.com/rilldata/rill/runtime/server","Test":"TestOne"}`,
		`{"Action":"skip","Package":"github.com/rilldata/rill/runtime/drivers/mssql","Test":"TestOLAP"}`,
		`{"Action":"fail","Package":"github.com/rilldata/rill/runtime","Test":"TestThree"}`,
		`{"Action":"pass","Package":"github.com/rilldata/rill/runtime/server","Test":"TestOne/subtest"}`,
		`{"Action":"pass","Package":"github.com/rilldata/rill/runtime/server"}`,
		`# github.com/rilldata/rill/runtime [build failed]`,
		"",
	}, "\n")

	var out strings.Builder
	results, err := scanEvents(strings.NewReader(in), &out, "github.com/rilldata/rill")
	require.NoError(t, err)
	require.Equal(t, map[testRef]outcome{
		{Pkg: "runtime/server", Test: "TestOne"}:         outcomePass,
		{Pkg: "runtime/drivers/mssql", Test: "TestOLAP"}: outcomeSkip,
		{Pkg: "runtime", Test: "TestThree"}:              outcomeFail,
	}, results)
	require.Equal(t, "=== RUN   TestOne\n# github.com/rilldata/rill/runtime [build failed]\n", out.String())
}

func TestCheckRan(t *testing.T) {
	own := map[string]map[string]bool{
		"runtime/server":        {"TestRuns": true, "TestSkipsDeclared": true, "TestSkipsUndeclared": true},
		"runtime/drivers/mssql": {"TestFilteredOut": true},
	}
	skips := map[testRef]string{
		{Pkg: "runtime/server", Test: "TestSkipsDeclared"}: "needs a real SQL Server",
	}
	results := map[testRef]outcome{
		{Pkg: "runtime/server", Test: "TestRuns"}:            outcomePass,
		{Pkg: "runtime/server", Test: "TestSkipsDeclared"}:   outcomeSkip,
		{Pkg: "runtime/server", Test: "TestSkipsUndeclared"}: outcomeSkip,
		// TestFilteredOut produced no event at all: -run or -skip left it out, or it did not build.
	}

	err := checkRan(own, skips, results)
	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime/server TestSkipsUndeclared")
	require.Contains(t, err.Error(), "skipped without a reason")
	require.Contains(t, err.Error(), "runtime/drivers/mssql TestFilteredOut")
	require.Contains(t, err.Error(), "produced no result at all")
	require.NotContains(t, err.Error(), "TestSkipsDeclared")
	require.NotContains(t, err.Error(), "TestRuns")

	// A failing test is go test's business, not this command's: it must not be reported twice.
	results[testRef{Pkg: "runtime/server", Test: "TestRuns"}] = outcomeFail
	delete(results, testRef{Pkg: "runtime/server", Test: "TestSkipsUndeclared"})
	own["runtime/server"] = map[string]bool{"TestRuns": true, "TestSkipsDeclared": true}
	own["runtime/drivers/mssql"] = map[string]bool{}
	require.NoError(t, checkRan(own, skips, results))
}

func TestParseArgs(t *testing.T) {
	rev, extra, err := parseArgs([]string{"FETCH_HEAD"})
	require.NoError(t, err)
	require.Equal(t, "FETCH_HEAD", rev)
	require.Empty(t, extra)

	rev, extra, err = parseArgs([]string{"FETCH_HEAD", "--", "-timeout", "45m", "-skip", "TestDedicatedConnector"})
	require.NoError(t, err)
	require.Equal(t, "FETCH_HEAD", rev)
	require.Equal(t, []string{"-timeout", "45m", "-skip", "TestDedicatedConnector"}, extra)

	_, _, err = parseArgs(nil)
	require.ErrorContains(t, err, "usage:")

	_, _, err = parseArgs([]string{"a", "b"})
	require.ErrorContains(t, err, "usage:")
}
