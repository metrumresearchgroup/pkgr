tags: log-file, plan-log, install-log, log-settings, log-level

result: FAIL

date_run: 12-03-2019

# Note:
During testing, we discovered several bugs. The failing behavior is documented here:
https://github.com/metrumresearchgroup/pkgr/issues/198

# Default
## Pkgr with "Logging--all: <...>" file set logs to appropriate files
![default1](default1.png)

## Pkgr appends subsequent commands to file
![default2](default2.png)

# Install Log
<Skipped>

# Overwrite Setting
<Skipped>
