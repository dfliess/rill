// Command kairos-proto-field-numbers fails if Kairos added a field numbered below 1000 to a proto message that also
// exists upstream (kairos-cloud#158). Upstream numbers its fields from 1 up, so a low number chosen by Kairos collides
// sooner or later with a field upstream adds, and messages that are entirely Kairos' own are not checked.
//
// It compares the .proto files under proto/ of two git revisions:
//
//	go run ./scripts/kairos-proto-field-numbers <upstream-rev> [<head-rev>]
//
// A field of <head-rev> is taken as added by Kairos when the same message in <upstream-rev> has no field with its
// name. <head-rev> defaults to HEAD. It must contain <upstream-rev> (as a rebased sync branch does): otherwise a field
// that upstream removed after the fork point looks as if Kairos had added it.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const minKairosFieldNumber = 1000

// field is a message field as declared in a .proto file.
type field struct {
	number int64
	file   string
	line   int
}

// messages maps a message's fully-qualified name to its fields by name.
type messages map[string]map[string]field

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <upstream-rev> [<head-rev>]\n", os.Args[0])
		os.Exit(2)
	}
	upstreamRev := os.Args[1]
	headRev := "HEAD"
	if len(os.Args) == 3 {
		headRev = os.Args[2]
	}

	upstream, err := loadMessages(upstreamRev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	if len(upstream) == 0 {
		fmt.Fprintf(os.Stderr, "error: no proto messages found under proto/ in %s\n", upstreamRev)
		os.Exit(2)
	}
	head, err := loadMessages(headRev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	added, violations := check(upstream, head)
	if len(violations) > 0 {
		fmt.Printf("Kairos fields added to upstream messages must be numbered from %d (kairos-cloud#158):\n", minKairosFieldNumber)
		for _, v := range violations {
			fmt.Printf("  %s\n", v)
		}
		os.Exit(1)
	}
	fmt.Printf("OK: %d Kairos fields in upstream messages (%s vs %s), all numbered from %d:\n", len(added), headRev, upstreamRev, minKairosFieldNumber)
	for _, a := range added {
		fmt.Printf("  %s\n", a)
	}
}

// check returns the fields of head that Kairos added to messages that exist in upstream, and those of them numbered
// below minKairosFieldNumber. Both are sorted and described as "file:line: message.field = number".
func check(upstream, head messages) (added, violations []string) {
	for msg, fields := range head {
		upstreamFields, ok := upstream[msg]
		if !ok {
			continue // A Kairos-only message numbers its fields from 1.
		}
		for name, f := range fields {
			if _, ok := upstreamFields[name]; ok {
				continue
			}
			desc := fmt.Sprintf("%s:%d: %s.%s = %d", f.file, f.line, msg, name, f.number)
			added = append(added, desc)
			if f.number < minKairosFieldNumber {
				violations = append(violations, desc)
			}
		}
	}
	sort.Strings(added)
	sort.Strings(violations)
	return added, violations
}

// loadMessages parses every .proto file under proto/ at the given git revision.
func loadMessages(rev string) (messages, error) {
	out, err := exec.Command("git", "ls-tree", "-r", "--name-only", rev, "--", "proto").Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-tree %s: %w", rev, err)
	}
	res := make(messages)
	for _, path := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if !strings.HasSuffix(path, ".proto") {
			continue
		}
		src, err := exec.Command("git", "show", rev+":"+path).Output()
		if err != nil {
			return nil, fmt.Errorf("git show %s:%s: %w", rev, path, err)
		}
		if err := parseFile(res, path, string(src)); err != nil {
			return nil, fmt.Errorf("%s at %s: %w", path, rev, err)
		}
	}
	return res, nil
}

type token struct {
	text string
	line int
}

