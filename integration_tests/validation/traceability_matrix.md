# UNOFFICIAL
# Traceability Matrix: pkgr 1.0.0

## Scope
This traceability matrix links product risk, test-names, and test results to specific user stories
for the proposed software release. User stories, including requirements and test specifications,
are listed in the [Requirements Specification](req_spec.md).

| issue title | risk | test tags | pass | date run |
| ----------- | ---- | --------- | ---- | -------- |
| 179: Pkgr will attempt to minimize duplicate work, such as redownloading, re-compiling, or re-installing packages that have been used in previous actions | medium | cache-local | 1/1 | 2099-11-22 |
|   | medium | cache-system | 1/1 | 2099-11-22 |


## Notes
In the above matrix, the `test-tags` column only lists the test-tags associated
 with the issue. The `pass` column lists the proportion of integration tests with
 the associated test-tag that passed. For example, if two integration tests had
 the "cache-local" tag, and both of those tests passed, then the `pass` column of
  row one in the table above would be `2/2`.

## Questions:
* How do we incporporate automated testing into this?
