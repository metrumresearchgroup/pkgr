package baseline

import (
	"context"
	"fmt"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const (
	cacheTest1="BSLNCHE-E2E-001"
	cacheTest2="BSLNCHE-E2E-002"
)

func TestClean(t *testing.T) {
	// We can't determine the specific global cache location up front because it varies by OS and might be a tempdir.
	// So this test is only going to verify that it is printed out and used.
	t.Run(MakeTestName(cacheTest1, "pkgr uses and cleans the global pkg/pkgdb caches"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")


		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			t.Errorf("error fetching user cache: %s \n Defaulting to temp dir", err)
			userCacheDir = os.TempDir()
		}
		userPkgrCacheDir := filepath.Join(userCacheDir, "pkgr")
		err = os.RemoveAll(userPkgrCacheDir)
		if err != nil {
			t.Fatalf("error when removing the global pkgr cache: %s", err)
		}

		ctx := context.TODO()
		installCmd := command.New()

		capture, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr-global-cache.yml", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr install: %s", err)
		}
		output := string(capture.Output)

		pkgDbRegex := `\{"level":"info","msg":"Database cache directory: .*\}` // We can't determine this up front because it varies by OS and might be a temp dir.
		pkgCacheRegex:= `\{"level":"info","msg":"Package installation cache directory: .*\}` // We can't determine this up front because it varies by OS and might be a temp dir.
		assert.Regexp(t, pkgDbRegex, output)
		assert.Regexp(t, pkgCacheRegex, output)

		// If that passed, we can assume this will work:
		rPkgDbs := regexp.MustCompile(`\{"level":"info","msg":"Database cache directory:  (.*)"\}`)
		pkgDbDir := rPkgDbs.FindStringSubmatch(output)[1]
		rPkgCache := regexp.MustCompile(`\{"level":"info","msg":"Package installation cache directory:  (.*)"\}`)
		pkgCacheDir := rPkgCache.FindStringSubmatch(output)[1]


		pkgDbDirContents, err := os.ReadDir(filepath.Clean(pkgDbDir))
		if err != nil {
			t.Fatalf("error attempting to read global pkgDb dir: %s", err)
		}
		assert.NotEmpty(t, pkgDbDirContents)

		pkgCacheDirContents, err := os.ReadDir(pkgCacheDir)
		if err != nil {
			t.Fatalf("error attempting to read global package cache dir: %s", err)
		}

		// Find repo folder prefixed with GLOBAL_CACHE_REPO as specified in pkgr-global.cache.yml
		found := false
		for _, entry := range pkgCacheDirContents {
			t.Log(entry.Name())
			if(entry.IsDir() && strings.Contains(entry.Name(), "GLOBAL_CACHE_REPO")) {
				found = true
			}
		}
		assert.True(t, found, "pkgr failed to create a directory in the global pkgr cache for GLOBAL_CACHE_REPO")

		cleanCmd := command.New()
		_, err = cleanCmd.Run(ctx, "pkgr", "clean", "--all", "--config=pkgr-global-cache.yml", "--logjson")
		if err != nil {
			t.Fatalf("error occurred running 'clean' command: %s", err)
		}

		pkgDbDirContents2, err := os.ReadDir(filepath.Clean(pkgDbDir))
		if err != nil {
			t.Fatalf("error attempting to read global pkgDb dir: %s", err)
		}
		assert.NotEmpty(t, pkgDbDirContents)

		pkgCacheDirContents2, err := os.ReadDir(pkgCacheDir)
		if err != nil {
			t.Fatalf("error attempting to read global package cache dir: %s", err)
		}

		assert.Empty(t, pkgDbDirContents2, fmt.Sprintf("pkg database at %s was not cleared", pkgDbDir))
		assert.Len(t, pkgCacheDirContents2, 1, "expected exactly item in the pkgr global cache but found more/less.")
		assert.Equal(t, pkgCacheDirContents2[0].Name(), "r_packagedb_caches", "expected r_packagedb_caches folder to remain in global cache")
	})


	t.Run(MakeTestName(cacheTest2, "pkgr uses and cleans local pkg cache"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")

		localPackageCacheDir := "test-cache"
		err := os.RemoveAll(localPackageCacheDir)
		if err != nil {
			t.Fatalf("error when removing the local package cache: %s", err)
		}

		ctx := context.TODO()
		installCmd := command.New()

		capture, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr.yml", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr install: %s", err)
		}
		output := string(capture.Output)

		pkgCacheMessage := fmt.Sprintf("Package installation cache directory:  %s", localPackageCacheDir)
		assert.Contains(t, output, pkgCacheMessage)


		pkgCacheDirContents, err := os.ReadDir(localPackageCacheDir)
		if err != nil {
			t.Fatalf("error attempting to read global package cache dir: %s", err)
		}

		found := false
		for _, entry := range pkgCacheDirContents {
			t.Log(entry.Name())
			if(entry.IsDir() && strings.Contains(entry.Name(), "CRAN")) {
				found = true
			}
		}
		assert.True(t, found, "pkgr failed to create a directory in the local cache for CRAN")

		cleanCmd := command.New()
		_, err = cleanCmd.Run(ctx, "pkgr", "clean", "cache", "--config=pkgr.yml", "--logjson")
		if err != nil {
			t.Fatalf("error occurred running 'clean' command: %s", err)
		}

		pkgCacheDirContents2, err := os.ReadDir(localPackageCacheDir)
		if err != nil {
			t.Fatalf("error attempting to read global package cache dir: %s", err)
		}

		assert.Empty(t, pkgCacheDirContents2, "found items in test-cache folder (expected empty after clean)")
	})
}