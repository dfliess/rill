// Command kairos-tests runs the Go tests Kairos owns and fails when one of them would go unnoticed
// (kairos-cloud#166).
//
//	go run ./scripts/kairos-tests <upstream-rev> [-- <extra go test flags>]
//
// A test belongs to Kairos when the _test.go file declaring it exists in HEAD and not in <upstream-rev>: Kairos writes
// its tests in files of its own, and deriving the set instead of listing it is what keeps the list from going stale.
// Tests Kairos added inside a file that also exists upstream are not covered by this, and are the one gap left.
//
// The command then:
//
//  1. compares the packages holding those files with packages.txt, and fails on any difference, so a test written in a
//     package nobody declared cannot pass unnoticed;
//  2. runs `go test -short` over the declared packages, streaming their output unchanged;
//  3. fails unless every Kairos test ran. A test that skipped must be declared in skips.txt with a reason; a test that
//     produced no result at all (filtered out by -run or -skip, excluded by a build tag, renamed) fails too, because
//     "it did not run" is exactly the silence this command exists to break.
//
// The rule behind it is the one Act learned the hard way: a green that proves nothing is worse than no green, because
// the list of covered packages then lies about what is covered.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// dir holds this command's own files, relative to the repository root.
const dir = "scripts/kairos-tests"

// testRef is a top-level test function, identified by the package that declares it.
type testRef struct {
	Pkg  string // package path relative to the repository root, e.g. "runtime/server"
	Test string
}

func (r testRef) String() string { return r.Pkg + " " + r.Test }

// outcome is what a test run reported for a test.
type outcome int

const (
	outcomeMissing outcome = iota // no result at all: filtered out, not built, or gone
	outcomePass
	outcomeSkip
	outcomeFail
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	upstreamRev, extra, err := parseArgs(args)
	if err != nil {
		return err
	}

	files, err := kairosTestFiles(upstreamRev)
	if err != nil {
		return err
	}
	sources := make(map[string]string, len(files))
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		sources[f] = string(b)
	}
	own, err := ownTests(sources)
	if err != nil {
		return err
	}
	if len(own) == 0 {
		return fmt.Errorf("no Kairos test files found against %s: is it an upstream revision?", upstreamRev)
	}

	declared, err := readList(path.Join(dir, "packages.txt"))
	if err != nil {
		return err
	}
	if err := checkPackages(own, declared); err != nil {
		return err
	}

	skips, err := readSkips(path.Join(dir, "skips.txt"))
	if err != nil {
		return err
	}
	if err := checkStaleSkips(skips, own); err != nil {
		return err
	}

	module, err := modulePath()
	if err != nil {
		return err
	}

	results, testErr := goTest(declared, extra, module, stdout)
	if err := checkRan(own, skips, results); err != nil {
		return err
	}
	if testErr != nil {
		return testErr
	}
	fmt.Fprintf(stdout, "\nOK: %d Kairos tests ran in %d packages (%d declared as skipped)\n", count(own), len(declared), len(skips))
	return nil
}

// parseArgs splits the upstream revision from the flags to pass on to go test.
func parseArgs(args []string) (upstreamRev string, extra []string, err error) {
	for i, a := range args {
		if a == "--" {
			extra = args[i+1:]
			args = args[:i]
			break
		}
	}
	if len(args) != 1 {
		return "", nil, fmt.Errorf("usage: kairos-tests <upstream-rev> [-- <extra go test flags>]")
	}
	return args[0], extra, nil
}

// kairosTestFiles returns the _test.go files that are in HEAD and not in the given revision.
func kairosTestFiles(upstreamRev string) ([]string, error) {
	head, err := treeFiles("HEAD")
	if err != nil {
		return nil, err
	}
	upstream, err := treeFiles(upstreamRev)
	if err != nil {
		return nil, err
	}
	var out []string
	for f := range head {
		if !strings.HasSuffix(f, "_test.go") || upstream[f] {
			continue
		}
		out = append(out, f)
	}
	sort.Strings(out)
	return out, nil
}

