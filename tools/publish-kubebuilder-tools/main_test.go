package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
			got, err := getRegistryAuthToken(path, "registry.ci.openshift.org/ocp/release")
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
		_, err := getRegistryAuthToken(filepath.Join(t.TempDir(), "missing"), "registry.ci.openshift.org/ocp/release")
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("error = %v, want missing file", err)
		}
	})
}

// TestImageReferences checks legacy and Quay repository paths and digests are
// preserved when constructing manifest and blob URLs.
func TestImageReferences(t *testing.T) {
	for _, repository := range []string{
		"registry.ci.openshift.org/ocp/release",
		"registry.ci.openshift.org/ocp/release-5",
		"registry.ci.openshift.org/ocp/kube-artifacts",
		quayRegistry + "/openshift/ci",
	} {
		t.Run(repository, func(t *testing.T) {
			wantDigest := fixtureDigest("component")
			image, digest := getImageStreamAndDigest(repository + "@" + wantDigest)
			if image != repository || digest != wantDigest {
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

// TestVerifyAPIServerVersion checks native command execution and semantic version
// matching. Only the dirty suffix and build metadata are ignored; other
// prereleases must match, and invalid versions or execution failures are rejected.
func TestVerifyAPIServerVersion(t *testing.T) {
	for _, tc := range []struct {
		name, expected, output, wantError string
	}{
		{"stable", "v1.37.1", "Kubernetes v1.37.1", ""},
		{"dirty", "v1.37.1", "Kubernetes v1.37.1-dirty", ""},
		{"build metadata", "v1.37.1+expected", "Kubernetes v1.37.1+actual", ""},
		{"dirty with metadata", "v1.37.1", "Kubernetes v1.37.1-dirty+build.123", ""},
		{"prerelease", "v1.37.1-rc.1", "Kubernetes v1.37.1-rc.1-dirty+build.123", ""},
		{"major mismatch", "v1.37.1", "Kubernetes v2.37.1", "does not match -version"},
		{"minor mismatch", "v1.37.1", "Kubernetes v1.38.1", "does not match -version"},
		{"patch mismatch", "v1.37.1", "Kubernetes v1.37.2", "does not match -version"},
		{"unexpected prerelease", "v1.37.1", "Kubernetes v1.37.1-rc.1", "does not match -version"},
		{"different prerelease", "v1.37.1-rc.1", "Kubernetes v1.37.1-rc.2", "does not match -version"},
		{"dirty is not a suffix", "v1.37.1", "Kubernetes v1.37.1-dirty.1", "does not match -version"},
		{"invalid expected version", "vgarbage", "Kubernetes v1.37.1", "failed to parse Kubernetes version"},
		{"invalid reported version", "v1.37.1", "Kubernetes unknown", "failed to parse Kubernetes version"},
		{"invalid build metadata", "v1.37.1", "Kubernetes v1.37.1+invalid!", "failed to parse Kubernetes version"},
		{"incomplete version", "v1.37.1", "Kubernetes v1.37", "failed to parse Kubernetes version"},
		{"empty output", "v1.37.1", "", "failed to parse Kubernetes version"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			binary := filepath.Join(dir, runtime.GOOS, runtime.GOARCH, "bin", "kube-apiserver")
			writeTestFile(t, binary, []byte(fmt.Sprintf("#!/bin/sh\n[ \"$1\" = '--version' ] || exit 1\nprintf '%%s\\n' '%s'\n", tc.output)))
			err := verifyAPIServerVersion(dir, tc.expected)
			if tc.wantError != "" {
				requireError(t, err, tc.wantError)
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
	t.Run("missing native binary", func(t *testing.T) {
		requireError(t, verifyAPIServerVersion(t.TempDir(), "v1.37.1"), "failed to check kube-apiserver version")
	})
	t.Run("command failure", func(t *testing.T) {
		dir := t.TempDir()
		writeTestFile(t, filepath.Join(dir, runtime.GOOS, runtime.GOARCH, "bin", "kube-apiserver"), []byte("#!/bin/sh\nexit 2\n"))
		requireError(t, verifyAPIServerVersion(dir, "v1.37.1"), "failed to check kube-apiserver version")
	})
}

const quayRegistry = "quay-proxy.ci.openshift.org"

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
