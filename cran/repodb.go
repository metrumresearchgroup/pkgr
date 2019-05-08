package cran

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/metrumresearchgroup/pkgr/desc"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

// NewRepoDb returns a new Repo database
func NewRepoDb(url RepoURL, dst SourceType, rc RepoConfig, rv RVersion) (*RepoDb, error) {
	repoDatabasePointer := &RepoDb{
		DescriptionsBySourceType: make(map[SourceType]map[string]desc.Desc),
		Time:                     time.Now(),
		Repo:                     url,
	}
	if rc.DefaultSourceType == Default {
		repoDatabasePointer.DefaultSourceType = dst
	} else {
		repoDatabasePointer.DefaultSourceType = rc.DefaultSourceType
	}

	if SupportsCranBinary() {
		repoDatabasePointer.DescriptionsBySourceType[Binary] = make(map[string]desc.Desc)
	}

	repoDatabasePointer.DescriptionsBySourceType[Source] = make(map[string]desc.Desc)

	return repoDatabasePointer, repoDatabasePointer.FetchPackages(rv)
}

// Decode decodes the package database
func (repoDb *RepoDb) Decode(file string) error {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("problem opening crandb", file)
		return err
	}
	d := gob.NewDecoder(f)
	return d.Decode(&repoDb.DescriptionsBySourceType)
}

// Encode encodes the PackageDatabase
func (repoDb *RepoDb) Encode(file string) error {
	err := os.MkdirAll(filepath.Dir(file), 0777)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	e := gob.NewEncoder(f)

	// Encoding the map
	err = e.Encode(repoDb.DescriptionsBySourceType)
	if err != nil {
		return err
	}
	return nil
}