// parseFile adds the messages declared in a .proto source to res. It understands just enough of the language to find
// message fields: blocks, statements, and option values (bracketed or aggregate) that may contain braces themselves.
func parseFile(res messages, path, src string) error {
	toks, err := tokenize(src)
	if err != nil {
		return err
	}

	type scope struct {
		kind string // "message", "oneof" or "other"
		name string // fully-qualified name, for messages
	}
	var pkg string
	var stack []scope
	var stmt []token
	depth := 0 // nesting of brackets, parentheses and aggregate values inside the current statement

	enclosingMessage := func() string {
		for i := len(stack) - 1; i >= 0; i-- {
			switch stack[i].kind {
			case "message":
				return stack[i].name
			case "other":
				return ""
			}
		}
		return ""
	}

	for _, tok := range toks {
		if depth > 0 {
			switch tok.text {
			case "[", "(", "{":
				depth++
			case "]", ")", "}":
				depth--
			}
			stmt = append(stmt, tok)
			continue
		}

		switch tok.text {
		case "[", "(":
			depth++
			stmt = append(stmt, tok)
		case "{":
			if hasToken(stmt, "=") { // Aggregate option value, e.g. option (google.api.http) = {...};
				depth++
				stmt = append(stmt, tok)
				continue
			}
			if len(stmt) == 0 {
				return fmt.Errorf("line %d: unexpected block", tok.line)
			}
			switch {
			case stmt[0].text == "message" && len(stmt) == 2:
				name := stmt[1].text
				if parent := enclosingMessage(); parent != "" {
					name = parent + "." + name
				} else if pkg != "" {
					name = pkg + "." + name
				}
				stack = append(stack, scope{kind: "message", name: name})
				if res[name] == nil {
					res[name] = make(map[string]field)
				}
			case stmt[0].text == "oneof":
				stack = append(stack, scope{kind: "oneof"})
			default:
				stack = append(stack, scope{kind: "other"})
			}
			stmt = nil
		case "}":
			if len(stack) == 0 || len(stmt) > 0 {
				return fmt.Errorf("line %d: unexpected }", tok.line)
			}
			stack = stack[:len(stack)-1]
		case ";":
			if len(stmt) >= 2 && stmt[0].text == "package" {
				pkg = stmt[1].text
			}
			if msg := enclosingMessage(); msg != "" {
				if name, f, ok := parseField(stmt); ok {
					f.file = path
					res[msg][name] = f
				}
			}
			stmt = nil
		default:
			stmt = append(stmt, tok)
		}
	}
	if len(stack) > 0 || depth > 0 || len(stmt) > 0 {
		return fmt.Errorf("unexpected end of file")
	}
	return nil
}

// parseField recognizes a field statement such as `repeated string tags = 3 [...]` or `map<string, string> m = 4`.
func parseField(stmt []token) (string, field, bool) {
	if len(stmt) < 4 {
		return "", field{}, false
	}
	switch stmt[0].text {
	case "option", "reserved", "extensions":
		return "", field{}, false
	}
	for i, tok := range stmt {
		if tok.text != "=" {
			continue
		}
		if i < 2 || i+1 >= len(stmt) {
			return "", field{}, false
		}
		n, err := strconv.ParseInt(stmt[i+1].text, 0, 64)
		if err != nil {
			return "", field{}, false
		}
		return stmt[i-1].text, field{number: n, line: stmt[i-1].line}, true
	}
	return "", field{}, false
}

func hasToken(toks []token, text string) bool {
	for _, t := range toks {
		if t.text == text {
			return true
		}
	}
	return false
}

// tokenize splits a .proto source into identifiers, numbers, string literals and single-character symbols,
// dropping whitespace and comments.
func tokenize(src string) ([]token, error) {
	var toks []token
	line := 1
	for i := 0; i < len(src); {
		c := src[i]
		switch {
		case c == '\n':
			line++
			i++
		case c == ' ' || c == '\t' || c == '\r':
			i++
		case strings.HasPrefix(src[i:], "//"):
			for i < len(src) && src[i] != '\n' {
				i++
			}
		case strings.HasPrefix(src[i:], "/*"):
			end := strings.Index(src[i+2:], "*/")
			if end < 0 {
				return nil, fmt.Errorf("line %d: unterminated comment", line)
			}
			line += strings.Count(src[i:i+2+end], "\n")
			i += 2 + end + 2
		case c == '"' || c == '\'':
			start := i
			i++
			for i < len(src) && src[i] != c {
				if src[i] == '\\' {
					i++
				}
				if i < len(src) && src[i] == '\n' {
					return nil, fmt.Errorf("line %d: unterminated string", line)
				}
				i++
			}
			if i >= len(src) {
				return nil, fmt.Errorf("line %d: unterminated string", line)
			}
			i++
			toks = append(toks, token{text: src[start:i], line: line})
		case isWordChar(c):
			start := i
			for i < len(src) && (isWordChar(src[i]) || src[i] == '.') {
				i++
			}
			toks = append(toks, token{text: src[start:i], line: line})
		default:
			toks = append(toks, token{text: string(c), line: line})
			i++
		}
	}
	return toks, nil
}

func isWordChar(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}
