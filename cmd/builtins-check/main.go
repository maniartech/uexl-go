// Command builtins-check is the UExL standard-library consistency checker/mapper.
//
// It loads the canonical manifest (docs/specs/standard-library.json), compares it
// against the live runtime registry (the core vm.Builtins plus every attachable
// standard-library family) and the books (docs/book), and emits a consistency
// matrix to docs/specs/standard-library.md.
//
//	go run ./cmd/builtins-check          # regenerate the matrix
//	go run ./cmd/builtins-check -check    # CI mode: no write, exit 1 on hard drift
//
// Hard drift (fails -check): a function implemented in the registry but absent from
// the manifest (orphan), or a manifest entry marked "implemented" that is missing
// from the registry. Soft drift (reported, not failed): "documented-only" and
// "planned" entries.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/maniartech/uexl/builtins/collections"
	"github.com/maniartech/uexl/builtins/conversion"
	"github.com/maniartech/uexl/builtins/datetime"
	"github.com/maniartech/uexl/builtins/introspection"
	jsonlib "github.com/maniartech/uexl/builtins/json"
	"github.com/maniartech/uexl/builtins/numbers"
	strlib "github.com/maniartech/uexl/builtins/strings"
	"github.com/maniartech/uexl/vm"
)

type arity struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type entry struct {
	Name        string            `json:"name"`
	Category    string            `json:"category"`
	Family      string            `json:"family"`
	Arity       arity             `json:"arity"`
	Params      []string          `json:"params"`
	Returns     string            `json:"returns"`
	OnInvalid   string            `json:"onInvalid"`
	Pure        bool              `json:"pure"`
	Profile     string            `json:"profile"`
	Status      string            `json:"status"`
	SpecRef     string            `json:"specRef"`
	DocRef      string            `json:"docRef"`
	Conformance []string          `json:"conformance"`
	Impl        map[string]string `json:"impl"`
	Notes       string            `json:"notes"`
}

type manifest struct {
	Version  string  `json:"version"`
	Builtins []entry `json:"builtins"`
}

// familyRule declares the ideal member set for a family of similar functions.
// The completeness check reports any expected member that is neither implemented
// nor planned — i.e. the empty cells in a family ("exists on some, not on others").
type familyRule struct {
	family string
	expect []string
	note   string
}

// familyRules encodes the symmetry audit: what *should* exist per family.
var familyRules = []familyRule{
	{"conversion (number)", []string{"parseNum", "tryParseNum", "formatNum"}, "strict/safe number parse plus formatting"},
	{"conversion (boolean)", []string{"parseBool", "tryParseBool"}, "strict/safe boolean parse"},
	{"conversion (string)", []string{"str", "formatNum"}, "value->string is the core str; formatNum covers numbers"},
	{"conversion (json)", []string{"parseJson", "toJson"}, "JSON string <-> value"},
	{"explode-reassemble", []string{"split", "join", "runes", "graphemes", "bytes"}, "split/join plus the unicode views"},
	{"length (unicode views)", []string{"len", "runeLen", "graphemeLen", "utf16Len"}, "utf16 view still missing (JS parity, optional)"},
	{"substring (unicode views)", []string{"substr", "runeSubstr", "graphemeSubstr", "utf16Substr"}, "utf16 view still missing (JS parity, optional)"},
	{"collection accessors", []string{"len", "get", "set", "has", "keys", "values", "remove", "merge"}, "complete; map/filter/sort are pipe-based by design"},
	{"string-ops", []string{"contains", "indexOf", "startsWith", "endsWith", "replace", "trim", "trimStart", "trimEnd", "upper", "lower", "split", "padStart", "padEnd", "repeat"}, "complete"},
	{"math", []string{"abs", "sign", "round", "floor", "ceil", "trunc", "min", "max", "sum", "avg", "pow", "sqrt", "mod", "clamp"}, "complete"},
	{"type-introspection", []string{"typeOf", "isNull", "isNumber", "isString", "isBool", "isArray", "isObject", "isDate", "isDuration", "isEmpty"}, "complete"},
}

type row struct {
	name        string
	category    string
	status      string
	implemented bool
	documented  bool
	specified   bool
	verdict     string
	hardDrift   bool
}

