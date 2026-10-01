package problems

import (
	"encoding/json"
	"strings"
	"testing"
	"testing/fstest"
)

const goodStatement = "Do the thing.\n\n## Constraints\n\n- small\n"

const goodTests = `{
  "examples": [
    {"label": "Example 1", "input": "1\n5", "expected": "5"},
    {"label": "Example 2", "input": "1\n6", "expected": "6"}
  ],
  "edge": [
    {"label": "Edge 1", "input": "1\n0", "expected": "0"},
    {"label": "Edge 2", "input": "1\n1", "expected": "1"},
    {"label": "Edge 3", "input": "0\n", "expected": "[]"}
  ]
}`

// pack adds one valid pack to fsys. change can edit the meta before it is written.
func pack(fsys fstest.MapFS, id string, change func(m *Meta)) {
	m := Meta{ID: id, Title: id, Difficulty: "easy", Sheet: "patterns", Group: "G", Order: 1, Generator: "pair-sum"}
	if change != nil {
		change(&m)
	}
	b, _ := json.Marshal(m)
	fsys[id+"/meta.json"] = &fstest.MapFile{Data: b}
	fsys[id+"/statement.md"] = &fstest.MapFile{Data: []byte(goodStatement)}
	fsys[id+"/tests.json"] = &fstest.MapFile{Data: []byte(goodTests)}
}

func validate(t *testing.T, fsys fstest.MapFS) error {
	t.Helper()
	lib, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	return Validate(lib)
}

func wantError(t *testing.T, err error, part string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), part) {
		t.Fatalf("want an error containing %q, got %v", part, err)
	}
}

func TestBuiltinPacksValidate(t *testing.T) {
	lib, err := Load(Builtin())
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(lib); err != nil {
		t.Fatal(err)
	}
}

func TestValidateAcceptsLinkedPacks(t *testing.T) {
	fsys := fstest.MapFS{}
	pack(fsys, "a", func(m *Meta) { m.LeadsTo = []string{"b"} })
	pack(fsys, "b", func(m *Meta) { m.BuildsOn = []string{"a"}; m.Order = 2 })
	if err := validate(t, fsys); err != nil {
		t.Fatal(err)
	}
}

func TestValidateRejectsMissingLink(t *testing.T) {
	fsys := fstest.MapFS{}
	pack(fsys, "a", func(m *Meta) { m.BuildsOn = []string{"ghost"} })
	wantError(t, validate(t, fsys), `buildsOn names "ghost"`)

	fsys = fstest.MapFS{}
	pack(fsys, "a", func(m *Meta) { m.LeadsTo = []string{"ghost"} })
	wantError(t, validate(t, fsys), `leadsTo names "ghost"`)
}

func TestValidateRejectsCycle(t *testing.T) {
	fsys := fstest.MapFS{}
	pack(fsys, "a", func(m *Meta) { m.BuildsOn = []string{"c"} })
	pack(fsys, "b", func(m *Meta) { m.BuildsOn = []string{"a"} })
	pack(fsys, "c", func(m *Meta) { m.BuildsOn = []string{"b"} })
	wantError(t, validate(t, fsys), "cycle")
}

func TestValidateRejectsOneSidedLeadsTo(t *testing.T) {
	fsys := fstest.MapFS{}
	pack(fsys, "a", func(m *Meta) { m.LeadsTo = []string{"b"} })
	pack(fsys, "b", nil)
	wantError(t, validate(t, fsys), `"b" does not build on it`)
}

func TestValidateRejectsBadMeta(t *testing.T) {
	cases := map[string]func(m *Meta){
		"group is not set":  func(m *Meta) { m.Group = "" },
		"order must be":     func(m *Meta) { m.Order = 0 },
		"sheet must be":     func(m *Meta) { m.Sheet = "other" },
		"is not registered": func(m *Meta) { m.Generator = "nope" },
	}
	for part, change := range cases {
		fsys := fstest.MapFS{}
		pack(fsys, "a", change)
		wantError(t, validate(t, fsys), part)
	}
}

func TestValidateRejectsBadStatementAndTests(t *testing.T) {
	fsys := fstest.MapFS{}
	pack(fsys, "a", nil)
	fsys["a/statement.md"] = &fstest.MapFile{Data: []byte("No limits given.\n")}
	wantError(t, validate(t, fsys), "Constraints")

	fsys = fstest.MapFS{}
	pack(fsys, "a", nil)
	fsys["a/tests.json"] = &fstest.MapFile{Data: []byte(`{"examples":[{"label":"E","input":"1","expected":"1"}],"edge":[]}`)}
	err := validate(t, fsys)
	wantError(t, err, "at least 2 examples")
	wantError(t, err, "at least 3 edge cases")

	for _, in := range []string{`1 \n5`, `1\n5\n\n`, `\n1\n5`} {
		fsys = fstest.MapFS{}
		pack(fsys, "a", nil)
		bad := strings.Replace(goodTests, `"input": "1\n5"`, `"input": "`+in+`"`, 1)
		fsys["a/tests.json"] = &fstest.MapFile{Data: []byte(bad)}
		wantError(t, validate(t, fsys), `"Example 1": input`)
	}
}

func TestValidateRealSheetRules(t *testing.T) {
	inRides := func(order int, buildsOn ...string) func(m *Meta) {
		return func(m *Meta) { m.Sheet, m.Group, m.Order, m.BuildsOn = "real", "Rides", order, buildsOn }
	}
	withFollowUp := func(fsys fstest.MapFS, id string) {
		fsys[id+"/statement.md"] = &fstest.MapFile{Data: []byte(goodStatement + "\n## Follow-up\n\nWhat if it grows?\n")}
	}

	fsys := fstest.MapFS{}
	pack(fsys, "basics", nil)
	pack(fsys, "first", inRides(101))
	pack(fsys, "second", inRides(102, "first", "basics"))
	withFollowUp(fsys, "first")
	withFollowUp(fsys, "second")
	if err := validate(t, fsys); err != nil {
		t.Fatalf("the first of a group may skip the patterns link: %v", err)
	}

	pack(fsys, "third", inRides(103, "second"))
	withFollowUp(fsys, "third")
	wantError(t, validate(t, fsys), `third: real-interview problem needs a buildsOn link`)

	fsys = fstest.MapFS{}
	pack(fsys, "first", inRides(101))
	wantError(t, validate(t, fsys), `first: real-interview statement has no "## Follow-up"`)
}
