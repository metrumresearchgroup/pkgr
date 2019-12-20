tags: automated

result: PASS

date_run: 12-20-2019

# Automated tests pass in IDE and in command-line:

![output](output.png)

Test output (pasted here from GoLand):
```
GOROOT=/usr/local/go #gosetup
GOPATH=/Users/johncarlos/go #gosetup
/usr/local/go/bin/go test -json ./... #gosetup
=== RUN   TestIndentationCollapse
--- PASS: TestIndentationCollapse (0.00s)
=== RUN   TestParsePackageReqs
--- PASS: TestParsePackageReqs (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/packrat	(cached)
=== RUN   TestLineScanning
--- PASS: TestLineScanning (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/rcmd/rp	(cached)
=== RUN   TestAppendGraph
time="2019-12-12T17:58:08-05:00" level=warning msg="error downloading repo CRAN, type: binary, with information: failed fetching PACKAGES file from https://cran.microsoft.com/snapshot/2018-11-11/bin/macosx/el-capitan/contrib/0.0/PACKAGES, with status 404 Not Found\n"
--- PASS: TestAppendGraph (8.73s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/gpsr	(cached)
=== RUN   TestHashing
--- PASS: TestHashing (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/rpkg	(cached)
=== RUN   TestDepTestSuite
--- PASS: TestDepTestSuite (0.00s)
=== RUN   TestDepTestSuite/TestDepToString_EqualsConstaint
    --- PASS: TestDepTestSuite/TestDepToString_EqualsConstaint (0.00s)
=== RUN   TestDepTestSuite/TestDepToString_GTConstraint
    --- PASS: TestDepTestSuite/TestDepToString_GTConstraint (0.00s)
=== RUN   TestDepTestSuite/TestDepToString_GTEConstraint
    --- PASS: TestDepTestSuite/TestDepToString_GTEConstraint (0.00s)
=== RUN   TestDepTestSuite/TestDepToString_LTConstaint
    --- PASS: TestDepTestSuite/TestDepToString_LTConstaint (0.00s)
=== RUN   TestDepTestSuite/TestDepToString_LTEConstaint
    --- PASS: TestDepTestSuite/TestDepToString_LTEConstaint (0.00s)
=== RUN   TestDepTestSuite/TestDepToString_VersionContainsDev
    --- PASS: TestDepTestSuite/TestDepToString_VersionContainsDev (0.00s)
=== RUN   TestDescParsing
--- PASS: TestDescParsing (0.00s)
=== RUN   TestVersion
--- PASS: TestVersion (0.00s)
=== RUN   TestComparisonString
--- PASS: TestComparisonString (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/desc	(cached)
=== RUN   Test_scratch
--- PASS: Test_scratch (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/scratch	(cached)
=== RUN   TestSetType
--- PASS: TestSetType (0.00s)
=== RUN   TestSetType2
time="2019-12-12T17:58:08-05:00" level=warning msg="error downloading repo CRAN_2018_11_11, type: binary, with information: failed fetching PACKAGES file from https://cran.microsoft.com/snapshot/2018-11-11/bin/macosx/el-capitan/contrib/0.0/PACKAGES, with status 404 Not Found\n"
--- PASS: TestSetType2 (7.45s)
=== RUN   TestSetRepo
time="2019-12-12T17:58:16-05:00" level=warning msg="error downloading repo CRAN_2018_11_12, type: binary, with information: failed fetching PACKAGES file from https://cran.microsoft.com/snapshot/2018-11-12/bin/macosx/el-capitan/contrib/0.0/PACKAGES, with status 404 Not Found\n"
--- PASS: TestSetRepo (4.19s)
=== RUN   TestRVersion
--- PASS: TestRVersion (0.00s)
=== RUN   TestRepoDbTestSuite
--- PASS: TestRepoDbTestSuite (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_EqualsInstalledVersionEquals
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_EqualsInstalledVersionEquals (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_GTEInstalledVersionEquals
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_GTEInstalledVersionEquals (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_GTEInstalledVersionGreaterThan
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_GTEInstalledVersionGreaterThan (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidEqualsInstalledVersionDifferent
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidEqualsInstalledVersionDifferent (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidGTEInstalledVersionLower
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidGTEInstalledVersionLower (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidGTInstalledVersionEquals
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidGTInstalledVersionEquals (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidGTInstalledVersionLower
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidGTInstalledVersionLower (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidLTEInstalledVersionHigher
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidLTEInstalledVersionHigher (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidLTInstalledVersionEqual
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidLTInstalledVersionEqual (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidLTInstalledVersionHigher
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_InvalidLTInstalledVersionHigher (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_LTEInstalledVersionEquals
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_LTEInstalledVersionEquals (0.00s)
=== RUN   TestRepoDbTestSuite/TestIsRVersionCompatible_LTInstalledVersionLessThan
    --- PASS: TestRepoDbTestSuite/TestIsRVersionCompatible_LTInstalledVersionLessThan (0.00s)
=== RUN   TestDownloadFMethod_GetMegabytes
--- PASS: TestDownloadFMethod_GetMegabytes (0.00s)
=== RUN   TestUrlHash
--- PASS: TestUrlHash (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/cran	(cached)
=== RUN   TestLibPathsEnv
--- PASS: TestLibPathsEnv (0.00s)
=== RUN   TestParseVersionData
--- PASS: TestParseVersionData (0.00s)
=== RUN   TestRMethod
--- PASS: TestRMethod (0.03s)
=== RUN   TestRunRBatch
--- PASS: TestRunRBatch (0.02s)
=== RUN   TestRunRBatch/Test_R_Version
    --- PASS: TestRunRBatch/Test_R_Version (0.01s)
=== RUN   TestConfigureArgs
--- PASS: TestConfigureArgs (0.01s)
=== RUN   TestBinaryName
--- PASS: TestBinaryName (0.00s)
=== RUN   TestBinaryExt
--- PASS: TestBinaryExt (0.00s)
=== RUN   TestInstallArgs
--- PASS: TestInstallArgs (0.00s)
=== RUN   TestUpdateDescriptionInfoByLines_RepoUpdated
--- PASS: TestUpdateDescriptionInfoByLines_RepoUpdated (0.00s)
=== RUN   TestUpdateDescriptionInfoByLines_RepoUpdated/Repository_Upated
    --- PASS: TestUpdateDescriptionInfoByLines_RepoUpdated/Repository_Upated (0.00s)
=== RUN   TestUpdateDescriptionInfoByLines_RepoTheSame
--- PASS: TestUpdateDescriptionInfoByLines_RepoTheSame (0.00s)
=== RUN   TestUpdateDescriptionInfoByLines_RepoTheSame/Pkgr_Info_Updated
    --- PASS: TestUpdateDescriptionInfoByLines_RepoTheSame/Pkgr_Info_Updated (0.00s)
=== RUN   TestUpdateDescriptionInfoByLines_RepoTheSame/Pkgr_Info_Partially_Updated
    --- PASS: TestUpdateDescriptionInfoByLines_RepoTheSame/Pkgr_Info_Partially_Updated (0.00s)
=== RUN   TestUpdateDescriptionInfoByLines_RepoTheSame/Pkgr_Info_Not_Upated
    --- PASS: TestUpdateDescriptionInfoByLines_RepoTheSame/Pkgr_Info_Not_Upated (0.00s)
=== RUN   TestUpdateDescriptionInfoByLines_RepoTheSame/Repository_Upated
    --- PASS: TestUpdateDescriptionInfoByLines_RepoTheSame/Repository_Upated (0.00s)
=== RUN   TestAppend
--- PASS: TestAppend (0.00s)
=== RUN   TestAppendNvp
--- PASS: TestAppendNvp (0.00s)
=== RUN   TestGet
--- PASS: TestGet (0.00s)
=== RUN   TestGetString
--- PASS: TestGetString (0.00s)
=== RUN   TestGetNvp
--- PASS: TestGetNvp (0.00s)
=== RUN   TestUpdate
--- PASS: TestUpdate (0.00s)
=== RUN   TestRemove
--- PASS: TestRemove (0.00s)
=== RUN   TestRemove_First
--- PASS: TestRemove_First (0.00s)
=== RUN   TestRemove_Last
--- PASS: TestRemove_Last (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/rcmd	(cached)
=== RUN   TestUtilsTestSuite
--- PASS: TestUtilsTestSuite (0.01s)
=== RUN   TestUtilsTestSuite/TestGetOutdatedPackages_DoesNotFlagOlderPackage
    --- PASS: TestUtilsTestSuite/TestGetOutdatedPackages_DoesNotFlagOlderPackage (0.00s)
=== RUN   TestUtilsTestSuite/TestGetOutdatedPackages_FindsOutdatedPackage
    --- PASS: TestUtilsTestSuite/TestGetOutdatedPackages_FindsOutdatedPackage (0.00s)
=== RUN   TestUtilsTestSuite/TestRollbackUpdatePackages_OverwritesFreshInstallation
time="2019-12-20T14:43:01-05:00" level=warning msg="did not update package, restoring last-installed version" failed to update to=2 pkg=CatsAndOranges rolling back to=1
    --- PASS: TestUtilsTestSuite/TestRollbackUpdatePackages_OverwritesFreshInstallation (0.00s)
=== RUN   TestUtilsTestSuite/TestRollbackUpdatePackages_RestoresWhenNoActiveInstallation
time="2019-12-20T14:43:01-05:00" level=warning msg="did not update package, restoring last-installed version" failed to update to=2 pkg=CatsAndOranges rolling back to=1
    --- PASS: TestUtilsTestSuite/TestRollbackUpdatePackages_RestoresWhenNoActiveInstallation (0.00s)
=== RUN   TestUtilsTestSuite/TestScanInstalledPackage_ReturnsNilWhenNoDescriptionFileFound
time="2019-12-20T14:43:01-05:00" level=warning msg="DESCRIPTION missing from installed package." file=test-library/CatsAndOranges/DESCRIPTION
    --- PASS: TestUtilsTestSuite/TestScanInstalledPackage_ReturnsNilWhenNoDescriptionFileFound (0.00s)
=== RUN   TestUtilsTestSuite/TestScanInstalledPackage_ScansReleventFieldsForOutdatedComparison
    --- PASS: TestUtilsTestSuite/TestScanInstalledPackage_ScansReleventFieldsForOutdatedComparison (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/pacman	0.029s
=== RUN   TestOperationsTestSuite
--- PASS: TestOperationsTestSuite (0.01s)
=== RUN   TestOperationsTestSuite/TestRollbackPackageEnvironment_DeletesMultiplePackages
Starting test with working directory /Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/rollback
    --- PASS: TestOperationsTestSuite/TestRollbackPackageEnvironment_DeletesMultiplePackages (0.00s)
=== RUN   TestOperationsTestSuite/TestRollbackPackageEnvironment_DeletesOnlyNewPackages
Starting test with working directory /Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/rollback
    --- PASS: TestOperationsTestSuite/TestRollbackPackageEnvironment_DeletesOnlyNewPackages (0.00s)
=== RUN   TestOperationsTestSuite/TestRollbackPackageEnvironment_HandlesEmptyListOfPackages
Starting test with working directory /Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/rollback
    --- PASS: TestOperationsTestSuite/TestRollbackPackageEnvironment_HandlesEmptyListOfPackages (0.00s)
=== RUN   TestTypesTestSuite
--- PASS: TestTypesTestSuite (0.00s)
=== RUN   TestTypesTestSuite/TestDiscernNewPackages_AllPackagesPreinstalled
    --- PASS: TestTypesTestSuite/TestDiscernNewPackages_AllPackagesPreinstalled (0.00s)
=== RUN   TestTypesTestSuite/TestDiscernNewPackages_PackagesAreCaseSensitive
    --- PASS: TestTypesTestSuite/TestDiscernNewPackages_PackagesAreCaseSensitive (0.00s)
=== RUN   TestTypesTestSuite/TestDiscernNewPackages_SomePackagesPreinstalled
    --- PASS: TestTypesTestSuite/TestDiscernNewPackages_SomePackagesPreinstalled (0.00s)
=== RUN   TestTypesTestSuite/TestDiscernNewPackages_SomePackagesPreinstalled2
    --- PASS: TestTypesTestSuite/TestDiscernNewPackages_SomePackagesPreinstalled2 (0.00s)
=== RUN   TestTypesTestSuite/TestDiscernNewPackages_ToInstallCanOutnumberPreinstalled
    --- PASS: TestTypesTestSuite/TestDiscernNewPackages_ToInstallCanOutnumberPreinstalled (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/rollback	0.031s
=== RUN   TestAddRemovePackage
--- PASS: TestAddRemovePackage (0.01s)
=== RUN   TestRemoveWhitespace
--- PASS: TestRemoveWhitespace (0.00s)
=== RUN   TestNewConfigPackrat
--- PASS: TestNewConfigPackrat (0.03s)
=== RUN   TestNewConfigNoLockfile
--- PASS: TestNewConfigNoLockfile (0.00s)
=== RUN   TestGetLibraryPath
--- PASS: TestGetLibraryPath (0.00s)
=== RUN   TestSetCustomizations
--- PASS: TestSetCustomizations (0.00s)
=== RUN   TestSetCfgCustomizations
--- PASS: TestSetCfgCustomizations (0.00s)
=== RUN   TestSetViperCustomizations
--- PASS: TestSetViperCustomizations (0.08s)
=== RUN   TestSetViperCustomizations2
--- PASS: TestSetViperCustomizations2 (0.14s)
=== RUN   TestSetPkgConfig
--- PASS: TestSetPkgConfig (0.13s)
=== RUN   TestFormat
--- PASS: TestFormat (0.00s)
=== RUN   TestViper
Version: 1
# top level packages
Packages:
  - R6
  - pillar

# any repositories, order matters
Repos:
  - CRAN: "https://cran.microsoft.com/snapshot/2019-05-01"


Library: "test-library"

--- PASS: TestViper (0.00s)
PASS
ok  	github.com/metrumresearchgroup/pkgr/configlib	0.424s
=== RUN   Test_rAddAndDelete
--- PASS: Test_rAddAndDelete (0.58s)
    add_test.go:85: testing add ...
    add_test.go:96: testing remove ...
    add_test.go:105: testing pkgr.yml for difference ...
    add_test.go:110: restoring yml file ...
=== RUN   TestCleanCacheTestSuite
--- PASS: TestCleanCacheTestSuite (0.08s)
=== RUN   TestCleanCacheTestSuite/TestCleanCache_CleansRepoFoldersWhenEmpty
    --- PASS: TestCleanCacheTestSuite/TestCleanCache_CleansRepoFoldersWhenEmpty (0.03s)
=== RUN   TestCleanCacheTestSuite/TestCleanCache_DeletesSpecificRepos
    --- PASS: TestCleanCacheTestSuite/TestCleanCache_DeletesSpecificRepos (0.02s)
=== RUN   TestCleanCacheTestSuite/TestCleanCache_DoesNotDeleteNonEmptyRepoFolders
    --- PASS: TestCleanCacheTestSuite/TestCleanCache_DoesNotDeleteNonEmptyRepoFolders (0.02s)
=== RUN   TestPlanTestSuite
--- PASS: TestPlanTestSuite (0.01s)
=== RUN   TestPlanTestSuite/TestGetPriorInstalledPackages_BasicTest
Starting test with working directory /Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/cmd
    --- PASS: TestPlanTestSuite/TestGetPriorInstalledPackages_BasicTest (0.00s)
=== RUN   TestPlanTestSuite/TestGetPriorInstalledPackages_NoPreinstalledPackages
Starting test with working directory /Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/cmd
    --- PASS: TestPlanTestSuite/TestGetPriorInstalledPackages_NoPreinstalledPackages (0.00s)
=== RUN   TestGetRestrictedWorkerCount
WARN[0000] number of workers exceeds the number of threads on machine by at least 2, this may result in degraded performance
WARN[0000] number of workers exceeds the number of threads on machine by at least 2, this may result in degraded performance
--- PASS: TestGetRestrictedWorkerCount (0.00s)
    utils_test.go:43: context: 1 CPU, no thread set
    utils_test.go:43: context: 2 CPUs, no thread set
    utils_test.go:43: context: 3 CPUs, no thread set
    utils_test.go:43: context: 4 CPUs, no thread set
    utils_test.go:43: context: 8 CPUs, no thread set
    utils_test.go:43: context: 9 CPUs, no thread set
    utils_test.go:43: context: 16 CPUs, no thread set
    utils_test.go:43: context: 4 CPUs, 2 threads set
    utils_test.go:43: context: 4 CPUs, 4 threads set
    utils_test.go:43: context: 4 CPUs, 5 threads set
    utils_test.go:43: context: 4 CPUs, 6 threads set
    utils_test.go:43: context: 4 CPUs, 7 threads set
    utils_test.go:43: context: 8 CPUs, 8 threads set
    utils_test.go:43: context: 8 CPUs, 9 threads set
    utils_test.go:43: context: 8 CPUs, 10 threads set
    utils_test.go:43: context: 8 CPUs, 16 threads set
PASS
ok  	github.com/metrumresearchgroup/pkgr/cmd	0.681s
?   	github.com/metrumresearchgroup/pkgr/cmd/pkgr	[no test files]
?   	github.com/metrumresearchgroup/pkgr/logger	[no test files]
?   	github.com/metrumresearchgroup/pkgr/pkgreqs	[no test files]
?   	github.com/metrumresearchgroup/pkgr/rpackage	[no test files]
?   	github.com/metrumresearchgroup/pkgr/testhelper	[no test files]

Process finished with exit code 0

```
