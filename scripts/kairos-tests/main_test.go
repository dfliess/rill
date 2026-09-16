package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// ownTestsFixture is a test file as the command finds one on disk. Everything in it other than TestOne and TestTwo is
// something the command must not count: a commented-out declaration, a helper, a benchmark, a method, a function whose
// signature is not a test's, the binary's entry point, and a test declared inside a string literal, which is what this
// command's own tests hold.
const ownTestsFixture = "package server_test\n" + `
import "testing"

// func TestCommentedOut(t *testing.T) is not a declaration
func TestOne(t *testing.T) {}

func TestTwo (t *testing.T) {
	t.Run("subtest", func(t *testing.T) {})
}

func testHelper(t *testing.T)               {}
func TestingSomething(t *testing.T)         {}
func BenchmarkThree(b *testing.B)           {}
func TestWithTwoParams(t *testing.T, n int) {}

func (s *suite) TestMethod(t *testing.T) {}

func TestMain(m *testing.M) {}

const embedded = "func TestInsideAStringLiteral(t *testing.T) {}"
`

func TestOwnTests(t *testing.T) {
	own, err := ownTests(map[string]string{"runtime/server/agents_test.go": ownTestsFixture})
	require.NoError(t, err)
	require.Equal(t, map[string]map[string]bool{
		"runtime/server": {"TestOne": true, "TestTwo": true},
	}, own)
}

func TestOwnTestsGroupsFilesByPackage(t *testing.T) {
	own, err := ownTests(map[string]string{
		"runtime/server/a_test.go": "package server_test\nfunc TestA(t *testing.T) {}\n",
		"runtime/server/b_test.go": "package server_test\nfunc TestB(t *testing.T) {}\n",
		"runtime/c_test.go":        "package runtime_test\nfunc TestC(t *testing.T) {}\n",
	})
	require.NoError(t, err)
	require.Equal(t, map[string]map[string]bool{
		"runtime/server": {"TestA": true, "TestB": true},
		"runtime":        {"TestC": true},
	}, own)
}

func TestOwnTestsRejectsUnparseableSource(t *testing.T) {
	_, err := ownTests(map[string]string{"runtime/server/broken_test.go": "package server_test\nfunc TestA(\n"})
	require.ErrorContains(t, err, "parse runtime/server/broken_test.go")
}

func TestCheckPackages(t *testing.T) {
	own := map[string]map[string]bool{
		"runtime/server": {"TestA": true, "TestB": true},
		"runtime/act":    {"TestC": true},
	}
	declared := []pkg{{Path: "runtime/act"}, {Path: "runtime/server"}}

	require.NoError(t, checkPackages(own, declared))

	err := checkPackages(own, []pkg{{Path: "runtime/act"}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime/server (2 tests)")
	require.Contains(t, err.Error(), "nothing runs them")

	// A package whose tests upstream absorbed: red while nobody says why it stays, fine once someone does.
	err = checkPackages(own, append(declared, pkg{Path: "runtime/drivers/gemini"}))
	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime/drivers/gemini")
	require.Contains(t, err.Error(), "hold no Kairos test any more")

	reasoned := append(declared, pkg{Path: "runtime/drivers/gemini", Reason: "upstream absorbed its tests; the fork still carries changes in the driver"})
	require.NoError(t, checkPackages(own, reasoned))
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

func TestParsePackages(t *testing.T) {
	require.Equal(t, []pkg{
		{Path: "runtime/act"},
		{Path: "runtime/server", Reason: "the agent API"},
	}, parsePackages(`
# A section heading, not a package
runtime/act

runtime/server  # the agent API
`))
	require.Nil(t, parsePackages("# only comments\n\n"))
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