func treeFiles(rev string) (map[string]bool, error) {
	out, err := exec.Command("git", "ls-tree", "-r", "--name-only", rev).Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-tree %s: %w", rev, err)
	}
	files := map[string]bool{}
	for _, l := range strings.Split(string(out), "\n") {
		if l != "" {
			files[l] = true
		}
	}
	return files, nil
}

// ownTests maps each package to the top-level tests declared in the given sources, keyed by file path. It parses the
// sources rather than matching text, because a test file may hold Go source inside a string literal, as this command's
// own tests do. TestMain is not a test: it is the binary's entry point and reports no result of its own.
func ownTests(sources map[string]string) (map[string]map[string]bool, error) {
	own := map[string]map[string]bool{}
	fset := token.NewFileSet()
	for file, src := range sources {
		f, err := parser.ParseFile(fset, file, src, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", file, err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !isTestName(fn.Name.Name) {
				continue
			}
			if fn.Type.Params == nil || len(fn.Type.Params.List) != 1 { // A test takes one parameter, *testing.T
				continue
			}
			pkg := path.Dir(file)
			if own[pkg] == nil {
				own[pkg] = map[string]bool{}
			}
			own[pkg][fn.Name.Name] = true
		}
	}
	return own, nil
}

// isTestName reports whether a function name is a test's, by the rule `go test` applies: Test followed by something
// that does not start with a lower-case letter. TestMain is excluded: it runs the binary, and reports nothing.
func isTestName(name string) bool {
	rest, ok := strings.CutPrefix(name, "Test")
	if !ok || name == "TestMain" {
		return false
	}
	if rest == "" {
		return true
	}
	r, _ := utf8.DecodeRuneInString(rest)
	return !unicode.IsLower(r)
}

func count(own map[string]map[string]bool) int {
	n := 0
	for _, tests := range own {
		n += len(tests)
	}
	return n
}

