package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"gopkg.in/yaml.v3"
	utilversion "k8s.io/apimachinery/pkg/util/version"
)

func main() {
	pullSecretFile := flag.String("pull-secret", "", "The pull secret to use for the kubebuilder tools")
	version := flag.String("version", "", "The version of the kubebuilder tools to publish. This should be a Kubernetes version that the build is based upon.")
	outputDir := flag.String("output-dir", "", "The output directory to write the kubebuilder tools to")
	payload := flag.String("payload", "", "The CI payload: registry.ci.openshift.org/ocp/release[-<number>]:<tag> or quay-proxy.ci.openshift.org/openshift/ci:rc_payload__<version>")
	skipUpload := flag.Bool("skip-upload", false, "Skip uploading the artifacts created to the openshift-gce-devel/openshift-kubebuilder-tools bucket")
	indexFile := flag.String("index-file", "envtest-releases.yaml", "The index file to use for the kubebuilder tools")

	flag.Parse()

	if *pullSecretFile == "" {
		panic("pull-secret is required")
	}

	if *version == "" {
		panic("version is required")
	}

	if *outputDir == "" {
		panic("output is required")
	}

	// We only expect Kubernetes versions that begin with the v prefix.
	if !strings.HasPrefix(*version, "v") {
		panic("Kubernetes version must begin with the v prefix. Example: v1.33.2")
	}

	// We only accept valid semver versions that Kubernetes uses.
	if _, err := utilversion.ParseSemantic(*version); err != nil {
		panic(fmt.Errorf("failed to parse Kubernetes version %q: %w", *version, err))
	}

	// Accept both the legacy release repositories and the quay-proxy CI repository.
	matches := regexp.MustCompile(`^((?:registry\.ci\.openshift\.org/ocp/release(?:-\d+)?|quay-proxy\.ci\.openshift\.org/openshift/ci)):(.+)$`).FindStringSubmatch(*payload)
	if *payload == "" || matches == nil {
		panic(fmt.Errorf("failed to validate payload: payload is required and must be a valid payload: registry.ci.openshift.org/ocp/release[-<number>]:<tag> or quay-proxy.ci.openshift.org/openshift/ci:rc_payload__<version>; got %q", *payload))
	}

	releaseImage := matches[1]
	payloadVersion := matches[2]
	if strings.HasPrefix(releaseImage, "quay-proxy.ci.openshift.org/") && (!strings.HasPrefix(payloadVersion, "rc_payload__") || payloadVersion == "rc_payload__") {
		panic(fmt.Errorf("failed to validate payload tag: Quay-proxy payload tags must start with rc_payload__ followed by a version; got %q", payloadVersion))
	}

	// Resolve credentials once; the Quay payload and artifact images share a repository.
	registryAuthToken, err := getRegistryAuthToken(*pullSecretFile, releaseImage)
	if err != nil {
		panic(fmt.Errorf("failed to get registry authentication for %q: %w", releaseImage, err))
	}

	// Download the image-references and convert to a map of image name to digest
	manifests, err := getReleaseImages(releaseImage, payloadVersion, registryAuthToken)
	if err != nil {
		panic(err)
	}

	// Extract the kube-apiserver binaries from the installer-kube-apiserver-artifacts image
	if err := getKubeAPIServerBins(*outputDir, manifests, registryAuthToken); err != nil {
		panic(err)
	}
	if err := verifyAPIServerVersion(*outputDir, *version); err != nil {
		panic(fmt.Errorf("failed to verify extracted kube-apiserver version: %w", err))
	}

	// Extract the etcd binaries from the installer-etcd-artifacts image
	if err := getEtcdBins(*outputDir, manifests, registryAuthToken); err != nil {
		panic(err)
	}

	// Build the envtest archives for each os and arch combination
	if err := buildEnvtestTars(*outputDir, *version); err != nil {
		panic(err)
	}

	if *skipUpload {
		fmt.Printf("Archives written to %s\n", *outputDir)
		return
	}

	// Upload the tars created earlier to the public GCS bucket for general consumption
	if err := uploadArchives(*outputDir, *version); err != nil {
		panic(err)
	}

	// Update the index file with the new version
	if err := updateIndexFile(*outputDir, *version, *indexFile); err != nil {
		panic(err)
	}

	fmt.Printf("Archives uploaded to openshift-gce-devel/openshift-kubebuilder-tools for version %s\n", *version)
}

