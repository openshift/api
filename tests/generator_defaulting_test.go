package tests

import (
    "testing"
)

func TestDefaultOnCreateExpected(t *testing.T) {
    initial := "apiVersion: v1\nkind: Pod\nmetadata:\n  name: foo\n"

    // When expected is empty, should default to initial
    got := defaultOnCreateExpected(initial, "")
    if string(got) != initial {
        t.Fatalf("expected default to be initial YAML; got %q", string(got))
    }

    // When expected provided, should use it
    expectedYAML := "apiVersion: v1\nkind: Pod\nmetadata:\n  name: bar\n"
    got = defaultOnCreateExpected(initial, expectedYAML)
    if string(got) != expectedYAML {
        t.Fatalf("expected to use provided expected YAML; got %q", string(got))
    }
}

func TestDefaultOnUpdateExpected(t *testing.T) {
    updated := "apiVersion: v1\nkind: Pod\nmetadata:\n  name: updated\n"

    // When expected is empty, should default to updated
    got := defaultOnUpdateExpected(updated, "")
    if string(got) != updated {
        t.Fatalf("expected default to be updated YAML; got %q", string(got))
    }

    // When expected provided, should use it
    expectedYAML := "apiVersion: v1\nkind: Pod\nmetadata:\n  name: something\n"
    got = defaultOnUpdateExpected(updated, expectedYAML)
    if string(got) != expectedYAML {
        t.Fatalf("expected to use provided expected YAML; got %q", string(got))
    }
}
