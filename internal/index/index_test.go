package index

import "testing"

func TestIndexAddAndAll(t *testing.T) {
	var idx Index
	idx.Add(Profile{Username: "a"})
	idx.Add(Profile{Username: "b"})

	all := idx.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(all))
	}

	// All() must return a copy: mutating the result must not affect the index.
	all[0].Username = "mutated"
	if idx.Profiles[0].Username != "a" {
		t.Fatalf("All() did not return a copy; index mutated to %q", idx.Profiles[0].Username)
	}
}