func getRegistryAuthToken(pullSecretFile, image string) (string, error) {
	registry, repository, _ := strings.Cut(image, "/")
	pullSecretRaw, err := os.ReadFile(pullSecretFile)
	if err != nil {
		return "", fmt.Errorf("failed to read pull secret %q: %w", pullSecretFile, err)
	}

	var secret struct {
		Auths map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
	}

	if err := json.Unmarshal(pullSecretRaw, &secret); err != nil {
		return "", fmt.Errorf("failed to decode pull secret %q: %w", pullSecretFile, err)
	}

	registryAuth, ok := secret.Auths[registry]
	if !ok {
		return "", fmt.Errorf("failed to find registry credentials: registry %q not found in pull secret", registry)
	}

	registryAuthToken, err := base64.StdEncoding.DecodeString(registryAuth.Auth)
	if err != nil {
		return "", fmt.Errorf("failed to decode credentials for registry %q: %w", registry, err)
	}

	// Passwords may contain colons, particularly when used as Basic credentials.
	credentials := strings.SplitN(string(registryAuthToken), ":", 2)
	if len(credentials) != 2 {
		return "", fmt.Errorf("failed to decode credentials: password not found in pull secret for registry %q", registry)
	}
	if registry != "quay-proxy.ci.openshift.org" {
		return credentials[1], nil
	}

	// Quay-proxy requires a scoped token rather than the legacy bearer password.
	query := url.Values{"service": {registry}, "scope": {"repository:" + repository + ":pull"}}
	req, err := http.NewRequest("GET", "https://"+registry+"/v2/auth?"+query.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create token request for registry %q: %w", registry, err)
	}
	req.SetBasicAuth(credentials[0], credentials[1])
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange credentials for registry %q: %w", registry, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to exchange credentials: token exchange bad status for registry %q: %s", registry, resp.Status)
	}
	var token struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("failed to decode token response for registry %q: %w", registry, err)
	}
	if token.Token == "" {
		token.Token = token.AccessToken
	}
	if token.Token == "" {
		return "", fmt.Errorf("failed to exchange credentials: token exchange returned no token for registry %q", registry)
	}
	return token.Token, nil
}

func getReleaseImages(releaseImage string, version string, registryToken string) (map[string]string, error) {
	releaseManifestRaw, err := downloadJSON(getRegistryURL(releaseImage, "manifests", version), registryToken)
	if err != nil {
		return nil, err
	}

	manifest := struct {
		FSLayers []struct {
			BlobSum string `json:"blobSum"`
		} `json:"fsLayers"`
	}{}

	if err := json.Unmarshal(releaseManifestRaw, &manifest); err != nil {
		return nil, err
	}

	if len(manifest.FSLayers) == 0 {
		return nil, errors.New("no fsLayers found in release manifest")
	}

	// The first fsLayer is the release image manifests, which has the image digests for the release images
	releaseImageLayer, close, err := downloadArchive(getRegistryURL(releaseImage, "blobs", manifest.FSLayers[0].BlobSum), registryToken)
	if err != nil {
		return nil, err
	}
	defer close()

	imageReferencesRaw, err := getFileFromArchive(releaseImageLayer, "release-manifests/image-references")
	if err != nil {
		return nil, err
	}

	imageReferences := struct {
		Spec struct {
			Tags []struct {
				Name string
				From struct {
					Name string `json:"name"`
				} `json:"from"`
			} `json:"tags"`
		} `json:"spec"`
	}{}

	if err := json.Unmarshal(imageReferencesRaw, &imageReferences); err != nil {
		return nil, err
	}

	images := map[string]string{}
	for _, tag := range imageReferences.Spec.Tags {
		images[tag.Name] = tag.From.Name
	}

	return images, nil
}

func getRegistryURL(image, kind, digest string) string {
	registry, repository, _ := strings.Cut(image, "/")
	return fmt.Sprintf("https://%s/v2/%s/%s/%s", registry, repository, kind, digest)
}

