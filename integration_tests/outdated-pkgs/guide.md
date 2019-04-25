# outdated-pkgs

## Description

Environment to help test the "--update" [flag or config setting]. The environment is configured for the same packages as [simple](../simple/guide.md), plus the base-R package "Matrix"*.

** DO NOT MODIFY THE [outdated-library](outdated-library) DIRECTORY OR YOU WILL DESTROY THIS ENVIRONMENT **

## Expected behavior

You can modify the [pkgr.yml](pkgr.yml) file by setting the `Update` value to `true` or `false`.

###

\* We included the "Matrix" package here because a user had trouble using this feature with Matrix. It's mainly here as a regression test.
