# logging-config

## Description

Logging-config contains three distinct test environments, each designed to test a specific feature of the logging framework. The logging options tested are:
* `all`: outputs all log statements to a given file,
* `install`: outputs all log statements from `pkgr install` commands to a given file. Compatible with `all`
* `overwrite`: "resets" log files on re-run by overwriting the file contents.

## Expected behavior

### All Tests
* `pkgr plan` and `pkgr install` make install plans/install packages with the same results as the [simple](../simple/guide.md) test.

### default
* Logging statements from **all pkgr commands** are output to `default/logs/all.log`.
* Subsequent pkgr commands append their logging statements to `default/logs/all.log`

### install-log
* Logging statements from **all pkgr commands** are output to `install-log/logs/all.log`.
* Subsequent pkgr commands append their logging statements to `install-log/logs/all.log`
* *Additionally*, logging statements from `pkgr install <args>` commands are output to `install-log/logs/install.log`
* Subsequent `pkgr install <args>` calls append their logging statements to `install-log/logs/install.log`

### overwrite-setting
* Logging statements from **all pkgr commands** are output to `overwrite-settng/logs/all.log`.
* Subsequent pkgr commands cause `overwrite-setting/logs/all.log` to be **completely overwritten** by the subsequent command's logging statements.
* *Additionally*, logging statements from `pkgr install <args>` commands are output to `overwrite-setting/logs/install.log`
* Subsequent `pkgr install <args>` commands cause `overwrite-setting/logs/install.log` to be **completely overwritten** by the subsequent `pkgr install <args>` logging statements.