func downloadJSON(url string, registryToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", registryToken))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// From an in-memory archive, traverse and find the file named and return it as a slice of bytes.
func getFileFromArchive(archive *tar.Reader, filename string) ([]byte, error) {
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if header.Name == filename {
			buf := bytes.NewBuffer(nil)
			_, err := io.Copy(buf, archive)
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("file %s not found in archive", filename)
}

// Fetch a tar.gz (container image layer) into memory so that we can extract files from it.
func downloadArchive(url string, registryToken string) (*tar.Reader, func() error, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", registryToken))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, nil, err
	}

	close := resp.Body.Close

	if resp.StatusCode != http.StatusOK {
		close()
		return nil, nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	uncompressedStream, err := gzip.NewReader(resp.Body)
	if err != nil {
		close()
		return nil, nil, err
	}

	return tar.NewReader(uncompressedStream), close, nil
}

func getKubeAPIServerBins(dir string, manifests map[string]string, registryToken string) error {
	return getMultiArchBinariesFromImage(dir, manifests, registryToken, "installer-kube-apiserver-artifacts", "kube-apiserver")

}

func getEtcdBins(dir string, manifests map[string]string, registryToken string) error {
	return getMultiArchBinariesFromImage(dir, manifests, registryToken, "installer-etcd-artifacts", "etcd")
}

func getMultiArchBinariesFromImage(dir string, manifests map[string]string, registryToken string, imageName string, binaryName string) error {
	image, ok := manifests[imageName]
	if !ok {
		return fmt.Errorf("%q not found in release images", imageName)
	}

	imageStream, digest := getImageStreamAndDigest(image)

	imageManifestRaw, err := downloadJSON(getRegistryURL(imageStream, "manifests", digest), registryToken)
	if err != nil {
		return err
	}

	imageManifest := struct {
		Manifests []struct {
			Digest string `json:"digest"`
		} `json:"manifests"`
	}{}

	if err := json.Unmarshal(imageManifestRaw, &imageManifest); err != nil {
		return err
	}

	if len(imageManifest.Manifests) != 1 {
		return fmt.Errorf("expected 1 image manifest for image stream %q, got 0 or more than 1", imageName)
	}

	imageLayerManifestRaw, err := downloadJSON(getRegistryURL(imageStream, "manifests", imageManifest.Manifests[0].Digest), registryToken)
	if err != nil {
		return err
	}

	imageLayerManifest := struct {
		Layers []struct {
			Digest string `json:"digest"`
		} `json:"layers"`
	}{}

	if err := json.Unmarshal(imageLayerManifestRaw, &imageLayerManifest); err != nil {
		return err
	}

	// The last layer is the layer containing the binaries.
	imageLayer, close, err := downloadArchive(getRegistryURL(imageStream, "blobs", imageLayerManifest.Layers[len(imageLayerManifest.Layers)-1].Digest), registryToken)
	if err != nil {
		return err
	}
	defer close()

	for _, goos := range []string{"darwin", "linux"} {
		for _, arch := range []string{"amd64", "arm64"} {
			binDir := filepath.Join(dir, goos, arch, "bin")
			if err := os.MkdirAll(binDir, 0755); err != nil {
				return err
			}

			binRaw, err := getFileFromArchive(imageLayer, fmt.Sprintf("usr/share/openshift/%s/%s/%s", goos, arch, binaryName))
			if err != nil {
				return err
			}

			if err := os.WriteFile(filepath.Join(binDir, binaryName), binRaw, 0755); err != nil {
				return err
			}
		}
	}

	return nil
}

func getImageStreamAndDigest(image string) (string, string) {
	parts := strings.Split(image, "@")
	return parts[0], parts[1]
}

// Verify the native kube-apiserver binary before creating archives or publishing them.
// The version reported by the binary must match the one specified in -version.
// Before comparison the version is normalized by dropping -dirty and build metadata.
func verifyAPIServerVersion(dir, expected string) error {
	binary := filepath.Join(dir, runtime.GOOS, runtime.GOARCH, "bin", "kube-apiserver")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, binary, "--version").Output()
	if err != nil {
		return fmt.Errorf("failed to check kube-apiserver version for %s/%s: %w", runtime.GOOS, runtime.GOARCH, err)
	}
	actual := strings.TrimSpace(strings.TrimPrefix(string(output), "Kubernetes "))
	var normalized [2]string
	for i, value := range []string{expected, actual} {
		parsed, err := utilversion.ParseSemantic(value)
		if err != nil {
			return fmt.Errorf("failed to parse Kubernetes version %q: %w", value, err)
		}
		withoutMetadata, _, _ := strings.Cut(parsed.String(), "+")
		normalized[i] = strings.TrimSuffix(withoutMetadata, "-dirty")
	}
	if normalized[0] != normalized[1] {
		return fmt.Errorf("failed to verify kube-apiserver version: kube-apiserver version %q does not match -version %q", actual, expected)
	}
	return nil
}

func buildEnvtestTars(dir string, version string) error {
	for _, goos := range []string{"darwin", "linux"} {
		for _, arch := range []string{"amd64", "arm64"} {
			out, err := os.Create(filepath.Join(dir, fmt.Sprintf("envtest-%s-%s-%s.tar.gz", version, goos, arch)))
			if err != nil {
				return err
			}
			defer out.Close()

			gzArchive := gzip.NewWriter(out)
			defer gzArchive.Close()

			binFS := os.DirFS(filepath.Join(dir, goos, arch))
			tarArchive := tar.NewWriter(gzArchive)
			defer tarArchive.Close()

			if err := addFS(tarArchive, binFS); err != nil {
				return err
			}
		}
	}

	return nil
}

// addFS copied from std library for Tar, introduced in Go 1.22. Copy/paste for now.
func addFS(tw *tar.Writer, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		// TODO(#49580): Handle symlinks when fs.ReadLinkFS is available.
		if !info.Mode().IsRegular() {
			return errors.New("tar: cannot add non-regular file")
		}
		h, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		h.Name = name
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		f, err := fsys.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return err
	})
}

