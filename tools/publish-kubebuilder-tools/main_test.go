package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGetRegistryAuthToken checks legacy bearer credentials are selected and
// decoded from the pull secret, and unreadable or malformed secrets fail.
func TestGetRegistryAuthToken(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("user:" + testRegistryToken))
	for _, tc := range []struct {
		name, secret, wantError string
	}{
		{"legacy credentials", `{"auths":{"registry.ci.openshift.org":{"auth":"` + encoded + `"},"quay.io":{"auth":"ignored"}}}`, ""},
		{"missing registry", `{"auths":{"quay.io":{"auth":"` + encoded + `"}}}`, "not found in pull secret"},
		{"invalid JSON", `{`, "unexpected end"},
		{"invalid base64", `{"auths":{"registry.ci.openshift.org":{"auth":"!"}}}`, "base64"},
		{"missing separator", `{"auths":{"registry.ci.openshift.org":{"auth":"` + base64.StdEncoding.EncodeToString([]byte("user")) + `"}}}`, "password not found"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "pull-secret.json")
			writeTestFile(t, path, []byte(tc.secret))
			got, err := getRegistryAuthToken(path)
			if tc.wantError != "" {
				requireError(t, err, tc.wantError)
				return
			}
			if err != nil || got != testRegistryToken {
				t.Fatalf("token = %q, error = %v", got, err)
			}
		})
	}
	t.Run("missing file", func(t *testing.T) {
		_, err := getRegistryAuthToken(filepath.Join(t.TempDir(), "missing"))
		if !os.IsNotExist(err) {
			t.Fatalf("error = %v, want missing file", err)
		}
	})
}

// TestImageReferences checks legacy repository paths and digests are
// preserved when constructing manifest and blob URLs.
func TestImageReferences(t *testing.T) {
	for _, repository := range []string{
		"registry.ci.openshift.org/ocp/release",
		"registry.ci.openshift.org/ocp/release-5",
		"registry.ci.openshift.org/ocp/kube-artifacts",
	} {
		t.Run(repository, func(t *testing.T) {
			wantDigest := fixtureDigest("component")
			image, digest := getImageStreamAndDigest(repository + "@" + wantDigest)
			if image != strings.TrimPrefix(repository, "registry.ci.openshift.org/ocp/") || digest != wantDigest {
				t.Fatalf("image = %q, digest = %q", image, digest)
			}
			for _, kind := range []string{"manifests", "blobs"} {
				registry, repo, _ := strings.Cut(repository, "/")
				want := "https://" + registry + "/v2/" + repo + "/" + kind + "/" + wantDigest
				if got := getRegistryURL(image, kind, digest); got != want {
					t.Errorf("URL = %q, want %q", got, want)
				}
			}
		})
	}
}

const testRegistryToken = "fixture-token"

func fixtureDigest(name string) string {
	return fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(name)))
}

func writeTestFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0755); err != nil {
		t.Fatal(err)
	}
}

func requireError(t *testing.T, err error, contains string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), contains) {
		t.Fatalf("error = %v, want error containing %q", err, contains)
	}
}
