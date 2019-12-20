# automated-tests
tags: automated

## Description

This is simply a descriptive guide to what we expect from the automated test. This guide is used primarily for test-result scraping and should not be treated like a regular integration test.

## How to run automated tests

You can run automated tests in two ways:
Option 1) In a terminal, navigate to the home directory of this project and run `go test ./...`

Option 2) Use your IDE of choice (for example, GoLand) and set up a run configuration that runs all of the unit tests in a project.

That's it!

## Expected behavior

All tests pass. The tests listed below are included in the automated test suite. Tests that are not directly related to validation have been omitted.


* pkgr/cmd/CleanCacheSuite
    * TestCleanCache_CleansRepoFoldersWhenEmpty
    * TestCleanCache_DoesNotDeleteNonEmptyRepoFolders
    * TestCleanCache_DeletesSpecificRepos

* pkgr/pacman/TestUtilsTestSuite
    * TestGetOutdatedPackages_FindsOutdatedPackage
    * TestGetOutdatedPackages_DoesNotFlagOlderPackage
        * (This test checks that pkgr doesn’t count a package on MPN (for example) as an “updated version” when the installed package has a higher version number)

* pkgr/rollback/TestOperationsTestSuite
    * TestRollbackPackageEnvironment_DeletesOnlyNewPackages
    * TestRollbackPackageEnvironment_DeletesMultiplePackages
    * TestRollbackPackageEnvironment_HandlesEmptyListOfPackages
* pkgr/pacman/TestUtilsTestSuite
    * TestRollbackUpdatePackages_RestoresWhenNoActiveInstallation
    * TestRollbackUpdatePackages_OverwritesFreshInstallation

- cmd/add_test.go
  - Test_rAddAndDelete

* pkgr/configlib/config_test.go
  - TestNewConfigPackrat
  - TestNewConfigNoPackrat
  - TestGetLibraryPath

* pkgr/configlib/config_test.go
  * TestNewConfigPackrat (tests Packrat and RENV)
  * TestNewConfigNoLockfile