// checkPackages fails when the declared list and the packages holding Kairos tests have drifted apart.
func checkPackages(own map[string]map[string]bool, declared []string) error {
	declaredSet := map[string]bool{}
	for _, p := range declared {
		declaredSet[p] = true
	}
	var missing []string
	for pkg := range own {
		if !declaredSet[pkg] {
			missing = append(missing, pkg)
		}
	}
	var extra []string
	for _, p := range declared {
		if own[p] == nil {
			extra = append(extra, p)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)

	var b strings.Builder
	if len(missing) > 0 {
		fmt.Fprintf(&b, "These packages hold Kairos tests and are not in %s/packages.txt, so nothing runs them:\n", dir)
		for _, p := range missing {
			fmt.Fprintf(&b, "  %s (%d tests)\n", p, len(own[p]))
		}
	}
	if len(extra) > 0 {
		fmt.Fprintf(&b, "These packages are in %s/packages.txt and hold no Kairos test any more:\n", dir)
		for _, p := range extra {
			fmt.Fprintf(&b, "  %s\n", p)
		}
	}
	if b.Len() > 0 {
		return fmt.Errorf("%s", strings.TrimRight(b.String(), "\n"))
	}
	return nil
}

// checkStaleSkips fails when a declared skip names a test that no longer exists.
func checkStaleSkips(skips map[testRef]string, own map[string]map[string]bool) error {
	var stale []string
	for ref := range skips {
		if !own[ref.Pkg][ref.Test] {
			stale = append(stale, ref.String())
		}
	}
	if len(stale) == 0 {
		return nil
	}
	sort.Strings(stale)
	return fmt.Errorf("%s/skips.txt declares tests that no longer exist:\n  %s", dir, strings.Join(stale, "\n  "))
}

// checkRan fails unless every Kairos test produced a result, and every skip was declared with a reason.
func checkRan(own map[string]map[string]bool, skips map[testRef]string, results map[testRef]outcome) error {
	var undeclared, missing []string
	for pkg, tests := range own {
		for test := range tests {
			ref := testRef{Pkg: pkg, Test: test}
			switch results[ref] {
			case outcomeSkip:
				if _, ok := skips[ref]; !ok {
					undeclared = append(undeclared, ref.String())
				}
			case outcomeMissing:
				missing = append(missing, ref.String())
			}
		}
	}
	sort.Strings(undeclared)
	sort.Strings(missing)

	var b strings.Builder
	if len(undeclared) > 0 {
		fmt.Fprintf(&b, "These Kairos tests skipped without a reason declared in %s/skips.txt:\n  %s\n", dir, strings.Join(undeclared, "\n  "))
	}
	if len(missing) > 0 {
		fmt.Fprintf(&b, "These Kairos tests produced no result at all (filtered out, not built, or renamed):\n  %s\n", strings.Join(missing, "\n  "))
	}
	if b.Len() > 0 {
		return fmt.Errorf("%s", strings.TrimRight(b.String(), "\n"))
	}
	return nil
}

// goTest runs the declared packages and returns what each test reported, streaming the output as go test prints it.
func goTest(packages, extra []string, module string, stdout io.Writer) (map[testRef]outcome, error) {
	args := append([]string{"test", "-short", "-count=1", "-json"}, extra...)
	for _, p := range packages {
		args = append(args, "./"+p)
	}
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	results, scanErr := scanEvents(out, stdout, module)
	waitErr := cmd.Wait()
	if scanErr != nil {
		return results, scanErr
	}
	if waitErr != nil {
		return results, fmt.Errorf("go test failed: %w", waitErr)
	}
	return results, nil
}

// event is the part of `go test -json` output this command reads.
type event struct {
	Action  string
	Package string
	Test    string
	Output  string
}

// scanEvents echoes the test output and records the outcome of each top-level test.
func scanEvents(r io.Reader, stdout io.Writer, module string) (map[testRef]outcome, error) {
	results := map[testRef]outcome{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		var ev event
		if err := json.Unmarshal(sc.Bytes(), &ev); err != nil {
			// Not a JSON line: go test writes build errors straight to stdout.
			fmt.Fprintln(stdout, sc.Text())
			continue
		}
		if ev.Action == "output" {
			fmt.Fprint(stdout, ev.Output)
			continue
		}
		if ev.Test == "" || strings.Contains(ev.Test, "/") { // package-level event, or a subtest
			continue
		}
		var o outcome
		switch ev.Action {
		case "pass":
			o = outcomePass
		case "skip":
			o = outcomeSkip
		case "fail":
			o = outcomeFail
		default:
			continue
		}
		results[testRef{Pkg: strings.TrimPrefix(ev.Package, module+"/"), Test: ev.Test}] = o
	}
	return results, sc.Err()
}

func modulePath() (string, error) {
	b, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}
	for _, l := range strings.Split(string(b), "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(l), "module "); ok {
			return strings.TrimSpace(rest), nil
		}
	}
	return "", fmt.Errorf("no module line in go.mod")
}

// readList reads a file of one entry per line, ignoring blank lines and # comments.
func readList(file string) ([]string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return parseList(string(b)), nil
}

func parseList(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if i := strings.Index(l, "#"); i >= 0 {
			l = l[:i]
		}
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// readSkips reads the declared skips: "<package> <Test> # <reason>", one per line. The reason is mandatory: an
// exclusion nobody justified is one nobody can review later.
func readSkips(file string) (map[testRef]string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return parseSkips(string(b))
}

func parseSkips(s string) (map[testRef]string, error) {
	skips := map[testRef]string{}
	for i, l := range strings.Split(s, "\n") {
		line := strings.TrimSpace(l)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entry, reason, ok := strings.Cut(line, "#")
		if !ok || strings.TrimSpace(reason) == "" {
			return nil, fmt.Errorf("%s/skips.txt:%d: every skip needs a reason after #: %q", dir, i+1, line)
		}
		fields := strings.Fields(entry)
		if len(fields) != 2 {
			return nil, fmt.Errorf("%s/skips.txt:%d: expected \"<package> <Test> # <reason>\", got %q", dir, i+1, line)
		}
		skips[testRef{Pkg: fields[0], Test: fields[1]}] = strings.TrimSpace(reason)
	}
	return skips, nil
}