func main() {
	var (
		manifestPath = flag.String("manifest", "docs/specs/standard-library.json", "path to the manifest JSON")
		bookDir      = flag.String("book", "docs/book", "directory of the books to scan for documented names")
		outPath      = flag.String("out", "docs/specs/standard-library.md", "path to write the consistency matrix")
		checkOnly    = flag.Bool("check", false, "CI mode: do not write; exit 1 on hard drift")
	)
	flag.Parse()

	m, err := loadManifest(*manifestPath)
	if err != nil {
		fail("load manifest: %v", err)
	}

	registry := registryNames()
	documented := scanDocumented(*bookDir, candidateNames(m, registry))

	rows, hardDrift := buildRows(m, registry, documented)

	if *checkOnly {
		if hardDrift > 0 {
			fmt.Fprintf(os.Stderr, "builtins-check: %d hard-drift issue(s) found\n", hardDrift)
			printDrift(os.Stderr, rows)
			os.Exit(1)
		}
		fmt.Println("builtins-check: OK (no hard drift)")
		return
	}

	present := presentSet(m, registry)
	md := renderMatrix(m, rows) + renderCompleteness(present)
	if err := os.WriteFile(*outPath, []byte(md), 0o644); err != nil {
		fail("write matrix: %v", err)
	}
	fmt.Printf("wrote %s (%d builtins; %d hard-drift, %d soft-drift, %d family gaps)\n",
		*outPath, len(rows), hardDrift, countSoft(rows), countFamilyGaps(present))
}

func loadManifest(path string) (*manifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// registryNames is the set of function names the runtime can actually provide: the always-on core
// (vm.Builtins) plus every attachable standard-library family (registered outside vm.Builtins, via
// uexl.WithStdlib / the per-family options). now()/today() are clock-injected and registered separately
// from datetime.Builtins, so they are added explicitly.
func registryNames() map[string]bool {
	out := map[string]bool{}
	for name := range vm.Builtins {
		out[name] = true
	}
	addKeys(out, numbers.Builtins)
	addKeys(out, conversion.Builtins)
	addKeys(out, introspection.Builtins)
	addKeys(out, strlib.Builtins)
	addKeys(out, collections.Builtins)
	addKeys(out, jsonlib.Builtins)
	addKeys(out, datetime.Builtins)
	out["now"] = true
	out["today"] = true
	return out
}

// addKeys records every key of m as present in out (value type ignored; families use different fn types).
func addKeys[V any](out map[string]bool, m map[string]V) {
	for name := range m {
		out[name] = true
	}
}

func candidateNames(m *manifest, registry map[string]bool) []string {
	seen := map[string]bool{}
	for _, e := range m.Builtins {
		seen[e.Name] = true
	}
	for n := range registry {
		seen[n] = true
	}
	out := make([]string, 0, len(seen))
	for n := range seen {
		out = append(out, n)
	}
	return out
}

// scanDocumented walks bookDir and reports which candidate names appear as a
// function call (name followed by "(") in any .md file.
func scanDocumented(bookDir string, names []string) map[string]bool {
	res := map[string]bool{}
	pats := map[string]*regexp.Regexp{}
	for _, n := range names {
		pats[n] = regexp.MustCompile(`\b` + regexp.QuoteMeta(n) + `\s*\(`)
	}
	_ = filepath.WalkDir(bookDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(b)
		for n, re := range pats {
			if !res[n] && re.MatchString(text) {
				res[n] = true
			}
		}
		return nil
	})
	return res
}

func buildRows(m *manifest, registry, documented map[string]bool) ([]row, int) {
	inManifest := map[string]bool{}
	var rows []row
	hard := 0

	for _, e := range m.Builtins {
		inManifest[e.Name] = true
		r := row{
			name:        e.Name,
			category:    e.category(),
			status:      e.Status,
			implemented: registry[e.Name],
			documented:  documented[e.Name],
			specified:   e.SpecRef != "",
		}
		switch {
		case r.implemented && e.Status == "implemented":
			r.verdict = "OK"
		case r.implemented && e.Status != "implemented":
			r.verdict = "registry ahead of manifest (update status)"
		case !r.implemented && e.Status == "implemented":
			r.verdict = "MISSING IMPL"
			r.hardDrift = true
		case !r.implemented && e.Status == "documented-only":
			r.verdict = "DRIFT: documented, not implemented"
		case !r.implemented && e.Status == "planned":
			r.verdict = "planned"
		case !r.implemented && e.Status == "gap":
			r.verdict = "gap (expected by symmetry)"
		case e.Status == "deprecated":
			r.verdict = "deprecated"
		default:
			r.verdict = "review"
		}
		if r.hardDrift {
			hard++
		}
		rows = append(rows, r)
	}

	// Orphans: implemented in the registry but missing from the manifest.
	for name := range registry {
		if inManifest[name] {
			continue
		}
		rows = append(rows, row{
			name:        name,
			category:    "?",
			status:      "(none)",
			implemented: true,
			documented:  documented[name],
			specified:   false,
			verdict:     "ORPHAN: implemented, unmanifested",
			hardDrift:   true,
		})
		hard++
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].category != rows[j].category {
			return rows[i].category < rows[j].category
		}
		return rows[i].name < rows[j].name
	})
	return rows, hard
}

func (e entry) category() string {
	if e.Category == "" {
		return "?"
	}
	return e.Category
}

func tick(b bool) string {
	if b {
		return "✓"
	}
	return "✗"
}

