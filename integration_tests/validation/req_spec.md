# UNOFFICIAL
# Requirements Specification: pkgr 1.0.0

## Scope
The purpose of this document is to define the requirements for pkgr and document
which tests have been created to ensure that individual requirements are met.

---

## 179: Pkgr will attempt to minimize duplicate work, such as redownloading, re-compiling, or re-installing packages that have been used in previous actions

**Product risk:** medium

* Pkgr caches downloaded packages and stores them in a temp folder for faster future installations.
  * Tests:
    * cache-local, cache-system
* Pkgr maintains a cache of Package Databases (pkgdb) in a temp folder. Pkgdbs contain the information in the PACKAGES file for all repos listed in a pkgr.yml file. When pkgdbs are available in the cache, pkgr will not have to re-parse the PACKAGES file for those repos.
  * Tests:
    * cache-system
* Users should be able to specify where the cache is created. If users do not specify, the cache should use System temp folders.
  * Tests:
    * cache-local, cache-system

| Test Tag | Test Location | File | Link Test | Link Result |
| -------- | --------------| -----| --------- | ----------- |
| cache-local | pkgr/integration_tests/simple-suggests | guide.md | [Test](../simple-suggests/guide.md) | [Results](./simple-suggests/results.md) |
| cache-system | pkgr/integration_tests/simple | guide.md | [Test](../simple/guide.md) | [Results](./simple/results.md) |



---
---

## Note
If a test-tag appears in multiple tests, then that test-tag will have multiple entries in the table above, with each entry having it's own test-location. For example, if the "cache-local" tag was also found in the "repo-local" integration
test, then the table would look like this:

| Test Tag | Test Location | File | Link Test | Link Result |
| -------- | --------------| -----| --------- | ----------- |
| cache-local | pkgr/integration_tests/simple-suggests | guide.md | [Test](../simple-suggests/guide.md) | [Results](./simple-suggests/results.md) |
| cache-local | pkgr/integration_tests/repo-local | guide.md | [Test](../repo-local/guide.md) | [Results](./repo-local/results.md) |
| cache-system | pkgr/integration_tests/simple | guide.md | [Test](../simple/guide.md) | [Results](./simple/results.md) |


## To consider:
* How to fit automated tests into this table.
  - Possibly just need to point to automated test using `Test Location` and `File`, then instead of linking to a `guide.yml`, we just give the name of the automated test in `Link Test` and point to the overall automated test-results in `Link Result`