func uploadArchives(dir string, version string) error {
	gcsClient, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}

	for _, goos := range []string{"darwin", "linux"} {
		for _, arch := range []string{"amd64", "arm64"} {
			archivePath := filepath.Join(dir, fmt.Sprintf("envtest-%s-%s-%s.tar.gz", version, goos, arch))
			if err := uploadArchive(gcsClient, archivePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func uploadArchive(gcsClient *storage.Client, archivePath string) error {
	gcsObj := gcsClient.Bucket("openshift-kubebuilder-tools").Object(filepath.Base(archivePath))
	gcsWriter := gcsObj.NewWriter(context.Background())
	defer gcsWriter.Close()

	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(gcsWriter, f); err != nil {
		return err
	}

	return nil
}

type indexFile struct {
	Hash     string `yaml:"hash"`
	SelfLink string `yaml:"selfLink"`
}

// updateIndexFile adds the new version to the existing index file.
// The index file is used by the setup-envtest tool to find the download links for the envtest archives.
func updateIndexFile(dir, version, indexFileName string) error {
	indexFileRaw, err := os.ReadFile(indexFileName)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read index file: %w", err)
	}

	index := struct {
		Releases map[string]map[string]indexFile `json:"releases"`
	}{}

	if indexFileRaw != nil {
		if err := yaml.Unmarshal(indexFileRaw, &index); err != nil {
			return fmt.Errorf("failed to unmarshal index file: %w", err)
		}
	} else {
		index.Releases = make(map[string]map[string]indexFile)
	}

	releaseIndexes := make(map[string]indexFile)
	for _, goos := range []string{"darwin", "linux"} {
		for _, arch := range []string{"amd64", "arm64"} {
			name := fmt.Sprintf("envtest-%s-%s-%s.tar.gz", version, goos, arch)

			archive, err := os.Open(filepath.Join(dir, name))
			if err != nil {
				return fmt.Errorf("failed to open archive %s: %w", name, err)
			}
			defer archive.Close()

			hash, err := hashFile(archive)
			if err != nil {
				return fmt.Errorf("failed to hash archive %s: %w", name, err)
			}

			releaseIndexes[name] = indexFile{
				Hash:     hash,
				SelfLink: fmt.Sprintf("https://storage.googleapis.com/openshift-kubebuilder-tools/%s", name),
			}
		}
	}

	index.Releases[version] = releaseIndexes

	indexRaw, err := yaml.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to marshal index file: %w", err)
	}

	if err := os.WriteFile(indexFileName, indexRaw, 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

func hashFile(f *os.File) (string, error) {
	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