func renderMatrix(m *manifest, rows []row) string {
	var b strings.Builder
	b.WriteString("# UExL Standard-Library Consistency Matrix\n\n")
	b.WriteString("> Generated by `cmd/builtins-check` from `docs/specs/standard-library.json`.\n")
	b.WriteString("> Do not edit by hand — edit the manifest and regenerate.\n\n")
	b.WriteString(fmt.Sprintf("Manifest version: `%s`\n\n", m.Version))
	b.WriteString("Columns: **Impl** = present in `vm.Builtins`; **Doc** = referenced in `docs/book`; ")
	b.WriteString("**Spec** = has a spec reference; **Status** = manifest intent.\n\n")
	b.WriteString("| Builtin | Category | Impl | Doc | Spec | Status | Verdict |\n")
	b.WriteString("|---------|----------|:----:|:---:|:----:|--------|---------|\n")
	for _, r := range rows {
		b.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s | %s | %s |\n",
			r.name, r.category, tick(r.implemented), tick(r.documented), tick(r.specified), r.status, r.verdict))
	}
	b.WriteString("\n## Drift summary\n\n")
	var hard, soft []row
	for _, r := range rows {
		switch {
		case r.hardDrift:
			hard = append(hard, r)
		case strings.HasPrefix(r.verdict, "DRIFT") || strings.Contains(r.verdict, "ahead"):
			soft = append(soft, r)
		}
	}
	if len(hard) == 0 && len(soft) == 0 {
		b.WriteString("No drift. \n")
	}
	if len(hard) > 0 {
		b.WriteString("**Hard drift (must fix):**\n\n")
		for _, r := range hard {
			b.WriteString(fmt.Sprintf("- `%s` — %s\n", r.name, r.verdict))
		}
		b.WriteString("\n")
	}
	if len(soft) > 0 {
		b.WriteString("**Soft drift (reconcile):**\n\n")
		for _, r := range soft {
			b.WriteString(fmt.Sprintf("- `%s` — %s\n", r.name, r.verdict))
		}
		b.WriteString("\n")
	}
	return b.String()
}

// presentSet returns the names that are actually available: implemented in the
// registry, or acknowledged as implemented/planned in the manifest. Names that
// are only documented-only or gap are NOT present.
func presentSet(m *manifest, registry map[string]bool) map[string]bool {
	present := map[string]bool{}
	for n := range registry {
		present[n] = true
	}
	for _, e := range m.Builtins {
		if e.Status == "implemented" || e.Status == "planned" {
			present[e.Name] = true
		}
	}
	return present
}

// renderCompleteness produces the "exists on some, not on others" report: for
// each family rule, which expected members are present vs missing.
func renderCompleteness(present map[string]bool) string {
	var b strings.Builder
	b.WriteString("\n## Family completeness (symmetry audit)\n\n")
	b.WriteString("For each family of similar functions, which expected members exist vs. are missing. ")
	b.WriteString("`present` = implemented or planned; `missing` = neither.\n\n")
	totalMissing := 0
	for _, fr := range familyRules {
		var have, miss []string
		for _, name := range fr.expect {
			if present[name] {
				have = append(have, name)
			} else {
				miss = append(miss, name)
			}
		}
		totalMissing += len(miss)
		b.WriteString(fmt.Sprintf("### %s\n\n", fr.family))
		if fr.note != "" {
			b.WriteString(fmt.Sprintf("_%s_\n\n", fr.note))
		}
		b.WriteString(fmt.Sprintf("- present: %s\n", listOrNone(have)))
		if len(miss) > 0 {
			b.WriteString(fmt.Sprintf("- **missing: %s**\n", listCode(miss)))
		} else {
			b.WriteString("- missing: — (complete)\n")
		}
		b.WriteString("\n")
	}
	b.WriteString(fmt.Sprintf("**Total missing siblings across families: %d.**\n", totalMissing))
	return b.String()
}

func listCode(xs []string) string {
	if len(xs) == 0 {
		return "—"
	}
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = "`" + x + "`"
	}
	return strings.Join(out, ", ")
}

func listOrNone(xs []string) string {
	if len(xs) == 0 {
		return "— (none)"
	}
	return listCode(xs)
}

func countFamilyGaps(present map[string]bool) int {
	n := 0
	for _, fr := range familyRules {
		for _, name := range fr.expect {
			if !present[name] {
				n++
			}
		}
	}
	return n
}

func countSoft(rows []row) int {
	n := 0
	for _, r := range rows {
		if !r.hardDrift && (strings.HasPrefix(r.verdict, "DRIFT") || strings.Contains(r.verdict, "ahead")) {
			n++
		}
	}
	return n
}

func printDrift(w *os.File, rows []row) {
	for _, r := range rows {
		if r.hardDrift {
			fmt.Fprintf(w, "  - %s: %s\n", r.name, r.verdict)
		}
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "builtins-check: "+format+"\n", args...)
	os.Exit(2)
}