// Hash provides a hash based on the RepoDb sources
func (repoDb *RepoDb) Hash() string {
	h := md5.New()
	// want to get the unique elements in the DescriptionsBySourceType so the cache
	// will be representative of the config. Eg if set to only source
	// vs Source/Binary
	stsum := Source
	for st := range repoDb.DescriptionsBySourceType {
		stsum += st + 1
	}
	io.WriteString(h, repoDb.Repo.Name+repoDb.Repo.URL+string(stsum))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetPackagesFileURL provides the base URL for a package in a cran-like repo given the source type and version of R
func GetPackagesFileURL(r RepoURL, st SourceType, rv RVersion) string {
	if st == Source {
		return fmt.Sprintf("%s/src/contrib/PACKAGES", strings.TrimSuffix(r.URL, "/"))
		// TODO: fix so isn't hard coded to 3.5 binaries
	}
	return fmt.Sprintf("%s/bin/%s/contrib/%s/PACKAGES", strings.TrimSuffix(r.URL, "/"), cranBinaryURL(), rv.ToString())
}

// FetchPackages gets the packages for  RepoDb
// R_AVAILABLE_PACKAGES_CACHE_CONTROL_MAX_AGE controls the timing to requery the cache in R
func (repoDb *RepoDb) FetchPackages(rVersion RVersion) error {
	var err error
	pkgdbFile := repoDb.GetRepoDbCacheFilePath()

	if fi, err := os.Stat(pkgdbFile); !os.IsNotExist(err) {
		if fi.ModTime().Add(1*time.Hour).Unix() > time.Now().Unix() {
			// only read if was cached in the last hour
			return repoDb.Decode(pkgdbFile)
		}
		err := os.Remove(pkgdbFile)
		if err != nil {
			fmt.Println("error removing cache ", pkgdbFile, err)
		}
	}

	type downloadDatabase struct {
		St   SourceType
		Stdb map[string]desc.Desc
		Err  error
	}

	downloadChannel := make(chan downloadDatabase, len(repoDb.DescriptionsBySourceType))
	defer close(downloadChannel)

	for sourceType := range repoDb.DescriptionsBySourceType {
		go func(st SourceType) {
			descriptionMap := make(map[string]desc.Desc)
			pkgURL := GetPackagesFileURL(repoDb.Repo, st, rVersion)

			var body []byte

			if strings.HasPrefix(pkgURL, "http") {
				res, err := http.Get(pkgURL)
				if res.StatusCode != 200 {
					downloadChannel <- downloadDatabase{
						St: st,
						Stdb: descriptionMap,
						Err:  fmt.Errorf("failed fetching PACKAGES file from %s, with status %s", pkgURL, res.Status)}
					return
				}

				if err != nil {
					err = fmt.Errorf("problem getting packages from url %s: %s", pkgURL, err)
					downloadChannel <- downloadDatabase{St: st, Stdb: descriptionMap, Err: err}
					return
				}

				defer res.Body.Close()
				body, err = ioutil.ReadAll(res.Body)
				if err != nil {
					err = fmt.Errorf("error reading body: %s", err)
					downloadChannel <- downloadDatabase{St: st, Stdb: descriptionMap, Err: err}
					return
				}

			} else {
				pkgdir, _ := homedir.Expand(pkgURL)
				pkgdir, _ = filepath.Abs(pkgdir)
				if fi, err := os.Open(pkgdir); !os.IsNotExist(err) {
					body, err = ioutil.ReadAll(fi)
				} else {
					err = fmt.Errorf("no package file found at: %s", pkgdir)
					downloadChannel <- downloadDatabase{St: st, Stdb: descriptionMap, Err: err}
					return
				}
			}

			parsedPackagesFile := bytes.Split(body, []byte("\n\n"))
			for _, pkg := range parsedPackagesFile {
				if len(pkg) == 0 {
					// end of file might have double spaces
					// and thus will be one split, so want
					// to skip that
					//todo: trim this before the loop
					continue
				}
				reader := bytes.NewReader(pkg)
				pkgDesc, err := desc.ParseDesc(reader)
				descriptionMap[pkgDesc.Package] = pkgDesc
				if err != nil {
					fmt.Println("problem parsing package with info ", string(pkg))
					fmt.Println(err)
					downloadChannel <- downloadDatabase{St: st, Stdb: descriptionMap, Err: err}
					return
				}
			}

			downloadChannel <- downloadDatabase{St: st, Stdb: descriptionMap, Err: err}
		}(sourceType)

	}
	nerr := 0
	var lasterr error
	for i := 0; i < len(repoDb.DescriptionsBySourceType); i++ {
		result := <-downloadChannel
		if result.Err != nil {
			log.Warnf("error downloading repo %s, type: %s, with information: %s\n", repoDb.Repo.Name, result.St, result.Err)
			nerr++
			lasterr = result.Err
			// if one repo fails should return the error and not continue
			// as don't want a partial repodb as it might cause improperly pulled packages
		} else {
			repoDb.DescriptionsBySourceType[result.St] = result.Stdb
		}
	}
	// if only one source fails, this could be because it isn't present - eg if have binary/source but only source available
	if len(repoDb.DescriptionsBySourceType) > 1 && nerr == len(repoDb.DescriptionsBySourceType) {
		return lasterr
	}
	return repoDb.Encode(pkgdbFile)
}

//GetRepoDbCacheFilePath Get the filename of the file in the cache that will store this RepoDB
func (repoDb *RepoDb) GetRepoDbCacheFilePath() string {
	cdir, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("could not use user cache dir, using temp dir")
		cdir = os.TempDir()
	}
	return (filepath.Join(cdir, "pkgr", "r_packagedb_caches", repoDb.Hash()))
}

// GetPackageDbFilePath get the filepath for the cached pkgdbs
func (repoDb *RepoDb) GetPackageDbFilePath() string {
	cdir, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("could not use user cache dir, using temp dir")
		cdir = os.TempDir()
	}
	pkgdbHash := repoDb.Hash()
	return filepath.Join(cdir, "pkgr", "r_packagedb_caches", pkgdbHash)
}