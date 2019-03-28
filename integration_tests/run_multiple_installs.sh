#!/bin/bash

rm -rf ./testsets/simple
mkdir ./testsets/simple
pkgr install --config pkgr-simple.yml

rm -rf ./testsets/simple-suggests
mkdir ./testsets/simple-suggests
pkgr install --config pkgr-simple-suggests.yml

rm -rf ./testsets/master
mkdir ./testsets/master
pkgr install --config master-pkgr.yml

rm -rf ./testsets/internal
mkdir ./testsets/internal
pkgr install --config pkgr-internal.yml

rm -rf ./testsets/multirepo/
mkdir ./testsets/multirepo
pkgr install --config pkgr-multirepo.yml
