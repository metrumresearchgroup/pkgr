# bad-customizations
tags: invalid-yml

## Description

Environment to help test malformed customizations and invalid pkgr.yml files.

## Expected Behavior

* `pkgr plan` should fail with an error about not being able to find repo `DoesNotExist`
* `pkgr install` should fail with an error about not being able to find repo `DoesNotExist`. Nothing should be installed.


Sample error:
`FATA[0006] error finding custom repo to set              pkg=R6 repo=DoesNotExist`
