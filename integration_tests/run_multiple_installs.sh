#!/bin/bash

pkgr install --config pkgr-simple.yml
pkgr install --config pkgr-simple-suggests.yml
pkgr install --config mixed-source.

rm -rf ./testsets/multirepo/
mkdir ./testsets/multirepo
pkgr install --config pkgr-multirepo.yml
